package eim

import (
	"time"
)

// CreatePolicyInput ...
type CreatePolicyInput struct {
	_              struct{} `type:"structure"`
	Description    *string  `type:"string"`
	Path           *string  `type:"string"`
	PolicyDocument *string  `min:"1" type:"string" required:"true"`
	PolicyName     *string  `min:"1" type:"string" required:"true"`
}

// CreatePolicyOutput ...
type CreatePolicyOutput struct {
	CreatePolicyResult *CreatePolicyResult `type:"structure"`
	ResponseMetadata   *ResponseMetadata   `type:"structure"`
}

// CreatePolicyResult ...
type CreatePolicyResult struct {
	_      struct{} `type:"structure"`
	Policy *Policy  `type:"structure"`
}

// Policy ...
type Policy struct {
	_                struct{}   `type:"structure"`
	Arn              *string    `min:"20" type:"string"`
	AttachmentCount  *int64     `type:"integer"`
	CreateDate       *time.Time `type:"timestamp" timestampFormat:"iso8601"`
	DefaultVersionId *string    `type:"string"`
	Description      *string    `type:"string"`
	IsAttachable     *bool      `type:"boolean"`
	Path             *string    `type:"string"`
	PolicyId         *string    `min:"16" type:"string"`
	PolicyName       *string    `min:"1" type:"string"`
	UpdateDate       *time.Time `type:"timestamp" timestampFormat:"iso8601"`
}

// GetPolicyInput ...
type GetPolicyInput struct {
	_         struct{} `type:"structure"`
	PolicyArn *string  `min:"20" type:"string" required:"true"`
	GroupName *string  `type:"string"`
	IsLinked  *bool    `type:"boolean"`
	Path      *string  `type:"string"`
	UserName  *string  `type:"string"`
}

// GetPolicyOutput ...
type GetPolicyOutput struct {
	GetPolicyResult  *GetPolicyResult  `type:"structure"`
	ResponseMetadata *ResponseMetadata `type:"structure"`
}

// GetPolicyResult ...
type GetPolicyResult struct {
	_      struct{} `type:"structure"`
	Policy *Policy  `type:"structure"`
}

// GetPolicyVersionInput ...
type GetPolicyVersionInput struct {
	_         struct{} `type:"structure"`
	PolicyArn *string  `min:"20" type:"string" required:"true"`
	VersionId *string  `type:"string" required:"true"`
}

// GetPolicyVersionOutput ...
type GetPolicyVersionOutput struct {
	_                      struct{}                `type:"structure"`
	GetPolicyVersionResult *GetPolicyVersionResult `type:"structure"`
	ResponseMetadata       *ResponseMetadata       `type:"structure"`
}

// GetPolicyVersionResult ...
type GetPolicyVersionResult struct {
	_             struct{}       `type:"structure"`
	PolicyVersion *PolicyVersion `type:"structure"`
}

// PolicyVersion ...
type PolicyVersion struct {
	_                struct{}   `type:"structure"`
	CreateDate       *time.Time `type:"timestamp" timestampFormat:"iso8601"`
	Document         *string    `min:"1" type:"string"`
	IsDefaultVersion *bool      `type:"boolean"`
	VersionId        *string    `type:"string"`
}

// DeletePolicyInput ...
type DeletePolicyInput struct {
	_         struct{} `type:"structure"`
	PolicyArn *string  `min:"20" type:"string" required:"true"`
}

// DeletePolicyOutput ...
type DeletePolicyOutput struct {
	_ struct{} `type:"structure"`
}

// DeletePolicyVersionInput ...
type DeletePolicyVersionInput struct {
	_         struct{} `type:"structure"`
	PolicyArn *string  `min:"20" type:"string" required:"true"`
	VersionId *string  `type:"string" required:"true"`
}

// DeletePolicyVersionOutput ...
type DeletePolicyVersionOutput struct {
	_ struct{} `type:"structure"`
}

// ListPolicyVersionsInput ...
type ListPolicyVersionsInput struct {
	_         struct{} `type:"structure"`
	Marker    *string  `min:"1" type:"string"`
	MaxItems  *int64   `min:"1" type:"integer"`
	PolicyArn *string  `min:"20" type:"string" required:"true"`
}

// ListPolicyVersionsOutput ...
type ListPolicyVersionsOutput struct {
	_           struct{}         `type:"structure"`
	IsTruncated *bool            `type:"boolean"`
	Marker      *string          `min:"1" type:"string"`
	Versions    []*PolicyVersion `type:"list"`
}

