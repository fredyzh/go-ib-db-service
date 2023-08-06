package main

import (
	"encoding/json"
	"log"
	"testing"
)

func TestJsonUnMashallDate(t *testing.T) {
	// var daily *models.DailyHistoricalStock

	var rawStrings map[string]string
	jsonstr := "{'stockId': 7,'date': [2000,11,30],'open': 0.595,'high': 0.605,'low': 0.575,'close': 0.59,'wap': 0.5892857142857142,'volume': 4048184,'count': 1}"
	jsonByte := []byte(jsonstr)
	// log.Println([]byte(jsonByte))
	// daily = &models.DailyHistoricalStock{}
	err := json.Unmarshal(jsonByte, &rawStrings)

	if err != nil {
		log.Println(err)
	}
}
