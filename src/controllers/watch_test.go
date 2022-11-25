package controllers

import (
	"github.com/kubesys/client-go/pkg/kubesys"
	"testing"
)

func TestWatch(t *testing.T) {
	client := kubesys.NewKubernetesClient("https://39.107.70.206:6443", "eyJhbGciOiJSUzI1NiIsImtpZCI6InJDT1BvcEZEejdqTllGWlZpbFZPNFludFdIeHFXeG5yNG94TjNILTlmazAifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudC10b2tlbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjFhY2JmNWYyLWFlMTEtNGI5Ni04M2EwLTlhODc3NzIxNTAzZiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprdWJlcm5ldGVzLWNsaWVudCJ9.rRzuTuPegVloxcydZ2jfbNqB8UraYi3cFJrLc8jqNnU9Scc9oQ9tX9dw66iPqw8rbbi0GnLSv79Ykba9TecoZNlKnsi5JkLJmPYWW3cACBtMu5VGsCXgCT0u9-Gik5fCBB6kywyBgKGULzr1isAuVHtbrHJJHuSTsDWONwskjMsLsrntlQTBrsgzLo9sRDX-76ofxPvPYSF_D8BYcsJg_kJln8wqdngmSURgsdOQE5WgWqIyxge3P2G3M3IPHp8XWv-AQLcKFZ4bYZB_CqJ1Q9Mttv0Ugun0aKGy2cCSmnrgbfLerwL3nO5MsnZOwTo_UArJKsrduu6x5S6X9WKsaw")
	client.Init()
	configPath := "D:\\GolangProjects\\src\\cloudctl\\config\\crd_configs\\openstack"
	err := LoadCrdConfigs(configPath)
	if err != nil {
		t.Error(err)
	}
	for _, crdConfig := range CrdConfigs {
		// todo go func
		DoWatch(client, crdConfig)
	}
}
