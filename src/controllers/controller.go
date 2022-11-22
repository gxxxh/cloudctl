package controllers

import (
	"github.com/kubesys/client-go/pkg/kubesys"
)

func DoWatch(client *kubesys.KubernetesClient, crdName string) {
	handler := NewCrdWatchHandler(crdName, client)
	watcher := kubesys.NewKubernetesWatcher(client, handler)
	client.WatchResources(crdName, "", watcher)
}
