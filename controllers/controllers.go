package controllers

import (
	"encoding/json"
	"ibdatabase/models"
	"ibdatabase/services"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Controllers struct {
	Services services.Service
	Port     string
}

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type StockReqDuration struct {
	Symbols []string `json:"symbols,omitempty"`
	Start   string   `json:"start"`
	End     string   `json:"end"`
}

func (ctr *Controllers) GetStocks(c *fiber.Ctx) error {
	stks, err := ctr.Services.GetStocks()
	if err != nil {
		errorJSON(c, err, 400)
		return err
	}

	return c.Status(200).JSON(stks)
}

func (ctr *Controllers) CreateStocks(c *fiber.Ctx) error {
	var stks []*models.Stock

	if err := c.BodyParser(&stks); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	err := ctr.Services.CreateStocks(stks)
	if err != nil {
		errorJSON(c, err, 400)
		return err
	}

	successJSON(c, "stocks added.")
	return nil
}

func (ctr *Controllers) CreateOrEditStock(c *fiber.Ctx) error {
	var stk *models.Stock

	if err := c.BodyParser(&stk); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	err := ctr.Services.SaveOrUpdateStock(stk)
	if err != nil {
		errorJSON(c, err, 400)
		return err
	}

	successJSON(c, "stocks added or update.")
	return nil
}

// unable handle spring boot date, manually parsing
func (ctr *Controllers) SaveHistoricalStocks(c *fiber.Ctx) error {
	var dailys []*models.DailyHistoricalStock

	var rawStrings []map[string]interface{}

	//parse to a map
	err := json.Unmarshal(c.Body(), &rawStrings)
	if err != nil {
		log.Println(err)
	}

	//decode map to daily stock
	decodeDailyStocks(rawStrings, &dailys)

	// log.Println(dailys)

	err = ctr.Services.CreateDailyStocks(dailys)
	if err != nil {
		errorJSON(c, err, 400)
		return err
	}

	successJSON(c, "daily stocks saved.")
	return nil
}

func (ctr *Controllers) SaveOrEditHistoricalStock(c *fiber.Ctx) error {
	var daily *models.DailyHistoricalStock

	if err := c.BodyParser(&daily); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	err := ctr.Services.SaveOrUdpdateDailyStock(daily)
	if err != nil {
		errorJSON(c, err, 400)
		return err
	}

	successJSON(c, "daily stock saved or update.")
	return nil
}

func (ctr *Controllers) GetDailyStocksByDurations(c *fiber.Ctx) error {
	var stockReq StockReqDuration

	if err := c.BodyParser(&stockReq); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	start, err := convertDateStr(stockReq.Start)

	if err != nil {
		errorJSON(c, err, 400)
		return err
	}

	end, err := convertDateStr(stockReq.End)

	daillys, err := ctr.Services.RetrieveDailyByDuration(stockReq.Symbols, *start, *end)

	if err != nil {
		errorJSON(c, err, 400)
		return err
	}

	return c.Status(200).JSON(daillys)
}

func convertDateStr(datestr string) (*time.Time, error) {
	y, err := strconv.Atoi(datestr[0:4])

	if err != nil {
		return nil, err
	}

	m, err := strconv.Atoi(datestr[4:6])
	if err != nil {
		return nil, err
	}

	d, err := strconv.Atoi(datestr[6:])
	if err != nil {
		return nil, err
	}

	dt := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	return &dt, nil
}

func errorJSON(c *fiber.Ctx, err error, status int) error {
	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	resps := []JSONResponse{payload}

	return c.Status(status).JSON(resps)
}

func successJSON(c *fiber.Ctx, msg string) error {
	var payload JSONResponse
	payload.Error = false
	payload.Message = msg

	resps := []JSONResponse{payload}

	return c.Status(200).JSON(resps)
}

// parse json
func decodeDailyStocks(jsonMaps []map[string]interface{}, out *[]*models.DailyHistoricalStock) {
	for _, m := range jsonMaps {
		daily := decodeDaily(m)
		*out = append(*out, daily)
	}
}

func decodeDaily(rawMap map[string]interface{}) *models.DailyHistoricalStock {
	daily := &models.DailyHistoricalStock{}
	daily.Close = rawMap["close"].(float64)
	daily.Open = rawMap["open"].(float64)
	daily.High = rawMap["high"].(float64)
	daily.Low = rawMap["low"].(float64)
	daily.StockId = uint(rawMap["stockId"].(float64))
	daily.Volume = int64(rawMap["volume"].(float64))
	daily.Count = uint(rawMap["count"].(float64))
	daily.Wap = rawMap["wap"].(float64)
	dt := rawMap["date"].([]interface{})
	y := int(dt[0].(float64))
	m := time.Month(int(dt[1].(float64)))
	d := int(dt[2].(float64))
	daily.Date = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	return daily
}
