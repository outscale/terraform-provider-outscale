package fwmodifyplan

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/nav-inc/datetime"
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
	// Do nothing if there is no state value or an unknown configuration value.
	if req.StateValue.IsNull() || req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	newExpirDate, err := datetime.Parse(req.ConfigValue.ValueString(), time.UTC)
	if err != nil {
		resp.Diagnostics.AddError(
			m.Description(ctx),
			"Unable to parse configuration expiration date value: "+err.Error(),
		)
	}
	oldExpirDate, _ := datetime.Parse(req.StateValue.ValueString(), time.UTC)
	if err != nil {
		resp.Diagnostics.AddError(
			m.Description(ctx),
			"Unable to parse state expiration date value: "+err.Error(),
		)
	}
	if newExpirDate.Equal(oldExpirDate) && req.ConfigValue.ValueString() != req.StateValue.ValueString() {
		resp.Diagnostics.AddError(
			m.Description(ctx),
			"The new expiration_date should be after the old one."+
				" If the new expiration_date has been update outside of terraform plugin,"+
				" copy the expiration_date state value in your terraform configration file.",
		)
	}
}

func CkeckExpirationDate() planmodifier.String {
	return datePlanModify{}
}
