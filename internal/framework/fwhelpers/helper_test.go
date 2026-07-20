package fwhelpers_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
)

type testListenerModel struct {
	BackendPort          types.Int64  `tfsdk:"backend_port"`
	BackendProtocol      types.String `tfsdk:"backend_protocol"`
	LoadBalancerPort     types.Int64  `tfsdk:"load_balancer_port"`
	LoadBalancerProtocol types.String `tfsdk:"load_balancer_protocol"`
	ServerCertificateId  types.String `tfsdk:"server_certificate_id"`
	PolicyNames          types.List   `tfsdk:"policy_names"`
}

func TestNestedSetDifference(t *testing.T) {
	t.Parallel()

	attrTypes := map[string]attr.Type{
		"backend_port":           types.Int64Type,
		"backend_protocol":       types.StringType,
		"load_balancer_port":     types.Int64Type,
		"load_balancer_protocol": types.StringType,
		"server_certificate_id":  types.StringType,
		"policy_names":           types.ListType{ElemType: types.StringType},
	}

	emptyPolicies := types.ListValueMust(types.StringType, []attr.Value{})
	oldSet, diags := types.SetValueFrom(t.Context(), types.ObjectType{AttrTypes: attrTypes}, []testListenerModel{{
		BackendPort:          types.Int64Value(80),
		BackendProtocol:      types.StringValue("TCP"),
		LoadBalancerPort:     types.Int64Value(80),
		LoadBalancerProtocol: types.StringValue("TCP"),
		ServerCertificateId:  types.StringValue(""),
		PolicyNames:          emptyPolicies,
	}})
	if diags.HasError() {
		t.Fatalf("building old set: %v", diags)
	}

	newSet, diags := types.SetValueFrom(t.Context(), types.ObjectType{AttrTypes: attrTypes}, []testListenerModel{
		{
			BackendPort:          types.Int64Value(80),
			BackendProtocol:      types.StringValue("TCP"),
			LoadBalancerPort:     types.Int64Value(80),
			LoadBalancerProtocol: types.StringValue("TCP"),
			ServerCertificateId:  types.StringValue(""),
			PolicyNames:          emptyPolicies,
		},
		{
			BackendPort:          types.Int64Value(90),
			BackendProtocol:      types.StringValue("TCP"),
			LoadBalancerPort:     types.Int64Value(90),
			LoadBalancerProtocol: types.StringValue("TCP"),
			ServerCertificateId:  types.StringValue(""),
			PolicyNames:          emptyPolicies,
		},
	})
	if diags.HasError() {
		t.Fatalf("building new set: %v", diags)
	}

	toCreate, toRemove, diags := fwhelpers.Difference(t.Context(), oldSet, newSet)
	if diags.HasError() {
		t.Fatalf("diffing sets: %v", diags)
	}

	if len(toRemove.Elements()) != 0 {
		t.Fatalf("expected no elements to remove, got %d", len(toRemove.Elements()))
	}
	if len(toCreate.Elements()) != 1 {
		t.Fatalf("expected one element to create, got %d", len(toCreate.Elements()))
	}
}
