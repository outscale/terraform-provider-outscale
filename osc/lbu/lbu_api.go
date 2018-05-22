package lbu

import (
	"context"
	"net/http"

	"github.com/terraform-providers/terraform-provider-outscale/osc"
)

//Operations defines all the operations needed for LBU
type Operations struct {
	client *osc.Client
}

//Service all the necessary actions for them LBU service
type Service interface {
	CreateLoadBalancer(input *CreateLoadBalancerInput) (*CreateLoadBalancerOutput, error)
	DescribeLoadBalancers(input *DescribeLoadBalancersInput) (*DescribeLoadBalancersOutput, error)
	DescribeLoadBalancerAttributes(input *DescribeLoadBalancerAttributesInput) (*DescribeLoadBalancerAttributesOutput, error)
	DeleteLoadBalancerListeners(input *DeleteLoadBalancerListenersInput) (*DeleteLoadBalancerListenersOutput, error)
	CreateLoadBalancerListeners(input *CreateLoadBalancerListenersInput) (*CreateLoadBalancerListenersOutput, error)
	ConfigureHealthCheck(input *ConfigureHealthCheckInput) (*ConfigureHealthCheckOutput, error)
	ApplySecurityGroupsToLoadBalancer(input *ApplySecurityGroupsToLoadBalancerInput) (*ApplySecurityGroupsToLoadBalancerOutput, error)
	EnableAvailabilityZonesForLoadBalancer(input *EnableAvailabilityZonesForLoadBalancerInput) (*EnableAvailabilityZonesForLoadBalancerOutput, error)
	DisableAvailabilityZonesForLoadBalancer(input *DisableAvailabilityZonesForLoadBalancerInput) (*DisableAvailabilityZonesForLoadBalancerOutput, error)
	AttachLoadBalancerToSubnets(input *AttachLoadBalancerToSubnetsInput) (*AttachLoadBalancerToSubnetsOutput, error)
	DeleteLoadBalancer(input *DeleteLoadBalancerInput) (*DeleteLoadBalancerOutput, error)
	RegisterInstancesWithLoadBalancer(input *RegisterInstancesWithLoadBalancerInput) (*RegisterInstancesWithLoadBalancerOutput, error)
	DeregisterInstancesFromLoadBalancer(input *DeregisterInstancesFromLoadBalancerInput) (*DeregisterInstancesFromLoadBalancerOutput, error)
	DetachLoadBalancerFromSubnets(input *DetachLoadBalancerFromSubnetsInput) (*DetachLoadBalancerFromSubnetsOutput, error)
	CreateLBCookieStickinessPolicy(input *CreateLBCookieStickinessPolicyInput) (*CreateLBCookieStickinessPolicyOutput, error)
	CreateAppCookieStickinessPolicy(input *CreateAppCookieStickinessPolicyInput) (*CreateAppCookieStickinessPolicyOutput, error)
	SetLoadBalancerPoliciesOfListener(input *SetLoadBalancerPoliciesOfListenerInput) (*SetLoadBalancerPoliciesOfListenerOutput, error)
	DescribeLoadBalancerPolicies(input *DescribeLoadBalancerPoliciesInput) (*DescribeLoadBalancerPoliciesOutput, error)
	DeleteLoadBalancerPolicy(input *DeleteLoadBalancerPolicyInput) (*DeleteLoadBalancerPolicyOutput, error)
}

