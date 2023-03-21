package controllers

import (
	"github.com/tidwall/gjson"
	"strings"
)

type NotificationInfo struct {
	ID    string
	Event string
}

// get instanceIDFrom notification Delivery Body
func ParseNotificationInfo(data []byte) *NotificationInfo {
	bodyJson := gjson.ParseBytes(data)
	messageJson := gjson.Parse(bodyJson.Get("oslo\\.message").Str)
	eventType := messageJson.Get("event_type")
	instanceID := ""
	if strings.Contains(eventType.String(), "instance.") {
		payLoad := messageJson.Get("payload")
		// instance UUID's seem to hide in a lot of odd places.
		if payLoad.Get("instance_id").Exists() {
			instanceID = payLoad.Get("instance_id").String()
		} else if payLoad.Get("instance_uuid").Exists() {
			instanceID = payLoad.Get("instacnce_uuid").String()
		} else if payLoad.Get("exception.kwargs.uuid").Exists() {
			instanceID = payLoad.Get("exception.kwargs.uuid").String()
		} else if payLoad.Get("instance.uuid").Exists() {
			instanceID = payLoad.Get("instance.uuid").String()
		}
	}
	if instanceID == "" {
		return nil
	}
	return &NotificationInfo{
		ID:    instanceID,
		Event: eventType.String(),
	}
}
