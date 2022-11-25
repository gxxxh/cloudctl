package controllers

import (
	"github.com/kubesys/client-go/pkg/kubesys"
)

func DoWatch(client *kubesys.KubernetesClient, crdConfig *CrdConfig) {
	handler := NewCrdWatchHandler(crdConfig, client)
	watcher := kubesys.NewKubernetesWatcher(client, handler)
	client.WatchResources(crdConfig.CrdName, "", watcher)
}
