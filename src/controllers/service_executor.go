package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/kube-stack/cloudctl/src/constants"
	"github.com/kube-stack/cloudctl/src/utils"
	cloudservice "github.com/kube-stack/multicloud_service/src/service"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"strings"
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
	lifeCycle := gjson.GetBytes(crdJson, constants.LifeCycleJsonPath).String()
	if strings.Contains(lifeCycle, "Create") {
		return true
	}
	return false
}

func (e *Executor) SetMetaByLifecycle(lifeCycle []byte, crdJson []byte) ([]byte, error) {
	var (
		newCrdJson []byte
		err        error
	)
	newCrdJson = make([]byte, len(crdJson), cap(crdJson))
	copy(newCrdJson, crdJson)
	lifeCycleMap := gjson.ParseBytes(lifeCycle).Map()
	if len(lifeCycleMap) == 0 {
		err := fmt.Errorf(", the lifecycle is empty")
		e.logger.Error(err, "SetMetaByLifeCycle error")
		return nil, err
	}
	for _, metaInfo := range e.crdConfig.GetMetaInfos() {
		for _, paraInfos := range lifeCycleMap {
			for paraName, paraValue := range paraInfos.Map() {
				if metaInfo.GetCloudParaName() == paraName {
					newCrdJson, err = sjson.SetBytes(newCrdJson, constants.SpecJsonPath+metaInfo.GetSpecName(), paraValue.String())
					if err != nil {
						e.logger.Error(err, "SetMetaByLifecycle SetJson error")
						return nil, err
					}
				}
			}
		}
	}
	return newCrdJson, nil
}

// 根据创建返回的resp设置元数据
func (e *Executor) SetMetaByResp(resp []byte, crdJson []byte) ([]byte, error) {
	var (
		newCrdJson []byte
		err        error
	)
	newCrdJson = make([]byte, len(crdJson), cap(crdJson))
	copy(newCrdJson, crdJson)
	for _, metaInfo := range e.crdConfig.GetMetaInfos() {
		newCrdJson, err = sjson.SetBytes(newCrdJson, constants.SpecJsonPath+metaInfo.GetSpecName(), gjson.GetBytes(resp, metaInfo.GetInitRespJsonPath()).String())
		if err != nil {
			e.logger.Error(err, "SetMetaByResp SetJson error")
			return nil, err
		}
	}
	return newCrdJson, nil
}

// fill crd init json with call parameters
func (e *Executor) initCrdDeleteJson(specInfo []byte) ([]byte, error) {
	var (
		deleteBytes []byte
		err         error
	)
	deleteBytes = e.crdConfig.GetDeleteJson()
	for _, metaInfo := range e.crdConfig.GetMetaInfos() {
		paraString := gjson.GetBytes(specInfo, metaInfo.GetSpecName()).String()
		deleteBytes, err = sjson.SetBytes(deleteBytes, metaInfo.GetDeleteJsonPath(), paraString)
		if err != nil {
			e.logger.Error(err, "SetJson error")
			return nil, err
		}
	}
	return deleteBytes, err
}

// fill crd init json with call parameters
func (e *Executor) initCrdInitInfo(specInfo []byte) ([]byte, error) {
	var (
		initBytes []byte
		err       error
	)
	//initBytes, err = e.crdConfig.GetInitJson().MarshalJSON()
	initBytes = e.crdConfig.GetInitJson()
	for _, metaInfo := range e.crdConfig.GetMetaInfos() {
		paraString := gjson.GetBytes(specInfo, metaInfo.GetSpecName()).String()
		if metaInfo.GetIsArray() {
			initBytes, err = sjson.SetBytes(initBytes, metaInfo.GetInitJsonPath(), "[\""+paraString+"\"]")
		} else {
			initBytes, err = sjson.SetBytes(initBytes, metaInfo.GetInitJsonPath(), paraString)
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

// 调用get方法
func (e *Executor) CallInit(crdJsonBytes []byte) ([]byte, error) {
	specInfo := []byte(gjson.GetBytes(crdJsonBytes, "spec").String())
	crdInitInfo, err := e.initCrdInitInfo(specInfo)
	if err != nil {
		return nil, err
	}
	resp, err := e.ServiceCall(crdInitInfo)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 调用删除方法
func (e *Executor) CallDelete(crdJsonBytes []byte) ([]byte, error) {
	specInfo := []byte(gjson.GetBytes(crdJsonBytes, "spec").String())
	crdDeleteInfo, err := e.initCrdDeleteJson(specInfo)
	if err != nil {
		return nil, err
	}
	resp, err := e.ServiceCall(crdDeleteInfo)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 检查资源是否删除
func (e *Executor) CheckExist(crdJsonBytes []byte) bool {
	resp, err := e.CallInit(crdJsonBytes)
	if err != nil {
		return false
	}
	if gjson.GetBytes(resp, e.crdConfig.GetDomainJsonPath()).Exists() {
		return true
	}
	return false
}

func (e *Executor) updateCrdJson(crdJsonBytes []byte) ([]byte, error) {
	resp, err := e.CallInit(crdJsonBytes)
	if err != nil {
		return nil, err
	}
	//set domain
	domain := gjson.GetBytes(resp, e.crdConfig.GetDomainJsonPath()).Raw
	newCrd, err := sjson.SetRawBytes(crdJsonBytes, constants.DomainJsonPath, []byte(domain))
	if err != nil {
		e.logger.Error(err, "setting domain err")
		return nil, err
	}
	//set meta info, which may be empty in other operations
	for _, metaInfo := range e.crdConfig.MetaInfos {
		metaJson := gjson.GetBytes(resp, metaInfo.GetInitRespJsonPath()).Raw
		newCrd, err = sjson.SetRawBytes(newCrd, constants.SpecJsonPath+metaInfo.GetSpecName(), []byte(metaJson))
	}

	return newCrd, nil
}
