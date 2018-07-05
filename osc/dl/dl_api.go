package dl

import (
	"context"
	"net/http"

	"github.com/terraform-providers/terraform-provider-outscale/osc"
)

//Operations defines all the operations needed for FCU VMs
type Operations struct {
	client *osc.Client
}

//Service all the necessary actions for them VM service
type Service interface {
	CreateConnection(input *CreateConnectionInput) (*Connection, error)
	DescribeConnections(input *DescribeConnectionsInput) (*Connections, error)
	DeleteConnection(input *DeleteConnectionInput) (*Connection, error)
	CreatePrivateVirtualInterface(input *CreatePrivateVirtualInterfaceInput) (*VirtualInterface, error)
	DescribeVirtualInterfaces(input *DescribeVirtualInterfacesInput) (*DescribeVirtualInterfacesOutput, error)
	DeleteVirtualInterface(input *DeleteVirtualInterfaceInput) (*DeleteVirtualInterfaceOutput, error)
}

// CreateConnection ...
func (v Operations) CreateConnection(input *CreateConnectionInput) (*Connection, error) {
	inURL := "/"
	endpoint := "CreateConnection"
	output := &Connection{}

	if input == nil {
		input = &CreateConnectionInput{}
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

// DescribeConnections ...
func (v Operations) DescribeConnections(input *DescribeConnectionsInput) (*Connections, error) {
	inURL := "/"
	endpoint := "DescribeConnections"
	output := &Connections{}

	if input == nil {
		input = &DescribeConnectionsInput{}
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

// DeleteConnection ...
func (v Operations) DeleteConnection(input *DeleteConnectionInput) (*Connection, error) {
	inURL := "/"
	endpoint := "DeleteConnection"
	output := &Connection{}

	if input == nil {
		input = &DeleteConnectionInput{}
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

// CreatePrivateVirtualInterface ...
func (v Operations) CreatePrivateVirtualInterface(input *CreatePrivateVirtualInterfaceInput) (*VirtualInterface, error) {
	inURL := "/"
	endpoint := "CreatePrivateVirtualInterface"
	output := &VirtualInterface{}

	if input == nil {
		input = &CreatePrivateVirtualInterfaceInput{}
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

// DescribeVirtualInterfaces ...
func (v Operations) DescribeVirtualInterfaces(input *DescribeVirtualInterfacesInput) (*DescribeVirtualInterfacesOutput, error) {
	inURL := "/"
	endpoint := "DescribeVirtualInterfaces"
	output := &DescribeVirtualInterfacesOutput{}

	if input == nil {
		input = &DescribeVirtualInterfacesInput{}
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

// DeleteVirtualInterface ...
func (v Operations) DeleteVirtualInterface(input *DeleteVirtualInterfaceInput) (*DeleteVirtualInterfaceOutput, error) {
	inURL := "/"
	endpoint := "DeleteVirtualInterface"
	output := &DeleteVirtualInterfaceOutput{}

	if input == nil {
		input = &DeleteVirtualInterfaceInput{}
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
