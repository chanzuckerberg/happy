package tf

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
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

const (
	inputVariablesAccessor = "var"
)

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
			return errors.Wrap(diags.Errs()[0], "terraform code has errors")
		}

		for _, block := range content.Blocks {
			if block.Type != "module" {
				continue
			}

			attrs, diags := block.Body.JustAttributes()
			if diags.HasErrors() {
				return errors.Wrap(diags.Errs()[0], "terraform code has errors")
			}
			var sourceAttr *hcl.Attribute
			var ok bool
			if sourceAttr, ok = attrs["source"]; !ok {
				// Module without a source
				continue
			}

			source, diags := sourceAttr.Expr.(*hclsyntax.TemplateExpr).Parts[0].Value(nil)
			if diags.HasErrors() {
				return errors.Wrap(diags.Errs()[0], "terraform code has errors")
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

func (tf TfParser) ParseModuleCall(happyProjectRoot, tfDirPath string) (ModuleCall, error) {
	dir := filepath.Join(happyProjectRoot, tfDirPath)
	moduleCall := ModuleCall{Parameters: map[string]any{}}

	excludedAttributes := map[string]bool{
		"image_tag":        true,
		"image_tags":       true,
		"k8s_namespace":    true,
		"stack_name":       true,
		"deployment_stage": true,
		"stack_prefix":     true,
	}

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
		relativePath, err := filepath.Rel(happyProjectRoot, path)
		if err != nil {
			return err
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
			return errors.Wrap(diags.Errs()[0], "terraform code has errors")
		}

		for _, block := range content.Blocks {
			if block.Type != "module" {
				continue
			}

			attrs, diags := block.Body.JustAttributes()
			if diags.HasErrors() {
				return errors.Wrap(diags.Errs()[0], "terraform code has errors")
			}
			var sourceAttr *hcl.Attribute
			var ok bool
			if sourceAttr, ok = attrs["source"]; !ok {
				// Module without a source
				continue
			}

			source, diags := sourceAttr.Expr.(*hclsyntax.TemplateExpr).Parts[0].Value(nil)
			if diags.HasErrors() {
				return errors.Wrap(diags.Errs()[0], "terraform code has errors")
			}

			if !strings.Contains(source.AsString(), "modules/happy-stack-") {
				// Not a happy stack module
				continue
			}

			tempDir, err := os.MkdirTemp("", "happy-stack-module")
			if err != nil {
				return errors.Wrap(err, "Unable to create temp directory")
			}
			defer os.RemoveAll(tempDir)

			// Download the module source
			err = getter.GetAny(tempDir, source.AsString())
			if err != nil {
				return errors.Wrap(err, "Unable to download module source")
			}

			mod, d := tfconfig.LoadModule(tempDir)
			if d.HasErrors() {
				return errors.Wrapf(d.Err(), "Unable to parse out variables or outputs from the module %s", source.AsString())
			}
			gen := NewTfGenerator(nil)
			variables := gen.PreprocessVars(mod.Variables)

			varMap := map[string]ModuleVariable{}

			for _, variable := range variables {
				varMap[variable.Name] = variable
				if _, ok := attrs[variable.Name]; !ok {
					if variable.Default.IsNull() {
						log.Warnf("Variable '%s' value is not specified in the module call", variable.Name)
					}
				}
			}

			for _, attr := range attrs {

				if _, ok := varMap[attr.Name]; !ok {
					if attr.Name != "source" {
						log.Warnf("%s(%d:%d): attribute '%s' is not a variable of a module", relativePath, attr.Range.Start.Line, attr.Range.Start.Column, attr.Name)
					}
				}

				if _, ok := excludedAttributes[attr.Name]; ok {
					// These variables below are managed by the generator, and we don't have a need to read them, interpret them or store them.
					continue
				}

				inputVariables := map[string]cty.Value{}
				for _, variable := range []string{"image_tag", "image_tags", "k8s_namespace", "stack_name", "deployment_stage", "stack_prefix"} {
					inputVariables[variable] = cty.StringVal("${var." + variable + "}")
				}

				ectx := &hcl.EvalContext{
					Variables: map[string]cty.Value{
						inputVariablesAccessor: cty.ObjectVal(inputVariables),
					},
				}

				value, diag := attr.Expr.Value(ectx)
				if diag.HasErrors() {
					// Referencing other variables
					continue
				}

				v, err := decodeValue(value)
				if err != nil {
					log.Warnf("%s(%d:%d): unable to decode value for attribute '%s': %s", relativePath, attr.Range.Start.Line, attr.Range.Start.Column, attr.Name, err.Error())
					continue
				}

				if v, ok := varMap[attr.Name]; ok {
					err = isFunctionallyCompatible(v.Type, value.Type())
					if err != nil {
						log.Warnf("%s(%d:%d): provided value for attribute '%s' doesn't match the one required by the module: %s", path, attr.Range.Start.Line, attr.Range.Start.Column, attr.Name, err.Error())
						continue
					}
				}

				if v != nil {
					moduleCall.Parameters[attr.Name] = v
				}
			}
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
				return nil, errors.Wrap(err, "failed to decode map/object value")
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
		var list []any
		if ctyValue.Type().IsListType() || ctyValue.Type().IsTupleType() {
			ctyValue.ForEachElement(func(key cty.Value, val cty.Value) (stop bool) {
				v, err := decodeValue(val)
				if err != nil {
					return true
				}
				list = append(list, v)
				return false
			})
			return list, nil
		}
		return nil, errors.Errorf("unsupported type %s", ctyValue.Type().FriendlyName())
	}
}

// Validates if a value of type t2 can be used in a module invocation for an attribute of type t1, and vice versa.
// No validation is performed for the values of the attributes, only their types.
func isFunctionallyCompatible(t1 cty.Type, t2 cty.Type) error {
	// Primitives are always compatible, as long as they are of the same type
	if t1.IsPrimitiveType() && t2.IsPrimitiveType() {
		if t1.FriendlyName() == t2.FriendlyName() {
			return nil
		}
		return errors.Errorf("expected: %s, got: %s", t1.FriendlyName(), t2.FriendlyName())
	}

	// Lists types are compatible if their element types are compatible
	if t1.IsListType() && t2.IsListType() {
		return isFunctionallyCompatible(t1.ElementType(), t2.ElementType())
	}

	// Collection types are compatible if their element types are compatible
	if t1.IsCollectionType() && t2.IsCollectionType() {
		return isFunctionallyCompatible(t1.ElementType(), t2.ElementType())
	}

	// Set types are compatible if their element types are compatible
	if t1.IsSetType() && t2.IsSetType() {
		return isFunctionallyCompatible(t1.ElementType(), t2.ElementType())
	}

	// Map types are compatible if their element types are compatible
	if t1.IsMapType() && t2.IsMapType() {
		return isFunctionallyCompatible(t1.ElementType(), t2.ElementType())
	}

	// Two cases below deal with the same scenario: module has a variable of type map, but the parser treats the passed value
	// as an object (from the parser stand point they are not different). We use the element type of a map and validate it
	// against the value of every attribute of the object.
	if t1.IsMapType() && t2.IsObjectType() {
		for name, attrType := range t2.AttributeTypes() {
			err := isFunctionallyCompatible(t1.ElementType(), attrType)
			if err != nil {
				return errors.Errorf("type mismatch for member '%s': %s", name, err.Error())
			}
		}
		return nil
	}

	if t1.IsObjectType() && t2.IsMapType() {
		for name, attrType := range t1.AttributeTypes() {
			err := isFunctionallyCompatible(attrType, t2.ElementType())
			if err != nil {
				return errors.Errorf("type mismatch for member '%s': %s", name, err.Error())
			}
		}
		return nil
	}

	// Object types are compatible if their attributes are compatible (if present)
	if t1.IsObjectType() && t2.IsObjectType() {
		attrs1 := t1.AttributeTypes()
		attrs2 := t2.AttributeTypes()
		for k1, v1 := range attrs1 {
			if v2, ok := attrs2[k1]; ok {
				err := isFunctionallyCompatible(v1, v2)
				if err != nil {
					return errors.Errorf("type mismatch for attribute '%s': %s", k1, err.Error())
				}
			}
			// TODO: Check for missing or extra attributes
		}
		for k2, v2 := range attrs2 {
			if v1, ok := attrs1[k2]; ok {
				err := isFunctionallyCompatible(v1, v2)
				if err != nil {
					return errors.Errorf("type mismatch for attribute '%s': %s", k2, err.Error())
				}
			}
			// TODO: Check for missing or extra attributes
		}
		return nil
	}

	// Typle types are compatible if their element types are compatible
	if t1.IsTupleType() && t2.IsTupleType() {
		u1 := t1.TupleElementTypes()
		u2 := t2.TupleElementTypes()
		if len(u1) != len(u2) {
			return errors.New("tuple types have different lengths")
		}
		for i := range u1 {
			err := isFunctionallyCompatible(u1[i], u2[i])
			if err != nil {
				return errors.Errorf("type mismatch for tuple element %d: %s", i, err.Error())
			}
		}
		return nil
	}

	// Tuple types are compatible with lists and collections if their element types are compatible
	if (t1.IsCollectionType() || t1.IsListType()) && t2.IsTupleType() {
		if len(t2.TupleElementTypes()) == 0 {
			return nil
		}
		elementTypes := distinctTypes(t2.TupleElementTypes())
		if len(elementTypes) == 1 {
			return isFunctionallyCompatible(t1.ElementType(), t2.TupleElementTypes()[0])
		}
	}

	if t1.IsTupleType() && (t2.IsCollectionType() || t2.IsListType()) {
		if len(t1.TupleElementTypes()) == 0 {
			return nil
		}
		elementTypes := distinctTypes(t1.TupleElementTypes())
		if len(elementTypes) == 1 {
			return isFunctionallyCompatible(t1.TupleElementTypes()[0], t2.ElementType())
		}
	}

	// This is a rather unexpected scenario
	return errors.Errorf("Unable to compare types %s and %s", t1.FriendlyName(), t2.FriendlyName())
}

func distinctTypes(types []cty.Type) []cty.Type {
	m := map[string]cty.Type{}
	for _, t := range types {
		m[t.FriendlyName()] = t
	}
	var result []cty.Type
	for _, t := range m {
		result = append(result, t)
	}
	return result
}
