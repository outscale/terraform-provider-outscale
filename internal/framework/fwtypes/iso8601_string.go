package fwtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/outscale/osc-sdk-go/v3/pkg/iso8601"
)

var (
	_ basetypes.StringTypable  = ISO8601Type{}
	_ basetypes.StringValuable = ISO8601Value{}
)

func ISO8601String(val string) ISO8601Value {
	return ISO8601Value{
		StringValue: types.StringValue(val),
	}
}

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
	return "ISO8601StringType"
}

func (t ISO8601Type) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := ISO8601Value{
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
	return ISO8601Value{}
}

type ISO8601Value struct {
	basetypes.StringValue
}

func (v ISO8601Value) Equal(o attr.Value) bool {
	other, ok := o.(ISO8601Value)
	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v ISO8601Value) Type(ctx context.Context) attr.Type {
	return ISO8601Type{}
}

func (v ISO8601Value) ValidateAttribute(ctx context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	// Allow empty string to support removing expiration dates
	if v.ValueString() == "" {
		return
	}

	if _, err := iso8601.ParseString(v.ValueString()); err != nil {
		resp.Diagnostics.Append(diag.WithPath(req.Path, iso8601InvalidStringDiagnostic(v.ValueString(), err)))

		return
	}
}

func (v ISO8601Value) ValidateParameter(ctx context.Context, req function.ValidateParameterRequest, resp *function.ValidateParameterResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	// Allow empty string to support removing expiration dates
	if v.ValueString() == "" {
		return
	}

	if _, err := iso8601.ParseString(v.ValueString()); err != nil {
		resp.Error = function.NewArgumentFuncError(
			req.Position,
			"Invalid ISO8601 String Value: "+
				"A string value was provided that is not valid ISO8601 string format.\n\n"+
				"Given Value: "+v.ValueString()+"\n"+
				"Error: "+err.Error(),
		)

		return
	}
}

func (v ISO8601Value) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(ISO8601Value)
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

	// Treat empty string and null as semantically equal
	// This allows removing expiration dates by setting expiration_date = ""
	vIsEmpty := v.IsNull() || v.ValueString() == ""
	newIsEmpty := newValue.IsNull() || newValue.ValueString() == ""

	if vIsEmpty && newIsEmpty {
		return true, diags
	}

	// If only one is empty, they are not equal
	if vIsEmpty != newIsEmpty {
		return false, diags
	}

	oldTime, err := iso8601.ParseString(v.ValueString())
	if err != nil {
		diags.AddError(
			"ISO8601 Parse Error",
			"Failed to parse old value as ISO8601.\n\n"+
				"Value: "+v.ValueString()+"\n"+
				"Error: "+err.Error(),
		)
		return false, diags
	}

	newTime, err := iso8601.ParseString(newValue.ValueString())
	if err != nil {
		diags.AddError(
			"ISO8601 Parse Error",
			"Failed to parse new value as ISO8601.\n\n"+
				"Value: "+newValue.ValueString()+"\n"+
				"Error: "+err.Error(),
		)
		return false, diags
	}

	// Compare the underlying time.Time values instead of the iso8601.Time.Equal()
	// method, because iso8601.Time.Equal() checks both the timestamp and the format
	return oldTime.Time.Equal(newTime.Time), diags
}

func (v ISO8601Value) ValueISO8601Time() (iso8601.Time, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v.IsNull() {
		diags.Append(diag.NewErrorDiagnostic("ISO8601 Error", "ISO8601 string value is null"))
		return iso8601.Time{}, diags
	}

	if v.IsUnknown() {
		diags.Append(diag.NewErrorDiagnostic("ISO8601 Error", "ISO8601 string value is unknown"))
		return iso8601.Time{}, diags
	}

	iso8601Time, err := iso8601.ParseString(v.ValueString())
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic("ISO8601 Error", err.Error()))
		return iso8601.Time{}, diags
	}

	return iso8601Time, nil
}

func NewISO8601Null() ISO8601Value {
	return ISO8601Value{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewISO8601Unknown() ISO8601Value {
	return ISO8601Value{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewISO8601TimeValue(value iso8601.Time) ISO8601Value {
	return ISO8601Value{
		StringValue: basetypes.NewStringValue(value.String()),
	}
}

func NewISO8601TimePointerValue(value *iso8601.Time) ISO8601Value {
	if value == nil {
		return NewISO8601Null()
	}

	return ISO8601Value{
		StringValue: basetypes.NewStringValue(value.String()),
	}
}

func NewISO8601Value(value string) (ISO8601Value, diag.Diagnostics) {
	_, err := iso8601.ParseString(value)
	if err != nil {
		return NewISO8601Unknown(), diag.Diagnostics{iso8601InvalidStringDiagnostic(value, err)}
	}

	return ISO8601Value{
		StringValue: basetypes.NewStringValue(value),
	}, nil
}

func NewISO8601PointerValue(value *string) (ISO8601Value, diag.Diagnostics) {
	if value == nil {
		return NewISO8601Null(), nil
	}

	return NewISO8601Value(*value)
}

func iso8601InvalidStringDiagnostic(value string, err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Invalid ISO8601 String Value",
		"A string value was provided that is not valid ISO8601 string format.\n\n"+
			"Given Value: "+value+"\n"+
			"Error: "+err.Error(),
	)
}
