package cmd

import (
	"context"
	"fmt"

	"github.com/chanzuckerberg/happy/api/pkg/dbutil"
	"github.com/chanzuckerberg/happy/api/pkg/ent"
	"github.com/chanzuckerberg/happy/api/pkg/ent/appconfig"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type Config interface {
	SetConfigValue(*model.AppConfigPayload) (*model.AppConfig, error)
	GetAllAppConfigs(*model.AppMetadata) ([]*model.AppConfig, error)
	GetAppConfigsForEnv(*model.AppMetadata) ([]*model.ResolvedAppConfig, error)
	GetAppConfigsForStack(*model.AppMetadata) ([]*model.ResolvedAppConfig, error)
	GetResolvedAppConfig(*model.AppConfigLookupPayload) (*model.ResolvedAppConfig, error)
	DeleteAppConfig(*model.AppConfigLookupPayload) (*model.AppConfig, error)
	CopyAppConfig(*model.CopyAppConfigPayload) (*model.AppConfig, error)
	AppConfigDiff(*model.AppConfigDiffPayload) ([]string, error)
	CopyAppConfigDiff(*model.AppConfigDiffPayload) ([]*model.AppConfig, error)
}

type dbConfig struct {
	DB *dbutil.DB
}

func MakeConfig(db *dbutil.DB) Config {
	return &dbConfig{
		DB: db,
	}
}

func MakeAppConfigFromEnt(in *ent.AppConfig) *model.AppConfig {
	var deletedAt gorm.DeletedAt
	if in.DeletedAt == nil {
		deletedAt = gorm.DeletedAt{
			Valid: false,
		}
	} else {
		deletedAt = gorm.DeletedAt{
			Time:  *in.DeletedAt,
			Valid: false,
		}
	}
	return &model.AppConfig{
		CommonDBFields: model.CommonDBFields{
			ID:        in.ID,
			CreatedAt: in.CreatedAt,
			UpdatedAt: in.UpdatedAt,
			DeletedAt: deletedAt,
		},
		AppConfigPayload: *model.NewAppConfigPayload(in.AppName, in.Environment, in.Stack, in.Key, in.Value),
	}
}

func (c *dbConfig) SetConfigValue(payload *model.AppConfigPayload) (*model.AppConfig, error) {
	db := c.DB.GetDBEnt()
	ctx := context.Background()
	tx, err := db.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("[SetConfigValue] starting a transaction: %w", err)
	}

	err = tx.AppConfig.Create().
		SetAppName(payload.AppName).
		SetEnvironment(payload.Environment).
		SetStack(payload.Stack).
		SetKey(payload.Key).
		SetValue(payload.Value).
		OnConflictColumns(appconfig.FieldAppName, appconfig.FieldEnvironment, appconfig.FieldStack, appconfig.FieldKey).
		UpdateNewValues().
		Exec(ctx)
	if err != nil {
		return nil, rollback(tx, errors.Wrap(err, "[SetConfigValue] unable to create app config"))
	}

	record, err := appEnvScopedQuery(tx.AppConfig, &payload.AppMetadata).
		Where(
			appconfig.Stack(payload.Stack),
			appconfig.Key(payload.Key),
		).Only(ctx)
	if err != nil {
		return nil, rollback(tx, errors.Wrap(err, "[SetConfigValue] unable to query app config"))
	}

	err = tx.Commit()
	return MakeAppConfigFromEnt(record), err
}

// Returns all env-level and stack-level configs for the given app and env (no overrides applied)
func (c *dbConfig) GetAllAppConfigs(payload *model.AppMetadata) ([]*model.AppConfig, error) {
	db := c.DB.GetDBEnt()
	records, err := appEnvScopedQuery(db.AppConfig, payload).All(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "[GetAllAppConfigs] unable to query app configs")
	}

	return entArrayToSharedModelArray(records), nil
}

// Returns only env-level configs for the given app and env
func (c *dbConfig) GetAppConfigsForEnv(payload *model.AppMetadata) ([]*model.ResolvedAppConfig, error) {
	db := c.DB.GetDBEnt()
	records, err := appEnvScopedQuery(db.AppConfig, payload).
		Where(appconfig.Stack("")).
		All(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "[GetAppConfigsForEnv] unable to query app configs")
	}

	results := entArrayToSharedModelArray(records)
	resolvedResults := []*model.ResolvedAppConfig{}
	for _, result := range results {
		resolvedResults = append(resolvedResults, model.NewResolvedAppConfig(result))
	}

	return resolvedResults, nil
}

// Returns resolved stack-level configs for the given app, env, and stack (with overrides applied)
func (c *dbConfig) GetAppConfigsForStack(payload *model.AppMetadata) ([]*model.ResolvedAppConfig, error) {
	// get all appconfigs for the app/env and order by key, then by stack DESC. Take the first item for each key
	db := c.DB.GetDBEnt()
	records, err := appEnvScopedQuery(db.AppConfig, payload).
		Where(appconfig.StackIn(payload.Stack, "")).
		Order(ent.Asc(appconfig.FieldKey), ent.Desc(appconfig.FieldStack)).
		All(context.Background())
	if err != nil {
		return nil, err
	}

	// we'll get at most 2 config records for each key (one for env and one for stack), so we'll use a map to dedupe
	// and select the stack record if it exists (since we order by stack DESC) and the env record otherwise
	resolvedMap := map[string]*ent.AppConfig{}
	for _, record := range records {
		if _, ok := resolvedMap[record.Key]; !ok {
			resolvedMap[record.Key] = record
		}
	}

	results := []*model.ResolvedAppConfig{}
	for _, record := range resolvedMap {
		results = append(results, model.NewResolvedAppConfig(MakeAppConfigFromEnt(record)))
	}

	return results, nil
}

