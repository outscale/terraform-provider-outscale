package icu

import (
	"context"
	"fmt"
	"net/http"

	"github.com/terraform-providers/terraform-provider-outscale/osc"
)

//VMOperations defines all the operations needed for FCU VMs
type ICUOperations struct {
	client *osc.Client
}

//VMService all the necessary actions for them VM service
type ICUService interface {
	CreateAccessKey(input *CreateAccessKeyInput) (*CreateAccessKeyOutput, error)
	DeleteAccessKey(input *DeleteAccessKeyInput) (*DeleteAccessKeyOutput, error)
	UpdateAccessKey(input *UpdateAccessKeyInput) (*UpdateAccessKeyOutput, error)
	ListAccessKeys(input *ListAccessKeysInput) (*ListAccessKeysOutput, error)
}

func (v ICUOperations) CreateAccessKey(input *CreateAccessKeyInput) (*CreateAccessKeyOutput, error) {
	inURL := "/"
	endpoint := "CreateAccessKey"
	output := &CreateAccessKeyOutput{}

	if input == nil {
		input = &CreateAccessKeyInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)
	fmt.Printf("[DEBUG ERROR] REQ => %+v => ERR %s", req, err)

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		fmt.Printf("[DEBUG ERROR] DO => %v => ERR %s", output, err)

		return nil, err
	}

	return output, nil
}
func (v ICUOperations) DeleteAccessKey(input *DeleteAccessKeyInput) (*DeleteAccessKeyOutput, error) {
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
func (v ICUOperations) UpdateAccessKey(input *UpdateAccessKeyInput) (*UpdateAccessKeyOutput, error) {
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
func (v ICUOperations) ListAccessKeys(input *ListAccessKeysInput) (*ListAccessKeysOutput, error) {
	inURL := "/"
	endpoint := "ListAccessKeys"
	output := &ListAccessKeysOutput{}

	if input == nil {
		input = &ListAccessKeysInput{}
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
