package icu

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/terraform-providers/terraform-provider-outscale/osc/common"
)

const (
	InstanceAttributeNameUserData = "userData"
)

type CreateAccessKeyInput struct {
	_ struct{} `type:"structure"`

	// The name of the IAM user that the new key will belong to.
	//
	// This parameter allows (per its regex pattern (http://wikipedia.org/wiki/regex))
	// a string of characters consisting of upper and lowercase alphanumeric characters
	// with no spaces. You can also include any of the following characters: _+=,.@-
	UserName        *string       `min:"1" type:"string"`
	AccessKeyId     *string       `type:"string"`
	SecretAccessKey *string       `type:"string"`
	Tag             []*common.Tag `locationName:"tag" locationNameList:"item" type:"list"`
}
type CreateApiKey struct {
	_ struct{} `type:"structure"`

	// The name of the IAM user that the new key will belong to.
	//
	// This parameter allows (per its regex pattern (http://wikipedia.org/wiki/regex))
	// a string of characters consisting of upper and lowercase alphanumeric characters
	// with no spaces. You can also include any of the following characters: _+=,.@-
	UserName        *string `min:"1" type:"string"`
	AccessKeyId     *string `type:"string"`
	SecretAccessKey *string `type:"string"`
}

// String returns the string representation
func (s CreateAccessKeyInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateAccessKeyInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CreateAccessKeyInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "CreateAccessKeyInput"}
	if s.UserName != nil && len(*s.UserName) < 1 {
		invalidParams.Add(request.NewErrParamMinLen("UserName", 1))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetUserName sets the UserName field's value.
func (s *CreateAccessKeyInput) SetUserName(v string) *CreateAccessKeyInput {
	s.UserName = &v
	return s
}

// Contains the response to a successful CreateAccessKey request.
type CreateAccessKeyOutput struct {
	_ struct{} `type:"structure"`

	// A structure with details about the access key.
	//
	// AccessKey is a required field
	AccessKey *AccessKey `json:"accessKey" type:"structure" required:"true"`
}

// String returns the string representation
func (s CreateAccessKeyOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateAccessKeyOutput) GoString() string {
	return s.String()
}

// SetAccessKey sets the AccessKey field's value.
func (s *CreateAccessKeyOutput) SetAccessKey(v *AccessKey) *CreateAccessKeyOutput {
	s.AccessKey = v
	return s
}

type AccessKey struct {
	_ struct{} `type:"structure"`

	// The ID for this access key.
	//
	// AccessKeyId is a required field
	AccessKeyId *string `min:"16" type:"string" required:"true"`

	// The date when the access key was created.
	CreateDate *time.Time `type:"timestamp" timestampFormat:"iso8601"`

	// The secret key used to sign requests.
	//
	// SecretAccessKey is a required field
	SecretAccessKey *string `type:"string" required:"true"`

	// The status of the access key. Active means that the key is valid for API
	// calls, while Inactive means it is not.
	//
	// Status is a required field
	Status *string `type:"string" required:"true" enum:"statusType"`

	// The name of the IAM user that the access key is associated with.
	//
	// UserName is a required field
	UserName *string `min:"1" type:"string" required:"true"`
}

// String returns the string representation
func (s AccessKey) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AccessKey) GoString() string {
	return s.String()
}

// SetAccessKeyId sets the AccessKeyId field's value.
func (s *AccessKey) SetAccessKeyId(v string) *AccessKey {
	s.AccessKeyId = &v
	return s
}

// SetCreateDate sets the CreateDate field's value.
func (s *AccessKey) SetCreateDate(v time.Time) *AccessKey {
	s.CreateDate = &v
	return s
}

// SetSecretAccessKey sets the SecretAccessKey field's value.
func (s *AccessKey) SetSecretAccessKey(v string) *AccessKey {
	s.SecretAccessKey = &v
	return s
}

// SetStatus sets the Status field's value.
func (s *AccessKey) SetStatus(v string) *AccessKey {
	s.Status = &v
	return s
}

// SetUserName sets the UserName field's value.
func (s *AccessKey) SetUserName(v string) *AccessKey {
	s.UserName = &v
	return s
}

type DeleteAccessKeyInput struct {
	_ struct{} `type:"structure"`

	// The access key ID for the access key ID and secret access key you want to
	// delete.
	//
	// This parameter allows (per its regex pattern (http://wikipedia.org/wiki/regex))
	// a string of characters that can consist of any upper or lowercased letter
	// or digit.
	//
	// AccessKeyId is a required field
	AccessKeyId *string `min:"16" type:"string" required:"true"`

	// The name of the user whose access key pair you want to delete.
	//
	// This parameter allows (per its regex pattern (http://wikipedia.org/wiki/regex))
	// a string of characters consisting of upper and lowercase alphanumeric characters
	// with no spaces. You can also include any of the following characters: _+=,.@-
	UserName *string `min:"1" type:"string"`
}

