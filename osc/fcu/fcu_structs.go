package fcu

import (
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

	RequesterId *string `locationName:"requestId" type:"string"`

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

	RequesterId *string `locationName:"requestId" type:"string"`

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

	RequesterId *string `type:"string"`

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
	KeyPairs []*KeyPairInfo `locationName:"keySet" locationNameList:"item" type:"list"`

	RequesterId *string `locationName:"requestId" type:"string"`
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
	_ struct{} `type:"structure"`

	Description *string `locationName:"groupDescription" type:"string"`

	GroupId *string `locationName:"groupId" type:"string"`

	GroupName *string `locationName:"groupName" type:"string"`

	IpPermissions []*IpPermission `locationName:"ipPermissions" locationNameList:"item" type:"list"`

	IpPermissionsEgress []*IpPermission `locationName:"ipPermissionsEgress" locationNameList:"item" type:"list"`

	OwnerId *string `locationName:"ownerId" type:"string"`

	Tags []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`

	VpcId *string `locationName:"vpcId" type:"string"`
}

type IpPermission struct {
	_ struct{} `type:"structure"`

	FromPort *int64 `locationName:"fromPort" type:"integer"`

	IpProtocol *string `locationName:"ipProtocol" type:"string"`

	IpRanges []*IpRange `locationName:"ipRanges" locationNameList:"item" type:"list"`

	Ipv6Ranges []*Ipv6Range `locationName:"ipv6Ranges" locationNameList:"item" type:"list"`

	PrefixListIds []*PrefixListId `locationName:"prefixListIds" locationNameList:"item" type:"list"`

	ToPort *int64 `locationName:"toPort" type:"integer"`

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

	RequesterId *string `locationName:"requestId" type:"string"`
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
	s.RequesterId = &v
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

	RequesterId *string `locationName:"requestId" type:"string"`
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

	RequesterId *string `locationName:"requestId" type:"string"`
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
	s.RequesterId = v
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

	RequesterId *string `locationName:"requestId" type:"string"`
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
	RequesterId      *string            `locationName:"requestId" type:"string"`
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

// String returns the string representation
func (s CreateNatGatewayInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateNatGatewayInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CreateNatGatewayInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "CreateNatGatewayInput"}
	if s.AllocationId == nil {
		invalidParams.Add(request.NewErrParamRequired("AllocationId"))
	}
	if s.SubnetId == nil {
		invalidParams.Add(request.NewErrParamRequired("SubnetId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAllocationId sets the AllocationId field's value.
func (s *CreateNatGatewayInput) SetAllocationId(v string) *CreateNatGatewayInput {
	s.AllocationId = &v
	return s
}

// SetClientToken sets the ClientToken field's value.
func (s *CreateNatGatewayInput) SetClientToken(v string) *CreateNatGatewayInput {
	s.ClientToken = &v
	return s
}

// SetSubnetId sets the SubnetId field's value.
func (s *CreateNatGatewayInput) SetSubnetId(v string) *CreateNatGatewayInput {
	s.SubnetId = &v
	return s
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

// String returns the string representation
func (s DetachInternetGatewayInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DetachInternetGatewayInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DetachInternetGatewayInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DetachInternetGatewayInput"}
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
func (s *DetachInternetGatewayInput) SetDryRun(v bool) *DetachInternetGatewayInput {
	s.DryRun = &v
	return s
}

// SetInternetGatewayId sets the InternetGatewayId field's value.
func (s *DetachInternetGatewayInput) SetInternetGatewayId(v string) *DetachInternetGatewayInput {
	s.InternetGatewayId = &v
	return s
}

// SetVpcId sets the VpcId field's value.
func (s *DetachInternetGatewayInput) SetVpcId(v string) *DetachInternetGatewayInput {
	s.VpcId = &v
	return s
}

type DetachInternetGatewayOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DetachInternetGatewayOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DetachInternetGatewayOutput) GoString() string {
	return s.String()
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
