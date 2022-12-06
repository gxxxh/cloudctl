package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/kube-stack/cloudctl/src/constants"
	"github.com/kube-stack/cloudctl/src/utils"
	"github.com/kubesys/client-go/pkg/kubesys"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type CrdWatchHandler struct {
	client        *kubesys.KubernetesClient
	crdConfig     *CrdConfig
	eventRecorder *EventRecorder
	logger        *Logger
}

func NewCrdWatchHandler(crdConfig *CrdConfig, client *kubesys.KubernetesClient) *CrdWatchHandler {
	logger := NewLogger().WithName("Controller").WithName(crdConfig.GetCrdName())
	return &CrdWatchHandler{
		client:        client,
		crdConfig:     crdConfig,
		eventRecorder: NewEventRecorder(client, logger),
		logger:        logger,
	}
}

func (handler *CrdWatchHandler) reconcile(obj map[string]interface{}) error {
	crdJsonBytes, err := json.Marshal(obj)
	if err != nil {
		handler.logger.Error(err, "Marshal object to crd error")
		return err
	}
	// get Secret Info
	executor, err := handler.getExecutor(crdJsonBytes)
	if err != nil {
		return err
	}

	isNewCreate := executor.isNewCreate(crdJsonBytes)
	oldLifeCycle := gjson.GetBytes(crdJsonBytes, constants.LifeCycleJsonPath).String()
	oldDomain := gjson.GetBytes(crdJsonBytes, constants.DomainJsonPath).String()

	//设置元数据，防止用户修改为空
	if isNewCreate == false && oldLifeCycle != "" && oldLifeCycle != "{}" {
		crdJsonBytes, err = executor.SetMetaByLifecycle([]byte(oldLifeCycle), crdJsonBytes)
		if err != nil {
			return err
		}
	}
	//无需处理
	if oldLifeCycle == "" || oldLifeCycle == "{}" {
		if oldDomain == "" || oldDomain == "{}" {
			msg := fmt.Sprintf("Add Crd %v to kubernetes cluster", utils.GetCrdInfo(crdJsonBytes))
			handler.eventRecorder.Event(crdJsonBytes, EventAPIVersion, "Add Resource to k8s", msg)
			handler.logger.Info(msg)
			//update domain, remote status may change
			if err = handler.updateCrd(executor, crdJsonBytes); err != nil {
				handler.logger.Error(err, "Update crd error")
				return nil
			}
		}
		handler.logger.Info(fmt.Sprintf("No need to operate on %v", utils.GetCrdInfo(crdJsonBytes)))
		return nil
	}

	//execute
	resp, err := executor.ServiceCall([]byte(oldLifeCycle))
	if err != nil {
		handler.eventRecorder.Event(crdJsonBytes, EventTypeError, "execute lifecycle error", err.Error())
		return nil
	}
	handler.eventRecorder.Event(crdJsonBytes, EventTypeNormal, "execute lifecycle succeed", string(resp))
	handler.logger.Info("call resp", "resp: ", string(resp))

	//update lifecycle to nil
	crdJsonBytes, err = sjson.SetBytes(crdJsonBytes, constants.LifeCycleJsonPath, nil)
	if err != nil {
		handler.logger.Error(err, "Set Crd Lifecycle to nil error.")
		return err
	}

	//set meta info from resp if the lifecycle is create
	if isNewCreate {
		crdJsonBytes, err = executor.SetMetaByResp(resp, crdJsonBytes)
		if err != nil {
			handler.logger.Error(err, "Set Crd Meta from create resp error")
			return err
		}
	}

	// update domain to new info
	err = handler.updateCrd(executor, crdJsonBytes)
	if err != nil {
		return nil
	}
	return nil
}

// 调用init更新crd的domain并提交给k8s
func (handler *CrdWatchHandler) updateCrd(executor *Executor, crdJsonBytes []byte) error {
	// update domain to new info
	newCrdJson, err := executor.updateCrdJson(crdJsonBytes)
	if err != nil {
		handler.logger.Error(err, "Update crd json error")
		return err
	}
	//update
	if _, err := handler.client.UpdateResource(string(newCrdJson)); err != nil {
		handler.logger.Error(err, "Update crd to k8s error")
		handler.eventRecorder.Event(newCrdJson, EventTypeError, "update crd domain fail", err.Error())
		return err
	}
	handler.eventRecorder.Event(newCrdJson, EventTypeNormal, "execute get succeed", "update crd domain succeed")
	return nil
}

// todo go func
func (handler *CrdWatchHandler) DoAdded(obj map[string]interface{}) {
	crdJsonBytes, err := json.Marshal(obj)
	if err != nil {
		handler.logger.Error(err, "Marshal object to crd error")
		return
	}
	err = handler.reconcile(obj)
	if err != nil {
		handler.eventRecorder.Event(crdJsonBytes, EventTypeError, "internal error", err.Error())
	}
}

func (handler *CrdWatchHandler) DoModified(obj map[string]interface{}) {
	crdJsonBytes, err := json.Marshal(obj)
	if err != nil {
		handler.logger.Error(err, "Marshal object to crd error")
		return
	}
	err = handler.reconcile(obj)
	if err != nil {
		handler.eventRecorder.Event(crdJsonBytes, EventTypeError, "internal error", err.Error())
	}
}

func (handler *CrdWatchHandler) DoDeleted(obj map[string]interface{}) {
	crdJsonBytes, err := json.Marshal(obj)
	if err != nil {
		handler.logger.Error(err, "Marshal object to crd error")
		return
	}

	// get Secret Info
	executor, err := handler.getExecutor(crdJsonBytes)
	if err != nil {
		return
	}
	//check if the resource is already delete
	if executor.CheckExist(crdJsonBytes) == false {
		handler.logger.Info(fmt.Sprintf("recource %v in cloud is already delete, ", utils.GetCrdInfo(crdJsonBytes)))
		return
	}
	// delete cloud resource
	_, err = executor.CallDelete(crdJsonBytes)
	if err != nil {
		handler.logger.Error(err, fmt.Sprintf("delete cloud resource %v error", utils.GetCrdInfo(crdJsonBytes)))
		return
	}
	handler.logger.Info(fmt.Sprintf("Delete Crd %v Succeed", utils.GetCrdInfo(crdJsonBytes)))
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
