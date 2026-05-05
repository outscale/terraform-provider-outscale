package validatorstring

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"k8s.io/client-go/util/jsonpath"
)

var _ validator.String = jsonPathValidator{}

type jsonPathValidator struct {
	toExpression func(string) string
}

func (v jsonPathValidator) Description(_ context.Context) string {
	return "Value must be a valid JSONPath expression"
}

func (v jsonPathValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v jsonPathValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if v.toExpression != nil {
		value = v.toExpression(value)
	}

	jp := jsonpath.New("validate").AllowMissingKeys(true)
	if err := jp.Parse(value); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid JSONPath expression",
			fmt.Sprintf("%q is not a valid JSONPath expression: %s", req.ConfigValue.ValueString(), err),
		)
	}
}

func IsJSONPath(toExpression func(string) string) validator.String {
	return jsonPathValidator{toExpression: toExpression}
}
