package fcu

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/request"
)

const (
	InstanceAttributeNameUserData = "userData"
)

type DescribeInstancesInput struct {
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	InstanceIds []*string `locationName:"InstanceId" locationNameList:"InstanceId" type:"list"`

	MaxResults *int64 `locationName:"maxResults" type:"integer"`

	NextToken *string `locationName:"nextToken" type:"string"`
}

type Filter struct {
	Name *string `type:"string"`

	Values []*string `locationName:"Value" locationNameList:"item" type:"list"`
}

type DescribeInstancesOutput struct {
	_ struct{} `type:"structure"`

	NextToken *string `locationName:"nextToken" type:"string"`

	OwnerId *string `locationName:"ownerId" locationNameList:"item" type:"string"`

	RequestId *string `locationName:"requestId" type:"string"`

	ReservationId *string `locationName:"reservationId" locationNameList:"item" type:"string"`

	Reservations []*Reservation `locationName:"reservationSet" locationNameList:"item" type:"list"`

	GroupSet []*GroupIdentifier `locationName:"groupSet" locationNameList:"item" type:"list"`
}

type GroupIdentifier struct {
	_ struct{} `type:"structure"`

	GroupId *string `locationName:"groupId" type:"string"`

	GroupName *string `locationName:"groupName" type:"string"`
}

type Reservation struct {
	_ struct{} `type:"structure"`

	Groups []*GroupIdentifier `locationName:"groupSet" locationNameList:"item" type:"list"`

	Instances []*Instance `locationName:"instancesSet" locationNameList:"item" type:"list"`

	OwnerId *string `locationName:"ownerId" type:"string"`

	RequestId *string `locationName:"requestId" type:"string"`

	ReservationId *string `locationName:"reservationId" type:"string"`
}

type Instance struct {
	AmiLaunchIndex *int64 `locationName:"amiLaunchIndex" type:"integer"`

	Architecture *string `locationName:"architecture" type:"string" enum:"ArchitectureValues"`

	BlockDeviceMappings []*InstanceBlockDeviceMapping `locationName:"blockDeviceMapping" locationNameList:"item" type:"list"`

	ClientToken *string `locationName:"clientToken" type:"string"`

	DnsName *string `locationName:"dnsName" type:"string"`

	EbsOptimized *bool `locationName:"ebsOptimized" type:"boolean"`

	GroupSet []*GroupIdentifier `locationName:"groupSet" locationNameList:"item" type:"list"`

	Hypervisor *string `locationName:"hypervisor" type:"string" enum:"HypervisorType"`

	IamInstanceProfile *IamInstanceProfile `locationName:"iamInstanceProfile" type:"structure"`

	ImageId *string `locationName:"imageId" type:"string"`

	InstanceId *string `locationName:"instanceId" type:"string"`

	InstanceLifecycle *string `locationName:"instanceLifecycle" type:"string" enum:"InstanceLifecycleType"`

	InstanceState *InstanceState `locationName:"instanceState" type:"structure"`

	InstanceType *string `locationName:"instanceType" type:"string" enum:"InstanceType"`

	IpAddress *string `locationName:"ipAddress" type:"string"`

	KernelId *string `locationName:"kernelId" type:"string"`

	KeyName *string `locationName:"keyName" type:"string"`

	Monitoring *Monitoring `locationName:"monitoring" type:"structure"`

	NetworkInterfaces []*InstanceNetworkInterface `locationName:"networkInterfaceSet" locationNameList:"item" type:"list"`

	Placement *Placement `locationName:"placement" type:"structure"`

	Platform *string `locationName:"platform" type:"string" enum:"PlatformValues"`

	PrivateDnsName *string `locationName:"privateDnsName" type:"string"`

	PrivateIpAddress *string `locationName:"privateIpAddress" type:"string"`

	ProductCodes []*ProductCode `locationName:"productCodes" locationNameList:"item" type:"list"`

	RamdiskId *string `locationName:"ramdiskId" type:"string"`

	Reason *string `locationName:"reason" type:"string"`

	RootDeviceName *string `locationName:"rootDeviceName" type:"string"`

	RootDeviceType *string `locationName:"rootDeviceType" type:"string" enum:"DeviceType"`

	SourceDestCheck *bool `locationName:"sourceDestCheck" type:"boolean"`

	SpotInstanceRequestId *string `locationName:"spotInstanceRequestId" type:"string"`

	SriovNetSupport *string `locationName:"sriovNetSupport" type:"string"`

	State *InstanceState `locationName:"instanceState" type:"structure"`

	StateReason *StateReason `locationName:"stateReason" type:"structure"`

	SubnetId *string `locationName:"subnetId" type:"string"`

	Tags []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`

	VirtualizationType *string `locationName:"virtualizationType" type:"string" enum:"VirtualizationType"`

	VpcId *string `locationName:"vpcId" type:"string"`
}

type InstanceBlockDeviceMapping struct {
	DeviceName *string `locationName:"deviceName" type:"string"`

	Ebs *EbsInstanceBlockDevice `locationName:"ebs" type:"structure"`
}

type InstanceBlockDeviceMappingSpecification struct {
	_ struct{} `type:"structure"`

	DeviceName *string `locationName:"deviceName" type:"string"`

	Ebs *EbsInstanceBlockDeviceSpecification `locationName:"ebs" type:"structure"`

	NoDevice *string `locationName:"noDevice" type:"string"`

	VirtualName *string `locationName:"virtualName" type:"string"`
}

type InstanceCapacity struct {
	_ struct{} `type:"structure"`

	AvailableCapacity *int64 `locationName:"availableCapacity" type:"integer"`

	InstanceType *string `locationName:"instanceType" type:"string"`

	TotalCapacity *int64 `locationName:"totalCapacity" type:"integer"`
}

type InstanceCount struct {
	_ struct{} `type:"structure"`

	InstanceCount *int64 `locationName:"instanceCount" type:"integer"`

	State *string `locationName:"state" type:"string" enum:"ListingState"`
}

type InstanceExportDetails struct {
	_ struct{} `type:"structure"`

	InstanceId *string `locationName:"instanceId" type:"string"`

	TargetEnvironment *string `locationName:"targetEnvironment" type:"string" enum:"ExportEnvironment"`
}

type InstanceMonitoring struct {
	_ struct{} `type:"structure"`

	InstanceId *string `locationName:"instanceId" type:"string"`

	Monitoring *Monitoring `locationName:"monitoring" type:"structure"`
}

type InstanceNetworkInterface struct {
	Association *InstanceNetworkInterfaceAssociation `locationName:"association" type:"structure"`

	Attachment *InstanceNetworkInterfaceAttachment `locationName:"attachment" type:"structure"`

	Description *string `locationName:"description" type:"string"`

	Groups []*GroupIdentifier `locationName:"groupSet" locationNameList:"item" type:"list"`

	MacAddress *string `locationName:"macAddress" type:"string"`

	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string"`

	OwnerId *string `locationName:"ownerId" type:"string"`

	PrivateDnsName *string `locationName:"privateDnsName" type:"string"`

	PrivateIpAddress *string `locationName:"privateIpAddress" type:"string"`

	PrivateIpAddresses []*InstancePrivateIpAddress `locationName:"privateIpAddressesSet" locationNameList:"item" type:"list"`

	SourceDestCheck *bool `locationName:"sourceDestCheck" type:"bool"`

	Status *string `locationName:"status" type:"string" enum:"NetworkInterfaceStatus"`

	SubnetId *string `locationName:"subnetId" type:"string"`

	VpcId *string `locationName:"vpcId" type:"string"`
}

type InstanceNetworkInterfaceAssociation struct {
	IpOwnerId *string `locationName:"ipOwnerId" type:"string"`

	PublicDnsName *string `locationName:"publicDnsName" type:"string"`

	PublicIp *string `locationName:"publicIp" type:"string"`
}

type InstanceNetworkInterfaceAttachment struct {
	AttachmentId *string `locationName:"attachmentId" type:"string"`

	DeleteOnTermination *bool `locationName:"deleteOnTermination" type:"boolean"`

	DeviceIndex *int64 `locationName:"deviceIndex" type:"integer"`

	Status *string `locationName:"status" type:"string" enum:"AttachmentStatus"`
}

type InstanceNetworkInterfaceSpecification struct {
	_ struct{} `type:"structure"`

	AssociatePublicIpAddress *bool `locationName:"associatePublicIpAddress" type:"boolean"`

	DeleteOnTermination *bool `locationName:"deleteOnTermination" type:"boolean"`

	Description *string `locationName:"description" type:"string"`

	DeviceIndex *int64 `locationName:"deviceIndex" type:"integer"`

	Groups []*string `locationName:"SecurityGroupId" locationNameList:"SecurityGroupId" type:"list"`

	Ipv6AddressCount *int64 `locationName:"ipv6AddressCount" type:"integer"`

	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string"`

	PrivateIpAddress *string `locationName:"privateIpAddress" type:"string"`

	PrivateIpAddresses []*PrivateIpAddressSpecification `locationName:"privateIpAddressesSet" queryName:"PrivateIpAddresses" locationNameList:"item" type:"list"`

	SecurityGroupIds []*string `locationName:"SecurityGroupId" locationNameList:"SecurityGroupId" type:"list"`

	SecondaryPrivateIpAddressCount *int64 `locationName:"secondaryPrivateIpAddressCount" type:"integer"`

	SubnetId *string `locationName:"subnetId" type:"string"`
}

type InstancePrivateIpAddress struct {
	Association *InstanceNetworkInterfaceAssociation `locationName:"association" type:"structure"`

	Primary *bool `locationName:"primary" type:"boolean"`

	PrivateDnsName *string `locationName:"privateDnsName" type:"string"`

	PrivateIpAddress *string `locationName:"privateIpAddress" type:"string"`
}

type InstanceState struct {
	Code *int64 `locationName:"code" type:"integer"`

	Name *string `locationName:"name" type:"string" enum:"InstanceStateName"`
}

type InstanceStateChange struct {
	_ struct{} `type:"structure"`

	CurrentState *InstanceState `locationName:"currentState" type:"structure"`

	InstanceId *string `locationName:"instanceId" type:"string"`

	PreviousState *InstanceState `locationName:"previousState" type:"structure"`
}

type InstanceStatus struct {
	_ struct{} `type:"structure"`

	AvailabilityZone *string `locationName:"availabilityZone" type:"string"`

	Events []*InstanceStatusEvent `locationName:"eventsSet" locationNameList:"item" type:"list"`

	InstanceId *string `locationName:"instanceId" type:"string"`

	InstanceState *InstanceState `locationName:"instanceState" type:"structure"`

	InstanceStatus *InstanceStatusSummary `locationName:"instanceStatus" type:"structure"`

	SystemStatus *InstanceStatusSummary `locationName:"systemStatus" type:"structure"`
}

type InstanceStatusDetails struct {
	_ struct{} `type:"structure"`

	ImpairedSince *time.Time `locationName:"impairedSince" type:"timestamp" timestampFormat:"iso8601"`

	Name *string `locationName:"name" type:"string" enum:"StatusName"`

	Status *string `locationName:"status" type:"string" enum:"StatusType"`
}

type InstanceStatusEvent struct {
	_ struct{} `type:"structure"`

	Code *string `locationName:"code" type:"string" enum:"EventCode"`

	Description *string `locationName:"description" type:"string"`

	NotAfter *time.Time `locationName:"notAfter" type:"timestamp" timestampFormat:"iso8601"`

	NotBefore *time.Time `locationName:"notBefore" type:"timestamp" timestampFormat:"iso8601"`
}

type InstanceStatusSummary struct {
	_ struct{} `type:"structure"`

	Details []*InstanceStatusDetails `locationName:"details" locationNameList:"item" type:"list"`

	Status *string `locationName:"status" type:"string" enum:"SummaryStatus"`
}

type EbsInstanceBlockDevice struct {
	AttachTime *time.Time `locationName:"attachTime" type:"timestamp" timestampFormat:"iso8601"`

	DeleteOnTermination *bool `locationName:"deleteOnTermination" type:"boolean"`

	Status *string `locationName:"status" type:"string" enum:"AttachmentStatus"`

	VolumeId *string `locationName:"volumeId" type:"string"`
}

type EbsInstanceBlockDeviceSpecification struct {
	_ struct{} `type:"structure"`

	DeleteOnTermination *bool `locationName:"deleteOnTermination" type:"boolean"`

	VolumeId *string `locationName:"volumeId" type:"string"`
}

type IamInstanceProfile struct {
	Arn *string `locationName:"arn" type:"string"`

	Id *string `locationName:"id" type:"string"`
}

type Monitoring struct {
	_ struct{} `type:"structure"`

	State *string `locationName:"state" type:"string" enum:"MonitoringState"`
}

type Placement struct {
	Affinity *string `locationName:"affinity" type:"string"`

	AvailabilityZone *string `locationName:"availabilityZone" type:"string"`

	GroupName *string `locationName:"groupName" type:"string"`

	HostId *string `locationName:"hostId" type:"string"`

	Tenancy *string `locationName:"tenancy" type:"string" enum:"Tenancy"`
}

type ProductCode struct {
	ProductCode *string `locationName:"productCode" type:"string"`

	Type *string `locationName:"type" type:"string" enum:"ProductCodeValues"`
}

type StateReason struct {
	Code    *string `locationName:"code" type:"string"`
	Message *string `locationName:"message" type:"string"`
}

type Tag struct {
	_ struct{} `type:"structure"`

	Key *string `locationName:"key" type:"string"`

	Value *string `locationName:"value" type:"string"`
}

type PrivateIpAddressSpecification struct {
	_ struct{} `type:"structure"`

	Primary *bool `locationName:"primary" type:"boolean"`

	PrivateIpAddress *string `locationName:"privateIpAddress" type:"string" required:"true"`
}

type DescribeInstanceAttributeInput struct {
	_ struct{} `type:"structure"`

	Attribute *string `locationName:"attribute" type:"string" required:"true" enum:"InstanceAttributeName"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	InstanceId *string `locationName:"instanceId" type:"string" required:"true"`
}

type DescribeInstanceAttributeOutput struct {
	_ struct{} `type:"structure"`

	BlockDeviceMappings []*InstanceBlockDeviceMapping `locationName:"blockDeviceMapping" locationNameList:"item" type:"list"`

	DisableApiTermination *AttributeBooleanValue `locationName:"disableApiTermination" type:"structure"`

	EbsOptimized *AttributeBooleanValue `locationName:"ebsOptimized" type:"structure"`

	EnaSupport *AttributeBooleanValue `locationName:"enaSupport" type:"structure"`

	Groups []*GroupIdentifier `locationName:"groupSet" locationNameList:"item" type:"list"`

	InstanceId *string `locationName:"instanceId" type:"string"`

	InstanceInitiatedShutdownBehavior *AttributeValue `locationName:"instanceInitiatedShutdownBehavior" type:"structure"`

	InstanceType *AttributeValue `locationName:"instanceType" type:"structure"`

	KernelId *AttributeValue `locationName:"kernel" type:"structure"`

	ProductCodes []*ProductCode `locationName:"productCodes" locationNameList:"item" type:"list"`

	RamdiskId *AttributeValue `locationName:"ramdisk" type:"structure"`

	RootDeviceName *AttributeValue `locationName:"rootDeviceName" type:"structure"`

	SourceDestCheck *AttributeBooleanValue `locationName:"sourceDestCheck" type:"structure"`

	SriovNetSupport *AttributeValue `locationName:"sriovNetSupport" type:"structure"`

	UserData *AttributeValue `locationName:"userData" type:"structure"`

	RequestId *string `locationName:"requestId" type:"string"`
}

type AttributeBooleanValue struct {
	_ struct{} `type:"structure"`

	Value *bool `locationName:"value" type:"boolean"`
}

type AttributeValue struct {
	_ struct{} `type:"structure"`

	Value *string `locationName:"value" type:"string"`
}

type RunInstancesInput struct {
	_ struct{} `type:"structure"`

	BlockDeviceMappings []*BlockDeviceMapping `locationName:"BlockDeviceMapping" locationNameList:"BlockDeviceMapping" type:"list"`

	ClientToken *string `locationName:"clientToken" type:"string"`

	DisableApiTermination *bool `locationName:"disableApiTermination" type:"boolean"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	EbsOptimized *bool `locationName:"ebsOptimized" type:"boolean"`

	ImageId *string `type:"string"`

	InstanceInitiatedShutdownBehavior *string `locationName:"instanceInitiatedShutdownBehavior" type:"string" enum:"ShutdownBehavior"`

	InstanceType *string `type:"string" enum:"InstanceType"`

	InstanceName *string `type:"string" enum:"InstanceName"`

	KeyName *string `locationName:"keyName" type:"string"`

	MaxCount *int64 `type:"integer" required:"true"`

	MinCount *int64 `type:"integer" required:"true"`

	NetworkInterfaces []*InstanceNetworkInterfaceSpecification `locationName:"networkInterface" locationNameList:"item" type:"list"`

	Placement *Placement `type:"structure"`

	PrivateIPAddress *string `locationName:"privateIpAddress" type:"string"`

	PrivateIPAddresses *string `locationName:"privateIpAddresses" type:"string"`

	RamdiskId *string `type:"string"`

	SecurityGroupIds []*string `locationName:"SecurityGroupId" locationNameList:"SecurityGroupId" type:"list"`

	SecurityGroups []*string `locationName:"SecurityGroup" locationNameList:"SecurityGroup" type:"list"`

	SubnetId *string `type:"string"`

	TagSpecifications []*TagSpecification `locationName:"TagSpecification" locationNameList:"item" type:"list"`

	UserData *string `type:"string"`

	OwnerId *string `type:"string"`

	RequestId *string `locationName:"requestId" type:"string"`

	ReservationId *string `type:"string"`

	PasswordData *string `type:"string"`
}

type BlockDeviceMapping struct {
	_ struct{} `type:"structure"`

	DeviceName *string `locationName:"deviceName" type:"string"`

	Ebs *EbsBlockDevice `locationName:"ebs" type:"structure"`

	NoDevice *string `locationName:"noDevice" type:"string"`

	VirtualName *string `locationName:"virtualName" type:"string"`
}

type PrivateIPAddressSpecification struct {
	_ struct{} `type:"structure"`

	Primary *bool `locationName:"primary" type:"boolean"`

	PrivateIPAddress *string `locationName:"privateIpAddress" type:"string" required:"true"`
}

type ModifyInstanceKeyPairInput struct {
	_ struct{} `type:"structure"`

	InstanceId *string `locationName:"instanceId" type:"string"`

	KeyName *string `locationName:"keyName" type:"string"`
}

type EbsBlockDevice struct {
	_ struct{} `type:"structure"`

	DeleteOnTermination *bool `locationName:"deleteOnTermination" type:"boolean"`

	Encrypted *bool `locationName:"encrypted" type:"boolean"`

	Iops *int64 `locationName:"iops" type:"integer"`

	KmsKeyId *string `type:"string"`

	SnapshotId *string `locationName:"snapshotId" type:"string"`

	VolumeSize *int64 `locationName:"volumeSize" type:"integer"`

	VolumeType *string `locationName:"volumeType" type:"string" enum:"VolumeType"`
}

type GetPasswordDataInput struct {
	_ struct{} `type:"structure"`

	InstanceId *string `type:"string" required:"true"`
}

type GetPasswordDataOutput struct {
	_ struct{} `type:"structure"`

	InstanceId *string `locationName:"instanceId" type:"string"`

	PasswordData *string `locationName:"passwordData" type:"string"`

	Timestamp *time.Time `locationName:"timestamp" type:"timestamp" timestampFormat:"iso8601"`
}

type TerminateInstancesInput struct {
	InstanceIds []*string `locationName:"InstanceId" locationNameList:"InstanceId" type:"list" required:"true"`
}

type TerminateInstancesOutput struct {
	_ struct{} `type:"structure"`

	TerminatingInstances []*InstanceStateChange `locationName:"instancesSet" locationNameList:"item" type:"list"`
}
type PublicIP struct {
	AllocationId             *string `locationName:"allocationId" type:"string"`
	AssociationId            *string `locationName:"associationId" type:"string"`
	Domain                   *string `locationName:"domain" type:"string"`
	InstanceId               *string `locationName:"instanceId" type:"string"`
	NetworkInterfaceId       *string `locationName:"networkInterfaceId" type:"string"`
	NetworkInterface_ownerId *string `locationName:"networkInterface_ownerId" type:"string"`
	PrivateIpAddress         *string `locationName:"privateIpAddress" type:"string"`
	PublicIp                 *string `locationName:"publicIp" type:"string"`
}

type AllocateAddressInput struct {
	_ struct{} `type:"structure"`

	Domain *string `type:"string" enum:"DomainType"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`
}

type AllocateAddressOutput struct {
	_ struct{} `type:"structure"`

	AllocationId *string `locationName:"allocationId" type:"string"`

	Domain *string `locationName:"domain" type:"string" enum:"DomainType"`

	PublicIp *string `locationName:"publicIp" type:"string"`
}

type DescribeAddressesInput struct {
	_ struct{} `type:"structure"`

	AllocationIds []*string `locationName:"AllocationId" locationNameList:"AllocationId" type:"list"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	PublicIps []*string `locationName:"PublicIp" locationNameList:"PublicIp" type:"list"`
}

type DescribeAddressesOutput struct {
	_ struct{} `type:"structure"`

	Addresses []*Address `locationName:"addressesSet" locationNameList:"item" type:"list"`

	RequestId *string `locationName:"requestId" type:"string"`
}

func (s DescribeAddressesOutput) String() string {
	return awsutil.Prettify(s)
}

func (s DescribeAddressesOutput) GoString() string {
	return s.String()
}

func (s *DescribeAddressesOutput) SetAddresses(v []*Address) *DescribeAddressesOutput {
	s.Addresses = v
	return s
}

func (s *DescribeAddressesOutput) SetRequestId(v string) *DescribeAddressesOutput {
	s.RequestId = &v
	return s
}

type Address struct {
	_ struct{} `type:"structure"`

	AllocationId *string `locationName:"allocationId" type:"string"`

	AssociationId *string `locationName:"associationId" type:"string"`

	AllowReassociation *bool `locationName:"allowReassociation" type:"bool"`

	Domain *string `locationName:"domain" type:"string" enum:"DomainType"`

	InstanceId *string `locationName:"instanceId" type:"string"`

	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string"`

	NetworkInterfaceOwnerId *string `locationName:"networkInterfaceOwnerId" type:"string"`

	PrivateIpAddress *string `locationName:"privateIpAddress" type:"string"`

	PublicIp *string `locationName:"publicIp" type:"string"`
}

type ModifyInstanceAttributeInput struct {
	_ struct{} `type:"structure"`

	Attribute *string `locationName:"attribute" type:"string" enum:"InstanceAttributeName"`

	BlockDeviceMappings []*BlockDeviceMapping `locationName:"blockDeviceMapping" locationNameList:"item" type:"list"`

	DisableApiTermination *AttributeBooleanValue `locationName:"disableApiTermination" type:"structure"`

	DeleteOnTermination *AttributeBooleanValue `locationName:"deleteOnTermination" type:"structure"`

	EbsOptimized *AttributeBooleanValue `locationName:"ebsOptimized" type:"structure"`

	Groups []*string `locationName:"GroupId" locationNameList:"groupId" type:"list"`

	InstanceId *string `locationName:"instanceId" type:"string" required:"true"`

	InstanceInitiatedShutdownBehavior *AttributeValue `locationName:"instanceInitiatedShutdownBehavior" type:"structure"`

	InstanceType *AttributeValue `locationName:"instanceType" type:"structure"`

	SourceDestCheck *AttributeBooleanValue `type:"structure"`

	UserData *BlobAttributeValue `locationName:"userData" type:"structure"`

	Value *string `locationName:"value" type:"string"`
}

func (s ModifyInstanceAttributeInput) String() string {
	return awsutil.Prettify(s)
}

func (s ModifyInstanceAttributeInput) GoString() string {
	return s.String()
}

func (s *ModifyInstanceAttributeInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "ModifyInstanceAttributeInput"}
	if s.InstanceId == nil {
		invalidParams.Add(request.NewErrParamRequired("InstanceId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

func (s *ModifyInstanceAttributeInput) SetAttribute(v string) *ModifyInstanceAttributeInput {
	s.Attribute = &v
	return s
}

func (s *ModifyInstanceAttributeInput) SetBlockDeviceMappings(v []*BlockDeviceMapping) *ModifyInstanceAttributeInput {
	s.BlockDeviceMappings = v
	return s
}

func (s *ModifyInstanceAttributeInput) SetDisableApiTermination(v *AttributeBooleanValue) *ModifyInstanceAttributeInput {
	s.DisableApiTermination = v
	return s
}

func (s *ModifyInstanceAttributeInput) SetEbsOptimized(v *AttributeBooleanValue) *ModifyInstanceAttributeInput {
	s.EbsOptimized = v
	return s
}

func (s *ModifyInstanceAttributeInput) SetGroups(v []*string) *ModifyInstanceAttributeInput {
	s.Groups = v
	return s
}

func (s *ModifyInstanceAttributeInput) SetInstanceId(v string) *ModifyInstanceAttributeInput {
	s.InstanceId = &v
	return s
}

func (s *ModifyInstanceAttributeInput) SetInstanceInitiatedShutdownBehavior(v *AttributeValue) *ModifyInstanceAttributeInput {
	s.InstanceInitiatedShutdownBehavior = v
	return s
}

func (s *ModifyInstanceAttributeInput) SetInstanceType(v *AttributeValue) *ModifyInstanceAttributeInput {
	s.InstanceType = v
	return s
}

func (s *ModifyInstanceAttributeInput) SetSourceDestCheck(v *AttributeBooleanValue) *ModifyInstanceAttributeInput {
	s.SourceDestCheck = v
	return s
}

func (s *ModifyInstanceAttributeInput) SetUserData(v *BlobAttributeValue) *ModifyInstanceAttributeInput {
	s.UserData = v
	return s
}

func (s *ModifyInstanceAttributeInput) SetValue(v string) *ModifyInstanceAttributeInput {
	s.Value = &v
	return s
}

type BlobAttributeValue struct {
	_ struct{} `type:"structure"`

	Value []byte `locationName:"value" type:"blob"`
}

type StopInstancesInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	Force *bool `locationName:"force" type:"boolean"`

	InstanceIds []*string `locationName:"InstanceId" locationNameList:"InstanceId" type:"list" required:"true"`
}

type StopInstancesOutput struct {
	_ struct{} `type:"structure"`

	StoppingInstances []*InstanceStateChange `locationName:"instancesSet" locationNameList:"item" type:"list"`
}
type ModifyInstanceAttributeOutput struct {
	_ struct{} `type:"structure"`
}

type StartInstancesInput struct {
	_ struct{} `type:"structure"`

	AdditionalInfo *string `locationName:"additionalInfo" type:"string"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	InstanceIds []*string `locationName:"InstanceId" locationNameList:"InstanceId" type:"list" required:"true"`
}

type StartInstancesOutput struct {
	_ struct{} `type:"structure"`

	StartingInstances []*InstanceStateChange `locationName:"instancesSet" locationNameList:"item" type:"list"`
}

type AssociateAddressInput struct {
	_ struct{} `type:"structure"`

	AllocationId *string `type:"string"`

	AllowReassociation *bool `locationName:"allowReassociation" type:"boolean"`

	InstanceId *string `type:"string"`

	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string"`

	PrivateIpAddress *string `locationName:"privateIpAddress" type:"string"`

	PublicIp *string `type:"string"`
}

type AssociateAddressOutput struct {
	_ struct{} `type:"structure"`

	AssociationId *string `locationName:"associationId" type:"string"`

	RequestId *string `locationName:"requestId" type:"string"`
}

type DisassociateAddressInput struct {
	_ struct{} `type:"structure"`

	AssociationId *string `type:"string"`

	PublicIp *string `type:"string"`
}

type DisassociateAddressOutput struct {
	_ struct{} `type:"structure"`

	RequestId *string `locationName:"requestId" type:"string"`
	Return    *bool   `locationName:"return" type:"boolean"`
}

type ReleaseAddressInput struct {
	_ struct{} `type:"structure"`

	AllocationId *string `type:"string"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	PublicIp *string `type:"string"`
}

type ReleaseAddressOutput struct {
	_ struct{} `type:"structure"`
}
type RegisterImageInput struct {
	_ struct{} `type:"structure"`

	Architecture *string `locationName:"architecture" type:"string" enum:"ArchitectureValues"`

	BillingProducts []*string `locationName:"BillingProduct" locationNameList:"item" type:"list"`

	BlockDeviceMappings []*BlockDeviceMapping `locationName:"BlockDeviceMapping" locationNameList:"BlockDeviceMapping" type:"list"`

	Description *string `locationName:"description" type:"string"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	EnaSupport *bool `locationName:"enaSupport" type:"boolean"`

	ImageLocation *string `type:"string"`

	InstanceId *string `type:"string"`

	NoReboot *bool `type:"boolean"`

	KernelId *string `locationName:"kernelId" type:"string"`

	Name *string `locationName:"name" type:"string" required:"true"`

	RamdiskId *string `locationName:"ramdiskId" type:"string"`

	RootDeviceName *string `locationName:"rootDeviceName" type:"string"`

	SriovNetSupport *string `locationName:"sriovNetSupport" type:"string"`

	VirtualizationType *string `locationName:"virtualizationType" type:"string"`
}

type RegisterImageOutput struct {
	_ struct{} `type:"structure"`

	ImageId *string `locationName:"imageId" type:"string"`
}

type DeregisterImageInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	ImageId *string `type:"string" required:"true"`
}

type DeregisterImageOutput struct {
	_ struct{} `type:"structure"`
}

type Image struct {
	_ struct{} `type:"structure"`

	Architecture *string `locationName:"architecture" type:"string" enum:"ArchitectureValues"`

	ClientToken *string `locationName:"clientToken" type:"string"`

	BlockDeviceMappings []*BlockDeviceMapping `locationName:"blockDeviceMapping" locationNameList:"item" type:"list"`

	CreationDate *string `locationName:"creationDate" type:"string"`

	Description *string `locationName:"description" type:"string"`

	EnaSupport *bool `locationName:"enaSupport" type:"boolean"`

	Hypervisor *string `locationName:"hypervisor" type:"string" enum:"HypervisorType"`

	ImageId *string `locationName:"imageId" type:"string"`

	ImageLocation *string `locationName:"imageLocation" type:"string"`

	ImageOwnerAlias *string `locationName:"imageOwnerAlias" type:"string"`

	ImageType *string `locationName:"imageType" type:"string" enum:"ImageTypeValues"`

	KernelId *string `locationName:"kernelId" type:"string"`

	Name *string `locationName:"name" type:"string"`

	OwnerId *string `locationName:"imageOwnerId" type:"string"`

	Platform *string `locationName:"platform" type:"string" enum:"PlatformValues"`

	ProductCodes []*ProductCode `locationName:"productCodes" locationNameList:"item" type:"list"`

	Public *bool `locationName:"isPublic" type:"boolean"`

	RamdiskId *string `locationName:"ramdiskId" type:"string"`

	RootDeviceName *string `locationName:"rootDeviceName" type:"string"`

	RootDeviceType *string `locationName:"rootDeviceType" type:"string" enum:"DeviceType"`

	SriovNetSupport *string `locationName:"sriovNetSupport" type:"string"`

	State *string `locationName:"imageState" type:"string" enum:"ImageState"`

	StateReason *StateReason `locationName:"stateReason" type:"structure"`

	Tags []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`

	VirtualizationType *string `locationName:"virtualizationType" type:"string" enum:"VirtualizationType"`
}

type DescribeImagesInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	ExecutableUsers []*string `locationName:"ExecutableBy" locationNameList:"ExecutableBy" type:"list"`

	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	ImageIds []*string `locationName:"ImageId" locationNameList:"ImageId" type:"list"`

	Owners []*string `locationName:"Owner" locationNameList:"Owner" type:"list"`
}

type DescribeImagesOutput struct {
	_ struct{} `type:"structure"`

	Images []*Image `locationName:"imagesSet" locationNameList:"item" type:"list"`

	RequestId *string `locationName:"requestId" type:"String"`
}

func (s DescribeImagesOutput) String() string {
	return awsutil.Prettify(s)
}

func (s DescribeImagesOutput) GoString() string {
	return s.String()
}

func (s *DescribeImagesOutput) SetImages(v []*Image) *DescribeImagesOutput {
	s.Images = v
	return s
}
func (s *DescribeImagesOutput) SetRequestId(v *string) *DescribeImagesOutput {
	s.RequestId = v
	return s
}

type ModifyImageAttributeInput struct {
	_ struct{} `type:"structure"`

	Attribute *string `type:"string"`

	Description *AttributeValue `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	ImageId *string `type:"string" required:"true"`

	LaunchPermission *LaunchPermissionModifications `type:"structure"`

	OperationType *string `type:"string" enum:"OperationType"`

	ProductCodes []*string `locationName:"ProductCode" locationNameList:"ProductCode" type:"list"`

	UserGroups []*string `locationName:"UserGroup" locationNameList:"UserGroup" type:"list"`

	UserIds []*string `locationName:"UserId" locationNameList:"UserId" type:"list"`

	Value *string `type:"string"`
}

type ModifyImageAttributeOutput struct {
	_ struct{} `type:"structure"`
}

type LaunchPermissionModifications struct {
	_ struct{} `type:"structure"`

	Add []*LaunchPermission `locationNameList:"item" type:"list"`

	Remove []*LaunchPermission `locationNameList:"item" type:"list"`
}

type LaunchPermission struct {
	_ struct{} `type:"structure"`

	Group *string `locationName:"group" type:"string" enum:"PermissionGroup"`

	UserId *string `locationName:"userId" type:"string"`
}

type DeleteTagsInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	Resources []*string `locationName:"resourceId" type:"list" required:"true"`

	Tags []*Tag `locationName:"tag" locationNameList:"item" type:"list"`
}

type DeleteTagsOutput struct {
	_ struct{} `type:"structure"`
}

type CreateTagsInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	Resources []*string `locationName:"ResourceId" type:"list" required:"true"`

	Tags []*Tag `locationName:"Tag" locationNameList:"item" type:"list" required:"true"`
}

type CreateTagsOutput struct {
	_ struct{} `type:"structure"`
}

type DescribeTagsInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	MaxResults *int64 `locationName:"maxResults" type:"integer"`

	NextToken *string `locationName:"nextToken" type:"string"`
}

type DescribeTagsOutput struct {
	_ struct{} `type:"structure"`

	// The token to use to retrieve the next page of results. This value is null
	// when there are no more results to return..
	NextToken *string `locationName:"nextToken" type:"string"`

	// A list of tags.
	Tags []*TagDescription `locationName:"tagSet" locationNameList:"item" type:"list"`

	RequestId *string `locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s DescribeTagsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeTagsOutput) GoString() string {
	return s.String()
}

// SetNextToken sets the NextToken field's value.
func (s *DescribeTagsOutput) SetNextToken(v string) *DescribeTagsOutput {
	s.NextToken = &v
	return s
}
func (s *DescribeTagsOutput) SetRequestId(v string) *DescribeTagsOutput {
	s.RequestId = &v
	return s
}

// SetTags sets the Tags field's value.
func (s *DescribeTagsOutput) SetTags(v []*TagDescription) *DescribeTagsOutput {
	s.Tags = v
	return s
}

type TagDescription struct {
	_ struct{} `type:"structure"`

	Key *string `locationName:"key" type:"string"`

	ResourceId *string `locationName:"resourceId" type:"string"`

	ResourceType *string `locationName:"resourceType" type:"string" enum:"ResourceType"`

	Value *string `locationName:"value" type:"string"`
}

type TagSpecification struct {
	_ struct{} `type:"structure"`

	ResourceType *string `locationName:"resourceType" type:"string" enum:"ResourceType"`

	Tags []*Tag `locationName:"Tag" locationNameList:"item" type:"list"`
}

// Contains the parameters for ImportKeyPair.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/ImportKeyPairRequest
type ImportKeyPairInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// A unique name for the key pair.
	//
	// KeyName is a required field
	KeyName *string `locationName:"keyName" type:"string" required:"true"`

	// The public key. For API calls, the text must be base64-encoded. For command
	// line tools, base64 encoding is performed for you.
	//
	// PublicKeyMaterial is automatically base64 encoded/decoded by the SDK.
	//
	// PublicKeyMaterial is a required field
	PublicKeyMaterial []byte `locationName:"publicKeyMaterial" type:"blob" required:"true"`
}

// String returns the string representation

// SetDryRun sets the DryRun field's value.

// Contains the output of ImportKeyPair.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/ImportKeyPairResult
type ImportKeyPairOutput struct {
	_ struct{} `type:"structure"`

	// The MD5 public key fingerprint as specified in section 4 of RFC 4716.
	KeyFingerprint *string `locationName:"keyFingerprint" type:"string"`

	// The key pair name you provided.
	KeyName *string `locationName:"keyName" type:"string"`
}

// Contains the parameters for DescribeKeyPairs.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeKeyPairsRequest
type DescribeKeyPairsInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more filters.
	//
	//    * fingerprint - The fingerprint of the key pair.
	//
	//    * key-name - The name of the key pair.
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// One or more key pair names.
	//
	// Default: Describes all your key pairs.
	KeyNames []*string `locationName:"KeyName" locationNameList:"KeyName" type:"list"`
}

// Contains the output of DescribeKeyPairs.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeKeyPairsResult
type DescribeKeyPairsOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more key pairs.
	KeyPairs  []*KeyPairInfo `locationName:"keySet" locationNameList:"item" type:"list"`
	RequestId *string        `locationName:"requestId" type:"String"`
}

// String returns the string representation

// Describes a key pair.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/KeyPairInfo
type KeyPairInfo struct {
	_ struct{} `type:"structure"`

	// If you used CreateKeyPair to create the key pair, this is the SHA-1 digest
	// of the DER encoded private key. If you used ImportKeyPair to provide AWS
	// the public key, this is the MD5 public key fingerprint as specified in section
	// 4 of RFC4716.
	KeyFingerprint *string `locationName:"keyFingerprint" type:"string"`

	// The name of the key pair.
	KeyName *string `locationName:"keyName" type:"string"`
}

// String returns the string representation

// Contains the parameters for DeleteKeyPair.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteKeyPairRequest
type DeleteKeyPairInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The name of the key pair.
	//
	// KeyName is a required field
	KeyName *string `type:"string" required:"true"`
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteKeyPairOutput
type DeleteKeyPairOutput struct {
	_ struct{} `type:"structure"`
}

type CreateKeyPairInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// A unique name for the key pair.
	//
	// Constraints: Up to 255 ASCII characters
	//
	// KeyName is a required field
	KeyName *string `type:"string" required:"true"`
}

// Describes a key pair.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/KeyPair
type CreateKeyPairOutput struct {
	_ struct{} `type:"structure"`

	// The SHA-1 digest of the DER encoded private key.
	KeyFingerprint *string `locationName:"keyFingerprint" type:"string"`

	// An unencrypted PEM encoded RSA private key.
	KeyMaterial *string `locationName:"keyMaterial" type:"string"`

	// The name of the key pair.
	KeyName *string `locationName:"keyName" type:"string"`

	// The name of the Request ID
	RequestId *string `locationName:"requestId" type:"String"`
}

type CreateSecurityGroupInput struct {
	_ struct{} `type:"structure"`

	Description *string `locationName:"GroupDescription" type:"string" required:"true"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	GroupName *string `type:"string" required:"true"`

	VpcId *string `type:"string"`
}

