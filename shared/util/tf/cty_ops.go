package tf

import (
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

var (
	hclUnescaper = strings.NewReplacer("$${", "${", "%%{", "%{")
)

// Merge two cty.Value objects, preferring the second if both are non-null (first one is used as a source of defaults)
func mergeCtyValues(v1, v2 cty.Value) cty.Value {
	if v1.Type().IsPrimitiveType() && v2.Type().IsPrimitiveType() {
		if !v2.IsNull() {
			return v2
		}
	}

	if v1.Type().IsListType() && v2.Type().IsListType() {
		if !v2.IsNull() && v2.LengthInt() > 0 {
			return v2
		}
	}

	if v1.Type().IsCollectionType() && v2.Type().IsCollectionType() {
		if !v2.IsNull() && v2.LengthInt() > 0 {
			return v2
		}
	}

	if v1.Type().IsSetType() && v2.Type().IsSetType() {
		if !v2.IsNull() && v2.LengthInt() > 0 {
			return v2
		}
	}

	if v1.Type().IsMapType() && v2.Type().IsMapType() {
		if !v2.IsNull() && v2.LengthInt() > 0 {
			return v2
		}
	}

	if v1.Type().IsTupleType() && v2.Type().IsTupleType() {
		if !v2.IsNull() && v2.LengthInt() > 0 {
			return v2
		}
	}

	if v1.Type().IsObjectType() && v2.Type().IsObjectType() {
		v1m := v1.AsValueMap()
		v2m := v2.AsValueMap()
		for k, v := range v2m {
			if !v.IsNull() {
				v1m[k] = mergeCtyValues(v1m[k], v)
			}
			if v1m[k].IsNull() {
				delete(v1m, k)
			}
		}
		return cty.ObjectVal(v1m)
	}
	return v1
}

// Remove keys with null values from the referenced cty.Value (if applicable)
func cleanupCtyValue(val cty.Value) cty.Value {
	if val.IsNull() {
		return val
	}

	if val.Type().IsObjectType() || val.Type().IsMapType() {
		vm := val.AsValueMap()
		for k, v := range vm {
			v = cleanupCtyValue(v)
			if v.IsNull() {
				delete(vm, k)
				continue
			}
			vm[k] = v
		}
		return cty.ObjectVal(vm)
	}

	return val
}

func cleanupTokens(tokens hclwrite.Tokens) hclwrite.Tokens {
	var out hclwrite.Tokens
	for _, token := range tokens {
		if token.Type != hclsyntax.TokenQuotedLit {
			out = append(out, token)
			continue
		}
		newToken := hclwrite.Token{
			Type:  hclsyntax.TokenQuotedLit,
			Bytes: []byte(unescapeHclString(string(token.Bytes))),
		}

		out = append(out, &newToken)
	}
	return out
}

func unescapeHclString(val string) string {
	return hclUnescaper.Replace(val)
}