// CreateLoadBalancer ...
func (v Operations) CreateLoadBalancer(input *CreateLoadBalancerInput) (*CreateLoadBalancerOutput, error) {
	inURL := "/"
	endpoint := "CreateLoadBalancer"
	output := &CreateLoadBalancerOutput{}

	if input == nil {
		input = &CreateLoadBalancerInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DescribeLoadBalancers ...
func (v Operations) DescribeLoadBalancers(input *DescribeLoadBalancersInput) (*DescribeLoadBalancersOutput, error) {
	inURL := "/"
	endpoint := "DescribeLoadBalancers"
	output := &DescribeLoadBalancersOutput{}

	if input == nil {
		input = &DescribeLoadBalancersInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DescribeLoadBalancerAttributes ...
func (v Operations) DescribeLoadBalancerAttributes(input *DescribeLoadBalancerAttributesInput) (*DescribeLoadBalancerAttributesOutput, error) {
	inURL := "/"
	endpoint := "DescribeLoadBalancerAttributes"
	output := &DescribeLoadBalancerAttributesOutput{}

	if input == nil {
		input = &DescribeLoadBalancerAttributesInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DeleteLoadBalancerListeners ...
func (v Operations) DeleteLoadBalancerListeners(input *DeleteLoadBalancerListenersInput) (*DeleteLoadBalancerListenersOutput, error) {
	inURL := "/"
	endpoint := "DeleteLoadBalancerListeners"
	output := &DeleteLoadBalancerListenersOutput{}

	if input == nil {
		input = &DeleteLoadBalancerListenersInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// CreateLoadBalancerListeners ...
func (v Operations) CreateLoadBalancerListeners(input *CreateLoadBalancerListenersInput) (*CreateLoadBalancerListenersOutput, error) {
	inURL := "/"
	endpoint := "CreateLoadBalancerListeners"
	output := &CreateLoadBalancerListenersOutput{}

	if input == nil {
		input = &CreateLoadBalancerListenersInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// ConfigureHealthCheck ...
func (v Operations) ConfigureHealthCheck(input *ConfigureHealthCheckInput) (*ConfigureHealthCheckOutput, error) {
	inURL := "/"
	endpoint := "ConfigureHealthCheck"
	output := &ConfigureHealthCheckOutput{}

	if input == nil {
		input = &ConfigureHealthCheckInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// ApplySecurityGroupsToLoadBalancer ...
func (v Operations) ApplySecurityGroupsToLoadBalancer(input *ApplySecurityGroupsToLoadBalancerInput) (*ApplySecurityGroupsToLoadBalancerOutput, error) {
	inURL := "/"
	endpoint := "ApplySecurityGroupsToLoadBalancer"
	output := &ApplySecurityGroupsToLoadBalancerOutput{}

	if input == nil {
		input = &ApplySecurityGroupsToLoadBalancerInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// EnableAvailabilityZonesForLoadBalancer ...
func (v Operations) EnableAvailabilityZonesForLoadBalancer(input *EnableAvailabilityZonesForLoadBalancerInput) (*EnableAvailabilityZonesForLoadBalancerOutput, error) {
	inURL := "/"
	endpoint := "EnableAvailabilityZonesForLoadBalancer"
	output := &EnableAvailabilityZonesForLoadBalancerOutput{}

	if input == nil {
		input = &EnableAvailabilityZonesForLoadBalancerInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DisableAvailabilityZonesForLoadBalancer ...
func (v Operations) DisableAvailabilityZonesForLoadBalancer(input *DisableAvailabilityZonesForLoadBalancerInput) (*DisableAvailabilityZonesForLoadBalancerOutput, error) {
	inURL := "/"
	endpoint := "DisableAvailabilityZonesForLoadBalancer"
	output := &DisableAvailabilityZonesForLoadBalancerOutput{}

	if input == nil {
		input = &DisableAvailabilityZonesForLoadBalancerInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// AttachLoadBalancerToSubnets ...
func (v Operations) AttachLoadBalancerToSubnets(input *AttachLoadBalancerToSubnetsInput) (*AttachLoadBalancerToSubnetsOutput, error) {
	inURL := "/"
	endpoint := "AttachLoadBalancerToSubnets"
	output := &AttachLoadBalancerToSubnetsOutput{}

	if input == nil {
		input = &AttachLoadBalancerToSubnetsInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DeleteLoadBalancer ...
func (v Operations) DeleteLoadBalancer(input *DeleteLoadBalancerInput) (*DeleteLoadBalancerOutput, error) {
	inURL := "/"
	endpoint := "DeleteLoadBalancer"
	output := &DeleteLoadBalancerOutput{}

	if input == nil {
		input = &DeleteLoadBalancerInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// RegisterInstancesWithLoadBalancer ...
func (v Operations) RegisterInstancesWithLoadBalancer(input *RegisterInstancesWithLoadBalancerInput) (*RegisterInstancesWithLoadBalancerOutput, error) {
	inURL := "/"
	endpoint := "RegisterInstancesWithLoadBalancer"
	output := &RegisterInstancesWithLoadBalancerOutput{}

	if input == nil {
		input = &RegisterInstancesWithLoadBalancerInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DeregisterInstancesFromLoadBalancer ...
func (v Operations) DeregisterInstancesFromLoadBalancer(input *DeregisterInstancesFromLoadBalancerInput) (*DeregisterInstancesFromLoadBalancerOutput, error) {
	inURL := "/"
	endpoint := "DeregisterInstancesFromLoadBalancer"
	output := &DeregisterInstancesFromLoadBalancerOutput{}

	if input == nil {
		input = &DeregisterInstancesFromLoadBalancerInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DetachLoadBalancerFromSubnets ...
func (v Operations) DetachLoadBalancerFromSubnets(input *DetachLoadBalancerFromSubnetsInput) (*DetachLoadBalancerFromSubnetsOutput, error) {
	inURL := "/"
	endpoint := "DetachLoadBalancerFromSubnets"
	output := &DetachLoadBalancerFromSubnetsOutput{}

	if input == nil {
		input = &DetachLoadBalancerFromSubnetsInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// CreateLBCookieStickinessPolicy ...
func (v Operations) CreateLBCookieStickinessPolicy(input *CreateLBCookieStickinessPolicyInput) (*CreateLBCookieStickinessPolicyOutput, error) {
	inURL := "/"
	endpoint := "CreateLBCookieStickinessPolicy"
	output := &CreateLBCookieStickinessPolicyOutput{}

	if input == nil {
		input = &CreateLBCookieStickinessPolicyInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// CreateAppCookieStickinessPolicy ...
func (v Operations) CreateAppCookieStickinessPolicy(input *CreateAppCookieStickinessPolicyInput) (*CreateAppCookieStickinessPolicyOutput, error) {
	inURL := "/"
	endpoint := "CreateAppCookieStickinessPolicy"
	output := &CreateAppCookieStickinessPolicyOutput{}

	if input == nil {
		input = &CreateAppCookieStickinessPolicyInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// SetLoadBalancerPoliciesOfListener ...
func (v Operations) SetLoadBalancerPoliciesOfListener(input *SetLoadBalancerPoliciesOfListenerInput) (*SetLoadBalancerPoliciesOfListenerOutput, error) {
	inURL := "/"
	endpoint := "SetLoadBalancerPoliciesOfListener"
	output := &SetLoadBalancerPoliciesOfListenerOutput{}

	if input == nil {
		input = &SetLoadBalancerPoliciesOfListenerInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DescribeLoadBalancerPolicies ...
func (v Operations) DescribeLoadBalancerPolicies(input *DescribeLoadBalancerPoliciesInput) (*DescribeLoadBalancerPoliciesOutput, error) {
	inURL := "/"
	endpoint := "DescribeLoadBalancerPolicies"
	output := &DescribeLoadBalancerPoliciesOutput{}

	if input == nil {
		input = &DescribeLoadBalancerPoliciesInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DeleteLoadBalancerPolicy ...
func (v Operations) DeleteLoadBalancerPolicy(input *DeleteLoadBalancerPolicyInput) (*DeleteLoadBalancerPolicyOutput, error) {
	inURL := "/"
	endpoint := "DeleteLoadBalancerPolicy"
	output := &DeleteLoadBalancerPolicyOutput{}

	if input == nil {
		input = &DeleteLoadBalancerPolicyInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}
