package dbutil

import (
	"fmt"
	"sync"

	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	once         sync.Once
	dbConnection *gorm.DB
	Config       setup.DatabaseConfiguration
}

func resolveDriver(dbCfg setup.DatabaseConfiguration) gorm.Dialector {
	switch dbCfg.Driver {
	case setup.Sqlite:
		if dbCfg.DataSourceName == ":memory:" {
			return sqlite.Open(fmt.Sprintf("file:memdb%s?mode=memory&cache=shared", uuid.NewString()))
		}
		return sqlite.Open(dbCfg.DataSourceName)
	case setup.Postgres:
		return postgres.Open(dbCfg.DataSourceName)
	default:
		logrus.Fatal("Configuration did not provide valid database driver and data_source_name")
	}
	return nil
}

func resolveLogLevel(dbCfg setup.DatabaseConfiguration) logger.LogLevel {
	switch dbCfg.LogLevel {
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "silent":
		return logger.Silent
	default:
		return logger.Info
	}
}

// MakeDB returns a pointer because we store sync.Once inside it
// so we don't need to initialize the database many times.
func MakeDB(cfg setup.DatabaseConfiguration) *DB {
	db := &DB{
		Config: cfg,
	}
	return db
}

func (d *DB) GetDB() *gorm.DB {
	d.once.Do(func() {
		var err error
		d.dbConnection, err = gorm.Open(resolveDriver(d.Config), &gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel(resolveLogLevel(d.Config))),
		})
		if err != nil {
			logrus.Fatal("Failed to connect to the DB")
		}
	})

	return d.dbConnection
}

// To get a new table added to the DB, add the model to this list
func allModels() []interface{} {
	return []interface{}{
		model.AppConfig{},
	}
}

func (d *DB) AutoMigrate() error {
	db := d.GetDB()
	return db.AutoMigrate(allModels()...)
}
