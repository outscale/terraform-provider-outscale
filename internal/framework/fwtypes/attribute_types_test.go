package fwtypes_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwtypes"
	"github.com/stretchr/testify/require"
)

func TestObjectAttributeTypesAtPath(t *testing.T) {
	t.Parallel()

	granteeTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	grantTypes := map[string]attr.Type{
		"grantee": types.ObjectType{
			AttrTypes: granteeTypes,
		},
	}
	root := types.ObjectType{AttrTypes: map[string]attr.Type{
		"permissions": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"grants": types.SetType{
					ElemType: types.ObjectType{
						AttrTypes: grantTypes,
					},
				},
			},
		},
	}}

	require.Equal(t, grantTypes, fwtypes.AttributeTypesAtPath(root, "permissions", "grants"))
	require.Equal(t, granteeTypes, fwtypes.AttributeTypesAtPath(root, "permissions", "grants", "grantee"))
}
