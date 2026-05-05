package validatorstring

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = regexValidator{}

type regexValidator struct{}

func (v regexValidator) Description(_ context.Context) string {
	return "Value must be a valid regular expression"
}

func (v regexValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v regexValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if _, err := regexp.Compile(value); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid regular expression",
			fmt.Sprintf("%q is not a valid regular expression: %s", value, err),
		)
	}
}

func IsRegex() validator.String {
	return regexValidator{}
}
