package models

import "time"

type Stock struct {
	Id          uint   `json:"id"`
	StockSymbol string `json:"stockSymbol" orm:"unique"`
}

type DailyHistoricalStock struct {
	Id      uint      `json:"id"`
	Close   float64   `json:"close"`
	Open    float64   `json:"open"`
	High    float64   `json:"high"`
	Low     float64   `json:"low"`
	Date    time.Time `json:"date"`
	StockId uint      `json:"stockId"`
	Count   uint      `json:"count"`
	Volume  int64     `json:"volume"`
	Wap     float64   `json:"wap"`
}

// type DailyCompatible interface {
// 	uint | int | int64 | float64
// }
