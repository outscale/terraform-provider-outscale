package from

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/osc-sdk-go/v3/pkg/iso8601"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/samber/lo"
)

func ISO8601[T iso8601.Time | *iso8601.Time | time.Time | *time.Time](v T) string {
	switch v := any(v).(type) {
	case iso8601.Time:
		return v.String()
	case *iso8601.Time:
		if v == nil {
			return ""
		}
		return v.String()
	case time.Time:
		return iso8601.Time{Time: v}.String()
	case *time.Time:
		if v == nil {
			return ""
		}
		return iso8601.Time{Time: *v}.String()
	default:
		panic(fmt.Sprintf("unsupported type %T", v))
	}
}

func RFC3339(v types.String) (*time.Time, error) {
	if !fwhelpers.IsSet(v) {
		return nil, nil
	}
	expires, err := time.Parse(time.RFC3339, v.ValueString())
	if err != nil {
		return nil, err
	}

	return &expires, nil
}

func Map[T any](ctx context.Context, v types.Map) (map[string]T, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v.IsNull() || v.IsUnknown() {
		return nil, diags
	}
	var result map[string]T
	diags.Append(v.ElementsAs(ctx, &result, false)...)

	return result, diags
}

func Diag(diags diag.Diagnostics) error {
	errs := lo.Map(diags.Errors(), func(d diag.Diagnostic, _ int) error {
		if d.Detail() != "" {
			return fmt.Errorf("%s: %s", d.Summary(), d.Detail())
		}
		return errors.New(d.Summary())
	})

	return errors.Join(errs...)
}
