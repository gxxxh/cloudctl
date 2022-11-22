package utils

import (
	"encoding/json"
	"github.com/pkg/errors"
)

func Jsonbyte2Map(data []byte) (map[string]interface{}, error) {
	var res map[string]interface{}
	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil, errors.Wrap(err, "jsonByte2Map: ")
	}
	return res, err
}
