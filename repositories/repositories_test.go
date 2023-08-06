package repositories

import (
	"ibdatabase/models"
	"ibdatabase/repositories/mysqldb"
	"log"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var MySqlRepo mysqldb.MysqlDBRepo = mysqldb.MysqlDBRepo{
	Host:      "localhost",
	Port:      "3306",
	Password:  "pw1234",
	User:      "root",
	DefaultDB: "stockv3",
}

func TestDBConnection(t *testing.T) {

	Convey("Test DB Access", t, func() {
		m := MySqlRepo
		err := m.Connection()
		So(err, ShouldBeNil)
	})
}

func TestDBCreateTables(t *testing.T) {

	Convey("Test DB Access", t, func() {
		m := MySqlRepo
		err := m.Connection()
		So(err, ShouldBeNil)
	})
}

func TestAddStock(t *testing.T) {

	Convey("Test Stocks", t, func() {
		m := MySqlRepo
		err := m.Connection()
		So(err, ShouldBeNil)

		stock := models.Stock{
			StockSymbol: "TTT4",
		}

		err = m.SaveOrUpdateStock(&stock)
		So(err, ShouldBeNil)
	})
}

func TestUpdateStock(t *testing.T) {

	Convey("Test Stocks", t, func() {
		m := MySqlRepo
		err := m.Connection()
		So(err, ShouldBeNil)

		stock := models.Stock{
			Id:          1,
			StockSymbol: "TTT10",
		}

		err = m.SaveOrUpdateStock(&stock)
		So(err, ShouldBeNil)
	})
}

func TestAddStocks(t *testing.T) {

	Convey("Test Stocks", t, func() {
		m := MySqlRepo
		err := m.Connection()
		So(err, ShouldBeNil)

		stock1 := models.Stock{
			StockSymbol: "TTT2",
		}

		stock2 := models.Stock{
			StockSymbol: "TTT3",
		}

		stocks := []*models.Stock{&stock1, &stock2}

		err = m.InsertStocks(stocks)
		So(err, ShouldBeNil)
	})
}

func TestGetStocks(t *testing.T) {

	Convey("Test Stocks", t, func() {
		m := MySqlRepo
		err := m.Connection()
		So(err, ShouldBeNil)

		stocks, err := m.GetStocks()
		So(err, ShouldBeNil)
		So(len(stocks), ShouldBeGreaterThan, 0)
		log.Println(stocks[0].StockSymbol)
	})
}

func TestDeleteStocks(t *testing.T) {

	Convey("Test delete Stocks", t, func() {
		m := MySqlRepo
		err := m.Connection()
		So(err, ShouldBeNil)

		err = m.DeleteStockById(7)
		So(err, ShouldBeNil)
	})
}

func TestAddDailyStocks(t *testing.T) {

	Convey("Test Daily Stocks", t, func() {
		m := MySqlRepo
		err := m.Connection()
		So(err, ShouldBeNil)

		stock1 := models.DailyHistoricalStock{
			Open:    12.30,
			Close:   12.35,
			High:    12.80,
			Low:     12.00,
			Date:    time.Date(2023, 8, 3, 0, 0, 0, 0, time.UTC),
			StockId: 1,
			Volume:  123000,
			Count:   3200,
			Wap:     3.1,
		}

		stock2 := models.DailyHistoricalStock{
			Open:    12.30,
			Close:   12.35,
			High:    12.80,
			Low:     12.00,
			Date:    time.Date(2023, 8, 3, 0, 0, 0, 0, time.UTC),
			StockId: 2,
			Volume:  123000,
			Count:   3200,
			Wap:     3.1,
		}

		stocks := []*models.DailyHistoricalStock{&stock1, &stock2}

		err = m.InsertDailyStocks(stocks)
		So(err, ShouldBeNil)
	})
}

func TestFindDailyStocks(t *testing.T) {

	Convey("Test find daily Stocks", t, func() {
		m := MySqlRepo
		err := m.Connection()
		So(err, ShouldBeNil)

		var t1 uint = 1
		var t2 uint = 2
		ids := []*uint{&t1, &t2}
		start := time.Date(2023, 8, 2, 0, 0, 0, 0, time.UTC)
		end := time.Date(2023, 8, 6, 0, 0, 0, 0, time.UTC)
		stks, err := m.FindDailyByDuration(ids, start, end)
		So(err, ShouldBeNil)
		log.Println(*stks[0])
		log.Println(*stks[1])
		So(len(stks), ShouldBeGreaterThan, 0)
	})

}
