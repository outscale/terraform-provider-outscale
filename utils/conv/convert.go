package conv

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ToString[T ~string | *string](v T) types.String {
	switch v := any(v).(type) {
	case string:
		return types.StringValue(v)
	case *string:
		if v == nil {
			return types.StringNull()
		}
		return types.StringValue(*v)
	default:
		return types.StringNull()
	}
}

func ToInt64[T ~int | ~int64 | *int | *int64](v T) types.Int64 {
	switch v := any(v).(type) {
	case int:
		return types.Int64Value(int64(v))
	case int64:
		return types.Int64Value(v)
	case *int:
		if v == nil {
			return types.Int64Null()
		}
		return types.Int64Value(int64(*v))
	case *int64:
		if v == nil {
			return types.Int64Null()
		}
		return types.Int64Value(*v)
	default:
		return types.Int64Null()
	}
}

func ToBool[T ~bool | *bool](v T) types.Bool {
	switch v := any(v).(type) {
	case bool:
		return types.BoolValue(v)
	case *bool:
		if v == nil {
			return types.BoolNull()
		}
		return types.BoolValue(*v)
	default:
		return types.BoolNull()
	}
}

func ToFloat64[T ~float64 | *float64](v T) types.Float64 {
	switch v := any(v).(type) {
	case float64:
		return types.Float64Value(v)
	case *float64:
		if v == nil {
			return types.Float64Null()
		}
		return types.Float64Value(*v)
	default:
		return types.Float64Null()
	}
}

func ToRFC3339[T time.Time | *time.Time](v T) timetypes.RFC3339 {
	switch v := any(v).(type) {
	case time.Time:
		return timetypes.NewRFC3339TimeValue(v)
	case *time.Time:
		if v == nil {
			return timetypes.NewRFC3339Null()
		}
		return timetypes.NewRFC3339TimeValue(*v)
	default:
		return timetypes.NewRFC3339Null()
	}
}

func ToSlice[T any, C types.List | types.Set](ctx context.Context, v C) ([]T, diag.Diagnostics) {
	var diags diag.Diagnostics

	switch collection := any(v).(type) {
	case types.List:
		if collection.IsNull() || collection.IsUnknown() {
			return nil, diags
		}
		var result []T
		diags.Append(collection.ElementsAs(ctx, &result, false)...)
		return result, diags
	case types.Set:
		if collection.IsNull() || collection.IsUnknown() {
			return nil, diags
		}
		var result []T
		diags.Append(collection.ElementsAs(ctx, &result, false)...)
		return result, diags
	default:
		return nil, diags
	}
}
