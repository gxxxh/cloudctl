package controllers

import (
	"encoding/base64"
	"encoding/json"
	"github.com/kube-stack/cloudctl/src/constants"
	"github.com/kube-stack/cloudctl/src/utils"
	cloudservice "github.com/kube-stack/multicloud_service/src/service"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// todo panic
type Executor struct {
	service   *cloudservice.MultiCloudService
	crdConfig *CrdConfig
	logger    *Logger
}

func NewExecutor(crdConfig *CrdConfig, logger *Logger, secretInfo []byte) (*Executor, error) {
	e := &Executor{
		service:   nil,
		crdConfig: crdConfig,
		logger:    logger,
	}
	secretData := gjson.GetBytes(secretInfo, "data").Map()
	params := make(map[string]string)
	for key, value := range secretData {
		realValue, err := base64.StdEncoding.DecodeString(value.String())
		if err != nil {
			logger.Error(err, "Decode secret data error")
			return nil, err
		}
		params[key] = string(realValue)
	}
	mcs, err := cloudservice.NewMultiCloudService(params)
	if err != nil {
		logger.Error(err, "InitServiceBySecret err")
		return nil, err
	}
	e.service = mcs
	return e, nil
}

// 根据元数据是否为空判断是否为新创建的
func (e *Executor) isNewCreate(crdJson []byte) bool {
	for _, metaInfo := range e.crdConfig.GetMetaInfos() {
		if gjson.GetBytes(crdJson, constants.SpecJsonPath+metaInfo.GetSpecName()).String() == "" {
			return true
		}
	}
	return false
}

//// 从describe中获取元数据
//func (e *Executor) GetMetaFromDomain(domainBytes []byte) map[string]string {
//	params := make(map[string]string)
//	for _, metaInfo := range e.crdConfig.GetMetaInfos() {
//		params[metaInfo.GetInitName()] = gjson.GetBytes(domainBytes, metaInfo.GetDomainName()).String()
//	}
//	return params
//}

func (e *Executor) SetMetaByResp(resp []byte, crdJson []byte) ([]byte, error) {
	if !e.isNewCreate(crdJson) {
		return crdJson, nil
	}
	var (
		newCrdJson []byte
		err        error
	)
	newCrdJson = make([]byte, len(crdJson), cap(crdJson))
	copy(newCrdJson, crdJson)
	for _, metaInfo := range e.crdConfig.GetMetaInfos() {
		newCrdJson, err = sjson.SetBytes(newCrdJson, constants.SpecJsonPath+metaInfo.GetSpecName(), gjson.GetBytes(resp, metaInfo.GetRespJsonPath()).String())
		if err != nil {
			e.logger.Error(err, "SetJson error")
			return nil, err
		}
	}
	return newCrdJson, nil
}

func (e *Executor) initCrdInitInfo(specInfo []byte) ([]byte, error) {
	var (
		initBytes []byte
		err       error
	)
	// 传入inst.Spec，从中获取对应类型的元数据
	params := make(map[string]string)
	for _, metaInfo := range e.crdConfig.GetMetaInfos() {
		params[metaInfo.GetInitName()] = gjson.GetBytes(specInfo, metaInfo.GetSpecName()).String()
	}
	//initBytes, err = e.crdConfig.GetInitJson().MarshalJSON()
	initBytes = e.crdConfig.GetInitJson()
	for _, metaInfo := range e.crdConfig.GetMetaInfos() {

		if metaInfo.GetIsArray() {
			initBytes, err = sjson.SetBytes(initBytes, metaInfo.GetInitJsonPath(), "[\""+params[metaInfo.GetInitName()]+"\"]")
		} else {
			initBytes, err = sjson.SetBytes(initBytes, metaInfo.GetInitJsonPath(), params[metaInfo.GetInitName()])
		}
		if err != nil {
			e.logger.Error(err, "SetJson error")
			return nil, err
		}
	}
	return initBytes, err
}

func (e *Executor) ServiceCall(requestInfo []byte) ([]byte, error) {
	requestMap, err := utils.Jsonbyte2Map(requestInfo)
	if err != nil {
		e.logger.Error(err, "Marshal requestinfo to map error:")
		return nil, err
	}
	//only on request
	for APIName, APIParameters := range requestMap {
		jsonBytes, err := json.Marshal(APIParameters)
		if err != nil {
			e.logger.Error(err, "Marshal parameters to json error")
			return nil, err
		}
		resp, err := e.service.CallCloudAPI(APIName, jsonBytes)
		if err != nil {
			e.logger.Error(err, "Call Cloud API Error")
			return nil, err
		}
		return resp, err
	}
	return nil, nil
}

func (e *Executor) UpdateCrdDomain(crdJsonBytes []byte) ([]byte, error) {
	specInfo := []byte(gjson.GetBytes(crdJsonBytes, "spec").String())
	crdInitInfo, err := e.initCrdInitInfo(specInfo)
	if err != nil {
		return nil, err
	}
	resp, err := e.ServiceCall(crdInitInfo)
	if err != nil {
		return nil, err
	}
	domain := gjson.GetBytes(resp, e.crdConfig.GetDomainJsonPath()).Raw
	newCrd, err := sjson.SetRawBytes(crdJsonBytes, constants.DomainJsonPath, []byte(domain))
	if err != nil {
		e.logger.Error(err, "setting domain err")
		return nil, err
	}
	return newCrd, nil
}
