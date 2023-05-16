package tf

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type ModuleCall struct {
	Parameters map[string]any
}

type TfParser struct {
}

func NewTfParser() TfParser {
	return TfParser{}
}

var moduleBlockSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "module",
			LabelNames: []string{"name"},
		},
	},
}

func (tf TfParser) ParseServices(dir string) (map[string]bool, error) {
	var services map[string]bool = map[string]bool{}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".terraform" || d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".tf" {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		f, diags := hclsyntax.ParseConfig(b, path, hcl.Pos{Line: 1, Column: 1})
		if diags.HasErrors() {
			return errors.Wrapf(diags.Errs()[0], "failed to parse %s", path)
		}

		content, _, diags := f.Body.PartialContent(moduleBlockSchema)
		if diags.HasErrors() {
			return errors.New("Terraform code has errors")
		}

		for _, block := range content.Blocks {
			if block.Type != "module" {
				continue
			}

			attrs, diags := block.Body.JustAttributes()
			if diags.HasErrors() {
				return errors.New("Terraform code has errors")
			}
			var sourceAttr *hcl.Attribute
			var ok bool
			if sourceAttr, ok = attrs["source"]; !ok {
				// Module without a source
				continue
			}

			source, diags := sourceAttr.Expr.(*hclsyntax.TemplateExpr).Parts[0].Value(nil)
			if diags.HasErrors() {
				return errors.New("Terraform code has errors")
			}

			if !strings.Contains(source.AsString(), "modules/happy-stack-") {
				// Not a happy stack module
				continue
			}

			if servicesAttr, ok := attrs["services"]; ok {
				switch servicesAttr.Expr.(type) {
				case *hclsyntax.ObjectConsExpr:
					for _, item := range servicesAttr.Expr.(*hclsyntax.ObjectConsExpr).Items {
						key, _ := item.KeyExpr.Value(nil)
						services[key.AsString()] = true
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		return services, errors.Wrap(err, "failed to parse terraform files")
	}
	return services, nil
}

func (tf TfParser) ParseModuleCall(dir string) (ModuleCall, error) {
	moduleCall := ModuleCall{Parameters: map[string]any{}}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".terraform" || d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".tf" {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		f, diags := hclsyntax.ParseConfig(b, path, hcl.Pos{Line: 1, Column: 1})
		if diags.HasErrors() {
			return errors.Wrapf(diags.Errs()[0], "failed to parse %s", path)
		}

		content, _, diags := f.Body.PartialContent(moduleBlockSchema)
		if diags.HasErrors() {
			return errors.New("Terraform code has errors")
		}

		for _, block := range content.Blocks {
			if block.Type != "module" {
				continue
			}

			attrs, diags := block.Body.JustAttributes()
			if diags.HasErrors() {
				return errors.New("Terraform code has errors")
			}
			var sourceAttr *hcl.Attribute
			var ok bool
			if sourceAttr, ok = attrs["source"]; !ok {
				// Module without a source
				continue
			}

			source, diags := sourceAttr.Expr.(*hclsyntax.TemplateExpr).Parts[0].Value(nil)
			if diags.HasErrors() {
				return errors.New("Terraform code has errors")
			}

			if !strings.Contains(source.AsString(), "modules/happy-stack-") {
				// Not a happy stack module
				continue
			}

			for _, attr := range attrs {
				value, diag := attr.Expr.Value(nil)
				if diag.HasErrors() {
					log.Warnf("Attribute %s cannot be read properly: %s", attr.Name, diag.Errs()[0].Error())
					//continue
				}

				v, err := decodeValue(value)
				if err != nil {
					continue
				}
				if v != nil {
					moduleCall.Parameters[attr.Name] = v
				}
			}

			// These variables below are managed by the generator, and we don't have a need to read them, interpret them or store them.
			delete(moduleCall.Parameters, "image_tag")
			delete(moduleCall.Parameters, "image_tags")
			delete(moduleCall.Parameters, "k8s_namespace")
			delete(moduleCall.Parameters, "stack_name")
			delete(moduleCall.Parameters, "deployment_stage")
		}

		return nil
	})
	if err != nil {
		return moduleCall, errors.Wrap(err, "failed to parse terraform files")
	}
	return moduleCall, nil
}

func decodeValue(ctyValue cty.Value) (any, error) {
	if ctyValue.IsNull() {
		return nil, nil
	}

	if ctyValue.Type().IsMapType() || ctyValue.Type().IsObjectType() {
		m := map[string]any{}
		for key, value := range ctyValue.AsValueMap() {

			v, err := decodeValue(value)
			if err != nil {
				return nil, err
			}
			m[key] = v
		}
		return m, nil
	}

	switch ctyValue.Type() {
	case cty.Bool:
		var v bool
		err := gocty.FromCtyValue(ctyValue, &v)
		return v, err
	case cty.Number:
		var v int64
		err := gocty.FromCtyValue(ctyValue, &v)
		return v, err
	case cty.String:
		var v string
		err := gocty.FromCtyValue(ctyValue, &v)
		return v, err
	case cty.DynamicPseudoType:
		return nil, nil
	default:
		return nil, errors.Errorf("unsupported type %s", ctyValue.Type().FriendlyName())
	}
}
