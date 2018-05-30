package fcu

import (
	"context"
	"net/http"

	"github.com/terraform-providers/terraform-provider-outscale/osc"
)

//VMOperations defines all the operations needed for FCU VMs
type VMOperations struct {
	client *osc.Client
}

//VMService all the necessary actions for them VM service
type VMService interface {
	RunInstance(input *RunInstancesInput) (*Reservation, error)
	DescribeInstances(input *DescribeInstancesInput) (*DescribeInstancesOutput, error)
	GetPasswordData(input *GetPasswordDataInput) (*GetPasswordDataOutput, error)
	ModifyInstanceKeyPair(input *ModifyInstanceKeyPairInput) error
	ModifyInstanceAttribute(input *ModifyInstanceAttributeInput) (*ModifyInstanceAttributeOutput, error)
	TerminateInstances(input *TerminateInstancesInput) (*TerminateInstancesOutput, error)
	AllocateAddress(input *AllocateAddressInput) (*AllocateAddressOutput, error)
	DescribeAddressesRequest(input *DescribeAddressesInput) (*DescribeAddressesOutput, error)
	StopInstances(input *StopInstancesInput) (*StopInstancesOutput, error)
	StartInstances(input *StartInstancesInput) (*StartInstancesOutput, error)
	ImportKeyPair(input *ImportKeyPairInput) (*ImportKeyPairOutput, error)
	DescribeKeyPairs(input *DescribeKeyPairsInput) (*DescribeKeyPairsOutput, error)
	DeleteKeyPairs(input *DeleteKeyPairInput) (*DeleteKeyPairOutput, error)
	CreateKeyPair(input *CreateKeyPairInput) (*CreateKeyPairOutput, error)
	AssociateAddress(input *AssociateAddressInput) (*AssociateAddressOutput, error)
	DisassociateAddress(input *DisassociateAddressInput) (*DisassociateAddressOutput, error)
	ReleaseAddress(input *ReleaseAddressInput) (*ReleaseAddressOutput, error)
	RegisterImage(input *RegisterImageInput) (*RegisterImageOutput, error)
	DescribeImages(input *DescribeImagesInput) (*DescribeImagesOutput, error)
	ModifyImageAttribute(input *ModifyImageAttributeInput) (*ModifyImageAttributeOutput, error)
	DeleteTags(input *DeleteTagsInput) (*DeleteTagsOutput, error)
	CreateTags(input *CreateTagsInput) (*CreateTagsOutput, error)
	DeregisterImage(input *DeregisterImageInput) (*DeregisterImageOutput, error)
	DescribeTags(input *DescribeTagsInput) (*DescribeTagsOutput, error)
	CreateSecurityGroup(input *CreateSecurityGroupInput) (*CreateSecurityGroupOutput, error)
	DescribeSecurityGroups(input *DescribeSecurityGroupsInput) (*DescribeSecurityGroupsOutput, error)
	RevokeSecurityGroupEgress(input *RevokeSecurityGroupEgressInput) (*RevokeSecurityGroupEgressOutput, error)
	RevokeSecurityGroupIngress(input *RevokeSecurityGroupIngressInput) (*RevokeSecurityGroupIngressOutput, error)
	AuthorizeSecurityGroupEgress(input *AuthorizeSecurityGroupEgressInput) (*AuthorizeSecurityGroupEgressOutput, error)
	AuthorizeSecurityGroupIngress(input *AuthorizeSecurityGroupIngressInput) (*AuthorizeSecurityGroupIngressOutput, error)
	DeleteSecurityGroup(input *DeleteSecurityGroupInput) (*DeleteSecurityGroupOutput, error)
	CreateVolume(input *CreateVolumeInput) (*Volume, error)
	DeleteVolume(input *DeleteVolumeInput) (*DeleteVolumeOutput, error)
	DescribeVolumes(input *DescribeVolumesInput) (*DescribeVolumesOutput, error)
	AttachVolume(input *AttachVolumeInput) (*VolumeAttachment, error)
	DetachVolume(input *DetachVolumeInput) (*VolumeAttachment, error)
	CreateSubNet(input *CreateSubnetInput) (*CreateSubnetOutput, error)
	DeleteSubNet(input *DeleteSubnetInput) (*DeleteSubnetOutput, error)
	DescribeSubNet(input *DescribeSubnetsInput) (*DescribeSubnetsOutput, error)
	DescribeInstanceAttribute(input *DescribeInstanceAttributeInput) (*DescribeInstanceAttributeOutput, error)
	DescribeInstanceStatus(input *DescribeInstanceStatusInput) (*DescribeInstanceStatusOutput, error)
	CreateInternetGateway(input *CreateInternetGatewayInput) (*CreateInternetGatewayOutput, error)
	DescribeInternetGateways(input *DescribeInternetGatewaysInput) (*DescribeInternetGatewaysOutput, error)
	DeleteInternetGateway(input *DeleteInternetGatewayInput) (*DeleteInternetGatewayOutput, error)
	CreateNatGateway(input *CreateNatGatewayInput) (*CreateNatGatewayOutput, error)
	DescribeNatGateways(input *DescribeNatGatewaysInput) (*DescribeNatGatewaysOutput, error)
	DeleteNatGateway(input *DeleteNatGatewayInput) (*DeleteNatGatewayOutput, error)
	CreateVpc(input *CreateVpcInput) (*CreateVpcOutput, error)
	DescribeVpcs(input *DescribeVpcsInput) (*DescribeVpcsOutput, error)
	DeleteVpc(input *DeleteVpcInput) (*DeleteVpcOutput, error)
	AttachInternetGateway(input *AttachInternetGatewayInput) (*AttachInternetGatewayOutput, error)
	DetachInternetGateway(input *DetachInternetGatewayInput) (*DetachInternetGatewayOutput, error)
	ModifyVpcAttribute(input *ModifyVpcAttributeInput) (*ModifyVpcAttributeOutput, error)
	DescribeVpcAttribute(input *DescribeVpcAttributeInput) (*DescribeVpcAttributeOutput, error)
	CreateAccessKey(input *CreateAccessKeyInput) (*CreateAccessKeyOutput, error)
	DescribeAccessKey(input *DescribeAccessKeyInput) (*DescribeAccessKeyOutput, error)
	DeleteAccessKey(input *DeleteAccessKeyInput) (*DeleteAccessKeyOutput, error)
	UpdateAccessKey(input *UpdateAccessKeyInput) (*UpdateAccessKeyOutput, error)
	DeleteDhcpOptions(input *DeleteDhcpOptionsInput) (*DeleteDhcpOptionsOutput, error)
	CreateDhcpOptions(input *CreateDhcpOptionsInput) (*CreateDhcpOptionsOutput, error)
	DescribeDhcpOptions(input *DescribeDhcpOptionsInput) (*DescribeDhcpOptionsOutput, error)
	AssociateDhcpOptions(input *AssociateDhcpOptionsInput) (*AssociateDhcpOptionsOutput, error)
	DescribeCustomerGateways(input *DescribeCustomerGatewaysInput) (*DescribeCustomerGatewaysOutput, error)
	DeleteCustomerGateway(input *DeleteCustomerGatewayInput) (*DeleteCustomerGatewayOutput, error)
	CreateCustomerGateway(input *CreateCustomerGatewayInput) (*CreateCustomerGatewayOutput, error)
	CreateRoute(input *CreateRouteInput) (*CreateRouteOutput, error)
	ReplaceRoute(input *ReplaceRouteInput) (*ReplaceRouteOutput, error)
	DeleteRoute(input *DeleteRouteInput) (*DeleteRouteOutput, error)
	DescribeRouteTables(input *DescribeRouteTablesInput) (*DescribeRouteTablesOutput, error)
	CreateRouteTable(input *CreateRouteTableInput) (*CreateRouteTableOutput, error)
	DisableVgwRoutePropagation(input *DisableVgwRoutePropagationInput) (*DisableVgwRoutePropagationOutput, error)
	EnableVgwRoutePropagation(input *EnableVgwRoutePropagationInput) (*EnableVgwRoutePropagationOutput, error)
	DisassociateRouteTable(input *DisassociateRouteTableInput) (*DisassociateRouteTableOutput, error)
	DeleteRouteTable(input *DeleteRouteTableInput) (*DeleteRouteTableOutput, error)
	AssociateRouteTable(input *AssociateRouteTableInput) (*AssociateRouteTableOutput, error)
	ReplaceRouteTableAssociation(input *ReplaceRouteTableAssociationInput) (*ReplaceRouteTableAssociationOutput, error)
	CopyImage(input *CopyImageInput) (*CopyImageOutput, error)
	DescribeSnapshots(input *DescribeSnapshotsInput) (*DescribeSnapshotsOutput, error)
	CreateVpnConnection(input *CreateVpnConnectionInput) (*CreateVpnConnectionOutput, error)
	DescribeVpnConnections(input *DescribeVpnConnectionsInput) (*DescribeVpnConnectionsOutput, error)
	DeleteVpnConnection(input *DeleteVpnConnectionInput) (*DeleteVpnConnectionOutput, error)
	CreateVpnGateway(input *CreateVpnGatewayInput) (*CreateVpnGatewayOutput, error)
	DescribeVpnGateways(input *DescribeVpnGatewaysInput) (*DescribeVpnGatewaysOutput, error)
	DeleteVpnGateway(input *DeleteVpnGatewayInput) (*DeleteVpnGatewayOutput, error)
	AttachVpnGateway(input *AttachVpnGatewayInput) (*AttachVpnGatewayOutput, error)
	DetachVpnGateway(input *DetachVpnGatewayInput) (*DetachVpnGatewayOutput, error)
	CreateImageExportTask(input *CreateImageExportTaskInput) (*CreateImageExportTaskOutput, error)
	DescribeImageExportTasks(input *DescribeImageExportTasksInput) (*DescribeImageExportTasksOutput, error)
	CreateVpnConnectionRoute(input *CreateVpnConnectionRouteInput) (*CreateVpnConnectionRouteOutput, error)
	DeleteVpnConnectionRoute(input *DeleteVpnConnectionRouteInput) (*DeleteVpnConnectionRouteOutput, error)
	DescribeAvailabilityZones(input *DescribeAvailabilityZonesInput) (*DescribeAvailabilityZonesOutput, error)
	DescribePrefixLists(input *DescribePrefixListsInput) (*DescribePrefixListsOutput, error)
	DescribeQuotas(input *DescribeQuotasInput) (*DescribeQuotasOutput, error)
	DescribeRegions(input *DescribeRegionsInput) (*DescribeRegionsOutput, error)
	CreateNetworkInterface(input *CreateNetworkInterfaceInput) (*CreateNetworkInterfaceOutput, error)
	DeleteNetworkInterface(input *DeleteNetworkInterfaceInput) (*DeleteNetworkInterfaceOutput, error)
	DescribeNetworkInterfaces(input *DescribeNetworkInterfacesInput) (*DescribeNetworkInterfacesOutput, error)
	ModifyNetworkInterfaceAttribute(input *ModifyNetworkInterfaceAttributeInput) (*ModifyNetworkInterfaceAttributeOutput, error)
	DescribeNetworkInterfaceAttribute(input *DescribeNetworkInterfaceAttributeInput) (*DescribeNetworkInterfaceAttributeOutput, error)
	DetachNetworkInterface(input *DetachNetworkInterfaceInput) (*DetachNetworkInterfaceOutput, error)
	AttachNetworkInterface(input *AttachNetworkInterfaceInput) (*AttachNetworkInterfaceOutput, error)
	AssignPrivateIpAddresses(input *AssignPrivateIpAddressesInput) (*AssignPrivateIpAddressesOutput, error)
	UnassignPrivateIpAddresses(input *UnassignPrivateIpAddressesInput) (*UnassignPrivateIpAddressesOutput, error)
	CreateSnapshotExportTask(input *CreateSnapshotExportTaskInput) (*CreateSnapshotExportTaskOutput, error)
	DescribeSnapshotExportTasks(input *DescribeSnapshotExportTasksInput) (*DescribeSnapshotExportTasksOutput, error)
	CreateSnapshot(input *CreateSnapshotInput) (*Snapshot, error)
	DeleteSnapshot(input *DeleteSnapshotInput) (*DeleteSnapshotOutput, error)
	DescribeProductTypes(input *DescribeProductTypesInput) (*DescribeProductTypesOutput, error)
	DescribeReservedInstances(input *DescribeReservedInstancesInput) (*DescribeReservedInstancesOutput, error)
	DescribeInstanceTypes(input *DescribeInstanceTypesInput) (*DescribeInstanceTypesOutput, error)
	DescribeReservedInstancesOfferings(input *DescribeReservedInstancesOfferingsInput) (*DescribeReservedInstancesOfferingsOutput, error)
	DescribeImageAttribute(input *DescribeImageAttributeInput) (*DescribeImageAttributeOutput, error)
	CreateVpcPeeringConnection(input *CreateVpcPeeringConnectionInput) (*CreateVpcPeeringConnectionOutput, error)
	DescribeVpcPeeringConnections(input *DescribeVpcPeeringConnectionsInput) (*DescribeVpcPeeringConnectionsOutput, error)
	AcceptVpcPeeringConnection(input *AcceptVpcPeeringConnectionInput) (*AcceptVpcPeeringConnectionOutput, error)
	ModifyVpcPeeringConnectionOptions(input *ModifyVpcPeeringConnectionOptionsInput) (*ModifyVpcPeeringConnectionOptionsOutput, error)
	DeleteVpcPeeringConnection(input *DeleteVpcPeeringConnectionInput) (*DeleteVpcPeeringConnectionOutput, error)
	PurchaseReservedInstancesOffering(input *PurchaseReservedInstancesOfferingInput) (*PurchaseReservedInstancesOfferingOutput, error)
}

