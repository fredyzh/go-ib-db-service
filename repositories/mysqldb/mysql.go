package mysqldb

import (
	"errors"
	"fmt"
	"ibdatabase/models"
	"log"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

const (
	Default   string = "default"
	DriveName string = "mysql"
	DSN       string = "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	MaxIdle   int    = 30
	MaxConn   int    = 30
	DBId      string = "id"
)

type MysqlDBRepo struct {
	DB        orm.Ormer
	User      string
	Password  string
	Host      string
	Port      string
	DefaultDB string
}

// Connection returns underlying connection pool.
func (m *MysqlDBRepo) Connection() error {
	orm.RegisterDriver(DriveName, orm.DRMySQL)
	dsn := fmt.Sprintf(DSN, m.User, m.Password, m.Host, m.Port, m.DefaultDB)

	orm.RegisterDataBase(Default, DriveName, dsn)
	orm.SetMaxIdleConns(Default, MaxIdle)
	orm.SetMaxOpenConns(Default, MaxConn)
	orm.DefaultTimeLoc = time.UTC
	orm.RegisterModel(new(models.Stock))
	orm.RegisterModel(new(models.DailyHistoricalStock))

	// err := orm.RunSyncdb(Default, true, true)
	// if err != nil {
	// 	log.Println(err)
	// }

	m.DB = orm.NewOrmUsingDB(Default)
	if m.DB.DBStats().OpenConnections > 0 {
		log.Println("db connected")
	} else {
		log.Println("Failed to connect to database")
		return errors.New("failed to connect to database")
	}
	return nil
}

// access db repo, laze init
func (m *MysqlDBRepo) GetDBRepo() interface{} {
	if m.DB == nil {
		err := m.Connection()
		if err != nil {
			log.Println("Failed to connect to database: ", err)
			return nil
		}
	}

	return m
}

// retrieve stock symbol list
func (m *MysqlDBRepo) GetStocks() ([]*models.Stock, error) {
	if m.DB == nil {
		m.GetDBRepo()
	}

	var stocks []*models.Stock
	_, err := m.DB.QueryTable(models.Stock{}).All(&stocks)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return stocks, nil
}

// add stock symbol
func (m *MysqlDBRepo) SaveOrUpdateStock(stk *models.Stock) error {
	if m.DB == nil {
		m.GetDBRepo()
	}

	if stk.Id == 0 {
		_, err := m.DB.Insert(stk)
		if err != nil {
			return err
		}
	} else {
		_, err := m.DB.Update(stk)

		if err != nil {
			return err
		}
	}

	return nil
}

// add stock symbols
func (m *MysqlDBRepo) InsertStocks(stks []*models.Stock) error {
	if m.DB == nil {
		m.GetDBRepo()
	}

	_, err := m.DB.InsertMulti(len(stks), stks)
	if err != nil {
		if !strings.Contains(err.Error(), "Duplicate") {
			return err
		}
	}

	return nil
}

func (m *MysqlDBRepo) DeleteStockById(id uint) error {
	if m.DB == nil {
		m.GetDBRepo()
	}

	_, err := m.DB.Delete(&models.Stock{Id: id})
	if err != nil {
		return err
	}

	return nil
}

// daily stock records
func (m *MysqlDBRepo) InsertDailyStocks(ds []*models.DailyHistoricalStock) error {
	if m.DB == nil {
		m.GetDBRepo()
	}

	_, err := m.DB.InsertMulti(len(ds), ds)
	if err != nil {
		if !strings.Contains(err.Error(), "Duplicate") {
			return err
		}
	}

	return nil
}

// daily stock records
func (m *MysqlDBRepo) SaveOrUpdateDailyStock(ds *models.DailyHistoricalStock) error {
	if m.DB == nil {
		m.GetDBRepo()
	}

	_, err := m.DB.InsertOrUpdate(ds)
	if err != nil {
		if !strings.Contains(err.Error(), "Duplicate") {
			return err
		}
	}

	return nil
}

func (m *MysqlDBRepo) FindDailyByDuration(ids []*uint, start time.Time, end time.Time) ([]*models.DailyHistoricalStock, error) {
	if m.DB == nil {
		m.GetDBRepo()
	}

	var dailys []*models.DailyHistoricalStock
	_, err := m.DB.QueryTable(models.DailyHistoricalStock{}).Filter("id__in", ids).Filter("date__gt", start).Filter("date__lt", end).All(&dailys)
	if err != nil {
		return nil, err
	}

	return dailys, nil
}
