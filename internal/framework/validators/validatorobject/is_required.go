package validatorobject

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Object = isRequiredValidator{}

type isRequiredValidator struct{}

func (isRequiredValidator) Description(_ context.Context) string {
	return "must have a configuration value"
}

func (v isRequiredValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (isRequiredValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if !req.ConfigValue.IsNull() {
		return
	}

	name := req.Path.String()
	lastStep, _ := req.Path.Steps().LastStep()
	if attributeName, ok := lastStep.(path.PathStepAttributeName); ok {
		name = string(attributeName)
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Missing required argument",
		fmt.Sprintf("The argument %q is required, but no definition was found.", name),
	)
}

// IsRequired returns a validator that requires an object configuration value.
func IsRequired() validator.Object {
	return isRequiredValidator{}
}
