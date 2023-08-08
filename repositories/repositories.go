package repositories

import (
	"ibdatabase/models"
	"time"
)

// an interface to de-coupling mysql implementation
type DatabaseRepo interface {
	Connection() error
	GetDBRepo() interface{}
	GetStocks() ([]*models.Stock, error)
	//save or update
	SaveOrUpdateStock(stk *models.Stock) error
	//save or update
	InsertStocks(stks []*models.Stock) error
	//delete a stock
	DeleteStockById(id uint) error

	//insert dail stock data
	InsertDailyStocks(ds []*models.DailyHistoricalStock) error
	//find certain stocks by duration
	FindDailyByDuration(ids []uint, start *time.Time, end *time.Time) ([]*models.DailyHistoricalStock, error)
	SaveOrUpdateDailyStock(ds *models.DailyHistoricalStock) error
}