// CreateGroupInput ...
type CreateGroupInput struct {
	_         struct{} `type:"structure"`
	GroupName *string  `min:"1" type:"string" required:"true"`
	Path      *string  `min:"1" type:"string"`
}

// CreateGroupResult ...
type CreateGroupResult struct {
	_     struct{} `type:"structure"`
	Group *Group   `type:"structure" required:"true"`
}

// CreateGroupOutput ...
type CreateGroupOutput struct {
	CreateGroupResult *CreateGroupResult `type:"structure"`
	ResponseMetadata  *ResponseMetadata  `type:"structure"`
}

// Group ...
type Group struct {
	_          struct{}   `type:"structure"`
	Arn        *string    `min:"20" type:"string" required:"true"`
	CreateDate *time.Time `type:"timestamp" timestampFormat:"iso8601" required:"true"`
	GroupId    *string    `min:"16" type:"string" required:"true"`
	GroupName  *string    `min:"1" type:"string" required:"true"`
	UserName   *string    `min:"1" type:"string" required:"true"`
	Path       *string    `min:"1" type:"string" required:"true"`
}

// GetGroupInput ...
type GetGroupInput struct {
	_         struct{} `type:"structure"`
	GroupName *string  `min:"1" type:"string" required:"true"`
	Marker    *string  `min:"1" type:"string"`
	MaxItems  *int64   `min:"1" type:"integer"`
}

// GetGroupResult ...
type GetGroupResult struct {
	_           struct{} `type:"structure"`
	Group       *Group   `type:"structure" required:"true"`
	IsTruncated *bool    `type:"boolean"`
	Marker      *string  `min:"1" type:"string"`
	Users       []*User  `type:"list" required:"true"`
}

// GetGroupOutput ...
type GetGroupOutput struct {
	GetGroupResult   *GetGroupResult   `type:"structure"`
	ResponseMetadata *ResponseMetadata `type:"structure"`
}

// User ...
type User struct {
	_                struct{}   `type:"structure"`
	Arn              *string    `min:"20" type:"string" required:"true"`
	CreateDate       *time.Time `type:"timestamp" timestampFormat:"iso8601" required:"true"`
	PasswordLastUsed *time.Time `type:"timestamp" timestampFormat:"iso8601"`
	Path             *string    `min:"1" type:"string" required:"true"`
	UserId           *string    `min:"16" type:"string" required:"true"`
	UserName         *string    `min:"1" type:"string" required:"true"`
}

// UpdateGroupInput ...
type UpdateGroupInput struct {
	_            struct{} `type:"structure"`
	GroupName    *string  `min:"1" type:"string" required:"true"`
	NewGroupName *string  `min:"1" type:"string"`
	NewPath      *string  `min:"1" type:"string"`
}

// UpdateGroupOutput ...
type UpdateGroupOutput struct {
	_ struct{} `type:"structure"`
}

// DeleteGroupInput ...
type DeleteGroupInput struct {
	_         struct{} `type:"structure"`
	GroupName *string  `min:"1" type:"string" required:"true"`
}

// DeleteGroupOutput ...
type DeleteGroupOutput struct {
	_ struct{} `type:"structure"`
}

// AddUserToGroupInput ...
type AddUserToGroupInput struct {
	_         struct{} `type:"structure"`
	GroupName *string  `min:"1" type:"string" required:"true"`
	UserName  *string  `min:"1" type:"string" required:"true"`
}

// AddUserToGroupOutput ...
type AddUserToGroupOutput struct {
	_ struct{} `type:"structure"`
}

// RemoveUserFromGroupInput ...
type RemoveUserFromGroupInput struct {
	_         struct{} `type:"structure"`
	GroupName *string  `min:"1" type:"string" required:"true"`
	UserName  *string  `min:"1" type:"string" required:"true"`
}

// RemoveUserFromGroupOutput ...
type RemoveUserFromGroupOutput struct {
	_ struct{} `type:"structure"`
}

// CreateUserInput ...
type CreateUserInput struct {
	_        struct{} `type:"structure"`
	Path     *string  `min:"1" type:"string"`
	UserName *string  `min:"1" type:"string" required:"true"`
}

// CreateUserOutput ...
type CreateUserOutput struct {
	_    struct{} `type:"structure"`
	User *User    `type:"structure"`
}

