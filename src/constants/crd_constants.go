package constants

var ConfigPath = "/root/go/src/cloudctl/config"

const (
	//CRD Group & Version
	DefaultGroup     = "doslab.io"
	DefaultVersion   = "v1"
	DefaultNamespace = "default"

	SpecJsonPath               = "spec."
	LifeCycleJsonPath          = "spec.lifeCycle"
	DomainJsonPath             = "spec.domain"
	SecretRefNamespaceJsonPath = "spec.secretRef.namespace"
	SecretRefNameJsonPath      = "spec.secretRef.name"
	//CRD CrdName
	OpenstackServer   = "OpenstackServer"
	OpenstackImage    = "OpenstackImage"
	OpenstackSnapshot = "OpenstackSnapshot"
	OpenstackDisk     = "OpenstackDisk"
	AliyunEcsInstance = "AliyunEcsInstance"
	AliyunEcsSnapshot = "AliyunEcsSnapshot"
)
