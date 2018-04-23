package fcu

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
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

func (v VMOperations) CreateRouteTable(input *CreateRouteTableInput) (*CreateRouteTableOutput, error) {
	inURL := "/"
	endpoint := "CreateRouteTable"
	output := &CreateRouteTableOutput{}

	if input == nil {
		input = &CreateRouteTableInput{}
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

func (v VMOperations) DisableVgwRoutePropagation(input *DisableVgwRoutePropagationInput) (*DisableVgwRoutePropagationOutput, error) {
	inURL := "/"
	endpoint := "DisableVgwRoutePropagation"
	output := &DisableVgwRoutePropagationOutput{}

	if input == nil {
		input = &DisableVgwRoutePropagationInput{}
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

func (v VMOperations) EnableVgwRoutePropagation(input *EnableVgwRoutePropagationInput) (*EnableVgwRoutePropagationOutput, error) {
	inURL := "/"
	endpoint := "EnableVgwRoutePropagation"
	output := &EnableVgwRoutePropagationOutput{}

	if input == nil {
		input = &EnableVgwRoutePropagationInput{}
	}
	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	requestDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (v VMOperations) DisassociateRouteTable(input *DisassociateRouteTableInput) (*DisassociateRouteTableOutput, error) {
	inURL := "/"
	endpoint := "DisassociateRouteTable"
	output := &DisassociateRouteTableOutput{}

	if input == nil {
		input = &DisassociateRouteTableInput{}
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

func (v VMOperations) DeleteRouteTable(input *DeleteRouteTableInput) (*DeleteRouteTableOutput, error) {
	inURL := "/"
	endpoint := "DeleteRouteTable"
	output := &DeleteRouteTableOutput{}

	if input == nil {
		input = &DeleteRouteTableInput{}
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

func (v VMOperations) AssociateRouteTable(input *AssociateRouteTableInput) (*AssociateRouteTableOutput, error) {
	inURL := "/"
	endpoint := "AssociateRouteTable"
	output := &AssociateRouteTableOutput{}

	if input == nil {
		input = &AssociateRouteTableInput{}
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

func (v VMOperations) ReplaceRouteTableAssociation(input *ReplaceRouteTableAssociationInput) (*ReplaceRouteTableAssociationOutput, error) {
	inURL := "/"
	endpoint := "ReplaceRouteTableAssociation"
	output := &ReplaceRouteTableAssociationOutput{}

	if input == nil {
		input = &ReplaceRouteTableAssociationInput{}
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
