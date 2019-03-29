package icu

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
	CreateAccessKey(input *CreateAccessKeyInput) (*CreateAccessKeyOutput, error)
	DeleteAccessKey(input *DeleteAccessKeyInput) (*DeleteAccessKeyOutput, error)
	UpdateAccessKey(input *UpdateAccessKeyInput) (*UpdateAccessKeyOutput, error)
	ListAccessKeys(input *ListAccessKeysInput) (*ListAccessKeysOutput, error)
	ReadCatalog(input *ReadCatalogInput) (*ReadCatalogOutput, error)
	ReadPublicCatalog(input *ReadCatalogInput) (*ReadCatalogOutput, error)
	ReadConsumptionAccount(input *ReadConsumptionAccountInput) (*ReadConsumptionAccountOutput, error)
	GetAccount(input *ReadAccountInput) (*ReadAccountOutput, error)
}

// CreateAccessKey ...
func (v Operations) CreateAccessKey(input *CreateAccessKeyInput) (*CreateAccessKeyOutput, error) {
	inURL := "/"
	endpoint := "CreateAccessKey"
	output := &CreateAccessKeyOutput{}

	if input == nil {
		input = &CreateAccessKeyInput{}
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

// DeleteAccessKey ...
func (v Operations) DeleteAccessKey(input *DeleteAccessKeyInput) (*DeleteAccessKeyOutput, error) {
	inURL := "/"
	endpoint := "DeleteAccessKey"
	output := &DeleteAccessKeyOutput{}

	if input == nil {
		input = &DeleteAccessKeyInput{}
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

// UpdateAccessKey ...
func (v Operations) UpdateAccessKey(input *UpdateAccessKeyInput) (*UpdateAccessKeyOutput, error) {
	inURL := "/"
	endpoint := "UpdateAccessKey"
	output := &UpdateAccessKeyOutput{}

	if input == nil {
		input = &UpdateAccessKeyInput{}
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

// ListAccessKeys ...
func (v Operations) ListAccessKeys(input *ListAccessKeysInput) (*ListAccessKeysOutput, error) {
	inURL := "/"
	endpoint := "ListAccessKeys"
	output := &ListAccessKeysOutput{}

	if input == nil {
		input = &ListAccessKeysInput{}
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

// ReadCatalog ...
func (v Operations) ReadCatalog(input *ReadCatalogInput) (*ReadCatalogOutput, error) {
	inURL := "/"
	endpoint := "ReadCatalog"
	output := &ReadCatalogOutput{}

	if input == nil {
		input = &ReadCatalogInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodPost, inURL, input)
	req.Header.Set("Content-Type", "application/x-amz-json-1.1")

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// ReadPublicCatalog ...
func (v Operations) ReadPublicCatalog(input *ReadCatalogInput) (*ReadCatalogOutput, error) {
	inURL := "/"
	endpoint := "ReadPublicCatalog"
	output := &ReadCatalogOutput{}

	if input == nil {
		input = &ReadCatalogInput{}
	}

	req, err := v.client.NewRequest(context.TODO(), endpoint, http.MethodPost, inURL, input)
	req.Header.Set("Content-Type", "application/x-amz-json-1.1")

	if err != nil {
		return nil, err
	}

	err = v.client.Do(context.TODO(), req, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

//ReadConsumptionAccount ...
func (v Operations) ReadConsumptionAccount(input *ReadConsumptionAccountInput) (*ReadConsumptionAccountOutput, error) {
	inURL := "/"
	endpoint := "ReadConsumptionAccount"
	output := &ReadConsumptionAccountOutput{}

	if input == nil {
		input = &ReadConsumptionAccountInput{}
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

// GetAccount gets information about the account that sent the request.
func (v Operations) GetAccount(input *ReadAccountInput) (*ReadAccountOutput, error) {
	inURL := "/"
	endpoint := "GetAccount"

	output := &ReadAccountOutput{}

	if input == nil {
		input = &ReadAccountInput{}
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
