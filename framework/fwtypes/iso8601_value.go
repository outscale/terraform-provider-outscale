package fwtypes

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/outscale/osc-sdk-go/v3/pkg/iso8601"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ basetypes.StringValuable                   = ISO8601{}
	_ basetypes.StringValuableWithSemanticEquals = ISO8601{}
)

type ISO8601 struct {
	basetypes.StringValue
}

func (v ISO8601) Equal(o attr.Value) bool {
	other, ok := o.(ISO8601)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v ISO8601) Type(ctx context.Context) attr.Type {
	// ISO8601Type defined in the schema type section
	return ISO8601Type{}
}

func (v ISO8601) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	// The framework should always pass the correct value type, but always check
	newValue, ok := newValuable.(ISO8601)

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

	if v.IsNull() || newValue.IsNull() || v.IsUnknown() || newValue.IsUnknown() {
		return false, diags
	}

	// Skipping error checking if ISO8601 already implemented ISO8601 validation
	priorTime, _ := iso8601.Parse([]byte(v.ValueString()))

	// Skipping error checking if ISO8601 already implemented ISO8601 validation
	newTime, _ := iso8601.Parse([]byte(newValue.ValueString()))

	// If the times are equivalent, keep the prior value
	return priorTime.Equal(newTime), diags
}

// ValidateAttribute implements attribute value validation. This type requires the value to be a String value that
// is valid ISO8601 format.
func (v ISO8601) ValidateAttribute(ctx context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	if _, err := iso8601.Parse([]byte(v.ValueString())); err != nil {
		resp.Diagnostics.Append(diag.WithPath(
			req.Path,
			iso8601InvalidStringDiagnostic(v.ValueString(), err),
		))

		return
	}
}

// ValidateParameter implements provider-defined function parameter value validation. This type requires the value to
// be a String value that is valid ISO8601 format.
func (v ISO8601) ValidateParameter(ctx context.Context, req function.ValidateParameterRequest, resp *function.ValidateParameterResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	if _, err := iso8601.Parse([]byte(v.ValueString())); err != nil {
		resp.Error = function.NewArgumentFuncError(
			req.Position,
			"Invalid ISO8601 String Value: "+
				"A string value was provided that is not valid ISO8601 date string format.\n\n"+
				"Given Value: "+v.ValueString()+"\n"+
				"Error: "+err.Error(),
		)

		return
	}
}

// ValueISO8601Time creates a new time.Time instance with the ISO8601 StringValue. A null or unknown value will produce an error diagnostic.
func (v ISO8601) ValueISO8601Time() (time.Time, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v.IsNull() {
		diags.Append(diag.NewErrorDiagnostic("ISO8601 ValueISO8601Time Error", "ISO8601 string value is null"))
		return time.Time{}, diags
	}

	if v.IsUnknown() {
		diags.Append(diag.NewErrorDiagnostic("ISO8601 ValueISO8601Time Error", "ISO8601 string value is unknown"))
		return time.Time{}, diags
	}

	iso8601Time, err := iso8601.Parse([]byte(v.ValueString()))
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic("ISO8601 ValueISO8601Time Error", err.Error()))
		return time.Time{}, diags
	}

	return iso8601Time, nil
}

// NewISO8601Null creates an ISO8601 with a null value. Determine whether the value is null via IsNull method.
func NewISO8601Null() ISO8601 {
	return ISO8601{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewISO8601Unknown creates an ISO8601 with an unknown value. Determine whether the value is unknown via IsUnknown method.
func NewISO8601Unknown() ISO8601 {
	return ISO8601{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewISO8601Value creates an ISO8601 with a known value or raises an error
// diagnostic if the string is not ISO8601 format.
func NewISO8601Value(value string) (ISO8601, diag.Diagnostics) {
	if value == "" {
		return NewISO8601Null(), nil
	}
	_, err := iso8601.Parse([]byte(value))
	if err != nil {
		// Returning an unknown value will guarantee that, as a last resort,
		// Terraform will return an error if attempting to store into state.
		return NewISO8601Unknown(), diag.Diagnostics{iso8601InvalidStringDiagnostic(value, err)}
	}

	return ISO8601{
		StringValue: basetypes.NewStringValue(value),
	}, nil
}

// NewISO8601ValueMust creates an ISO8601 with a known value or raises a panic
// if the string is not ISO8601 format.
//
// This creation function is only recommended to create ISO8601 values which
// either will not potentially affect practitioners, such as testing, or within
// exhaustively tested provider logic.
func NewISO8601ValueMust(value string) ISO8601 {
	if value == "" {
		return NewISO8601Null()
	}
	_, err := iso8601.Parse([]byte(value))
	if err != nil {
		panic(fmt.Sprintf("Invalid ISO8601 String Value (%s): %s", value, err))
	}

	return ISO8601{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewISO8601PointerValue creates an ISO8601 with a null value if nil, a known
// value, or raises an error diagnostic if the string is not ISO8601 format.
func NewISO8601PointerValue(value *string) (ISO8601, diag.Diagnostics) {
	if value == nil || *value == "" {
		return NewISO8601Null(), nil
	}

	return NewISO8601Value(*value)
}

// NewISO8601PointerValueMust creates an ISO8601 with a null value if nil, a
// known value, or raises a panic if the string is not ISO8601 format.
//
// This creation function is only recommended to create ISO8601 values which
// either will not potentially affect practitioners, such as testing, or within
// exhaustively tested provider logic.
func NewISO8601PointerValueMust(value *string) ISO8601 {
	if value == nil || *value == "" {
		return NewISO8601Null()
	}

	return NewISO8601ValueMust(*value)
}

// iso8601InvalidStringDiagnostic returns an error diagnostic intended to report
// when a string is not ISO8601 format.
func iso8601InvalidStringDiagnostic(value string, err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Invalid ISO8601 String Value",
		"A string value was provided that is not valid ISO8601 string format.\n\n"+
			"Given Value: "+value+"\n"+
			"Error: "+err.Error(),
	)
}
