package fcu

import (
	"context"
	"net/http"
)

func (v VMOperations) CreateRoute(input *CreateRouteInput) (*CreateRouteOutput, error) {
	inURL := "/"
	endpoint := "CreateRoute"
	output := &CreateRouteOutput{}

	if input == nil {
		input = &CreateRouteInput{}
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

func (v VMOperations) ReplaceRoute(input *ReplaceRouteInput) (*ReplaceRouteOutput, error) {
	inURL := "/"
	endpoint := "ReplaceRoute"
	output := &ReplaceRouteOutput{}

	if input == nil {
		input = &ReplaceRouteInput{}
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

func (v VMOperations) DeleteRoute(input *DeleteRouteInput) (*DeleteRouteOutput, error) {
	inURL := "/"
	endpoint := "DeleteRoute"
	output := &DeleteRouteOutput{}

	if input == nil {
		input = &DeleteRouteInput{}
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

func (v VMOperations) DescribeRouteTables(input *DescribeRouteTablesInput) (*DescribeRouteTablesOutput, error) {
	inURL := "/"
	endpoint := "DescribeRouteTables"
	output := &DescribeRouteTablesOutput{}

	if input == nil {
		input = &DescribeRouteTablesInput{}
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
