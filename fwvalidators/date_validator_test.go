package fwvalidators

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nav-inc/datetime"
)

func TestFwDatevalidators(t *testing.T) {
	t.Parallel()

	oldDate, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		t.Errorf("Unable to parse date: %v ", err.Error())
	}
	newDate, err := datetime.Parse(oldDate.AddDate(0, 1, 10).Format(time.RFC3339), time.UTC)
	if err != nil {
		t.Errorf("Unable to parse date: %v ", err.Error())
	}

	currentDate := oldDate.Format(time.RFC3339)
	updatetime := newDate.Format(time.RFC3339)
	cases := map[string]struct {
		ConfigValue   types.String
		StateValue    types.String
		ExpectedError bool
	}{
		"valide_date_updating": {
			ConfigValue:   types.StringValue(updatetime),
			ExpectedError: false,
		},
		"invalid_date_current_date": {
			ConfigValue:   types.StringValue(currentDate),
			ExpectedError: true,
		},
		"valid_date_unknown_values": {
			ConfigValue:   types.StringUnknown(),
			ExpectedError: false,
		},
		"valid_date_configValue": {
			ConfigValue:   types.StringValue(updatetime),
			ExpectedError: false,
		},
		"valid_date_unset_Values": {
			ConfigValue:   types.StringNull(),
			ExpectedError: false,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: tc.ConfigValue,
			}

			resp := validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			DateValidator().ValidateString(context.Background(), req, &resp)
			if !tc.ExpectedError && resp.Diagnostics.HasError() {
				t.Errorf("got unexpected error: %s", resp.Diagnostics.Errors())
			}
			if tc.ExpectedError && !resp.Diagnostics.HasError() {
				t.Error("expected error, got none")
			}
		})
	}
}
