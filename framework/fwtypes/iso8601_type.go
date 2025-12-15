package fwtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementation satisfies the expected interfaces
var _ basetypes.StringTypable = ISO8601Type{}

type ISO8601Type struct {
	basetypes.StringType
}

func (t ISO8601Type) Equal(o attr.Type) bool {
	other, ok := o.(ISO8601Type)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t ISO8601Type) String() string {
	return "ISO8601Type"
}

func (t ISO8601Type) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := ISO8601{
		StringValue: in,
	}

	return value, nil
}

func (t ISO8601Type) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

func (t ISO8601Type) ValueType(ctx context.Context) attr.Value {
	return ISO8601{}
}
