package validatorlist

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.List = alsoRequiresIfValidator{}

type alsoRequiresIfValidator struct {
	pathExpression path.Expression
	condition      func(context.Context, validator.ListRequest) bool
}

func (v alsoRequiresIfValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("also requires %s if condition is met", v.pathExpression)
}

func (v alsoRequiresIfValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("also requires %s if condition is met", v.pathExpression)
}

func (v alsoRequiresIfValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if !v.condition(ctx, req) {
		return
	}

	alsoRequiresValidator := listvalidator.AlsoRequires(v.pathExpression)
	alsoRequiresValidator.ValidateList(ctx, req, resp)
}

func AlsoRequiresIf(pathExpression path.Expression, condition func(context.Context, validator.ListRequest) bool) validator.List {
	return alsoRequiresIfValidator{
		pathExpression: pathExpression,
		condition:      condition,
	}
}
