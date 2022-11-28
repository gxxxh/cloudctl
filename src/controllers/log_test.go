package controllers

import "testing"

func TestLog(t *testing.T) {
	logger := NewLogger()
	logger.Info("test Info", "key1", "value1")
	logger.Info("test Info", "key1", "value1", "key2", "value2")
}
