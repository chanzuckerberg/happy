package config

import (
	"github.com/chanzuckerberg/happy-api/pkg/dbutil"
	"github.com/chanzuckerberg/happy-api/pkg/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SetConfigValue(payload *model.AppConfigPayload) {
	db := dbutil.GetDB()

	record := model.AppConfig{AppConfigPayload: *payload}
	db.Clauses(clause.OnConflict{
		// in order to make this ON CONFLICT work we must not allow nulls for
		// stack values when dealing with an environment-level config,
		// thus the stack column defaults to emptry string and enforces NOT NULL
		Columns: []clause.Column{
			{Name: "app_name"},
			{Name: "environment"},
			{Name: "stack"},
			{Name: "key"},
		},
		UpdateAll: true,
	}).Create(&record)
}

// Returns env-level and stack-level configs for the given app and env
func GetAllAppConfigs(payload *model.AppMetadata) []*model.AppConfig {
	var records []*model.AppConfig
	criteria := dbutil.StructToMap(payload)
	delete(criteria, "stack")

	db := dbutil.GetDB()
	db.Where(criteria).Find(&records)

	return records
}

// Returns only env-level configs for the given app and env
func GetAppConfigsForEnv(payload *model.AppMetadata) []*model.AppConfigResponse {
	var records []*model.AppConfig
	criteria := dbutil.StructToMap(payload)
	criteria["stack"] = ""

	db := dbutil.GetDB()
	db.Where(criteria).Find(&records)

	return createAppConfigResponses(&records, "environment")
}

// Returns only stack-level configs for the given app, env, and stack
func GetAppConfigsForStack(payload *model.AppMetadata) []*model.AppConfigResponse {
	envConfigs := GetAppConfigsForEnv(payload)

	var records []*model.AppConfig
	criteria := dbutil.StructToMap(payload)

	db := dbutil.GetDB()
	db.Where(criteria).Find(&records)

	stackConfigs := createAppConfigResponses(&records, "stack")

	var resolvedConfigs []*model.AppConfigResponse
	for _, cfg := range envConfigs {
		stackOverrideIdx := findInStackConfigs(&stackConfigs, cfg.Key)

		if stackOverrideIdx >= 0 {
			cfg = stackConfigs[stackOverrideIdx]
			// reomve the item from the slice
			stackConfigs = append(stackConfigs[:stackOverrideIdx], stackConfigs[stackOverrideIdx+1:]...)
		}
		resolvedConfigs = append(resolvedConfigs, cfg)
	}
	resolvedConfigs = append(resolvedConfigs, stackConfigs...)

	return resolvedConfigs
}

func createAppConfigResponses(records *[]*model.AppConfig, source string) []*model.AppConfigResponse {
	var cfgs []*model.AppConfigResponse
	for _, record := range *records {
		cfgs = append(
			cfgs,
			&model.AppConfigResponse{AppConfig: *record, Source: source},
		)
	}
	return cfgs
}

func findInStackConfigs(stackConfigs *[]*model.AppConfigResponse, key string) int {
	stackOverrideIdx := -1
	for idx := range *stackConfigs {
		if (*stackConfigs)[idx].Key == key {
			stackOverrideIdx = idx
			break
		}
	}
	return stackOverrideIdx
}

func GetAppConfig(payload *model.AppConfigLookupPayload) *model.AppConfigResponse {
	var record model.AppConfig
	criteria := dbutil.StructToMap(payload)
	var result *gorm.DB

	db := dbutil.GetDB()
	if _, ok := criteria["stack"]; ok {
		result = db.Where(criteria).First(&record)
		if result.RowsAffected > 0 {
			return &model.AppConfigResponse{AppConfig: record, Source: "stack"}
		}
	}

	criteria["stack"] = ""
	result = db.Where(criteria).First(&record)
	if result.RowsAffected > 0 {
		return &model.AppConfigResponse{AppConfig: record, Source: "environment"}
	}

	return nil
}

func DeleteAppConfig(payload *model.AppConfigLookupPayload) *model.AppConfig {
	db := dbutil.GetDB()
	criteria := dbutil.StructToMap(payload)
	if _, ok := criteria["stack"]; !ok {
		criteria["stack"] = ""
	}
	record := model.AppConfig{}
	result := db.Clauses(clause.Returning{}).Where(criteria).Delete(&record)

	if result.RowsAffected == 0 {
		return nil
	}
	return &record
}
