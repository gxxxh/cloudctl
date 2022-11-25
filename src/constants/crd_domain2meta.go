package constants
// domain字段映射到spec中自定义的信息字段
// spec key to domain key
var Spec2Domain = map[string]map[string]string{
	AliyunEcsInstance: map[string]string{
		"regionId":"RegionId",
		"instanceId": "InstanceId",
	},
	AliyunEcsSnapshot: map[string]string{
		"regionId": "RegionId",
		"snapshotId":"SnapshotId",
	},
	OpenstackServer: map[string]string{
		"id": "id",
	},
}