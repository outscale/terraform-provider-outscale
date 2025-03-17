package fwvalidators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/outscale/terraform-provider-outscale/utils"
)

var _ validator.String = dateValidator{}

type dateValidator struct{}

func (v dateValidator) Description(_ context.Context) string {
	return "Invalid date format "
}

func (v dateValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v dateValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {

	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	if err := utils.CheckDateFormat(request.ConfigValue.ValueString()); err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx)+err.Error(),
			request.ConfigValue.ValueString(),
		))
		return
	}

}

func DateValidator() validator.String {
	return dateValidator{}
}
