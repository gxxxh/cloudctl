package controllers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// meta info in spec, which will be used in cloud resource init
// todo add is in response, using
type MetaInfo struct {
	SpecName         string `protobuf:"bytes,1,opt,name=SpecName,proto3" json:"SpecName,omitempty"`
	DomainName       string `protobuf:"bytes,2,opt,name=DomainName,proto3" json:"DomainName,omitempty"`
	InitJsonPath     string `protobuf:"bytes,5,opt,name=InitJsonPath,proto3" json:"InitJsonPath,omitempty"`
	DeleteJsonPath   string `protobuf:"bytes,5,opt,name=DeleteJsonPath,proto3" json:"DeleteJsonPath,omitempty"`
	InitRespJsonPath string `protobuf:"bytes,6,opt,name=InitRespJsonPath,proto3" json:"InitRespJsonPath,omitempty"`
	IsArray          bool   `protobuf:"varint,7,opt,name=IsArray,proto3" json:"IsArray,omitempty"`
}

func (x *MetaInfo) GetSpecName() string {
	if x != nil {
		return x.SpecName
	}
	return ""
}

func (x *MetaInfo) GetDomainName() string {
	if x != nil {
		return x.DomainName
	}
	return ""
}

func (x *MetaInfo) GetInitJsonPath() string {
	if x != nil {
		return x.InitJsonPath
	}
	return ""
}

func (x *MetaInfo) GetDeleteJsonPath() string {
	if x != nil {
		return x.DeleteJsonPath
	}
	return ""
}

func (x *MetaInfo) GetInitRespJsonPath() string {
	if x != nil {
		return x.InitRespJsonPath
	}
	return ""
}

func (x *MetaInfo) GetIsArray() bool {
	if x != nil {
		return x.IsArray
	}
	return false
}

func (x *CrdConfig) GetCrdName() string {
	if x != nil {
		return x.CrdName
	}
	return ""
}

func (x *CrdConfig) GetMetaInfos() []*MetaInfo {
	if x != nil {
		return x.MetaInfos
	}
	return nil
}

func (x *CrdConfig) GetInitJson() []byte {
	if x != nil {
		res, err := x.InitJson.MarshalJSON()
		if err != nil {
			panic(err)
		}
		return res
	}
	return nil
}

func (x *CrdConfig) GetDeleteJson() []byte {
	if x != nil {
		res, err := x.DeleteJson.MarshalJSON()
		if err != nil {
			panic(err)
		}
		return res
	}
	return nil
}

func (x *CrdConfig) GetDomainJsonPath() string {
	if x != nil {
		return x.DomainJsonPath
	}
	return ""
}

type CrdConfig struct {
	CrdName        string          `protobuf:"bytes,1,opt,name=CrdName,proto3" json:"CrdName,omitempty"`
	MetaInfos      []*MetaInfo     `protobuf:"bytes,2,rep,name=MetaInfos,proto3" json:"MetaInfos,omitempty"`
	InitJson       json.RawMessage `protobuf:"bytes,3,opt,name=InitJson,proto3" json:"InitJson,omitempty"`
	DeleteJson     json.RawMessage `protobuf:"bytes,4,opt,name=DeleteJson,proto3" json:"DeleteJson,omitempty"`
	DomainJsonPath string          `protobuf:"bytes,5,opt,name=DomainJsonPath,proto3" json:"DomainJsonPath,omitempty"`
}

var CrdConfigs = make([]*CrdConfig, 0, 0)

func LoadCrdConfigs(path string) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		crdConfigJson, err := os.ReadFile(path)
		if err != nil {
			log.Panicf("read crdConfig %v error, %v\n", path, err)
		}
		crdConfig := &CrdConfig{}
		err = json.Unmarshal(crdConfigJson, crdConfig)
		if err != nil {
			log.Panicf("unmarshal crdConfig %v error, %v\n", path, err)
		}
		CrdConfigs = append(CrdConfigs, crdConfig)
		return nil
	})
	return err
}
