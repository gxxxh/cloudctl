package controllers

import (
	"github.com/kubesys/client-go/pkg/kubesys"
	"log"
	"testing"
)

func TestWatch(t *testing.T) {
	client := kubesys.NewKubernetesClient("https://192.168.56.103:6443", "eyJhbGciOiJSUzI1NiIsImtpZCI6Ikg5dWZoWjQ0bzNzVjRyRmJIUGZ4YnpnNkRoaXNtQkhKYWNOalFteVpFbGsifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudC10b2tlbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImJmMDRhZTU2LWRhODctNGRkNy1hNDNlLTI5ZjA4NGY5MGYxMyIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprdWJlcm5ldGVzLWNsaWVudCJ9.bWm1XblyMWL6-fOzPRsp0G1nGQHDQJoq6vzYM6IsQz47kbXNTsamN-CyoaZ3mcgNhqGXvafjeuhjo9ow5FY9uTcrvFNpNOVu-gWNcp7z_FYCySTu1LeUJBEweZrvEJCHW3sRu93XdzlumLHkMvVo_emPw6WxZtRIHQ__SAS-0hmn7IBrTFh_a9rh1v4hSpiIhZr4Sd8kSuosXVIL45Gas30qlLV1oxP5-qrMYfOMmpirn4taRG3sWCEx3Bu7TJQbrnhOuEes-XRllv7R0tu52-MO4k8uMjrZd1b1w9NybETZDr--ThYvn2jKuCm_QRWeEy_Hndq-bN6GQq4QXDIFaw")
	client.Init()
	configPath := "/root/go/src/cloudctl/config/crd_configs/openstack/openstack_server.json"

	log.Println("config is in %v", configPath)
	err := LoadCrdConfigs(configPath)
	if err != nil {
		panic(err)
	}
	for _, crdConfig := range CrdConfigs {
		DoWatch(client, crdConfig)
	}
}
