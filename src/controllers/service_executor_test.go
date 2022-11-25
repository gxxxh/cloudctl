package controllers

import (
	"encoding/json"
	"log"
	"testing"
)

func TestInitCrdInitInfo(t *testing.T) {
	initJson := `{
	"GetComputeV2Servers": {
		"Id": ""
	}
	}`
	//paraMap := map[structpb.Value_StringValue]structpb.Value_StringValue{
	//	structpb.Value_StringValue{StringValue: "Id"}: structpb.Value_StringValue{StringValue: ""},
	//}
	//initMap := map[string]interface{}{
	//	"GetComputeV2Servers": paraMap,
	//}
	//initJson, err := structpb.NewStruct(initMap)
	//if err != nil {
	//	t.Error(err)
	//}
	var newInit json.RawMessage
	newInit.UnmarshalJSON([]byte(initJson))
	executor := &Executor{
		service: nil,
		crdConfig: &CrdConfig{
			CrdName: "OpenstackServer",
			MetaInfos: []*MetaInfo{
				&MetaInfo{
					SpecName:     "id",
					DomainName:   "id",
					InitName:     "Id",
					InitJsonPath: "GetComputeV2Servers.Id",
					IsArray:      false,
				},
			},
			InitJson:       newInit,
			DeleteJson:     nil,
			DomainJsonPath: "",
		},
		logger: NewLogger().WithName("Test"),
	}
	specInfo := `{
		"id": "111111",
	}`

	initInfo, err := executor.initCrdInitInfo([]byte(specInfo))
	if err != nil {
		t.Error(err)
	}
	log.Println(string(initInfo))
}
