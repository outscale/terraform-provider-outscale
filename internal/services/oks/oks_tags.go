package oks

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type OKSTagsModel struct {
	Tags types.Map `tfsdk:"tags"`
}

func OKSTagsSchema() schema.MapAttribute {
	return schema.MapAttribute{
		Computed:    true,
		Optional:    true,
		Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
		ElementType: types.StringType,
	}
}

func expandOKSTags(ctx context.Context, data OKSTagsModel) (map[string]string, diag.Diagnostics) {
	var tags map[string]string
	diags := data.Tags.ElementsAs(ctx, &tags, false)
	if diags.HasError() {
		diags.AddError("Tags conversion error", "Unable to convert Tags into the SDK Model.")
		return tags, diags
	}
	return tags, nil
}

func flattenOKSTags(ctx context.Context, oksTags any) (basetypes.MapValue, diag.Diagnostics) {
	tags, diags := types.MapValueFrom(ctx, types.StringType, oksTags)
	if diags.HasError() {
		diags.AddError("Tags conversion error", "Unable to convert Tags into the Schema Model.")
		return tags, diags
	}
	return tags, nil
}

func cmpOKSTags(ctx context.Context, plan, state OKSTagsModel) (map[string]string, diag.Diagnostics) {
	if plan.Tags.Equal(state.Tags) {
		return nil, nil
	}

	var tags map[string]string
	diags := plan.Tags.ElementsAs(ctx, &tags, false)

	return tags, diags
}
