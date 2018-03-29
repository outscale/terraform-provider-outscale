package icu

const (
	InstanceAttributeNameUserData = "userData"
)

type CreateAccessKeyInput struct {
	_ struct{} `type:"structure"`

	AccessKeyId     *string `type:"string" required:"false"`
	SecretAccessKey *string `type:"string" required:"false"`
}
type CreateAccessKeyOutput struct {
	_                struct{} `type:"structure"`
	AccessKey        *string  `locationName:"accessKey" type:"structure"`
	SecretAccessKey  *string  `locationName:"accessKey" type:"structure"`
	ResponseMetadata *string  `locationName:"responseMetaData" type:"structure"`
}
type DeleteAccessKeyInput struct {
	_           struct{} `type:"structure"`
	AccessKeyId *string  `type:"string" required:"true"`
}
type DeleteAccessKeyOutput struct {
	_                struct{} `type:"structure"`
	ResponseMetadata *string  `locationName:"responseMetaData" type:"structure"`
	Return           *bool    `locationName:"deleteAccessKey" type:"boolean"`
}
type UpdateAccessKeyInput struct {
	_           struct{} `type:"structure"`
	AccessKeyId *string  `type:"string" required:"true"`
	Status      *string  `locationName:"status" type:"string" enum:"StatusType"`
}
type UpdateAccessKeyOutput struct {
	_                struct{} `type:"structure"`
	ResponseMetadata *string  `locationName:"responseMetaData" type:"structure"`
	Return           *bool    `locationName:"updateAccessKey" type:"boolean"`
}
type DescribeAccessKeyInput struct {
	_               struct{} `type:"structure"`
	AccessKeyId     *string  `type:"string" required:"false"`
	SecretAccessKey *string  `type:"string" required:"false"`
}
type DescribeAccessKeyOutput struct {
	_                struct{} `type:"structure"`
	AccessKey        *string  `locationName:"accessKey" type:"structure"`
	ResponseMetadata *string  `locationName:"responseMetaData" type:"structure"`
}
