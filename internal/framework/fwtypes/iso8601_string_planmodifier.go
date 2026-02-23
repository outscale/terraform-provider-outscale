package fwtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// EmptyStringAsNull returns a plan modifier that ensures empty string and null are treated consistently
func EmptyStringAsNull() planmodifier.String {
	return emptyStringAsNullModifier{}
}

type emptyStringAsNullModifier struct{}

func (m emptyStringAsNullModifier) Description(ctx context.Context) string {
	return "Treats empty string and null as semantically equivalent"
}

func (m emptyStringAsNullModifier) MarkdownDescription(ctx context.Context) string {
	return "Treats empty string and null as semantically equivalent"
}

func (m emptyStringAsNullModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If config is empty string, keep it as empty string in the plan.
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() && req.ConfigValue.ValueString() == "" {
		resp.PlanValue = types.StringValue("")
		return
	}

	// If config is null, preserve the state to avoid forcing changes
	if req.ConfigValue.IsNull() && !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
		resp.PlanValue = req.StateValue
		return
	}
}
