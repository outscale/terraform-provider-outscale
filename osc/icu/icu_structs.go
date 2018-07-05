package icu

import (
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/common"
)

// CreateAccessKeyInput ...
type CreateAccessKeyInput struct {
	_               struct{}      `type:"structure"`
	UserName        *string       `json:"UserName,omitempty" type:"string"`
	AccessKeyID     *string       `json:"AccessKeyId,omitempty" type:"string"`
	SecretAccessKey *string       `json:"SecretAccessKey,omitempty" type:"string"`
	Tag             []*common.Tag `json:"Tag,omitempty"`
}

// CreateApiKey ...
type CreateApiKey struct {
	_               struct{} `type:"structure"`
	UserName        *string  `min:"1" type:"string"`
	AccessKeyId     *string  `type:"string"`
	SecretAccessKey *string  `type:"string"`
}

// CreateAccessKeyOutput ...
type CreateAccessKeyOutput struct {
	_         struct{}   `type:"structure"`
	AccessKey *AccessKey `json:"accessKey" type:"structure" required:"true"`
}

// AccessKey ...
type AccessKey struct {
	_               struct{}   `type:"structure"`
	AccessKeyId     *string    `min:"16" type:"string" required:"true"`
	CreateDate      *time.Time `type:"timestamp" timestampFormat:"iso8601"`
	SecretAccessKey *string    `type:"string" required:"true"`
	Status          *string    `type:"string" required:"true" enum:"statusType"`
	UserName        *string    `min:"1" type:"string" required:"true"`
}

// DeleteAccessKeyInput ...
type DeleteAccessKeyInput struct {
	_           struct{} `type:"structure"`
	AccessKeyId *string  `min:"16" type:"string" required:"true"`
	UserName    *string  `min:"1" type:"string"`
}

// DeleteAccessKeyOutput ...
type DeleteAccessKeyOutput struct {
	_ struct{} `type:"structure"`
}

// UpdateAccessKeyInput ...
type UpdateAccessKeyInput struct {
	_           struct{} `type:"structure"`
	AccessKeyId *string  `min:"16" type:"string" required:"true"`
	Status      *string  `type:"string" required:"true" enum:"statusType"`
	UserName    *string  `min:"1" type:"string"`
}

// UpdateAccessKeyOutput ...
type UpdateAccessKeyOutput struct {
	_ struct{} `type:"structure"`
}

// ListAccessKeysInput ...
type ListAccessKeysInput struct {
	_ struct{} `type:"structure"`
}

// ListAccessKeysOutput ...
type ListAccessKeysOutput struct {
	_                 struct{}             `type:"structure"`
	AccessKeyMetadata []*AccessKeyMetadata `json:"accessKeys" locationName:"accessKeys" type:"list" required:"true"`
	IsTruncated       *bool                `type:"boolean"`
	Marker            *string              `min:"1" type:"string"`
	ResponseMetadata  RequestID            `json:"ResponseMetadata" locationName:"requestId" type:"string"`
}

// RequestID ...
type RequestID struct {
	RequestID *string `json:"RequestId"  locationName:"requestId" type:"string"`
}

// AccessKeyMetadata ...
type AccessKeyMetadata struct {
	_               struct{}   `type:"structure"`
	AccessKeyID     *string    `json:"accessKeyId" type:"string"`
	CreateDate      *time.Time `json:"createDate" type:"timestamp" timestampFormat:"iso8601"`
	Status          *string    `json:"status" type:"string"`
	UserName        *string    `json:"userName" type:"string"`
	OwnerID         *string    `json:"ownerId" type:"string"`
	SecretAccessKey *string    `json:"secretAccessKey" type:"string"`
	Tags            []*Tag     `json:"tags"`
}

// Tag ...
type Tag struct {
	_     struct{} `type:"structure"`
	Key   *string  `json:"key" type:"string"`
	Value *string  `json:"value" type:"string"`
}
