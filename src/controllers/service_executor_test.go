package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func NewTestExecutor() *Executor {
	initJson := `{
	"GetComputeV2Servers": {
		"Id": ""
	}
	}`
	var newInit json.RawMessage
	newInit.UnmarshalJSON([]byte(initJson))
	return &Executor{
		service: nil,
		crdConfig: &CrdConfig{
			CrdName: "OpenstackServer",
			MetaInfos: []*MetaInfo{
				&MetaInfo{
					SpecName:         "id",
					DomainName:       "id",
					CloudParaName:    "id",
					InitJsonPath:     "GetComputeV2Servers.Id",
					DeleteJsonPath:   "DeleteComputeV2Servers.Id",
					InitRespJsonPath: "server.id",
					IsArray:          false,
				},
			},
			InitJson:       newInit,
			DeleteJson:     nil,
			DomainJsonPath: "",
		},
		logger: NewLogger().WithName("Test"),
	}
}
func TestExecutor_CallInit(t *testing.T) {
	executor := NewTestExecutor()
	specInfo := `{
		"id": "111111",
	}`

	initInfo, err := executor.initCrdInitInfo([]byte(specInfo))
	if err != nil {
		t.Error(err)
	}
	log.Println(string(initInfo))
}

func TestExecutor_SetMetaByLifecycle(t *testing.T) {
	lifeCycleJson := ` {
      "GetComputeV2Servers": {
        "id": "72634d76-8545-4586-95a9-34cbd57a294b"
      }
    }`
	crdJson := `{
  "apiVersion": "doslab.io/v1",
  "kind": "OpenstackServer",
  "metadata": {
    "name": "openstack-server-create"
  },
  "spec": {
    "lifeCycle": {
    },
    "id": "",
    "secretRef": {
      "namespace": "default",
      "name": "openstack-compute-secret"
    }
  }
}`
	executor := NewTestExecutor()
	newCrdJson, err := executor.SetMetaByLifecycle([]byte(lifeCycleJson), []byte(crdJson))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(newCrdJson))
}
