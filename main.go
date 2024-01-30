package main

import (
	"fmt"
	"github.com/kube-stack/cloudctl/src/constants"
	"github.com/kube-stack/cloudctl/src/controllers"
	"github.com/kubesys/client-go/pkg/kubesys"
	"log"
	"os"
)

func main() {
	//used in cluster
	client := kubesys.NewKubernetesClientInCluster()
	configPath := os.Getenv(constants.CloudCtlConfigPath)
	if configPath == "" {
		panic("project crd_configs path is empty")
	}
	client.Init()
	//using for test
	//client := kubesys.NewKubernetesClient("https://192.168.56.103:6443", "eyJhbGciOiJSUzI1NiIsImtpZCI6IkN2SUlhdm5tb0JZN2phc1BNT3VSWVJVZUxpX1gzMjZDeUhfdTFuZFJXOFUifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudC10b2tlbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImUwYTdjOGI1LTY1MzQtNDliZS05OGI0LTg4NTg0ZWNjZGZkMSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprdWJlcm5ldGVzLWNsaWVudCJ9.tEfxNwXag4F8tN_hYof4FEDdKZc2B79CR7X4rQKcu7itk0scjezhnL2iIewrFOdkKfQeLTVmIiOjRGd5nZFAV8ntAdFzuMkEbFqfWbLPU3LaF-KRSpA7sf3EHDfMeUqFq_w8r-Ojle37FIuG_F-NsNRXzqz1WH9GPhx9LyZPcUMWmx84nLCGHkHAy6dZwsn8TljjzlsZgaq3OPFGa2CYe8y3xg43k-8kwfBjjamKQYVAOqPtm2RafxVrr8wxDEaAhf6ML8dU0UyL9OHYmsy4Vk87EFTTElzJMbX1-sf-A5PvoDmUyBMKAEV9Spr0bdoTPoeb4vq4oLpnyRjaPnqkLw")
	//client.Init()
	//configPath := "/root/go/src/cloudctl/config/crd_configs/openstack/"

	fullKinds := client.GetFullKinds()
	for _, kind := range fullKinds {
		fmt.Println(kind)
	}
	log.Println("config is in %v", configPath)
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
