package hclmanager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/util/tf"
	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type HclManager struct {
	HappyConfig *config.HappyConfig
}

func NewHclManager() HclManager {
	return HclManager{}
}

func (h HclManager) WithHappyConfig(happyConfig *config.HappyConfig) HclManager {
	h.HappyConfig = happyConfig
	return h
}

func (h HclManager) Generate(ctx context.Context) error {
	stackConfig, err := h.HappyConfig.GetStackConfig()
	if err != nil {
		return errors.Wrap(err, "Unable to get stack config")
	}
	moduleSource := "git@github.com:chanzuckerberg/happy//terraform/modules/happy-stack-%s?ref=main"
	if h.HappyConfig.TaskLaunchType() == util.LaunchTypeK8S {
		moduleSource = fmt.Sprintf(moduleSource, "eks")
	} else {
		moduleSource = fmt.Sprintf(moduleSource, "ecs")
	}

	if source, ok := stackConfig["source"]; ok {
		moduleSource = source.(string)
	}

	_, modulePath, _, err := tf.ParseModuleSource(moduleSource)
	if err != nil {
		return errors.Wrap(err, "unable to parse module path out")
	}
	modulePathParts := strings.Split(modulePath, "/")
	moduleName := modulePathParts[len(modulePathParts)-1]

	tempDir, err := os.MkdirTemp("", moduleName)
	if err != nil {
		return errors.Wrap(err, "Unable to create temp directory")
	}
	defer os.RemoveAll(tempDir)

	// Download the module source
	err = getter.GetAny(tempDir, moduleSource)
	if err != nil {
		return errors.Wrap(err, "Unable to download module source")
	}

	mod, diags := tfconfig.LoadModule(tempDir)
	if diags.HasErrors() {
		return errors.Wrap(err, "Unable to parse out variables or outputs from the module")
	}

	tfDirPath := h.HappyConfig.TerraformDirectory()

	happyProjectRoot := h.HappyConfig.GetProjectRoot()
	srcDir := filepath.Join(happyProjectRoot, tfDirPath)

	gen := tf.NewTfGenerator(h.HappyConfig)

	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		err = os.MkdirAll(srcDir, 0777)
		if err != nil {
			return errors.Wrapf(err, "unable to create terraform directory: %s", srcDir)
		}
	}

	log.Debugf("Generating terraform files in %s", srcDir)

	err = gen.GenerateMain(srcDir, moduleSource, mod.Variables)
	if err != nil {
		return errors.Wrap(err, "Unable to generate main.tf")
	}

	err = gen.GenerateProviders(srcDir)
	if err != nil {
		return errors.Wrap(err, "Unable to generate providers.tf")
	}

	err = gen.GenerateVersions(srcDir)
	if err != nil {
		return errors.Wrap(err, "Unable to generate versions.tf")
	}

	err = gen.GenerateOutputs(srcDir, mod.Outputs)
	if err != nil {
		return errors.Wrap(err, "Unable to generate outputs.tf")
	}

	err = gen.GenerateVariables(srcDir)
	if err != nil {
		return errors.Wrap(err, "Unable to generate variables.tf")
	}

	return nil
}