const opRunInstances = "RunInstances"

func (v VMOperations) RunInstance(input *RunInstancesInput) (*Reservation, error) {
	req, err := v.client.NewRequest(context.Background(), opRunInstances, http.MethodGet, "/", input)
	if err != nil {
		return nil, err
	}

	output := Reservation{}

	err = v.client.Do(context.Background(), req, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

const opDescribeInstances = "DescribeInstances"

// DescribeInstances method
func (v VMOperations) DescribeInstances(input *DescribeInstancesInput) (*DescribeInstancesOutput, error) {
	inURL := "/"
	endpoint := "DescribeInstances"
	output := &DescribeInstancesOutput{}

	if input == nil {
		input = &DescribeInstancesInput{}
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

// DescribeInstances method
func (v VMOperations) ModifyInstanceKeyPair(input *ModifyInstanceKeyPairInput) error {
	inURL := "/"
	endpoint := "ModifyInstanceKeypair"

	if input == nil {
		input = &ModifyInstanceKeyPairInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodPost, inURL, input)

	if err != nil {
		return err
	}

	err = v.client.Do(context.TODO(), req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (v VMOperations) ModifyInstanceAttribute(input *ModifyInstanceAttributeInput) (*ModifyInstanceAttributeOutput, error) {
	inURL := "/"
	endpoint := "ModifyInstanceAttribute"
	output := &ModifyInstanceAttributeOutput{}

	if input == nil {
		input = &ModifyInstanceAttributeInput{}
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

func (v VMOperations) GetPasswordData(input *GetPasswordDataInput) (*GetPasswordDataOutput, error) {
	inURL := "/"
	endpoint := "GetPasswordData"
	output := &GetPasswordDataOutput{}

	if input == nil {
		input = &GetPasswordDataInput{}
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

// DescribeInstances method
func (v VMOperations) TerminateInstances(input *TerminateInstancesInput) (*TerminateInstancesOutput, error) {
	inURL := "/"
	endpoint := "TerminateInstances"
	output := &TerminateInstancesOutput{}

	if input == nil {
		input = &TerminateInstancesInput{}
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

func (v VMOperations) AllocateAddress(input *AllocateAddressInput) (*AllocateAddressOutput, error) {
	inURL := "/"
	endpoint := "AllocateAddress"
	output := &AllocateAddressOutput{}

	if input == nil {
		input = &AllocateAddressInput{}
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

func (v VMOperations) StopInstances(input *StopInstancesInput) (*StopInstancesOutput, error) {
	inURL := "/"
	endpoint := "StopInstances"
	output := &StopInstancesOutput{}

	if input == nil {
		input = &StopInstancesInput{}
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

//DescribeAddresses
func (v VMOperations) DescribeAddressesRequest(input *DescribeAddressesInput) (*DescribeAddressesOutput, error) {
	inURL := "/"
	endpoint := "DescribeAddresses"
	output := &DescribeAddressesOutput{}

	if input == nil {
		input = &DescribeAddressesInput{}
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

func (v VMOperations) StartInstances(input *StartInstancesInput) (*StartInstancesOutput, error) {
	inURL := "/"
	endpoint := "StartInstances"
	output := &StartInstancesOutput{}

	if input == nil {
		input = &StartInstancesInput{}
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

func (v VMOperations) AssociateAddress(input *AssociateAddressInput) (*AssociateAddressOutput, error) {
	inURL := "/"
	endpoint := "AssociateAddress"
	output := &AssociateAddressOutput{}

	if input == nil {
		input = &AssociateAddressInput{}
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

func (v VMOperations) DisassociateAddress(input *DisassociateAddressInput) (*DisassociateAddressOutput, error) {
	inURL := "/"
	endpoint := "DisassociateAddress"
	output := &DisassociateAddressOutput{}

	if input == nil {
		input = &DisassociateAddressInput{}
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

func (v VMOperations) ReleaseAddress(input *ReleaseAddressInput) (*ReleaseAddressOutput, error) {
	inURL := "/"
	endpoint := "ReleaseAddress"
	output := &ReleaseAddressOutput{}

	if input == nil {
		input = &ReleaseAddressInput{}
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

func (v VMOperations) RegisterImage(input *RegisterImageInput) (*RegisterImageOutput, error) {
	inURL := "/"
	endpoint := "CreateImage"
	output := &RegisterImageOutput{}

	if input == nil {
		input = &RegisterImageInput{}
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

func (v VMOperations) DescribeImages(input *DescribeImagesInput) (*DescribeImagesOutput, error) {
	inURL := "/"
	endpoint := "DescribeImages"
	output := &DescribeImagesOutput{}

	if input == nil {
		input = &DescribeImagesInput{}
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

func (v VMOperations) ModifyImageAttribute(input *ModifyImageAttributeInput) (*ModifyImageAttributeOutput, error) {
	inURL := "/"
	endpoint := "ModifyImageAttribute"
	output := &ModifyImageAttributeOutput{}

	if input == nil {
		input = &ModifyImageAttributeInput{}
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

func (v VMOperations) DeleteTags(input *DeleteTagsInput) (*DeleteTagsOutput, error) {
	inURL := "/"
	endpoint := "DeleteTags"
	output := &DeleteTagsOutput{}

	if input == nil {
		input = &DeleteTagsInput{}
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

func (v VMOperations) CreateTags(input *CreateTagsInput) (*CreateTagsOutput, error) {
	inURL := "/"
	endpoint := "CreateTags"
	output := &CreateTagsOutput{}

	if input == nil {
		input = &CreateTagsInput{}
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

func (v VMOperations) DeregisterImage(input *DeregisterImageInput) (*DeregisterImageOutput, error) {
	inURL := "/"
	endpoint := "DeregisterImage"
	output := &DeregisterImageOutput{}

	if input == nil {
		input = &DeregisterImageInput{}
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

func (v VMOperations) DescribeTags(input *DescribeTagsInput) (*DescribeTagsOutput, error) {
	inURL := "/"
	endpoint := "DescribeTags"
	output := &DescribeTagsOutput{}

	if input == nil {
		input = &DescribeTagsInput{}
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

func (v VMOperations) CreateSecurityGroup(input *CreateSecurityGroupInput) (*CreateSecurityGroupOutput, error) {
	inURL := "/"
	endpoint := "CreateSecurityGroup"
	output := &CreateSecurityGroupOutput{}

	if input == nil {
		input = &CreateSecurityGroupInput{}
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

func (v VMOperations) ImportKeyPair(input *ImportKeyPairInput) (*ImportKeyPairOutput, error) {
	inURL := "/"
	endpoint := "ImportKeyPair"
	output := &ImportKeyPairOutput{}

	if input == nil {
		input = &ImportKeyPairInput{}
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

func (v VMOperations) DescribeSecurityGroups(input *DescribeSecurityGroupsInput) (*DescribeSecurityGroupsOutput, error) {
	inURL := "/"
	endpoint := "DescribeSecurityGroups"
	output := &DescribeSecurityGroupsOutput{}

	if input == nil {
		input = &DescribeSecurityGroupsInput{}

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

func (v VMOperations) DescribeKeyPairs(input *DescribeKeyPairsInput) (*DescribeKeyPairsOutput, error) {
	inURL := "/"
	endpoint := "DescribeKeyPairs"
	output := &DescribeKeyPairsOutput{}

	if input == nil {
		input = &DescribeKeyPairsInput{}
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

func (v VMOperations) RevokeSecurityGroupEgress(input *RevokeSecurityGroupEgressInput) (*RevokeSecurityGroupEgressOutput, error) {
	inURL := "/"
	endpoint := "RevokeSecurityGroupEgress"
	output := &RevokeSecurityGroupEgressOutput{}

	if input == nil {
		input = &RevokeSecurityGroupEgressInput{}
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

func (v VMOperations) RevokeSecurityGroupIngress(input *RevokeSecurityGroupIngressInput) (*RevokeSecurityGroupIngressOutput, error) {
	inURL := "/"
	endpoint := "RevokeSecurityGroupIngress"
	output := &RevokeSecurityGroupIngressOutput{}

	if input == nil {
		input = &RevokeSecurityGroupIngressInput{}
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

func (v VMOperations) AuthorizeSecurityGroupEgress(input *AuthorizeSecurityGroupEgressInput) (*AuthorizeSecurityGroupEgressOutput, error) {
	inURL := "/"
	endpoint := "AuthorizeSecurityGroupEgress"
	output := &AuthorizeSecurityGroupEgressOutput{}

	if input == nil {
		input = &AuthorizeSecurityGroupEgressInput{}
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

func (v VMOperations) DeleteKeyPairs(input *DeleteKeyPairInput) (*DeleteKeyPairOutput, error) {
	inURL := "/"
	endpoint := "DeleteKeyPair"
	output := &DeleteKeyPairOutput{}

	if input == nil {
		input = &DeleteKeyPairInput{}
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

func (v VMOperations) AuthorizeSecurityGroupIngress(input *AuthorizeSecurityGroupIngressInput) (*AuthorizeSecurityGroupIngressOutput, error) {
	inURL := "/"
	endpoint := "AuthorizeSecurityGroupIngress"
	output := &AuthorizeSecurityGroupIngressOutput{}

	if input == nil {
		input = &AuthorizeSecurityGroupIngressInput{}
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

func (v VMOperations) DeleteSecurityGroup(input *DeleteSecurityGroupInput) (*DeleteSecurityGroupOutput, error) {
	inURL := "/"
	endpoint := "DeleteSecurityGroup"
	output := &DeleteSecurityGroupOutput{}

	if input == nil {
		input = &DeleteSecurityGroupInput{}
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

func (v VMOperations) CreateKeyPair(input *CreateKeyPairInput) (*CreateKeyPairOutput, error) {
	inURL := "/"
	endpoint := "CreateKeyPair"
	output := &CreateKeyPairOutput{}

	if input == nil {
		input = &CreateKeyPairInput{}
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

func (v VMOperations) CreateVolume(input *CreateVolumeInput) (*Volume, error) {
	inURL := "/"
	endpoint := "CreateVolume"
	output := &Volume{}

	if input == nil {
		input = &CreateVolumeInput{}
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

func (v VMOperations) DeleteVolume(input *DeleteVolumeInput) (*DeleteVolumeOutput, error) {
	inURL := "/"
	endpoint := "DeleteVolume"
	output := &DeleteVolumeOutput{}

	if input == nil {
		input = &DeleteVolumeInput{}
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

func (v VMOperations) DescribeVolumes(input *DescribeVolumesInput) (*DescribeVolumesOutput, error) {
	inURL := "/"
	endpoint := "DescribeVolumes"
	output := &DescribeVolumesOutput{}

	if input == nil {
		input = &DescribeVolumesInput{}
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

func (v VMOperations) AttachVolume(input *AttachVolumeInput) (*VolumeAttachment, error) {
	inURL := "/"
	endpoint := "AttachVolume"
	output := &VolumeAttachment{}

	if input == nil {
		input = &AttachVolumeInput{}
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

func (v VMOperations) DetachVolume(input *DetachVolumeInput) (*VolumeAttachment, error) {
	inURL := "/"
	endpoint := "DetachVolume"
	output := &VolumeAttachment{}

	if input == nil {
		input = &DetachVolumeInput{}
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

func (v VMOperations) DescribeInstanceAttribute(input *DescribeInstanceAttributeInput) (*DescribeInstanceAttributeOutput, error) {
	inURL := "/"
	endpoint := "DescribeInstanceAttribute"
	output := &DescribeInstanceAttributeOutput{}

	if input == nil {
		input = &DescribeInstanceAttributeInput{}
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
func (v VMOperations) CreateNatGateway(input *CreateNatGatewayInput) (*CreateNatGatewayOutput, error) {
	inURL := "/"
	endpoint := "CreateNatGateway"
	output := &CreateNatGatewayOutput{}

	if input == nil {
		input = &CreateNatGatewayInput{}
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

func (v VMOperations) DescribeNatGateways(input *DescribeNatGatewaysInput) (*DescribeNatGatewaysOutput, error) {
	inURL := "/"
	endpoint := "DescribeNatGateways"
	output := &DescribeNatGatewaysOutput{}

	if input == nil {
		input = &DescribeNatGatewaysInput{}
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

func (v VMOperations) DescribeInstanceStatus(input *DescribeInstanceStatusInput) (*DescribeInstanceStatusOutput, error) {
	inURL := "/"
	endpoint := "DescribeInstanceStatus"
	output := &DescribeInstanceStatusOutput{}

	if input == nil {
		input = &DescribeInstanceStatusInput{}
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
func (v VMOperations) DeleteNatGateway(input *DeleteNatGatewayInput) (*DeleteNatGatewayOutput, error) {
	inURL := "/"
	endpoint := "DeleteNatGateway"
	output := &DeleteNatGatewayOutput{}

	if input == nil {
		input = &DeleteNatGatewayInput{}
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
func (v VMOperations) CreateSubNet(input *CreateSubnetInput) (*CreateSubnetOutput, error) {
	inURL := "/"
	endpoint := "CreateSubnet"
	output := &CreateSubnetOutput{}

	if input == nil {
		input = &CreateSubnetInput{}
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

func (v VMOperations) DeleteSubNet(input *DeleteSubnetInput) (*DeleteSubnetOutput, error) {
	inURL := "/"
	endpoint := "DeleteSubnet"
	output := &DeleteSubnetOutput{}

	if input == nil {
		input = &DeleteSubnetInput{}
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
func (v VMOperations) DescribeSubNet(input *DescribeSubnetsInput) (*DescribeSubnetsOutput, error) {
	inURL := "/"
	endpoint := "DescribeSubnets"
	output := &DescribeSubnetsOutput{}

	if input == nil {
		input = &DescribeSubnetsInput{}
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

func (v VMOperations) CreateAccessKey(input *CreateAccessKeyInput) (*CreateAccessKeyOutput, error) {
	inURL := "/"
	endpoint := "CreateAccessKey"
	output := &CreateAccessKeyOutput{}

	if input == nil {
		input = &CreateAccessKeyInput{}
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
func (v VMOperations) DeleteDhcpOptions(input *DeleteDhcpOptionsInput) (*DeleteDhcpOptionsOutput, error) {
	inURL := "/"
	endpoint := "DescribeDhcpOptions"
	output := &DeleteDhcpOptionsOutput{}

	if input == nil {
		input = &DeleteDhcpOptionsInput{}
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

func (v VMOperations) DescribeCustomerGateways(input *DescribeCustomerGatewaysInput) (*DescribeCustomerGatewaysOutput, error) {
	inURL := "/"
	endpoint := "DescribeCustomerGateways"
	output := &DescribeCustomerGatewaysOutput{}

	if input == nil {
		input = &DescribeCustomerGatewaysInput{}
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

func (v VMOperations) DescribeAccessKey(input *DescribeAccessKeyInput) (*DescribeAccessKeyOutput, error) {
	inURL := "/"
	endpoint := "GetAccessKey"
	output := &DescribeAccessKeyOutput{}

	if input == nil {
		input = &DescribeAccessKeyInput{}
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

func (v VMOperations) CreateDhcpOptions(input *CreateDhcpOptionsInput) (*CreateDhcpOptionsOutput, error) {
	inURL := "/"
	endpoint := "CreateDhcpOptions"
	output := &CreateDhcpOptionsOutput{}

	if input == nil {
		input = &CreateDhcpOptionsInput{}
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

func (v VMOperations) DeleteCustomerGateway(input *DeleteCustomerGatewayInput) (*DeleteCustomerGatewayOutput, error) {
	inURL := "/"
	endpoint := "DeleteCustomerGateway"
	output := &DeleteCustomerGatewayOutput{}

	if input == nil {
		input = &DeleteCustomerGatewayInput{}
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

func (v VMOperations) DeleteAccessKey(input *DeleteAccessKeyInput) (*DeleteAccessKeyOutput, error) {
	inURL := "/"
	endpoint := "DeleteAccessKey"
	output := &DeleteAccessKeyOutput{}

	if input == nil {
		input = &DeleteAccessKeyInput{}
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
func (v VMOperations) UpdateAccessKey(input *UpdateAccessKeyInput) (*UpdateAccessKeyOutput, error) {
	inURL := "/"
	endpoint := "UpdateAccessKey"
	output := &UpdateAccessKeyOutput{}

	if input == nil {
		input = &UpdateAccessKeyInput{}
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

func (v VMOperations) DescribeDhcpOptions(input *DescribeDhcpOptionsInput) (*DescribeDhcpOptionsOutput, error) {
	inURL := "/"
	endpoint := "DescribeDhcpOptions"
	output := &DescribeDhcpOptionsOutput{}

	if input == nil {
		input = &DescribeDhcpOptionsInput{}
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

func (v VMOperations) AssociateDhcpOptions(input *AssociateDhcpOptionsInput) (*AssociateDhcpOptionsOutput, error) {
	inURL := "/"
	endpoint := "AssociateDhcpOptions"
	output := &AssociateDhcpOptionsOutput{}

	if input == nil {
		input = &AssociateDhcpOptionsInput{}
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
func (v VMOperations) CreateCustomerGateway(input *CreateCustomerGatewayInput) (*CreateCustomerGatewayOutput, error) {

	inURL := "/"
	endpoint := "CreateCustomerGateway"
	output := &CreateCustomerGatewayOutput{}

	if input == nil {
		input = &CreateCustomerGatewayInput{}

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

func (v VMOperations) CreateImageExportTask(input *CreateImageExportTaskInput) (*CreateImageExportTaskOutput, error) {
	inURL := "/"
	endpoint := "CreateImageExportTask"
	output := &CreateImageExportTaskOutput{}

	if input == nil {
		input = &CreateImageExportTaskInput{}
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
func (v VMOperations) CopyImage(input *CopyImageInput) (*CopyImageOutput, error) {
	inURL := "/"
	endpoint := "CopyImage"
	output := &CopyImageOutput{}

	if input == nil {
		input = &CopyImageInput{}
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
func (v VMOperations) DescribeAvailabilityZones(input *DescribeAvailabilityZonesInput) (*DescribeAvailabilityZonesOutput, error) {
	inURL := "/"
	endpoint := "DescribeAvailabilityZones"
	output := &DescribeAvailabilityZonesOutput{}

	if input == nil {
		input = &DescribeAvailabilityZonesInput{}
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

func (v VMOperations) DescribeImageExportTasks(input *DescribeImageExportTasksInput) (*DescribeImageExportTasksOutput, error) {
	inURL := "/"
	endpoint := "DescribeImageExportTasks"
	output := &DescribeImageExportTasksOutput{}

	if input == nil {
		input = &DescribeImageExportTasksInput{}
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
func (v VMOperations) DescribeSnapshots(input *DescribeSnapshotsInput) (*DescribeSnapshotsOutput, error) {
	inURL := "/"
	endpoint := "DescribeSnapshots"
	output := &DescribeSnapshotsOutput{}

	if input == nil {
		input = &DescribeSnapshotsInput{}
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
func (v VMOperations) DescribePrefixLists(input *DescribePrefixListsInput) (*DescribePrefixListsOutput, error) {
	inURL := "/"
	endpoint := "DescribePrefixLists"
	output := &DescribePrefixListsOutput{}

	if input == nil {
		input = &DescribePrefixListsInput{}
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

func (v VMOperations) DescribeQuotas(input *DescribeQuotasInput) (*DescribeQuotasOutput, error) {
	inURL := "/"
	endpoint := "DescribeQuotas"
	output := &DescribeQuotasOutput{}

	if input == nil {
		input = &DescribeQuotasInput{}
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
func (v VMOperations) CreateSnapshotExportTask(input *CreateSnapshotExportTaskInput) (*CreateSnapshotExportTaskOutput, error) {
	inURL := "/"
	endpoint := "CreateSnapshotExportTask"
	output := &CreateSnapshotExportTaskOutput{}

	if input == nil {
		input = &CreateSnapshotExportTaskInput{}
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

func (v VMOperations) DescribeRegions(input *DescribeRegionsInput) (*DescribeRegionsOutput, error) {
	inURL := "/"
	endpoint := "DescribeRegions"
	output := &DescribeRegionsOutput{}

	if input == nil {
		input = &DescribeRegionsInput{}
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

func (v VMOperations) CreateNetworkInterface(input *CreateNetworkInterfaceInput) (*CreateNetworkInterfaceOutput, error) {
	inURL := "/"
	endpoint := "CreateNetworkInterface"
	output := &CreateNetworkInterfaceOutput{}

	if input == nil {
		input = &CreateNetworkInterfaceInput{}
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
func (v VMOperations) DescribeSnapshotExportTasks(input *DescribeSnapshotExportTasksInput) (*DescribeSnapshotExportTasksOutput, error) {
	inURL := "/"
	endpoint := "DescribeSnapshotExportTasks"
	output := &DescribeSnapshotExportTasksOutput{}

	if input == nil {
		input = &DescribeSnapshotExportTasksInput{}
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

func (v VMOperations) DeleteNetworkInterface(input *DeleteNetworkInterfaceInput) (*DeleteNetworkInterfaceOutput, error) {
	inURL := "/"
	endpoint := "DeleteNetworkInterface"
	output := &DeleteNetworkInterfaceOutput{}

	if input == nil {
		input = &DeleteNetworkInterfaceInput{}
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
func (v VMOperations) CreateSnapshot(input *CreateSnapshotInput) (*Snapshot, error) {
	inURL := "/"
	endpoint := "CreateSnapshot"
	output := &Snapshot{}

	if input == nil {
		input = &CreateSnapshotInput{}
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
func (v VMOperations) DescribeProductTypes(input *DescribeProductTypesInput) (*DescribeProductTypesOutput, error) {
	inURL := "/"
	endpoint := "DescribeProductTypes"
	output := &DescribeProductTypesOutput{}

	if input == nil {
		input = &DescribeProductTypesInput{}
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

func (v VMOperations) DescribeNetworkInterfaces(input *DescribeNetworkInterfacesInput) (*DescribeNetworkInterfacesOutput, error) {
	inURL := "/"
	endpoint := "DescribeNetworkInterfaces"
	output := &DescribeNetworkInterfacesOutput{}

	if input == nil {
		input = &DescribeNetworkInterfacesInput{}
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
func (v VMOperations) DescribeReservedInstances(input *DescribeReservedInstancesInput) (*DescribeReservedInstancesOutput, error) {
	inURL := "/"
	endpoint := "DescribeReservedInstances"
	output := &DescribeReservedInstancesOutput{}

	if input == nil {
		input = &DescribeReservedInstancesInput{}
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

func (v VMOperations) ModifyNetworkInterfaceAttribute(input *ModifyNetworkInterfaceAttributeInput) (*ModifyNetworkInterfaceAttributeOutput, error) {

	inURL := "/"
	endpoint := "ModifyNetworkInterfaceAttribute"
	output := &ModifyNetworkInterfaceAttributeOutput{}

	if input == nil {
		input = &ModifyNetworkInterfaceAttributeInput{}
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

func (v VMOperations) DescribeInstanceTypes(input *DescribeInstanceTypesInput) (*DescribeInstanceTypesOutput, error) {
	inURL := "/"
	endpoint := "DescribeInstanceTypes"
	output := &DescribeInstanceTypesOutput{}

	if input == nil {
		input = &DescribeInstanceTypesInput{}
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

func (v VMOperations) DeleteSnapshot(input *DeleteSnapshotInput) (*DeleteSnapshotOutput, error) {
	inURL := "/"
	endpoint := "DeleteSnapshot"
	output := &DeleteSnapshotOutput{}

	if input == nil {
		input = &DeleteSnapshotInput{}
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
func (v VMOperations) DescribeReservedInstancesOfferings(input *DescribeReservedInstancesOfferingsInput) (*DescribeReservedInstancesOfferingsOutput, error) {
	inURL := "/"
	endpoint := "DescribeReservedInstancesOfferings"
	output := &DescribeReservedInstancesOfferingsOutput{}

	if input == nil {
		input = &DescribeReservedInstancesOfferingsInput{}
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

func (v VMOperations) DescribeNetworkInterfaceAttribute(input *DescribeNetworkInterfaceAttributeInput) (*DescribeNetworkInterfaceAttributeOutput, error) {
	inURL := "/"
	endpoint := "DescribeNetworkInterfaceAttribute"
	output := &DescribeNetworkInterfaceAttributeOutput{}

	if input == nil {
		input = &DescribeNetworkInterfaceAttributeInput{}
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
func (v VMOperations) DescribeImageAttribute(input *DescribeImageAttributeInput) (*DescribeImageAttributeOutput, error) {
	inURL := "/"
	endpoint := "DescribeImageAttribute"
	output := &DescribeImageAttributeOutput{}

	if input == nil {
		input = &DescribeImageAttributeInput{}
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

func (v VMOperations) DetachNetworkInterface(input *DetachNetworkInterfaceInput) (*DetachNetworkInterfaceOutput, error) {
	inURL := "/"
	endpoint := "DetachNetworkInterface"
	output := &DetachNetworkInterfaceOutput{}

	if input == nil {
		input = &DetachNetworkInterfaceInput{}
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
func (v VMOperations) CreateVpcPeeringConnection(input *CreateVpcPeeringConnectionInput) (*CreateVpcPeeringConnectionOutput, error) {
	inURL := "/"
	endpoint := "CreateVpcPeeringConnection"
	output := &CreateVpcPeeringConnectionOutput{}

	if input == nil {
		input = &CreateVpcPeeringConnectionInput{}
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

func (v VMOperations) AttachNetworkInterface(input *AttachNetworkInterfaceInput) (*AttachNetworkInterfaceOutput, error) {
	inURL := "/"
	endpoint := "AttachNetworkInterface"
	output := &AttachNetworkInterfaceOutput{}

	if input == nil {
		input = &AttachNetworkInterfaceInput{}
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

func (v VMOperations) DescribeVpcPeeringConnections(input *DescribeVpcPeeringConnectionsInput) (*DescribeVpcPeeringConnectionsOutput, error) {
	inURL := "/"
	endpoint := "DescribeVpcPeeringConnections"
	output := &DescribeVpcPeeringConnectionsOutput{}

	if input == nil {
		input = &DescribeVpcPeeringConnectionsInput{}
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

func (v VMOperations) AcceptVpcPeeringConnection(input *AcceptVpcPeeringConnectionInput) (*AcceptVpcPeeringConnectionOutput, error) {
	inURL := "/"
	endpoint := "AcceptVpcPeeringConnection"
	output := &AcceptVpcPeeringConnectionOutput{}

	if input == nil {
		input = &AcceptVpcPeeringConnectionInput{}
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

func (v VMOperations) AssignPrivateIpAddresses(input *AssignPrivateIpAddressesInput) (*AssignPrivateIpAddressesOutput, error) {
	inURL := "/"
	endpoint := "AssignPrivateIpAddresses"
	output := &AssignPrivateIpAddressesOutput{}

	if input == nil {
		input = &AssignPrivateIpAddressesInput{}
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
func (v VMOperations) ModifyVpcPeeringConnectionOptions(input *ModifyVpcPeeringConnectionOptionsInput) (*ModifyVpcPeeringConnectionOptionsOutput, error) {
	inURL := "/"
	endpoint := "ModifyVpcPeeringConnectionOptions"
	output := &ModifyVpcPeeringConnectionOptionsOutput{}

	if input == nil {
		input = &ModifyVpcPeeringConnectionOptionsInput{}
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

func (v VMOperations) UnassignPrivateIpAddresses(input *UnassignPrivateIpAddressesInput) (*UnassignPrivateIpAddressesOutput, error) {
	inURL := "/"
	endpoint := "UnassignPrivateIpAddresses"
	output := &UnassignPrivateIpAddressesOutput{}

	if input == nil {
		input = &UnassignPrivateIpAddressesInput{}
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
func (v VMOperations) DeleteVpcPeeringConnection(input *DeleteVpcPeeringConnectionInput) (*DeleteVpcPeeringConnectionOutput, error) {
	inURL := "/"
	endpoint := "DeleteVpcPeeringConnection"
	output := &DeleteVpcPeeringConnectionOutput{}

	if input == nil {
		input = &DeleteVpcPeeringConnectionInput{}
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

func (v VMOperations) PurchaseReservedInstancesOffering(input *PurchaseReservedInstancesOfferingInput) (*PurchaseReservedInstancesOfferingOutput, error) {
	inURL := "/"
	endpoint := "PurchaseReservedInstancesOffering"
	output := &PurchaseReservedInstancesOfferingOutput{}

	if input == nil {
		input = &PurchaseReservedInstancesOfferingInput{}
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
