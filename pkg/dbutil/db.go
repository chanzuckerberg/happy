package dbutil

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/chanzuckerberg/happy-api/pkg/model"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var once sync.Once
var dbConnection *gorm.DB

type DbOptions struct {
	SqlDriver gorm.Dialector
	LogLevel  int
}

func getDBOptions(env string) DbOptions {
	options := DbOptions{LogLevel: int(logger.Info)}
	switch env {
	case "development":
		options.SqlDriver = sqlite.Open("gorm.db")
	case "test":
		options.SqlDriver = sqlite.Open("file::memory:?cache=shared")
		options.LogLevel = int(logger.Silent)
	case "staging":
	case "prod":
	}
	return options
}

func init() {
	once.Do(func() {
		// TODO: make the SQL driver configurable later
		// db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
		var err error
		options := getDBOptions(os.Getenv("APP_ENV"))
		dbConnection, err = gorm.Open(options.SqlDriver, &gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel(options.LogLevel)),
		})
		if err != nil {
			logrus.Fatal("Failed to connect to the DB")
		}
	})
}

func GetDB() *gorm.DB {
	return dbConnection
}

// To get a new table added to the DB, add the model to this list
func AllModels() []interface{} {
	return []interface{}{
		model.AppConfig{},
	}
}

func AutoMigrate() error {
	db := GetDB()
	for _, mod := range AllModels() {
		db.AutoMigrate(&mod)
	}
	return nil
}

func PurgeTables() error {
	db := GetDB()
	for _, mod := range AllModels() {
		stmt := &gorm.Statement{DB: db}
		stmt.Parse(&mod)
		tableName := stmt.Schema.Table

		db.Exec(fmt.Sprintf("DELETE FROM %s;", tableName))
	}
	return nil
}

func StructToMap(payload interface{}) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(payload)
	json.Unmarshal(inrec, &inInterface)
	return inInterface
}
