package controllers

import (
	"github.com/gxxxh/stacktack-go/src/worker"
	"log"
)

type NotificationMonitor struct {
	*worker.NovaConsumer
	NotificationInfoChan chan *NotificationInfo
	Logger               *Logger
}

func NewNotificationMonitor(config *worker.ConsumerConfig, notificationInfoChan chan *NotificationInfo) *NotificationMonitor {
	if config == nil {
		return nil
	}
	nc, err := worker.NewNovaConsumer(config)
	if err != nil {
		log.Fatal(err)
	}
	return &NotificationMonitor{
		NovaConsumer:         nc,
		NotificationInfoChan: notificationInfoChan,
		Logger:               NewLogger().WithName("Notification Monitor"),
	}
}

func (nm *NotificationMonitor) Run() {
	nm.Logger.Info("Start Monitoring")
	err := nm.Consumer.StartConsume(nm.NovaConsumer.ExchangeConfig, nm.NovaConsumer.QueueConfig, nm.NovaConsumer.ConsumerTag)
	if err != nil {
		nm.Logger.Error(err, "NotificationMonitor Consumer error")
		return
	}
	defer nm.NovaConsumer.CleanUp()
	for d := range nm.Deliveries {
		d.Ack(false)
		info := ParseNotificationInfo(d.Body)
		nm.Logger.Info("Handling Event ", "Info", info)
		if info != nil {
			nm.NotificationInfoChan <- info
		}
	}
}
