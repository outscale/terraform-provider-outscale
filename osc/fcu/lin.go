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

func (v VMOperations) DeleteVpc(input *DeleteVpcInput) (*DeleteVpcOutput, error) {
	inURL := "/"
	endpoint := "DeleteVpc"
	output := &DeleteVpcOutput{}

	if input == nil {
		input = &DeleteVpcInput{}
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

func (v VMOperations) AttachInternetGateway(input *AttachInternetGatewayInput) (*AttachInternetGatewayOutput, error) {
	inURL := "/"
	endpoint := "AttachInternetGateway"
	output := &AttachInternetGatewayOutput{}

	if input == nil {
		input = &AttachInternetGatewayInput{}
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

func (v VMOperations) DetachInternetGateway(input *DetachInternetGatewayInput) (*DetachInternetGatewayOutput, error) {
	inURL := "/"
	endpoint := "DetachInternetGateway"
	output := &DetachInternetGatewayOutput{}

	if input == nil {
		input = &DetachInternetGatewayInput{}
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

func (v VMOperations) ModifyVpcAttribute(input *ModifyVpcAttributeInput) (*ModifyVpcAttributeOutput, error) {
	inURL := "/"
	endpoint := "ModifyVpcAttribute"
	output := &ModifyVpcAttributeOutput{}

	if input == nil {
		input = &ModifyVpcAttributeInput{}
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

func (v VMOperations) DescribeVpcAttribute(input *DescribeVpcAttributeInput) (*DescribeVpcAttributeOutput, error) {
	inURL := "/"
	endpoint := "DescribeVpcAttribute"
	output := &DescribeVpcAttributeOutput{}

	if input == nil {
		input = &DescribeVpcAttributeInput{}
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

func (v VMOperations) CreateVpnConnection(input *CreateVpnConnectionInput) (*CreateVpnConnectionOutput, error) {
	inURL := "/"
	endpoint := "CreateVpnConnection"
	output := &CreateVpnConnectionOutput{}

	if input == nil {
		input = &CreateVpnConnectionInput{}
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

func (v VMOperations) CreateVpnGateway(input *CreateVpnGatewayInput) (*CreateVpnGatewayOutput, error) {
	inURL := "/"
	endpoint := "CreateVpnGateway"
	output := &CreateVpnGatewayOutput{}

	if input == nil {
		input = &CreateVpnGatewayInput{}
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

func (v VMOperations) DescribeVpnConnections(input *DescribeVpnConnectionsInput) (*DescribeVpnConnectionsOutput, error) {
	inURL := "/"
	endpoint := "DescribeVpnConnections"
	output := &DescribeVpnConnectionsOutput{}

	if input == nil {
		input = &DescribeVpnConnectionsInput{}
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
func (v VMOperations) DescribeVpnGateways(input *DescribeVpnGatewaysInput) (*DescribeVpnGatewaysOutput, error) {
	inURL := "/"
	endpoint := "DescribeVpnGateways"
	output := &DescribeVpnGatewaysOutput{}

	if input == nil {
		input = &DescribeVpnGatewaysInput{}
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

func (v VMOperations) DeleteVpnConnection(input *DeleteVpnConnectionInput) (*DeleteVpnConnectionOutput, error) {
	inURL := "/"
	endpoint := "DeleteVpnConnection"
	output := &DeleteVpnConnectionOutput{}

	if input == nil {
		input = &DeleteVpnConnectionInput{}
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
func (v VMOperations) DeleteVpnGateway(input *DeleteVpnGatewayInput) (*DeleteVpnGatewayOutput, error) {
	inURL := "/"
	endpoint := "DeleteVpnGateway"
	output := &DeleteVpnGatewayOutput{}

	if input == nil {
		input = &DeleteVpnGatewayInput{}
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

func (v VMOperations) AttachVpnGateway(input *AttachVpnGatewayInput) (*AttachVpnGatewayOutput, error) {
	inURL := "/"
	endpoint := "AttachVpnGateway"
	output := &AttachVpnGatewayOutput{}

	if input == nil {
		input = &AttachVpnGatewayInput{}
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

func (v VMOperations) DetachVpnGateway(input *DetachVpnGatewayInput) (*DetachVpnGatewayOutput, error) {
	inURL := "/"
	endpoint := "DetachVpnGateway"
	output := &DetachVpnGatewayOutput{}

	if input == nil {
		input = &DetachVpnGatewayInput{}
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

func (v VMOperations) CreateVpnConnectionRoute(input *CreateVpnConnectionRouteInput) (*CreateVpnConnectionRouteOutput, error) {
	inURL := "/"
	endpoint := "CreateVpnConnectionRoute"
	output := &CreateVpnConnectionRouteOutput{}

	if input == nil {
		input = &CreateVpnConnectionRouteInput{}
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

func (v VMOperations) DeleteVpnConnectionRoute(input *DeleteVpnConnectionRouteInput) (*DeleteVpnConnectionRouteOutput, error) {
	inURL := "/"
	endpoint := "DeleteVpnConnectionRoute"
	output := &DeleteVpnConnectionRouteOutput{}

	if input == nil {
		input = &DeleteVpnConnectionRouteInput{}
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
