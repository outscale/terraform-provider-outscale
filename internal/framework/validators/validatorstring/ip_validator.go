package validatorstring

import (
	"context"
	"net"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = ipValidator{}

type ipValidator struct{}

func (v ipValidator) Description(_ context.Context) string {
	return "Value must be a valid IP"
}

func (v ipValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v ipValidator) ValidateString(ctx context.Context, req validator.StringRequest, response *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if net.ParseIP(req.ConfigValue.ValueString()) == nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			req.ConfigValue.ValueString(),
		))
		return
	}
}

func IsIP() validator.String {
	return ipValidator{}
}
