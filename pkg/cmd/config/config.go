package config

import (
	"github.com/chanzuckerberg/happy-api/pkg/dbutil"
	"github.com/chanzuckerberg/happy-api/pkg/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SetConfigValue(payload *model.AppConfigPayload) (*model.AppConfig, error) {
	db := dbutil.GetDB()
	record := model.AppConfig{AppConfigPayload: *payload}
	res := db.Clauses(clause.OnConflict{
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

	return &record, res.Error
}

// Returns env-level and stack-level configs for the given app and env
func GetAllAppConfigs(payload *model.AppMetadata) ([]*model.AppConfig, error) {
	var records []*model.AppConfig
	criteria, err := dbutil.StructToMap(payload)
	if err != nil {
		return nil, err
	}
	delete(criteria, "stack")

	db := dbutil.GetDB()
	res := db.Where(criteria).Find(&records)

	return records, res.Error
}

// Returns only env-level configs for the given app and env
func GetAppConfigsForEnv(payload *model.AppMetadata) ([]*model.AppConfigResponse, error) {
	var records []*model.AppConfig
	criteria, err := dbutil.StructToMap(payload)
	if err != nil {
		return nil, err
	}
	criteria["stack"] = ""

	db := dbutil.GetDB()
	res := db.Where(criteria).Find(&records)
	if res.Error != nil {
		return nil, res.Error
	}

	return createAppConfigResponses(&records, "environment"), nil
}

// Returns only stack-level configs for the given app, env, and stack
func GetAppConfigsForStack(payload *model.AppMetadata) ([]*model.AppConfigResponse, error) {
	envConfigs, err := GetAppConfigsForEnv(payload)
	if err != nil {
		return nil, err
	}

	var records []*model.AppConfig
	criteria, err := dbutil.StructToMap(payload)
	if err != nil {
		return nil, err
	}

	db := dbutil.GetDB()
	res := db.Where(criteria).Find(&records)
	if res.Error != nil {
		return nil, res.Error
	}

	stackConfigs := createAppConfigResponses(&records, "stack")

	resolvedConfigs := []*model.AppConfigResponse{}
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

	return resolvedConfigs, nil
}

func createAppConfigResponses(records *[]*model.AppConfig, source string) []*model.AppConfigResponse {
	cfgs := []*model.AppConfigResponse{}
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

func GetResolvedAppConfig(payload *model.AppConfigLookupPayload) (*model.AppConfigResponse, error) {
	criteria, err := dbutil.StructToMap(payload)
	if err != nil {
		return nil, err
	}

	if _, ok := criteria["stack"]; ok {
		record, err := getAppConfig(&criteria)
		if err != nil {
			return nil, err
		}
		if record != nil {
			return &model.AppConfigResponse{AppConfig: *record, Source: "stack"}, nil
		}
	}

	criteria["stack"] = ""
	record, err := getAppConfig(&criteria)
	if err != nil {
		return nil, err
	}
	if record != nil {
		return &model.AppConfigResponse{AppConfig: *record, Source: "environment"}, nil
	}

	return nil, nil
}

func getAppConfig(criteria *map[string]interface{}) (*model.AppConfig, error) {
	record := &model.AppConfig{}
	db := dbutil.GetDB()
	res := db.Where(*criteria).First(record)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}

	return record, nil
}

func DeleteAppConfig(payload *model.AppConfigLookupPayload) (*model.AppConfig, error) {
	criteria, err := dbutil.StructToMap(payload)
	if err != nil {
		return nil, err
	}
	if _, ok := criteria["stack"]; !ok {
		criteria["stack"] = ""
	}
	db := dbutil.GetDB()
	record := &model.AppConfig{}
	res := db.Clauses(clause.Returning{}).Where(criteria).Delete(record)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return record, nil
}

func CopyAppConfig(payload *model.CopyAppConfigPayload) (*model.AppConfig, error) {
	srcAppConfigPayload := model.NewAppConfigLookupPayload(payload.AppName, payload.SrcEnvironment, payload.SrcStack, payload.Key)
	// GORM won't include "stack" in the generated WHERE clause if it's unset, so we need to convert to a map then manually set the stack
	criteria, err := dbutil.StructToMap(srcAppConfigPayload)
	if err != nil {
		return nil, err
	}
	criteria["stack"] = payload.SrcStack

	record, err := getAppConfig(&criteria)
	if err != nil || record == nil {
		return nil, err
	}

	record, err = SetConfigValue(
		model.NewAppConfigPayload(payload.AppName, payload.DstEnvironment, payload.DstStack, payload.Key, record.Value),
	)
	if err != nil {
		return nil, err
	}

	return record, nil
}
