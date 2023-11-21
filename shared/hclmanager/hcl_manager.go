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
	errs "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
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
		return errs.Wrap(err, "Unable to get stack config")
	}

	moduleSource := h.HappyConfig.GetModuleSource()

	if source, ok := stackConfig["source"]; ok {
		moduleSource = source.(string)
	}

	_, modulePath, _, err := tf.ParseModuleSource(moduleSource)
	if err != nil {
		return errs.Wrap(err, "unable to parse module path out")
	}
	modulePathParts := strings.Split(modulePath, "/")
	moduleName := modulePathParts[len(modulePathParts)-1]

	tempDir, err := os.MkdirTemp("", moduleName)
	if err != nil {
		return errs.Wrap(err, "Unable to create temp directory")
	}
	defer os.RemoveAll(tempDir)

	// Download the module source
	err = getter.GetAny(tempDir, moduleSource)
	if err != nil {
		return fmt.Errorf("%w: %w", err, tf.ErrUnableToDownloadModuleSource)
	}

	mod, diags := tfconfig.LoadModule(tempDir)
	if diags.HasErrors() {
		return errs.Wrap(err, "Unable to parse out variables or outputs from the module")
	}

	tfDirPath := h.HappyConfig.TerraformDirectory()

	happyProjectRoot := h.HappyConfig.GetProjectRoot()
	srcDir := filepath.Join(happyProjectRoot, tfDirPath)

	gen := tf.NewTfGenerator(h.HappyConfig)

	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		err = os.MkdirAll(srcDir, 0777)
		if err != nil {
			return errs.Wrapf(err, "unable to create terraform directory: %s", srcDir)
		}
	}

	log.Debugf("Generating terraform files in %s", srcDir)

	err = gen.GenerateMain(srcDir, moduleSource, mod.Variables)
	if err != nil {
		return errs.Wrap(err, "Unable to generate main.tf")
	}

	err = gen.GenerateProviders(srcDir)
	if err != nil {
		return errs.Wrap(err, "Unable to generate providers.tf")
	}

	err = gen.GenerateVersions(srcDir)
	if err != nil {
		return errs.Wrap(err, "Unable to generate versions.tf")
	}

	err = gen.GenerateOutputs(srcDir, mod.Outputs)
	if err != nil {
		return errs.Wrap(err, "Unable to generate outputs.tf")
	}

	err = gen.GenerateVariables(srcDir)
	if err != nil {
		return errs.Wrap(err, "Unable to generate variables.tf")
	}

	return nil
}

func (h HclManager) Ingest(ctx context.Context) error {
	stackDefaults := map[string]any{}
	moduleCalls := map[string]tf.ModuleCall{}

	// Read configuration from all environments
	for name, environment := range h.HappyConfig.GetData().Environments {
		tfDirPath := environment.TerraformDirectory

		happyProjectRoot := h.HappyConfig.GetProjectRoot()

		parser := tf.NewTfParser()
		moduleCall, err := parser.ParseModuleCall(happyProjectRoot, tfDirPath)
		if err != nil {
			return errs.Wrapf(err, "Unable to parse a stack module call for environment '%s'", name)
		}

		moduleCall.Parameters = util.DeepCleanup(moduleCall.Parameters)
		moduleCalls[name] = moduleCall
		stackDefaults = moduleCall.Parameters
	}

	// Determine common stack defaults
	for _, moduleCall := range moduleCalls {
		stackDefaults = util.DeepIntersect(stackDefaults, moduleCall.Parameters)
	}

	// Figure out stack overrides
	for name, moduleCall := range moduleCalls {
		stackOverrides := util.DeepCleanup(util.DeepDiff(stackDefaults, moduleCall.Parameters))
		environment := h.HappyConfig.GetData().Environments[name]
		environment.StackOverrides = stackOverrides
		h.HappyConfig.GetData().Environments[name] = environment
	}

	h.HappyConfig.SetStackDefaults(stackDefaults)
	return errs.Wrap(h.HappyConfig.Save(), "Unable to save happy config")
}

func (h HclManager) Validate(ctx context.Context) error {
	for name, environment := range h.HappyConfig.GetData().Environments {
		tfDirPath := environment.TerraformDirectory

		happyProjectRoot := h.HappyConfig.GetProjectRoot()

		parser := tf.NewTfParser()
		moduleCall, err := parser.ParseModuleCall(happyProjectRoot, tfDirPath)
		if err != nil {
			return errs.Wrapf(err, "Unable to parse a stack module call for environment '%s'", name)
		}

		if moduleCall.Parameters["source"] == nil {
			return errs.Errorf("module source is not set for terraform code in %s", filepath.Join(happyProjectRoot, tfDirPath))
		}

		moduleSource := moduleCall.Parameters["source"].(string)
		isLocalReference := strings.HasPrefix(moduleSource, "./modules/")
		if !isLocalReference {
			_, moduleName, _, err := tf.ParseModuleSource(moduleSource)
			if err != nil {
				return errs.Wrapf(err, "unable to parse module source for environment '%s'", moduleSource)
			}
			moduleName = strings.TrimPrefix(moduleName, "terraform/modules/")
			expectedModuleNames := h.HappyConfig.GetModuleNames()
			if _, ok := expectedModuleNames[moduleName]; !ok {
				return errs.Errorf("module name '%s' does not match, expected '%v'", moduleName, strings.Join(maps.Keys(expectedModuleNames), ","))
			}
		}
	}
	return nil
}

// "The best way to predict the future is to invent it." - Alan Kay
