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
}

// GetPolicyOutput ...
type GetPolicyOutput struct {
	_         struct{} `type:"structure"`
	Policy    *Policy  `type:"structure"`
	RequestId *string  `type:"string"`
}

// GetPolicyVersionInput ...
type GetPolicyVersionInput struct {
	_         struct{} `type:"structure"`
	PolicyArn *string  `min:"20" type:"string" required:"true"`
	VersionId *string  `type:"string" required:"true"`
}

// GetPolicyVersionOutput ...
type GetPolicyVersionOutput struct {
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
	_                         struct{}                   `type:"structure"`
	ServerCertificateMetadata *ServerCertificateMetadata `type:"structure"`
}

// ServerCertificateMetadata ...
type ServerCertificateMetadata struct {
	_                     struct{}   `type:"structure"`
	Arn                   *string    `min:"20" type:"string" required:"true"`
	Expiration            *time.Time `type:"timestamp" timestampFormat:"iso8601"`
	Path                  *string    `min:"1" type:"string" required:"true"`
	ServerCertificateId   *string    `min:"16" type:"string" required:"true"`
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
	_                 struct{}           `type:"structure"`
	ServerCertificate *ServerCertificate `type:"structure" required:"true"`
	ResponseMetadata  *ResponseMetadata  `type:"structure" required:"true"`
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
	_                             struct{}                     `type:"structure"`
	IsTruncated                   *bool                        `type:"boolean"`
	Marker                        *string                      `min:"1" type:"string"`
	ResponseMetadata              *ResponseMetadata            `type:"structure"`
	ServerCertificateMetadataList []*ServerCertificateMetadata `type:"list" required:"true"`
}

// ResponseMetadata ...
type ResponseMetadata struct {
	RequestId *string `min:"1" type:"string"`
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
