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

// 负责json的处理以及调用底层的云SDK函数
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

// 通过递归获取InitJson和DeleteJson中的函数名称
func (e *Executor) parseFuncName(jsonBytes []byte) string {
	funcInfo := gjson.ParseBytes(jsonBytes).Map()
	for funcName, _ := range funcInfo {
		return funcName
	}
	return ""
}

// 根据元数据是否为空判断是否为新创建的
func (e *Executor) IsNewCreate(crdJson []byte) bool {
	lifeCycle := gjson.GetBytes(crdJson, constants.LifeCycleJsonPath).String()
	if strings.Contains(lifeCycle, constants.CreateKeyWord) {
		return true
	}
	return false
}

func (e *Executor) IsDelete(crdJson []byte) bool {
	lifeCycle := gjson.GetBytes(crdJson, constants.LifeCycleJsonPath).String()
	deleteJsonBytes, _ := e.crdConfig.DeleteJson.MarshalJSON()
	deleteFuncName := e.parseFuncName(deleteJsonBytes)
	if strings.Contains(lifeCycle, deleteFuncName) {
		return true
	}
	return false
}

// 若id为空，根据lifecycle中的id填写到spec的id字段
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

// 根据创建返回的resp设置id，用于create
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

// 执行删除命令后，元数据需要置空
func (e *Executor) SetMetaEmpty(crdJson []byte) ([]byte, error) {
	var (
		newCrdJson []byte
		err        error
	)
	newCrdJson = make([]byte, len(crdJson), cap(crdJson))
	copy(newCrdJson, crdJson)
	for _, metaInfo := range e.crdConfig.GetMetaInfos() {
		newCrdJson, err = sjson.SetBytes(newCrdJson, constants.SpecJsonPath+metaInfo.GetSpecName(), "")
		if err != nil {
			e.logger.Error(err, "SetMetaByResp SetJson error")
			return nil, err
		}
	}
	return newCrdJson, nil
}

func (e *Executor) IsMetaEmpty(crdJson []byte) bool {
	for _, metaInfo := range e.crdConfig.GetMetaInfos() {
		if gjson.GetBytes(crdJson, constants.SpecJsonPath+metaInfo.GetSpecName()).String() == "" {
			return true
		}
	}
	return false
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
	specInfo := []byte(gjson.GetBytes(crdJsonBytes, constants.SpecName).String())
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
	specInfo := []byte(gjson.GetBytes(crdJsonBytes, constants.SpecName).String())
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
	//for image
	tmpMeta := e.crdConfig.MetaInfos[0]
	if gjson.GetBytes(resp, tmpMeta.GetInitRespJsonPath()).Exists() {
		return true
	}
	return false
}

func (e *Executor) updateCrdJson(crdJsonBytes []byte) ([]byte, error) {
	var (
		domain string
	)
	resp, err := e.CallInit(crdJsonBytes)
	if err != nil {
		return nil, err
	}
	//set domain
	if e.crdConfig.GetDomainJsonPath() != "" {
		domain = gjson.GetBytes(resp, e.crdConfig.GetDomainJsonPath()).Raw
	} else {
		domain = gjson.ParseBytes(resp).Raw
	}
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
