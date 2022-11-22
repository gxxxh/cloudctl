package constants

var CrdInitPath = ConfigPath + "/init"
var CrdInitJsonFilePath = map[string]string{
	AliyunEcsInstance: CrdInitPath + "ecs/descirbe_snapshot.json",
	AliyunEcsSnapshot: CrdInitPath + "ecs/describe_snapshot.json",
}
