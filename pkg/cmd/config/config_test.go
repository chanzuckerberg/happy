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

func TestSetConfigValueSucceed(t *testing.T) {
	testData := []struct {
		input    []*model.AppConfigPayload
		expected []*model.AppConfigPayload
	}{
		{
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// same key, different value
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val2"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val2"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// different key, different value
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// same key BUT with a stack name, different value
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val2"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val2"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// same key and stack name, different value
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val3"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val3"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// same key, different stack name, different value
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val3"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val3"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// different env, different value, everything else the same
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// different app, different value, everything else the same
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("anotherapp", "rdev", "", "KEY1", "val2"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("anotherapp", "rdev", "", "KEY1", "val2"),
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
				_, err := config.SetConfigValue(input)
				r.NoError(err)
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
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY2", "val-foo2"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:         model.NewAppConfigLookupPayload("testapp", "rdev", "", "KEY2"),
			expected:       model.NewAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
			expectedSource: "environment",
		},
		{
			// stack and env configs exist, stack specified -> return stack config
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY2", "val-foo2"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:         model.NewAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expected:       model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
			expectedSource: "stack",
		},
		{
			// only env configs exist, stack specified -> return env config
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:         model.NewAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expected:       model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			expectedSource: "environment",
		},
		{
			// only stack configs exist, stack specified -> return stack config
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:         model.NewAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expected:       model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			expectedSource: "stack",
		},
		{
			// only stack configs exist, stack not specified -> return null
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:   model.NewAppConfigLookupPayload("testapp", "rdev", "", "KEY1"),
			expected: nil,
		},
		{
			// no configs, stack not specified -> return null
			input:    []*model.AppConfigPayload{},
			lookup:   model.NewAppConfigLookupPayload("testapp", "rdev", "", "KEY1"),
			expected: nil,
		},
		{
			// no configs, stack specified -> return null
			input:    []*model.AppConfigPayload{},
			lookup:   model.NewAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expected: nil,
		},
		{
			// stack config exists, different stack specified -> return null
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:   model.NewAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
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
				_, err := config.SetConfigValue(input)
				r.NoError(err)
			}

			result, err := config.GetResolvedAppConfig(testCase.lookup)
			r.NoError(err)

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
			deleteCriteria:   model.NewAppConfigLookupPayload("testapp", "rdev", "", "KEY1"),
			expectedResult:   nil,
			remainingConfigs: []*model.AppConfigPayload{},
		},
		{
			// no configs exist, stack specified -> delete nothing
			input:            []*model.AppConfigPayload{},
			deleteCriteria:   model.NewAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expectedResult:   nil,
			remainingConfigs: []*model.AppConfigPayload{},
		},
		{
			// configs exist, no stack specified -> deletes env-level config
			input: []*model.AppConfigPayload{
				// order here is important, stack config needs to be created first to ensure
				// that this test fails when empty stack criteria is not added in the DeleteAppConfig
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			},
			deleteCriteria: model.NewAppConfigLookupPayload("testapp", "rdev", "", "KEY1"),
			expectedResult: model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			remainingConfigs: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
			},
		},
		{
			// configs exist, stack specified -> deletes stack-level config
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			deleteCriteria: model.NewAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expectedResult: model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
			remainingConfigs: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
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
				_, err := config.SetConfigValue(input)
				r.NoError(err)
			}

			result, err := config.DeleteAppConfig(testCase.deleteCriteria)
			r.NoError(err)

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
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			appMetadata: getAppMetadata("testapp", "rdev", "foo"),
		},
		{
			// configs exist, different env specified -> returns only configs for specified env
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY2", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY2", "val2"),
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
				_, err := config.SetConfigValue(val)
				r.NoError(err)
			}

			configs, err := config.GetAllAppConfigs(testCase.appMetadata)
			r.NoError(err)

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
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			},
			appMetadata: getAppMetadata("testapp", "rdev", "foo"),
		},
		{
			// configs exist, different env specified -> returns only configs for specified env
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY2", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY2", "val2"),
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
				_, err := config.SetConfigValue(val)
				r.NoError(err)
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
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
			},
			expected: []expected{
				{
					payload: model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
					source:  "stack",
				},
				{
					payload: model.NewAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
					source:  "environment",
				},
			},
			appMetadata: getAppMetadata("testapp", "rdev", "foo"),
		},
		{
			// configs exist, different env specified -> returns correct stack and env-level configs with overrides applied
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "staging", "stg", "KEY2", "val-stg"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY2", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			expected: []expected{
				{
					payload: model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
					source:  "environment",
				},
				{
					payload: model.NewAppConfigPayload("testapp", "staging", "stg", "KEY2", "val-stg"),
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
				_, err := config.SetConfigValue(val)
				r.NoError(err)
			}

			configs, err := config.GetAppConfigsForStack(testCase.appMetadata)
			r.NoError(err)

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

func TestCopyAppConfigSucceed(t *testing.T) {
	testData := []struct {
		seeds       []*model.AppConfigPayload
		copyPayload *model.CopyAppConfigPayload
		expected    []*model.AppConfigPayload
	}{
		{
			// no configs exist -> nothing gets copied
			seeds:       []*model.AppConfigPayload{},
			copyPayload: model.NewCopyAppConfigPayload("testapp", "rdev", "", "staging", "", "KEY1"),
			expected:    []*model.AppConfigPayload{},
		},
		{
			// configs exist but don't match -> nothing gets copied
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			},
			copyPayload: model.NewCopyAppConfigPayload("testapp", "rdev", "", "staging", "", "KEY1"),
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			},
		},
		{
			// configs exist but don't match -> nothing gets copied
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			},
			copyPayload: model.NewCopyAppConfigPayload("testapp", "rdev", "bar", "staging", "", "KEY1"),
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			},
		},
		{
			// matching env config exists -> copy env config
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			},
			copyPayload: model.NewCopyAppConfigPayload("testapp", "rdev", "", "staging", "", "KEY1"),
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val1"),
			},
		},
		{
			// matching stack config exists -> copy stack config
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val1"),
			},
			copyPayload: model.NewCopyAppConfigPayload("testapp", "rdev", "bar", "staging", "", "KEY1"),
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val1"),
			},
		},
		{
			// test copying to different stack in same env
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val2"),
			},
			copyPayload: model.NewCopyAppConfigPayload("testapp", "rdev", "bar", "rdev", "foo", "KEY1"),
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			},
		},
		{
			// test copying to same stack in same env
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val1"),
			},
			copyPayload: model.NewCopyAppConfigPayload("testapp", "rdev", "bar", "rdev", "bar", "KEY1"),
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val1"),
			},
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			err := dbutil.AutoMigrate()
			r.NoError(err)
			defer purgeTables(r)

			for _, input := range testCase.seeds {
				_, err := config.SetConfigValue(input)
				r.NoError(err)
			}

			_, err = config.CopyAppConfig(testCase.copyPayload)
			r.NoError(err)

			db := dbutil.GetDB()
			configs := []*model.AppConfig{}
			db.Find(&configs)

			results := []*model.AppConfigPayload{}
			for _, config := range configs {
				results = append(results, &config.AppConfigPayload)
			}

			r.EqualValues(testCase.expected, results)
		})
	}
}
