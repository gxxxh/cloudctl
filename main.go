package main

import (
	"github.com/kube-stack/cloudctl/src/controllers"
	"github.com/kubesys/client-go/pkg/kubesys"
	"log"
)

func main() {
	//client := kubesys.NewKubernetesClientWithDefaultKubeConfig()
	//client.Init()
	client := kubesys.NewKubernetesClient("https://39.107.70.206:6443", "eyJhbGciOiJSUzI1NiIsImtpZCI6IlVfMWgwWjBsajRieERHSG5ZTFFLcGtnbktvRnRiNHl3NHNscWZ4NVlZVU0ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudC10b2tlbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6Ijc3NmY2ZTlmLTRjOGUtNDFhMC1iODQxLTY4MzQwMmI5YTRiOSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprdWJlcm5ldGVzLWNsaWVudCJ9.HkdTWtzk-Rr-Lccf0FRbMkhgs9u5QQe0YD-pk1NUQroeWqZkX-ilh-b2c5udAJfh16B3vfrpqppdny4PIiT2JEmXj_xUrMV23QzW8sKbXKw5Mox6cXMaI9_8nI_FAHFWQUKRbY6eeBEs9MWLzI-EO3IeDVdBaCoglMMSWljfvdQHWxIvgATT75eNcg7JtQYsF95O9Dpgi6_HsJKGfFWAmDJ6AitExPmP-fbjieXDbF7dycY5IFkaXlh3zbABSY-ZlcG2pG8T7dkODF3i0qaQyZmTFgAnJaINfX3IeiwFqQrXLn9HeOafUzPNPKmnsMymjyRQ1vHu_wF0wx_flXEqNg")
	client.Init()
	//configPath := "/root/go/src/cloudctl/config/crd_configs/openstack/test"
	configPath := "/root/go/src/cloudctl/config/crd_configs/test"
	err := controllers.LoadCrdConfigs(configPath)
	if err != nil {
		panic(err)
	}
	for _, crdConfig := range controllers.CrdConfigs {
		log.Println("watch resource ", crdConfig.CrdName)
		// todo go func
		controllers.DoWatch(client, crdConfig)
	}
}
