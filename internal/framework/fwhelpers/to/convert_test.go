package to_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/stretchr/testify/require"
)

type objectModel struct {
	Name types.String `tfsdk:"name"`
}

func (objectModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{"name": types.StringType}
}

func TestObjectConversionsUseAttributeTypesProvider(t *testing.T) {
	t.Parallel()

	null := to.ObjectNull[objectModel]()
	require.True(t, null.IsNull())
	require.Equal(t, objectModel{}.AttributeTypes(), null.AttributeTypes(t.Context()))

	object, diags := to.Object(t.Context(), objectModel{Name: types.StringValue("example")})
	require.False(t, diags.HasError(), diags.Errors())
	require.Equal(t, types.StringValue("example"), object.Attributes()["name"])

	set, diags := to.SetObject(t.Context(), []objectModel{{Name: types.StringValue("example")}})
	require.False(t, diags.HasError(), diags.Errors())
	require.True(t, set.ElementType(t.Context()).Equal(types.ObjectType{AttrTypes: objectModel{}.AttributeTypes()}))
}
