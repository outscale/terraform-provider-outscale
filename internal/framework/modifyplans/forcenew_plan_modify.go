package modifyplans

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var _ planmodifier.String = forceNewPlanModify{}

type forceNewPlanModify struct{}

func (f forceNewPlanModify) Description(_ context.Context) string {
	return "Changing the parameter will destroy and recreate the resource."
}

func (f forceNewPlanModify) MarkdownDescription(ctx context.Context) string {
	return f.Description(ctx)
}

// Only force new the attribute when State and Config values differ.
// This is necessary when RequiresReplace is combined with Computed and Optional,
// see https://github.com/hashicorp/terraform-plugin-framework/issues/187
func (f forceNewPlanModify) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	resp.RequiresReplace = false

	// Do nothing if there is no state value or an unknown configuration value.
	if req.StateValue.IsNull() || req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	if req.StateValue.ValueString() != req.ConfigValue.ValueString() {
		resp.RequiresReplace = true
	}
}

func ForceNewFramework() planmodifier.String {
	return forceNewPlanModify{}
}
