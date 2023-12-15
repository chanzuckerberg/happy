package util

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/require"
)

func TestGetAWSProfiles(t *testing.T) {
	// Create a temporary config file for testing
	configFile := createTempConfigFile()
	defer removeTempConfigFile(configFile)

	// Add some test profiles to the config file
	addTestProfilesToConfigFile(configFile)

	// Call the function being tested
	profiles, err := GetAWSProfiles(func(lo *config.LoadOptions) {
		lo.SharedConfigFiles = []string{configFile}
	})

	// Check for any errors
	require.NoError(t, err)

	// Check the expected profiles
	expectedProfiles := []string{"profile1", "profile2", "profile3"}
	require.ElementsMatch(t, expectedProfiles, profiles)
}

// Helper functions for test setup and teardown
func createTempConfigFile() string {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "config")
	if err != nil {
		panic(err)
	}

	// Get the file path
	filePath := tempFile.Name()

	// Close the file
	tempFile.Close()

	return filePath
}

func removeTempConfigFile(configFile string) {
	// Remove the temporary file
	err := os.Remove(configFile)
	if err != nil {
		panic(err)
	}
}

func addTestProfilesToConfigFile(configFile string) {
	// Open the config file in append mode
	file, err := os.OpenFile(configFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write the test profiles to the config file
	profiles := []string{"profile1", "profile2", "profile3"}
	for _, profile := range profiles {
		_, err := file.WriteString("[profile " + profile + "]\n")
		if err != nil {
			panic(err)
		}
	}
}
