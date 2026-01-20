package validatorstring

import (
	"context"
	"net"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = cidrValidator{}

type cidrValidator struct{}

func (v cidrValidator) Description(_ context.Context) string {
	return "Value must be a valid CIDR"
}

func (v cidrValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v cidrValidator) ValidateString(ctx context.Context, req validator.StringRequest, response *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if _, _, err := net.ParseCIDR(req.ConfigValue.ValueString()); err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			req.ConfigValue.ValueString(),
		))
		return
	}
}

func IsCIDR() validator.String {
	return cidrValidator{}
}
