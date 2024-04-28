package main

import (
	"fmt"
	"github.com/kube-stack/cloudctl/src/controllers"
	"github.com/kubesys/client-go/pkg/kubesys"
	"log"
)

func main() {
	//used in cluster
	//client := kubesys.NewKubernetesClientInCluster()
	//configPath := os.Getenv(constants.CloudCtlConfigPath)
	//if configPath == "" {
	//	panic("project crd_configs path is empty")
	//}
	//client.Init()

	//using for test
	client := kubesys.NewKubernetesClient("https://192.168.56.103:6443", "eyJhbGciOiJSUzI1NiIsImtpZCI6IkxKb3RDYjhoS0IyaGxzcFF2MUZrNDF5Y1pJQU40SG9KZzZuZWhsV0NNMWcifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudC10b2tlbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjZlNTg1YTBhLTY5MTMtNDkxYy1hNzJmLTc1Y2M5M2FlZjM4OSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprdWJlcm5ldGVzLWNsaWVudCJ9.kAgct_KRWLddx8sQ87rF9XXzX5Ayd9En7IVrGj-GnwGOzu4pYJksI_BQkiFDpespQFcdvvs_gHyWaA8ijnfHRBspxNTadmz87MHRMclpbb9TPh2-xKG-uwtz8ALL6GmBrOoxbNOTT7HsKgXESFY-Y64fo3TBoEJEjzPWezfpR8g9QhlmMvJIVYfyxQqc4rK7nbVsG2ItA0hSUNPLj4vPRpr1m1Fx8iYQrkQRsfUmDAzKio51Fv6R-WYw2lpk-QMgVKFUggTDyP6xel69KqrQ9d7duBi3qmBaijzwQFyb6Eh5TbW078dhReZXQ1H11wtNTd35R3_hu8PkbzLC75-WxA")
	client.Init()
	configPath := "/root/go/src/cloudctl/config/crd_configs/openstack/"

	fullKinds := client.GetFullKinds()
	for _, kind := range fullKinds {
		fmt.Println(kind)
	}
	log.Printf("config is in %v\n", configPath)
	err := controllers.LoadCrdConfigs(configPath)
	if err != nil {
		panic(err)
	}
	log.Print("load config done")
	for _, crdConfig := range controllers.CrdConfigs {
		log.Println("watch resource ", crdConfig.CrdName)
		//go controllers.DoWatch(client, crdConfig)
		controllers.DoWatch(client, crdConfig)
	}
	for {

	}
}
