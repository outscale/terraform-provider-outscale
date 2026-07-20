package validatorstring

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = alsoRequiresIfValidator{}

type alsoRequiresIfValidator struct {
	pathExpression path.Expression
	condition      func(context.Context, validator.StringRequest) bool
	summary        string
	detail         string
}

func (v alsoRequiresIfValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("also requires %s if condition is met", v.pathExpression)
}

func (v alsoRequiresIfValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("also requires %s if condition is met", v.pathExpression)
}

func (v alsoRequiresIfValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if !v.condition(ctx, req) {
		return
	}

	tmpResp := &validator.StringResponse{}
	alsoRequiresValidator := stringvalidator.AlsoRequires(v.pathExpression)
	alsoRequiresValidator.ValidateString(ctx, req, tmpResp)

	// If the validation failed, we return the custom diagnostic
	if tmpResp.Diagnostics.HasError() {
		resp.Diagnostics.AddAttributeError(req.Path, v.summary, v.detail)
		return
	}
}

func AlsoRequiresIf(pathExpression path.Expression, condition func(context.Context, validator.StringRequest) bool, summary string, detail string) validator.String {
	return alsoRequiresIfValidator{
		pathExpression: pathExpression,
		condition:      condition,
		summary:        summary,
		detail:         detail,
	}
}
