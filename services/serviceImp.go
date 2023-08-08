package services

import (
	"ibdatabase/models"
	"ibdatabase/repositories"
	"time"
)

// a service layer to de-coupling backend
type Service struct {
	GetDBRepo   repositories.DatabaseRepo
	IdSymbolMap map[string]uint
}

// get stock symbols and cache
func (s *Service) GetStocks() ([]*models.Stock, error) {
	stocks, err := s.GetDBRepo.GetStocks()
	if err != nil {
		return nil, err
	}

	if s.IdSymbolMap == nil {
		s.IdSymbolMap = map[string]uint{}
		for _, stock := range stocks {
			s.IdSymbolMap[stock.StockSymbol] = stock.Id
		}
	}

	return stocks, nil
}

func (s *Service) SaveOrUpdateStock(stk *models.Stock) error {
	return s.GetDBRepo.SaveOrUpdateStock(stk)
}

func (s *Service) CreateStocks(stks []*models.Stock) error {
	return s.GetDBRepo.InsertStocks(stks)
}

func (s *Service) RemoveStock(id uint) error {
	return s.GetDBRepo.DeleteStockById(id)
}

func (s *Service) CreateDailyStocks(ds []*models.DailyHistoricalStock) error {
	return s.GetDBRepo.InsertDailyStocks(ds)
}

func (s *Service) SaveOrUdpdateDailyStock(ds *models.DailyHistoricalStock) error {
	return s.GetDBRepo.SaveOrUpdateDailyStock(ds)
}

func (s Service) RetrieveDailyByDuration(strs []string, start time.Time, end time.Time) ([]*models.DailyHistoricalStock, error) {
	ids, _ := s.getIdsBySymbols(strs)

	return s.GetDBRepo.FindDailyByDuration(ids, &start, &end)
}

func (s *Service) getIdsBySymbols(symbs []string) ([]uint, error) {
	ids := []uint{}
	for _, symb := range symbs {
		ids = append(ids, s.IdSymbolMap[symb])
	}

	return ids, nil
}
