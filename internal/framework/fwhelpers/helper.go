package fwhelpers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/outscale/terraform-provider-outscale/internal/framework/fwtypes"

	"github.com/samber/lo"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func WaitForResource[T any](ctx context.Context, conf *retry.StateChangeConf) (*T, error) {
	respRaw, err := conf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}

	resp, ok := respRaw.(*T)
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", respRaw)
	}

	return resp, nil
}

func CheckDiags[T *resource.CreateResponse | *resource.UpdateResponse | *resource.DeleteResponse | *resource.ReadResponse | *resource.ModifyPlanResponse | *resource.ImportStateResponse | *datasource.ReadResponse | *resource.ValidateConfigResponse | *ephemeral.OpenResponse | *provider.ConfigureResponse](resp T, diags diag.Diagnostics) bool {
	switch r := any(resp).(type) {
	case *resource.DeleteResponse:
		r.Diagnostics.Append(diags...)
		return r.Diagnostics.HasError()
	case *resource.ReadResponse:
		r.Diagnostics.Append(diags...)
		return r.Diagnostics.HasError()
	case *resource.UpdateResponse:
		r.Diagnostics.Append(diags...)
		return r.Diagnostics.HasError()
	case *resource.CreateResponse:
		r.Diagnostics.Append(diags...)
		return r.Diagnostics.HasError()
	case *resource.ModifyPlanResponse:
		r.Diagnostics.Append(diags...)
		return r.Diagnostics.HasError()
	case *resource.ImportStateResponse:
		r.Diagnostics.Append(diags...)
		return r.Diagnostics.HasError()
	case *datasource.ReadResponse:
		r.Diagnostics.Append(diags...)
		return r.Diagnostics.HasError()
	case *datasource.ValidateConfigResponse:
		r.Diagnostics.Append(diags...)
		return r.Diagnostics.HasError()
	case *provider.ConfigureResponse:
		r.Diagnostics.Append(diags...)
		return r.Diagnostics.HasError()
	default:
		return true
	}
}

func IsSet(v attr.Value) bool {
	return !v.IsNull() && !v.IsUnknown()
}

func HasChange(planValue, stateValue attr.Value) bool {
	return IsSet(planValue) && !planValue.Equal(stateValue)
}

func GetAttrTypes(model any) map[string]attr.Type {
	attrTypes := make(map[string]attr.Type)

	v := reflect.TypeOf(model)

	for field := range v.Fields() {
		tfsdkTag := field.Tag.Get("tfsdk")
		if tfsdkTag == "" {
			continue
		}

		switch field.Type {
		case reflect.TypeFor[types.String]():
			attrTypes[tfsdkTag] = types.StringType
		case reflect.TypeFor[types.Bool]():
			attrTypes[tfsdkTag] = types.BoolType
		case reflect.TypeFor[types.Int64]():
			attrTypes[tfsdkTag] = types.Int64Type
		case reflect.TypeFor[types.Float64]():
			attrTypes[tfsdkTag] = types.Float64Type
		case reflect.TypeFor[types.Int32]():
			attrTypes[tfsdkTag] = types.Int32Type
		case reflect.TypeFor[fwtypes.CaseInsensitiveStringValue]():
			attrTypes[tfsdkTag] = fwtypes.CaseInsensitiveStringType{}
		default:
			panic(fmt.Sprintf("unhandled field type: %v for field: %s", field.Type, field.Name))
		}
	}

	return attrTypes
}

func GetSliceFromFwtypeSet(ctx context.Context, dataTypeSet types.Set) ([]string, diag.Diagnostics) {
	sliceAttribute := []string{}
	diags := dataTypeSet.ElementsAs(ctx, &sliceAttribute, false)
	if diags.HasError() {
		return sliceAttribute, diags
	}
	return sliceAttribute, diags
}

func GetSlicesFromTypesSetForUpdating(ctx context.Context, stateTypeSet, planTypeSet types.Set) ([]string, []string, diag.Diagnostics) {
	var toAdd, toRemove []string
	diags := planTypeSet.ElementsAs(ctx, &toAdd, false)
	if diags.HasError() {
		return toAdd, toRemove, diags
	}
	diags = stateTypeSet.ElementsAs(ctx, &toRemove, false)
	if diags.HasError() {
		return toAdd, toRemove, diags
	}

	toAdd, toRemove = lo.Difference(toAdd, toRemove)
	return toAdd, toRemove, diags
}

func Difference[C types.List | types.Set](ctx context.Context, oldCollection, newCollection C) (C, C, diag.Diagnostics) {
	var diags diag.Diagnostics
	var zero C

	switch oldCollection := any(oldCollection).(type) {
	case types.List:
		newCollection, ok := any(newCollection).(types.List)
		if !ok {
			return zero, zero, diags
		}
		if !oldCollection.ElementType(ctx).Equal(newCollection.ElementType(ctx)) {
			return zero, zero, diags
		}

		toCreate, toRemove := diffValues(oldCollection.Elements(), newCollection.Elements())
		createList, d := types.ListValue(newCollection.ElementType(ctx), toCreate)
		diags.Append(d...)
		removeList, d := types.ListValue(oldCollection.ElementType(ctx), toRemove)
		diags.Append(d...)

		return any(createList).(C), any(removeList).(C), diags
	case types.Set:
		newCollection, ok := any(newCollection).(types.Set)
		if !ok {
			return zero, zero, diags
		}
		if !oldCollection.ElementType(ctx).Equal(newCollection.ElementType(ctx)) {
			return zero, zero, diags
		}

		toCreate, toRemove := diffValues(oldCollection.Elements(), newCollection.Elements())
		createSet, d := types.SetValue(newCollection.ElementType(ctx), toCreate)
		diags.Append(d...)
		removeSet, d := types.SetValue(oldCollection.ElementType(ctx), toRemove)
		diags.Append(d...)

		return any(createSet).(C), any(removeSet).(C), diags
	default:
		panic(fmt.Sprintf("unsupported type %T", oldCollection))
	}
}

func diffValues(oldValues, newValues []attr.Value) ([]attr.Value, []attr.Value) {
	contains := func(values []attr.Value, target attr.Value) bool {
		return lo.ContainsBy(values, func(value attr.Value) bool {
			return value.Equal(target)
		})
	}

	toCreate := lo.Filter(newValues, func(newValue attr.Value, _ int) bool {
		return !contains(oldValues, newValue)
	})

	toRemove := lo.Filter(oldValues, func(oldValue attr.Value, _ int) bool {
		return !contains(newValues, oldValue)
	})

	return toCreate, toRemove
}
