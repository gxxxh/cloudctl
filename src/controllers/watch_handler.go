package controllers

import (
	"encoding/json"
	"github.com/go-logr/logr"
	"github.com/kube-stack/cloudctl/src/interfaces"
	"github.com/kubesys/client-go/pkg/kubesys"
	"github.com/tidwall/gjson"
)

type CrdWatchHandler struct {
	client  kubesys.KubernetesClient
	CrdName string
	Log     logr.Logger
}

func NewCrdWatchHandler(crdName string, client *kubesys.KubernetesClient) *CrdWatchHandler {
	return &CrdWatchHandler{
		CrdName: crdName,
		Log:     logr.Logger{}.WithName("Controller").WithName(crdName),
	}
}

// todo go func
func (handler *CrdWatchHandler) DoAdded(obj map[string]interface{}) {
	crdJsonBytes, err := json.Marshal(obj)
	if err != nil {
		handler.Log.Error(err, "Marshal object to crd error")
		return
	}
	crd := &interfaces.K8sCrd{}
	err = json.Unmarshal(crdJsonBytes, crd)
	if err != nil {
		handler.Log.Error(err, "Unmarshal json into crd error")
		return
	}
	// get Secret Info
	executor, err := handler.getExecutor(crdJsonBytes)
	if err != nil {
		return
	}

	//execute
}

func (handler *CrdWatchHandler) getExecutor(crdJsonBytes []byte) (*Executor, error) {
	secret, err := handler.client.GetResource("secret",
		gjson.GetBytes(crdJsonBytes, "spec.secretRef.namespace").String(),
		gjson.GetBytes(crdJsonBytes, "spec.secretRef.name").String(),
	)
	if err != nil {
		handler.Log.Error(err, "Not found secret for the cloud")
		return nil, err
	}
	executor, err := NewExecutor(handler.CrdName, handler.Log, secret)
	if err != nil {
		handler.Log.Error(err, "Init Executor error")
		return nil, err
	}
	return executor, nil
}

func (handler *CrdWatchHandler) DoModified(obj map[string]interface{}) {

}

func (handler *CrdWatchHandler) DoDeleted(obj map[string]interface{}) {

}
