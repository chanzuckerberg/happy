package config_test

import (
	"fmt"
	"testing"

	"github.com/chanzuckerberg/happy-api/pkg/cmd/config"
	"github.com/chanzuckerberg/happy-api/pkg/dbutil"
	"github.com/chanzuckerberg/happy-api/pkg/model"
	"github.com/stretchr/testify/require"
)

func purgeTables(r *require.Assertions) {
	err := dbutil.PurgeTables()
	r.NoError(err)
}

func getAppConfigPayload(appName string, env string, stack string, key string, value string) *model.AppConfigPayload {
	return &model.AppConfigPayload{
		AppMetadata: model.AppMetadata{
			AppName:     appName,
			Environment: env,
			Stack:       stack,
		},
		ConfigValue: model.ConfigValue{
			Value: value,
			ConfigKey: model.ConfigKey{
				Key: key,
			},
		},
	}
}

func TestSetConfigValueSucceed(t *testing.T) {
	testData := []struct {
		input    []*model.AppConfigPayload
		expected []*model.AppConfigPayload
	}{
		{
			input: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			},
			expected: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// same key, different value
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val2"),
			},
			expected: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val2"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// different key, different value
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
			},
			expected: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// same key BUT with a stack name, different value
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val2"),
			},
			expected: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val2"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// same key and stack name, different value
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val2"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val3"),
			},
			expected: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val3"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// same key, different stack name, different value
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val2"),
				getAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val3"),
			},
			expected: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val2"),
				getAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val3"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// different env, different value, everything else the same
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
			},
			expected: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// different app, different value, everything else the same
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("anotherapp", "rdev", "", "KEY1", "val2"),
			},
			expected: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("anotherapp", "rdev", "", "KEY1", "val2"),
			},
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			err := dbutil.AutoMigrate()
			r.NoError(err)
			defer purgeTables(r)

			for _, input := range testCase.input {
				config.SetConfigValue(input)
			}

			db := dbutil.GetDB()
			configs := []*model.AppConfig{}
			db.Find(&configs)

			results := []*model.AppConfigPayload{}
			for _, config := range configs {
				results = append(results, &config.AppConfigPayload)
			}

			r.EqualValues(results, testCase.expected)
		})
	}
}

func getAppConfigLookupPayload(app string, env string, stack string, key string) *model.AppConfigLookupPayload {
	return &model.AppConfigLookupPayload{
		AppMetadata: model.AppMetadata{
			AppName:     app,
			Environment: env,
			Stack:       stack,
		},
		ConfigKey: model.ConfigKey{
			Key: key,
		},
	}
}

func TestGetAppConfigSucceed(t *testing.T) {
	testData := []struct {
		input          []*model.AppConfigPayload
		lookup         *model.AppConfigLookupPayload
		expected       *model.AppConfigPayload
		expectedSource string
	}{
		{
			// stack and env configs exist, no stack specified -> return env config
			input: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY2", "val-foo2"),
				getAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:         getAppConfigLookupPayload("testapp", "rdev", "", "KEY2"),
			expected:       getAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
			expectedSource: "environment",
		},
		{
			// stack and env configs exist, stack specified -> return stack config
			input: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY2", "val-foo2"),
				getAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:         getAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expected:       getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
			expectedSource: "stack",
		},
		{
			// only env configs exist, stack specified -> return env config
			input: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:         getAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expected:       getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			expectedSource: "environment",
		},
		{
			// only stack configs exist, stack specified -> return stack config
			input: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:         getAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expected:       getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			expectedSource: "stack",
		},
		{
			// only stack configs exist, stack not specified -> return null
			input: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:   getAppConfigLookupPayload("testapp", "rdev", "", "KEY1"),
			expected: nil,
		},
		{
			// no configs, stack not specified -> return null
			input:    []*model.AppConfigPayload{},
			lookup:   getAppConfigLookupPayload("testapp", "rdev", "", "KEY1"),
			expected: nil,
		},
		{
			// no configs, stack specified -> return null
			input:    []*model.AppConfigPayload{},
			lookup:   getAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expected: nil,
		},
		{
			// stack config exists, different stack specified -> return null
			input: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:   getAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expected: nil,
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			err := dbutil.AutoMigrate()
			r.NoError(err)
			defer purgeTables(r)

			for _, input := range testCase.input {
				config.SetConfigValue(input)
			}

			result := config.GetAppConfig(testCase.lookup)

			if testCase.expected == nil {
				r.Nil(result)
			} else {
				r.EqualValues(result.AppConfigPayload, *testCase.expected)
				r.EqualValues(result.Source, testCase.expectedSource)
			}
		})
	}
}

