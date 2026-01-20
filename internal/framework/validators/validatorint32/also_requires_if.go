package validatorint32

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Int32 = alsoRequiresIfValidator{}

type alsoRequiresIfValidator struct {
	pathExpression path.Expression
	condition      func(context.Context, validator.Int32Request) bool
}

func (v alsoRequiresIfValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("also requires %s if condition is met", v.pathExpression)
}

func (v alsoRequiresIfValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("also requires %s if condition is met", v.pathExpression)
}

func (v alsoRequiresIfValidator) ValidateInt32(ctx context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	if !v.condition(ctx, req) {
		return
	}

	alsoRequiresValidator := int32validator.AlsoRequires(v.pathExpression)
	alsoRequiresValidator.ValidateInt32(ctx, req, resp)
}

func AlsoRequiresIf(pathExpression path.Expression, condition func(context.Context, validator.Int32Request) bool) validator.Int32 {
	return alsoRequiresIfValidator{
		pathExpression: pathExpression,
		condition:      condition,
	}
}
