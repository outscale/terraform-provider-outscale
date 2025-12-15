package fwmodifyplan

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/outscale/osc-sdk-go/v3/pkg/iso8601"
)

var _ planmodifier.String = datePlanModify{}

type datePlanModify struct{}

func (m datePlanModify) Description(_ context.Context) string {
	return "Invalid 'expiration_date' value"
}

func (m datePlanModify) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m datePlanModify) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Do nothing if there is no state value.
	if req.StateValue.IsNull() {
		return
	}

	// If config is null or unknown, keep the current state value.
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		resp.PlanValue = req.StateValue
		return
	}

	configDate, err := iso8601.Parse([]byte(req.ConfigValue.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			m.Description(ctx),
			"Unable to parse configuration expiration date value: "+err.Error(),
		)
	}
	stateDate, err := iso8601.Parse([]byte(req.StateValue.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			m.Description(ctx),
			"Unable to parse state expiration date value: "+err.Error(),
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// if configDate.Equal(stateDate) {
	// 	resp.PlanValue = req.StateValue
	// } else if configDate.Before(stateDate) {
	if configDate.Before(stateDate) {
		resp.Diagnostics.AddError(
			m.Description(ctx),
			"The new expiration_date should be after the old one."+
				" If the new expiration_date has been update outside of terraform plugin,"+
				" copy the expiration_date state value in your terraform configration file.",
		)
	}
}

func CheckExpirationDate() planmodifier.String {
	return datePlanModify{}
}
