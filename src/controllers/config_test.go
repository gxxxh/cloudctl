package controllers

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	configPath := "D:\\GolangProjects\\src\\cloudctl\\config\\crd_configs\\test"
	err := LoadCrdConfigs(configPath)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(CrdConfigs[0])
}