// GetUserInput ...
type GetUserInput struct {
	_        struct{} `type:"structure"`
	UserName *string  `min:"1" type:"string"`
}

// GetUserOutput ...
type GetUserOutput struct {
	_    struct{} `type:"structure"`
	User *User    `type:"structure" required:"true"`
}

// UpdateUserInput ...
type UpdateUserInput struct {
	_           struct{} `type:"structure"`
	NewPath     *string  `min:"1" type:"string"`
	NewUserName *string  `min:"1" type:"string"`
	UserName    *string  `min:"1" type:"string" required:"true"`
}

// UpdateUserOutput ...
type UpdateUserOutput struct {
	_ struct{} `type:"structure"`
}

// ListGroupsForUserInput ...
type ListGroupsForUserInput struct {
	_        struct{} `type:"structure"`
	Marker   *string  `min:"1" type:"string"`
	MaxItems *int64   `min:"1" type:"integer"`
	UserName *string  `min:"1" type:"string" required:"true"`
}

// ListGroupsForUserOutput ...
type ListGroupsForUserOutput struct {
	_           struct{} `type:"structure"`
	Groups      []*Group `type:"list" required:"true"`
	IsTruncated *bool    `type:"boolean"`
	Marker      *string  `min:"1" type:"string"`
}

// DeleteUserInput ...
type DeleteUserInput struct {
	_        struct{} `type:"structure"`
	UserName *string  `min:"1" type:"string" required:"true"`
}

// DeleteUserOutput ...
type DeleteUserOutput struct {
	_ struct{} `type:"structure"`
}

// SetDefaultPolicyVersionInput ...
type SetDefaultPolicyVersionInput struct {
	_         struct{} `type:"structure"`
	PolicyArn *string  `min:"20" type:"string" required:"true"`
	VersionId *string  `type:"string" required:"true"`
}

// SetDefaultPolicyVersionOutput ...
type SetDefaultPolicyVersionOutput struct {
	_ struct{} `type:"structure"`
}

// AttachUserPolicyInput ...
type AttachUserPolicyInput struct {
	_         struct{} `type:"structure"`
	PolicyArn *string  `min:"20" type:"string" required:"true"`
	UserName  *string  `min:"1" type:"string" required:"true"`
}

// AttachUserPolicyOutput ...
type AttachUserPolicyOutput struct {
	_ struct{} `type:"structure"`
}

// ListAttachedUserPoliciesInput ...
type ListAttachedUserPoliciesInput struct {
	_          struct{} `type:"structure"`
	Marker     *string  `min:"1" type:"string"`
	MaxItems   *int64   `min:"1" type:"integer"`
	PathPrefix *string  `type:"string"`
	UserName   *string  `min:"1" type:"string" required:"true"`
}

// ListAttachedUserPoliciesResult ...
type ListAttachedUserPoliciesResult struct {
	_                struct{}          `type:"structure"`
	AttachedPolicies []*AttachedPolicy `type:"list"`
	IsTruncated      *bool             `type:"boolean"`
	Marker           *string           `min:"1" type:"string"`
}

// ListAttachedUserPoliciesOutput ...
type ListAttachedUserPoliciesOutput struct {
	_                              struct{}                        `type:"structure"`
	ListAttachedUserPoliciesResult *ListAttachedUserPoliciesResult `type:"structure"`
	ResponseMetadata               *ResponseMetadata               `type:"structure"`
}

// AttachedPolicy ...
type AttachedPolicy struct {
	_          struct{} `type:"structure"`
	PolicyArn  *string  `min:"20" type:"string"`
	PolicyName *string  `min:"1" type:"string"`
}

// DetachUserPolicyInput ...
type DetachUserPolicyInput struct {
	_         struct{} `type:"structure"`
	PolicyArn *string  `min:"20" type:"string" required:"true"`
	UserName  *string  `min:"1" type:"string" required:"true"`
}

// DetachUserPolicyOutput ...
type DetachUserPolicyOutput struct {
	_ struct{} `type:"structure"`
}

// PutUserPolicyInput ...
type PutUserPolicyInput struct {
	_              struct{} `type:"structure"`
	PolicyDocument *string  `min:"1" type:"string" required:"true"`
	PolicyName     *string  `min:"1" type:"string" required:"true"`
	UserName       *string  `min:"1" type:"string" required:"true"`
}

// PutUserPolicyOutput ...
type PutUserPolicyOutput struct {
	_ struct{} `type:"structure"`
}

