package constants

var ConfigPath = "/root/go/src/cloudctl/config"

const (
	//CRD Group & Version
	DefaultGroup     = "doslab.io"
	DefaultVersion   = "v1"
	DefaultNamespace = "default"

	//CRD CrdName
	OpenstackVm       = "OpenstackVirtualMachine"
	OpenstackImage    = "OpenstackImage"
	OpenstackSnapshot = "OpenstackSnapshot"
	OpenstackDisk     = "OpenstackDisk"
	AliyunEcsInstance = "AliyunEcsInstance"
	AliyunEcsSnapshot = "AliyunEcsSnapshot"
)
