package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/kube-stack/cloudctl/src/constants"
	"github.com/kube-stack/cloudctl/src/interfaces"
	"github.com/kube-stack/cloudctl/src/utils"
	"github.com/kubesys/client-go/pkg/kubesys"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type CrdWatchHandler struct {
	client    *kubesys.KubernetesClient
	crdConfig *CrdConfig
	logger    *Logger
}

func NewCrdWatchHandler(crdConfig *CrdConfig, client *kubesys.KubernetesClient) *CrdWatchHandler {

	return &CrdWatchHandler{
		client:    client,
		crdConfig: crdConfig,
		logger:    NewLogger().WithName("Controller").WithName(crdConfig.GetCrdName()),
	}
}

func (handler *CrdWatchHandler) reconcile(obj map[string]interface{}) {
	crdJsonBytes, err := json.Marshal(obj)
	if err != nil {
		handler.logger.Error(err, "Marshal object to crd error")
		return
	}
	crd := &interfaces.K8sCrd{}
	err = json.Unmarshal(crdJsonBytes, crd)
	if err != nil {
		handler.logger.Error(err, "Unmarshal json into crd error")
		return
	}
	// get Secret Info
	executor, err := handler.getExecutor(crdJsonBytes)
	if err != nil {
		return
	}

	oldLifeCycle := gjson.GetBytes(crdJsonBytes, constants.LifeCycleJsonPath).String()
	oldDomain := gjson.GetBytes(crdJsonBytes, constants.DomainJsonPath).String()

	//无需处理
	if oldLifeCycle == "" || oldLifeCycle == "{}" {
		if oldDomain == "" || oldDomain == "{}" {
			//add cloud resource to k8s
			newCrdJson, err := executor.UpdateCrdDomain(crdJsonBytes)
			if err != nil {
				return
			}
			//update
			resp, err := handler.client.UpdateResource(string(newCrdJson))
			if err != nil {
				handler.logger.Error(err, "Update crd error")
				return
			}
			//todo add event
			handler.logger.Info(fmt.Sprintf("Add Crd %v to kubernetes cluster", utils.GetCrdInfo(resp)))
		}
		handler.logger.Info(fmt.Sprintf("No need to operate on %v", utils.GetCrdInfo(crdJsonBytes)))
		return
	}
	//execute
	resp, err := executor.ServiceCall([]byte(oldLifeCycle))
	if err != nil {
		//todo add event
		return
	}

	//todo add succeed event
	handler.logger.Info("call resp", "resp: ", string(resp))
	//update lifecycle to nil
	//crdJsonBytes, err = sjson.SetBytes(crdJsonBytes, constants.LifeCycleJsonPath, "")
	crdJsonBytes, err = sjson.SetBytes(crdJsonBytes, constants.LifeCycleJsonPath, nil)
	//crdJsonBytes, err = sjson.SetRawBytes(crdJsonBytes, constants.LifeCycleJsonPath, nil)
	if err != nil {
		handler.logger.Error(err, "Set Crd Lifecycle to nil error.")
		return
	}
	crdJsonBytes, err = executor.SetMetaByResp(resp, crdJsonBytes)
	if err != nil {
		handler.logger.Error(err, "Set Crd Meta from create resp error")
		return
	}
	// update domain to new info
	//add cloud resource to k8s
	newCrdJson, err := executor.UpdateCrdDomain(crdJsonBytes)
	if err != nil {
		return
	}
	//update
	if _, err := handler.client.UpdateResource(string(newCrdJson)); err != nil {
		handler.logger.Error(err, "Update crd error")
		return
	}
}

// todo go func
func (handler *CrdWatchHandler) DoAdded(obj map[string]interface{}) {
	handler.reconcile(obj)
}

func (handler *CrdWatchHandler) DoModified(obj map[string]interface{}) {
	handler.reconcile(obj)
}

func (handler *CrdWatchHandler) DoDeleted(obj map[string]interface{}) {
	//todo call delete
	handler.reconcile(obj)
}

func (handler *CrdWatchHandler) getExecutor(crdJsonBytes []byte) (*Executor, error) {
	secret, err := handler.client.GetResource("Secret",
		gjson.GetBytes(crdJsonBytes, constants.SecretRefNamespaceJsonPath).String(),
		gjson.GetBytes(crdJsonBytes, constants.SecretRefNameJsonPath).String(),
	)
	if err != nil {
		handler.logger.Error(err, "Not found secret for the cloud")
		return nil, err
	}
	executor, err := NewExecutor(handler.crdConfig, handler.logger, secret)
	if err != nil {
		handler.logger.Error(err, "Init Executor error")
		return nil, err
	}
	return executor, nil
}
