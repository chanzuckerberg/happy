package dbutil

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/chanzuckerberg/happy-api/pkg/model"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var once sync.Once
var dbConnection *gorm.DB

func init() {
	once.Do(func() {
		// TODO: make the SQL driver configurable later
		// db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
		var err error
		dbConnection, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
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
