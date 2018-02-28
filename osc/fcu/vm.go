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
	inURL := "/?Action=ModifyInstanceKeypair"
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
