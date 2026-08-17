package validatorstring

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = jsonValidator{}

type jsonValidator struct{}

func (v jsonValidator) Description(_ context.Context) string {
	return "Value must be a valid JSON Document"
}

func (v jsonValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v jsonValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if !json.Valid([]byte(value)) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid JSON Document",
			fmt.Sprintf("%q is not a valid JSON Document", req.ConfigValue.ValueString()),
		)
	}
}

func IsJSON() validator.String {
	return jsonValidator{}
}
