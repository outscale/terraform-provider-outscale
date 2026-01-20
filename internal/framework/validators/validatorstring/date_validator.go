package validatorstring

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/outscale/osc-sdk-go/v3/pkg/iso8601"
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

	if err := checkDateFormat(request.ConfigValue.ValueString()); err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx)+err.Error(),
			request.ConfigValue.ValueString(),
		))
		return
	}
}

func checkDateFormat(dateFormat string) error {
	if dateFormat == "" {
		return nil
	}
	currentDate := time.Now()

	settingDate, err := iso8601.Parse([]byte(dateFormat))
	if err != nil {
		return err
	}
	if currentDate.After(settingDate) {
		return fmt.Errorf("expiration date: '%s' should be after current date '%s'", settingDate, currentDate)
	}

	return nil
}

func DateValidator() validator.String {
	return dateValidator{}
}
