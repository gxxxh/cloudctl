package controllers

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"log"
	"testing"
)

func TestTime_MarshalJSON(t *testing.T) {
	cur := Now()
	jsonBytes, err := cur.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	info := make([]byte, 0, 0)
	info, err = sjson.SetBytes(info, "time", jsonBytes)
	if err != nil {
		t.Error(err)
	}
	info, err = sjson.SetRawBytes(info, "time1", jsonBytes)
	if err != nil {
		t.Error(err)
	}
	tmpInfo, err := sjson.Set(string(info), "time2", jsonBytes)
	if err != nil {
		t.Error(err)
	}
	log.Println(tmpInfo)
	parseTime := gjson.GetBytes(info, "time").String()
	tmp := Time{}
	err = tmp.UnmarshalJSON([]byte(parseTime))
	if err != nil {
		t.Error(err)
	}
}