func TestDeleteAppConfigSucceed(t *testing.T) {
	testData := []struct {
		input            []*model.AppConfigPayload
		deleteCriteria   *model.AppConfigLookupPayload
		expectedResult   *model.AppConfigPayload
		remainingConfigs []*model.AppConfigPayload
	}{
		{
			// no configs exist, no stack specified -> delete nothing
			input:            []*model.AppConfigPayload{},
			deleteCriteria:   getAppConfigLookupPayload("testapp", "rdev", "", "KEY1"),
			expectedResult:   nil,
			remainingConfigs: []*model.AppConfigPayload{},
		},
		{
			// no configs exist, stack specified -> delete nothing
			input:            []*model.AppConfigPayload{},
			deleteCriteria:   getAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expectedResult:   nil,
			remainingConfigs: []*model.AppConfigPayload{},
		},
		{
			// configs exist, no stack specified -> deletes env-level config
			input: []*model.AppConfigPayload{
				// order here is important, stack config needs to be created first to ensure
				// that this test fails when empty stack criteria is not added in the DeleteAppConfig
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			},
			deleteCriteria: getAppConfigLookupPayload("testapp", "rdev", "", "KEY1"),
			expectedResult: getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			remainingConfigs: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
			},
		},
		{
			// configs exist, stack specified -> deletes stack-level config
			input: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				getAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			deleteCriteria: getAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expectedResult: getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
			remainingConfigs: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			err := dbutil.AutoMigrate()
			r.NoError(err)
			defer purgeTables(r)

			for _, input := range testCase.input {
				config.SetConfigValue(input)
			}

			result := config.DeleteAppConfig(testCase.deleteCriteria)

			if testCase.expectedResult == nil {
				r.Nil(result)
			} else {
				r.EqualValues(result.AppConfigPayload, *testCase.expectedResult)
			}

			db := dbutil.GetDB()
			configs := []*model.AppConfig{}
			db.Find(&configs)

			results := []*model.AppConfigPayload{}
			for _, config := range configs {
				results = append(results, &config.AppConfigPayload)
			}

			r.EqualValues(results, testCase.remainingConfigs)
		})
	}
}

func getAppMetadata(app string, env string, stack string) *model.AppMetadata {
	return &model.AppMetadata{
		AppName:     app,
		Environment: env,
		Stack:       stack,
	}
}

func TestGetAllAppConfigsSucceed(t *testing.T) {
	testData := []struct {
		input       []*model.AppConfigPayload
		expected    []*model.AppConfigPayload
		appMetadata *model.AppMetadata
	}{
		{
			// no configs exist -> return empty array
			input:       []*model.AppConfigPayload{},
			expected:    []*model.AppConfigPayload{},
			appMetadata: getAppMetadata("testapp", "rdev", ""),
		},
		{
			// configs exist, stack specified -> ignores stack and returns all configs for app and env
			input: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				getAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			expected: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				getAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			appMetadata: getAppMetadata("testapp", "rdev", "foo"),
		},
		{
			// configs exist, different env specified -> returns only configs for specified env
			input: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				getAppConfigPayload("testapp", "staging", "", "KEY2", "val2"),
				getAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			expected: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
				getAppConfigPayload("testapp", "staging", "", "KEY2", "val2"),
			},
			appMetadata: getAppMetadata("testapp", "staging", ""),
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			err := dbutil.AutoMigrate()
			r.NoError(err)
			defer purgeTables(r)

			for _, val := range testCase.input {
				config.SetConfigValue(val)
			}

			configs := config.GetAllAppConfigs(testCase.appMetadata)

			results := []*model.AppConfigPayload{}
			for _, config := range configs {
				results = append(results, &config.AppConfigPayload)
			}

			r.ElementsMatch(results, testCase.expected)
		})
	}
}

