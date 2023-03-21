package constants

var ConfigPath = "/root/go/src/cloudctl/crd_configs"

const (
	//CRD Group & Version
	DefaultGroup               = "doslab.io"
	DefaultVersion             = "v1"
	DefaultNamespace           = "default"
	SpecName                   = "spec"
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

const CreateKeyWord = "Create"

const (
	KubernetesConfigPathEnv = "KUBERNETES_CONFIG_PATH"
	CloudCtlConfigPath      = "CLOUDCTL_CONFIG_PATH"
	NotificationConfigPath  = "NOTIFICATION_CONFIG_PATH"
)
