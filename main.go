package main

import "log"

func main() {
	////used in cluster
	//client := kubesys.NewKubernetesClientInCluster()
	////using for test
	////client := kubesys.NewKubernetesClient("https://39.107.70.206:6443", "eyJhbGciOiJSUzI1NiIsImtpZCI6IlVfMWgwWjBsajRieERHSG5ZTFFLcGtnbktvRnRiNHl3NHNscWZ4NVlZVU0ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJpbi1jbHVzdGVyLW5zIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImluLWNsdXN0ZXItc2VjcmV0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImluLWNsdXN0ZXItc2EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI0Y2ExMGM5Yi1mOWRiLTQwOGItOGZiNy05M2JiMWM3NDAyZDciLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6aW4tY2x1c3Rlci1uczppbi1jbHVzdGVyLXNhIn0.12qL9-pR4esVkMi9b3FoQ16dH5NjfhfN8OL_-15XhurA3vt-zLfXG_HW7dTFBme3MsBmgcmSpu932ssTRGrd9f8sdWVSpEXbbcLfmzwdn8FxdcUAZRTP1JHCQwks8g1HP2afsuKUfPBvPUC6VUJziVe9wD0LPrnPvqVLs0nnQ5uQ5T7HAypUs-5MxmRDsqeDnGIJK32RP21I_9kDKM2UHx4sa-umsi58kx1O_pyW2vfeXp9uC46TZG2u7LWOshLqKk64Ao-aX8C74N1czr28yTzmZGc3mtHAXsotrlVs5WD-9aDVUFkP91qoR4rzWvdGc-CHQ2vCRlaa45YNgU7l7A")
	//client.Init()
	//configPath := os.Getenv(constants.CloudCtlConfigPath)
	//if configPath == "" {
	//	panic("project crd_configs path is empty")
	//}
	//err := controllers.LoadCrdConfigs(configPath)
	//if err != nil {
	//	panic(err)
	//}
	//for _, crdConfig := range controllers.CrdConfigs {
	//	log.Println("watch resource ", crdConfig.CrdName)
	//	go controllers.DoWatch(client, crdConfig)
	//}
	log.Println("test")
}