type CreateSecurityGroupOutput struct {
	_ struct{} `type:"structure"`

	GroupId *string `locationName:"groupId" type:"string"`
}

type DescribeSecurityGroupsInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	GroupIds []*string `locationName:"GroupId" locationNameList:"groupId" type:"list"`

	GroupNames []*string `locationName:"GroupName" locationNameList:"GroupName" type:"list"`
}

type DescribeSecurityGroupsOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more security groups.
	SecurityGroups []*SecurityGroup `locationName:"securityGroupInfo" locationNameList:"item" type:"list"`

	RequestId *string `locationName:"requestId" type:"String"`
}

// String returns the string representation
func (s DescribeSecurityGroupsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeSecurityGroupsOutput) GoString() string {
	return s.String()
}

// SetSecurityGroups sets the SecurityGroups field's value.
func (s *DescribeSecurityGroupsOutput) SetSecurityGroups(v []*SecurityGroup) *DescribeSecurityGroupsOutput {
	s.SecurityGroups = v
	return s
}

func (s *DescribeSecurityGroupsOutput) SetRequestId(v string) *DescribeSecurityGroupsOutput {
	s.RequestId = &v
	return s
}

type SecurityGroup struct {
	_                   struct{}        `type:"structure"`
	Description         *string         `locationName:"groupDescription" type:"string"`
	GroupId             *string         `locationName:"groupId" type:"string"`
	GroupName           *string         `locationName:"groupName" type:"string"`
	IpPermissions       []*IpPermission `locationName:"ipPermissions" locationNameList:"item" type:"list"`
	IpPermissionsEgress []*IpPermission `locationName:"ipPermissionsEgress" locationNameList:"item" type:"list"`
	OwnerId             *string         `locationName:"ownerId" type:"string"`
	Tags                []*Tag          `locationName:"tagSet" locationNameList:"item" type:"list"`
	VpcId               *string         `locationName:"vpcId" type:"string"`
}

type IpPermission struct {
	_                struct{}           `type:"structure"`
	FromPort         *int64             `locationName:"fromPort" type:"integer"`
	IpProtocol       *string            `locationName:"ipProtocol" type:"string"`
	IpRanges         []*IpRange         `locationName:"ipRanges" locationNameList:"item" type:"list"`
	Ipv6Ranges       []*Ipv6Range       `locationName:"ipv6Ranges" locationNameList:"item" type:"list"`
	PrefixListIds    []*PrefixListId    `locationName:"prefixListIds" locationNameList:"item" type:"list"`
	ToPort           *int64             `locationName:"toPort" type:"integer"`
	UserIdGroupPairs []*UserIdGroupPair `locationName:"groups" locationNameList:"item" type:"list"`
}

type IpRange struct {
	_ struct{} `type:"structure"`

	CidrIp *string `locationName:"cidrIp" type:"string"`
}

type Ipv6Range struct {
	_ struct{} `type:"structure"`

	CidrIpv6 *string `locationName:"cidrIpv6" type:"string"`
}

type PrefixListId struct {
	_ struct{} `type:"structure"`

	PrefixListId *string `locationName:"prefixListId" type:"string"`
}

type UserIdGroupPair struct {
	_ struct{} `type:"structure"`

	GroupId *string `locationName:"groupId" type:"string"`

	GroupName *string `locationName:"groupName" type:"string"`

	PeeringStatus *string `locationName:"peeringStatus" type:"string"`

	UserId *string `locationName:"userId" type:"string"`

	VpcId *string `locationName:"vpcId" type:"string"`

	VpcPeeringConnectionId *string `locationName:"vpcPeeringConnectionId" type:"string"`
}

type RevokeSecurityGroupEgressInput struct {
	_ struct{} `type:"structure"`

	CidrIp *string `locationName:"cidrIp" type:"string"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	FromPort *int64 `locationName:"fromPort" type:"integer"`

	GroupId *string `locationName:"groupId" type:"string" required:"true"`

	IpPermissions []*IpPermission `locationName:"ipPermissions" locationNameList:"item" type:"list"`

	IpProtocol *string `locationName:"ipProtocol" type:"string"`

	SourceSecurityGroupName *string `locationName:"sourceSecurityGroupName" type:"string"`

	SourceSecurityGroupOwnerId *string `locationName:"sourceSecurityGroupOwnerId" type:"string"`

	ToPort *int64 `locationName:"toPort" type:"integer"`
}

type RevokeSecurityGroupEgressOutput struct {
	_ struct{} `type:"structure"`
}

type RevokeSecurityGroupIngressInput struct {
	_ struct{} `type:"structure"`

	CidrIp *string `type:"string"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	FromPort *int64 `type:"integer"`

	GroupId *string `type:"string"`

	GroupName *string `type:"string"`

	IpPermissions []*IpPermission `locationNameList:"item" type:"list"`

	IpProtocol *string `type:"string"`

	SourceSecurityGroupName *string `type:"string"`

	SourceSecurityGroupOwnerId *string `type:"string"`

	ToPort *int64 `type:"integer"`
}

type RevokeSecurityGroupIngressOutput struct {
	_ struct{} `type:"structure"`
}

type AuthorizeSecurityGroupEgressInput struct {
	_ struct{} `type:"structure"`

	CidrIp *string `locationName:"cidrIp" type:"string"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	FromPort *int64 `locationName:"fromPort" type:"integer"`

	GroupId *string `locationName:"groupId" type:"string" required:"true"`

	IpPermissions []*IpPermission `locationName:"ipPermissions" locationNameList:"item" type:"list"`

	IpProtocol *string `locationName:"ipProtocol" type:"string"`

	SourceSecurityGroupName *string `locationName:"sourceSecurityGroupName" type:"string"`

	SourceSecurityGroupOwnerId *string `locationName:"sourceSecurityGroupOwnerId" type:"string"`

	ToPort *int64 `locationName:"toPort" type:"integer"`
}

type AuthorizeSecurityGroupEgressOutput struct {
	_ struct{} `type:"structure"`
}

type AuthorizeSecurityGroupIngressInput struct {
	_ struct{} `type:"structure"`

	CidrIp *string `type:"string"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	FromPort *int64 `type:"integer"`

	GroupId *string `type:"string"`

	GroupName *string `type:"string"`

	IpPermissions []*IpPermission `locationNameList:"item" type:"list"`

	IpProtocol *string `type:"string"`

	SourceSecurityGroupName *string `type:"string"`

	SourceSecurityGroupOwnerId *string `type:"string"`

	ToPort *int64 `type:"integer"`
}

type AuthorizeSecurityGroupIngressOutput struct {
	_ struct{} `type:"structure"`
}

type DeleteSecurityGroupInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	GroupId *string `type:"string"`

	GroupName *string `type:"string"`
}

type DeleteSecurityGroupOutput struct {
	_ struct{} `type:"structure"`
}

// Contains the parameters for CreateVolume.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateVolumeRequest
type CreateVolumeInput struct {
	_ struct{} `type:"structure"`

	// The Availability Zone in which to create the volume. Use DescribeAvailabilityZones
	// to list the Availability Zones that are currently available to you.
	//
	// AvailabilityZone is a required field
	AvailabilityZone *string `type:"string" required:"true"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// Specifies whether the volume should be encrypted. Encrypted Amazon EBS volumes
	// may only be attached to instances that support Amazon EBS encryption. Volumes
	// that are created from encrypted snapshots are automatically encrypted. There
	// is no way to create an encrypted volume from an unencrypted snapshot or vice
	// versa. If your AMI uses encrypted volumes, you can only launch it on supported
	// instance types. For more information, see Amazon EBS Encryption (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSEncryption.html)
	// in the Amazon Elastic Compute Cloud User Guide.
	Encrypted *bool `locationName:"encrypted" type:"boolean"`

	// Only valid for Provisioned IOPS SSD volumes. The number of I/O operations
	// per second (IOPS) to provision for the volume, with a maximum ratio of 50
	// IOPS/GiB.
	//
	// Constraint: Range is 100 to 20000 for Provisioned IOPS SSD volumes
	Iops *int64 `type:"integer"`

	// The full ARN of the AWS Key Management Service (AWS KMS) customer master
	// key (CMK) to use when creating the encrypted volume. This parameter is only
	// required if you want to use a non-default CMK; if this parameter is not specified,
	// the default CMK for EBS is used. The ARN contains the arn:aws:kms namespace,
	// followed by the region of the CMK, the AWS account ID of the CMK owner, the
	// key namespace, and then the CMK ID. For example, arn:aws:kms:us-east-1:012345678910:key/abcd1234-a123-456a-a12b-a123b4cd56ef.
	// If a KmsKeyId is specified, the Encrypted flag must also be set.
	KmsKeyId *string `type:"string"`

	// The size of the volume, in GiBs.
	//
	// Constraints: 1-16384 for gp2, 4-16384 for io1, 500-16384 for st1, 500-16384
	// for sc1, and 1-1024 for standard. If you specify a snapshot, the volume size
	// must be equal to or larger than the snapshot size.
	//
	// Default: If you're creating the volume from a snapshot and don't specify
	// a volume size, the default is the snapshot size.
	Size *int64 `type:"integer"`

	// The snapshot from which to create the volume.
	SnapshotId *string `type:"string"`

	// The tags to apply to the volume during creation.
	TagSpecifications []*TagSpecification `locationName:"TagSpecification" locationNameList:"item" type:"list"`

	// The volume type. This can be gp2 for General Purpose SSD, io1 for Provisioned
	// IOPS SSD, st1 for Throughput Optimized HDD, sc1 for Cold HDD, or standard
	// for Magnetic volumes.
	//
	// Default: standard
	VolumeType *string `type:"string" enum:"VolumeType"`
}

// Contains the parameters for DeleteVolume.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteVolumeRequest
type DeleteVolumeInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the volume.
	//
	// VolumeId is a required field
	VolumeId *string `type:"string" required:"true"`
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteVolumeOutput
type DeleteVolumeOutput struct {
	_ struct{} `type:"structure"`
}

// Contains the parameters for DescribeVolumes.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeVolumesRequest
type DescribeVolumesInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more filters.
	//
	//    * attachment.attach-time - The time stamp when the attachment initiated.
	//
	//    * attachment.delete-on-termination - Whether the volume is deleted on
	//    instance termination.
	//
	//    * attachment.device - The device name that is exposed to the instance
	//    (for example, /dev/sda1).
	//
	//    * attachment.instance-id - The ID of the instance the volume is attached
	//    to.
	//
	//    * attachment.status - The attachment state (attaching | attached | detaching
	//    | detached).
	//
	//    * availability-zone - The Availability Zone in which the volume was created.
	//
	//    * create-time - The time stamp when the volume was created.
	//
	//    * encrypted - The encryption status of the volume.
	//
	//    * size - The size of the volume, in GiB.
	//
	//    * snapshot-id - The snapshot from which the volume was created.
	//
	//    * status - The status of the volume (creating | available | in-use | deleting
	//    | deleted | error).
	//
	//    * tag:key=value - The key/value combination of a tag assigned to the resource.
	//    Specify the key of the tag in the filter name and the value of the tag
	//    in the filter value. For example, for the tag Purpose=X, specify tag:Purpose
	//    for the filter name and X for the filter value.
	//
	//    * tag-key - The key of a tag assigned to the resource. This filter is
	//    independent of the tag-value filter. For example, if you use both the
	//    filter "tag-key=Purpose" and the filter "tag-value=X", you get any resources
	//    assigned both the tag key Purpose (regardless of what the tag's value
	//    is), and the tag value X (regardless of what the tag's key is). If you
	//    want to list only resources where Purpose is X, see the tag:key=value
	//    filter.
	//
	//    * tag-value - The value of a tag assigned to the resource. This filter
	//    is independent of the tag-key filter.
	//
	//    * volume-id - The volume ID.
	//
	//    * volume-type - The Amazon EBS volume type. This can be gp2 for General
	//    Purpose SSD, io1 for Provisioned IOPS SSD, st1 for Throughput Optimized
	//    HDD, sc1 for Cold HDD, or standard for Magnetic volumes.
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// The maximum number of volume results returned by DescribeVolumes in paginated
	// output. When this parameter is used, DescribeVolumes only returns MaxResults
	// results in a single page along with a NextToken response element. The remaining
	// results of the initial request can be seen by sending another DescribeVolumes
	// request with the returned NextToken value. This value can be between 5 and
	// 500; if MaxResults is given a value larger than 500, only 500 results are
	// returned. If this parameter is not used, then DescribeVolumes returns all
	// results. You cannot specify this parameter and the volume IDs parameter in
	// the same request.
	MaxResults *int64 `locationName:"maxResults" type:"integer"`

	// The NextToken value returned from a previous paginated DescribeVolumes request
	// where MaxResults was used and the results exceeded the value of that parameter.
	// Pagination continues from the end of the previous results that returned the
	// NextToken value. This value is null when there are no more results to return.
	NextToken *string `locationName:"nextToken" type:"string"`

	// One or more volume IDs.
	VolumeIds []*string `locationName:"VolumeId" locationNameList:"VolumeId" type:"list"`
}

type DescribeVolumesOutput struct {
	_ struct{} `type:"structure"`

	// The NextToken value to include in a future DescribeVolumes request. When
	// the results of a DescribeVolumes request exceed MaxResults, this value can
	// be used to retrieve the next page of results. This value is null when there
	// are no more results to return.
	NextToken *string `locationName:"nextToken" type:"string"`

	// Information about the volumes.
	Volumes []*Volume `locationName:"volumeSet" locationNameList:"item" type:"list"`

	RequestId *string `locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s DescribeVolumesOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeVolumesOutput) GoString() string {
	return s.String()
}

// SetNextToken sets the NextToken field's value.
func (s *DescribeVolumesOutput) SetNextToken(v string) *DescribeVolumesOutput {
	s.NextToken = &v
	return s
}
func (s *DescribeVolumesOutput) SetRequesterId(v string) *DescribeVolumesOutput {
	s.RequestId = &v
	return s
}

// SetVolumes sets the Volumes field's value.
func (s *DescribeVolumesOutput) SetVolumes(v []*Volume) *DescribeVolumesOutput {
	s.Volumes = v
	return s
}

// Describes a volume.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/Volume
type Volume struct {
	_ struct{} `type:"structure"`

	// Information about the volume attachments.
	Attachments []*VolumeAttachment `locationName:"attachmentSet" locationNameList:"item" type:"list"`

	// The Availability Zone for the volume.
	AvailabilityZone *string `locationName:"availabilityZone" type:"string"`

	// The time stamp when volume creation was initiated.
	CreateTime *time.Time `locationName:"createTime" type:"timestamp" timestampFormat:"iso8601"`

	// Indicates whether the volume will be encrypted.
	Encrypted *bool `locationName:"encrypted" type:"boolean"`

	Iops *int64 `locationName:"iops" type:"integer"`

	// The full ARN of the AWS Key Management Service (AWS KMS) customer master
	// key (CMK) that was used to protect the volume encryption key for the volume.
	KmsKeyId *string `locationName:"kmsKeyId" type:"string"`

	// The size of the volume, in GiBs.
	Size *int64 `locationName:"size" type:"integer"`

	// The snapshot from which the volume was created, if applicable.
	SnapshotId *string `locationName:"snapshotId" type:"string"`

	// The volume state.
	State *string `locationName:"status" type:"string" enum:"VolumeState"`

	// Any tags assigned to the volume.
	Tags []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`

	// The ID of the volume.
	VolumeId *string `locationName:"volumeId" type:"string"`

	// The volume type. This can be gp2 for General Purpose SSD, io1 for Provisioned
	// IOPS SSD, st1 for Throughput Optimized HDD, sc1 for Cold HDD, or standard
	// for Magnetic volumes.
	VolumeType *string `locationName:"volumeType" type:"string" enum:"VolumeType"`
}

// String returns the string representation
func (s Volume) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Volume) GoString() string {
	return s.String()
}

// SetAttachments sets the Attachments field's value.
func (s *Volume) SetAttachments(v []*VolumeAttachment) *Volume {
	s.Attachments = v
	return s
}

// SetAvailabilityZone sets the AvailabilityZone field's value.
func (s *Volume) SetAvailabilityZone(v string) *Volume {
	s.AvailabilityZone = &v
	return s
}

// SetCreateTime sets the CreateTime field's value.
func (s *Volume) SetCreateTime(v time.Time) *Volume {
	s.CreateTime = &v
	return s
}

// SetEncrypted sets the Encrypted field's value.
func (s *Volume) SetEncrypted(v bool) *Volume {
	s.Encrypted = &v
	return s
}

// SetIops sets the Iops field's value.
func (s *Volume) SetIops(v int64) *Volume {
	s.Iops = &v
	return s
}

// SetKmsKeyId sets the KmsKeyId field's value.
func (s *Volume) SetKmsKeyId(v string) *Volume {
	s.KmsKeyId = &v
	return s
}

// SetSize sets the Size field's value.
func (s *Volume) SetSize(v int64) *Volume {
	s.Size = &v
	return s
}

// SetSnapshotId sets the SnapshotId field's value.
func (s *Volume) SetSnapshotId(v string) *Volume {
	s.SnapshotId = &v
	return s
}

// SetState sets the State field's value.
func (s *Volume) SetState(v string) *Volume {
	s.State = &v
	return s
}

// SetTags sets the Tags field's value.
func (s *Volume) SetTags(v []*Tag) *Volume {
	s.Tags = v
	return s
}

// SetVolumeId sets the VolumeId field's value.
func (s *Volume) SetVolumeId(v string) *Volume {
	s.VolumeId = &v
	return s
}

// SetVolumeType sets the VolumeType field's value.
func (s *Volume) SetVolumeType(v string) *Volume {
	s.VolumeType = &v
	return s
}

type VolumeAttachment struct {
	_ struct{} `type:"structure"`

	// The time stamp when the attachment initiated.
	AttachTime *time.Time `locationName:"attachTime" type:"timestamp" timestampFormat:"iso8601"`

	// Indicates whether the EBS volume is deleted on instance termination.
	DeleteOnTermination *bool `locationName:"deleteOnTermination" type:"boolean"`

	// The device name.
	Device *string `locationName:"device" type:"string"`

	// The ID of the instance.
	InstanceId *string `locationName:"instanceId" type:"string"`

	// The attachment state of the volume.
	State *string `locationName:"status" type:"string" enum:"VolumeAttachmentState"`

	// The ID of the volume.
	VolumeId *string `locationName:"volumeId" type:"string"`
}

// String returns the string representation
func (s VolumeAttachment) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s VolumeAttachment) GoString() string {
	return s.String()
}

// SetAttachTime sets the AttachTime field's value.
func (s *VolumeAttachment) SetAttachTime(v time.Time) *VolumeAttachment {
	s.AttachTime = &v
	return s
}

// SetDeleteOnTermination sets the DeleteOnTermination field's value.
func (s *VolumeAttachment) SetDeleteOnTermination(v bool) *VolumeAttachment {
	s.DeleteOnTermination = &v
	return s
}

// SetDevice sets the Device field's value.
func (s *VolumeAttachment) SetDevice(v string) *VolumeAttachment {
	s.Device = &v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *VolumeAttachment) SetInstanceId(v string) *VolumeAttachment {
	s.InstanceId = &v
	return s
}

// SetState sets the State field's value.
func (s *VolumeAttachment) SetState(v string) *VolumeAttachment {
	s.State = &v
	return s
}

// SetVolumeId sets the VolumeId field's value.
func (s *VolumeAttachment) SetVolumeId(v string) *VolumeAttachment {
	s.VolumeId = &v
	return s
}

type AttachVolumeInput struct {
	_ struct{} `type:"structure"`

	Device *string `type:"string" required:"true"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the instance.
	//
	// InstanceId is a required field
	InstanceId *string `type:"string" required:"true"`

	VolumeId *string `type:"string" required:"true"`
}

type DetachVolumeInput struct {
	_ struct{} `type:"structure"`

	// The device name.
	Device *string `type:"string"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	Force *bool `type:"boolean"`

	// The ID of the instance.
	InstanceId *string `type:"string"`

	VolumeId *string `type:"string" required:"true"`
}
type CreateSubnetInput struct {
	_ struct{} `type:"structure"`

	// The Availability Zone for the subnet.
	//
	// Default: AWS selects one for you. If you create more than one subnet in your
	// VPC, we may not necessarily select a different zone for each subnet.
	AvailabilityZone *string `type:"string"`

	// The IPv4 network range for the subnet, in CIDR notation. For example, 10.0.0.0/24.
	//
	// CidrBlock is a required field
	CidrBlock *string `type:"string" required:"true"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The IPv6 network range for the subnet, in CIDR notation. The subnet size
	// must use a /64 prefix length.
	Ipv6CidrBlock *string `type:"string"`

	// The ID of the VPC.
	//
	// VpcId is a required field
	VpcId *string `type:"string" required:"true"`
}

type CreateSubnetOutput struct {
	_ struct{} `type:"structure"`

	// Information about the subnet.
	Subnet *Subnet `locationName:"subnet" type:"structure"`
}

type DescribeInstanceStatusInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	IncludeAllInstances *bool `locationName:"includeAllInstances" type:"boolean"`

	InstanceIds []*string `locationName:"InstanceId" locationNameList:"InstanceId" type:"list"`

	MaxResults *int64 `type:"integer"`

	NextToken *string `type:"string"`
}

func (s DescribeInstanceStatusInput) String() string {
	return awsutil.Prettify(s)
}

func (s DescribeInstanceStatusInput) GoString() string {
	return s.String()
}

func (s *DescribeInstanceStatusInput) SetDryRun(v bool) *DescribeInstanceStatusInput {
	s.DryRun = &v
	return s
}

func (s *DescribeInstanceStatusInput) SetFilters(v []*Filter) *DescribeInstanceStatusInput {
	s.Filters = v
	return s
}

func (s *DescribeInstanceStatusInput) SetIncludeAllInstances(v bool) *DescribeInstanceStatusInput {
	s.IncludeAllInstances = &v
	return s
}

func (s *DescribeInstanceStatusInput) SetInstanceIds(v []*string) *DescribeInstanceStatusInput {
	s.InstanceIds = v
	return s
}

func (s *DescribeInstanceStatusInput) SetMaxResults(v int64) *DescribeInstanceStatusInput {
	s.MaxResults = &v
	return s
}

func (s *DescribeInstanceStatusInput) SetNextToken(v string) *DescribeInstanceStatusInput {
	s.NextToken = &v
	return s
}

type DescribeInstanceStatusOutput struct {
	_ struct{} `type:"structure"`

	InstanceStatuses []*InstanceStatus `locationName:"instanceStatusSet" locationNameList:"item" type:"list"`

	NextToken *string `locationName:"nextToken" type:"string"`

	RequestId *string `locationName:"requestId" type:"string"`
}

func (s DescribeInstanceStatusOutput) String() string {
	return awsutil.Prettify(s)
}

func (s DescribeInstanceStatusOutput) GoString() string {
	return s.String()
}

func (s *DescribeInstanceStatusOutput) SetInstanceStatuses(v []*InstanceStatus) *DescribeInstanceStatusOutput {
	s.InstanceStatuses = v
	return s
}

func (s *DescribeInstanceStatusOutput) SetNextToken(v string) *DescribeInstanceStatusOutput {
	s.NextToken = &v
	return s
}

//CreateInternetGatewayInput Contains the parameters for CreateInternetGateway.
type CreateInternetGatewayInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the subnet.
	//
	// SubnetId is a required field
	SubnetId *string `type:"string" required:"true"`
}

type DeleteSubnetInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the subnet.
	//
	// SubnetId is a required field
	SubnetId *string `type:"string" required:"true"`
}