// GetUserPolicyInput ...
type GetUserPolicyInput struct {
	_          struct{} `type:"structure"`
	PolicyName *string  `min:"1" type:"string" required:"true"`
	UserName   *string  `min:"1" type:"string" required:"true"`
}

// GetUserPolicyOutput ...
type GetUserPolicyOutput struct {
	_              struct{} `type:"structure"`
	PolicyDocument *string  `min:"1" type:"string" required:"true"`
	PolicyName     *string  `min:"1" type:"string" required:"true"`
	UserName       *string  `min:"1" type:"string" required:"true"`
}

// ListUsersInput ...
type ListUsersInput struct {
	_          struct{} `type:"structure"`
	Marker     *string  `min:"1" type:"string"`
	MaxItems   *int64   `min:"1" type:"integer"`
	PathPrefix *string  `min:"1" type:"string"`
}

// ListUsersOutput ...
type ListUsersOutput struct {
	_           struct{} `type:"structure"`
	IsTruncated *bool    `type:"boolean"`
	Marker      *string  `min:"1" type:"string"`
	Users       []*User  `type:"list" required:"true"`
}

// ListGroupsInput ...
type ListGroupsInput struct {
	_          struct{} `type:"structure"`
	Marker     *string  `min:"1" type:"string"`
	MaxItems   *int64   `min:"1" type:"integer"`
	PathPrefix *string  `min:"1" type:"string"`
}

// ListGroupsOutput ...
type ListGroupsOutput struct {
	_           struct{} `type:"structure"`
	Groups      []*Group `type:"list" required:"true"`
	IsTruncated *bool    `type:"boolean"`
	Marker      *string  `min:"1" type:"string"`
}

// UploadServerCertificateInput ...
type UploadServerCertificateInput struct {
	_                     struct{} `type:"structure"`
	CertificateBody       *string  `min:"1" type:"string" required:"true"`
	CertificateChain      *string  `min:"1" type:"string"`
	Path                  *string  `min:"1" type:"string"`
	PrivateKey            *string  `min:"1" type:"string" required:"true"`
	ServerCertificateName *string  `min:"1" type:"string" required:"true"`
}

//UploadServerCertificateOutput Contains the response to a successful UploadServerCertificate request.
type UploadServerCertificateOutput struct {
	UploadServerCertificateResult *UploadServerCertificateResult `type:"structure"`
	ResponseMetadata              *ResponseMetadata              `type:"structure" required:"true"`
}

//UploadServerCertificateResult Contains the response to a successful UploadServerCertificate request.
type UploadServerCertificateResult struct {
	_                         struct{}                   `type:"structure"`
	ServerCertificateMetadata *ServerCertificateMetadata `type:"structure"`
}

// ServerCertificateMetadata ...
type ServerCertificateMetadata struct {
	_                     struct{}   `type:"structure"`
	Arn                   *string    `min:"20" type:"string" required:"true"`
	Expiration            *time.Time `type:"timestamp" timestampFormat:"iso8601"`
	Path                  *string    `min:"1" type:"string" required:"true"`
	ServerCertificateID   *string    `locationName:"ServerCertificateId" min:"16" type:"string" required:"true"`
	ServerCertificateName *string    `min:"1" type:"string" required:"true"`
	UploadDate            *time.Time `type:"timestamp" timestampFormat:"iso8601"`
}

// GetServerCertificateInput ...
type GetServerCertificateInput struct {
	_                     struct{} `type:"structure"`
	ServerCertificateName *string  `min:"1" type:"string" required:"true"`
}

//GetServerCertificateOutput Contains the response to a successful GetServerCertificate request.
type GetServerCertificateOutput struct {
	GetServerCertificateResult *GetServerCertificateResult `type:"structure"`
	ResponseMetadata           *ResponseMetadata           `type:"structure" required:"true"`
}

//GetServerCertificateResult Contains the response to a successful GetServerCertificate request.
type GetServerCertificateResult struct {
	_                 struct{}           `type:"structure"`
	ServerCertificate *ServerCertificate `type:"structure" required:"true"`
}

// ServerCertificate Contains information about a server certificate.
type ServerCertificate struct {
	_                         struct{}                   `type:"structure"`
	CertificateBody           *string                    `min:"1" type:"string" required:"true"` // The contents of the public key certificate.
	CertificateChain          *string                    `min:"1" type:"string"`                 // The contents of the public key certificate chain.
	ServerCertificateMetadata *ServerCertificateMetadata `type:"structure" required:"true"`      // The meta information of the server certificate, such as its name, path, ID, and ARN.
}

