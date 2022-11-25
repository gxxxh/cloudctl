package constants


//spec中字段和初始化文件的json字段的对应
//spec key to init api key
var Spec2Init = map[string]map[string]string{
	AliyunEcsInstance: map[string]string{
		"regionId":"RegionId",
		"instanceId": "InstanceId",
	},
	AliyunEcsSnapshot: map[string]string{
		"regionId": "RegionId",
		"snapshotId":"SnapshotId",
	},
	OpenstackServer: map[string]string{
		"id": "Id",
	},
}