func TestGetAppConfigsForEnvSucceed(t *testing.T) {
	testData := []struct {
		input       []*model.AppConfigPayload
		expected    []*model.AppConfigPayload
		appMetadata *model.AppMetadata
	}{
		{
			// no configs exist -> return empty array
			input:       []*model.AppConfigPayload{},
			expected:    []*model.AppConfigPayload{},
			appMetadata: getAppMetadata("testapp", "rdev", ""),
		},
		{
			// configs exist, stack specified -> ignores stack and returns only env-level configs
			input: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				getAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			expected: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			},
			appMetadata: getAppMetadata("testapp", "rdev", "foo"),
		},
		{
			// configs exist, different env specified -> returns only configs for specified env
			input: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				getAppConfigPayload("testapp", "staging", "", "KEY2", "val2"),
				getAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			expected: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
				getAppConfigPayload("testapp", "staging", "", "KEY2", "val2"),
			},
			appMetadata: getAppMetadata("testapp", "staging", ""),
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			err := dbutil.AutoMigrate()
			r.NoError(err)
			defer purgeTables(r)

			for _, val := range testCase.input {
				config.SetConfigValue(val)
			}

			configs, err := config.GetAppConfigsForEnv(testCase.appMetadata)
			r.NoError(err)

			results := []*model.AppConfigPayload{}
			for _, config := range configs {
				results = append(results, &config.AppConfigPayload)
				r.Equal(config.Source, "environment")
			}

			r.ElementsMatch(results, testCase.expected)
		})
	}
}

func TestGetAppConfigsForStackSucceed(t *testing.T) {
	type expected struct {
		payload *model.AppConfigPayload
		source  string
	}
	testData := []struct {
		input       []*model.AppConfigPayload
		expected    []expected
		appMetadata *model.AppMetadata
	}{
		{
			// no configs exist -> return empty array
			input:       []*model.AppConfigPayload{},
			expected:    []expected{},
			appMetadata: getAppMetadata("testapp", "rdev", ""),
		},
		{
			// configs exist -> returns stack and env-level configs with overrides applied
			input: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				getAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
				getAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
			},
			expected: []expected{
				{
					payload: getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
					source:  "stack",
				},
				{
					payload: getAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
					source:  "environment",
				},
			},
			appMetadata: getAppMetadata("testapp", "rdev", "foo"),
		},
		{
			// configs exist, different env specified -> returns correct stack and env-level configs with overrides applied
			input: []*model.AppConfigPayload{
				getAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				getAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
				getAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				getAppConfigPayload("testapp", "staging", "stg", "KEY2", "val-stg"),
				getAppConfigPayload("testapp", "staging", "", "KEY2", "val2"),
				getAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			expected: []expected{
				{
					payload: getAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
					source:  "environment",
				},
				{
					payload: getAppConfigPayload("testapp", "staging", "stg", "KEY2", "val-stg"),
					source:  "stack",
				},
			},
			appMetadata: getAppMetadata("testapp", "staging", "stg"),
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			err := dbutil.AutoMigrate()
			r.NoError(err)
			defer purgeTables(r)

			for _, val := range testCase.input {
				config.SetConfigValue(val)
			}

			configs := config.GetAppConfigsForStack(testCase.appMetadata)

			results := []expected{}
			for _, config := range configs {
				results = append(results, expected{
					payload: &config.AppConfigPayload,
					source:  config.Source,
				})
			}

			r.ElementsMatch(results, testCase.expected)
		})
	}
}
