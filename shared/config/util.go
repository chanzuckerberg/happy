package config

import (
	"github.com/spf13/cobra"
)

func GetHappyConfigForCmd(cmd *cobra.Command) (*HappyConfig, error) {
	bootstrapConfig, err := NewBootstrapConfig(cmd)
	if err != nil {
		return nil, err
	}
	happyConfig, err := NewHappyConfig(bootstrapConfig)
	if err != nil {
		return nil, err
	}
	return happyConfig, nil
}
