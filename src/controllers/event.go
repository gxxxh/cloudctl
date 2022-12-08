package controllers

import (
	"fmt"
	"github.com/kubesys/client-go/pkg/kubesys"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	EventTypeNormal string = "Normal"
	EventTypeWaring string = "Warning"
	EventTypeError  string = "Error"
)

const (
	EventKind       string = "Event"
	EventAPIVersion string = "v1"
)

type EventRecorder struct {
	client *kubesys.KubernetesClient
	logger *Logger
}

func NewEventRecorder(client *kubesys.KubernetesClient, logger *Logger) *EventRecorder {
	return &EventRecorder{
		client: client,
		logger: logger,
	}
}

func (r *EventRecorder) Event(crdJsonBytes []byte, eventType, reason, msg string) {
	refJsonBytes, err := r.getReference(crdJsonBytes)
	if err != nil {
		return
	}
	eventJsonBytes, err := r.makeEvent(refJsonBytes, eventType, reason, msg)
	if err != nil {
		return
	}
	if err = r.pushEvent(eventJsonBytes); err != nil {
		return
	}
	return
}

func (r *EventRecorder) pushEvent(eventJsonBytes []byte) error {
	jsonStr := gjson.ParseBytes(eventJsonBytes).String()
	resp, err := r.client.CreateResource(jsonStr)
	if err != nil {
		r.logger.Error(err, fmt.Sprintf("Create Event error, resp is %v", string(resp)))
		return err
	}
	r.logger.Info(fmt.Sprintf("Create Event succeed, resp is %v", string(resp)))
	return nil
}

func (r *EventRecorder) makeEvent(refJsonBytes []byte, eventType, reason, msg string) ([]byte, error) {
	var (
		err            error
		eventJsonBytes []byte
	)
	eventJsonBytes = make([]byte, 0)
	cur := Now()
	curJson, err := cur.MarshalJSON()
	if err != nil {
		r.logger.Error(err, "marshal time error")
		return nil, err
	}
	setEventJsonString := func(eventJsonBytes []byte, jsonPath string, content interface{}) []byte {
		newEventJsonBytes, err := sjson.SetBytes(eventJsonBytes, jsonPath, content)
		if err != nil {
			r.logger.Error(err, fmt.Sprintf("set event %v error", jsonPath))
			return eventJsonBytes
		}
		return newEventJsonBytes
	}
	setEventJsonRawBytes := func(eventJsonBytes []byte, jsonPath string, content []byte) []byte {
		newEventJsonBytes, err := sjson.SetRawBytes(eventJsonBytes, jsonPath, content)
		if err != nil {
			r.logger.Error(err, fmt.Sprintf("set event %v error", jsonPath))
			return eventJsonBytes
		}
		return newEventJsonBytes
	}
	//typemeta

	eventJsonBytes = setEventJsonString(eventJsonBytes, "kind", EventKind)
	eventJsonBytes = setEventJsonString(eventJsonBytes, "apiVersion", EventAPIVersion)
	// metadata
	refName := gjson.GetBytes(refJsonBytes, "name").String()
	eventJsonBytes = setEventJsonString(eventJsonBytes, "metadata.name", fmt.Sprintf("%v.%v", refName, cur.UnixNano()))
	refNamespace := gjson.GetBytes(refJsonBytes, "namespace").String()
	if refNamespace == "" {
		refNamespace = "default"
	}
	eventJsonBytes = setEventJsonString(eventJsonBytes, "metadata.namespace", refNamespace)
	//other info
	eventJsonBytes = setEventJsonRawBytes(eventJsonBytes, "involvedObject", refJsonBytes)
	eventJsonBytes = setEventJsonString(eventJsonBytes, "reason", reason)
	eventJsonBytes = setEventJsonString(eventJsonBytes, "message", msg)
	eventJsonBytes = setEventJsonRawBytes(eventJsonBytes, "firstTimestamp", curJson)
	eventJsonBytes = setEventJsonRawBytes(eventJsonBytes, "lastTimestamp", curJson)
	eventJsonBytes = setEventJsonString(eventJsonBytes, "count", 1)
	eventJsonBytes = setEventJsonString(eventJsonBytes, "type", eventType)

	return eventJsonBytes, nil
}

// get event object reference
func (r *EventRecorder) getReference(crdJsonBytes []byte) ([]byte, error) {
	var (
		err          error
		refJsonBytes []byte
	)
	refJsonBytes = make([]byte, 0)
	jsonPathMap := map[string]string{
		"kind":            "kind",
		"apiversion":      "apiVersion",
		"name":            "metadata.name",
		"namespace":       "metadata.namespace",
		"uid":             "metadata.uid",
		"resourceVersion": "metadata.resourceVersion",
	}
	for refJsonPath, crdJsonPath := range jsonPathMap {
		refJsonBytes, err = sjson.SetBytes(refJsonBytes, refJsonPath, gjson.GetBytes(crdJsonBytes, crdJsonPath).String())
		if err != nil {
			r.logger.Error(err, "Init crd event error")
			return nil, err
		}
	}
	return refJsonBytes, nil
}
