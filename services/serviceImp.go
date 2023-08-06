package services

import (
	"ibdatabase/models"
	"ibdatabase/repositories/mysqldb"
	"time"
)

// a service layer to de-coupling backend
type Service struct {
	MysqlDB *mysqldb.MysqlDBRepo
}

func (s *Service) GetStocks() ([]*models.Stock, error) {
	return s.MysqlDB.GetStocks()
}

func (s *Service) SaveOrUpdateStock(stk *models.Stock) error {
	return s.MysqlDB.SaveOrUpdateStock(stk)
}

func (s *Service) CreateStocks(stks []*models.Stock) error {
	return s.MysqlDB.InsertStocks(stks)
}

func (s *Service) RemoveStock(id uint) error {
	return s.MysqlDB.DeleteStockById(id)
}

func (s *Service) CreateDailyStocks(ds []*models.DailyHistoricalStock) error {
	return s.MysqlDB.InsertDailyStocks(ds)
}

func (s *Service) SaveOrUdpdateDailyStock(ds *models.DailyHistoricalStock) error {
	return s.MysqlDB.SaveOrUpdateDailyStock(ds)
}

func (s *Service) RetrieveDailyByDuration(ids []*uint, start time.Time, end time.Time) ([]*models.DailyHistoricalStock, error) {
	return s.MysqlDB.FindDailyByDuration(ids, start, end)
}
