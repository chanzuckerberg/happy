package tf

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestCtyMerge(t *testing.T) {
	r := require.New(t)

	v1 := mergeCtyValues(cty.StringVal("foo"), cty.StringVal("bar"))
	r.Equal(cty.StringVal("bar"), v1)

	v1 = mergeCtyValues(cty.StringVal("foo"), cty.NullVal(cty.String))
	r.Equal(cty.StringVal("foo"), v1)

	v2 := mergeCtyValues(cty.BoolVal(true), cty.BoolVal(false))
	r.Equal(cty.BoolVal(false), v2)

	v2 = mergeCtyValues(cty.BoolVal(true), cty.NullVal(cty.Bool))
	r.Equal(cty.BoolVal(true), v2)

	v3 := mergeCtyValues(cty.NumberIntVal(1), cty.NumberIntVal(2))
	r.Equal(cty.NumberIntVal(2), v3)

	v4 := mergeCtyValues(cty.ObjectVal(map[string]cty.Value{"foo": cty.StringVal("bar")}), cty.ObjectVal(map[string]cty.Value{"foo": cty.StringVal("baz")}))
	r.Equal(cty.ObjectVal(map[string]cty.Value{"foo": cty.StringVal("baz")}), v4)

	v4 = mergeCtyValues(cty.ObjectVal(map[string]cty.Value{"foo": cty.ListVal([]cty.Value{cty.StringVal("bar")})}), cty.ObjectVal(map[string]cty.Value{"foo": cty.ListVal([]cty.Value{cty.StringVal("baz")})}))
	r.Equal(cty.ObjectVal(map[string]cty.Value{"foo": cty.ListVal([]cty.Value{cty.StringVal("baz")})}), v4)

	v4 = mergeCtyValues(cty.ObjectVal(map[string]cty.Value{"foo": cty.ListVal([]cty.Value{cty.StringVal("bar")})}), cty.ObjectVal(nil))
	r.Equal(cty.ObjectVal(map[string]cty.Value{"foo": cty.ListVal([]cty.Value{cty.StringVal("bar")})}), v4)
}

func TestCtyCleanup(t *testing.T) {
	r := require.New(t)

	v1 := cleanupCtyValue(cty.StringVal("foo"))
	r.Equal(cty.StringVal("foo"), v1)

	v1 = cleanupCtyValue(cty.NullVal(cty.String))
	r.Equal(cty.NullVal(cty.String), v1)

	v2 := cleanupCtyValue(cty.BoolVal(true))
	r.Equal(cty.BoolVal(true), v2)

	v2 = cleanupCtyValue(cty.NullVal(cty.Bool))
	r.Equal(cty.NullVal(cty.Bool), v2)

	v3 := cleanupCtyValue(cty.ObjectVal(nil))
	r.Equal(cty.ObjectVal(nil), v3)

	v3 = cleanupCtyValue(cty.ObjectVal(map[string]cty.Value{"foo": cty.NullVal(cty.String)}))
	r.Equal(cty.ObjectVal(nil), v3)

	v3 = cleanupCtyValue(cty.ObjectVal(map[string]cty.Value{"foo": cty.NullVal(cty.String), "bar": cty.NumberIntVal(1)}))
	r.Equal(cty.ObjectVal(map[string]cty.Value{"bar": cty.NumberIntVal(1)}), v3)
}

func TestCtyEscape(t *testing.T) {
	r := require.New(t)
	val := cty.StringVal("{frontend}:${var.tag}")
	tk1 := string(unescape(hclwrite.TokensForValue(val)).Bytes())
	tk2 := tokens(val.AsString())

	r.Equal(len(tk1), len(tk2))
}