// String returns the string representation
func (s DeleteSubnetInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteSubnetInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteSubnetInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeleteSubnetInput"}
	if s.SubnetId == nil {
		invalidParams.Add(request.NewErrParamRequired("SubnetId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDryRun sets the DryRun field's value.
func (s *DeleteSubnetInput) SetDryRun(v bool) *DeleteSubnetInput {
	s.DryRun = &v
	return s
}

// SetSubnetId sets the SubnetId field's value.
func (s *DeleteSubnetInput) SetSubnetId(v string) *DeleteSubnetInput {
	s.SubnetId = &v
	return s
}

type DeleteSubnetOutput struct {
	_ struct{} `type:"structure"`
}

type Subnet struct {
	_ struct{} `type:"structure"`

	// Indicates whether a network interface created in this subnet (including a
	// network interface created by RunInstances) receives an IPv6 address.
	AssignIpv6AddressOnCreation *bool `locationName:"assignIpv6AddressOnCreation" type:"boolean"`

	// The Availability Zone of the subnet.
	AvailabilityZone *string `locationName:"availabilityZone" type:"string"`

	// The number of unused private IPv4 addresses in the subnet. Note that the
	// IPv4 addresses for any stopped instances are considered unavailable.
	AvailableIpAddressCount *int64 `locationName:"availableIpAddressCount" type:"integer"`

	// The IPv4 CIDR block assigned to the subnet.
	CidrBlock *string `locationName:"cidrBlock" type:"string"`

	// Indicates whether this is the default subnet for the Availability Zone.
	DefaultForAz *bool `locationName:"defaultForAz" type:"boolean"`

	// Information about the IPv6 CIDR blocks associated with the subnet.

	// Indicates whether instances launched in this subnet receive a public IPv4
	// address.
	MapPublicIpOnLaunch *bool `locationName:"mapPublicIpOnLaunch" type:"boolean"`

	// The current state of the subnet.
	State *string `locationName:"state" type:"string" enum:"SubnetState"`

	// The ID of the subnet.
	SubnetId *string `locationName:"subnetId" type:"string"`

	// Any tags assigned to the subnet.
	Tags []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`

	// The ID of the VPC the subnet is in.
	VpcId *string `locationName:"vpcId" type:"string"`
}

type DescribeSubnetsInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more filters.
	//
	//    * availabilityZone - The Availability Zone for the subnet. You can also
	//    use availability-zone as the filter name.
	//
	//    * available-ip-address-count - The number of IPv4 addresses in the subnet
	//    that are available.
	//
	//    * cidrBlock - The IPv4 CIDR block of the subnet. The CIDR block you specify
	//    must exactly match the subnet's CIDR block for information to be returned
	//    for the subnet. You can also use cidr or cidr-block as the filter names.
	//
	//    * defaultForAz - Indicates whether this is the default subnet for the
	//    Availability Zone. You can also use default-for-az as the filter name.
	//
	//    * ipv6-cidr-block-association.ipv6-cidr-block - An IPv6 CIDR block associated
	//    with the subnet.
	//
	//    * ipv6-cidr-block-association.association-id - An association ID for an
	//    IPv6 CIDR block associated with the subnet.
	//
	//    * ipv6-cidr-block-association.state - The state of an IPv6 CIDR block
	//    associated with the subnet.
	//
	//    * state - The state of the subnet (pending | available).
	//
	//    * subnet-id - The ID of the subnet.
	//
	//    * tag:key=value - The key/value combination of a tag assigned to the resource.
	//    Specify the key of the tag in the filter name and the value of the tag
	//    in the filter value. For example, for the tag Purpose=X, specify tag:Purpose
	//    for the filter name and X for the filter value.
	//
	//    * tag-key - The key of a tag assigned to the resource. This filter is
	//    independent of the tag-value filter. For example, if you use both the
	//    filter "tag-key=Purpose" and the filter "tag-value=X", you get any resources
	//    assigned both the tag key Purpose (regardless of what the tag's value
	//    is), and the tag value X (regardless of what the tag's key is). If you
	//    want to list only resources where Purpose is X, see the tag:key=value
	//    filter.
	//
	//    * tag-value - The value of a tag assigned to the resource. This filter
	//    is independent of the tag-key filter.
	//
	//    * vpc-id - The ID of the VPC for the subnet.
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// One or more subnet IDs.
	//
	// Default: Describes all your subnets.
	SubnetIds []*string `locationName:"SubnetId" locationNameList:"SubnetId" type:"list"`
}
type DescribeSubnetsOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more subnets.
	Subnets []*Subnet `locationName:"subnetSet" locationNameList:"item" type:"list"`

	RequestId *string `locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s DescribeSubnetsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeSubnetsOutput) GoString() string {
	return s.String()
}

// SetSubnets sets the Subnets field's value.
func (s *DescribeSubnetsOutput) SetSubnets(v []*Subnet) *DescribeSubnetsOutput {
	s.Subnets = v
	return s
}
func (s *DescribeSubnetsOutput) SetRequesterId(v *string) *DescribeSubnetsOutput {
	s.RequestId = v
	return s
}

// Contains the output of CreateNatGateway.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateNatGatewayResult
type CreateNatGatewayOutput struct {
	_ struct{} `type:"structure"`

	// Unique, case-sensitive identifier to ensure the idempotency of the request.
	// Only returned if a client token was provided in the request.
	ClientToken *string `locationName:"clientToken" type:"string"`

	// Information about the NAT gateway.
	NatGateway *NatGateway `locationName:"natGateway" type:"structure"`
}

// Describes a NAT gateway.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/NatGateway
type NatGateway struct {
	_ struct{} `type:"structure"`

	// The date and time the NAT gateway was created.
	CreateTime *time.Time `locationName:"createTime" type:"timestamp" timestampFormat:"iso8601"`

	// The date and time the NAT gateway was deleted, if applicable.
	DeleteTime *time.Time `locationName:"deleteTime" type:"timestamp" timestampFormat:"iso8601"`

	// If the NAT gateway could not be created, specifies the error code for the
	// failure. (InsufficientFreeAddressesInSubnet | Gateway.NotAttached | InvalidAllocationID.NotFound
	// | Resource.AlreadyAssociated | InternalError | InvalidSubnetID.NotFound)
	FailureCode *string `locationName:"failureCode" type:"string"`

	// If the NAT gateway could not be created, specifies the error message for
	// the failure, that corresponds to the error code.
	//
	//    * For InsufficientFreeAddressesInSubnet: "Subnet has insufficient free
	//    addresses to create this NAT gateway"
	//
	//    * For Gateway.NotAttached: "Network vpc-xxxxxxxx has no Internet gateway
	//    attached"
	//
	//    * For InvalidAllocationID.NotFound: "Elastic IP address eipalloc-xxxxxxxx
	//    could not be associated with this NAT gateway"
	//
	//    * For Resource.AlreadyAssociated: "Elastic IP address eipalloc-xxxxxxxx
	//    is already associated"
	//
	//    * For InternalError: "Network interface eni-xxxxxxxx, created and used
	//    internally by this NAT gateway is in an invalid state. Please try again."
	//
	//    * For InvalidSubnetID.NotFound: "The specified subnet subnet-xxxxxxxx
	//    does not exist or could not be found."
	FailureMessage *string `locationName:"failureMessage" type:"string"`

	// Information about the IP addresses and network interface associated with
	// the NAT gateway.
	NatGatewayAddresses []*NatGatewayAddress `locationName:"natGatewayAddressSet" locationNameList:"item" type:"list"`

	// The ID of the NAT gateway.
	NatGatewayId *string `locationName:"natGatewayId" type:"string"`

	// Reserved. If you need to sustain traffic greater than the documented limits
	// (http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/vpc-nat-gateway.html),
	// contact us through the Support Center (https://console.aws.amazon.com/support/home?).
	ProvisionedBandwidth *ProvisionedBandwidth `locationName:"provisionedBandwidth" type:"structure"`

	// The state of the NAT gateway.
	//
	//    * pending: The NAT gateway is being created and is not ready to process
	//    traffic.
	//
	//    * failed: The NAT gateway could not be created. Check the failureCode
	//    and failureMessage fields for the reason.
	//
	//    * available: The NAT gateway is able to process traffic. This status remains
	//    until you delete the NAT gateway, and does not indicate the health of
	//    the NAT gateway.
	//
	//    * deleting: The NAT gateway is in the process of being terminated and
	//    may still be processing traffic.
	//
	//    * deleted: The NAT gateway has been terminated and is no longer processing
	//    traffic.
	State *string `locationName:"state" type:"string" enum:"NatGatewayState"`

	// The ID of the subnet in which the NAT gateway is located.
	SubnetId *string `locationName:"subnetId" type:"string"`

	// The ID of the VPC in which the NAT gateway is located.
	VpcId *string `locationName:"vpcId" type:"string"`
}

// Describes the IP addresses and network interface associated with a NAT gateway.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/NatGatewayAddress
type NatGatewayAddress struct {
	_ struct{} `type:"structure"`

	// The allocation ID of the Elastic IP address that's associated with the NAT
	// gateway.
	AllocationId *string `locationName:"allocationId" type:"string"`

	// The ID of the network interface associated with the NAT gateway.
	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string"`

	// The private IP address associated with the Elastic IP address.
	PrivateIp *string `locationName:"privateIp" type:"string"`

	// The Elastic IP address associated with the NAT gateway.
	PublicIp *string `locationName:"publicIp" type:"string"`
}

type ProvisionedBandwidth struct {
	_ struct{} `type:"structure"`

	// Reserved. If you need to sustain traffic greater than the documented limits
	// (http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/vpc-nat-gateway.html),
	// contact us through the Support Center (https://console.aws.amazon.com/support/home?).
	ProvisionTime *time.Time `locationName:"provisionTime" type:"timestamp" timestampFormat:"iso8601"`

	// Reserved. If you need to sustain traffic greater than the documented limits
	// (http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/vpc-nat-gateway.html),
	// contact us through the Support Center (https://console.aws.amazon.com/support/home?).
	Provisioned *string `locationName:"provisioned" type:"string"`

	// Reserved. If you need to sustain traffic greater than the documented limits
	// (http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/vpc-nat-gateway.html),
	// contact us through the Support Center (https://console.aws.amazon.com/support/home?).
	RequestTime *time.Time `locationName:"requestTime" type:"timestamp" timestampFormat:"iso8601"`

	// Reserved. If you need to sustain traffic greater than the documented limits
	// (http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/vpc-nat-gateway.html),
	// contact us through the Support Center (https://console.aws.amazon.com/support/home?).
	Requested *string `locationName:"requested" type:"string"`

	// Reserved. If you need to sustain traffic greater than the documented limits
	// (http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/vpc-nat-gateway.html),
	// contact us through the Support Center (https://console.aws.amazon.com/support/home?).
	Status *string `locationName:"status" type:"string"`
}
type DescribeNatGatewaysInput struct {
	_ struct{} `type:"structure"`

	// One or more filters.
	//
	//    * nat-gateway-id - The ID of the NAT gateway.
	//
	//    * state - The state of the NAT gateway (pending | failed | available |
	//    deleting | deleted).
	//
	//    * subnet-id - The ID of the subnet in which the NAT gateway resides.
	//
	//    * vpc-id - The ID of the VPC in which the NAT gateway resides.
	Filter []*Filter `locationNameList:"Filter" type:"list"`

	// The maximum number of items to return for this request. The request returns
	// a token that you can specify in a subsequent call to get the next set of
	// results.
	//
	// Constraint: If the value specified is greater than 1000, we return only 1000
	// items.
	MaxResults *int64 `type:"integer"`

	// One or more NAT gateway IDs.
	NatGatewayIds []*string `locationName:"NatGatewayId" locationNameList:"item" type:"list"`

	// The token to retrieve the next page of results.
	NextToken *string `type:"string"`
}

type DescribeNatGatewaysOutput struct {
	_ struct{} `type:"structure"`

	// Information about the NAT gateways.
	NatGateways []*NatGateway `locationName:"natGatewaySet" locationNameList:"item" type:"list"`

	// The token to use to retrieve the next page of results. This value is null
	// when there are no more results to return.
	NextToken *string `locationName:"nextToken" type:"string"`

	RequestId *string `locationName:"requestId" type:"string"`
}

type DeleteNatGatewayInput struct {
	_ struct{} `type:"structure"`

	// The ID of the NAT gateway.
	//
	// NatGatewayId is a required field
	NatGatewayId *string `type:"string" required:"true"`
}

type DeleteNatGatewayOutput struct {
	_ struct{} `type:"structure"`

	// The ID of the NAT gateway.
	NatGatewayId *string `locationName:"natGatewayId" type:"string"`
}

// Contains the parameters for CreateVpc.
type CreateVpcInput struct {
	_ struct{} `type:"structure"`

	// Requests an Amazon-provided IPv6 CIDR block with a /56 prefix length for
	// the VPC. You cannot specify the range of IP addresses, or the size of the
	// CIDR block.
	AmazonProvidedIpv6CidrBlock *bool `locationName:"amazonProvidedIpv6CidrBlock" type:"boolean"`

	// The IPv4 network range for the VPC, in CIDR notation. For example, 10.0.0.0/16.
	//
	// CidrBlock is a required field
	CidrBlock *string `type:"string" required:"true"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The tenancy options for instances launched into the VPC. For default, instances
	// are launched with shared tenancy by default. You can launch instances with
	// any tenancy into a shared tenancy VPC. For dedicated, instances are launched
	// as dedicated tenancy instances by default. You can only launch instances
	// with a tenancy of dedicated or host into a dedicated tenancy VPC.
	//
	// Important: The host value cannot be used with this parameter. Use the default
	// or dedicated values only.
	//
	// Default: default
	InstanceTenancy *string `locationName:"instanceTenancy" type:"string" enum:"Tenancy"`
}

// Contains the output of CreateVpc.
type CreateVpcOutput struct {
	_ struct{} `type:"structure"`

	// Information about the VPC.
	Vpc *Vpc `locationName:"vpc" type:"structure"`
}

// Contains the parameters for DescribeVpcs.
type DescribeVpcsInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more filters.
	//
	//    * cidr - The primary IPv4 CIDR block of the VPC. The CIDR block you specify
	//    must exactly match the VPC's CIDR block for information to be returned
	//    for the VPC. Must contain the slash followed by one or two digits (for
	//    example, /28).
	//
	//    * cidr-block-association.cidr-block - An IPv4 CIDR block associated with
	//    the VPC.
	//
	//    * cidr-block-association.association-id - The association ID for an IPv4
	//    CIDR block associated with the VPC.
	//
	//    * cidr-block-association.state - The state of an IPv4 CIDR block associated
	//    with the VPC.
	//
	//    * dhcp-options-id - The ID of a set of DHCP options.
	//
	//    * ipv6-cidr-block-association.ipv6-cidr-block - An IPv6 CIDR block associated
	//    with the VPC.
	//
	//    * ipv6-cidr-block-association.association-id - The association ID for
	//    an IPv6 CIDR block associated with the VPC.
	//
	//    * ipv6-cidr-block-association.state - The state of an IPv6 CIDR block
	//    associated with the VPC.
	//
	//    * isDefault - Indicates whether the VPC is the default VPC.
	//
	//    * state - The state of the VPC (pending | available).
	//
	//    * tag:key=value - The key/value combination of a tag assigned to the resource.
	//    Specify the key of the tag in the filter name and the value of the tag
	//    in the filter value. For example, for the tag Purpose=X, specify tag:Purpose
	//    for the filter name and X for the filter value.
	//
	//    * tag-key - The key of a tag assigned to the resource. This filter is
	//    independent of the tag-value filter. For example, if you use both the
	//    filter "tag-key=Purpose" and the filter "tag-value=X", you get any resources
	//    assigned both the tag key Purpose (regardless of what the tag's value
	//    is), and the tag value X (regardless of what the tag's key is). If you
	//    want to list only resources where Purpose is X, see the tag:key=value
	//    filter.
	//
	//    * tag-value - The value of a tag assigned to the resource. This filter
	//    is independent of the tag-key filter.
	//
	//    * vpc-id - The ID of the VPC.
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// One or more VPC IDs.
	//
	// Default: Describes all your VPCs.
	VpcIds []*string `locationName:"VpcId" locationNameList:"VpcId" type:"list"`
}

// Contains the output of DescribeVpcs.
type DescribeVpcsOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more VPCs.
	Vpcs []*Vpc `locationName:"vpcSet" locationNameList:"item" type:"list"`

	RequestId *string `locationName:"requestId" type:"string"`
}

// Describes a VPC.
type Vpc struct {
	_ struct{} `type:"structure"`

	// The primary IPv4 CIDR block for the VPC.
	CidrBlock *string `locationName:"cidrBlock" type:"string"`

	// Information about the IPv4 CIDR blocks associated with the VPC.
	CidrBlockAssociationSet []*VpcCidrBlockAssociation `locationName:"cidrBlockAssociationSet" locationNameList:"item" type:"list"`

	// The ID of the set of DHCP options you've associated with the VPC (or default
	// if the default options are associated with the VPC).
	DhcpOptionsId *string `locationName:"dhcpOptionsId" type:"string"`

	// The allowed tenancy of instances launched into the VPC.
	InstanceTenancy *string `locationName:"instanceTenancy" type:"string" enum:"Tenancy"`

	// Information about the IPv6 CIDR blocks associated with the VPC.
	Ipv6CidrBlockAssociationSet []*VpcIpv6CidrBlockAssociation `locationName:"ipv6CidrBlockAssociationSet" locationNameList:"item" type:"list"`

	// Indicates whether the VPC is the default VPC.
	IsDefault *bool `locationName:"isDefault" type:"boolean"`

	// The current state of the VPC.
	State *string `locationName:"state" type:"string" enum:"VpcState"`

	// Any tags assigned to the VPC.
	Tags []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`

	// The ID of the VPC.
	VpcId *string `locationName:"vpcId" type:"string"`
}

// Describes an IPv4 CIDR block associated with a VPC.
type VpcCidrBlockAssociation struct {
	_ struct{} `type:"structure"`

	// The association ID for the IPv4 CIDR block.
	AssociationId *string `locationName:"associationId" type:"string"`

	// The IPv4 CIDR block.
	CidrBlock *string `locationName:"cidrBlock" type:"string"`

	// Information about the state of the CIDR block.
	CidrBlockState *VpcCidrBlockState `locationName:"cidrBlockState" type:"structure"`
}

// Describes the state of a CIDR block.
type VpcCidrBlockState struct {
	_ struct{} `type:"structure"`

	// The state of the CIDR block.
	State *string `locationName:"state" type:"string" enum:"VpcCidrBlockStateCode"`

	// A message about the status of the CIDR block, if applicable.
	StatusMessage *string `locationName:"statusMessage" type:"string"`
}

// Describes an IPv6 CIDR block associated with a VPC.
type VpcIpv6CidrBlockAssociation struct {
	_ struct{} `type:"structure"`

	// The association ID for the IPv6 CIDR block.
	AssociationId *string `locationName:"associationId" type:"string"`

	// The IPv6 CIDR block.
	Ipv6CidrBlock *string `locationName:"ipv6CidrBlock" type:"string"`

	// Information about the state of the CIDR block.
	Ipv6CidrBlockState *VpcCidrBlockState `locationName:"ipv6CidrBlockState" type:"structure"`
}

// Contains the parameters for DeleteVpc.
type DeleteVpcInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the VPC.
	//
	// VpcId is a required field
	VpcId *string `type:"string" required:"true"`
}

type DeleteVpcOutput struct {
	_ struct{} `type:"structure"`
}

//CreateInternetGatewayOutput Contains the output of CreateInternetGateway.
type CreateInternetGatewayOutput struct {
	_ struct{} `type:"structure"`

	// Information about the Internet gateway.
	InternetGateway *InternetGateway `locationName:"internetGateway" type:"structure"`
}

//InternetGateway Describes an Internet gateway.
type InternetGateway struct {
	_ struct{} `type:"structure"`

	// Any VPCs attached to the Internet gateway.
	Attachments []*InternetGatewayAttachment `locationName:"attachmentSet" locationNameList:"item" type:"list"`

	// The ID of the Internet gateway.
	InternetGatewayId *string `locationName:"internetGatewayId" type:"string"`

	// Any tags assigned to the Internet gateway.
	Tags []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`
}

//InternetGatewayAttachment Describes the attachment of a VPC to an Internet gateway or an egress-only
// Internet gateway.
type InternetGatewayAttachment struct {
	_ struct{} `type:"structure"`

	// The current state of the attachment. For an Internet gateway, the state is
	// available when attached to a VPC; otherwise, this value is not returned.
	State *string `locationName:"state" type:"string" enum:"AttachmentStatus"`

	// The ID of the VPC.
	VpcId *string `locationName:"vpcId" type:"string"`
}

type ModifyVpcAttributeInput struct {
	_ struct{} `type:"structure"`

	// Indicates whether the instances launched in the VPC get DNS hostnames. If
	// enabled, instances in the VPC get DNS hostnames; otherwise, they do not.
	//
	// You cannot modify the DNS resolution and DNS hostnames attributes in the
	// same request. Use separate requests for each attribute. You can only enable
	// DNS hostnames if you've enabled DNS support.
	EnableDnsHostnames *AttributeBooleanValue `type:"structure"`

	// Indicates whether the DNS resolution is supported for the VPC. If enabled,
	// queries to the Amazon provided DNS server at the 169.254.169.253 IP address,
	// or the reserved IP address at the base of the VPC network range "plus two"
	// will succeed. If disabled, the Amazon provided DNS service in the VPC that
	// resolves public DNS hostnames to IP addresses is not enabled.
	//
	// You cannot modify the DNS resolution and DNS hostnames attributes in the
	// same request. Use separate requests for each attribute.
	EnableDnsSupport *AttributeBooleanValue `type:"structure"`

	// The ID of the VPC.
	//
	// VpcId is a required field
	VpcId *string `locationName:"vpcId" type:"string" required:"true"`
}

// String returns the string representation
func (s ModifyVpcAttributeInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ModifyVpcAttributeInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *ModifyVpcAttributeInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "ModifyVpcAttributeInput"}
	if s.VpcId == nil {
		invalidParams.Add(request.NewErrParamRequired("VpcId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetEnableDnsHostnames sets the EnableDnsHostnames field's value.
func (s *ModifyVpcAttributeInput) SetEnableDnsHostnames(v *AttributeBooleanValue) *ModifyVpcAttributeInput {
	s.EnableDnsHostnames = v
	return s
}

// SetEnableDnsSupport sets the EnableDnsSupport field's value.
func (s *ModifyVpcAttributeInput) SetEnableDnsSupport(v *AttributeBooleanValue) *ModifyVpcAttributeInput {
	s.EnableDnsSupport = v
	return s
}

// SetVpcId sets the VpcId field's value.
func (s *ModifyVpcAttributeInput) SetVpcId(v string) *ModifyVpcAttributeInput {
	s.VpcId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/ModifyVpcAttributeOutput
type ModifyVpcAttributeOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s ModifyVpcAttributeOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ModifyVpcAttributeOutput) GoString() string {
	return s.String()
}

// Contains the parameters for DescribeInternetGateways.
type DescribeInternetGatewaysInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more filters.
	//
	//    * attachment.state - The current state of the attachment between the gateway
	//    and the VPC (available). Present only if a VPC is attached.
	//
	//    * attachment.vpc-id - The ID of an attached VPC.
	//
	//    * internet-gateway-id - The ID of the Internet gateway.
	//
	//    * tag:key=value - The key/value combination of a tag assigned to the resource.
	//    Specify the key of the tag in the filter name and the value of the tag
	//    in the filter value. For example, for the tag Purpose=X, specify tag:Purpose
	//    for the filter name and X for the filter value.
	//
	//    * tag-key - The key of a tag assigned to the resource. This filter is
	//    independent of the tag-value filter. For example, if you use both the
	//    filter "tag-key=Purpose" and the filter "tag-value=X", you get any resources
	//    assigned both the tag key Purpose (regardless of what the tag's value
	//    is), and the tag value X (regardless of what the tag's key is). If you
	//    want to list only resources where Purpose is X, see the tag:key=value
	//    filter.
	//
	//    * tag-value - The value of a tag assigned to the resource. This filter
	//    is independent of the tag-key filter.
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// One or more Internet gateway IDs.
	//
	// Default: Describes all your Internet gateways.
	InternetGatewayIds []*string `locationName:"internetGatewayId" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s DescribeInternetGatewaysInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeInternetGatewaysInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeInternetGatewaysInput) SetDryRun(v bool) *DescribeInternetGatewaysInput {
	s.DryRun = &v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribeInternetGatewaysInput) SetFilters(v []*Filter) *DescribeInternetGatewaysInput {
	s.Filters = v
	return s
}

// SetInternetGatewayIds sets the InternetGatewayIds field's value.
func (s *DescribeInternetGatewaysInput) SetInternetGatewayIds(v []*string) *DescribeInternetGatewaysInput {
	s.InternetGatewayIds = v
	return s
}

//DescribeInternetGatewaysOutput Contains the output of DescribeInternetGateways.
type DescribeInternetGatewaysOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more Internet gateways.
	InternetGateways []*InternetGateway `locationName:"internetGatewaySet" locationNameList:"item" type:"list"`
	RequestId        *string            `locationName:"requestId" type:"string"`
}

type DescribeVpcAttributeInput struct {
	_ struct{} `type:"structure"`

	// The VPC attribute.
	//
	// Attribute is a required field
	Attribute *string `type:"string" required:"true" enum:"VpcAttributeName"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the VPC.
	//
	// VpcId is a required field
	VpcId *string `type:"string" required:"true"`

	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`
}

func (s *DescribeVpcAttributeInput) SetFilters(v []*Filter) *DescribeVpcAttributeInput {
	s.Filters = v
	return s
}

// String returns the string representation
func (s DescribeVpcAttributeInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeVpcAttributeInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DescribeVpcAttributeInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DescribeVpcAttributeInput"}
	if s.Attribute == nil {
		invalidParams.Add(request.NewErrParamRequired("Attribute"))
	}
	if s.VpcId == nil {
		invalidParams.Add(request.NewErrParamRequired("VpcId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAttribute sets the Attribute field's value.
func (s *DescribeVpcAttributeInput) SetAttribute(v string) *DescribeVpcAttributeInput {
	s.Attribute = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeVpcAttributeInput) SetDryRun(v bool) *DescribeVpcAttributeInput {
	s.DryRun = &v
	return s
}

// SetVpcId sets the VpcId field's value.
func (s *DescribeVpcAttributeInput) SetVpcId(v string) *DescribeVpcAttributeInput {
	s.VpcId = &v
	return s
}

type DescribeVpcAttributeOutput struct {
	_ struct{} `type:"structure"`

	// Indicates whether the instances launched in the VPC get DNS hostnames. If
	// this attribute is true, instances in the VPC get DNS hostnames; otherwise,
	// they do not.
	EnableDnsHostnames *AttributeBooleanValue `locationName:"enableDnsHostnames" type:"structure"`

	// Indicates whether DNS resolution is enabled for the VPC. If this attribute
	// is true, the Amazon DNS server resolves DNS hostnames for your instances
	// to their corresponding IP addresses; otherwise, it does not.
	EnableDnsSupport *AttributeBooleanValue `locationName:"enableDnsSupport" type:"structure"`

	// The ID of the VPC.
	VpcId *string `locationName:"vpcId" type:"string"`

	RequestId *string `locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s DescribeVpcAttributeOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeVpcAttributeOutput) GoString() string {
	return s.String()
}

// SetEnableDnsHostnames sets the EnableDnsHostnames field's value.
func (s *DescribeVpcAttributeOutput) SetEnableDnsHostnames(v *AttributeBooleanValue) *DescribeVpcAttributeOutput {
	s.EnableDnsHostnames = v
	return s
}

// SetEnableDnsSupport sets the EnableDnsSupport field's value.
func (s *DescribeVpcAttributeOutput) SetRequesterId(v *string) *DescribeVpcAttributeOutput {
	s.RequestId = v
	return s
}

func (s *DescribeVpcAttributeOutput) SetEnableDnsSupport(v *AttributeBooleanValue) *DescribeVpcAttributeOutput {
	s.EnableDnsSupport = v
	return s
}

// SetVpcId sets the VpcId field's value.
func (s *DescribeVpcAttributeOutput) SetVpcId(v string) *DescribeVpcAttributeOutput {
	s.VpcId = &v
	return s
}

type AttachInternetGatewayInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the Internet gateway.
	//
	// InternetGatewayId is a required field
	InternetGatewayId *string `locationName:"internetGatewayId" type:"string" required:"true"`

	// The ID of the VPC.
	//
	// VpcId is a required field
	VpcId *string `locationName:"vpcId" type:"string" required:"true"`
}

// String returns the string representation
func (s AttachInternetGatewayInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AttachInternetGatewayInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *AttachInternetGatewayInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "AttachInternetGatewayInput"}
	if s.InternetGatewayId == nil {
		invalidParams.Add(request.NewErrParamRequired("InternetGatewayId"))
	}
	if s.VpcId == nil {
		invalidParams.Add(request.NewErrParamRequired("VpcId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDryRun sets the DryRun field's value.
func (s *AttachInternetGatewayInput) SetDryRun(v bool) *AttachInternetGatewayInput {
	s.DryRun = &v
	return s
}

// SetInternetGatewayId sets the InternetGatewayId field's value.
func (s *AttachInternetGatewayInput) SetInternetGatewayId(v string) *AttachInternetGatewayInput {
	s.InternetGatewayId = &v
	return s
}

// SetVpcId sets the VpcId field's value.
func (s *AttachInternetGatewayInput) SetVpcId(v string) *AttachInternetGatewayInput {
	s.VpcId = &v
	return s
}

type DeleteInternetGatewayInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the Internet gateway.
	//
	// InternetGatewayId is a required field
	InternetGatewayId *string `locationName:"internetGatewayId" type:"string" required:"true"`
}

// String returns the string representation
func (s DeleteInternetGatewayInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteInternetGatewayInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteInternetGatewayInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeleteInternetGatewayInput"}
	if s.InternetGatewayId == nil {
		invalidParams.Add(request.NewErrParamRequired("InternetGatewayId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDryRun sets the DryRun field's value.
func (s *DeleteInternetGatewayInput) SetDryRun(v bool) *DeleteInternetGatewayInput {
	s.DryRun = &v
	return s
}

// SetInternetGatewayId sets the InternetGatewayId field's value.
func (s *DeleteInternetGatewayInput) SetInternetGatewayId(v string) *DeleteInternetGatewayInput {
	s.InternetGatewayId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteInternetGatewayOutput
type DeleteInternetGatewayOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeleteInternetGatewayOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteInternetGatewayOutput) GoString() string {
	return s.String()
}

// Contains the parameters for CreateNatGateway.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateNatGatewayRequest
type CreateNatGatewayInput struct {
	_ struct{} `type:"structure"`

	// The allocation ID of an Elastic IP address to associate with the NAT gateway.
	// If the Elastic IP address is associated with another resource, you must first
	// disassociate it.
	//
	// AllocationId is a required field
	AllocationId *string `type:"string" required:"true"`

	// Unique, case-sensitive identifier you provide to ensure the idempotency of
	// the request. For more information, see How to Ensure Idempotency (http://docs.aws.amazon.com/AWSEC2/latest/APIReference/Run_Instance_Idempotency.html).
	//
	// Constraint: Maximum 64 ASCII characters.
	ClientToken *string `type:"string"`

	// The subnet in which to create the NAT gateway.
	//
	// SubnetId is a required field
	SubnetId *string `type:"string" required:"true"`
}

type AttachInternetGatewayOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s AttachInternetGatewayOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AttachInternetGatewayOutput) GoString() string {
	return s.String()
}

type DetachInternetGatewayOutput struct {
	_ struct{} `type:"structure"`
}

type DetachInternetGatewayInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the Internet gateway.
	//
	// InternetGatewayId is a required field
	InternetGatewayId *string `locationName:"internetGatewayId" type:"string" required:"true"`

	// The ID of the VPC.
	//
	// VpcId is a required field
	VpcId *string `locationName:"vpcId" type:"string" required:"true"`
}

type CreateAccessKeyInput struct {
	_ struct{} `type:"structure"`

	AccessKeyId     *string `type:"string" required:"false"`
	SecretAccessKey *string `type:"string" required:"false"`
	Tags            []*Tag  `locationName:"tagSet" locationNameList:"item" type:"list"`
}
type CreateAccessKeyOutput struct {
	_                struct{} `type:"structure"`
	AccessKey        *string  `locationName:"accessKey" type:"structure"`
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
	Tags            []*Tag   `locationName:"tagSet" locationNameList:"item" type:"list"`
}
type DescribeAccessKeyOutput struct {
	_                struct{} `type:"structure"`
	AccessKey        *string  `locationName:"accessKey" type:"structure"`
	ResponseMetadata *string  `locationName:"responseMetaData" type:"structure"`
}

// Contains the parameters for DeleteDhcpOptions.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteDhcpOptionsRequest
type DeleteDhcpOptionsInput struct {
	_ struct{} `type:"structure"`

	// The ID of the DHCP options set.
	//
	// DhcpOptionsId is a required field
	DhcpOptionsId *string `type:"string" required:"true"`

	// The Internet-routable IP address for the customer gateway's outside interface.
	// The address must be static.
	//
	// PublicIp is a required field
	PublicIp *string `locationName:"IpAddress" type:"string" required:"true"`

	// The type of VPN connection that this customer gateway supports (ipsec.1).
	//
	// Type is a required field
	Type *string `type:"string" required:"true" enum:"GatewayType"`
}

// String returns the string representation
func (s DeleteDhcpOptionsInput) String() string {
	return awsutil.Prettify(s)
}

// String returns the string representation
func (s CreateCustomerGatewayInput) String() string {
	return awsutil.Prettify(s)
}

type CreateCustomerGatewayInput struct {
	_ struct{} `type:"structure"`

	// For devices that support BGP, the customer gateway's BGP ASN.
	//
	// Default: 65000
	//
	// BgpAsn is a required field
	BgpAsn *int64 `type:"integer" required:"true"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// PublicIp is a required field
	PublicIp *string `locationName:"IpAddress" type:"string" required:"true"`

	// The type of VPN connection that this customer gateway supports (ipsec.1).
	//
	// Type is a required field
	Type *string `type:"string" required:"true" enum:"GatewayType"`
}

// GoString returns the string representation
func (s DeleteDhcpOptionsInput) GoString() string {
	return s.String()
}
func (s CreateCustomerGatewayInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteDhcpOptionsInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeleteDhcpOptionsInput"}
	if s.DhcpOptionsId == nil {
		invalidParams.Add(request.NewErrParamRequired("DhcpOptionsId"))
	}
	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}
func (s *CreateCustomerGatewayInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "CreateCustomerGatewayInput"}
	if s.BgpAsn == nil {
		invalidParams.Add(request.NewErrParamRequired("BgpAsn"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDhcpOptionsId sets the DhcpOptionsId field's value.
func (s *DeleteDhcpOptionsInput) SetDhcpOptionsId(v string) *DeleteDhcpOptionsInput {
	s.DhcpOptionsId = &v
	return s
}

// SetBgpAsn sets the BgpAsn field's value.
func (s *CreateCustomerGatewayInput) SetBgpAsn(v int64) *CreateCustomerGatewayInput {
	s.BgpAsn = &v
	return s
}

func (s *CreateCustomerGatewayInput) SetDryRun(v bool) *CreateCustomerGatewayInput {
	s.DryRun = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteDhcpOptionsOutput
type DeleteDhcpOptionsOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeleteDhcpOptionsOutput) String() string {
	return awsutil.Prettify(s)
}

// SetPublicIp sets the PublicIp field's value.
func (s *CreateCustomerGatewayInput) SetPublicIp(v string) *CreateCustomerGatewayInput {
	s.PublicIp = &v
	return s
}

// SetType sets the Type field's value.
func (s *CreateCustomerGatewayInput) SetType(v string) *CreateCustomerGatewayInput {
	s.Type = &v
	return s
}

// Contains the output of CreateCustomerGateway.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateCustomerGatewayResult
type CreateCustomerGatewayOutput struct {
	_ struct{} `type:"structure"`

	// Information about the customer gateway.
	CustomerGateway *CustomerGateway `locationName:"customerGateway" type:"structure"`
}

// String returns the string representation
func (s CreateCustomerGatewayOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteDhcpOptionsOutput) GoString() string {
	return s.String()
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/NewDhcpConfiguration
type NewDhcpConfiguration struct {
	_ struct{} `type:"structure"`

	Key *string `locationName:"key" type:"string"`

	Values []*string `locationName:"Value" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s NewDhcpConfiguration) String() string {
	return awsutil.Prettify(s)
}
func (s CreateCustomerGatewayOutput) GoString() string {
	return s.String()
}

// SetCustomerGateway sets the CustomerGateway field's value.
func (s *CreateCustomerGatewayOutput) SetCustomerGateway(v *CustomerGateway) *CreateCustomerGatewayOutput {
	s.CustomerGateway = v
	return s
}

type CustomerGateway struct {
	_ struct{} `type:"structure"`

	// The customer gateway's Border Gateway Protocol (BGP) Autonomous System Number
	// (ASN).
	BgpAsn *string `locationName:"bgpAsn" type:"string"`

	// The ID of the customer gateway.
	CustomerGatewayId *string `locationName:"customerGatewayId" type:"string"`

	// The Internet-routable IP address of the customer gateway's outside interface.
	IpAddress *string `locationName:"ipAddress" type:"string"`

	// The current state of the customer gateway (pending | available | deleting
	// | deleted).
	State *string `locationName:"state" type:"string"`

	// Any tags assigned to the customer gateway.
	Tags []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`

	// The type of VPN connection the customer gateway supports (ipsec.1).
	Type *string `locationName:"type" type:"string"`
}

// String returns the string representation
func (s CustomerGateway) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s NewDhcpConfiguration) GoString() string {
	return s.String()
}

// SetKey sets the Key field's value.
func (s *NewDhcpConfiguration) SetKey(v string) *NewDhcpConfiguration {
	s.Key = &v
	return s
}

// SetValues sets the Values field's value.
func (s *NewDhcpConfiguration) SetValues(v []*string) *NewDhcpConfiguration {
	s.Values = v
	return s
}

// Contains the parameters for CreateDhcpOptions.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateDhcpOptionsRequest
type CreateDhcpOptionsInput struct {
	_ struct{} `type:"structure"`

	// A DHCP configuration option.
	//
	// DhcpConfigurations is a required field
	DhcpConfigurations []*NewDhcpConfiguration `locationName:"dhcpConfiguration" locationNameList:"item" type:"list" required:"true"`
	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`
}

func (s CustomerGateway) GoString() string {
	return s.String()
}

// SetBgpAsn sets the BgpAsn field's value.
func (s *CustomerGateway) SetBgpAsn(v string) *CustomerGateway {
	s.BgpAsn = &v
	return s
}

// SetCustomerGatewayId sets the CustomerGatewayId field's value.
func (s *CustomerGateway) SetCustomerGatewayId(v string) *CustomerGateway {
	s.CustomerGatewayId = &v
	return s
}

// SetIpAddress sets the IpAddress field's value.
func (s *CustomerGateway) SetIpAddress(v string) *CustomerGateway {
	s.IpAddress = &v
	return s
}

// SetState sets the State field's value.
func (s *CustomerGateway) SetState(v string) *CustomerGateway {
	s.State = &v
	return s
}

// SetTags sets the Tags field's value.
func (s *CustomerGateway) SetTags(v []*Tag) *CustomerGateway {
	s.Tags = v
	return s
}

// SetType sets the Type field's value.
func (s *CustomerGateway) SetType(v string) *CustomerGateway {
	s.Type = &v
	return s
}

// Contains the parameters for DeleteCustomerGateway.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteCustomerGatewayRequest
type DeleteCustomerGatewayInput struct {
	_ struct{} `type:"structure"`

	// The ID of the customer gateway.
	//
	// CustomerGatewayId is a required field
	CustomerGatewayId *string `type:"string" required:"true"`
	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`
}

// String returns the string representation
func (s CreateDhcpOptionsInput) String() string {
	return awsutil.Prettify(s)
}
func (s DeleteCustomerGatewayInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateDhcpOptionsInput) GoString() string {
	return s.String()
}
func (s DeleteCustomerGatewayInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CreateDhcpOptionsInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "CreateDhcpOptionsInput"}
	if s.DhcpConfigurations == nil {
		invalidParams.Add(request.NewErrParamRequired("DhcpConfigurations"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}
func (s *DeleteCustomerGatewayInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeleteCustomerGatewayInput"}
	if s.CustomerGatewayId == nil {
		invalidParams.Add(request.NewErrParamRequired("CustomerGatewayId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDhcpConfigurations sets the DhcpConfigurations field's value.
func (s *CreateDhcpOptionsInput) SetDhcpConfigurations(v []*NewDhcpConfiguration) *CreateDhcpOptionsInput {
	s.DhcpConfigurations = v
	return s
}

// SetCustomerGatewayId sets the CustomerGatewayId field's value.
func (s *DeleteCustomerGatewayInput) SetCustomerGatewayId(v string) *DeleteCustomerGatewayInput {
	s.CustomerGatewayId = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *CreateDhcpOptionsInput) SetDryRun(v bool) *CreateDhcpOptionsInput {
	s.DryRun = &v
	return s
}
func (s *DeleteCustomerGatewayInput) SetDryRun(v bool) *DeleteCustomerGatewayInput {
	s.DryRun = &v
	return s
}

// Contains the output of CreateDhcpOptions.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateDhcpOptionsResult
type CreateDhcpOptionsOutput struct {
	_ struct{} `type:"structure"`

	// A set of DHCP options.
	DhcpOptions *DhcpOptions `locationName:"dhcpOptions" type:"structure"`
}

// String returns the string representation
func (s CreateDhcpOptionsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateDhcpOptionsOutput) GoString() string {
	return s.String()
}

// SetDhcpOptions sets the DhcpOptions field's value.
func (s *CreateDhcpOptionsOutput) SetDhcpOptions(v *DhcpOptions) *CreateDhcpOptionsOutput {
	s.DhcpOptions = v
	return s
}

// Describes a set of DHCP options.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DhcpOptions
type DhcpOptions struct {
	_ struct{} `type:"structure"`

	// One or more DHCP options in the set.
	DhcpConfigurations []*DhcpConfiguration `locationName:"dhcpConfigurationSet" locationNameList:"item" type:"list"`

	// The ID of the set of DHCP options.
	DhcpOptionsId *string `locationName:"dhcpOptionsId" type:"string"`

	// Any tags assigned to the DHCP options set.
	Tags []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`
}

// Describes a DHCP configuration option.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DhcpConfiguration
type DhcpConfiguration struct {
	_ struct{} `type:"structure"`

	// The name of a DHCP option.
	Key *string `locationName:"key" type:"string"`

	// One or more values for the DHCP option.
	Values []*AttributeValue `locationName:"valueSet" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s DhcpConfiguration) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DhcpConfiguration) GoString() string {
	return s.String()
}

// SetKey sets the Key field's value.
func (s *DhcpConfiguration) SetKey(v string) *DhcpConfiguration {
	s.Key = &v
	return s
}

// SetValues sets the Values field's value.
func (s *DhcpConfiguration) SetValues(v []*AttributeValue) *DhcpConfiguration {
	s.Values = v
	return s
}

// String returns the string representation
func (s DhcpOptions) String() string {
	return awsutil.Prettify(s)
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteCustomerGatewayOutput
type DeleteCustomerGatewayOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeleteCustomerGatewayOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DhcpOptions) GoString() string {
	return s.String()
}

// SetDhcpConfigurations sets the DhcpConfigurations field's value.
func (s *DhcpOptions) SetDhcpConfigurations(v []*DhcpConfiguration) *DhcpOptions {
	s.DhcpConfigurations = v
	return s
}

// SetDhcpOptionsId sets the DhcpOptionsId field's value.
func (s *DhcpOptions) SetDhcpOptionsId(v string) *DhcpOptions {
	s.DhcpOptionsId = &v
	return s
}

// SetTags sets the Tags field's value.
func (s *DhcpOptions) SetTags(v []*Tag) *DhcpOptions {
	s.Tags = v
	return s
}

// Contains the parameters for DescribeDhcpOptions.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeDhcpOptionsRequest
type DescribeDhcpOptionsInput struct {
	_ struct{} `type:"structure"`

	// The IDs of one or more DHCP options sets.
	//
	// Default: Describes all your DHCP options sets.
	DhcpOptionsIds []*string `locationName:"DhcpOptionsId" locationNameList:"DhcpOptionsId" type:"list"`
	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more filters.
	//
	//    * dhcp-options-id - The ID of a set of DHCP options.
	//
	//    * key - The key for one of the options (for example, domain-name).
	//
	//    * value - The value for one of the options.
	//    * bgp-asn - The customer gateway's Border Gateway Protocol (BGP) Autonomous
	//    System Number (ASN).
	//
	//    * customer-gateway-id - The ID of the customer gateway.
	//
	//    * ip-address - The IP address of the customer gateway's Internet-routable
	//    external interface.
	//
	//    * state - The state of the customer gateway (pending | available | deleting
	//    | deleted).
	//
	//    * type - The type of customer gateway. Currently, the only supported type
	//    is ipsec.1.
	//
	//    * tag:key=value - The key/value combination of a tag assigned to the resource.
	//    Specify the key of the tag in the filter name and the value of the tag
	//    in the filter value. For example, for the tag Purpose=X, specify tag:Purpose
	//    for the filter name and X for the filter value.
	//
	//    * tag-key - The key of a tag assigned to the resource. This filter is
	//    independent of the tag-value filter. For example, if you use both the
	//    filter "tag-key=Purpose" and the filter "tag-value=X", you get any resources
	//    assigned both the tag key Purpose (regardless of what the tag's value
	//    is), and the tag value X (regardless of what the tag's key is). If you
	//    want to list only resources where Purpose is X, see the tag:key=value
	//    filter.
	//
	//    * tag-value - The value of a tag assigned to the resource. This filter
	//    is independent of the tag-key filter.
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`
}

func (s DeleteCustomerGatewayOutput) GoString() string {
	return s.String()
}

type DescribeCustomerGatewaysInput struct {
	_ struct{} `type:"structure"`

	// One or more customer gateway IDs.
	//
	// Default: Describes all your customer gateways.
	CustomerGatewayIds []*string `locationName:"CustomerGatewayId" locationNameList:"CustomerGatewayId" type:"list"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more filters.
	//
	//    * dhcp-options-id - The ID of a set of DHCP options.
	//
	//    * key - The key for one of the options (for example, domain-name).
	//
	//    * value - The value for one of the options.
	//    * bgp-asn - The customer gateway's Border Gateway Protocol (BGP) Autonomous
	//    System Number (ASN).
	//
	//    * customer-gateway-id - The ID of the customer gateway.
	//
	//    * ip-address - The IP address of the customer gateway's Internet-routable
	//    external interface.
	//
	//    * state - The state of the customer gateway (pending | available | deleting
	//    | deleted).
	//
	//    * type - The type of customer gateway. Currently, the only supported type
	//    is ipsec.1.
	//
	//    * tag:key=value - The key/value combination of a tag assigned to the resource.
	//    Specify the key of the tag in the filter name and the value of the tag
	//    in the filter value. For example, for the tag Purpose=X, specify tag:Purpose
	//    for the filter name and X for the filter value.
	//
	//    * tag-key - The key of a tag assigned to the resource. This filter is
	//    independent of the tag-value filter. For example, if you use both the
	//    filter "tag-key=Purpose" and the filter "tag-value=X", you get any resources
	//    assigned both the tag key Purpose (regardless of what the tag's value
	//    is), and the tag value X (regardless of what the tag's key is). If you
	//    want to list only resources where Purpose is X, see the tag:key=value
	//    filter.
	//
	//    * tag-value - The value of a tag assigned to the resource. This filter
	//    is independent of the tag-key filter.
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`
}

// String returns the string representation
func (s DescribeDhcpOptionsInput) String() string {
	return awsutil.Prettify(s)
}
func (s DescribeCustomerGatewaysInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeDhcpOptionsInput) GoString() string {
	return s.String()
}

// SetDhcpOptionsIds sets the DhcpOptionsIds field's value.
func (s *DescribeDhcpOptionsInput) SetDhcpOptionsIds(v []*string) *DescribeDhcpOptionsInput {
	s.DhcpOptionsIds = v
	return s
}
func (s DescribeCustomerGatewaysInput) GoString() string {
	return s.String()
}

// SetCustomerGatewayIds sets the CustomerGatewayIds field's value.
func (s *DescribeCustomerGatewaysInput) SetCustomerGatewayIds(v []*string) *DescribeCustomerGatewaysInput {
	s.CustomerGatewayIds = v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeDhcpOptionsInput) SetDryRun(v bool) *DescribeDhcpOptionsInput {
	s.DryRun = &v
	return s
}
func (s *DescribeCustomerGatewaysInput) SetDryRun(v bool) *DescribeCustomerGatewaysInput {
	s.DryRun = &v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribeDhcpOptionsInput) SetFilters(v []*Filter) *DescribeDhcpOptionsInput {
	s.Filters = v
	return s
}
func (s *DescribeCustomerGatewaysInput) SetFilters(v []*Filter) *DescribeCustomerGatewaysInput {
	s.Filters = v
	return s
}

// Contains the output of DescribeDhcpOptions.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeDhcpOptionsResult
type DescribeDhcpOptionsOutput struct {
	_         struct{} `type:"structure"`
	RequestId *string  `locationName:"requestId" type:"string"`

	// Information about one or more DHCP options sets.
	DhcpOptions []*DhcpOptions `locationName:"dhcpOptionsSet" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s DescribeDhcpOptionsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeDhcpOptionsOutput) GoString() string {
	return s.String()
}

// SetDhcpOptions sets the DhcpOptions field's value.
func (s *DescribeDhcpOptionsOutput) SetDhcpOptions(v []*DhcpOptions) *DescribeDhcpOptionsOutput {
	s.DhcpOptions = v
	return s
}

// Contains the parameters for AssociateDhcpOptions.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/AssociateDhcpOptionsRequest
type AssociateDhcpOptionsInput struct {
	_ struct{} `type:"structure"`

	// The ID of the DHCP options set, or default to associate no DHCP options with
	// the VPC.
	//
	// DhcpOptionsId is a required field
	DhcpOptionsId *string `type:"string" required:"true"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the VPC.
	//
	// VpcId is a required field
	VpcId *string `type:"string" required:"true"`
}

// String returns the string representation
func (s AssociateDhcpOptionsInput) String() string {
	return awsutil.Prettify(s)
}

// Contains the output of DescribeCustomerGateways.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeCustomerGatewaysResult
type DescribeCustomerGatewaysOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more customer gateways.
	CustomerGateways []*CustomerGateway `locationName:"customerGatewaySet" locationNameList:"item" type:"list"`

	RequestId *string `locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s DescribeCustomerGatewaysOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AssociateDhcpOptionsInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *AssociateDhcpOptionsInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "AssociateDhcpOptionsInput"}
	if s.DhcpOptionsId == nil {
		invalidParams.Add(request.NewErrParamRequired("DhcpOptionsId"))
	}
	if s.VpcId == nil {
		invalidParams.Add(request.NewErrParamRequired("VpcId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDhcpOptionsId sets the DhcpOptionsId field's value.
func (s *AssociateDhcpOptionsInput) SetDhcpOptionsId(v string) *AssociateDhcpOptionsInput {
	s.DhcpOptionsId = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *AssociateDhcpOptionsInput) SetDryRun(v bool) *AssociateDhcpOptionsInput {
	s.DryRun = &v
	return s
}

// SetVpcId sets the VpcId field's value.
func (s *AssociateDhcpOptionsInput) SetVpcId(v string) *AssociateDhcpOptionsInput {
	s.VpcId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/AssociateDhcpOptionsOutput
type AssociateDhcpOptionsOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s AssociateDhcpOptionsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AssociateDhcpOptionsOutput) GoString() string {
	return s.String()
}
func (s DescribeCustomerGatewaysOutput) GoString() string {
	return s.String()
}

// SetCustomerGateways sets the CustomerGateways field's value.
func (s *DescribeCustomerGatewaysOutput) SetCustomerGateways(v []*CustomerGateway) *DescribeCustomerGatewaysOutput {
	s.CustomerGateways = v
	return s
}
func (s *DescribeCustomerGatewaysOutput) SetRequesterId(v *string) *DescribeCustomerGatewaysOutput {
	s.RequestId = v
	return s
}

type CreateRouteInput struct {
	_ struct{} `type:"structure"`

	// The IPv4 CIDR address block used for the destination match. Routing decisions
	// are based on the most specific match.
	DestinationCidrBlock *string `locationName:"destinationCidrBlock" type:"string"`

	// The IPv6 CIDR block used for the destination match. Routing decisions are
	// based on the most specific match.
	DestinationIpv6CidrBlock *string `locationName:"destinationIpv6CidrBlock" type:"string"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// [IPv6 traffic only] The ID of an egress-only Internet gateway.
	EgressOnlyInternetGatewayId *string `locationName:"egressOnlyInternetGatewayId" type:"string"`

	// The ID of an Internet gateway or virtual private gateway attached to your
	// VPC.
	GatewayId *string `locationName:"gatewayId" type:"string"`

	// The ID of a NAT instance in your VPC. The operation fails if you specify
	// an instance ID unless exactly one network interface is attached.
	InstanceId *string `locationName:"instanceId" type:"string"`

	// [IPv4 traffic only] The ID of a NAT gateway.
	NatGatewayId *string `locationName:"natGatewayId" type:"string"`

	// The ID of a network interface.
	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string"`

	// The ID of the route table for the route.
	//
	// RouteTableId is a required field
	RouteTableId *string `locationName:"routeTableId" type:"string" required:"true"`

	// The ID of a VPC peering connection.
	VpcPeeringConnectionId *string `locationName:"vpcPeeringConnectionId" type:"string"`
}

// String returns the string representation
func (s CreateRouteInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateRouteInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CreateRouteInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "CreateRouteInput"}
	if s.RouteTableId == nil {
		invalidParams.Add(request.NewErrParamRequired("RouteTableId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDestinationCidrBlock sets the DestinationCidrBlock field's value.
func (s *CreateRouteInput) SetDestinationCidrBlock(v string) *CreateRouteInput {
	s.DestinationCidrBlock = &v
	return s
}

// SetDestinationIpv6CidrBlock sets the DestinationIpv6CidrBlock field's value.
func (s *CreateRouteInput) SetDestinationIpv6CidrBlock(v string) *CreateRouteInput {
	s.DestinationIpv6CidrBlock = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *CreateRouteInput) SetDryRun(v bool) *CreateRouteInput {
	s.DryRun = &v
	return s
}

// SetEgressOnlyInternetGatewayId sets the EgressOnlyInternetGatewayId field's value.
func (s *CreateRouteInput) SetEgressOnlyInternetGatewayId(v string) *CreateRouteInput {
	s.EgressOnlyInternetGatewayId = &v
	return s
}

// SetGatewayId sets the GatewayId field's value.
func (s *CreateRouteInput) SetGatewayId(v string) *CreateRouteInput {
	s.GatewayId = &v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *CreateRouteInput) SetInstanceId(v string) *CreateRouteInput {
	s.InstanceId = &v
	return s
}

// SetNatGatewayId sets the NatGatewayId field's value.
func (s *CreateRouteInput) SetNatGatewayId(v string) *CreateRouteInput {
	s.NatGatewayId = &v
	return s
}

// SetNetworkInterfaceId sets the NetworkInterfaceId field's value.
func (s *CreateRouteInput) SetNetworkInterfaceId(v string) *CreateRouteInput {
	s.NetworkInterfaceId = &v
	return s
}

// SetRouteTableId sets the RouteTableId field's value.
func (s *CreateRouteInput) SetRouteTableId(v string) *CreateRouteInput {
	s.RouteTableId = &v
	return s
}

// SetVpcPeeringConnectionId sets the VpcPeeringConnectionId field's value.
func (s *CreateRouteInput) SetVpcPeeringConnectionId(v string) *CreateRouteInput {
	s.VpcPeeringConnectionId = &v
	return s
}

// Contains the output of CreateRoute.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateRouteResult
type CreateRouteOutput struct {
	_ struct{} `type:"structure"`

	// Returns true if the request succeeds; otherwise, it returns an error.
	Return *bool `locationName:"return" type:"boolean"`
}

// String returns the string representation
func (s CreateRouteOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateRouteOutput) GoString() string {
	return s.String()
}

// SetReturn sets the Return field's value.
func (s *CreateRouteOutput) SetReturn(v bool) *CreateRouteOutput {
	s.Return = &v
	return s
}

type Route struct {
	_ struct{} `type:"structure"`

	// The IPv4 CIDR block used for the destination match.
	DestinationCidrBlock *string `locationName:"destinationCidrBlock" type:"string"`

	// The IPv6 CIDR block used for the destination match.
	DestinationIpv6CidrBlock *string `locationName:"destinationIpv6CidrBlock" type:"string"`

	// The prefix of the AWS service.
	DestinationPrefixListId *string `locationName:"destinationPrefixListId" type:"string"`

	// The ID of the egress-only Internet gateway.
	EgressOnlyInternetGatewayId *string `locationName:"egressOnlyInternetGatewayId" type:"string"`

	// The ID of a gateway attached to your VPC.
	GatewayId *string `locationName:"gatewayId" type:"string"`

	// The ID of a NAT instance in your VPC.
	InstanceId *string `locationName:"instanceId" type:"string"`

	// The AWS account ID of the owner of the instance.
	InstanceOwnerId *string `locationName:"instanceOwnerId" type:"string"`

	// The ID of a NAT gateway.
	NatGatewayId *string `locationName:"natGatewayId" type:"string"`

	// The ID of the network interface.
	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string"`

	// Describes how the route was created.
	//
	//    * CreateRouteTable - The route was automatically created when the route
	//    table was created.
	//
	//    * CreateRoute - The route was manually added to the route table.
	//
	//    * EnableVgwRoutePropagation - The route was propagated by route propagation.
	Origin *string `locationName:"origin" type:"string" enum:"RouteOrigin"`

	// The state of the route. The blackhole state indicates that the route's target
	// isn't available (for example, the specified gateway isn't attached to the
	// VPC, or the specified NAT instance has been terminated).
	State *string `locationName:"state" type:"string" enum:"RouteState"`

	// The ID of the VPC peering connection.
	VpcPeeringConnectionId *string `locationName:"vpcPeeringConnectionId" type:"string"`
}

// String returns the string representation
func (s Route) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Route) GoString() string {
	return s.String()
}

// SetDestinationCidrBlock sets the DestinationCidrBlock field's value.
func (s *Route) SetDestinationCidrBlock(v string) *Route {
	s.DestinationCidrBlock = &v
	return s
}

// SetDestinationIpv6CidrBlock sets the DestinationIpv6CidrBlock field's value.
func (s *Route) SetDestinationIpv6CidrBlock(v string) *Route {
	s.DestinationIpv6CidrBlock = &v
	return s
}

// SetDestinationPrefixListId sets the DestinationPrefixListId field's value.
func (s *Route) SetDestinationPrefixListId(v string) *Route {
	s.DestinationPrefixListId = &v
	return s
}

// SetEgressOnlyInternetGatewayId sets the EgressOnlyInternetGatewayId field's value.
func (s *Route) SetEgressOnlyInternetGatewayId(v string) *Route {
	s.EgressOnlyInternetGatewayId = &v
	return s
}

// SetGatewayId sets the GatewayId field's value.
func (s *Route) SetGatewayId(v string) *Route {
	s.GatewayId = &v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *Route) SetInstanceId(v string) *Route {
	s.InstanceId = &v
	return s
}

// SetInstanceOwnerId sets the InstanceOwnerId field's value.
func (s *Route) SetInstanceOwnerId(v string) *Route {
	s.InstanceOwnerId = &v
	return s
}

// SetNatGatewayId sets the NatGatewayId field's value.
func (s *Route) SetNatGatewayId(v string) *Route {
	s.NatGatewayId = &v
	return s
}

// SetNetworkInterfaceId sets the NetworkInterfaceId field's value.
func (s *Route) SetNetworkInterfaceId(v string) *Route {
	s.NetworkInterfaceId = &v
	return s
}

// SetOrigin sets the Origin field's value.
func (s *Route) SetOrigin(v string) *Route {
	s.Origin = &v
	return s
}

// SetState sets the State field's value.
func (s *Route) SetState(v string) *Route {
	s.State = &v
	return s
}

// SetVpcPeeringConnectionId sets the VpcPeeringConnectionId field's value.
func (s *Route) SetVpcPeeringConnectionId(v string) *Route {
	s.VpcPeeringConnectionId = &v
	return s
}

type ReplaceRouteInput struct {
	_ struct{} `type:"structure"`

	// The IPv4 CIDR address block used for the destination match. The value you
	// provide must match the CIDR of an existing route in the table.
	DestinationCidrBlock *string `locationName:"destinationCidrBlock" type:"string"`

	// The IPv6 CIDR address block used for the destination match. The value you
	// provide must match the CIDR of an existing route in the table.
	DestinationIpv6CidrBlock *string `locationName:"destinationIpv6CidrBlock" type:"string"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// [IPv6 traffic only] The ID of an egress-only Internet gateway.
	EgressOnlyInternetGatewayId *string `locationName:"egressOnlyInternetGatewayId" type:"string"`

	// The ID of an Internet gateway or virtual private gateway.
	GatewayId *string `locationName:"gatewayId" type:"string"`

	// The ID of a NAT instance in your VPC.
	InstanceId *string `locationName:"instanceId" type:"string"`

	// [IPv4 traffic only] The ID of a NAT gateway.
	NatGatewayId *string `locationName:"natGatewayId" type:"string"`

	// The ID of a network interface.
	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string"`

	// The ID of the route table.
	//
	// RouteTableId is a required field
	RouteTableId *string `locationName:"routeTableId" type:"string" required:"true"`

	// The ID of a VPC peering connection.
	VpcPeeringConnectionId *string `locationName:"vpcPeeringConnectionId" type:"string"`
}

// String returns the string representation
func (s ReplaceRouteInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ReplaceRouteInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *ReplaceRouteInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "ReplaceRouteInput"}
	if s.RouteTableId == nil {
		invalidParams.Add(request.NewErrParamRequired("RouteTableId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDestinationCidrBlock sets the DestinationCidrBlock field's value.
func (s *ReplaceRouteInput) SetDestinationCidrBlock(v string) *ReplaceRouteInput {
	s.DestinationCidrBlock = &v
	return s
}

// SetDestinationIpv6CidrBlock sets the DestinationIpv6CidrBlock field's value.
func (s *ReplaceRouteInput) SetDestinationIpv6CidrBlock(v string) *ReplaceRouteInput {
	s.DestinationIpv6CidrBlock = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *ReplaceRouteInput) SetDryRun(v bool) *ReplaceRouteInput {
	s.DryRun = &v
	return s
}

// SetEgressOnlyInternetGatewayId sets the EgressOnlyInternetGatewayId field's value.
func (s *ReplaceRouteInput) SetEgressOnlyInternetGatewayId(v string) *ReplaceRouteInput {
	s.EgressOnlyInternetGatewayId = &v
	return s
}

// SetGatewayId sets the GatewayId field's value.
func (s *ReplaceRouteInput) SetGatewayId(v string) *ReplaceRouteInput {
	s.GatewayId = &v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *ReplaceRouteInput) SetInstanceId(v string) *ReplaceRouteInput {
	s.InstanceId = &v
	return s
}

// SetNatGatewayId sets the NatGatewayId field's value.
func (s *ReplaceRouteInput) SetNatGatewayId(v string) *ReplaceRouteInput {
	s.NatGatewayId = &v
	return s
}

// SetNetworkInterfaceId sets the NetworkInterfaceId field's value.
func (s *ReplaceRouteInput) SetNetworkInterfaceId(v string) *ReplaceRouteInput {
	s.NetworkInterfaceId = &v
	return s
}

// SetRouteTableId sets the RouteTableId field's value.
func (s *ReplaceRouteInput) SetRouteTableId(v string) *ReplaceRouteInput {
	s.RouteTableId = &v
	return s
}

// SetVpcPeeringConnectionId sets the VpcPeeringConnectionId field's value.
func (s *ReplaceRouteInput) SetVpcPeeringConnectionId(v string) *ReplaceRouteInput {
	s.VpcPeeringConnectionId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/ReplaceRouteOutput
type ReplaceRouteOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s ReplaceRouteOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ReplaceRouteOutput) GoString() string {
	return s.String()
}

type DeleteRouteInput struct {
	_ struct{} `type:"structure"`

	// The IPv4 CIDR range for the route. The value you specify must match the CIDR
	// for the route exactly.
	DestinationCidrBlock *string `locationName:"destinationCidrBlock" type:"string"`

	// The IPv6 CIDR range for the route. The value you specify must match the CIDR
	// for the route exactly.
	DestinationIpv6CidrBlock *string `locationName:"destinationIpv6CidrBlock" type:"string"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the route table.
	//
	// RouteTableId is a required field
	RouteTableId *string `locationName:"routeTableId" type:"string" required:"true"`
}

// String returns the string representation
func (s DeleteRouteInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteRouteInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteRouteInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeleteRouteInput"}
	if s.RouteTableId == nil {
		invalidParams.Add(request.NewErrParamRequired("RouteTableId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDestinationCidrBlock sets the DestinationCidrBlock field's value.
func (s *DeleteRouteInput) SetDestinationCidrBlock(v string) *DeleteRouteInput {
	s.DestinationCidrBlock = &v
	return s
}

// SetDestinationIpv6CidrBlock sets the DestinationIpv6CidrBlock field's value.
func (s *DeleteRouteInput) SetDestinationIpv6CidrBlock(v string) *DeleteRouteInput {
	s.DestinationIpv6CidrBlock = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *DeleteRouteInput) SetDryRun(v bool) *DeleteRouteInput {
	s.DryRun = &v
	return s
}

// SetRouteTableId sets the RouteTableId field's value.
func (s *DeleteRouteInput) SetRouteTableId(v string) *DeleteRouteInput {
	s.RouteTableId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteRouteOutput
type DeleteRouteOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeleteRouteOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteRouteOutput) GoString() string {
	return s.String()
}

type DescribeRouteTablesInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	//    * vpc-id - The ID of the VPC for the route table.
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// One or more route table IDs.
	//
	// Default: Describes all your route tables.
	RouteTableIds []*string `locationName:"RouteTableId" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s DescribeRouteTablesInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeRouteTablesInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeRouteTablesInput) SetDryRun(v bool) *DescribeRouteTablesInput {
	s.DryRun = &v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribeRouteTablesInput) SetFilters(v []*Filter) *DescribeRouteTablesInput {
	s.Filters = v
	return s
}

// SetRouteTableIds sets the RouteTableIds field's value.
func (s *DescribeRouteTablesInput) SetRouteTableIds(v []*string) *DescribeRouteTablesInput {
	s.RouteTableIds = v
	return s
}

// Contains the output of DescribeRouteTables.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeRouteTablesResult
type DescribeRouteTablesOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more route tables.
	RouteTables []*RouteTable `locationName:"routeTableSet" locationNameList:"item" type:"list"`
	RequestId   *string       `locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s DescribeRouteTablesOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeRouteTablesOutput) GoString() string {
	return s.String()
}

// SetRouteTables sets the RouteTables field's value.
func (s *DescribeRouteTablesOutput) SetRouteTables(v []*RouteTable) *DescribeRouteTablesOutput {
	s.RouteTables = v
	return s
}

type RouteTable struct {
	_ struct{} `type:"structure"`

	// The associations between the route table and one or more subnets.
	Associations []*RouteTableAssociation `locationName:"associationSet" locationNameList:"item" type:"list"`

	// Any virtual private gateway (VGW) propagating routes.
	PropagatingVgws []*PropagatingVgw `locationName:"propagatingVgwSet" locationNameList:"item" type:"list"`

	// The ID of the route table.
	RouteTableId *string `locationName:"routeTableId" type:"string"`

	// The routes in the route table.
	Routes []*Route `locationName:"routeSet" locationNameList:"item" type:"list"`

	// Any tags assigned to the route table.
	Tags []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`

	// The ID of the VPC.
	VpcId *string `locationName:"vpcId" type:"string"`
}

// String returns the string representation
func (s RouteTable) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s RouteTable) GoString() string {
	return s.String()
}

// SetAssociations sets the Associations field's value.
func (s *RouteTable) SetAssociations(v []*RouteTableAssociation) *RouteTable {
	s.Associations = v
	return s
}

// SetPropagatingVgws sets the PropagatingVgws field's value.
func (s *RouteTable) SetPropagatingVgws(v []*PropagatingVgw) *RouteTable {
	s.PropagatingVgws = v
	return s
}

// SetRouteTableId sets the RouteTableId field's value.
func (s *RouteTable) SetRouteTableId(v string) *RouteTable {
	s.RouteTableId = &v
	return s
}

// SetRoutes sets the Routes field's value.
func (s *RouteTable) SetRoutes(v []*Route) *RouteTable {
	s.Routes = v
	return s
}

// SetTags sets the Tags field's value.
func (s *RouteTable) SetTags(v []*Tag) *RouteTable {
	s.Tags = v
	return s
}

// SetVpcId sets the VpcId field's value.
func (s *RouteTable) SetVpcId(v string) *RouteTable {
	s.VpcId = &v
	return s
}

// Describes an association between a route table and a subnet.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/RouteTableAssociation
type RouteTableAssociation struct {
	_ struct{} `type:"structure"`

	// Indicates whether this is the main route table.
	Main *bool `locationName:"main" type:"boolean"`

	// The ID of the association between a route table and a subnet.
	RouteTableAssociationId *string `locationName:"routeTableAssociationId" type:"string"`

	// The ID of the route table.
	RouteTableId *string `locationName:"routeTableId" type:"string"`

	// The ID of the subnet. A subnet ID is not returned for an implicit association.
	SubnetId *string `locationName:"subnetId" type:"string"`
}

// String returns the string representation
func (s RouteTableAssociation) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s RouteTableAssociation) GoString() string {
	return s.String()
}

// SetMain sets the Main field's value.
func (s *RouteTableAssociation) SetMain(v bool) *RouteTableAssociation {
	s.Main = &v
	return s
}

// SetRouteTableAssociationId sets the RouteTableAssociationId field's value.
func (s *RouteTableAssociation) SetRouteTableAssociationId(v string) *RouteTableAssociation {
	s.RouteTableAssociationId = &v
	return s
}

// SetRouteTableId sets the RouteTableId field's value.
func (s *RouteTableAssociation) SetRouteTableId(v string) *RouteTableAssociation {
	s.RouteTableId = &v
	return s
}

// SetSubnetId sets the SubnetId field's value.
func (s *RouteTableAssociation) SetSubnetId(v string) *RouteTableAssociation {
	s.SubnetId = &v
	return s
}

type PropagatingVgw struct {
	_ struct{} `type:"structure"`

	// The ID of the virtual private gateway (VGW).
	GatewayId *string `locationName:"gatewayId" type:"string"`
}

// String returns the string representation
func (s PropagatingVgw) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s PropagatingVgw) GoString() string {
	return s.String()
}

// SetGatewayId sets the GatewayId field's value.
func (s *PropagatingVgw) SetGatewayId(v string) *PropagatingVgw {
	s.GatewayId = &v
	return s
}

type CreateRouteTableInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the VPC.
	//
	// VpcId is a required field
	VpcId *string `locationName:"vpcId" type:"string" required:"true"`
}

// String returns the string representation
func (s CreateRouteTableInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateRouteTableInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CreateRouteTableInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "CreateRouteTableInput"}
	if s.VpcId == nil {
		invalidParams.Add(request.NewErrParamRequired("VpcId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDryRun sets the DryRun field's value.
func (s *CreateRouteTableInput) SetDryRun(v bool) *CreateRouteTableInput {
	s.DryRun = &v
	return s
}

// SetVpcId sets the VpcId field's value.
func (s *CreateRouteTableInput) SetVpcId(v string) *CreateRouteTableInput {
	s.VpcId = &v
	return s
}

// Contains the output of CreateRouteTable.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateRouteTableResult
type CreateRouteTableOutput struct {
	_ struct{} `type:"structure"`

	// Information about the route table.
	RouteTable *RouteTable `locationName:"routeTable" type:"structure"`

	RequestId *string `locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s CreateRouteTableOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateRouteTableOutput) GoString() string {
	return s.String()
}

// SetRouteTable sets the RouteTable field's value.
func (s *CreateRouteTableOutput) SetRouteTable(v *RouteTable) *CreateRouteTableOutput {
	s.RouteTable = v
	return s
}

type DisableVgwRoutePropagationInput struct {
	_ struct{} `type:"structure"`

	// The ID of the virtual private gateway.
	//
	// GatewayId is a required field
	GatewayId *string `type:"string" required:"true"`

	// The ID of the route table.
	//
	// RouteTableId is a required field
	RouteTableId *string `type:"string" required:"true"`
}

// String returns the string representation
func (s DisableVgwRoutePropagationInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DisableVgwRoutePropagationInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DisableVgwRoutePropagationInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DisableVgwRoutePropagationInput"}
	if s.GatewayId == nil {
		invalidParams.Add(request.NewErrParamRequired("GatewayId"))
	}
	if s.RouteTableId == nil {
		invalidParams.Add(request.NewErrParamRequired("RouteTableId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetGatewayId sets the GatewayId field's value.
func (s *DisableVgwRoutePropagationInput) SetGatewayId(v string) *DisableVgwRoutePropagationInput {
	s.GatewayId = &v
	return s
}

// SetRouteTableId sets the RouteTableId field's value.
func (s *DisableVgwRoutePropagationInput) SetRouteTableId(v string) *DisableVgwRoutePropagationInput {
	s.RouteTableId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DisableVgwRoutePropagationOutput
type DisableVgwRoutePropagationOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DisableVgwRoutePropagationOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DisableVgwRoutePropagationOutput) GoString() string {
	return s.String()
}

type EnableVgwRoutePropagationInput struct {
	_ struct{} `type:"structure"`

	// The ID of the virtual private gateway.
	//
	// GatewayId is a required field
	GatewayId *string `type:"string" required:"true"`

	// The ID of the route table.
	//
	// RouteTableId is a required field
	RouteTableId *string `type:"string" required:"true"`
}

// String returns the string representation
func (s EnableVgwRoutePropagationInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s EnableVgwRoutePropagationInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *EnableVgwRoutePropagationInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "EnableVgwRoutePropagationInput"}
	if s.GatewayId == nil {
		invalidParams.Add(request.NewErrParamRequired("GatewayId"))
	}
	if s.RouteTableId == nil {
		invalidParams.Add(request.NewErrParamRequired("RouteTableId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetGatewayId sets the GatewayId field's value.
func (s *EnableVgwRoutePropagationInput) SetGatewayId(v string) *EnableVgwRoutePropagationInput {
	s.GatewayId = &v
	return s
}

// SetRouteTableId sets the RouteTableId field's value.
func (s *EnableVgwRoutePropagationInput) SetRouteTableId(v string) *EnableVgwRoutePropagationInput {
	s.RouteTableId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/EnableVgwRoutePropagationOutput
type EnableVgwRoutePropagationOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s EnableVgwRoutePropagationOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s EnableVgwRoutePropagationOutput) GoString() string {
	return s.String()
}

type DisassociateRouteTableInput struct {
	_ struct{} `type:"structure"`

	// The association ID representing the current association between the route
	// table and subnet.
	//
	// AssociationId is a required field
	AssociationId *string `locationName:"associationId" type:"string" required:"true"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`
}

// String returns the string representation
func (s DisassociateRouteTableInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DisassociateRouteTableInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DisassociateRouteTableInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DisassociateRouteTableInput"}
	if s.AssociationId == nil {
		invalidParams.Add(request.NewErrParamRequired("AssociationId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAssociationId sets the AssociationId field's value.
func (s *DisassociateRouteTableInput) SetAssociationId(v string) *DisassociateRouteTableInput {
	s.AssociationId = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *DisassociateRouteTableInput) SetDryRun(v bool) *DisassociateRouteTableInput {
	s.DryRun = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DisassociateRouteTableOutput
type DisassociateRouteTableOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DisassociateRouteTableOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DisassociateRouteTableOutput) GoString() string {
	return s.String()
}

type DeleteRouteTableInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the route table.
	//
	// RouteTableId is a required field
	RouteTableId *string `locationName:"routeTableId" type:"string" required:"true"`
}

// String returns the string representation
func (s DeleteRouteTableInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteRouteTableInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteRouteTableInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeleteRouteTableInput"}
	if s.RouteTableId == nil {
		invalidParams.Add(request.NewErrParamRequired("RouteTableId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDryRun sets the DryRun field's value.
func (s *DeleteRouteTableInput) SetDryRun(v bool) *DeleteRouteTableInput {
	s.DryRun = &v
	return s
}

// SetRouteTableId sets the RouteTableId field's value.
func (s *DeleteRouteTableInput) SetRouteTableId(v string) *DeleteRouteTableInput {
	s.RouteTableId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteRouteTableOutput
type DeleteRouteTableOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeleteRouteTableOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteRouteTableOutput) GoString() string {
	return s.String()
}

type AssociateRouteTableInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the route table.
	//
	// RouteTableId is a required field
	RouteTableId *string `locationName:"routeTableId" type:"string" required:"true"`

	// The ID of the subnet.
	//
	// SubnetId is a required field
	SubnetId *string `locationName:"subnetId" type:"string" required:"true"`
}

// String returns the string representation
func (s AssociateRouteTableInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AssociateRouteTableInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *AssociateRouteTableInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "AssociateRouteTableInput"}
	if s.RouteTableId == nil {
		invalidParams.Add(request.NewErrParamRequired("RouteTableId"))
	}
	if s.SubnetId == nil {
		invalidParams.Add(request.NewErrParamRequired("SubnetId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDryRun sets the DryRun field's value.
func (s *AssociateRouteTableInput) SetDryRun(v bool) *AssociateRouteTableInput {
	s.DryRun = &v
	return s
}

// SetRouteTableId sets the RouteTableId field's value.
func (s *AssociateRouteTableInput) SetRouteTableId(v string) *AssociateRouteTableInput {
	s.RouteTableId = &v
	return s
}

// SetSubnetId sets the SubnetId field's value.
func (s *AssociateRouteTableInput) SetSubnetId(v string) *AssociateRouteTableInput {
	s.SubnetId = &v
	return s
}

// Contains the output of AssociateRouteTable.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/AssociateRouteTableResult
type AssociateRouteTableOutput struct {
	_ struct{} `type:"structure"`

	// The route table association ID (needed to disassociate the route table).
	AssociationId *string `locationName:"associationId" type:"string"`
}

// String returns the string representation
func (s AssociateRouteTableOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AssociateRouteTableOutput) GoString() string {
	return s.String()
}

// SetAssociationId sets the AssociationId field's value.
func (s *AssociateRouteTableOutput) SetAssociationId(v string) *AssociateRouteTableOutput {
	s.AssociationId = &v
	return s
}

type ReplaceRouteTableAssociationInput struct {
	_ struct{} `type:"structure"`

	// The association ID.
	//
	// AssociationId is a required field
	AssociationId *string `locationName:"associationId" type:"string" required:"true"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the new route table to associate with the subnet.
	//
	// RouteTableId is a required field
	RouteTableId *string `locationName:"routeTableId" type:"string" required:"true"`
}

// String returns the string representation
func (s ReplaceRouteTableAssociationInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ReplaceRouteTableAssociationInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *ReplaceRouteTableAssociationInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "ReplaceRouteTableAssociationInput"}
	if s.AssociationId == nil {
		invalidParams.Add(request.NewErrParamRequired("AssociationId"))
	}
	if s.RouteTableId == nil {
		invalidParams.Add(request.NewErrParamRequired("RouteTableId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAssociationId sets the AssociationId field's value.
func (s *ReplaceRouteTableAssociationInput) SetAssociationId(v string) *ReplaceRouteTableAssociationInput {
	s.AssociationId = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *ReplaceRouteTableAssociationInput) SetDryRun(v bool) *ReplaceRouteTableAssociationInput {
	s.DryRun = &v
	return s
}

// SetRouteTableId sets the RouteTableId field's value.
func (s *ReplaceRouteTableAssociationInput) SetRouteTableId(v string) *ReplaceRouteTableAssociationInput {
	s.RouteTableId = &v
	return s
}

// Contains the output of ReplaceRouteTableAssociation.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/ReplaceRouteTableAssociationResult
type ReplaceRouteTableAssociationOutput struct {
	_ struct{} `type:"structure"`

	// The ID of the new association.
	NewAssociationId *string `locationName:"newAssociationId" type:"string"`
}

// String returns the string representation
func (s ReplaceRouteTableAssociationOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ReplaceRouteTableAssociationOutput) GoString() string {
	return s.String()
}

// SetNewAssociationId sets the NewAssociationId field's value.
func (s *ReplaceRouteTableAssociationOutput) SetNewAssociationId(v string) *ReplaceRouteTableAssociationOutput {
	s.NewAssociationId = &v
	return s
}

// String returns the string representation
func (s DescribeVpcsInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeVpcsInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeVpcsInput) SetDryRun(v bool) *DescribeVpcsInput {
	s.DryRun = &v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribeVpcsInput) SetFilters(v []*Filter) *DescribeVpcsInput {
	s.Filters = v
	return s
}

// SetVpcIds sets the VpcIds field's value.
func (s *DescribeVpcsInput) SetVpcIds(v []*string) *DescribeVpcsInput {
	s.VpcIds = v
	return s
}

// // Contains the output of DescribeVpcs.
// // Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeVpcsResult
// type DescribeVpcsOutput struct {
// 	_ struct{} `type:"structure"`

// 	// Information about one or more VPCs.
// 	Vpcs []*Vpc `locationName:"vpcSet" locationNameList:"item" type:"list"`
// }

// // String returns the string representation
// func (s DescribeVpcsOutput) String() string {
// 	return awsutil.Prettify(s)
// }

// GoString returns the string representation
// func (s DescribeVpcsOutput) GoString() string {
// 	return s.String()
// }

// SetVpcs sets the Vpcs field's value.
func (s *DescribeVpcsOutput) SetVpcs(v []*Vpc) *DescribeVpcsOutput {
	s.Vpcs = v
	return s
}

type VpnConnectionOptionsSpecification struct {
	_ struct{} `type:"structure"`

	// Indicates whether the VPN connection uses static routes only. Static routes
	// must be used for devices that don't support BGP.
	StaticRoutesOnly *bool `locationName:"staticRoutesOnly" type:"boolean"`
}

// String returns the string representation
func (s VpnConnectionOptionsSpecification) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s VpnConnectionOptionsSpecification) GoString() string {
	return s.String()
}

// SetStaticRoutesOnly sets the StaticRoutesOnly field's value.
func (s *VpnConnectionOptionsSpecification) SetStaticRoutesOnly(v bool) *VpnConnectionOptionsSpecification {
	s.StaticRoutesOnly = &v
	return s
}

type CreateVpnGatewayInput struct {
	_ struct{} `type:"structure"`

	// The Availability Zone for the virtual private gateway.
	AvailabilityZone *string `type:"string"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The type of VPN connection this virtual private gateway supports.
	//
	// Type is a required field
	Type *string `type:"string" required:"true" enum:"GatewayType"`
}

// String returns the string representation
func (s CreateVpnGatewayInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateVpnGatewayInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CreateVpnGatewayInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "CreateVpnGatewayInput"}
	if s.Type == nil {
		invalidParams.Add(request.NewErrParamRequired("Type"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAvailabilityZone sets the AvailabilityZone field's value.
func (s *CreateVpnGatewayInput) SetAvailabilityZone(v string) *CreateVpnGatewayInput {
	s.AvailabilityZone = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *CreateVpnGatewayInput) SetDryRun(v bool) *CreateVpnGatewayInput {
	s.DryRun = &v
	return s
}

// SetType sets the Type field's value.
func (s *CreateVpnGatewayInput) SetType(v string) *CreateVpnGatewayInput {
	s.Type = &v
	return s
}

type CreateVpnConnectionInput struct {
	_ struct{} `type:"structure"`

	// The ID of the customer gateway.
	//
	// CustomerGatewayId is a required field
	CustomerGatewayId *string `type:"string" required:"true"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// Indicates whether the VPN connection requires static routes. If you are creating
	// a VPN connection for a device that does not support BGP, you must specify
	// true.
	//
	// Default: false
	Options *VpnConnectionOptionsSpecification `locationName:"options" type:"structure"`

	// The type of VPN connection (ipsec.1).
	//
	// Type is a required field
	Type *string `type:"string" required:"true"`

	// The ID of the virtual private gateway.
	//
	// VpnGatewayId is a required field
	VpnGatewayId *string `type:"string" required:"true"`
}

// String returns the string representation
func (s CreateVpnConnectionInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateVpnConnectionInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CreateVpnConnectionInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "CreateVpnConnectionInput"}
	if s.CustomerGatewayId == nil {
		invalidParams.Add(request.NewErrParamRequired("CustomerGatewayId"))
	}
	if s.Type == nil {
		invalidParams.Add(request.NewErrParamRequired("Type"))
	}
	if s.VpnGatewayId == nil {
		invalidParams.Add(request.NewErrParamRequired("VpnGatewayId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetCustomerGatewayId sets the CustomerGatewayId field's value.
func (s *CreateVpnConnectionInput) SetCustomerGatewayId(v string) *CreateVpnConnectionInput {
	s.CustomerGatewayId = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *CreateVpnConnectionInput) SetDryRun(v bool) *CreateVpnConnectionInput {
	s.DryRun = &v
	return s
}

// SetOptions sets the Options field's value.
func (s *CreateVpnConnectionInput) SetOptions(v *VpnConnectionOptionsSpecification) *CreateVpnConnectionInput {
	s.Options = v
	return s
}

// SetType sets the Type field's value.
func (s *CreateVpnConnectionInput) SetType(v string) *CreateVpnConnectionInput {
	s.Type = &v
	return s
}

// SetVpnGatewayId sets the VpnGatewayId field's value.
func (s *CreateVpnConnectionInput) SetVpnGatewayId(v string) *CreateVpnConnectionInput {
	s.VpnGatewayId = &v
	return s
}

// Contains the output of CreateVpnConnection.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateVpnConnectionResult
type CreateVpnConnectionOutput struct {
	_ struct{} `type:"structure"`

	// Information about the VPN connection.
	VpnConnection *VpnConnection `locationName:"vpnConnection" type:"structure"`
}

// String returns the string representation
func (s CreateVpnConnectionOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateVpnConnectionOutput) GoString() string {
	return s.String()
}

// SetVpnConnection sets the VpnConnection field's value.
func (s *CreateVpnConnectionOutput) SetVpnConnection(v *VpnConnection) *CreateVpnConnectionOutput {
	s.VpnConnection = v
	return s
}

// Contains the output of CreateVpnGateway.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateVpnGatewayResult
type CreateVpnGatewayOutput struct {
	_ struct{} `type:"structure"`

	// Information about the virtual private gateway.
	VpnGateway *VpnGateway `locationName:"vpnGateway" type:"structure"`
}

// String returns the string representation
func (s CreateVpnGatewayOutput) String() string {
	return awsutil.Prettify(s)
}

type VpnConnection struct {
	_ struct{} `type:"structure"`

	// The configuration information for the VPN connection's customer gateway (in
	// the native XML format). This element is always present in the CreateVpnConnection
	// response; however, it's present in the DescribeVpnConnections response only
	// if the VPN connection is in the pending or available state.
	CustomerGatewayConfiguration *string `locationName:"customerGatewayConfiguration" type:"string"`

	// The ID of the customer gateway at your end of the VPN connection.
	CustomerGatewayId *string `locationName:"customerGatewayId" type:"string"`

	// The VPN connection options.
	Options *VpnConnectionOptions `locationName:"options" type:"structure"`

	// The static routes associated with the VPN connection.
	Routes []*VpnStaticRoute `locationName:"routes" locationNameList:"item" type:"list"`

	// The current state of the VPN connection.
	State *string `locationName:"state" type:"string" enum:"VpnState"`

	// Any tags assigned to the VPN connection.
	Tags []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`

	// The type of VPN connection.
	Type *string `locationName:"type" type:"string" enum:"GatewayType"`

	// Information about the VPN tunnel.
	VgwTelemetry []*VgwTelemetry `locationName:"vgwTelemetry" locationNameList:"item" type:"list"`

	// The ID of the VPN connection.
	VpnConnectionId *string `locationName:"vpnConnectionId" type:"string"`

	// The ID of the virtual private gateway at the AWS side of the VPN connection.
	VpnGatewayId *string `locationName:"vpnGatewayId" type:"string"`
}

// String returns the string representation
func (s VpnConnection) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s VpnConnection) GoString() string {
	return s.String()
}

// SetCustomerGatewayConfiguration sets the CustomerGatewayConfiguration field's value.
func (s *VpnConnection) SetCustomerGatewayConfiguration(v string) *VpnConnection {
	s.CustomerGatewayConfiguration = &v
	return s
}

// SetCustomerGatewayId sets the CustomerGatewayId field's value.
func (s *VpnConnection) SetCustomerGatewayId(v string) *VpnConnection {
	s.CustomerGatewayId = &v
	return s
}

// SetOptions sets the Options field's value.
func (s *VpnConnection) SetOptions(v *VpnConnectionOptions) *VpnConnection {
	s.Options = v
	return s
}

// SetRoutes sets the Routes field's value.
func (s *VpnConnection) SetRoutes(v []*VpnStaticRoute) *VpnConnection {
	s.Routes = v
	return s
}

// SetState sets the State field's value.
func (s *VpnConnection) SetState(v string) *VpnConnection {
	s.State = &v
	return s
}

// SetTags sets the Tags field's value.
func (s *VpnConnection) SetTags(v []*Tag) *VpnConnection {
	s.Tags = v
	return s
}

// SetType sets the Type field's value.
func (s *VpnConnection) SetType(v string) *VpnConnection {
	s.Type = &v
	return s
}

// SetVgwTelemetry sets the VgwTelemetry field's value.
func (s *VpnConnection) SetVgwTelemetry(v []*VgwTelemetry) *VpnConnection {
	s.VgwTelemetry = v
	return s
}

// SetVpnConnectionId sets the VpnConnectionId field's value.
func (s *VpnConnection) SetVpnConnectionId(v string) *VpnConnection {
	s.VpnConnectionId = &v
	return s
}

// SetVpnGatewayId sets the VpnGatewayId field's value.
func (s *VpnConnection) SetVpnGatewayId(v string) *VpnConnection {
	s.VpnGatewayId = &v
	return s
}
func (s CreateVpnGatewayOutput) GoString() string {
	return s.String()
}

// SetVpnGateway sets the VpnGateway field's value.
func (s *CreateVpnGatewayOutput) SetVpnGateway(v *VpnGateway) *CreateVpnGatewayOutput {
	s.VpnGateway = v
	return s
}

type VpnGateway struct {
	_ struct{} `type:"structure"`

	// The Availability Zone where the virtual private gateway was created, if applicable.
	// This field may be empty or not returned.
	AvailabilityZone *string `locationName:"availabilityZone" type:"string"`

	// The current state of the virtual private gateway.
	State *string `locationName:"state" type:"string" enum:"VpnState"`

	// Any tags assigned to the virtual private gateway.
	Tags []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`

	// The type of VPN connection the virtual private gateway supports.
	Type *string `locationName:"type" type:"string" enum:"GatewayType"`

	// Any VPCs attached to the virtual private gateway.
	VpcAttachments []*VpcAttachment `locationName:"attachments" locationNameList:"item" type:"list"`

	// The ID of the virtual private gateway.
	VpnGatewayId *string `locationName:"vpnGatewayId" type:"string"`
}

func (s VpnGateway) String() string {
	return awsutil.Prettify(s)
}

func (s VpnGateway) GoString() string {
	return s.String()
}

// SetAvailabilityZone sets the AvailabilityZone field's value.
func (s *VpnGateway) SetAvailabilityZone(v string) *VpnGateway {
	s.AvailabilityZone = &v
	return s
}

func (s *VpnGateway) SetState(v string) *VpnGateway {
	s.State = &v
	return s
}

func (s *VpnGateway) SetTags(v []*Tag) *VpnGateway {
	s.Tags = v
	return s
}

func (s *VpnGateway) SetType(v string) *VpnGateway {
	s.Type = &v
	return s
}

// SetVpcAttachments sets the VpcAttachments field's value.
func (s *VpnGateway) SetVpcAttachments(v []*VpcAttachment) *VpnGateway {
	s.VpcAttachments = v
	return s
}

func (s *VpnGateway) SetVpnGatewayId(v string) *VpnGateway {
	s.VpnGatewayId = &v
	return s
}

type VpnConnectionOptions struct {
	_ struct{} `type:"structure"`

	// Indicates whether the VPN connection uses static routes only. Static routes
	// must be used for devices that don't support BGP.
	StaticRoutesOnly *bool `locationName:"staticRoutesOnly" type:"boolean"`
}

type VpcAttachment struct {
	_ struct{} `type:"structure"`

	// The current state of the attachment.
	State *string `locationName:"state" type:"string" enum:"AttachmentStatus"`

	// The ID of the VPC.
	VpcId *string `locationName:"vpcId" type:"string"`
}

// String returns the string representation
func (s VpcAttachment) String() string {
	return awsutil.Prettify(s)
}

// SetStaticRoutesOnly sets the StaticRoutesOnly field's value.
func (s *VpnConnectionOptions) SetStaticRoutesOnly(v bool) *VpnConnectionOptions {
	s.StaticRoutesOnly = &v
	return s
}

type VpnStaticRoute struct {
	_ struct{} `type:"structure"`

	// The CIDR block associated with the local subnet of the customer data center.
	DestinationCidrBlock *string `locationName:"destinationCidrBlock" type:"string"`

	// Indicates how the routes were provided.
	Source *string `locationName:"source" type:"string" enum:"VpnStaticRouteSource"`

	// The current state of the static route.
	State *string `locationName:"state" type:"string" enum:"VpnState"`
}

// String returns the string representation
func (s VpnStaticRoute) String() string {
	return s.String()
}
func (s VpcAttachment) GoString() string {
	return s.String()
}

// SetState sets the State field's value.
func (s *VpcAttachment) SetState(v string) *VpcAttachment {
	s.State = &v
	return s
}

// SetVpcId sets the VpcId field's value.
func (s *VpcAttachment) SetVpcId(v string) *VpcAttachment {
	s.VpcId = &v
	return s
}

type DescribeVpnGatewaysInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more filters.
	//
	//    * attachment.state - The current state of the attachment between the gateway
	//    and the VPC (attaching | attached | detaching | detached).
	//
	//    * attachment.vpc-id - The ID of an attached VPC.
	//
	//    * availability-zone - The Availability Zone for the virtual private gateway
	//    (if applicable).
	//
	//    * state - The state of the virtual private gateway (pending | available
	//    | deleting | deleted).
	//
	//    * tag:key=value - The key/value combination of a tag assigned to the resource.
	//    Specify the key of the tag in the filter name and the value of the tag
	//    in the filter value. For example, for the tag Purpose=X, specify tag:Purpose
	//    for the filter name and X for the filter value.
	//
	//    * tag-key - The key of a tag assigned to the resource. This filter is
	//    independent of the tag-value filter. For example, if you use both the
	//    filter "tag-key=Purpose" and the filter "tag-value=X", you get any resources
	//    assigned both the tag key Purpose (regardless of what the tag's value
	//    is), and the tag value X (regardless of what the tag's key is). If you
	//    want to list only resources where Purpose is X, see the tag:key=value
	//    filter.
	//
	//    * tag-value - The value of a tag assigned to the resource. This filter
	//    is independent of the tag-key filter.
	//
	//    * type - The type of virtual private gateway. Currently the only supported
	//    type is ipsec.1.
	//
	//    * vpn-gateway-id - The ID of the virtual private gateway.
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// One or more virtual private gateway IDs.
	//
	// Default: Describes all your virtual private gateways.
	VpnGatewayIds []*string `locationName:"VpnGatewayId" locationNameList:"VpnGatewayId" type:"list"`
}

// String returns the string representation
func (s DescribeVpnGatewaysInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s VpnStaticRoute) GoString() string {
	return s.String()
}

// SetDestinationCidrBlock sets the DestinationCidrBlock field's value.
func (s *VpnStaticRoute) SetDestinationCidrBlock(v string) *VpnStaticRoute {
	s.DestinationCidrBlock = &v
	return s
}

// SetSource sets the Source field's value.
func (s *VpnStaticRoute) SetSource(v string) *VpnStaticRoute {
	s.Source = &v
	return s
}

// SetState sets the State field's value.
func (s *VpnStaticRoute) SetState(v string) *VpnStaticRoute {
	s.State = &v
	return s
}

type VgwTelemetry struct {
	_ struct{} `type:"structure"`

	// The number of accepted routes.
	AcceptedRouteCount *int64 `locationName:"acceptedRouteCount" type:"integer"`

	// The date and time of the last change in status.
	LastStatusChange *time.Time `locationName:"lastStatusChange" type:"timestamp" timestampFormat:"iso8601"`

	// The Internet-routable IP address of the virtual private gateway's outside
	// interface.
	OutsideIpAddress *string `locationName:"outsideIpAddress" type:"string"`

	// The status of the VPN tunnel.
	Status *string `locationName:"status" type:"string" enum:"TelemetryStatus"`

	// If an error occurs, a description of the error.
	StatusMessage *string `locationName:"statusMessage" type:"string"`
}

// String returns the string representation
func (s VgwTelemetry) String() string {
	return s.String()
}
func (s DescribeVpnGatewaysInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeVpnGatewaysInput) SetDryRun(v bool) *DescribeVpnGatewaysInput {
	s.DryRun = &v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribeVpnGatewaysInput) SetFilters(v []*Filter) *DescribeVpnGatewaysInput {
	s.Filters = v
	return s
}

// SetVpnGatewayIds sets the VpnGatewayIds field's value.
func (s *DescribeVpnGatewaysInput) SetVpnGatewayIds(v []*string) *DescribeVpnGatewaysInput {
	s.VpnGatewayIds = v
	return s
}

// Contains the output of DescribeVpnGateways.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeVpnGatewaysResult
type DescribeVpnGatewaysOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more virtual private gateways.
	VpnGateways []*VpnGateway `locationName:"vpnGatewaySet" locationNameList:"item" type:"list"`
	RequestId   *string       `locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s DescribeVpnGatewaysOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeVpnGatewaysOutput) GoString() string {
	return s.String()
}

// SetVpnGateways sets the VpnGateways field's value.
func (s *DescribeVpnGatewaysOutput) SetVpnGateways(v []*VpnGateway) *DescribeVpnGatewaysOutput {
	s.VpnGateways = v
	return s
}

type DeleteVpnGatewayInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the virtual private gateway.
	//
	// VpnGatewayId is a required field
	VpnGatewayId *string `type:"string" required:"true"`
}

// String returns the string representation
func (s DeleteVpnGatewayInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s VgwTelemetry) GoString() string {
	return s.String()
}

// SetAcceptedRouteCount sets the AcceptedRouteCount field's value.
func (s *VgwTelemetry) SetAcceptedRouteCount(v int64) *VgwTelemetry {
	s.AcceptedRouteCount = &v
	return s
}

// SetLastStatusChange sets the LastStatusChange field's value.
func (s *VgwTelemetry) SetLastStatusChange(v time.Time) *VgwTelemetry {
	s.LastStatusChange = &v
	return s
}

// SetOutsideIpAddress sets the OutsideIpAddress field's value.
func (s *VgwTelemetry) SetOutsideIpAddress(v string) *VgwTelemetry {
	s.OutsideIpAddress = &v
	return s
}

// SetStatus sets the Status field's value.
func (s *VgwTelemetry) SetStatus(v string) *VgwTelemetry {
	s.Status = &v
	return s
}

// SetStatusMessage sets the StatusMessage field's value.
func (s *VgwTelemetry) SetStatusMessage(v string) *VgwTelemetry {
	s.StatusMessage = &v
	return s
}

type DescribeVpnConnectionsInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more filters.
	//
	//    * customer-gateway-configuration - The configuration information for the
	//    customer gateway.
	//
	//    * customer-gateway-id - The ID of a customer gateway associated with the
	//    VPN connection.
	//
	//    * state - The state of the VPN connection (pending | available | deleting
	//    | deleted).
	//
	//    * option.static-routes-only - Indicates whether the connection has static
	//    routes only. Used for devices that do not support Border Gateway Protocol
	//    (BGP).
	//
	//    * route.destination-cidr-block - The destination CIDR block. This corresponds
	//    to the subnet used in a customer data center.
	//
	//    * bgp-asn - The BGP Autonomous System Number (ASN) associated with a BGP
	//    device.
	//
	//    * tag:key=value - The key/value combination of a tag assigned to the resource.
	//    Specify the key of the tag in the filter name and the value of the tag
	//    in the filter value. For example, for the tag Purpose=X, specify tag:Purpose
	//    for the filter name and X for the filter value.
	//
	//    * tag-key - The key of a tag assigned to the resource. This filter is
	//    independent of the tag-value filter. For example, if you use both the
	//    filter "tag-key=Purpose" and the filter "tag-value=X", you get any resources
	//    assigned both the tag key Purpose (regardless of what the tag's value
	//    is), and the tag value X (regardless of what the tag's key is). If you
	//    want to list only resources where Purpose is X, see the tag:key=value
	//    filter.
	//
	//    * tag-value - The value of a tag assigned to the resource. This filter
	//    is independent of the tag-key filter.
	//
	//    * type - The type of VPN connection. Currently the only supported type
	//    is ipsec.1.
	//
	//    * vpn-connection-id - The ID of the VPN connection.
	//
	//    * vpn-gateway-id - The ID of a virtual private gateway associated with
	//    the VPN connection.
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// One or more VPN connection IDs.
	//
	// Default: Describes your VPN connections.
	VpnConnectionIds []*string `locationName:"VpnConnectionId" locationNameList:"VpnConnectionId" type:"list"`
}

// String returns the string representation
func (s DescribeVpnConnectionsInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeVpnConnectionsInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeVpnConnectionsInput) SetDryRun(v bool) *DescribeVpnConnectionsInput {
	s.DryRun = &v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribeVpnConnectionsInput) SetFilters(v []*Filter) *DescribeVpnConnectionsInput {
	s.Filters = v
	return s
}

// SetVpnConnectionIds sets the VpnConnectionIds field's value.
func (s *DescribeVpnConnectionsInput) SetVpnConnectionIds(v []*string) *DescribeVpnConnectionsInput {
	s.VpnConnectionIds = v
	return s
}

// Contains the output of DescribeVpnConnections.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeVpnConnectionsResult
type DescribeVpnConnectionsOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more VPN connections.
	VpnConnections []*VpnConnection `locationName:"vpnConnectionSet" locationNameList:"item" type:"list"`
	RequestId      *string          `locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s DescribeVpnConnectionsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeVpnConnectionsOutput) GoString() string {
	return s.String()
}

// SetVpnConnections sets the VpnConnections field's value.
func (s *DescribeVpnConnectionsOutput) SetVpnConnections(v []*VpnConnection) *DescribeVpnConnectionsOutput {
	s.VpnConnections = v
	return s
}

func (s DeleteVpnGatewayInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteVpnGatewayInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeleteVpnGatewayInput"}
	if s.VpnGatewayId == nil {
		invalidParams.Add(request.NewErrParamRequired("VpnGatewayId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDryRun sets the DryRun field's value.
func (s *DeleteVpnGatewayInput) SetDryRun(v bool) *DeleteVpnGatewayInput {
	s.DryRun = &v
	return s
}

// SetVpnGatewayId sets the VpnGatewayId field's value.
func (s *DeleteVpnGatewayInput) SetVpnGatewayId(v string) *DeleteVpnGatewayInput {
	s.VpnGatewayId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteVpnGatewayOutput
type DeleteVpnGatewayOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeleteVpnGatewayOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteVpnGatewayOutput) GoString() string {
	return s.String()
}

type AttachVpnGatewayInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the VPC.
	//
	// VpcId is a required field
	VpcId *string `type:"string" required:"true"`

	// The ID of the virtual private gateway.
	//
	// VpnGatewayId is a required field
	VpnGatewayId *string `type:"string" required:"true"`
}

// String returns the string representation
func (s AttachVpnGatewayInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AttachVpnGatewayInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *AttachVpnGatewayInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "AttachVpnGatewayInput"}
	if s.VpcId == nil {
		invalidParams.Add(request.NewErrParamRequired("VpcId"))
	}
	if s.VpnGatewayId == nil {
		invalidParams.Add(request.NewErrParamRequired("VpnGatewayId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDryRun sets the DryRun field's value.
func (s *AttachVpnGatewayInput) SetDryRun(v bool) *AttachVpnGatewayInput {
	s.DryRun = &v
	return s
}

// SetVpcId sets the VpcId field's value.
func (s *AttachVpnGatewayInput) SetVpcId(v string) *AttachVpnGatewayInput {
	s.VpcId = &v
	return s
}

// SetVpnGatewayId sets the VpnGatewayId field's value.
func (s *AttachVpnGatewayInput) SetVpnGatewayId(v string) *AttachVpnGatewayInput {
	s.VpnGatewayId = &v
	return s
}

// Contains the output of AttachVpnGateway.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/AttachVpnGatewayResult
type AttachVpnGatewayOutput struct {
	_ struct{} `type:"structure"`

	// Information about the attachment.
	VpcAttachment *VpcAttachment `locationName:"attachment" type:"structure"`
}

// String returns the string representation
func (s AttachVpnGatewayOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AttachVpnGatewayOutput) GoString() string {
	return s.String()
}

// SetVpcAttachment sets the VpcAttachment field's value.
func (s *AttachVpnGatewayOutput) SetVpcAttachment(v *VpcAttachment) *AttachVpnGatewayOutput {
	s.VpcAttachment = v
	return s
}

type DeleteVpnConnectionInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the VPN connection.
	//
	// VpnConnectionId is a required field
	VpnConnectionId *string `type:"string" required:"true"`
}

// String returns the string representation
func (s DeleteVpnConnectionInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteVpnConnectionInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteVpnConnectionInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeleteVpnConnectionInput"}
	if s.VpnConnectionId == nil {
		invalidParams.Add(request.NewErrParamRequired("VpnConnectionId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDryRun sets the DryRun field's value.
func (s *DeleteVpnConnectionInput) SetDryRun(v bool) *DeleteVpnConnectionInput {
	s.DryRun = &v
	return s
}

// SetVpnConnectionId sets the VpnConnectionId field's value.
func (s *DeleteVpnConnectionInput) SetVpnConnectionId(v string) *DeleteVpnConnectionInput {
	s.VpnConnectionId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteVpnConnectionOutput
type DeleteVpnConnectionOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeleteVpnConnectionOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteVpnConnectionOutput) GoString() string {
	return s.String()
}

type DetachVpnGatewayInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the VPC.
	//
	// VpcId is a required field
	VpcId *string `type:"string" required:"true"`

	// The ID of the virtual private gateway.
	//
	// VpnGatewayId is a required field
	VpnGatewayId *string `type:"string" required:"true"`
}

// String returns the string representation
func (s DetachVpnGatewayInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DetachVpnGatewayInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DetachVpnGatewayInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DetachVpnGatewayInput"}
	if s.VpcId == nil {
		invalidParams.Add(request.NewErrParamRequired("VpcId"))
	}
	if s.VpnGatewayId == nil {
		invalidParams.Add(request.NewErrParamRequired("VpnGatewayId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDryRun sets the DryRun field's value.
func (s *DetachVpnGatewayInput) SetDryRun(v bool) *DetachVpnGatewayInput {
	s.DryRun = &v
	return s
}

// SetVpcId sets the VpcId field's value.
func (s *DetachVpnGatewayInput) SetVpcId(v string) *DetachVpnGatewayInput {
	s.VpcId = &v
	return s
}

// SetVpnGatewayId sets the VpnGatewayId field's value.
func (s *DetachVpnGatewayInput) SetVpnGatewayId(v string) *DetachVpnGatewayInput {
	s.VpnGatewayId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DetachVpnGatewayOutput
type DetachVpnGatewayOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DetachVpnGatewayOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DetachVpnGatewayOutput) GoString() string {
	return s.String()
}

type DescribeSnapshotExportTasksInput struct {
	SnapshotExportTaskId []*string `locationName:"snapshotExportTaskId" locationNameList:"item" type:"list"`
}
type DescribeSnapshotExportTasksOutput struct {
	SnapshotExportTask []*SnapshotExportTask `locationName:"snapshotExportTaskSet" locationNameList:"item" type:"list"`
	RequestId          *string               `locationName:"requestId" type:"string"`
}

type CreateSnapshotExportTaskInput struct {
	_           struct{}                      `type:"structure"`
	ExportToOsu *ExportToOsuTaskSpecification `locationName:"exportToOsu" type:"structure"`
	SnapshotId  *string                       `locationName:"snapshotId" type:"string"`
}

type CreateSnapshotExportTaskOutput struct {
	_                  struct{}            `type:"structure"`
	SnapshotExportTask *SnapshotExportTask `locationName:"snapshotExportTask" type:"structure"`
	RequestId          *string             `locationName:"requestId" type:"string"`
}

type SnapshotExportTask struct {
	_                    struct{}                      `type:"structure"`
	Completion           *int64                        `locationName:"completion" type:"string"`
	ExportToOsu          *ExportToOsuTaskSpecification `locationName:"exportToOsu" type:"structure"`
	SnapshotExport       *SnapshotExport               `locationName:"snapshotExport" type:"structure"`
	SnapshotExportTaskId *string                       `locationName:"snapshotExportTaskId" type:"string"`
	SnapshotId           *string                       `locationName:"SnapshotId" type:"string"`
	State                *string                       `locationName:"state" type:"string"`
	StatusMessage        *string                       `locationName:"statusMessage" type:"string"`
}

type SnapshotExport struct {
	SnapshotId *string `locationName:"snapshotId" type:"string"`
}

type ExportToOsuTaskSpecification struct {
	_               struct{}                           `type:"structure"`
	DiskImageFormat *string                            `locationName:"diskImageFormat" type:"string"`
	AkSk            *ExportToOsuAccessKeySpecification `locationName:"akSk" type:"structure"`
	OsuBucket       *string                            `locationName:"osuBucket" type:"string"`
	OsuKey          *string                            `locationName:"osuKey" type:"string"`
	OsuPrefix       *string                            `locationName:"osuPrefix" type:"string"`
}

type CreateImageExportTaskInput struct {
	_           struct{}                           `type:"structure"`
	ExportToOsu *ImageExportToOsuTaskSpecification `locationName:"exportToOsu" type:"structure"`
	ImageId     *string                            `locationName:"imageId" type:"string"`
}

type ImageExportToOsuTaskSpecification struct {
	_               struct{}                           `type:"structure"`
	DiskImageFormat *string                            `locationName:"diskImageFormat" type:"string"`
	OsuAkSk         *ExportToOsuAccessKeySpecification `locationName:"osuAkSk" type:"structure"`
	OsuBucket       *string                            `locationName:"osuBucket" type:"string"`
	OsuManifestUrl  *string                            `locationName:"osuManifestUrl" type:"string"`
	OsuPrefix       *string                            `locationName:"osuPrefix" type:"string"`
}

type ExportToOsuAccessKeySpecification struct {
	_         struct{} `type:"structure"`
	AccessKey *string  `locationName:"accessKey" type:"string"`
	SecretKey *string  `locationName:"secretKey" type:"string"`
}

type CreateImageExportTaskOutput struct {
	_               struct{}         `type:"structure"`
	ImageExportTask *ImageExportTask `locationName:"imageExportTask" type:"structure"`
	RequestId       *string          `locationName:"requestId" type:"string"`
}

type ImageExportTask struct {
	_                 struct{}                           `type:"structure"`
	Completion        *int64                             `locationName:"completion" type:"string"`
	ExportToOsu       *ImageExportToOsuTaskSpecification `locationName:"exportToOsu" type:"structure"`
	ImageExport       *ImageExport                       `locationName:"imageExport" type:"structure"`
	ImageExportTaskId *string                            `locationName:"imageExportTaskId" type:"string"`
	ImageId           *string                            `locationName:"imageId" type:"string"`
	State             *string                            `locationName:"state" type:"string"`
	StatusMessage     *string                            `locationName:"statusMessage" type:"string"`
}

type ImageExport struct {
	_       struct{} `type:"structure"`
	ImageId *string  `locationName:"imageId" type:"string"`
}

type DescribeImageExportTasksInput struct {
	_                 struct{}  `type:"structure"`
	ImageExportTaskId []*string `locationName:"imageExportTaskId" locationNameList:"item" type:"list"`
}

type DescribeImageExportTasksOutput struct {
	_               struct{}           `type:"structure"`
	ImageExportTask []*ImageExportTask `locationName:"imageExportTask" locationNameList:"item" type:"list"`
	RequestId       *string            `locationName:"requestId" type:"string"`
}

// Contains the parameters for CopyImage.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CopyImageRequest
type CopyImageInput struct {
	_ struct{} `type:"structure"`

	// Unique, case-sensitive identifier you provide to ensure idempotency of the
	// request. For more information, see How to Ensure Idempotency (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Run_Instance_Idempotency.html)
	// in the Amazon Elastic Compute Cloud User Guide.
	ClientToken *string `type:"string"`

	// A description for the new AMI in the destination region.
	Description *string `type:"string"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// Specifies whether the destination snapshots of the copied image should be
	// encrypted. The default CMK for EBS is used unless a non-default AWS Key Management
	// Service (AWS KMS) CMK is specified with KmsKeyId. For more information, see
	// Amazon EBS Encryption (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSEncryption.html)
	// in the Amazon Elastic Compute Cloud User Guide.
	Encrypted *bool `locationName:"encrypted" type:"boolean"`

	// The full ARN of the AWS Key Management Service (AWS KMS) CMK to use when
	// encrypting the snapshots of an image during a copy operation. This parameter
	// is only required if you want to use a non-default CMK; if this parameter
	// is not specified, the default CMK for EBS is used. The ARN contains the arn:aws:kms
	// namespace, followed by the region of the CMK, the AWS account ID of the CMK
	// owner, the key namespace, and then the CMK ID. For example, arn:aws:kms:us-east-1:012345678910:key/abcd1234-a123-456a-a12b-a123b4cd56ef.
	// The specified CMK must exist in the region that the snapshot is being copied
	// to. If a KmsKeyId is specified, the Encrypted flag must also be set.
	KmsKeyId *string `locationName:"kmsKeyId" type:"string"`

	// The name of the new AMI in the destination region.
	//
	// Name is a required field
	Name *string `type:"string" required:"true"`

	// The ID of the AMI to copy.
	//
	// SourceImageId is a required field
	SourceImageId *string `type:"string" required:"true"`

	// The name of the region that contains the AMI to copy.
	//
	// SourceRegion is a required field
	SourceRegion *string `type:"string" required:"true"`
}

// String returns the string representation
func (s CopyImageInput) String() string {
	return awsutil.Prettify(s)
}

type CreateVpnConnectionRouteInput struct {
	_ struct{} `type:"structure"`

	// The CIDR block associated with the local subnet of the customer network.
	//
	// DestinationCidrBlock is a required field
	DestinationCidrBlock *string `type:"string" required:"true"`

	// The ID of the VPN connection.
	//
	// VpnConnectionId is a required field
	VpnConnectionId *string `type:"string" required:"true"`
}

// String returns the string representation
func (s CreateVpnConnectionRouteInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CopyImageInput) GoString() string {
	return s.String()
}
func (s CreateVpnConnectionRouteInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CopyImageInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "CopyImageInput"}
	if s.Name == nil {
		invalidParams.Add(request.NewErrParamRequired("Name"))
	}
	if s.SourceImageId == nil {
		invalidParams.Add(request.NewErrParamRequired("SourceImageId"))
	}
	if s.SourceRegion == nil {
		invalidParams.Add(request.NewErrParamRequired("SourceRegion"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}
func (s *CreateVpnConnectionRouteInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "CreateVpnConnectionRouteInput"}
	if s.DestinationCidrBlock == nil {
		invalidParams.Add(request.NewErrParamRequired("DestinationCidrBlock"))
	}
	if s.VpnConnectionId == nil {
		invalidParams.Add(request.NewErrParamRequired("VpnConnectionId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetClientToken sets the ClientToken field's value.
func (s *CopyImageInput) SetClientToken(v string) *CopyImageInput {
	s.ClientToken = &v
	return s
}

// SetDescription sets the Description field's value.
func (s *CopyImageInput) SetDescription(v string) *CopyImageInput {
	s.Description = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *CopyImageInput) SetDryRun(v bool) *CopyImageInput {
	s.DryRun = &v
	return s
}

// SetEncrypted sets the Encrypted field's value.
func (s *CopyImageInput) SetEncrypted(v bool) *CopyImageInput {
	s.Encrypted = &v
	return s
}

// SetKmsKeyId sets the KmsKeyId field's value.
func (s *CopyImageInput) SetKmsKeyId(v string) *CopyImageInput {
	s.KmsKeyId = &v
	return s
}

// SetName sets the Name field's value.
func (s *CopyImageInput) SetName(v string) *CopyImageInput {
	s.Name = &v
	return s
}

// SetSourceImageId sets the SourceImageId field's value.
func (s *CopyImageInput) SetSourceImageId(v string) *CopyImageInput {
	s.SourceImageId = &v
	return s
}

// SetSourceRegion sets the SourceRegion field's value.
func (s *CopyImageInput) SetSourceRegion(v string) *CopyImageInput {
	s.SourceRegion = &v
	return s
}

// Contains the output of CopyImage.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CopyImageResult
type CopyImageOutput struct {
	_ struct{} `type:"structure"`

	// The ID of the new AMI.
	ImageId *string `locationName:"imageId" type:"string"`
}

// String returns the string representation
func (s CopyImageOutput) String() string {
	return awsutil.Prettify(s)
}

// SetDestinationCidrBlock sets the DestinationCidrBlock field's value.
func (s *CreateVpnConnectionRouteInput) SetDestinationCidrBlock(v string) *CreateVpnConnectionRouteInput {
	s.DestinationCidrBlock = &v
	return s
}

// SetVpnConnectionId sets the VpnConnectionId field's value.
func (s *CreateVpnConnectionRouteInput) SetVpnConnectionId(v string) *CreateVpnConnectionRouteInput {
	s.VpnConnectionId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateVpnConnectionRouteOutput
type CreateVpnConnectionRouteOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s CreateVpnConnectionRouteOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CopyImageOutput) GoString() string {
	return s.String()
}

// SetImageId sets the ImageId field's value.
func (s *CopyImageOutput) SetImageId(v string) *CopyImageOutput {
	s.ImageId = &v
	return s
}

// Contains the parameters for DescribeSnapshots.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeSnapshotsRequest
type DescribeSnapshotsInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more filters.
	//
	//    * description - A description of the snapshot.
	//
	//    * owner-alias - Value from an Amazon-maintained list (amazon | aws-marketplace
	//    | microsoft) of snapshot owners. Not to be confused with the user-configured
	//    AWS account alias, which is set from the IAM consolew.
	//
	//    * owner-id - The ID of the AWS account that owns the snapshot.
	//
	//    * progress - The progress of the snapshot, as a percentage (for example,
	//    80%).
	//
	//    * snapshot-id - The snapshot ID.
	//
	//    * start-time - The time stamp when the snapshot was initiated.
	//
	//    * status - The status of the snapshot (pending | completed | error).
	//
	//    * tag:key=value - The key/value combination of a tag assigned to the resource.
	//    Specify the key of the tag in the filter name and the value of the tag
	//    in the filter value. For example, for the tag Purpose=X, specify tag:Purpose
	//    for the filter name and X for the filter value.
	//
	//    * tag-key - The key of a tag assigned to the resource. This filter is
	//    independent of the tag-value filter. For example, if you use both the
	//    filter "tag-key=Purpose" and the filter "tag-value=X", you get any resources
	//    assigned both the tag key Purpose (regardless of what the tag's value
	//    is), and the tag value X (regardless of what the tag's key is). If you
	//    want to list only resources where Purpose is X, see the tag:key=value
	//    filter.
	//
	//    * tag-value - The value of a tag assigned to the resource. This filter
	//    is independent of the tag-key filter.
	//
	//    * volume-id - The ID of the volume the snapshot is for.
	//
	//    * volume-size - The size of the volume, in GiB.
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// The maximum number of snapshot results returned by DescribeSnapshots in paginated
	// output. When this parameter is used, DescribeSnapshots only returns MaxResults
	// results in a single page along with a NextToken response element. The remaining
	// results of the initial request can be seen by sending another DescribeSnapshots
	// request with the returned NextToken value. This value can be between 5 and
	// 1000; if MaxResults is given a value larger than 1000, only 1000 results
	// are returned. If this parameter is not used, then DescribeSnapshots returns
	// all results. You cannot specify this parameter and the snapshot IDs parameter
	// in the same request.
	MaxResults *int64 `type:"integer"`

	// The NextToken value returned from a previous paginated DescribeSnapshots
	// request where MaxResults was used and the results exceeded the value of that
	// parameter. Pagination continues from the end of the previous results that
	// returned the NextToken value. This value is null when there are no more results
	// to return.
	NextToken *string `type:"string"`

	// Returns the snapshots owned by the specified owner. Multiple owners can be
	// specified.
	OwnerIds []*string `locationName:"Owner" locationNameList:"Owner" type:"list"`

	// One or more AWS accounts IDs that can create volumes from the snapshot.
	RestorableByUserIds []*string `locationName:"RestorableBy" type:"list"`

	// One or more snapshot IDs.
	//
	// Default: Describes snapshots for which you have launch permissions.
	SnapshotIds []*string `locationName:"SnapshotId" locationNameList:"SnapshotId" type:"list"`
}

// String returns the string representation
func (s DescribeSnapshotsInput) String() string {
	return s.String()
}
func (s CreateVpnConnectionRouteOutput) GoString() string {
	return s.String()
}

type DeleteVpnConnectionRouteInput struct {
	_ struct{} `type:"structure"`

	// The CIDR block associated with the local subnet of the customer network.
	//
	// DestinationCidrBlock is a required field
	DestinationCidrBlock *string `type:"string" required:"true"`

	// The ID of the VPN connection.
	//
	// VpnConnectionId is a required field
	VpnConnectionId *string `type:"string" required:"true"`
}

// String returns the string representation
func (s DeleteVpnConnectionRouteInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeSnapshotsInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeSnapshotsInput) SetDryRun(v bool) *DescribeSnapshotsInput {
	s.DryRun = &v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribeSnapshotsInput) SetFilters(v []*Filter) *DescribeSnapshotsInput {
	s.Filters = v
	return s
}

// SetMaxResults sets the MaxResults field's value.
func (s *DescribeSnapshotsInput) SetMaxResults(v int64) *DescribeSnapshotsInput {
	s.MaxResults = &v
	return s
}

// SetNextToken sets the NextToken field's value.
func (s *DescribeSnapshotsInput) SetNextToken(v string) *DescribeSnapshotsInput {
	s.NextToken = &v
	return s
}

// SetOwnerIds sets the OwnerIds field's value.
func (s *DescribeSnapshotsInput) SetOwnerIds(v []*string) *DescribeSnapshotsInput {
	s.OwnerIds = v
	return s
}

// SetRestorableByUserIds sets the RestorableByUserIds field's value.
func (s *DescribeSnapshotsInput) SetRestorableByUserIds(v []*string) *DescribeSnapshotsInput {
	s.RestorableByUserIds = v
	return s
}

// SetSnapshotIds sets the SnapshotIds field's value.
func (s *DescribeSnapshotsInput) SetSnapshotIds(v []*string) *DescribeSnapshotsInput {
	s.SnapshotIds = v
	return s
}

// Contains the output of DescribeSnapshots.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeSnapshotsResult
type DescribeSnapshotsOutput struct {
	_ struct{} `type:"structure"`

	// The NextToken value to include in a future DescribeSnapshots request. When
	// the results of a DescribeSnapshots request exceed MaxResults, this value
	// can be used to retrieve the next page of results. This value is null when
	// there are no more results to return.
	NextToken *string `locationName:"nextToken" type:"string"`

	// Information about the snapshots.
	Snapshots []*Snapshot `locationName:"snapshotSet" locationNameList:"item" type:"list"`

	RequestId *string `locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s DescribeSnapshotsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeSnapshotsOutput) GoString() string {
	return s.String()
}

// SetNextToken sets the NextToken field's value.
func (s *DescribeSnapshotsOutput) SetNextToken(v string) *DescribeSnapshotsOutput {
	s.NextToken = &v
	return s
}

// SetSnapshots sets the Snapshots field's value.
func (s *DescribeSnapshotsOutput) SetSnapshots(v []*Snapshot) *DescribeSnapshotsOutput {
	s.Snapshots = v
	return s
}

// Describes a snapshot.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/Snapshot
type Snapshot struct {
	_ struct{} `type:"structure"`

	// The data encryption key identifier for the snapshot. This value is a unique
	// identifier that corresponds to the data encryption key that was used to encrypt
	// the original volume or snapshot copy. Because data encryption keys are inherited
	// by volumes created from snapshots, and vice versa, if snapshots share the
	// same data encryption key identifier, then they belong to the same volume/snapshot
	// lineage. This parameter is only returned by the DescribeSnapshots API operation.
	DataEncryptionKeyId *string `locationName:"dataEncryptionKeyId" type:"string"`

	// The description for the snapshot.
	Description *string `locationName:"description" type:"string"`

	// Indicates whether the snapshot is encrypted.
	Encrypted *bool `locationName:"encrypted" type:"boolean"`

	// The full ARN of the AWS Key Management Service (AWS KMS) customer master
	// key (CMK) that was used to protect the volume encryption key for the parent
	// volume.
	KmsKeyId *string `locationName:"kmsKeyId" type:"string"`

	// Value from an Amazon-maintained list (amazon | aws-marketplace | microsoft)
	// of snapshot owners. Not to be confused with the user-configured AWS account
	// alias, which is set from the IAM console.
	OwnerAlias *string `locationName:"ownerAlias" type:"string"`

	// The AWS account ID of the EBS snapshot owner.
	OwnerId *string `locationName:"ownerId" type:"string"`

	// The progress of the snapshot, as a percentage.
	Progress *string `locationName:"progress" type:"string"`

	// The ID of the snapshot. Each snapshot receives a unique identifier when it
	// is created.
	SnapshotId *string `locationName:"snapshotId" type:"string"`

	// The time stamp when the snapshot was initiated.
	StartTime *time.Time `locationName:"startTime" type:"timestamp" timestampFormat:"iso8601"`

	// The snapshot state.
	State *string `locationName:"status" type:"string" enum:"SnapshotState"`

	// Encrypted Amazon EBS snapshots are copied asynchronously. If a snapshot copy
	// operation fails (for example, if the proper AWS Key Management Service (AWS
	// KMS) permissions are not obtained) this field displays error state details
	// to help you diagnose why the error occurred. This parameter is only returned
	// by the DescribeSnapshots API operation.
	StateMessage *string `locationName:"statusMessage" type:"string"`

	// Any tags assigned to the snapshot.
	Tags []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`

	// The ID of the volume that was used to create the snapshot. Snapshots created
	// by the CopySnapshot action have an arbitrary volume ID that should not be
	// used for any purpose.
	VolumeId *string `locationName:"volumeId" type:"string"`

	// The size of the volume, in GiB.
	VolumeSize *int64 `locationName:"volumeSize" type:"integer"`
}

// String returns the string representation
func (s Snapshot) String() string {
	return awsutil.Prettify(s)
}
func (s DeleteVpnConnectionRouteInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteVpnConnectionRouteInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeleteVpnConnectionRouteInput"}
	if s.DestinationCidrBlock == nil {
		invalidParams.Add(request.NewErrParamRequired("DestinationCidrBlock"))
	}
	if s.VpnConnectionId == nil {
		invalidParams.Add(request.NewErrParamRequired("VpnConnectionId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDestinationCidrBlock sets the DestinationCidrBlock field's value.
func (s *DeleteVpnConnectionRouteInput) SetDestinationCidrBlock(v string) *DeleteVpnConnectionRouteInput {
	s.DestinationCidrBlock = &v
	return s
}

// SetVpnConnectionId sets the VpnConnectionId field's value.
func (s *DeleteVpnConnectionRouteInput) SetVpnConnectionId(v string) *DeleteVpnConnectionRouteInput {
	s.VpnConnectionId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteVpnConnectionRouteOutput
type DeleteVpnConnectionRouteOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeleteVpnConnectionRouteOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Snapshot) GoString() string {
	return s.String()
}

// SetDataEncryptionKeyId sets the DataEncryptionKeyId field's value.
func (s *Snapshot) SetDataEncryptionKeyId(v string) *Snapshot {
	s.DataEncryptionKeyId = &v
	return s
}

// SetDescription sets the Description field's value.
func (s *Snapshot) SetDescription(v string) *Snapshot {
	s.Description = &v
	return s
}

// SetEncrypted sets the Encrypted field's value.
func (s *Snapshot) SetEncrypted(v bool) *Snapshot {
	s.Encrypted = &v
	return s
}

// SetKmsKeyId sets the KmsKeyId field's value.
func (s *Snapshot) SetKmsKeyId(v string) *Snapshot {
	s.KmsKeyId = &v
	return s
}

// SetOwnerAlias sets the OwnerAlias field's value.
func (s *Snapshot) SetOwnerAlias(v string) *Snapshot {
	s.OwnerAlias = &v
	return s
}

// SetOwnerId sets the OwnerId field's value.
func (s *Snapshot) SetOwnerId(v string) *Snapshot {
	s.OwnerId = &v
	return s
}

// SetProgress sets the Progress field's value.
func (s *Snapshot) SetProgress(v string) *Snapshot {
	s.Progress = &v
	return s
}

// SetSnapshotId sets the SnapshotId field's value.
func (s *Snapshot) SetSnapshotId(v string) *Snapshot {
	s.SnapshotId = &v
	return s
}

// SetStartTime sets the StartTime field's value.
func (s *Snapshot) SetStartTime(v time.Time) *Snapshot {
	s.StartTime = &v
	return s
}

// SetState sets the State field's value.
func (s *Snapshot) SetState(v string) *Snapshot {
	s.State = &v
	return s
}

// SetStateMessage sets the StateMessage field's value.
func (s *Snapshot) SetStateMessage(v string) *Snapshot {
	s.StateMessage = &v
	return s
}

// SetTags sets the Tags field's value.
func (s *Snapshot) SetTags(v []*Tag) *Snapshot {
	s.Tags = v
	return s
}

// SetVolumeId sets the VolumeId field's value.
func (s *Snapshot) SetVolumeId(v string) *Snapshot {
	s.VolumeId = &v
	return s
}

// SetVolumeSize sets the VolumeSize field's value.
func (s *Snapshot) SetVolumeSize(v int64) *Snapshot {
	s.VolumeSize = &v
	return s
}
func (s DeleteVpnConnectionRouteOutput) GoString() string {
	return s.String()
}

type DescribeAvailabilityZonesInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more filters.
	//
	//    * message - Information about the Availability Zone.
	//
	//    * region-name - The name of the region for the Availability Zone (for
	//    example, us-east-1).
	//
	//    * state - The state of the Availability Zone (available | information
	//    | impaired | unavailable).
	//
	//    * zone-name - The name of the Availability Zone (for example, us-east-1a).
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// The names of one or more Availability Zones.
	ZoneNames []*string `locationName:"ZoneName" locationNameList:"ZoneName" type:"list"`
}

// String returns the string representation
func (s DescribeAvailabilityZonesInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeAvailabilityZonesInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeAvailabilityZonesInput) SetDryRun(v bool) *DescribeAvailabilityZonesInput {
	s.DryRun = &v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribeAvailabilityZonesInput) SetFilters(v []*Filter) *DescribeAvailabilityZonesInput {
	s.Filters = v
	return s
}

// SetZoneNames sets the ZoneNames field's value.
func (s *DescribeAvailabilityZonesInput) SetZoneNames(v []*string) *DescribeAvailabilityZonesInput {
	s.ZoneNames = v
	return s
}

// Contains the output of DescribeAvailabiltyZones.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeAvailabilityZonesResult
type DescribeAvailabilityZonesOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more Availability Zones.
	AvailabilityZones []*AvailabilityZone `locationName:"availabilityZoneInfo" locationNameList:"item" type:"list"`

	RequestId *string `locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s DescribeAvailabilityZonesOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeAvailabilityZonesOutput) GoString() string {
	return s.String()
}

// SetAvailabilityZones sets the AvailabilityZones field's value.
func (s *DescribeAvailabilityZonesOutput) SetAvailabilityZones(v []*AvailabilityZone) *DescribeAvailabilityZonesOutput {
	s.AvailabilityZones = v
	return s
}

type AvailabilityZone struct {
	_ struct{} `type:"structure"`

	// Any messages about the Availability Zone.
	Messages []*AvailabilityZoneMessage `locationName:"messageSet" locationNameList:"item" type:"list"`

	// The name of the region.
	RegionName *string `locationName:"regionName" type:"string"`

	// The state of the Availability Zone.
	State *string `locationName:"zoneState" type:"string" enum:"AvailabilityZoneState"`

	// The name of the Availability Zone.
	ZoneName *string `locationName:"zoneName" type:"string"`
}

// String returns the string representation
func (s AvailabilityZone) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AvailabilityZone) GoString() string {
	return s.String()
}

// SetMessages sets the Messages field's value.
func (s *AvailabilityZone) SetMessages(v []*AvailabilityZoneMessage) *AvailabilityZone {
	s.Messages = v
	return s
}

// SetRegionName sets the RegionName field's value.
func (s *AvailabilityZone) SetRegionName(v string) *AvailabilityZone {
	s.RegionName = &v
	return s
}

// SetState sets the State field's value.
func (s *AvailabilityZone) SetState(v string) *AvailabilityZone {
	s.State = &v
	return s
}

// SetZoneName sets the ZoneName field's value.
func (s *AvailabilityZone) SetZoneName(v string) *AvailabilityZone {
	s.ZoneName = &v
	return s
}

type AvailabilityZoneMessage struct {
	_ struct{} `type:"structure"`

	// The message about the Availability Zone.
	Message *string `locationName:"message" type:"string"`
}

// String returns the string representation
func (s AvailabilityZoneMessage) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AvailabilityZoneMessage) GoString() string {
	return s.String()
}

// SetMessage sets the Message field's value.
func (s *AvailabilityZoneMessage) SetMessage(v string) *AvailabilityZoneMessage {
	s.Message = &v
	return s
}

type DescribePrefixListsInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `type:"boolean"`

	// One or more filters.
	//
	//    * prefix-list-id: The ID of a prefix list.
	//
	//    * prefix-list-name: The name of a prefix list.
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// The maximum number of items to return for this request. The request returns
	// a token that you can specify in a subsequent call to get the next set of
	// results.
	//
	// Constraint: If the value specified is greater than 1000, we return only 1000
	// items.
	MaxResults *int64 `type:"integer"`

	// The token for the next set of items to return. (You received this token from
	// a prior call.)
	NextToken *string `type:"string"`

	// One or more prefix list IDs.
	PrefixListIds []*string `locationName:"PrefixListId" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s DescribePrefixListsInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribePrefixListsInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *DescribePrefixListsInput) SetDryRun(v bool) *DescribePrefixListsInput {
	s.DryRun = &v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribePrefixListsInput) SetFilters(v []*Filter) *DescribePrefixListsInput {
	s.Filters = v
	return s
}

// SetMaxResults sets the MaxResults field's value.
func (s *DescribePrefixListsInput) SetMaxResults(v int64) *DescribePrefixListsInput {
	s.MaxResults = &v
	return s
}

// SetNextToken sets the NextToken field's value.
func (s *DescribePrefixListsInput) SetNextToken(v string) *DescribePrefixListsInput {
	s.NextToken = &v
	return s
}

// SetPrefixListIds sets the PrefixListIds field's value.
func (s *DescribePrefixListsInput) SetPrefixListIds(v []*string) *DescribePrefixListsInput {
	s.PrefixListIds = v
	return s
}

// Contains the output of DescribePrefixLists.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribePrefixListsResult
type DescribePrefixListsOutput struct {
	_ struct{} `type:"structure"`

	// The token to use when requesting the next set of items. If there are no additional
	// items to return, the string is empty.
	NextToken *string `locationName:"nextToken" type:"string"`

	// All available prefix lists.
	PrefixLists []*PrefixList `locationName:"prefixListSet" locationNameList:"item" type:"list"`

	RequestId *string `locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s DescribePrefixListsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribePrefixListsOutput) GoString() string {
	return s.String()
}

// SetNextToken sets the NextToken field's value.
func (s *DescribePrefixListsOutput) SetNextToken(v string) *DescribePrefixListsOutput {
	s.NextToken = &v
	return s
}

// SetPrefixLists sets the PrefixLists field's value.
func (s *DescribePrefixListsOutput) SetPrefixLists(v []*PrefixList) *DescribePrefixListsOutput {
	s.PrefixLists = v
	return s
}

type PrefixList struct {
	_ struct{} `type:"structure"`

	// The IP address range of the AWS service.
	Cidrs []*string `locationName:"cidrSet" locationNameList:"item" type:"list"`

	// The ID of the prefix.
	PrefixListId *string `locationName:"prefixListId" type:"string"`

	// The name of the prefix.
	PrefixListName *string `locationName:"prefixListName" type:"string"`
}

// String returns the string representation
func (s PrefixList) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s PrefixList) GoString() string {
	return s.String()
}

// SetCidrs sets the Cidrs field's value.
func (s *PrefixList) SetCidrs(v []*string) *PrefixList {
	s.Cidrs = v
	return s
}

// SetPrefixListId sets the PrefixListId field's value.
func (s *PrefixList) SetPrefixListId(v string) *PrefixList {
	s.PrefixListId = &v
	return s
}

// SetPrefixListName sets the PrefixListName field's value.
func (s *PrefixList) SetPrefixListName(v string) *PrefixList {
	s.PrefixListName = &v
	return s
}

type DescribeQuotasInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `type:"boolean"`

	// One or more filters.
	//
	//    * prefix-list-id: The ID of a prefix list.
	//
	//    * prefix-list-name: The name of a prefix list.
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// The maximum number of items to return for this request. The request returns
	// a token that you can specify in a subsequent call to get the next set of
	// results.
	//
	// Constraint: If the value specified is greater than 1000, we return only 1000
	// items.
	MaxResults *int64 `type:"integer"`

	// The token for the next set of items to return. (You received this token from
	// a prior call.)
	NextToken *string `type:"string"`

	// One or more prefix list IDs.
	QuotaName []*string `locationName:"QuotaName" locationNameList:"item" type:"list"`
}

type DescribeQuotasOutput struct {
	_ struct{} `type:"structure"`

	// The token to use when requesting the next set of items. If there are no additional
	// items to return, the string is empty.
	NextToken *string `locationName:"nextToken" type:"string"`

	// All available prefix lists.
	ReferenceQuotaSet []*ReferenceQuota `locationName:"referenceQuotaSet" locationNameList:"item" type:"list"`
	RequestId         *string           `locationName:"requestId" type:"string"`
}

type DescribeRegionsInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more filters.
	//
	//    * endpoint - The endpoint of the region (for example, ec2.us-east-1.amazonaws.com).
	//
	//    * region-name - The name of the region (for example, us-east-1).
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// The names of one or more regions.
	RegionNames []*string `locationName:"RegionName" locationNameList:"RegionName" type:"list"`
}

// String returns the string representation
func (s DescribeRegionsInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeRegionsInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeRegionsInput) SetDryRun(v bool) *DescribeRegionsInput {
	s.DryRun = &v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribeRegionsInput) SetFilters(v []*Filter) *DescribeRegionsInput {
	s.Filters = v
	return s
}

// SetRegionNames sets the RegionNames field's value.
func (s *DescribeRegionsInput) SetRegionNames(v []*string) *DescribeRegionsInput {
	s.RegionNames = v
	return s
}

// Contains the output of DescribeRegions.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeRegionsResult
type DescribeRegionsOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more regions.
	Regions   []*Region `locationName:"regionInfo" locationNameList:"item" type:"list"`
	RequestId *string   `locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s DescribeRegionsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeRegionsOutput) GoString() string {
	return s.String()
}

// SetRegions sets the Regions field's value.
func (s *DescribeRegionsOutput) SetRegions(v []*Region) *DescribeRegionsOutput {
	s.Regions = v
	return s
}

type ReferenceQuota struct {
	QuotaSet  []*QuotaSet `locationName:"quotaSet" locationNameList:"item" type:"list"`
	Reference *string     `locationName:"reference" type:"string"`
}

type QuotaSet struct {
	Description    *string `locationName:"description" type:"string"`
	DisplayName    *string `locationName:"displayName" type:"string"`
	GroupName      *string `locationName:"groupName" type:"string"`
	MaxQuotaValue  *string `locationName:"maxQuotaValue" type:"string"`
	Name           *string `locationName:"name" type:"string"`
	OwnerId        *string `locationName:"ownerId" type:"string"`
	UsedQuotaValue *string `locationName:"usedQuotaValue" type:"string"`
}

type Region struct {
	_ struct{} `type:"structure"`

	// The region service endpoint.
	Endpoint *string `locationName:"regionEndpoint" type:"string"`

	// The name of the region.
	RegionName *string `locationName:"regionName" type:"string"`
}

// String returns the string representation
func (s Region) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Region) GoString() string {
	return s.String()
}

// SetEndpoint sets the Endpoint field's value.
func (s *Region) SetEndpoint(v string) *Region {
	s.Endpoint = &v
	return s
}

// SetRegionName sets the RegionName field's value.
func (s *Region) SetRegionName(v string) *Region {
	s.RegionName = &v
	return s
}

type CreateSnapshotInput struct {
	_ struct{} `type:"structure"`

	// A description for the snapshot.
	Description *string `type:"string"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the EBS volume.
	//
	// VolumeId is a required field
	VolumeId *string `type:"string" required:"true"`
}

// String returns the string representation
func (s CreateSnapshotInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateSnapshotInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CreateSnapshotInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "CreateSnapshotInput"}
	if s.VolumeId == nil {
		invalidParams.Add(request.NewErrParamRequired("VolumeId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDescription sets the Description field's value.
func (s *CreateSnapshotInput) SetDescription(v string) *CreateSnapshotInput {
	s.Description = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *CreateSnapshotInput) SetDryRun(v bool) *CreateSnapshotInput {
	s.DryRun = &v
	return s
}

// SetVolumeId sets the VolumeId field's value.
func (s *CreateSnapshotInput) SetVolumeId(v string) *CreateSnapshotInput {
	s.VolumeId = &v
	return s
}

// Describes the snapshot created from the imported disk.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/SnapshotDetail
type SnapshotDetail struct {
	_ struct{} `type:"structure"`

	// A description for the snapshot.
	Description *string `locationName:"description" type:"string"`

	// The block device mapping for the snapshot.
	DeviceName *string `locationName:"deviceName" type:"string"`

	// The size of the disk in the snapshot, in GiB.
	DiskImageSize *float64 `locationName:"diskImageSize" type:"double"`

	// The format of the disk image from which the snapshot is created.
	Format *string `locationName:"format" type:"string"`

	// The percentage of progress for the task.
	Progress *string `locationName:"progress" type:"string"`

	// The snapshot ID of the disk being imported.
	SnapshotId *string `locationName:"snapshotId" type:"string"`

	// A brief status of the snapshot creation.
	Status *string `locationName:"status" type:"string"`

	// A detailed status message for the snapshot creation.
	StatusMessage *string `locationName:"statusMessage" type:"string"`

	// The URL used to access the disk image.
	Url *string `locationName:"url" type:"string"`

	// The S3 bucket for the disk image.
	UserBucket *UserBucketDetails `locationName:"userBucket" type:"structure"`
}

// String returns the string representation
func (s SnapshotDetail) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s SnapshotDetail) GoString() string {
	return s.String()
}

// SetDescription sets the Description field's value.
func (s *SnapshotDetail) SetDescription(v string) *SnapshotDetail {
	s.Description = &v
	return s
}

// SetDeviceName sets the DeviceName field's value.
func (s *SnapshotDetail) SetDeviceName(v string) *SnapshotDetail {
	s.DeviceName = &v
	return s
}

// SetDiskImageSize sets the DiskImageSize field's value.
func (s *SnapshotDetail) SetDiskImageSize(v float64) *SnapshotDetail {
	s.DiskImageSize = &v
	return s
}

// SetFormat sets the Format field's value.
func (s *SnapshotDetail) SetFormat(v string) *SnapshotDetail {
	s.Format = &v
	return s
}

// SetProgress sets the Progress field's value.
func (s *SnapshotDetail) SetProgress(v string) *SnapshotDetail {
	s.Progress = &v
	return s
}

// SetSnapshotId sets the SnapshotId field's value.
func (s *SnapshotDetail) SetSnapshotId(v string) *SnapshotDetail {
	s.SnapshotId = &v
	return s
}

// SetStatus sets the Status field's value.
func (s *SnapshotDetail) SetStatus(v string) *SnapshotDetail {
	s.Status = &v
	return s
}

// SetStatusMessage sets the StatusMessage field's value.
func (s *SnapshotDetail) SetStatusMessage(v string) *SnapshotDetail {
	s.StatusMessage = &v
	return s
}

// SetUrl sets the Url field's value.
func (s *SnapshotDetail) SetUrl(v string) *SnapshotDetail {
	s.Url = &v
	return s
}

// SetUserBucket sets the UserBucket field's value.
func (s *SnapshotDetail) SetUserBucket(v *UserBucketDetails) *SnapshotDetail {
	s.UserBucket = v
	return s
}

// The disk container object for the import snapshot request.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/SnapshotDiskContainer
type SnapshotDiskContainer struct {
	_ struct{} `type:"structure"`

	// The description of the disk image being imported.
	Description *string `type:"string"`

	// The format of the disk image being imported.
	//
	// Valid values: RAW | VHD | VMDK | OVA
	Format *string `type:"string"`

	// The URL to the Amazon S3-based disk image being imported. It can either be
	// a https URL (https://..) or an Amazon S3 URL (s3://..).
	Url *string `type:"string"`

	// The S3 bucket for the disk image.
	UserBucket *UserBucket `type:"structure"`
}

// String returns the string representation
func (s SnapshotDiskContainer) String() string {
	return awsutil.Prettify(s)
}

type DescribeProductTypesInput struct {
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`
}

type DescribeProductTypesOutput struct {
	ProductTypeSet []*ProductType `locationName:"productTypeSet" locationNameList:"item" type:"list"`
	RequestId      *string        `locationName:"requestId" type:"string"`
}

type ProductType struct {
	Description   *string `locationName:"description" type:"string"`
	ProductTypeId *string `locationName:"productTypeId" type:"string"`
	Vendor        *string `locationName:"vendor" type:"string"`
}

type DescribeReservedInstancesInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// Describes whether the Reserved Instance is Standard or Convertible.
	AvailabilityZone *string `type:"string" enum:"AvailabilityZone"`

	OfferingClass *string `type:"string" enum:"OfferingClassType"`

	// The Reserved Instance offering type. If you are using tools that predate
	// the 2011-11-01 API version, you only have access to the Medium Utilization
	// Reserved Instance offering type.
	OfferingType *string `locationName:"offeringType" type:"string" enum:"OfferingTypeValues"`

	// One or more Reserved Instance IDs.
	//
	// Default: Describes all your Reserved Instances, or only those otherwise specified.
	ReservedInstancesIds []*string `locationName:"ReservedInstancesId" locationNameList:"ReservedInstancesId" type:"list"`
}

// String returns the string representation
func (s DescribeReservedInstancesInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s SnapshotDiskContainer) GoString() string {
	return s.String()
}

// SetDescription sets the Description field's value.
func (s *SnapshotDiskContainer) SetDescription(v string) *SnapshotDiskContainer {
	s.Description = &v
	return s
}

// SetFormat sets the Format field's value.
func (s *SnapshotDiskContainer) SetFormat(v string) *SnapshotDiskContainer {
	s.Format = &v
	return s
}

// SetUrl sets the Url field's value.
func (s *SnapshotDiskContainer) SetUrl(v string) *SnapshotDiskContainer {
	s.Url = &v
	return s
}

// SetUserBucket sets the UserBucket field's value.
func (s *SnapshotDiskContainer) SetUserBucket(v *UserBucket) *SnapshotDiskContainer {
	s.UserBucket = v
	return s
}

// Details about the import snapshot task.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/SnapshotTaskDetail
type SnapshotTaskDetail struct {
	_ struct{} `type:"structure"`

	// The description of the snapshot.
	Description *string `locationName:"description" type:"string"`

	// The size of the disk in the snapshot, in GiB.
	DiskImageSize *float64 `locationName:"diskImageSize" type:"double"`

	// The format of the disk image from which the snapshot is created.
	Format *string `locationName:"format" type:"string"`

	// The percentage of completion for the import snapshot task.
	Progress *string `locationName:"progress" type:"string"`

	// The snapshot ID of the disk being imported.
	SnapshotId *string `locationName:"snapshotId" type:"string"`

	// A brief status for the import snapshot task.
	Status *string `locationName:"status" type:"string"`

	// A detailed status message for the import snapshot task.
	StatusMessage *string `locationName:"statusMessage" type:"string"`

	// The URL of the disk image from which the snapshot is created.
	Url *string `locationName:"url" type:"string"`

	// The S3 bucket for the disk image.
	UserBucket *UserBucketDetails `locationName:"userBucket" type:"structure"`
}

// String returns the string representation
func (s SnapshotTaskDetail) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s SnapshotTaskDetail) GoString() string {
	return s.String()
}

// SetDescription sets the Description field's value.
func (s *SnapshotTaskDetail) SetDescription(v string) *SnapshotTaskDetail {
	s.Description = &v
	return s
}

// SetDiskImageSize sets the DiskImageSize field's value.
func (s *SnapshotTaskDetail) SetDiskImageSize(v float64) *SnapshotTaskDetail {
	s.DiskImageSize = &v
	return s
}

// SetFormat sets the Format field's value.
func (s *SnapshotTaskDetail) SetFormat(v string) *SnapshotTaskDetail {
	s.Format = &v
	return s
}

// SetProgress sets the Progress field's value.
func (s *SnapshotTaskDetail) SetProgress(v string) *SnapshotTaskDetail {
	s.Progress = &v
	return s
}

// SetSnapshotId sets the SnapshotId field's value.
func (s *SnapshotTaskDetail) SetSnapshotId(v string) *SnapshotTaskDetail {
	s.SnapshotId = &v
	return s
}

// SetStatus sets the Status field's value.
func (s *SnapshotTaskDetail) SetStatus(v string) *SnapshotTaskDetail {
	s.Status = &v
	return s
}

// SetStatusMessage sets the StatusMessage field's value.
func (s *SnapshotTaskDetail) SetStatusMessage(v string) *SnapshotTaskDetail {
	s.StatusMessage = &v
	return s
}

// SetUrl sets the Url field's value.
func (s *SnapshotTaskDetail) SetUrl(v string) *SnapshotTaskDetail {
	s.Url = &v
	return s
}

// SetUserBucket sets the UserBucket field's value.
func (s *SnapshotTaskDetail) SetUserBucket(v *UserBucketDetails) *SnapshotTaskDetail {
	s.UserBucket = v
	return s
}

type UserBucketDetails struct {
	_ struct{} `type:"structure"`

	// The S3 bucket from which the disk image was created.
	S3Bucket *string `locationName:"s3Bucket" type:"string"`

	// The file name of the disk image.
	S3Key *string `locationName:"s3Key" type:"string"`
}

// String returns the string representation
func (s UserBucketDetails) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s UserBucketDetails) GoString() string {
	return s.String()
}

// SetS3Bucket sets the S3Bucket field's value.
func (s *UserBucketDetails) SetS3Bucket(v string) *UserBucketDetails {
	s.S3Bucket = &v
	return s
}

// SetS3Key sets the S3Key field's value.
func (s *UserBucketDetails) SetS3Key(v string) *UserBucketDetails {
	s.S3Key = &v
	return s
}

type UserBucket struct {
	_ struct{} `type:"structure"`

	// The name of the S3 bucket where the disk image is located.
	S3Bucket *string `type:"string"`

	// The file name of the disk image.
	S3Key *string `type:"string"`
}

// String returns the string representation
func (s UserBucket) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s UserBucket) GoString() string {
	return s.String()
}

// SetS3Bucket sets the S3Bucket field's value.
func (s *UserBucket) SetS3Bucket(v string) *UserBucket {
	s.S3Bucket = &v
	return s
}

// SetS3Key sets the S3Key field's value.
func (s *UserBucket) SetS3Key(v string) *UserBucket {
	s.S3Key = &v
	return s
}

type DeleteSnapshotInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the EBS snapshot.
	//
	// SnapshotId is a required field
	SnapshotId *string `type:"string" required:"true"`
}

// String returns the string representation
func (s DeleteSnapshotInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteSnapshotInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteSnapshotInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeleteSnapshotInput"}
	if s.SnapshotId == nil {
		invalidParams.Add(request.NewErrParamRequired("SnapshotId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDryRun sets the DryRun field's value.
func (s *DeleteSnapshotInput) SetDryRun(v bool) *DeleteSnapshotInput {
	s.DryRun = &v
	return s
}

// SetSnapshotId sets the SnapshotId field's value.
func (s *DeleteSnapshotInput) SetSnapshotId(v string) *DeleteSnapshotInput {
	s.SnapshotId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteSnapshotOutput
type DeleteSnapshotOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeleteSnapshotOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteSnapshotOutput) GoString() string {
	return s.String()
}

func (s DescribeReservedInstancesInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeReservedInstancesInput) SetDryRun(v bool) *DescribeReservedInstancesInput {
	s.DryRun = &v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribeReservedInstancesInput) SetFilters(v []*Filter) *DescribeReservedInstancesInput {
	s.Filters = v
	return s
}

// SetOfferingClass sets the OfferingClass field's value.
func (s *DescribeReservedInstancesInput) SetOfferingClass(v string) *DescribeReservedInstancesInput {
	s.OfferingClass = &v
	return s
}

// SetOfferingType sets the OfferingType field's value.
func (s *DescribeReservedInstancesInput) SetOfferingType(v string) *DescribeReservedInstancesInput {
	s.OfferingType = &v
	return s
}

// SetReservedInstancesIds sets the ReservedInstancesIds field's value.
func (s *DescribeReservedInstancesInput) SetReservedInstancesIds(v []*string) *DescribeReservedInstancesInput {
	s.ReservedInstancesIds = v
	return s
}

type DescribeReservedInstancesOutput struct {
	_ struct{} `type:"structure"`

	// A list of Reserved Instances.
	ReservedInstances []*ReservedInstances `locationName:"reservedInstancesSet" locationNameList:"item" type:"list"`
	RequestId         *string              `locationName:"requestId" type:"string"`
}
type ReservedInstances struct {
	_ struct{} `type:"structure"`

	// The Availability Zone in which the Reserved Instance can be used.
	AvailabilityZone *string `locationName:"availabilityZone" type:"string"`

	// The currency of the Reserved Instance. It's specified using ISO 4217 standard
	// currency codes. At this time, the only supported currency is USD.
	CurrencyCode *string `locationName:"currencyCode" type:"string" enum:"CurrencyCodeValues"`

	// The duration of the Reserved Instance, in seconds.
	Duration *int64 `locationName:"duration" type:"long"`

	// The time when the Reserved Instance expires.
	End *time.Time `locationName:"end" type:"timestamp" timestampFormat:"iso8601"`

	// The purchase price of the Reserved Instance.
	FixedPrice *float64 `locationName:"fixedPrice" type:"float"`

	// The number of reservations purchased.
	InstanceCount *int64 `locationName:"instanceCount" type:"integer"`

	// The tenancy of the instance.
	InstanceTenancy *string `locationName:"instanceTenancy" type:"string" enum:"Tenancy"`

	// The instance type on which the Reserved Instance can be used.
	InstanceType *string `locationName:"instanceType" type:"string" enum:"InstanceType"`

	// The offering class of the Reserved Instance.
	OfferingClass *string `locationName:"offeringClass" type:"string" enum:"OfferingClassType"`

	// The Reserved Instance offering type.
	OfferingType *string `locationName:"offeringType" type:"string" enum:"OfferingTypeValues"`

	// The Reserved Instance product platform description.
	ProductDescription *string `locationName:"productDescription" type:"string" enum:"RIProductDescription"`

	// The recurring charge tag assigned to the resource.
	RecurringCharges []*RecurringCharge `locationName:"recurringCharges" locationNameList:"item" type:"list"`

	// The ID of the Reserved Instance.
	ReservedInstancesId *string `locationName:"reservedInstancesId" type:"string"`

	// The scope of the Reserved Instance.
	Scope *string `locationName:"scope" type:"string" enum:"scope"`

	// The date and time the Reserved Instance started.
	Start *time.Time `locationName:"start" type:"timestamp" timestampFormat:"iso8601"`

	// The state of the Reserved Instance purchase.
	State *string `locationName:"state" type:"string" enum:"ReservedInstanceState"`

	// Any tags assigned to the resource.
	Tags []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`

	// The usage price of the Reserved Instance, per hour.
	UsagePrice *float64 `locationName:"usagePrice" type:"float"`
}

type RecurringCharge struct {
	_ struct{} `type:"structure"`

	// The amount of the recurring charge.
	Amount *float64 `locationName:"amount" type:"double"`

	// The frequency of the recurring charge.
	Frequency *string `locationName:"frequency" type:"string" enum:"RecurringChargeFrequency"`
}

type DescribeInstanceTypesInput struct {
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`
}

type DescribeInstanceTypesOutput struct {
	InstanceTypeSet []*InstanceType `locationName:"instanceTypeSet" locationNameList:"item" type:"list"`
	RequestId       *string         `locationName:"requestId" type:"string"`
}

type InstanceType struct {
	EbsOptimizedAvailable *bool   `locationName:"ebsOptimizedAvailable" type:"bool"`
	MaxIpAddresses        *int64  `locationName:"maxIpAddresses" type:"int64"`
	Memory                *int64  `locationName:"memory" type:"int64"`
	Name                  *string `locationName:"name" type:"string"`
	StorageCount          *int64  `locationName:"storageCount" type:"int64"`
	StorageSize           *int64  `locationName:"storageSize" type:"int64"`
	Vcpu                  *int64  `locationName:"vcpu" type:"int64"`
}

type DescribeReservedInstancesOfferingsInput struct {
	AvailabilityZone             *string   `locationName:"availabilityZone" type:"string"`
	Filters                      []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`
	InstanceTenancy              *string   `locationName:"instanceTenancy" type:"string" enum:"Tenancy"`
	InstanceType                 *string   `locationName:"instanceType" type:"string" enum:"InstanceType"`
	OfferingType                 *string   `locationName:"offeringType" type:"string" enum:"OfferingTypeValues"`
	ProductDescription           *string   `locationName:"productDescription" type:"string" enum:"RIProductDescription"`
	ReservedInstancesOfferingIds []*string `locationName:"reservedInstancesOfferingId" type:"string"`
}

type DescribeReservedInstancesOfferingsOutput struct {
	ReservedInstancesOfferingsSet []*ReservedInstancesOffering `locationName:"reservedInstancesOfferingsSet" locationNameList:"item" type:"list"`
	RequestId                     *string                      `locationName:"requestId" type:"string"`
}

type ReservedInstancesOffering struct {
	AvailabilityZone            *string            `locationName:"availabilityZone" type:"string"`
	CurrencyCode                *string            `locationName:"currencyCode" type:"string"`
	Duration                    *string            `locationName:"duration" type:"string"`
	FixedPrice                  *int64             `locationName:"fixedPrice" type:"int64"`
	InstanceTenancy             *string            `locationName:"instanceTenancy" type:"string" enum:"Tenancy"`
	InstanceType                *string            `locationName:"instanceType" type:"string" enum:"InstanceType"`
	Martketplace                *bool              `locationName:"martketplace" type:"bool"`
	OfferingType                *string            `locationName:"offeringType" type:"string" enum:"OfferingTypeValues"`
	ProductDescription          *string            `locationName:"productDescription" type:"string" enum:"RIProductDescription"`
	PricingDetailsSet           []*PricingDetail   `locationName:"pricingDetail" locationNameList:"item" type:"list"`
	RecurringCharges            []*RecurringCharge `locationName:"recurringCharges" locationNameList:"item" type:"list"`
	ReservedInstancesOfferingId *string            `locationName:"reservedInstancesOfferingId" type:"string"`
	UsagePrice                  *int64             `locationName:"usagePrice" type:"int64"`
}

type PricingDetail struct {
	Count *int64 `locationName:"count" type:"int64"`
}

type DescribeImageAttributeInput struct {
	_ struct{} `type:"structure"`

	// The AMI attribute.
	//
	// Note: Depending on your account privileges, the blockDeviceMapping attribute
	// may return a Client.AuthFailure error. If this happens, use DescribeImages
	// to get information about the block device mapping for the AMI.
	//
	// Attribute is a required field
	Attribute *string `type:"string" required:"true" enum:"ImageAttributeName"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the AMI.
	//
	// ImageId is a required field
	ImageId *string `type:"string" required:"true"`
}

// String returns the string representation
func (s DescribeImageAttributeInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeImageAttributeInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DescribeImageAttributeInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DescribeImageAttributeInput"}
	if s.Attribute == nil {
		invalidParams.Add(request.NewErrParamRequired("Attribute"))
	}
	if s.ImageId == nil {
		invalidParams.Add(request.NewErrParamRequired("ImageId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAttribute sets the Attribute field's value.
func (s *DescribeImageAttributeInput) SetAttribute(v string) *DescribeImageAttributeInput {
	s.Attribute = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeImageAttributeInput) SetDryRun(v bool) *DescribeImageAttributeInput {
	s.DryRun = &v
	return s
}

// SetImageId sets the ImageId field's value.
func (s *DescribeImageAttributeInput) SetImageId(v string) *DescribeImageAttributeInput {
	s.ImageId = &v
	return s
}

// Describes an image attribute.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/ImageAttribute
type DescribeImageAttributeOutput struct {
	_ struct{} `type:"structure"`

	// One or more block device mapping entries.
	BlockDeviceMappings []*BlockDeviceMapping `locationName:"blockDeviceMapping" locationNameList:"item" type:"list"`

	// A description for the AMI.
	Description *AttributeValue `locationName:"description" type:"structure"`

	// The ID of the AMI.
	ImageId *string `locationName:"imageId" type:"string"`

	// The kernel ID.
	KernelId *AttributeValue `locationName:"kernel" type:"structure"`

	// One or more launch permissions.
	LaunchPermissions []*LaunchPermission `locationName:"launchPermission" locationNameList:"item" type:"list"`

	// One or more product codes.
	ProductCodes []*ProductCode `locationName:"productCodes" locationNameList:"item" type:"list"`

	// The RAM disk ID.
	RamdiskId *AttributeValue `locationName:"ramdisk" type:"structure"`

	// Indicates whether enhanced networking with the Intel 82599 Virtual Function
	// interface is enabled.
	SriovNetSupport *AttributeValue `locationName:"sriovNetSupport" type:"structure"`
	RequestId       *string         `locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s DescribeImageAttributeOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeImageAttributeOutput) GoString() string {
	return s.String()
}

// SetBlockDeviceMappings sets the BlockDeviceMappings field's value.
func (s *DescribeImageAttributeOutput) SetBlockDeviceMappings(v []*BlockDeviceMapping) *DescribeImageAttributeOutput {
	s.BlockDeviceMappings = v
	return s
}

// SetDescription sets the Description field's value.
func (s *DescribeImageAttributeOutput) SetDescription(v *AttributeValue) *DescribeImageAttributeOutput {
	s.Description = v
	return s
}

// SetImageId sets the ImageId field's value.
func (s *DescribeImageAttributeOutput) SetImageId(v string) *DescribeImageAttributeOutput {
	s.ImageId = &v
	return s
}

// SetKernelId sets the KernelId field's value.
func (s *DescribeImageAttributeOutput) SetKernelId(v *AttributeValue) *DescribeImageAttributeOutput {
	s.KernelId = v
	return s
}

// SetLaunchPermissions sets the LaunchPermissions field's value.
func (s *DescribeImageAttributeOutput) SetLaunchPermissions(v []*LaunchPermission) *DescribeImageAttributeOutput {
	s.LaunchPermissions = v
	return s
}

// SetProductCodes sets the ProductCodes field's value.
func (s *DescribeImageAttributeOutput) SetProductCodes(v []*ProductCode) *DescribeImageAttributeOutput {
	s.ProductCodes = v
	return s
}

// SetRamdiskId sets the RamdiskId field's value.
func (s *DescribeImageAttributeOutput) SetRamdiskId(v *AttributeValue) *DescribeImageAttributeOutput {
	s.RamdiskId = v
	return s
}

// SetSriovNetSupport sets the SriovNetSupport field's value.
func (s *DescribeImageAttributeOutput) SetSriovNetSupport(v *AttributeValue) *DescribeImageAttributeOutput {
	s.SriovNetSupport = v
	return s
}

type CreateVpcPeeringConnectionInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The AWS account ID of the owner of the peer VPC.
	//
	// Default: Your AWS account ID
	PeerOwnerId *string `locationName:"peerOwnerId" type:"string"`

	// The ID of the VPC with which you are creating the VPC peering connection.
	PeerVpcId *string `locationName:"peerVpcId" type:"string"`

	// The ID of the requester VPC.
	VpcId *string `locationName:"vpcId" type:"string"`
}

// String returns the string representation
func (s CreateVpcPeeringConnectionInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateVpcPeeringConnectionInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *CreateVpcPeeringConnectionInput) SetDryRun(v bool) *CreateVpcPeeringConnectionInput {
	s.DryRun = &v
	return s
}

// SetPeerOwnerId sets the PeerOwnerId field's value.
func (s *CreateVpcPeeringConnectionInput) SetPeerOwnerId(v string) *CreateVpcPeeringConnectionInput {
	s.PeerOwnerId = &v
	return s
}

// SetPeerVpcId sets the PeerVpcId field's value.
func (s *CreateVpcPeeringConnectionInput) SetPeerVpcId(v string) *CreateVpcPeeringConnectionInput {
	s.PeerVpcId = &v
	return s
}

// SetVpcId sets the VpcId field's value.
func (s *CreateVpcPeeringConnectionInput) SetVpcId(v string) *CreateVpcPeeringConnectionInput {
	s.VpcId = &v
	return s
}

// Contains the output of CreateVpcPeeringConnection.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateVpcPeeringConnectionResult
type CreateVpcPeeringConnectionOutput struct {
	_ struct{} `type:"structure"`

	// Information about the VPC peering connection.
	VpcPeeringConnection *VpcPeeringConnection `locationName:"vpcPeeringConnection" type:"structure"`
}

// String returns the string representation
func (s CreateVpcPeeringConnectionOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateVpcPeeringConnectionOutput) GoString() string {
	return s.String()
}

// SetVpcPeeringConnection sets the VpcPeeringConnection field's value.
func (s *CreateVpcPeeringConnectionOutput) SetVpcPeeringConnection(v *VpcPeeringConnection) *CreateVpcPeeringConnectionOutput {
	s.VpcPeeringConnection = v
	return s
}

type VpcPeeringConnection struct {
	_ struct{} `type:"structure"`

	// Information about the accepter VPC. CIDR block information is not returned
	// when creating a VPC peering connection, or when describing a VPC peering
	// connection that's in the initiating-request or pending-acceptance state.
	AccepterVpcInfo *VpcPeeringConnectionVpcInfo `locationName:"accepterVpcInfo" type:"structure"`

	// The time that an unaccepted VPC peering connection will expire.
	ExpirationTime *time.Time `locationName:"expirationTime" type:"timestamp" timestampFormat:"iso8601"`

	// Information about the requester VPC.
	RequesterVpcInfo *VpcPeeringConnectionVpcInfo `locationName:"requesterVpcInfo" type:"structure"`

	// The status of the VPC peering connection.
	Status *VpcPeeringConnectionStateReason `locationName:"status" type:"structure"`

	// Any tags assigned to the resource.
	Tags []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`

	// The ID of the VPC peering connection.
	VpcPeeringConnectionId *string `locationName:"vpcPeeringConnectionId" type:"string"`
}

type VpcPeeringConnectionVpcInfo struct {
	_ struct{} `type:"structure"`

	// The IPv4 CIDR block for the VPC.
	CidrBlock *string `locationName:"cidrBlock" type:"string"`

	// The IPv6 CIDR block for the VPC.
	Ipv6CidrBlockSet []*Ipv6CidrBlock `locationName:"ipv6CidrBlockSet" locationNameList:"item" type:"list"`

	// The AWS account ID of the VPC owner.
	OwnerId *string `locationName:"ownerId" type:"string"`

	// Information about the VPC peering connection options for the accepter or
	// requester VPC.
	PeeringOptions *VpcPeeringConnectionOptionsDescription `locationName:"peeringOptions" type:"structure"`

	// The ID of the VPC.
	VpcId *string `locationName:"vpcId" type:"string"`
}

type Ipv6CidrBlock struct {
	_ struct{} `type:"structure"`

	// The IPv6 CIDR block.
	Ipv6CidrBlock *string `locationName:"ipv6CidrBlock" type:"string"`
}

type VpcPeeringConnectionOptionsDescription struct {
	_ struct{} `type:"structure"`

	// Indicates whether a local VPC can resolve public DNS hostnames to private
	// IP addresses when queried from instances in a peer VPC.
	AllowDnsResolutionFromRemoteVpc *bool `locationName:"allowDnsResolutionFromRemoteVpc" type:"boolean"`

	// Indicates whether a local ClassicLink connection can communicate with the
	// peer VPC over the VPC peering connection.
	AllowEgressFromLocalClassicLinkToRemoteVpc *bool `locationName:"allowEgressFromLocalClassicLinkToRemoteVpc" type:"boolean"`

	// Indicates whether a local VPC can communicate with a ClassicLink connection
	// in the peer VPC over the VPC peering connection.
	AllowEgressFromLocalVpcToRemoteClassicLink *bool `locationName:"allowEgressFromLocalVpcToRemoteClassicLink" type:"boolean"`
}

type VpcPeeringConnectionStateReason struct {
	_ struct{} `type:"structure"`

	// The status of the VPC peering connection.
	Code *string `locationName:"code" type:"string" enum:"VpcPeeringConnectionStateReasonCode"`

	// A message that provides more information about the status, if applicable.
	Message *string `locationName:"message" type:"string"`
}

type DescribeVpcPeeringConnectionsInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more filters.
	//
	//    * accepter-vpc-info.cidr-block - The IPv4 CIDR block of the peer VPC.
	//
	//    * accepter-vpc-info.owner-id - The AWS account ID of the owner of the
	//    peer VPC.
	//
	//    * accepter-vpc-info.vpc-id - The ID of the peer VPC.
	//
	//    * expiration-time - The expiration date and time for the VPC peering connection.
	//
	//    * requester-vpc-info.cidr-block - The IPv4 CIDR block of the requester's
	//    VPC.
	//
	//    * requester-vpc-info.owner-id - The AWS account ID of the owner of the
	//    requester VPC.
	//
	//    * requester-vpc-info.vpc-id - The ID of the requester VPC.
	//
	//    * status-code - The status of the VPC peering connection (pending-acceptance
	//    | failed | expired | provisioning | active | deleted | rejected).
	//
	//    * status-message - A message that provides more information about the
	//    status of the VPC peering connection, if applicable.
	//
	//    * tag:key=value - The key/value combination of a tag assigned to the resource.
	//    Specify the key of the tag in the filter name and the value of the tag
	//    in the filter value. For example, for the tag Purpose=X, specify tag:Purpose
	//    for the filter name and X for the filter value.
	//
	//    * tag-key - The key of a tag assigned to the resource. This filter is
	//    independent of the tag-value filter. For example, if you use both the
	//    filter "tag-key=Purpose" and the filter "tag-value=X", you get any resources
	//    assigned both the tag key Purpose (regardless of what the tag's value
	//    is), and the tag value X (regardless of what the tag's key is). If you
	//    want to list only resources where Purpose is X, see the tag:key=value
	//    filter.
	//
	//    * tag-value - The value of a tag assigned to the resource. This filter
	//    is independent of the tag-key filter.
	//
	//    * vpc-peering-connection-id - The ID of the VPC peering connection.
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// One or more VPC peering connection IDs.
	//
	// Default: Describes all your VPC peering connections.
	VpcPeeringConnectionIds []*string `locationName:"VpcPeeringConnectionId" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s DescribeVpcPeeringConnectionsInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeVpcPeeringConnectionsInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeVpcPeeringConnectionsInput) SetDryRun(v bool) *DescribeVpcPeeringConnectionsInput {
	s.DryRun = &v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribeVpcPeeringConnectionsInput) SetFilters(v []*Filter) *DescribeVpcPeeringConnectionsInput {
	s.Filters = v
	return s
}

// SetVpcPeeringConnectionIds sets the VpcPeeringConnectionIds field's value.
func (s *DescribeVpcPeeringConnectionsInput) SetVpcPeeringConnectionIds(v []*string) *DescribeVpcPeeringConnectionsInput {
	s.VpcPeeringConnectionIds = v
	return s
}

// Contains the output of DescribeVpcPeeringConnections.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeVpcPeeringConnectionsResult
type DescribeVpcPeeringConnectionsOutput struct {
	_ struct{} `type:"structure"`

	// Information about the VPC peering connections.
	VpcPeeringConnections []*VpcPeeringConnection `locationName:"vpcPeeringConnectionSet" locationNameList:"item" type:"list"`

	RequestId *string `locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s DescribeVpcPeeringConnectionsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeVpcPeeringConnectionsOutput) GoString() string {
	return s.String()
}

// SetVpcPeeringConnections sets the VpcPeeringConnections field's value.
func (s *DescribeVpcPeeringConnectionsOutput) SetVpcPeeringConnections(v []*VpcPeeringConnection) *DescribeVpcPeeringConnectionsOutput {
	s.VpcPeeringConnections = v
	return s
}

type AcceptVpcPeeringConnectionInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the VPC peering connection.
	VpcPeeringConnectionId *string `locationName:"vpcPeeringConnectionId" type:"string"`
}

// String returns the string representation
func (s AcceptVpcPeeringConnectionInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AcceptVpcPeeringConnectionInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *AcceptVpcPeeringConnectionInput) SetDryRun(v bool) *AcceptVpcPeeringConnectionInput {
	s.DryRun = &v
	return s
}

// SetVpcPeeringConnectionId sets the VpcPeeringConnectionId field's value.
func (s *AcceptVpcPeeringConnectionInput) SetVpcPeeringConnectionId(v string) *AcceptVpcPeeringConnectionInput {
	s.VpcPeeringConnectionId = &v
	return s
}

// Contains the output of AcceptVpcPeeringConnection.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/AcceptVpcPeeringConnectionResult
type AcceptVpcPeeringConnectionOutput struct {
	_ struct{} `type:"structure"`

	// Information about the VPC peering connection.
	VpcPeeringConnection *VpcPeeringConnection `locationName:"vpcPeeringConnection" type:"structure"`
}

// String returns the string representation
func (s AcceptVpcPeeringConnectionOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AcceptVpcPeeringConnectionOutput) GoString() string {
	return s.String()
}

// SetVpcPeeringConnection sets the VpcPeeringConnection field's value.
func (s *AcceptVpcPeeringConnectionOutput) SetVpcPeeringConnection(v *VpcPeeringConnection) *AcceptVpcPeeringConnectionOutput {
	s.VpcPeeringConnection = v
	return s
}

type ModifyVpcPeeringConnectionOptionsInput struct {
	_ struct{} `type:"structure"`

	// The VPC peering connection options for the accepter VPC.
	AccepterPeeringConnectionOptions *PeeringConnectionOptionsRequest `type:"structure"`

	// Checks whether you have the required permissions for the operation, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `type:"boolean"`

	// The VPC peering connection options for the requester VPC.
	RequesterPeeringConnectionOptions *PeeringConnectionOptionsRequest `type:"structure"`

	// The ID of the VPC peering connection.
	//
	// VpcPeeringConnectionId is a required field
	VpcPeeringConnectionId *string `type:"string" required:"true"`
}

// String returns the string representation
func (s ModifyVpcPeeringConnectionOptionsInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ModifyVpcPeeringConnectionOptionsInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *ModifyVpcPeeringConnectionOptionsInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "ModifyVpcPeeringConnectionOptionsInput"}
	if s.VpcPeeringConnectionId == nil {
		invalidParams.Add(request.NewErrParamRequired("VpcPeeringConnectionId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAccepterPeeringConnectionOptions sets the AccepterPeeringConnectionOptions field's value.
func (s *ModifyVpcPeeringConnectionOptionsInput) SetAccepterPeeringConnectionOptions(v *PeeringConnectionOptionsRequest) *ModifyVpcPeeringConnectionOptionsInput {
	s.AccepterPeeringConnectionOptions = v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *ModifyVpcPeeringConnectionOptionsInput) SetDryRun(v bool) *ModifyVpcPeeringConnectionOptionsInput {
	s.DryRun = &v
	return s
}

// SetRequesterPeeringConnectionOptions sets the RequesterPeeringConnectionOptions field's value.
func (s *ModifyVpcPeeringConnectionOptionsInput) SetRequesterPeeringConnectionOptions(v *PeeringConnectionOptionsRequest) *ModifyVpcPeeringConnectionOptionsInput {
	s.RequesterPeeringConnectionOptions = v
	return s
}

// SetVpcPeeringConnectionId sets the VpcPeeringConnectionId field's value.
func (s *ModifyVpcPeeringConnectionOptionsInput) SetVpcPeeringConnectionId(v string) *ModifyVpcPeeringConnectionOptionsInput {
	s.VpcPeeringConnectionId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/ModifyVpcPeeringConnectionOptionsResult
type ModifyVpcPeeringConnectionOptionsOutput struct {
	_ struct{} `type:"structure"`

	// Information about the VPC peering connection options for the accepter VPC.
	AccepterPeeringConnectionOptions *PeeringConnectionOptions `locationName:"accepterPeeringConnectionOptions" type:"structure"`

	// Information about the VPC peering connection options for the requester VPC.
	RequesterPeeringConnectionOptions *PeeringConnectionOptions `locationName:"requesterPeeringConnectionOptions" type:"structure"`
}

// String returns the string representation
func (s ModifyVpcPeeringConnectionOptionsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ModifyVpcPeeringConnectionOptionsOutput) GoString() string {
	return s.String()
}

// SetAccepterPeeringConnectionOptions sets the AccepterPeeringConnectionOptions field's value.
func (s *ModifyVpcPeeringConnectionOptionsOutput) SetAccepterPeeringConnectionOptions(v *PeeringConnectionOptions) *ModifyVpcPeeringConnectionOptionsOutput {
	s.AccepterPeeringConnectionOptions = v
	return s
}

// SetRequesterPeeringConnectionOptions sets the RequesterPeeringConnectionOptions field's value.
func (s *ModifyVpcPeeringConnectionOptionsOutput) SetRequesterPeeringConnectionOptions(v *PeeringConnectionOptions) *ModifyVpcPeeringConnectionOptionsOutput {
	s.RequesterPeeringConnectionOptions = v
	return s
}

type PeeringConnectionOptions struct {
	_ struct{} `type:"structure"`

	// If true, enables a local VPC to resolve public DNS hostnames to private IP
	// addresses when queried from instances in the peer VPC.
	AllowDnsResolutionFromRemoteVpc *bool `locationName:"allowDnsResolutionFromRemoteVpc" type:"boolean"`

	// If true, enables outbound communication from an EC2-Classic instance that's
	// linked to a local VPC via ClassicLink to instances in a peer VPC.
	AllowEgressFromLocalClassicLinkToRemoteVpc *bool `locationName:"allowEgressFromLocalClassicLinkToRemoteVpc" type:"boolean"`

	// If true, enables outbound communication from instances in a local VPC to
	// an EC2-Classic instance that's linked to a peer VPC via ClassicLink.
	AllowEgressFromLocalVpcToRemoteClassicLink *bool `locationName:"allowEgressFromLocalVpcToRemoteClassicLink" type:"boolean"`
}

type PeeringConnectionOptionsRequest struct {
	_ struct{} `type:"structure"`

	// If true, enables a local VPC to resolve public DNS hostnames to private IP
	// addresses when queried from instances in the peer VPC.
	AllowDnsResolutionFromRemoteVpc *bool `type:"boolean"`

	// If true, enables outbound communication from an EC2-Classic instance that's
	// linked to a local VPC via ClassicLink to instances in a peer VPC.
	AllowEgressFromLocalClassicLinkToRemoteVpc *bool `type:"boolean"`

	// If true, enables outbound communication from instances in a local VPC to
	// an EC2-Classic instance that's linked to a peer VPC via ClassicLink.
	AllowEgressFromLocalVpcToRemoteClassicLink *bool `type:"boolean"`
}

type DeleteVpcPeeringConnectionInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the VPC peering connection.
	//
	// VpcPeeringConnectionId is a required field
	VpcPeeringConnectionId *string `locationName:"vpcPeeringConnectionId" type:"string" required:"true"`
}

// String returns the string representation
func (s DeleteVpcPeeringConnectionInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteVpcPeeringConnectionInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteVpcPeeringConnectionInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeleteVpcPeeringConnectionInput"}
	if s.VpcPeeringConnectionId == nil {
		invalidParams.Add(request.NewErrParamRequired("VpcPeeringConnectionId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDryRun sets the DryRun field's value.
func (s *DeleteVpcPeeringConnectionInput) SetDryRun(v bool) *DeleteVpcPeeringConnectionInput {
	s.DryRun = &v
	return s
}

// SetVpcPeeringConnectionId sets the VpcPeeringConnectionId field's value.
func (s *DeleteVpcPeeringConnectionInput) SetVpcPeeringConnectionId(v string) *DeleteVpcPeeringConnectionInput {
	s.VpcPeeringConnectionId = &v
	return s
}

// Contains the output of DeleteVpcPeeringConnection.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteVpcPeeringConnectionResult
type DeleteVpcPeeringConnectionOutput struct {
	_ struct{} `type:"structure"`

	// Returns true if the request succeeds; otherwise, it returns an error.
	Return *bool `locationName:"return" type:"boolean"`
}

// String returns the string representation
func (s DeleteVpcPeeringConnectionOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteVpcPeeringConnectionOutput) GoString() string {
	return s.String()
}

// SetReturn sets the Return field's value.
func (s *DeleteVpcPeeringConnectionOutput) SetReturn(v bool) *DeleteVpcPeeringConnectionOutput {
	s.Return = &v
	return s
}

//
//
// Create Network Interface

// Contains the parameters for CreateNetworkInterface.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateNetworkInterfaceRequest
type CreateNetworkInterfaceInput struct {
	_ struct{} `type:"structure"`

	// A description for the network interface.
	Description *string `locationName:"description" type:"string"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The IDs of one or more security groups.
	Groups []*string `locationName:"SecurityGroupId" locationNameList:"SecurityGroupId" type:"list"`

	// The number of IPv6 addresses to assign to a network interface. Amazon EC2
	// automatically selects the IPv6 addresses from the subnet range. You can't
	// use this option if specifying specific IPv6 addresses. If your subnet has
	// the AssignIpv6AddressOnCreation attribute set to true, you can specify 0
	// to override this setting.
	Ipv6AddressCount *int64 `locationName:"ipv6AddressCount" type:"integer"`

	// One or more specific IPv6 addresses from the IPv6 CIDR block range of your
	// subnet. You can't use this option if you're specifying a number of IPv6 addresses.
	Ipv6Addresses []*InstanceIpv6Address `locationName:"ipv6Addresses" locationNameList:"item" type:"list"`

	// The primary private IPv4 address of the network interface. If you don't specify
	// an IPv4 address, Amazon EC2 selects one for you from the subnet's IPv4 CIDR
	// range. If you specify an IP address, you cannot indicate any IP addresses
	// specified in privateIpAddresses as primary (only one IP address can be designated
	// as primary).
	PrivateIpAddress *string `locationName:"privateIpAddress" type:"string"`

	// One or more private IPv4 addresses.
	PrivateIpAddresses []*PrivateIpAddressSpecification `locationName:"privateIpAddresses" locationNameList:"item" type:"list"`

	// The number of secondary private IPv4 addresses to assign to a network interface.
	// When you specify a number of secondary IPv4 addresses, Amazon EC2 selects
	// these IP addresses within the subnet's IPv4 CIDR range. You can't specify
	// this option and specify more than one private IP address using privateIpAddresses.
	//
	// The number of IP addresses you can assign to a network interface varies by
	// instance type. For more information, see IP Addresses Per ENI Per Instance
	// Type (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-eni.html#AvailableIpPerENI)
	// in the Amazon Virtual Private Cloud User Guide.
	SecondaryPrivateIpAddressCount *int64 `locationName:"secondaryPrivateIpAddressCount" type:"integer"`

	// The ID of the subnet to associate with the network interface.
	//
	// SubnetId is a required field
	SubnetId *string `locationName:"subnetId" type:"string" required:"true"`
}

// String returns the string representation
func (s CreateNetworkInterfaceInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateNetworkInterfaceInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CreateNetworkInterfaceInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "CreateNetworkInterfaceInput"}
	if s.SubnetId == nil {
		invalidParams.Add(request.NewErrParamRequired("SubnetId"))
	}
	if s.PrivateIpAddresses != nil {
		for i, v := range s.PrivateIpAddresses {
			if v == nil {
				continue
			}
			if err := s.Validate(); err != nil {
				invalidParams.AddNested(fmt.Sprintf("%s[%v]", "PrivateIpAddresses", i), err.(request.ErrInvalidParams))
			}
		}
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDescription sets the Description field's value.
func (s *CreateNetworkInterfaceInput) SetDescription(v string) *CreateNetworkInterfaceInput {
	s.Description = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *CreateNetworkInterfaceInput) SetDryRun(v bool) *CreateNetworkInterfaceInput {
	s.DryRun = &v
	return s
}

// SetGroups sets the Groups field's value.
func (s *CreateNetworkInterfaceInput) SetGroups(v []*string) *CreateNetworkInterfaceInput {
	s.Groups = v
	return s
}

// SetIpv6AddressCount sets the Ipv6AddressCount field's value.
func (s *CreateNetworkInterfaceInput) SetIpv6AddressCount(v int64) *CreateNetworkInterfaceInput {
	s.Ipv6AddressCount = &v
	return s
}

// SetIpv6Addresses sets the Ipv6Addresses field's value.
func (s *CreateNetworkInterfaceInput) SetIpv6Addresses(v []*InstanceIpv6Address) *CreateNetworkInterfaceInput {
	s.Ipv6Addresses = v
	return s
}

// SetPrivateIpAddress sets the PrivateIpAddress field's value.
func (s *CreateNetworkInterfaceInput) SetPrivateIpAddress(v string) *CreateNetworkInterfaceInput {
	s.PrivateIpAddress = &v
	return s
}

// SetPrivateIpAddresses sets the PrivateIpAddresses field's value.
func (s *CreateNetworkInterfaceInput) SetPrivateIpAddresses(v []*PrivateIpAddressSpecification) *CreateNetworkInterfaceInput {
	s.PrivateIpAddresses = v
	return s
}

// SetSecondaryPrivateIpAddressCount sets the SecondaryPrivateIpAddressCount field's value.
func (s *CreateNetworkInterfaceInput) SetSecondaryPrivateIpAddressCount(v int64) *CreateNetworkInterfaceInput {
	s.SecondaryPrivateIpAddressCount = &v
	return s
}

// SetSubnetId sets the SubnetId field's value.
func (s *CreateNetworkInterfaceInput) SetSubnetId(v string) *CreateNetworkInterfaceInput {
	s.SubnetId = &v
	return s
}

// Contains the output of CreateNetworkInterface.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateNetworkInterfaceResult
type CreateNetworkInterfaceOutput struct {
	_ struct{} `type:"structure"`

	// Information about the network interface.
	NetworkInterface *NetworkInterface `locationName:"networkInterface" type:"structure"`
}

// String returns the string representation
func (s CreateNetworkInterfaceOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateNetworkInterfaceOutput) GoString() string {
	return s.String()
}

// SetNetworkInterface sets the NetworkInterface field's value.
func (s *CreateNetworkInterfaceOutput) SetNetworkInterface(v *NetworkInterface) *CreateNetworkInterfaceOutput {
	s.NetworkInterface = v
	return s
}

// Contains the parameters for DeleteNetworkInterface.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteNetworkInterfaceRequest
type DeleteNetworkInterfaceInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the network interface.
	//
	// NetworkInterfaceId is a required field
	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string" required:"true"`
}

// String returns the string representation
func (s DeleteNetworkInterfaceInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteNetworkInterfaceInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteNetworkInterfaceInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeleteNetworkInterfaceInput"}
	if s.NetworkInterfaceId == nil {
		invalidParams.Add(request.NewErrParamRequired("NetworkInterfaceId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDryRun sets the DryRun field's value.
func (s *DeleteNetworkInterfaceInput) SetDryRun(v bool) *DeleteNetworkInterfaceInput {
	s.DryRun = &v
	return s
}

// SetNetworkInterfaceId sets the NetworkInterfaceId field's value.
func (s *DeleteNetworkInterfaceInput) SetNetworkInterfaceId(v string) *DeleteNetworkInterfaceInput {
	s.NetworkInterfaceId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteNetworkInterfaceOutput
type DeleteNetworkInterfaceOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeleteNetworkInterfaceOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteNetworkInterfaceOutput) GoString() string {
	return s.String()
}

// Contains the parameters for DescribeNetworkInterfaces.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeNetworkInterfacesRequest
type DescribeNetworkInterfacesInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more filters.
	//
	//    * addresses.private-ip-address - The private IPv4 addresses associated
	//    with the network interface.
	//
	//    * addresses.primary - Whether the private IPv4 address is the primary
	//    IP address associated with the network interface.
	//
	//    * addresses.association.public-ip - The association ID returned when the
	//    network interface was associated with the Elastic IP address (IPv4).
	//
	//    * addresses.association.owner-id - The owner ID of the addresses associated
	//    with the network interface.
	//
	//    * association.association-id - The association ID returned when the network
	//    interface was associated with an IPv4 address.
	//
	//    * association.allocation-id - The allocation ID returned when you allocated
	//    the Elastic IP address (IPv4) for your network interface.
	//
	//    * association.ip-owner-id - The owner of the Elastic IP address (IPv4)
	//    associated with the network interface.
	//
	//    * association.public-ip - The address of the Elastic IP address (IPv4)
	//    bound to the network interface.
	//
	//    * association.public-dns-name - The public DNS name for the network interface
	//    (IPv4).
	//
	//    * attachment.attachment-id - The ID of the interface attachment.
	//
	//    * attachment.attach.time - The time that the network interface was attached
	//    to an instance.
	//
	//    * attachment.delete-on-termination - Indicates whether the attachment
	//    is deleted when an instance is terminated.
	//
	//    * attachment.device-index - The device index to which the network interface
	//    is attached.
	//
	//    * attachment.instance-id - The ID of the instance to which the network
	//    interface is attached.
	//
	//    * attachment.instance-owner-id - The owner ID of the instance to which
	//    the network interface is attached.
	//
	//    * attachment.nat-gateway-id - The ID of the NAT gateway to which the network
	//    interface is attached.
	//
	//    * attachment.status - The status of the attachment (attaching | attached
	//    | detaching | detached).
	//
	//    * availability-zone - The Availability Zone of the network interface.
	//
	//    * description - The description of the network interface.
	//
	//    * group-id - The ID of a security group associated with the network interface.
	//
	//    * group-name - The name of a security group associated with the network
	//    interface.
	//
	//    * ipv6-addresses.ipv6-address - An IPv6 address associated with the network
	//    interface.
	//
	//    * mac-address - The MAC address of the network interface.
	//
	//    * network-interface-id - The ID of the network interface.
	//
	//    * owner-id - The AWS account ID of the network interface owner.
	//
	//    * private-ip-address - The private IPv4 address or addresses of the network
	//    interface.
	//
	//    * private-dns-name - The private DNS name of the network interface (IPv4).
	//
	//    * requester-id - The ID of the entity that launched the instance on your
	//    behalf (for example, AWS Management Console, Auto Scaling, and so on).
	//
	//    * requester-managed - Indicates whether the network interface is being
	//    managed by an AWS service (for example, AWS Management Console, Auto Scaling,
	//    and so on).
	//
	//    * source-desk-check - Indicates whether the network interface performs
	//    source/destination checking. A value of true means checking is enabled,
	//    and false means checking is disabled. The value must be false for the
	//    network interface to perform network address translation (NAT) in your
	//    VPC.
	//
	//    * status - The status of the network interface. If the network interface
	//    is not attached to an instance, the status is available; if a network
	//    interface is attached to an instance the status is in-use.
	//
	//    * subnet-id - The ID of the subnet for the network interface.
	//
	//    * tag:key=value - The key/value combination of a tag assigned to the resource.
	//    Specify the key of the tag in the filter name and the value of the tag
	//    in the filter value. For example, for the tag Purpose=X, specify tag:Purpose
	//    for the filter name and X for the filter value.
	//
	//    * tag-key - The key of a tag assigned to the resource. This filter is
	//    independent of the tag-value filter. For example, if you use both the
	//    filter "tag-key=Purpose" and the filter "tag-value=X", you get any resources
	//    assigned both the tag key Purpose (regardless of what the tag's value
	//    is), and the tag value X (regardless of what the tag's key is). If you
	//    want to list only resources where Purpose is X, see the tag:key=value
	//    filter.
	//
	//    * tag-value - The value of a tag assigned to the resource. This filter
	//    is independent of the tag-key filter.
	//
	//    * vpc-id - The ID of the VPC for the network interface.
	Filters []*Filter `locationName:"filter" locationNameList:"Filter" type:"list"`

	// One or more network interface IDs.
	//
	// Default: Describes all your network interfaces.
	NetworkInterfaceIds []*string `locationName:"NetworkInterfaceId" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s DescribeNetworkInterfacesInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeNetworkInterfacesInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeNetworkInterfacesInput) SetDryRun(v bool) *DescribeNetworkInterfacesInput {
	s.DryRun = &v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribeNetworkInterfacesInput) SetFilters(v []*Filter) *DescribeNetworkInterfacesInput {
	s.Filters = v
	return s
}

// SetNetworkInterfaceIds sets the NetworkInterfaceIds field's value.
func (s *DescribeNetworkInterfacesInput) SetNetworkInterfaceIds(v []*string) *DescribeNetworkInterfacesInput {
	s.NetworkInterfaceIds = v
	return s
}

// Contains the output of DescribeNetworkInterfaces.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeNetworkInterfacesResult
type DescribeNetworkInterfacesOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more network interfaces.
	NetworkInterfaces []*NetworkInterface `locationName:"networkInterfaceSet" locationNameList:"item" type:"list"`

	RequestId *string `locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s DescribeNetworkInterfacesOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeNetworkInterfacesOutput) GoString() string {
	return s.String()
}

// SetNetworkInterfaces sets the NetworkInterfaces field's value.
func (s *DescribeNetworkInterfacesOutput) SetNetworkInterfaces(v []*NetworkInterface) *DescribeNetworkInterfacesOutput {
	s.NetworkInterfaces = v
	return s
}

// Contains the parameters for ModifyNetworkInterfaceAttribute.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/ModifyNetworkInterfaceAttributeRequest
type ModifyNetworkInterfaceAttributeInput struct {
	_ struct{} `type:"structure"`

	// Information about the interface attachment. If modifying the 'delete on termination'
	// attribute, you must specify the ID of the interface attachment.
	Attachment *NetworkInterfaceAttachmentChanges `locationName:"attachment" type:"structure"`

	// A description for the network interface.
	Description *AttributeValue `locationName:"description" type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// Changes the security groups for the network interface. The new set of groups
	// you specify replaces the current set. You must specify at least one group,
	// even if it's just the default security group in the VPC. You must specify
	// the ID of the security group, not the name.
	Groups []*string `locationName:"SecurityGroupId" locationNameList:"SecurityGroupId" type:"list"`

	// The ID of the network interface.
	//
	// NetworkInterfaceId is a required field
	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string" required:"true"`

	// Indicates whether source/destination checking is enabled. A value of true
	// means checking is enabled, and false means checking is disabled. This value
	// must be false for a NAT instance to perform NAT. For more information, see
	// NAT Instances (http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/VPC_NAT_Instance.html)
	// in the Amazon Virtual Private Cloud User Guide.
	SourceDestCheck *AttributeBooleanValue `locationName:"sourceDestCheck" type:"structure"`
}

// String returns the string representation
func (s ModifyNetworkInterfaceAttributeInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ModifyNetworkInterfaceAttributeInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *ModifyNetworkInterfaceAttributeInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "ModifyNetworkInterfaceAttributeInput"}
	if s.NetworkInterfaceId == nil {
		invalidParams.Add(request.NewErrParamRequired("NetworkInterfaceId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAttachment sets the Attachment field's value.
func (s *ModifyNetworkInterfaceAttributeInput) SetAttachment(v *NetworkInterfaceAttachmentChanges) *ModifyNetworkInterfaceAttributeInput {
	s.Attachment = v
	return s
}

// SetDescription sets the Description field's value.
func (s *ModifyNetworkInterfaceAttributeInput) SetDescription(v *AttributeValue) *ModifyNetworkInterfaceAttributeInput {
	s.Description = v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *ModifyNetworkInterfaceAttributeInput) SetDryRun(v bool) *ModifyNetworkInterfaceAttributeInput {
	s.DryRun = &v
	return s
}

// SetGroups sets the Groups field's value.
func (s *ModifyNetworkInterfaceAttributeInput) SetGroups(v []*string) *ModifyNetworkInterfaceAttributeInput {
	s.Groups = v
	return s
}

// SetNetworkInterfaceId sets the NetworkInterfaceId field's value.
func (s *ModifyNetworkInterfaceAttributeInput) SetNetworkInterfaceId(v string) *ModifyNetworkInterfaceAttributeInput {
	s.NetworkInterfaceId = &v
	return s
}

// SetSourceDestCheck sets the SourceDestCheck field's value.
func (s *ModifyNetworkInterfaceAttributeInput) SetSourceDestCheck(v *AttributeBooleanValue) *ModifyNetworkInterfaceAttributeInput {
	s.SourceDestCheck = v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/ModifyNetworkInterfaceAttributeOutput
type ModifyNetworkInterfaceAttributeOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s ModifyNetworkInterfaceAttributeOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ModifyNetworkInterfaceAttributeOutput) GoString() string {
	return s.String()
}

// Contains the parameters for DescribeNetworkInterfaceAttribute.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeNetworkInterfaceAttributeRequest
type DescribeNetworkInterfaceAttributeInput struct {
	_ struct{} `type:"structure"`

	// The attribute of the network interface. This parameter is required.
	Attribute *string `locationName:"attribute" type:"string" enum:"NetworkInterfaceAttribute"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the network interface.
	//
	// NetworkInterfaceId is a required field
	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string" required:"true"`
}

// String returns the string representation
func (s DescribeNetworkInterfaceAttributeInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeNetworkInterfaceAttributeInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DescribeNetworkInterfaceAttributeInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DescribeNetworkInterfaceAttributeInput"}
	if s.NetworkInterfaceId == nil {
		invalidParams.Add(request.NewErrParamRequired("NetworkInterfaceId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAttribute sets the Attribute field's value.
func (s *DescribeNetworkInterfaceAttributeInput) SetAttribute(v string) *DescribeNetworkInterfaceAttributeInput {
	s.Attribute = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeNetworkInterfaceAttributeInput) SetDryRun(v bool) *DescribeNetworkInterfaceAttributeInput {
	s.DryRun = &v
	return s
}

// SetNetworkInterfaceId sets the NetworkInterfaceId field's value.
func (s *DescribeNetworkInterfaceAttributeInput) SetNetworkInterfaceId(v string) *DescribeNetworkInterfaceAttributeInput {
	s.NetworkInterfaceId = &v
	return s
}

// Contains the output of DescribeNetworkInterfaceAttribute.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeNetworkInterfaceAttributeResult
type DescribeNetworkInterfaceAttributeOutput struct {
	_ struct{} `type:"structure"`

	// The attachment (if any) of the network interface.
	Attachment *NetworkInterfaceAttachment `locationName:"attachment" type:"structure"`

	// The description of the network interface.
	Description *AttributeValue `locationName:"description" type:"structure"`

	// The security groups associated with the network interface.
	Groups []*GroupIdentifier `locationName:"groupSet" locationNameList:"item" type:"list"`

	// The ID of the network interface.
	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string"`

	// Indicates whether source/destination checking is enabled.
	SourceDestCheck *AttributeBooleanValue `locationName:"sourceDestCheck" type:"structure"`
}

// String returns the string representation
func (s DescribeNetworkInterfaceAttributeOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeNetworkInterfaceAttributeOutput) GoString() string {
	return s.String()
}

// SetAttachment sets the Attachment field's value.
func (s *DescribeNetworkInterfaceAttributeOutput) SetAttachment(v *NetworkInterfaceAttachment) *DescribeNetworkInterfaceAttributeOutput {
	s.Attachment = v
	return s
}

// SetDescription sets the Description field's value.
func (s *DescribeNetworkInterfaceAttributeOutput) SetDescription(v *AttributeValue) *DescribeNetworkInterfaceAttributeOutput {
	s.Description = v
	return s
}

// SetGroups sets the Groups field's value.
func (s *DescribeNetworkInterfaceAttributeOutput) SetGroups(v []*GroupIdentifier) *DescribeNetworkInterfaceAttributeOutput {
	s.Groups = v
	return s
}

// SetNetworkInterfaceId sets the NetworkInterfaceId field's value.
func (s *DescribeNetworkInterfaceAttributeOutput) SetNetworkInterfaceId(v string) *DescribeNetworkInterfaceAttributeOutput {
	s.NetworkInterfaceId = &v
	return s
}

// SetSourceDestCheck sets the SourceDestCheck field's value.
func (s *DescribeNetworkInterfaceAttributeOutput) SetSourceDestCheck(v *AttributeBooleanValue) *DescribeNetworkInterfaceAttributeOutput {
	s.SourceDestCheck = v
	return s
}

// Describes a network interface.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/NetworkInterface
type NetworkInterface struct {
	_ struct{} `type:"structure"`

	// The association information for an Elastic IP address (IPv4) associated with
	// the network interface.
	Association *NetworkInterfaceAssociation `locationName:"association" type:"structure"`

	// The network interface attachment.
	Attachment *NetworkInterfaceAttachment `locationName:"attachment" type:"structure"`

	// The Availability Zone.
	AvailabilityZone *string `locationName:"availabilityZone" type:"string"`

	// A description.
	Description *string `locationName:"description" type:"string"`

	// Any security groups for the network interface.
	Groups []*GroupIdentifier `locationName:"groupSet" locationNameList:"item" type:"list"`

	// The type of interface.
	InterfaceType *string `locationName:"interfaceType" type:"string" enum:"NetworkInterfaceType"`

	// The IPv6 addresses associated with the network interface.
	Ipv6Addresses []*NetworkInterfaceIpv6Address `locationName:"ipv6AddressesSet" locationNameList:"item" type:"list"`

	// The MAC address.
	MacAddress *string `locationName:"macAddress" type:"string"`

	// The ID of the network interface.
	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string"`

	// The AWS account ID of the owner of the network interface.
	OwnerId *string `locationName:"ownerId" type:"string"`

	// The private DNS name.
	PrivateDnsName *string `locationName:"privateDnsName" type:"string"`

	// The IPv4 address of the network interface within the subnet.
	PrivateIpAddress *string `locationName:"privateIpAddress" type:"string"`

	// The private IPv4 addresses associated with the network interface.
	PrivateIpAddresses []*NetworkInterfacePrivateIpAddress `locationName:"privateIpAddressesSet" locationNameList:"item" type:"list"`

	// The ID of the entity that launched the instance on your behalf (for example,
	// AWS Management Console or Auto Scaling).
	RequesterId *string `locationName:"requesterId" type:"string"`

	// Indicates whether the network interface is being managed by AWS.
	RequesterManaged *bool `locationName:"requesterManaged" type:"boolean"`

	// Indicates whether traffic to or from the instance is validated.
	SourceDestCheck *bool `locationName:"sourceDestCheck" type:"boolean"`

	// The status of the network interface.
	Status *string `locationName:"status" type:"string" enum:"NetworkInterfaceStatus"`

	// The ID of the subnet.
	SubnetId *string `locationName:"subnetId" type:"string"`

	// Any tags assigned to the network interface.
	TagSet []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`

	// The ID of the VPC.
	VpcId *string `locationName:"vpcId" type:"string"`
}

// String returns the string representation
func (s NetworkInterface) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s NetworkInterface) GoString() string {
	return s.String()
}

// SetAssociation sets the Association field's value.
func (s *NetworkInterface) SetAssociation(v *NetworkInterfaceAssociation) *NetworkInterface {
	s.Association = v
	return s
}

// SetAttachment sets the Attachment field's value.
func (s *NetworkInterface) SetAttachment(v *NetworkInterfaceAttachment) *NetworkInterface {
	s.Attachment = v
	return s
}

// SetAvailabilityZone sets the AvailabilityZone field's value.
func (s *NetworkInterface) SetAvailabilityZone(v string) *NetworkInterface {
	s.AvailabilityZone = &v
	return s
}

// SetDescription sets the Description field's value.
func (s *NetworkInterface) SetDescription(v string) *NetworkInterface {
	s.Description = &v
	return s
}

// SetGroups sets the Groups field's value.
func (s *NetworkInterface) SetGroups(v []*GroupIdentifier) *NetworkInterface {
	s.Groups = v
	return s
}

// SetInterfaceType sets the InterfaceType field's value.
func (s *NetworkInterface) SetInterfaceType(v string) *NetworkInterface {
	s.InterfaceType = &v
	return s
}

// SetIpv6Addresses sets the Ipv6Addresses field's value.
func (s *NetworkInterface) SetIpv6Addresses(v []*NetworkInterfaceIpv6Address) *NetworkInterface {
	s.Ipv6Addresses = v
	return s
}

// SetMacAddress sets the MacAddress field's value.
func (s *NetworkInterface) SetMacAddress(v string) *NetworkInterface {
	s.MacAddress = &v
	return s
}

// SetNetworkInterfaceId sets the NetworkInterfaceId field's value.
func (s *NetworkInterface) SetNetworkInterfaceId(v string) *NetworkInterface {
	s.NetworkInterfaceId = &v
	return s
}

// SetOwnerId sets the OwnerId field's value.
func (s *NetworkInterface) SetOwnerId(v string) *NetworkInterface {
	s.OwnerId = &v
	return s
}

// SetPrivateDnsName sets the PrivateDnsName field's value.
func (s *NetworkInterface) SetPrivateDnsName(v string) *NetworkInterface {
	s.PrivateDnsName = &v
	return s
}

// SetPrivateIpAddress sets the PrivateIpAddress field's value.
func (s *NetworkInterface) SetPrivateIpAddress(v string) *NetworkInterface {
	s.PrivateIpAddress = &v
	return s
}

// SetPrivateIpAddresses sets the PrivateIpAddresses field's value.
func (s *NetworkInterface) SetPrivateIpAddresses(v []*NetworkInterfacePrivateIpAddress) *NetworkInterface {
	s.PrivateIpAddresses = v
	return s
}

// SetRequesterId sets the RequesterId field's value.
func (s *NetworkInterface) SetRequesterId(v string) *NetworkInterface {
	s.RequesterId = &v
	return s
}

// SetRequesterManaged sets the RequesterManaged field's value.
func (s *NetworkInterface) SetRequesterManaged(v bool) *NetworkInterface {
	s.RequesterManaged = &v
	return s
}

// SetSourceDestCheck sets the SourceDestCheck field's value.
func (s *NetworkInterface) SetSourceDestCheck(v bool) *NetworkInterface {
	s.SourceDestCheck = &v
	return s
}

// SetStatus sets the Status field's value.
func (s *NetworkInterface) SetStatus(v string) *NetworkInterface {
	s.Status = &v
	return s
}

// SetSubnetId sets the SubnetId field's value.
func (s *NetworkInterface) SetSubnetId(v string) *NetworkInterface {
	s.SubnetId = &v
	return s
}

// SetTagSet sets the TagSet field's value.
func (s *NetworkInterface) SetTagSet(v []*Tag) *NetworkInterface {
	s.TagSet = v
	return s
}

// SetVpcId sets the VpcId field's value.
func (s *NetworkInterface) SetVpcId(v string) *NetworkInterface {
	s.VpcId = &v
	return s
}

// Describes association information for an Elastic IP address (IPv4 only).
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/NetworkInterfaceAssociation
type NetworkInterfaceAssociation struct {
	_ struct{} `type:"structure"`

	// The allocation ID.
	AllocationId *string `locationName:"allocationId" type:"string"`

	// The association ID.
	AssociationId *string `locationName:"associationId" type:"string"`

	// The ID of the Elastic IP address owner.
	IpOwnerId *string `locationName:"ipOwnerId" type:"string"`

	// The public DNS name.
	PublicDnsName *string `locationName:"publicDnsName" type:"string"`

	// The address of the Elastic IP address bound to the network interface.
	PublicIp *string `locationName:"publicIp" type:"string"`
}

// String returns the string representation
func (s NetworkInterfaceAssociation) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s NetworkInterfaceAssociation) GoString() string {
	return s.String()
}

// SetAllocationId sets the AllocationId field's value.
func (s *NetworkInterfaceAssociation) SetAllocationId(v string) *NetworkInterfaceAssociation {
	s.AllocationId = &v
	return s
}

// SetAssociationId sets the AssociationId field's value.
func (s *NetworkInterfaceAssociation) SetAssociationId(v string) *NetworkInterfaceAssociation {
	s.AssociationId = &v
	return s
}

// SetIpOwnerId sets the IpOwnerId field's value.
func (s *NetworkInterfaceAssociation) SetIpOwnerId(v string) *NetworkInterfaceAssociation {
	s.IpOwnerId = &v
	return s
}

// SetPublicDnsName sets the PublicDnsName field's value.
func (s *NetworkInterfaceAssociation) SetPublicDnsName(v string) *NetworkInterfaceAssociation {
	s.PublicDnsName = &v
	return s
}

// SetPublicIp sets the PublicIp field's value.
func (s *NetworkInterfaceAssociation) SetPublicIp(v string) *NetworkInterfaceAssociation {
	s.PublicIp = &v
	return s
}

// Describes a network interface attachment.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/NetworkInterfaceAttachment
type NetworkInterfaceAttachment struct {
	_ struct{} `type:"structure"`

	// The timestamp indicating when the attachment initiated.
	AttachTime *time.Time `locationName:"attachTime" type:"timestamp" timestampFormat:"iso8601"`

	// The ID of the network interface attachment.
	AttachmentId *string `locationName:"attachmentId" type:"string"`

	// Indicates whether the network interface is deleted when the instance is terminated.
	DeleteOnTermination *bool `locationName:"deleteOnTermination" type:"boolean"`

	// The device index of the network interface attachment on the instance.
	DeviceIndex *int64 `locationName:"deviceIndex" type:"integer"`

	// The ID of the instance.
	InstanceId *string `locationName:"instanceId" type:"string"`

	// The AWS account ID of the owner of the instance.
	InstanceOwnerId *string `locationName:"instanceOwnerId" type:"string"`

	// The attachment state.
	Status *string `locationName:"status" type:"string" enum:"AttachmentStatus"`
}

// String returns the string representation
func (s NetworkInterfaceAttachment) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s NetworkInterfaceAttachment) GoString() string {
	return s.String()
}

// SetAttachTime sets the AttachTime field's value.
func (s *NetworkInterfaceAttachment) SetAttachTime(v time.Time) *NetworkInterfaceAttachment {
	s.AttachTime = &v
	return s
}

// SetAttachmentId sets the AttachmentId field's value.
func (s *NetworkInterfaceAttachment) SetAttachmentId(v string) *NetworkInterfaceAttachment {
	s.AttachmentId = &v
	return s
}

// SetDeleteOnTermination sets the DeleteOnTermination field's value.
func (s *NetworkInterfaceAttachment) SetDeleteOnTermination(v bool) *NetworkInterfaceAttachment {
	s.DeleteOnTermination = &v
	return s
}

// SetDeviceIndex sets the DeviceIndex field's value.
func (s *NetworkInterfaceAttachment) SetDeviceIndex(v int64) *NetworkInterfaceAttachment {
	s.DeviceIndex = &v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *NetworkInterfaceAttachment) SetInstanceId(v string) *NetworkInterfaceAttachment {
	s.InstanceId = &v
	return s
}

// SetInstanceOwnerId sets the InstanceOwnerId field's value.
func (s *NetworkInterfaceAttachment) SetInstanceOwnerId(v string) *NetworkInterfaceAttachment {
	s.InstanceOwnerId = &v
	return s
}

// SetStatus sets the Status field's value.
func (s *NetworkInterfaceAttachment) SetStatus(v string) *NetworkInterfaceAttachment {
	s.Status = &v
	return s
}

// Describes an IPv6 address associated with a network interface.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/NetworkInterfaceIpv6Address
type NetworkInterfaceIpv6Address struct {
	_ struct{} `type:"structure"`

	// The IPv6 address.
	Ipv6Address *string `locationName:"ipv6Address" type:"string"`
}

// String returns the string representation
func (s NetworkInterfaceIpv6Address) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s NetworkInterfaceIpv6Address) GoString() string {
	return s.String()
}

// SetIpv6Address sets the Ipv6Address field's value.
func (s *NetworkInterfaceIpv6Address) SetIpv6Address(v string) *NetworkInterfaceIpv6Address {
	s.Ipv6Address = &v
	return s
}

// Describes the private IPv4 address of a network interface.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/NetworkInterfacePrivateIpAddress
type NetworkInterfacePrivateIpAddress struct {
	_ struct{} `type:"structure"`

	// The association information for an Elastic IP address (IPv4) associated with
	// the network interface.
	Association *NetworkInterfaceAssociation `locationName:"association" type:"structure"`

	// Indicates whether this IPv4 address is the primary private IPv4 address of
	// the network interface.
	Primary *bool `locationName:"primary" type:"boolean"`

	// The private DNS name.
	PrivateDnsName *string `locationName:"privateDnsName" type:"string"`

	// The private IPv4 address.
	PrivateIpAddress *string `locationName:"privateIpAddress" type:"string"`
}

// String returns the string representation
func (s NetworkInterfacePrivateIpAddress) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s NetworkInterfacePrivateIpAddress) GoString() string {
	return s.String()
}

// SetAssociation sets the Association field's value.
func (s *NetworkInterfacePrivateIpAddress) SetAssociation(v *NetworkInterfaceAssociation) *NetworkInterfacePrivateIpAddress {
	s.Association = v
	return s
}

// SetPrimary sets the Primary field's value.
func (s *NetworkInterfacePrivateIpAddress) SetPrimary(v bool) *NetworkInterfacePrivateIpAddress {
	s.Primary = &v
	return s
}

// SetPrivateDnsName sets the PrivateDnsName field's value.
func (s *NetworkInterfacePrivateIpAddress) SetPrivateDnsName(v string) *NetworkInterfacePrivateIpAddress {
	s.PrivateDnsName = &v
	return s
}

// SetPrivateIpAddress sets the PrivateIpAddress field's value.
func (s *NetworkInterfacePrivateIpAddress) SetPrivateIpAddress(v string) *NetworkInterfacePrivateIpAddress {
	s.PrivateIpAddress = &v
	return s
}

// Describes an IPv6 address.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/InstanceIpv6Address
type InstanceIpv6Address struct {
	_ struct{} `type:"structure"`

	// The IPv6 address.
	Ipv6Address *string `locationName:"ipv6Address" type:"string"`
}

// String returns the string representation
func (s InstanceIpv6Address) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstanceIpv6Address) GoString() string {
	return s.String()
}

// SetIpv6Address sets the Ipv6Address field's value.
func (s *InstanceIpv6Address) SetIpv6Address(v string) *InstanceIpv6Address {
	s.Ipv6Address = &v
	return s
}

// Describes an attachment change.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/NetworkInterfaceAttachmentChanges
type NetworkInterfaceAttachmentChanges struct {
	_ struct{} `type:"structure"`

	// The ID of the network interface attachment.
	AttachmentId *string `locationName:"attachmentId" type:"string"`

	// Indicates whether the network interface is deleted when the instance is terminated.
	DeleteOnTermination *bool `locationName:"deleteOnTermination" type:"boolean"`
}

// String returns the string representation
func (s NetworkInterfaceAttachmentChanges) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s NetworkInterfaceAttachmentChanges) GoString() string {
	return s.String()
}

// SetAttachmentId sets the AttachmentId field's value.
func (s *NetworkInterfaceAttachmentChanges) SetAttachmentId(v string) *NetworkInterfaceAttachmentChanges {
	s.AttachmentId = &v
	return s
}

// SetDeleteOnTermination sets the DeleteOnTermination field's value.
func (s *NetworkInterfaceAttachmentChanges) SetDeleteOnTermination(v bool) *NetworkInterfaceAttachmentChanges {
	s.DeleteOnTermination = &v
	return s
}

// String returns the string representation
func (s PrivateIpAddressSpecification) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s PrivateIpAddressSpecification) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *PrivateIpAddressSpecification) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "PrivateIpAddressSpecification"}
	if s.PrivateIpAddress == nil {
		invalidParams.Add(request.NewErrParamRequired("PrivateIpAddress"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetPrimary sets the Primary field's value.
func (s *PrivateIpAddressSpecification) SetPrimary(v bool) *PrivateIpAddressSpecification {
	s.Primary = &v
	return s
}

// SetPrivateIpAddress sets the PrivateIpAddress field's value.
func (s *PrivateIpAddressSpecification) SetPrivateIpAddress(v string) *PrivateIpAddressSpecification {
	s.PrivateIpAddress = &v
	return s
}

// Contains the parameters for DetachNetworkInterface.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DetachNetworkInterfaceRequest
type DetachNetworkInterfaceInput struct {
	_ struct{} `type:"structure"`

	// The ID of the attachment.
	//
	// AttachmentId is a required field
	AttachmentId *string `locationName:"attachmentId" type:"string" required:"true"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// Specifies whether to force a detachment.
	Force *bool `locationName:"force" type:"boolean"`
}

// String returns the string representation
func (s DetachNetworkInterfaceInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DetachNetworkInterfaceInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DetachNetworkInterfaceInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DetachNetworkInterfaceInput"}
	if s.AttachmentId == nil {
		invalidParams.Add(request.NewErrParamRequired("AttachmentId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAttachmentId sets the AttachmentId field's value.
func (s *DetachNetworkInterfaceInput) SetAttachmentId(v string) *DetachNetworkInterfaceInput {
	s.AttachmentId = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *DetachNetworkInterfaceInput) SetDryRun(v bool) *DetachNetworkInterfaceInput {
	s.DryRun = &v
	return s
}

// SetForce sets the Force field's value.
func (s *DetachNetworkInterfaceInput) SetForce(v bool) *DetachNetworkInterfaceInput {
	s.Force = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DetachNetworkInterfaceOutput
type DetachNetworkInterfaceOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DetachNetworkInterfaceOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DetachNetworkInterfaceOutput) GoString() string {
	return s.String()
}

// Contains the parameters for AttachNetworkInterface.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/AttachNetworkInterfaceRequest
type AttachNetworkInterfaceInput struct {
	_ struct{} `type:"structure"`

	// The index of the device for the network interface attachment.
	//
	// DeviceIndex is a required field
	DeviceIndex *int64 `locationName:"deviceIndex" type:"integer" required:"true"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the instance.
	//
	// InstanceId is a required field
	InstanceId *string `locationName:"instanceId" type:"string" required:"true"`

	// The ID of the network interface.
	//
	// NetworkInterfaceId is a required field
	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string" required:"true"`
}

// String returns the string representation
func (s AttachNetworkInterfaceInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AttachNetworkInterfaceInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *AttachNetworkInterfaceInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "AttachNetworkInterfaceInput"}
	if s.DeviceIndex == nil {
		invalidParams.Add(request.NewErrParamRequired("DeviceIndex"))
	}
	if s.InstanceId == nil {
		invalidParams.Add(request.NewErrParamRequired("InstanceId"))
	}
	if s.NetworkInterfaceId == nil {
		invalidParams.Add(request.NewErrParamRequired("NetworkInterfaceId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDeviceIndex sets the DeviceIndex field's value.
func (s *AttachNetworkInterfaceInput) SetDeviceIndex(v int64) *AttachNetworkInterfaceInput {
	s.DeviceIndex = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *AttachNetworkInterfaceInput) SetDryRun(v bool) *AttachNetworkInterfaceInput {
	s.DryRun = &v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *AttachNetworkInterfaceInput) SetInstanceId(v string) *AttachNetworkInterfaceInput {
	s.InstanceId = &v
	return s
}

// SetNetworkInterfaceId sets the NetworkInterfaceId field's value.
func (s *AttachNetworkInterfaceInput) SetNetworkInterfaceId(v string) *AttachNetworkInterfaceInput {
	s.NetworkInterfaceId = &v
	return s
}

// Contains the output of AttachNetworkInterface.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/AttachNetworkInterfaceResult
type AttachNetworkInterfaceOutput struct {
	_ struct{} `type:"structure"`

	// The ID of the network interface attachment.
	AttachmentId *string `locationName:"attachmentId" type:"string"`
}

// String returns the string representation
func (s AttachNetworkInterfaceOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AttachNetworkInterfaceOutput) GoString() string {
	return s.String()
}

// SetAttachmentId sets the AttachmentId field's value.
func (s *AttachNetworkInterfaceOutput) SetAttachmentId(v string) *AttachNetworkInterfaceOutput {
	s.AttachmentId = &v
	return s
}

// Contains the parameters for AssignPrivateIpAddresses.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/AssignPrivateIpAddressesRequest
type AssignPrivateIpAddressesInput struct {
	_ struct{} `type:"structure"`

	// Indicates whether to allow an IP address that is already assigned to another
	// network interface or instance to be reassigned to the specified network interface.
	AllowReassignment *bool `locationName:"allowReassignment" type:"boolean"`

	// The ID of the network interface.
	//
	// NetworkInterfaceId is a required field
	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string" required:"true"`

	// One or more IP addresses to be assigned as a secondary private IP address
	// to the network interface. You can't specify this parameter when also specifying
	// a number of secondary IP addresses.
	//
	// If you don't specify an IP address, Amazon EC2 automatically selects an IP
	// address within the subnet range.
	PrivateIpAddresses []*string `locationName:"privateIpAddress" locationNameList:"PrivateIpAddress" type:"list"`

	// The number of secondary IP addresses to assign to the network interface.
	// You can't specify this parameter when also specifying private IP addresses.
	SecondaryPrivateIpAddressCount *int64 `locationName:"secondaryPrivateIpAddressCount" type:"integer"`
}

// String returns the string representation
func (s AssignPrivateIpAddressesInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AssignPrivateIpAddressesInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *AssignPrivateIpAddressesInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "AssignPrivateIpAddressesInput"}
	if s.NetworkInterfaceId == nil {
		invalidParams.Add(request.NewErrParamRequired("NetworkInterfaceId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAllowReassignment sets the AllowReassignment field's value.
func (s *AssignPrivateIpAddressesInput) SetAllowReassignment(v bool) *AssignPrivateIpAddressesInput {
	s.AllowReassignment = &v
	return s
}

// SetNetworkInterfaceId sets the NetworkInterfaceId field's value.
func (s *AssignPrivateIpAddressesInput) SetNetworkInterfaceId(v string) *AssignPrivateIpAddressesInput {
	s.NetworkInterfaceId = &v
	return s
}

// SetPrivateIpAddresses sets the PrivateIpAddresses field's value.
func (s *AssignPrivateIpAddressesInput) SetPrivateIpAddresses(v []*string) *AssignPrivateIpAddressesInput {
	s.PrivateIpAddresses = v
	return s
}

// SetSecondaryPrivateIpAddressCount sets the SecondaryPrivateIpAddressCount field's value.
func (s *AssignPrivateIpAddressesInput) SetSecondaryPrivateIpAddressCount(v int64) *AssignPrivateIpAddressesInput {
	s.SecondaryPrivateIpAddressCount = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/AssignPrivateIpAddressesOutput
type AssignPrivateIpAddressesOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s AssignPrivateIpAddressesOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AssignPrivateIpAddressesOutput) GoString() string {
	return s.String()
}

// Contains the parameters for UnassignPrivateIpAddresses.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/UnassignPrivateIpAddressesRequest
type UnassignPrivateIpAddressesInput struct {
	_ struct{} `type:"structure"`

	// The ID of the network interface.
	//
	// NetworkInterfaceId is a required field
	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string" required:"true"`

	// The secondary private IP addresses to unassign from the network interface.
	// You can specify this option multiple times to unassign more than one IP address.
	//
	// PrivateIpAddresses is a required field
	PrivateIpAddresses []*string `locationName:"privateIpAddress" locationNameList:"PrivateIpAddress" type:"list" required:"true"`
}

type PurchaseReservedInstancesOfferingInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The number of Reserved Instances to purchase.
	//
	// InstanceCount is a required field
	InstanceCount *int64 `type:"integer" required:"true"`

	// Specified for Reserved Instance Marketplace offerings to limit the total
	// order and ensure that the Reserved Instances are not purchased at unexpected
	// prices.
	LimitPrice *ReservedInstanceLimitPrice `locationName:"limitPrice" type:"structure"`

	// The ID of the Reserved Instance offering to purchase.
	//
	// ReservedInstancesOfferingId is a required field
	ReservedInstancesOfferingId *string `type:"string" required:"true"`
}

type PurchaseReservedInstancesOfferingOutput struct {
	_ struct{} `type:"structure"`

	// The IDs of the purchased Reserved Instances.
	ReservedInstancesId *string `locationName:"reservedInstancesId" type:"string"`
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/ReservedInstanceLimitPrice
type ReservedInstanceLimitPrice struct {
	_ struct{} `type:"structure"`

	// Used for Reserved Instance Marketplace offerings. Specifies the limit price
	// on the total order (instanceCount * price).
	Amount *float64 `locationName:"amount" type:"double"`

	// The currency in which the limitPrice amount is specified. At this time, the
	// only supported currency is USD.
	CurrencyCode *string `locationName:"currencyCode" type:"string" enum:"CurrencyCodeValues"`
}

type UnassignPrivateIpAddressesOutput struct {
	_ struct{} `type:"structure"`
}

// Contains the parameters for CreateVpcEndpoint.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateVpcEndpointRequest
type CreateVpcEndpointInput struct {
	_ struct{} `type:"structure"`

	// Unique, case-sensitive identifier you provide to ensure the idempotency of
	// the request. For more information, see How to Ensure Idempotency (http://docs.aws.amazon.com/AWSEC2/latest/APIReference/Run_Instance_Idempotency.html).
	ClientToken *string `type:"string"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `type:"boolean"`

	// A policy to attach to the endpoint that controls access to the service. The
	// policy must be in valid JSON format. If this parameter is not specified,
	// we attach a default policy that allows full access to the service.
	PolicyDocument *string `type:"string"`

	// One or more route table IDs.
	RouteTableIds []*string `locationName:"RouteTableId" locationNameList:"item" type:"list"`

	// The AWS service name, in the form com.amazonaws.region.service. To get a
	// list of available services, use the DescribeVpcEndpointServices request.
	//
	// ServiceName is a required field
	ServiceName *string `type:"string" required:"true"`

	// The ID of the VPC in which the endpoint will be used.
	//
	// VpcId is a required field
	VpcId *string `type:"string" required:"true"`
}

// Contains the output of CreateVpcEndpoint.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateVpcEndpointResult
type CreateVpcEndpointOutput struct {
	_ struct{} `type:"structure"`

	// Unique, case-sensitive identifier you provide to ensure the idempotency of
	// the request.
	ClientToken *string `locationName:"clientToken" type:"string"`

	// Information about the endpoint.
	VpcEndpoint *VpcEndpoint `locationName:"vpcEndpoint" type:"structure"`
}

// Describes a VPC endpoint.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/VpcEndpoint
type VpcEndpoint struct {
	_ struct{} `type:"structure"`

	// The date and time the VPC endpoint was created.
	CreationTimestamp *time.Time `locationName:"creationTimestamp" type:"timestamp" timestampFormat:"iso8601"`

	// The policy document associated with the endpoint.
	PolicyDocument *string `locationName:"policyDocument" type:"string"`

	// One or more route tables associated with the endpoint.
	RouteTableIds []*string `locationName:"routeTableIdSet" locationNameList:"item" type:"list"`

	// The name of the AWS service to which the endpoint is associated.
	ServiceName *string `locationName:"serviceName" type:"string"`

	// The state of the VPC endpoint.
	State *string `locationName:"state" type:"string" enum:"State"`

	// The ID of the VPC endpoint.
	VpcEndpointId *string `locationName:"vpcEndpointId" type:"string"`

	// The ID of the VPC to which the endpoint is associated.
	VpcId *string `locationName:"vpcId" type:"string"`
}

// Contains the parameters for DescribeVpcEndpoints.
type DescribeVpcEndpointsInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `type:"boolean"`

	// One or more filters.
	//
	//    * service-name: The name of the service.
	//
	//    * vpc-id: The ID of the VPC in which the endpoint resides.
	//
	//    * vpc-endpoint-id: The ID of the endpoint.
	//
	//    * vpc-endpoint-state: The state of the endpoint. (pending | available
	//    | deleting | deleted)
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// The maximum number of items to return for this request. The request returns
	// a token that you can specify in a subsequent call to get the next set of
	// results.
	//
	// Constraint: If the value is greater than 1000, we return only 1000 items.
	MaxResults *int64 `type:"integer"`

	// The token for the next set of items to return. (You received this token from
	// a prior call.)
	NextToken *string `type:"string"`

	// One or more endpoint IDs.
	VpcEndpointIds []*string `locationName:"VpcEndpointId" locationNameList:"item" type:"list"`
}

// Contains the output of DescribeVpcEndpoints.
type DescribeVpcEndpointsOutput struct {
	_ struct{} `type:"structure"`

	// The token to use when requesting the next set of items. If there are no additional
	// items to return, the string is empty.
	NextToken *string `locationName:"nextToken" type:"string"`

	// Information about the endpoints.
	VpcEndpoints []*VpcEndpoint `locationName:"vpcEndpointSet" locationNameList:"item" type:"list"`

	RequestId *string `locationName:"requestId" type:"string"`
}

type ModifyVpcEndpointInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun              *bool     `type:"boolean"`
	PolicyDocument      *string   `type:"string"`                                                         // A policy document to attach to the endpoint. The policy must be in valid JSON format.
	RemoveRouteTableIds []*string `locationName:"RemoveRouteTableId" locationNameList:"item" type:"list"` // One or more route table IDs to disassociate from the endpoint.
	ResetPolicy         *bool     `type:"boolean"`                                                        // Specify true to reset the policy document to the default policy. The default policy allows access to the service.
	VpcEndpointId       *string   `type:"string" required:"true"`                                         //The ID of the endpoint. VpcEndpointId is a required field
	AddRouteTableIds    []*string `locationName:"AddRouteTableId" locationNameList:"item" type:"list"`    // One or more route tables IDs to associate with the endpoint.
}

// Contains the output of ModifyVpcEndpoint.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/ModifyVpcEndpointResult
type ModifyVpcEndpointOutput struct {
	_      struct{} `type:"structure"`
	Return *bool    `locationName:"return" type:"boolean"` // Returns true if the request succeeds; otherwise, it returns an error.
}

type DeleteVpcEndpointsInput struct {
	_              struct{}  `type:"structure"`
	VpcEndpointIds []*string `locationName:"VpcEndpointId" locationNameList:"item" type:"list" required:"true"` // One or more endpoint IDs. VpcEndpointIds is a required field
	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `type:"boolean"`
}

// Contains the output of DeleteVpcEndpoints.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteVpcEndpointsResult
type DeleteVpcEndpointsOutput struct {
	_ struct{} `type:"structure"`
}

type ModifySnapshotAttributeInput struct {
	_ struct{} `type:"structure"`

	// The snapshot attribute to modify.
	//
	// Only volume creation permissions may be modified at the customer level.
	Attribute *string `type:"string" enum:"SnapshotAttributeName"`

	// A JSON representation of the snapshot attribute modification.
	CreateVolumePermission *CreateVolumePermissionModifications `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The group to modify for the snapshot.
	GroupNames []*string `locationName:"UserGroup" locationNameList:"GroupName" type:"list"`

	// The type of operation to perform to the attribute.
	OperationType *string `type:"string" enum:"OperationType"`

	// The ID of the snapshot.
	//
	// SnapshotId is a required field
	SnapshotId *string `type:"string" required:"true"`

	// The account ID to modify for the snapshot.
	UserIds []*string `locationName:"UserId" locationNameList:"UserId" type:"list"`
}

type CreateVolumePermissionModifications struct {
	_ struct{} `type:"structure"`

	// Adds a specific AWS account ID or group to a volume's list of create volume
	// permissions.
	Add []*CreateVolumePermission `locationNameList:"item" type:"list"`

	// Removes a specific AWS account ID or group from a volume's list of create
	// volume permissions.
	Remove []*CreateVolumePermission `locationNameList:"item" type:"list"`
}

type CreateVolumePermission struct {
	_ struct{} `type:"structure"`

	// The specific group that is to be added or removed from a volume's list of
	// create volume permissions.
	Group *string `locationName:"group" type:"string" enum:"PermissionGroup"`

	// The specific AWS account ID that is to be added or removed from a volume's
	// list of create volume permissions.
	UserId *string `locationName:"userId" type:"string"`
}

type ModifySnapshotAttributeOutput struct {
	_ struct{} `type:"structure"`
}

type DescribeSnapshotAttributeInput struct {
	_ struct{} `type:"structure"`

	// The snapshot attribute you would like to view.
	//
	// Attribute is a required field
	Attribute *string `type:"string" required:"true" enum:"SnapshotAttributeName"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the EBS snapshot.
	//
	// SnapshotId is a required field
	SnapshotId *string `type:"string" required:"true"`
}

type DescribeSnapshotAttributeOutput struct {
	_ struct{} `type:"structure"`

	// A list of permissions for creating volumes from the snapshot.
	CreateVolumePermissions []*CreateVolumePermission `locationName:"createVolumePermission" locationNameList:"item" type:"list"`

	// A list of product codes.
	ProductCodes []*ProductCode `locationName:"productCodes" locationNameList:"item" type:"list"`

	// The ID of the EBS snapshot.
	SnapshotId *string `locationName:"snapshotId" type:"string"`

	RequestId *string `locationName:"requestId" type:"string"`
}

type ImportSnapshotInput struct {
	_ struct{} `type:"structure"`

	// The client-specific data.
	ClientData *ClientData `type:"structure"`

	// Token to enable idempotency for VM import requests.
	ClientToken *string `type:"string"`

	// The description string for the import snapshot task.
	Description *string `type:"string"`

	// Information about the disk container.
	DiskContainer *SnapshotDiskContainer `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `type:"boolean"`

	// The name of the role to use when not using the default role, 'vmimport'.
	RoleName         *string `type:"string"`
	SnapshotLocation *string `type:"string"`
	SnapshotSize     *string `type:"string"`
}

type ImportSnapshotOutput struct {
	_ struct{} `type:"structure"`

	// A description of the import snapshot task.
	Description *string `locationName:"description" type:"string"`

	// The ID of the import snapshot task.
	ImportTaskId *string `locationName:"importTaskId" type:"string"`

	// Information about the import snapshot task.
	SnapshotTaskDetail *SnapshotTaskDetail `locationName:"snapshotTaskDetail" type:"structure"`

	Id *string `locationName:"id" type:"string"`
}

type ClientData struct {
	_ struct{} `type:"structure"`

	// A user-defined comment about the disk upload.
	Comment *string `type:"string"`

	// The time that the disk upload ends.
	UploadEnd *time.Time `type:"timestamp" timestampFormat:"iso8601"`

	// The size of the uploaded disk image, in GiB.
	UploadSize *float64 `type:"double"`

	// The time that the disk upload starts.
	UploadStart *time.Time `type:"timestamp" timestampFormat:"iso8601"`
}
