package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/spf13/cobra"
)

type A struct {
	ConfigVersion         string   `json:"config_version"`
	TerraformVersion      string   `json:"terraform_version"`
	DefaultEnv            string   `json:"default_env"`
	App                   string   `json:"app"`
	DefaultComposeEnvFile string   `json:"default_compose_env_file"`
	SliceDefaultTag       string   `json:"slice_default_tag"`
	Services              []string `json:"services"`
	Slices                struct {
		Frontend struct {
			BuildImages []string `json:"build_images"`
			Profile     string   `json:"profile"`
		} `json:"frontend"`
		Backend struct {
			BuildImages []string `json:"build_images"`
			Profile     string   `json:"profile"`
		} `json:"backend"`
		Fullstack struct {
			BuildImages []string `json:"build_images"`
			Profile     string   `json:"profile"`
		} `json:"fullstack"`
		Batch struct {
			BuildImages []string `json:"build_images"`
			Profile     string   `json:"profile"`
		} `json:"batch"`
		Nextstrain struct {
			BuildImages []string `json:"build_images"`
			Profile     string   `json:"profile"`
		} `json:"nextstrain"`
		Pangolin struct {
			BuildImages []string `json:"build_images"`
			Profile     string   `json:"profile"`
		} `json:"pangolin"`
		Gisaid struct {
			BuildImages []string `json:"build_images"`
			Profile     string   `json:"profile"`
		} `json:"gisaid"`
	} `json:"slices"`
	Environments struct {
		Rdev struct {
			AwsProfile         string `json:"aws_profile"`
			SecretArn          string `json:"secret_arn"`
			TerraformDirectory string `json:"terraform_directory"`
			LogGroupPrefix     string `json:"log_group_prefix"`
			TaskLaunchType     string `json:"task_launch_type"`
			AutoRunMigrations  bool   `json:"auto_run_migrations"`
		} `json:"rdev"`
		Staging struct {
			AwsProfile         string `json:"aws_profile"`
			SecretArn          string `json:"secret_arn"`
			TerraformDirectory string `json:"terraform_directory"`
			DeleteProtected    bool   `json:"delete_protected"`
			AutoRunMigrations  bool   `json:"auto_run_migrations"`
			LogGroupPrefix     string `json:"log_group_prefix"`
			TaskLaunchType     string `json:"task_launch_type"`
		} `json:"staging"`
		Prod struct {
			AwsProfile         string `json:"aws_profile"`
			SecretArn          string `json:"secret_arn"`
			TerraformDirectory string `json:"terraform_directory"`
			DeleteProtected    bool   `json:"delete_protected"`
			AutoRunMigrations  bool   `json:"auto_run_migrations"`
			LogGroupPrefix     string `json:"log_group_prefix"`
			TaskLaunchType     string `json:"task_launch_type"`
		} `json:"prod"`
	} `json:"environments"`
	Tasks struct {
		Migrate []string `json:"migrate"`
		Delete  []string `json:"delete"`
	} `json:"tasks"`
}

func init() {
	rootCmd.AddCommand(initCmd)
	config.ConfigureCmdWithBootstrapConfig(initCmd)
}

var initCmd = &cobra.Command{
	Use:          "init",
	Short:        "init repo",
	Long:         "Scaffold the repo",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		//ctx := cmd.Context()

		_, err := config.NewBootstrapConfig(cmd)
		if err == nil {
			return errors.New("this repo was previously initialized")
		}

		currentDir, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "unable to figure out current directory")
		}

		paths := []string{".happy", ".happy/terraform/envs/dev", ".happy/terraform/envs/staging", ".happy/terraform/envs/prod", ".happy/terraform/modules", ".github/workflows"}

		for _, path := range paths {
			absDir := filepath.Join(currentDir, path)
			err = os.MkdirAll(absDir, os.ModePerm)
			if err != nil {
				return errors.Wrapf(err, "unable to create a directory: %s", absDir)
			}
			_, err = os.Create(filepath.Join(absDir, ".gitignore"))
			if err != nil {
				return errors.Wrapf(err, "unable to create a .gitignore file in a %s directory", absDir)
			}
		}

		configFile, err := os.Create(filepath.Join(currentDir, ".happy/config.json"))
		if err != nil {
			return errors.Wrap(err, "unable to create a .happy/config.json")
		}

		_, appName := filepath.Split(currentDir)
		prompt := &survey.Input{Message: "App Name?", Default: appName}
		err = survey.AskOne(prompt, &appName)
		if err != nil {
			return errors.Wrap(err, "unable to prompt")
		}

		configData := config.ConfigData{
			ConfigVersion:         "v2",
			TerraformVersion:      "0.13.5",
			DefaultEnv:            "dev",
			App:                   appName,
			DefaultComposeEnvFile: ".env.ecr",
			Environments:          make(map[string]config.Environment),
			Tasks:                 make(map[string][]string),
			SliceDefaultTag:       "branch-trunk",
			Services:              make([]string, 0),
			Slices:                make(map[string]config.Slice),
		}

		bytes, err := json.MarshalIndent(configData, "", "    ")
		if err != nil {
			return errors.Wrap(err, "cannot serialize happy config")
		}

		_, err = configFile.Write(bytes)
		if err != nil {
			return errors.Wrap(err, "cannot write out a happy config")
		}

		err = configFile.Close()
		if err != nil {
			return errors.Wrap(err, "failed to close the happy config file")
		}

		// TODO: write out .env.ecr

		return nil
	},
}
