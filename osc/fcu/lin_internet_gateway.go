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
