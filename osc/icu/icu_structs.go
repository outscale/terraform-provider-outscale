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

//ReadCatalogInput ...
type ReadCatalogInput struct {
	_ struct{} `type:"structure"`
}

//ReadCatalogOutput ...
type ReadCatalogOutput struct {
	Catalog          *Catalog          `type:"structure"`
	ResponseMetadata *ResponseMetadata `type:"structure"`
}

//Catalog ...
type Catalog struct {
	Attributes []*CatalogAttribute `type:"list"`
	Entries    []*CatalogEntry     `type:"list"`
}

//CatalogAttribute ...
type CatalogAttribute struct {
	Key   *string `type:"string"`
	Value *string `type:"string"`
}

//ResponseMetadata ...
type ResponseMetadata struct {
	RequestID *string `locationName:"RequestId" type:"string"`
}

//CatalogEntry ...
type CatalogEntry struct {
	Attributes []*CatalogAttribute `type:"structure"`
	Key        *string             `type:"string"`
	Value      *int64              `type:"integer"`
	Title      *string             `type:"string"`
}

//ReadConsumptionAccountInput ...
type ReadConsumptionAccountInput struct {
	FromDate *string `type:"string"`
	ToDate   *string `type:"string"`
}

//ReadConsumptionAccountOutput ...
type ReadConsumptionAccountOutput struct {
	Entries          []*ConsumptionEntry `type:"structure"`
	ResponseMetadata *RequestID          `json:"ResponseMetadata" type:"structure"`
}

//ConsumptionEntry ...
type ConsumptionEntry struct {
	Category  *string  `type:"string"`
	FromDate  *string  `type:"string"`
	Operation *string  `type:"string"`
	Service   *string  `type:"string"`
	Title     *string  `type:"string"`
	ToDate    *string  `type:"string"`
	Type      *string  `type:"string"`
	Value     *float64 `type:"integer"`
}

// ReadAccountInput contains the GetAccount request.
type ReadAccountInput struct {
	_ struct{} `type:"structure"`
}

// ReadAccountOutput contains the response to a successful GetAccount request.
type ReadAccountOutput struct {
	_                struct{}   `type:"structure"`
	Account          *Account   `json:"Account" locationName:"account" type:"structure"`
	ResponseMetadata *RequestID `json:"ResponseMetadata" type:"structure"`
}

// Account contains the response to a successful ListAccessKeys request.
type Account struct {
	_            struct{} `type:"structure"`
	AccountPid   *string  `json:"AccountPid"  locationName:"accountPid" type:"string"`
	City         *string  `json:"City"  locationName:"city" type:"string"`
	CompanyName  *string  `json:"CompanyName"  locationName:"companyName" type:"string"`
	Country      *string  `json:"Country"  locationName:"country" type:"string"`
	CustomerId   *string  `json:"CustomerId"  locationName:"customerId" type:"string"`
	Email        *string  `json:"Email"  locationName:"email" type:"string"`
	FirstName    *string  `json:"FirstName"  locationName:"firstName" type:"string"`
	JobTitle     *string  `json:"JobTitle"  locationName:"jobTitle" type:"string"`
	LastName     *string  `json:"LastName"  locationName:"lastName" type:"string"`
	MobileNumber *string  `json:"MobileNumber"  locationName:"mobileNumber" type:"string"`
	PhoneNumber  *string  `json:"PhoneNumber"  locationName:"phoneNumber" type:"string"`
	State        *string  `json:"State"  locationName:"state" type:"string"`
	VatNumber    *string  `json:"VatNumber"  locationName:"vatNumber" type:"string"`
	Zipcode      *string  `json:"Zipcode"  locationName:"zipcode" type:"string"`
}
