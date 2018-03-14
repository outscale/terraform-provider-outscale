package fcu

import (
	"context"
	"net/http"
)

// CreateLinInternetGateway method
func (v VMOperations) CreateInternetGateway(input *CreateInternetGatewayInput) (*CreateInternetGatewayOutput, error) {
	inURL := "/"
	endpoint := "CreateInternetGateway"
	output := &CreateInternetGatewayOutput{}

	if input == nil {
		input = &CreateInternetGatewayInput{}
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

func (v VMOperations) DescribeInternetGateways(input *DescribeInternetGatewaysInput) (*DescribeInternetGatewaysOutput, error) {
	inURL := "/"
	endpoint := "DescribeInternetGateways"
	output := &DescribeInternetGatewaysOutput{}

	if input == nil {
		input = &DescribeInternetGatewaysInput{}
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

func (v VMOperations) DeleteInternetGateway(input *DeleteInternetGatewayInput) (*DeleteInternetGatewayOutput, error) {
	inURL := "/"
	endpoint := "DeleteInternetGateway"
	output := &DeleteInternetGatewayOutput{}

	if input == nil {
		input = &DeleteInternetGatewayInput{}
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
func (v VMOperations) CreateVpc(input *CreateVpcInput) (*CreateVpcOutput, error) {
	inURL := "/"
	endpoint := "CreateVpc"
	output := &CreateVpcOutput{}

	if input == nil {
		input = &CreateVpcInput{}
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
func (v VMOperations) DescribeVpcs(input *DescribeVpcsInput) (*DescribeVpcsOutput, error) {
	inURL := "/"
	endpoint := "DescribeVpcs"
	output := &DescribeVpcsOutput{}

	if input == nil {
		input = &DescribeVpcsInput{}
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
