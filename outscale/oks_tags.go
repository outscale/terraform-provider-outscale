package outscale

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

func expandOKSTags(ctx context.Context, data OKSTagsModel) (tags map[string]string, diags diag.Diagnostics) {
	diags = data.Tags.ElementsAs(ctx, &tags, false)
	if diags.HasError() {
		diags.AddError("Tags conversion error", "Unable to convert Tags into the SDK Model.")
		return tags, diags
	}
	return tags, nil
}

func flattenOKSTags(ctx context.Context, oksTags any) (tags basetypes.MapValue, diags diag.Diagnostics) {
	tags, diags = types.MapValueFrom(ctx, types.StringType, oksTags)
	if diags.HasError() {
		diags.AddError("Tags conversion error", "Unable to convert Tags into the Schema Model.")
		return tags, diags
	}
	return tags, nil
}

func cmpOKSTags(ctx context.Context, plan, state OKSTagsModel) (_ map[string]string, diags diag.Diagnostics) {
	if plan.Tags.Equal(state.Tags) {
		return nil, diags
	}

	var tags map[string]string
	diags.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)

	return tags, diags
}
