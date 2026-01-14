package fwtypes

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable  = CaseInsensitiveStringType{}
	_ basetypes.StringValuable = CaseInsensitiveStringValue{}
)

func CaseInsensitiveString(val string) CaseInsensitiveStringValue {
	return CaseInsensitiveStringValue{
		StringValue: types.StringValue(val),
	}
}

type CaseInsensitiveStringType struct {
	basetypes.StringType
}

func (t CaseInsensitiveStringType) Equal(o attr.Type) bool {
	other, ok := o.(CaseInsensitiveStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t CaseInsensitiveStringType) String() string {
	return "CaseInsensitiveStringType"
}

func (t CaseInsensitiveStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := CaseInsensitiveStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t CaseInsensitiveStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t CaseInsensitiveStringType) ValueType(ctx context.Context) attr.Value {
	return CaseInsensitiveStringValue{}
}

type CaseInsensitiveStringValue struct {
	basetypes.StringValue
}

func (v CaseInsensitiveStringValue) Equal(o attr.Value) bool {
	other, ok := o.(CaseInsensitiveStringValue)
	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v CaseInsensitiveStringValue) Type(ctx context.Context) attr.Type {
	return CaseInsensitiveStringType{}
}

func (v CaseInsensitiveStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(CaseInsensitiveStringValue)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	return strings.EqualFold(newValue.ValueString(), v.ValueString()), diags
}