// String returns the string representation
func (s DeleteAccessKeyInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteAccessKeyInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteAccessKeyInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeleteAccessKeyInput"}
	if s.AccessKeyId == nil {
		invalidParams.Add(request.NewErrParamRequired("AccessKeyId"))
	}
	if s.AccessKeyId != nil && len(*s.AccessKeyId) < 16 {
		invalidParams.Add(request.NewErrParamMinLen("AccessKeyId", 16))
	}
	if s.UserName != nil && len(*s.UserName) < 1 {
		invalidParams.Add(request.NewErrParamMinLen("UserName", 1))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAccessKeyId sets the AccessKeyId field's value.
func (s *DeleteAccessKeyInput) SetAccessKeyId(v string) *DeleteAccessKeyInput {
	s.AccessKeyId = &v
	return s
}

// SetUserName sets the UserName field's value.
func (s *DeleteAccessKeyInput) SetUserName(v string) *DeleteAccessKeyInput {
	s.UserName = &v
	return s
}

type DeleteAccessKeyOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeleteAccessKeyOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteAccessKeyOutput) GoString() string {
	return s.String()
}

type UpdateAccessKeyInput struct {
	_ struct{} `type:"structure"`

	// The access key ID of the secret access key you want to update.
	//
	// This parameter allows (per its regex pattern (http://wikipedia.org/wiki/regex))
	// a string of characters that can consist of any upper or lowercased letter
	// or digit.
	//
	// AccessKeyId is a required field
	AccessKeyId *string `min:"16" type:"string" required:"true"`

	// The status you want to assign to the secret access key. Active means that
	// the key can be used for API calls to AWS, while Inactive means that the key
	// cannot be used.
	//
	// Status is a required field
	Status *string `type:"string" required:"true" enum:"statusType"`

	// The name of the user whose key you want to update.
	//
	// This parameter allows (per its regex pattern (http://wikipedia.org/wiki/regex))
	// a string of characters consisting of upper and lowercase alphanumeric characters
	// with no spaces. You can also include any of the following characters: _+=,.@-
	UserName *string `min:"1" type:"string"`
}

// String returns the string representation
func (s UpdateAccessKeyInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s UpdateAccessKeyInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *UpdateAccessKeyInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "UpdateAccessKeyInput"}
	if s.AccessKeyId == nil {
		invalidParams.Add(request.NewErrParamRequired("AccessKeyId"))
	}
	if s.AccessKeyId != nil && len(*s.AccessKeyId) < 16 {
		invalidParams.Add(request.NewErrParamMinLen("AccessKeyId", 16))
	}
	if s.Status == nil {
		invalidParams.Add(request.NewErrParamRequired("Status"))
	}
	if s.UserName != nil && len(*s.UserName) < 1 {
		invalidParams.Add(request.NewErrParamMinLen("UserName", 1))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAccessKeyId sets the AccessKeyId field's value.
func (s *UpdateAccessKeyInput) SetAccessKeyId(v string) *UpdateAccessKeyInput {
	s.AccessKeyId = &v
	return s
}

// SetStatus sets the Status field's value.
func (s *UpdateAccessKeyInput) SetStatus(v string) *UpdateAccessKeyInput {
	s.Status = &v
	return s
}

// SetUserName sets the UserName field's value.
func (s *UpdateAccessKeyInput) SetUserName(v string) *UpdateAccessKeyInput {
	s.UserName = &v
	return s
}

type UpdateAccessKeyOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s UpdateAccessKeyOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s UpdateAccessKeyOutput) GoString() string {
	return s.String()
}

type ListAccessKeysInput struct {
	_ struct{} `type:"structure"`

	// Use this parameter only when paginating results and only after you receive
	// a response indicating that the results are truncated. Set it to the value
	// of the Marker element in the response that you received to indicate where
	// the next call should start.
	// Marker *string `min:"1" type:"string"`

	// (Optional) Use this only when paginating results to indicate the maximum
	// number of items you want in the response. If additional items exist beyond
	// the maximum you specify, the IsTruncated response element is true.
	//
	// If you do not include this parameter, it defaults to 100. Note that IAM might
	// return fewer results, even when there are more results available. In that
	// case, the IsTruncated response element returns true and Marker contains a
	// value to include in the subsequent call that tells the service where to continue
	// from.
	// MaxItems *int64 `min:"1" type:"integer"`

	// The name of the user.
	//
	// This parameter allows (per its regex pattern (http://wikipedia.org/wiki/regex))
	// a string of characters consisting of upper and lowercase alphanumeric characters
	// with no spaces. You can also include any of the following characters: _+=,.@-
	// UserName *string `min:"1" type:"string"`

	// Tags []*common.Tag `locationName:"tagSet" locationNameList:"item" type:"list"`
}

// SetTags sets the Tags field's value.
// func (s *ListAccessKeysInput) SetTags(v []*common.Tag) *ListAccessKeysInput {
// 	s.Tags = v
// 	return s
// }

// String returns the string representation
func (s ListAccessKeysInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ListAccessKeysInput) GoString() string {
	return s.String()
}

// SetMarker sets the Marker field's value.
// func (s *ListAccessKeysInput) SetMarker(v string) *ListAccessKeysInput {
// 	s.Marker = &v
// 	return s
// }

// // SetMaxItems sets the MaxItems field's value.
// func (s *ListAccessKeysInput) SetMaxItems(v int64) *ListAccessKeysInput {
// 	s.MaxItems = &v
// 	return s
// }

// // SetUserName sets the UserName field's value.
// func (s *ListAccessKeysInput) SetUserName(v string) *ListAccessKeysInput {
// 	s.UserName = &v
// 	return s
// }

// Contains the response to a successful ListAccessKeys request.
type ListAccessKeysOutput struct {
	_ struct{} `type:"structure"`

	// A list of objects containing metadata about the access keys.
	//
	// AccessKeyMetadata is a required field
	AccessKeyMetadata []*AccessKeyMetadata `json:"accessKeys" locationName:"accessKeys" type:"list" required:"true"`

	// A flag that indicates whether there are more items to return. If your results
	// were truncated, you can make a subsequent pagination request using the Marker
	// request parameter to retrieve more items. Note that IAM might return fewer
	// than the MaxItems number of results even when there are more results available.
	// We recommend that you check IsTruncated after every call to ensure that you
	// receive all of your results.
	IsTruncated *bool `type:"boolean"`

	// When IsTruncated is true, this element is present and contains the value
	// to use for the Marker parameter in a subsequent pagination request.
	Marker *string `min:"1" type:"string"`

	ResponseMetadata RequestId `json:"ResponseMetadata" locationName:"requestId" type:"string"`
}

type RequestId struct {
	RequestId *string `json:"RequestId"  locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s ListAccessKeysOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ListAccessKeysOutput) GoString() string {
	return s.String()
}

// SetAccessKeyMetadata sets the AccessKeyMetadata field's value.
func (s *ListAccessKeysOutput) SetAccessKeyMetadata(v []*AccessKeyMetadata) *ListAccessKeysOutput {
	s.AccessKeyMetadata = v
	return s
}

// SetIsTruncated sets the IsTruncated field's value.
func (s *ListAccessKeysOutput) SetIsTruncated(v bool) *ListAccessKeysOutput {
	s.IsTruncated = &v
	return s
}

// SetMarker sets the Marker field's value.
func (s *ListAccessKeysOutput) SetMarker(v string) *ListAccessKeysOutput {
	s.Marker = &v
	return s
}

type AccessKeyMetadata struct {
	_ struct{} `type:"structure"`

	// The ID for this access key.
	AccessKeyId *string `min:"16" type:"string"`

	// The date when the access key was created.
	CreateDate *time.Time `type:"timestamp" timestampFormat:"iso8601"`

	// The status of the access key. Active means the key is valid for API calls;
	// Inactive means it is not.
	Status *string `type:"string" enum:"statusType"`

	// The name of the IAM user that the key is associated with.
	UserName        *string       `min:"1" type:"string"`
	OwnerId         *string       `type:"string"`
	SecretAccessKey *string       `type:"string"`
	Tags            []*common.Tag `locationName:"tag_set" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s AccessKeyMetadata) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AccessKeyMetadata) GoString() string {
	return s.String()
}

// SetAccessKeyId sets the AccessKeyId field's value.
func (s *AccessKeyMetadata) SetAccessKeyId(v string) *AccessKeyMetadata {
	s.AccessKeyId = &v
	return s
}

// SetCreateDate sets the CreateDate field's value.
func (s *AccessKeyMetadata) SetCreateDate(v time.Time) *AccessKeyMetadata {
	s.CreateDate = &v
	return s
}

// SetStatus sets the Status field's value.
func (s *AccessKeyMetadata) SetStatus(v string) *AccessKeyMetadata {
	s.Status = &v
	return s
}

// SetUserName sets the UserName field's value.
func (s *AccessKeyMetadata) SetUserName(v string) *AccessKeyMetadata {
	s.UserName = &v
	return s
}