// DeleteServerCertificateInput ...
type DeleteServerCertificateInput struct {
	_                     struct{} `type:"structure"`
	ServerCertificateName *string  `min:"1" type:"string" required:"true"` // The name of the server certificate you want to delete.
}

// DeleteServerCertificateOutput ...
type DeleteServerCertificateOutput struct {
	_ struct{} `type:"structure"`
}

//ListServerCertificatesInput ...
type ListServerCertificatesInput struct {
	_          struct{} `type:"structure"`
	Marker     *string  `min:"1" type:"string"`  // Use this parameter only when paginating results and only after you receive a response indicating that the results are truncated. Set it to the value of the Marker element in the response that you received to indicate where the next call should start.
	MaxItems   *int64   `min:"1" type:"integer"` // (Optional) Use this only when paginating results to indicate the maximum number of items you want in the response. If additional items exist beyond the maximum you specify, the IsTruncated response element is true.
	PathPrefix *string  `min:"1" type:"string"`  // The path prefix for filtering the results. For example: /company/servercerts would get all server certificates for which the path starts with /company/servercerts.
}

//ListServerCertificatesOutput Contains the response to a successful ListServerCertificates request.
type ListServerCertificatesOutput struct {
	ListServerCertificatesResult *ListServerCertificatesResult `type:"structure"`
	ResponseMetadata             *ResponseMetadata             `type:"structure"`
}

//ListServerCertificatesResult Contains the response to a successful ListServerCertificates request.
type ListServerCertificatesResult struct {
	_                             struct{}                     `type:"structure"`
	IsTruncated                   *bool                        `type:"boolean"`
	Marker                        *string                      `min:"1" type:"string"`
	ServerCertificateMetadataList []*ServerCertificateMetadata `type:"list" required:"true"`
}

// ResponseMetadata ...
type ResponseMetadata struct {
	RequestID *string `locationName:"RequestId" min:"1" type:"string"`
}

// ListCertificatesOutput ...
type ListCertificatesOutput struct {
	_                      struct{}              `type:"structure"`
	CertificateSummaryList []*CertificateSummary `type:"list"`           // A list of ACM certificates.
	NextToken              *string               `min:"1" type:"string"` // When the list is truncated, this value is present and contains the value to use for the NextToken parameter in a subsequent pagination request.
}

// CertificateSummary ...
type CertificateSummary struct {
	_              struct{} `type:"structure"`
	CertificateArn *string  `min:"20" type:"string"`
	// Fully qualified domain name (FQDN), such as www.example.com or example.com,
	// for the certificate.
	DomainName *string `min:"1" type:"string"`
}

// UpdateServerCertificateInput ...
type UpdateServerCertificateInput struct {
	_                        struct{} `type:"structure"`
	NewPath                  *string  `min:"1" type:"string"`
	NewServerCertificateName *string  `min:"1" type:"string"`
	ServerCertificateName    *string  `min:"1" type:"string" required:"true"`
}

// UpdateServerCertificateOutput ...
type UpdateServerCertificateOutput struct {
	_ struct{} `type:"structure"`
}

// DeleteUserPolicyInput ...
type DeleteUserPolicyInput struct {
	_          struct{} `type:"structure"`
	PolicyName *string  `min:"1" type:"string" required:"true"`
	UserName   *string  `min:"1" type:"string" required:"true"`
}

// DeleteUserPolicyOutput ...
type DeleteUserPolicyOutput struct {
	_ struct{} `type:"structure"`
}

// GetRolePolicyInput ...
type GetRolePolicyInput struct {
	_          struct{} `type:"structure"`
	PolicyName *string  `min:"1" type:"string" required:"true"`
	RoleName   *string  `min:"1" type:"string" required:"true"`
}

// GetRolePolicyResult ...
type GetRolePolicyResult struct {
	_              struct{} `type:"structure"`
	PolicyDocument *string  `min:"1" type:"string" required:"true"`
	PolicyName     *string  `min:"1" type:"string" required:"true"`
	RoleName       *string  `min:"1" type:"string" required:"true"`
}

// GetRolePolicyOutput ...
type GetRolePolicyOutput struct {
	GetRolePolicyResult *GetRolePolicyResult `type:"structure"`
	ResponseMetadata    *ResponseMetadata    `type:"structure"`
}

