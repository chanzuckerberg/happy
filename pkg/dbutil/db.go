package dbutil

import (
	"fmt"
	"os"
	"sync"

	"github.com/chanzuckerberg/happy-api/pkg/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbOptions struct {
	DSN      string
	LogLevel logger.LogLevel
}

type DBOption func(d *DB)

type DB struct {
	once         sync.Once
	dbConnection *gorm.DB
	DbOptions
}

func WithInMemorySQLDriver() DBOption {
	return func(d *DB) {
		d.DbOptions.DSN = fmt.Sprintf("file:memdb%s?mode=memory&cache=shared", uuid.NewString())
	}
}

func WithTempFileSQLDriver() DBOption {
	return func(d *DB) {
		file, err := os.CreateTemp("", "gorm.db*")
		if err != nil {
			logrus.Fatal("unable to create tmp file", err)
		}
		d.DbOptions.DSN = file.Name()
	}
}

func WithInfoLogLevel() DBOption {
	return func(d *DB) {
		d.DbOptions.LogLevel = logger.Info
	}
}

func WithErrorLogLevel() DBOption {
	return func(d *DB) {
		d.DbOptions.LogLevel = logger.Error
	}
}

// MakeDB returns a pointer because we store sync.Once inside it
// so we don't need to initialize the database many times.
func MakeDB(opts ...DBOption) *DB {
	db := &DB{
		DbOptions: DbOptions{
			DSN:      "gorm.db",
			LogLevel: logger.Silent,
		},
	}

	for _, opt := range opts {
		opt(db)
	}

	return db
}

func (d *DB) GetDB() *gorm.DB {
	d.once.Do(func() {
		var err error
		d.dbConnection, err = gorm.Open(sqlite.Open(d.DSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel(d.LogLevel)),
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
		model.AppStack{},
	}
}

func (d *DB) AutoMigrate() error {
	db := d.GetDB()
	return db.AutoMigrate(allModels()...)
}

func (d *DB) PurgeTables() error {
	db := d.GetDB()
	for _, mod := range allModels() {
		stmt := &gorm.Statement{DB: db}
		err := stmt.Parse(&mod)
		if err != nil {
			return err
		}
		tableName := stmt.Schema.Table

		db.Exec(fmt.Sprintf("DELETE FROM %s;", tableName))
	}
	return nil
}
