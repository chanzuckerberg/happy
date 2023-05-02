package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

const (
	sourceDefault = "default"
)

var variableBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name: "description",
		},
		{
			Name: "default",
		},
		{
			Name: "type",
		},
		{
			Name: "sensitive",
		},
		{
			Name: "nullable",
		},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type: "validation",
		},
	},
}

func ParseServices(dir string) (map[string]bool, error) {
	var services map[string]bool = map[string]bool{}
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "module",
				LabelNames: []string{"name"},
			},
		},
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".terraform" || info.Name() == ".git" {
				return filepath.SkipDir
			}
		}
		if !info.IsDir() && filepath.Ext(path) == ".tf" {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			f, diags := hclsyntax.ParseConfig(b, path, hcl.Pos{Line: 1, Column: 1})
			if diags.HasErrors() {
				return errors.Wrapf(diags.Errs()[0], "failed to parse %s", path)
			}

			content, _, diags := f.Body.PartialContent(schema)
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

				if !strings.Contains(source.AsString(), "modules/happy-stack-eks") && !strings.Contains(source.AsString(), "modules/happy-stack-ecs") {
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
		}
		return nil
	})
	if err != nil {
		return services, errors.Wrap(err, "failed to parse terraform files")
	}
	return services, nil
}

func ParseVariables(dir string) ([]Variable, error) {
	variables := []Variable{}
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
			},
			{
				Type:       "validation",
				LabelNames: []string{"name"},
			},
		},
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".terraform" || info.Name() == ".git" {
				return filepath.SkipDir
			}
		}
		if !info.IsDir() && filepath.Ext(path) == ".tf" {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			f, diags := hclsyntax.ParseConfig(b, path, hcl.Pos{Line: 1, Column: 1})
			if diags.HasErrors() {
				return errors.Wrapf(diags.Errs()[0], "failed to parse %s", path)
			}

			content, _, diags := f.Body.PartialContent(schema)
			if diags.HasErrors() {
				return errors.New("Terraform code has errors")
			}
			for _, block := range content.Blocks {
				if block.Type != "variable" {
					continue
				}

				v, diags := decodeVariableBlock(block)
				if diags.HasErrors() {
					logrus.Infof("diagnostics: %v", diags)
					return errors.New("Terraform code has errors")
				}

				variables = append(variables, *v)
			}
		}
		return err
	})
	if err != nil {
		return variables, errors.Wrap(err, "failed to parse terraform files")
	}
	return variables, nil
}

func decodeVariableBlock(block *hcl.Block) (*Variable, hcl.Diagnostics) {
	v := &Variable{
		Name:      block.Labels[0],
		DeclRange: block.DefRange,
	}

	content, diags := block.Body.Content(variableBlockSchema)

	if attr, exists := content.Attributes["type"]; exists {
		ty, tyDefaults, tyDiags := decodeVariableType(attr.Expr)
		diags = append(diags, tyDiags...)
		v.ConstraintType = ty
		v.TypeDefaults = tyDefaults
		v.Type = ty.WithoutOptionalAttributesDeep()
		val, _ := attr.Expr.Value(nil)
		val = v.TypeDefaults.Apply(val)
	}

	if attr, exists := content.Attributes["default"]; exists {
		val, valDiags := attr.Expr.Value(nil)
		diags = append(diags, valDiags...)

		if v.ConstraintType != cty.NilType {
			var err error
			if v.TypeDefaults != nil && !val.IsNull() {
				val = v.TypeDefaults.Apply(val)
			}
			val, err = convert.Convert(val, v.ConstraintType)
			if err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid default value for variable",
					Detail:   fmt.Sprintf("This default value is not compatible with the variable's type constraint: %s.", err),
					Subject:  attr.Expr.Range().Ptr(),
				})
				val = cty.DynamicVal
			}
		}

		v.Default = val
	}

	return v, diags
}

func decodeVariableType(expr hcl.Expression) (cty.Type, *typeexpr.Defaults, hcl.Diagnostics) {
	switch hcl.ExprAsKeyword(expr) {
	case "list":
		return cty.List(cty.DynamicPseudoType), nil, nil
	case "map":
		return cty.Map(cty.DynamicPseudoType), nil, nil
	default:
	}

	ty, typeDefaults, diags := typeexpr.TypeConstraintWithDefaults(expr)
	if diags.HasErrors() {
		return cty.DynamicPseudoType, nil, diags
	}

	return ty, typeDefaults, diags
}

type Value struct {
	Value  cty.Value
	Source string
	Expr   hcl.Expression
	Range  hcl.Range
}

type Variable struct {
	Name        string
	Description string
	Default     cty.Value

	Type           cty.Type
	ConstraintType cty.Type
	TypeDefaults   *typeexpr.Defaults

	DeclRange hcl.Range
}