// PutGroupPolicyInput ...
type PutGroupPolicyInput struct {
	_              struct{} `type:"structure"`
	GroupName      *string  `min:"1" type:"string" required:"true"`
	PolicyDocument *string  `min:"1" type:"string" required:"true"`
	PolicyName     *string  `min:"1" type:"string" required:"true"`
}

// PutGroupPolicyOutput ...
type PutGroupPolicyOutput struct {
	_ struct{} `type:"structure"`
}

// GetGroupPolicyInput ...
type GetGroupPolicyInput struct {
	_          struct{} `type:"structure"`
	GroupName  *string  `min:"1" type:"string" required:"true"`
	PolicyName *string  `min:"1" type:"string" required:"true"`
}

// GetGroupPolicyResult ...
type GetGroupPolicyResult struct {
	_              struct{} `type:"structure"`
	GroupName      *string  `min:"1" type:"string" required:"true"`
	PolicyDocument *string  `min:"1" type:"string" required:"true"`
	PolicyName     *string  `min:"1" type:"string" required:"true"`
}

// GetGroupPolicyOutput ...
type GetGroupPolicyOutput struct {
	GetGroupPolicyResult *GetGroupPolicyResult `type:"structure"`
	ResponseMetadata     *ResponseMetadata     `type:"structure"`
}

// DeleteGroupPolicyInput ...
type DeleteGroupPolicyInput struct {
	_          struct{} `type:"structure"`
	GroupName  *string  `min:"1" type:"string" required:"true"`
	PolicyName *string  `min:"1" type:"string" required:"true"`
}

// DeleteGroupPolicyOutput ...
type DeleteGroupPolicyOutput struct {
	_ struct{} `type:"structure"`
}

// AttachGroupPolicyInput ...
type AttachGroupPolicyInput struct {
	_         struct{} `type:"structure"`
	GroupName *string  `min:"1" type:"string" required:"true"`
	PolicyArn *string  `min:"20" type:"string" required:"true"`
}

// AttachGroupPolicyOutput ...
type AttachGroupPolicyOutput struct {
	_ struct{} `type:"structure"`
}

// DetachGroupPolicyInput ...
type DetachGroupPolicyInput struct {
	_         struct{} `type:"structure"`
	GroupName *string  `min:"1" type:"string" required:"true"`
	PolicyArn *string  `min:"20" type:"string" required:"true"`
}

// DetachGroupPolicyOutput ...
type DetachGroupPolicyOutput struct {
	_ struct{} `type:"structure"`
}

// ListAttachedGroupPoliciesInput ...
type ListAttachedGroupPoliciesInput struct {
	_          struct{} `type:"structure"`
	GroupName  *string  `min:"1" type:"string" required:"true"`
	Marker     *string  `min:"1" type:"string"`
	MaxItems   *int64   `min:"1" type:"integer"`
	PathPrefix *string  `type:"string"`
}

// ListAttachedGroupPoliciesResult ...
type ListAttachedGroupPoliciesResult struct {
	_                struct{}          `type:"structure"`
	AttachedPolicies []*AttachedPolicy `type:"list"`
	IsTruncated      *bool             `type:"boolean"`
	Marker           *string           `min:"1" type:"string"`
}

// ListAttachedGroupPoliciesOutput ...
type ListAttachedGroupPoliciesOutput struct {
	ListAttachedGroupPoliciesResult *ListAttachedGroupPoliciesResult `type:"structure"`
	ResponseMetadata                *ResponseMetadata                `type:"structure"`
}

// CreatePolicyVersionInput ...
type CreatePolicyVersionInput struct {
	_              struct{} `type:"structure"`
	PolicyArn      *string  `min:"20" type:"string" required:"true"`
	PolicyDocument *string  `min:"1" type:"string" required:"true"`
	SetAsDefault   *bool    `type:"boolean"`
}

// CreatePolicyVersionOutput ...
type CreatePolicyVersionOutput struct {
	_                         struct{}                   `type:"structure"`
	CreatePolicyVersionResult *CreatePolicyVersionResult `type:"structure"`
	ResponseMetadata          *ResponseMetadata          `type:"structure"`
}

// CreatePolicyVersionResult ...
type CreatePolicyVersionResult struct {
	_             struct{}       `type:"structure"`
	PolicyVersion *PolicyVersion `type:"structure"`
}
