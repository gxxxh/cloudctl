package utils

import "github.com/tidwall/gjson"

func GetKind(crdJsonBytes []byte) string {
	return gjson.GetBytes(crdJsonBytes, "kind").String()
}

func GetNamespace(crdJsonBytes []byte) string {
	return gjson.GetBytes(crdJsonBytes, "metadata.namespace").String()
}
func GetName(crdJsonBytes []byte) string {
	return gjson.GetBytes(crdJsonBytes, "metadata.name").String()
}

func GetCrdInfo(crdJsonBytes []byte) string {
	return GetNamespace(crdJsonBytes) + ":" + GetKind(crdJsonBytes) + "." + GetName(crdJsonBytes)
}
