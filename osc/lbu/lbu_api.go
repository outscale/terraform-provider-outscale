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
}

// CreateLoadBalancer ...
func (v Operations) CreateLoadBalancer(input *CreateLoadBalancerInput) (*CreateLoadBalancerOutput, error) {
	inURL := "/"
	endpoint := "CreateLoadBalancer"
	output := &CreateLoadBalancerOutput{}

	if input == nil {
		input = &CreateLoadBalancerInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodPost, inURL, input)

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

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodPost, inURL, input)

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

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodPost, inURL, input)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}
