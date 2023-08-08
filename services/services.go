package services

import (
	"ibdatabase/models"
	"time"
)

// a service layer to de-coupling backend
type Services interface {
	GetStocks() ([]*models.Stock, error)
	SaveOrUpdateStock(stk *models.Stock) error
	CreateStocks(stks []*models.Stock) error
	RemoveStock(id uint) error
	CreateDailyStocks(ds []*models.DailyHistoricalStock) error
	//retrieve certain stocks by duration
	RetrieveDailyByDuration(ids []string, start time.Time, end time.Time) ([]*models.DailyHistoricalStock, error)
	SaveOrUdpdateDailyStock(ds *models.DailyHistoricalStock) error
}
