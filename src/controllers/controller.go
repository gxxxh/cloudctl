package controllers

import (
	"github.com/kubesys/client-go/pkg/kubesys"
)

func DoWatch(client *kubesys.KubernetesClient, crdConfig *CrdConfig) {
	handler := NewCrdWatchHandler(crdConfig, client)
	watcher := kubesys.NewKubernetesWatcher(client, handler)
	go client.WatchResources(crdConfig.CrdName, "", watcher)
	if crdConfig.ConsumerConfig != nil {
		go handler.notificationMonitor.Run()
		go handler.HandleNotifications()
	}
}