func appEnvScopedQuery(client *ent.AppConfigClient, payload *model.AppMetadata) *ent.AppConfigQuery {
	return client.Query().Where(
		appconfig.AppName(payload.AppName),
		appconfig.Environment(payload.Environment),
	)
}

// rollback calls to tx.Rollback and wraps the given error
// with the rollback error if occurred.
func rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}
	return err
}

func entArrayToSharedModelArray(records []*ent.AppConfig) []*model.AppConfig {
	results := make([]*model.AppConfig, len(records))
	for idx, record := range records {
		results[idx] = MakeAppConfigFromEnt(record)
	}
	return results
}

func (c *dbConfig) GetResolvedAppConfig(payload *model.AppConfigLookupPayload) (*model.ResolvedAppConfig, error) {
	db := c.DB.GetDBEnt()
	records, err := appEnvScopedQuery(db.AppConfig, &payload.AppMetadata).
		Where(
			appconfig.Key(payload.Key),
			appconfig.StackIn(payload.Stack, ""),
		).
		Order(ent.Desc(appconfig.FieldStack)).
		All(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "[GetResolvedAppConfig] unable to query app configs")
	}

	if len(records) == 0 {
		return nil, nil
	}

	// at most 2 records are defined and since we order by stack DESC, the first record is the stack-specific one if it exists
	result := MakeAppConfigFromEnt(records[0])
	return model.NewResolvedAppConfig(result), nil
}

func (c *dbConfig) DeleteAppConfig(payload *model.AppConfigLookupPayload) (*model.AppConfig, error) {
	db := c.DB.GetDBEnt()
	ctx := context.Background()
	tx, err := db.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("starting a transaction: %w", err)
	}

	records, err := appEnvScopedQuery(tx.AppConfig, &payload.AppMetadata).Where(
		appconfig.AppName(payload.AppName),
		appconfig.Environment(payload.Environment),
		appconfig.Stack(payload.Stack),
		appconfig.Key(payload.Key),
	).All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "[DeleteAppConfig] unable to query app configs")
	}

	if len(records) == 0 {
		return nil, nil
	}

	record := records[0]
	err = tx.AppConfig.DeleteOne(record).Exec(ctx)
	if err != nil {
		return nil, rollback(tx, errors.Wrap(err, "[DeleteAppConfig] unable to delete app config"))
	}
	err = tx.Commit()
	return MakeAppConfigFromEnt(record), err
}

func (c *dbConfig) CopyAppConfig(payload *model.CopyAppConfigPayload) (*model.AppConfig, error) {
	db := c.DB.GetDBEnt()
	srcRecord, err := db.AppConfig.Query().
		Where(
			appconfig.AppName(payload.AppName),
			appconfig.Environment(payload.SrcEnvironment),
			appconfig.Stack(payload.SrcStack),
			appconfig.Key(payload.Key),
		).
		First(context.Background())
	if err != nil || srcRecord == nil {
		return nil, err
	}

	resultRecord, err := c.SetConfigValue(
		model.NewAppConfigPayload(payload.AppName, payload.DstEnvironment, payload.DstStack, payload.Key, srcRecord.Value),
	)
	if err != nil {
		return nil, err
	}

	return resultRecord, nil
}

// returns array of keys that are present in Src and not in Dst
func (c *dbConfig) AppConfigDiff(payload *model.AppConfigDiffPayload) ([]string, error) {
	srcPayload := model.NewAppMetadata(payload.AppName, payload.SrcEnvironment, payload.SrcStack)
	srcConfigs, err := c.GetAppConfigsForStack(srcPayload)
	if err != nil {
		return nil, err
	}
	srcConfigKeys := []string{}
	for _, srcConfig := range srcConfigs {
		srcConfigKeys = append(srcConfigKeys, srcConfig.Key)
	}

	dstPayload := model.NewAppMetadata(payload.AppName, payload.DstEnvironment, payload.DstStack)
	dstConfigs, err := c.GetAppConfigsForStack(dstPayload)
	if err != nil {
		return nil, err
	}
	dstConfigKeys := []string{}
	for _, dstConfig := range dstConfigs {
		dstConfigKeys = append(dstConfigKeys, dstConfig.Key)
	}

	return lo.Without(srcConfigKeys, dstConfigKeys...), nil
}

func (c *dbConfig) CopyAppConfigDiff(payload *model.AppConfigDiffPayload) ([]*model.AppConfig, error) {
	missingKeys, err := c.AppConfigDiff(payload)
	if err != nil {
		return nil, err
	}

	results := []*model.AppConfig{}
	for _, key := range missingKeys {
		copyConfigPayload := model.NewCopyAppConfigPayload(
			payload.AppName,
			payload.SrcEnvironment,
			payload.SrcStack,
			payload.DstEnvironment,
			payload.DstStack,
			key,
		)
		appConfig, err := c.CopyAppConfig(copyConfigPayload)
		if err != nil {
			return nil, err
		}
		results = append(results, appConfig)
	}

	return results, nil
}
