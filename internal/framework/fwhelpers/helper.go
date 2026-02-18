package fwhelpers

import (
	"context"
	"fmt"
	"reflect"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwtypes"

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

	for i := range v.NumField() {
		field := v.Field(i)
		tfsdkTag := field.Tag.Get("tfsdk")
		if tfsdkTag == "" {
			continue
		}

		switch field.Type {
		case reflect.TypeOf(types.String{}):
			attrTypes[tfsdkTag] = types.StringType
		case reflect.TypeOf(types.Bool{}):
			attrTypes[tfsdkTag] = types.BoolType
		case reflect.TypeOf(types.Int64{}):
			attrTypes[tfsdkTag] = types.Int64Type
		case reflect.TypeOf(types.Float64{}):
			attrTypes[tfsdkTag] = types.Float64Type
		case reflect.TypeOf(types.Int32{}):
			attrTypes[tfsdkTag] = types.Int32Type
		case reflect.TypeOf(fwtypes.CaseInsensitiveStringValue{}):
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

	setIdsToAdd := mapset.NewSet[string]()
	setIdsToRemove := mapset.NewSet[string]()
	setIdsToAdd.Append(toAdd...)
	setIdsToRemove.Append(toRemove...)
	toAdd = setIdsToAdd.Difference(setIdsToRemove).ToSlice()
	toRemove = setIdsToRemove.Difference(setIdsToAdd).ToSlice()
	return toAdd, toRemove, diags
}
