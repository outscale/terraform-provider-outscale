package modifyplans

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// This test covers possible cases with Computed,
// Optionnal and the Plan Modification on an attribute
func TestFwForceNewmodifyplan(t *testing.T) {
	t.Parallel()

	s1 := "eu-west-2a"
	s2 := "eu-west-2b"
	cases := map[string]struct {
		ConfigValue     types.String
		StateValue      types.String
		RequiresReplace bool
	}{
		"same_values": {
			ConfigValue:     types.StringValue(s1),
			StateValue:      types.StringValue(s1),
			RequiresReplace: false,
		},
		"different_values": {
			ConfigValue:     types.StringValue(s1),
			StateValue:      types.StringValue(s2),
			RequiresReplace: true,
		},
		"config_null": {
			ConfigValue:     types.StringNull(),
			StateValue:      types.StringValue(s1),
			RequiresReplace: false,
		},
		"config_unknown": {
			ConfigValue:     types.StringUnknown(),
			StateValue:      types.StringValue(s1),
			RequiresReplace: false,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			req := planmodifier.StringRequest{
				ConfigValue: tc.ConfigValue,
				StateValue:  tc.StateValue,
			}
			resp := planmodifier.StringResponse{}

			ForceNewFramework().PlanModifyString(context.Background(), req, &resp)
			if !tc.RequiresReplace && resp.RequiresReplace {
				t.Errorf("got unexpected error: %s", resp.Diagnostics.Errors())
			}
			if tc.RequiresReplace && !resp.RequiresReplace {
				t.Error("expected replace, got none")
			}
		})
	}
}
