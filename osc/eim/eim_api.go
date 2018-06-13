package eim

import (
	"context"
	"net/http"

	"github.com/terraform-providers/terraform-provider-outscale/osc"
)

//Operations defines all the operations needed for EIM
type Operations struct {
	client *osc.Client
}

//Service all the necessary actions for them EIM service
type Service interface {
	CreatePolicy(input *CreatePolicyInput) (*CreatePolicyOutput, error)
	GetPolicy(input *GetPolicyInput) (*GetPolicyOutput, error)
	GetPolicyVersion(input *GetPolicyVersionInput) (*GetPolicyVersionOutput, error)
	DeletePolicy(input *DeletePolicyInput) (*DeletePolicyOutput, error)
	DeletePolicyVersion(input *DeletePolicyVersionInput) (*DeletePolicyVersionOutput, error)
	ListPolicyVersions(input *ListPolicyVersionsInput) (*ListPolicyVersionsOutput, error)
	UploadServerCertificate(input *UploadServerCertificateInput) (*UploadServerCertificateOutput, error)
	GetServerCertificate(input *GetServerCertificateInput) (*GetServerCertificateOutput, error)
	DeleteServerCertificate(input *DeleteServerCertificateInput) (*DeleteServerCertificateOutput, error)
}

// CreatePolicy ...
func (v Operations) CreatePolicy(input *CreatePolicyInput) (*CreatePolicyOutput, error) {
	inURL := "/"
	endpoint := "CreatePolicy"
	output := &CreatePolicyOutput{}

	if input == nil {
		input = &CreatePolicyInput{}
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

// GetPolicy ...
func (v Operations) GetPolicy(input *GetPolicyInput) (*GetPolicyOutput, error) {
	inURL := "/"
	endpoint := "GetPolicy"
	output := &GetPolicyOutput{}

	if input == nil {
		input = &GetPolicyInput{}
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

// GetPolicyVersion ...
func (v Operations) GetPolicyVersion(input *GetPolicyVersionInput) (*GetPolicyVersionOutput, error) {
	inURL := "/"
	endpoint := "GetPolicyVersion"
	output := &GetPolicyVersionOutput{}

	if input == nil {
		input = &GetPolicyVersionInput{}
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

// DeletePolicy ...
func (v Operations) DeletePolicy(input *DeletePolicyInput) (*DeletePolicyOutput, error) {
	inURL := "/"
	endpoint := "DeletePolicy"
	output := &DeletePolicyOutput{}

	if input == nil {
		input = &DeletePolicyInput{}
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

// DeletePolicyVersion ...
func (v Operations) DeletePolicyVersion(input *DeletePolicyVersionInput) (*DeletePolicyVersionOutput, error) {
	inURL := "/"
	endpoint := "DeletePolicyVersion"
	output := &DeletePolicyVersionOutput{}

	if input == nil {
		input = &DeletePolicyVersionInput{}
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

// ListPolicyVersions ...
func (v Operations) ListPolicyVersions(input *ListPolicyVersionsInput) (*ListPolicyVersionsOutput, error) {
	inURL := "/"
	endpoint := "ListPolicyVersions"
	output := &ListPolicyVersionsOutput{}

	if input == nil {
		input = &ListPolicyVersionsInput{}
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

// UploadServerCertificate Uploads a server certificate and its matching private key.
func (v Operations) UploadServerCertificate(input *UploadServerCertificateInput) (*UploadServerCertificateOutput, error) {
	inURL := "/"
	endpoint := "UploadServerCertificate"
	output := &UploadServerCertificateOutput{}

	if input == nil {
		input = &UploadServerCertificateInput{}
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

// GetServerCertificate Gets a server certificate and its matching private key.
func (v Operations) GetServerCertificate(input *GetServerCertificateInput) (*GetServerCertificateOutput, error) {
	inURL := "/"
	endpoint := "GetServerCertificate"
	output := &GetServerCertificateOutput{}

	if input == nil {
		input = &GetServerCertificateInput{}
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

// DeleteServerCertificate Deletes a server certificate and its matching private key.
func (v Operations) DeleteServerCertificate(input *DeleteServerCertificateInput) (*DeleteServerCertificateOutput, error) {
	inURL := "/"
	endpoint := "DeleteServerCertificate"
	output := &DeleteServerCertificateOutput{}

	if input == nil {
		input = &DeleteServerCertificateInput{}
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
