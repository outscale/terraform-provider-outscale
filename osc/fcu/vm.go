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
	StopInstances(input *StopInstancesInput) (*StopInstancesOutput, error)
	StartInstances(input *StartInstancesInput) (*StartInstancesOutput, error)
	ImportKeyPair(input *ImportKeyPairInput) (*ImportKeyPairOutput, error)
	DescribeKeyPairs(input *DescribeKeyPairsInput) (*DescribeKeyPairsOutput, error)
	DeleteKeyPairs(input *DeleteKeyPairInput) (*DeleteKeyPairOutput, error)
	CreateKeyPair(input *CreateKeyPairInput) (*CreateKeyPairOutput, error)
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
