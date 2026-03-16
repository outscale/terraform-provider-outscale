package modifyplans

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/osc-sdk-go/v3/pkg/iso8601"
)

func TestFwDatemodifyplan(t *testing.T) {
	t.Parallel()

	oldDate, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		t.Errorf("%v", err.Error())
	}
	newDate, err := iso8601.ParseString(oldDate.AddDate(0, 1, 10).Format(time.RFC3339))
	if err != nil {
		t.Errorf("%v", err.Error())
	}

	currentDate := oldDate.Format(time.RFC3339)
	cases := map[string]struct {
		ConfigValue   types.String
		StateValue    types.String
		ExpectedError bool
	}{
		"valide_date_updating": {
			ConfigValue:   types.StringValue(newDate.String()),
			StateValue:    types.StringValue(currentDate),
			ExpectedError: false,
		},
		"valid_date_plan": {
			ConfigValue:   types.StringValue(currentDate),
			StateValue:    types.StringValue(currentDate),
			ExpectedError: false,
		},
		"valid_date_unknown_values": {
			ConfigValue:   types.StringUnknown(),
			StateValue:    types.StringUnknown(),
			ExpectedError: false,
		},
		"valid_date_configValue": {
			ConfigValue:   types.StringValue(newDate.String()),
			StateValue:    types.StringNull(),
			ExpectedError: false,
		},
		"valid_date_unset_Values": {
			ConfigValue:   types.StringNull(),
			StateValue:    types.StringNull(),
			ExpectedError: false,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			req := planmodifier.StringRequest{
				ConfigValue: tc.ConfigValue,
				StateValue:  tc.StateValue,
			}
			resp := planmodifier.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}

			CheckExpirationDate().PlanModifyString(context.Background(), req, &resp)
			if !tc.ExpectedError && resp.Diagnostics.HasError() {
				t.Errorf("got unexpected error: %s", resp.Diagnostics.Errors())
			}
			if tc.ExpectedError && !resp.Diagnostics.HasError() {
				t.Error("expected error, got none")
			}
		})
	}
}
