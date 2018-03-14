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

func (s InstanceStatusSummary) String() string {
	return awsutil.Prettify(s)
}

func (s InstanceStatusSummary) GoString() string {
	return s.String()
}

func (s *InstanceStatusSummary) SetDetails(v []*InstanceStatusDetails) *InstanceStatusSummary {
	s.Details = v
	return s
}

func (s *InstanceStatusSummary) SetStatus(v string) *InstanceStatusSummary {
	s.Status = &v
	return s
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

func (s AllocateAddressInput) String() string {
	return awsutil.Prettify(s)
}

func (s AllocateAddressInput) GoString() string {
	return s.String()
}

func (s *AllocateAddressInput) SetDomain(v string) *AllocateAddressInput {
	s.Domain = &v
	return s
}

func (s *AllocateAddressInput) SetDryRun(v bool) *AllocateAddressInput {
	s.DryRun = &v
	return s
}

type AllocateAddressOutput struct {
	_ struct{} `type:"structure"`

	AllocationId *string `locationName:"allocationId" type:"string"`

	Domain *string `locationName:"domain" type:"string" enum:"DomainType"`

	PublicIp *string `locationName:"publicIp" type:"string"`
}

func (s AllocateAddressOutput) String() string {
	return awsutil.Prettify(s)
}

func (s AllocateAddressOutput) GoString() string {
	return s.String()
}

func (s *AllocateAddressOutput) SetAllocationId(v string) *AllocateAddressOutput {
	s.AllocationId = &v
	return s
}

func (s *AllocateAddressOutput) SetDomain(v string) *AllocateAddressOutput {
	s.Domain = &v
	return s
}

func (s *AllocateAddressOutput) SetPublicIp(v string) *AllocateAddressOutput {
	s.PublicIp = &v
	return s
}

type DescribeAddressesInput struct {
	_ struct{} `type:"structure"`

	AllocationIds []*string `locationName:"AllocationId" locationNameList:"AllocationId" type:"list"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	PublicIps []*string `locationName:"PublicIp" locationNameList:"PublicIp" type:"list"`
}

func (s DescribeAddressesInput) String() string {
	return awsutil.Prettify(s)
}

func (s DescribeAddressesInput) GoString() string {
	return s.String()
}

func (s *DescribeAddressesInput) SetAllocationIds(v []*string) *DescribeAddressesInput {
	s.AllocationIds = v
	return s
}

func (s *DescribeAddressesInput) SetDryRun(v bool) *DescribeAddressesInput {
	s.DryRun = &v
	return s
}

func (s *DescribeAddressesInput) SetFilters(v []*Filter) *DescribeAddressesInput {
	s.Filters = v
	return s
}

func (s *DescribeAddressesInput) SetPublicIps(v []*string) *DescribeAddressesInput {
	s.PublicIps = v
	return s
}

type DescribeAddressesOutput struct {
	_ struct{} `type:"structure"`

	Addresses []*Address `locationName:"addressesSet" locationNameList:"item" type:"list"`
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

func (s Address) String() string {
	return awsutil.Prettify(s)
}

func (s Address) GoString() string {
	return s.String()
}

func (s *Address) SetAllocationId(v string) *Address {
	s.AllocationId = &v
	return s
}

func (s *Address) SetAllowReassociation(v bool) *Address {
	s.AllowReassociation = &v
	return s
}

func (s *Address) SetAssociationId(v string) *Address {
	s.AssociationId = &v
	return s
}

func (s *Address) SetDomain(v string) *Address {
	s.Domain = &v
	return s
}

func (s *Address) SetInstanceId(v string) *Address {
	s.InstanceId = &v
	return s
}

func (s *Address) SetNetworkInterfaceId(v string) *Address {
	s.NetworkInterfaceId = &v
	return s
}

func (s *Address) SetNetworkInterfaceOwnerId(v string) *Address {
	s.NetworkInterfaceOwnerId = &v
	return s
}

func (s *Address) SetPrivateIpAddress(v string) *Address {
	s.PrivateIpAddress = &v
	return s
}

func (s *Address) SetPublicIp(v string) *Address {
	s.PublicIp = &v
	return s
}

type ModifyInstanceAttributeInput struct {
	_ struct{} `type:"structure"`

	Attribute *string `locationName:"attribute" type:"string" enum:"InstanceAttributeName"`

	BlockDeviceMappings []*InstanceBlockDeviceMappingSpecification `locationName:"blockDeviceMapping" locationNameList:"item" type:"list"`

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

func (s *ModifyInstanceAttributeInput) SetBlockDeviceMappings(v []*InstanceBlockDeviceMappingSpecification) *ModifyInstanceAttributeInput {
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

func (s BlobAttributeValue) String() string {
	return awsutil.Prettify(s)
}

func (s BlobAttributeValue) GoString() string {
	return s.String()
}

func (s *BlobAttributeValue) SetValue(v []byte) *BlobAttributeValue {
	s.Value = v
	return s
}

type StopInstancesInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	Force *bool `locationName:"force" type:"boolean"`

	InstanceIds []*string `locationName:"InstanceId" locationNameList:"InstanceId" type:"list" required:"true"`
}

func (s StopInstancesInput) String() string {
	return awsutil.Prettify(s)
}

func (s StopInstancesInput) GoString() string {
	return s.String()
}

func (s *StopInstancesInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "StopInstancesInput"}
	if s.InstanceIds == nil {
		invalidParams.Add(request.NewErrParamRequired("InstanceIds"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

func (s *StopInstancesInput) SetDryRun(v bool) *StopInstancesInput {
	s.DryRun = &v
	return s
}

func (s *StopInstancesInput) SetForce(v bool) *StopInstancesInput {
	s.Force = &v
	return s
}

func (s *StopInstancesInput) SetInstanceIds(v []*string) *StopInstancesInput {
	s.InstanceIds = v
	return s
}

type StopInstancesOutput struct {
	_ struct{} `type:"structure"`

	StoppingInstances []*InstanceStateChange `locationName:"instancesSet" locationNameList:"item" type:"list"`
}

func (s StopInstancesOutput) String() string {
	return awsutil.Prettify(s)
}

func (s StopInstancesOutput) GoString() string {
	return s.String()
}

func (s *StopInstancesOutput) SetStoppingInstances(v []*InstanceStateChange) *StopInstancesOutput {
	s.StoppingInstances = v
	return s
}

type ModifyInstanceAttributeOutput struct {
	_ struct{} `type:"structure"`
}

func (s ModifyInstanceAttributeOutput) String() string {
	return awsutil.Prettify(s)
}

func (s ModifyInstanceAttributeOutput) GoString() string {
	return s.String()
}

type StartInstancesInput struct {
	_ struct{} `type:"structure"`

	AdditionalInfo *string `locationName:"additionalInfo" type:"string"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	InstanceIds []*string `locationName:"InstanceId" locationNameList:"InstanceId" type:"list" required:"true"`
}

func (s StartInstancesInput) String() string {
	return awsutil.Prettify(s)
}

func (s StartInstancesInput) GoString() string {
	return s.String()
}

func (s *StartInstancesInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "StartInstancesInput"}
	if s.InstanceIds == nil {
		invalidParams.Add(request.NewErrParamRequired("InstanceIds"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

func (s *StartInstancesInput) SetAdditionalInfo(v string) *StartInstancesInput {
	s.AdditionalInfo = &v
	return s
}

func (s *StartInstancesInput) SetDryRun(v bool) *StartInstancesInput {
	s.DryRun = &v
	return s
}

func (s *StartInstancesInput) SetInstanceIds(v []*string) *StartInstancesInput {
	s.InstanceIds = v
	return s
}

type StartInstancesOutput struct {
	_ struct{} `type:"structure"`

	StartingInstances []*InstanceStateChange `locationName:"instancesSet" locationNameList:"item" type:"list"`
}

func (s StartInstancesOutput) String() string {
	return awsutil.Prettify(s)
}

func (s StartInstancesOutput) GoString() string {
	return s.String()
}

func (s *StartInstancesOutput) SetStartingInstances(v []*InstanceStateChange) *StartInstancesOutput {
	s.StartingInstances = v
	return s
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

func (s AssociateAddressInput) String() string {
	return awsutil.Prettify(s)
}

func (s AssociateAddressInput) GoString() string {
	return s.String()
}

func (s *AssociateAddressInput) SetAllocationId(v string) *AssociateAddressInput {
	s.AllocationId = &v
	return s
}

func (s *AssociateAddressInput) SetAllowReassociation(v bool) *AssociateAddressInput {
	s.AllowReassociation = &v
	return s
}

func (s *AssociateAddressInput) SetInstanceId(v string) *AssociateAddressInput {
	s.InstanceId = &v
	return s
}

func (s *AssociateAddressInput) SetNetworkInterfaceId(v string) *AssociateAddressInput {
	s.NetworkInterfaceId = &v
	return s
}

func (s *AssociateAddressInput) SetPrivateIpAddress(v string) *AssociateAddressInput {
	s.PrivateIpAddress = &v
	return s
}

func (s *AssociateAddressInput) SetPublicIp(v string) *AssociateAddressInput {
	s.PublicIp = &v
	return s
}

type AssociateAddressOutput struct {
	_ struct{} `type:"structure"`

	AssociationId *string `locationName:"associationId" type:"string"`

	RequestId *string `locationName:"requestId" type:"string"`
}

func (s AssociateAddressOutput) String() string {
	return awsutil.Prettify(s)
}

func (s AssociateAddressOutput) GoString() string {
	return s.String()
}

func (s *AssociateAddressOutput) SetAssociationId(v string) *AssociateAddressOutput {
	s.AssociationId = &v
	return s
}

func (s *AssociateAddressOutput) SetRequestId(v string) *AssociateAddressOutput {
	s.RequestId = &v
	return s
}

type DisassociateAddressInput struct {
	_ struct{} `type:"structure"`

	AssociationId *string `type:"string"`

	PublicIp *string `type:"string"`
}

func (s DisassociateAddressInput) String() string {
	return awsutil.Prettify(s)
}

func (s DisassociateAddressInput) GoString() string {
	return s.String()
}

func (s *DisassociateAddressInput) SetAssociationId(v string) *DisassociateAddressInput {
	s.AssociationId = &v
	return s
}

func (s *DisassociateAddressInput) SetPublicIp(v string) *DisassociateAddressInput {
	s.PublicIp = &v
	return s
}

type DisassociateAddressOutput struct {
	_ struct{} `type:"structure"`

	RequestId *string `locationName:"requestId" type:"string"`
	Return    *bool   `locationName:"return" type:"boolean"`
}

func (s DisassociateAddressOutput) String() string {
	return awsutil.Prettify(s)
}

func (s DisassociateAddressOutput) GoString() string {
	return s.String()
}

func (s *DisassociateAddressOutput) SetReturn(v bool) *DisassociateAddressOutput {
	s.Return = &v
	return s
}

func (s *DisassociateAddressOutput) SetRequestId(v string) *DisassociateAddressOutput {
	s.RequestId = &v
	return s
}

type ReleaseAddressInput struct {
	_ struct{} `type:"structure"`

	AllocationId *string `type:"string"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	PublicIp *string `type:"string"`
}

func (s ReleaseAddressInput) String() string {
	return awsutil.Prettify(s)
}

func (s ReleaseAddressInput) GoString() string {
	return s.String()
}

func (s *ReleaseAddressInput) SetAllocationId(v string) *ReleaseAddressInput {
	s.AllocationId = &v
	return s
}

func (s *ReleaseAddressInput) SetDryRun(v bool) *ReleaseAddressInput {
	s.DryRun = &v
	return s
}

func (s *ReleaseAddressInput) SetPublicIp(v string) *ReleaseAddressInput {
	s.PublicIp = &v
	return s
}

type ReleaseAddressOutput struct {
	_ struct{} `type:"structure"`
}

func (s ReleaseAddressOutput) String() string {
	return awsutil.Prettify(s)
}

func (s ReleaseAddressOutput) GoString() string {
	return s.String()
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

func (s RegisterImageInput) String() string {
	return awsutil.Prettify(s)
}

func (s RegisterImageInput) GoString() string {
	return s.String()
}

func (s *RegisterImageInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "RegisterImageInput"}
	if s.Name == nil {
		invalidParams.Add(request.NewErrParamRequired("Name"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

func (s *RegisterImageInput) SetArchitecture(v string) *RegisterImageInput {
	s.Architecture = &v
	return s
}

func (s *RegisterImageInput) SetBillingProducts(v []*string) *RegisterImageInput {
	s.BillingProducts = v
	return s
}

func (s *RegisterImageInput) SetBlockDeviceMappings(v []*BlockDeviceMapping) *RegisterImageInput {
	s.BlockDeviceMappings = v
	return s
}

func (s *RegisterImageInput) SetDescription(v string) *RegisterImageInput {
	s.Description = &v
	return s
}

func (s *RegisterImageInput) SetDryRun(v bool) *RegisterImageInput {
	s.DryRun = &v
	return s
}

func (s *RegisterImageInput) SetEnaSupport(v bool) *RegisterImageInput {
	s.EnaSupport = &v
	return s
}

func (s *RegisterImageInput) SetImageLocation(v string) *RegisterImageInput {
	s.ImageLocation = &v
	return s
}

func (s *RegisterImageInput) SetKernelId(v string) *RegisterImageInput {
	s.KernelId = &v
	return s
}

func (s *RegisterImageInput) SetName(v string) *RegisterImageInput {
	s.Name = &v
	return s
}

func (s *RegisterImageInput) SetRamdiskId(v string) *RegisterImageInput {
	s.RamdiskId = &v
	return s
}

func (s *RegisterImageInput) SetRootDeviceName(v string) *RegisterImageInput {
	s.RootDeviceName = &v
	return s
}

func (s *RegisterImageInput) SetSriovNetSupport(v string) *RegisterImageInput {
	s.SriovNetSupport = &v
	return s
}

func (s *RegisterImageInput) SetVirtualizationType(v string) *RegisterImageInput {
	s.VirtualizationType = &v
	return s
}

type RegisterImageOutput struct {
	_ struct{} `type:"structure"`

	ImageId *string `locationName:"imageId" type:"string"`
}

func (s RegisterImageOutput) String() string {
	return awsutil.Prettify(s)
}

func (s RegisterImageOutput) GoString() string {
	return s.String()
}

func (s *RegisterImageOutput) SetImageId(v string) *RegisterImageOutput {
	s.ImageId = &v
	return s
}

type DeregisterImageInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	ImageId *string `type:"string" required:"true"`
}

func (s DeregisterImageInput) String() string {
	return awsutil.Prettify(s)
}

func (s DeregisterImageInput) GoString() string {
	return s.String()
}

func (s *DeregisterImageInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeregisterImageInput"}
	if s.ImageId == nil {
		invalidParams.Add(request.NewErrParamRequired("ImageId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

func (s *DeregisterImageInput) SetDryRun(v bool) *DeregisterImageInput {
	s.DryRun = &v
	return s
}

func (s *DeregisterImageInput) SetImageId(v string) *DeregisterImageInput {
	s.ImageId = &v
	return s
}

type DeregisterImageOutput struct {
	_ struct{} `type:"structure"`
}

func (s DeregisterImageOutput) String() string {
	return awsutil.Prettify(s)
}

func (s DeregisterImageOutput) GoString() string {
	return s.String()
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

func (s Image) String() string {
	return awsutil.Prettify(s)
}

func (s Image) GoString() string {
	return s.String()
}

func (s *Image) SetArchitecture(v string) *Image {
	s.Architecture = &v
	return s
}

type DescribeImagesInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	ExecutableUsers []*string `locationName:"ExecutableBy" locationNameList:"ExecutableBy" type:"list"`

	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	ImageIds []*string `locationName:"ImageId" locationNameList:"ImageId" type:"list"`

	Owners []*string `locationName:"Owner" locationNameList:"Owner" type:"list"`
}

func (s DescribeImagesInput) String() string {
	return awsutil.Prettify(s)
}

func (s DescribeImagesInput) GoString() string {
	return s.String()
}

func (s *DescribeImagesInput) SetDryRun(v bool) *DescribeImagesInput {
	s.DryRun = &v
	return s
}

func (s *DescribeImagesInput) SetExecutableUsers(v []*string) *DescribeImagesInput {
	s.ExecutableUsers = v
	return s
}

func (s *DescribeImagesInput) SetFilters(v []*Filter) *DescribeImagesInput {
	s.Filters = v
	return s
}

func (s *DescribeImagesInput) SetImageIds(v []*string) *DescribeImagesInput {
	s.ImageIds = v
	return s
}

func (s *DescribeImagesInput) SetOwners(v []*string) *DescribeImagesInput {
	s.Owners = v
	return s
}

type DescribeImagesOutput struct {
	_ struct{} `type:"structure"`

	Images []*Image `locationName:"imagesSet" locationNameList:"item" type:"list"`
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

func (s ModifyImageAttributeInput) String() string {
	return awsutil.Prettify(s)
}

func (s ModifyImageAttributeInput) GoString() string {
	return s.String()
}

func (s *ModifyImageAttributeInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "ModifyImageAttributeInput"}
	if s.ImageId == nil {
		invalidParams.Add(request.NewErrParamRequired("ImageId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

func (s *ModifyImageAttributeInput) SetAttribute(v string) *ModifyImageAttributeInput {
	s.Attribute = &v
	return s
}

func (s *ModifyImageAttributeInput) SetDescription(v *AttributeValue) *ModifyImageAttributeInput {
	s.Description = v
	return s
}

func (s *ModifyImageAttributeInput) SetDryRun(v bool) *ModifyImageAttributeInput {
	s.DryRun = &v
	return s
}

func (s *ModifyImageAttributeInput) SetImageId(v string) *ModifyImageAttributeInput {
	s.ImageId = &v
	return s
}

func (s *ModifyImageAttributeInput) SetLaunchPermission(v *LaunchPermissionModifications) *ModifyImageAttributeInput {
	s.LaunchPermission = v
	return s
}

func (s *ModifyImageAttributeInput) SetOperationType(v string) *ModifyImageAttributeInput {
	s.OperationType = &v
	return s
}

func (s *ModifyImageAttributeInput) SetProductCodes(v []*string) *ModifyImageAttributeInput {
	s.ProductCodes = v
	return s
}

func (s *ModifyImageAttributeInput) SetUserGroups(v []*string) *ModifyImageAttributeInput {
	s.UserGroups = v
	return s
}

func (s *ModifyImageAttributeInput) SetUserIds(v []*string) *ModifyImageAttributeInput {
	s.UserIds = v
	return s
}

func (s *ModifyImageAttributeInput) SetValue(v string) *ModifyImageAttributeInput {
	s.Value = &v
	return s
}

type ModifyImageAttributeOutput struct {
	_ struct{} `type:"structure"`
}

func (s ModifyImageAttributeOutput) String() string {
	return awsutil.Prettify(s)
}

func (s ModifyImageAttributeOutput) GoString() string {
	return s.String()
}

type LaunchPermissionModifications struct {
	_ struct{} `type:"structure"`

	Add []*LaunchPermission `locationNameList:"item" type:"list"`

	Remove []*LaunchPermission `locationNameList:"item" type:"list"`
}

func (s LaunchPermissionModifications) String() string {
	return awsutil.Prettify(s)
}

func (s LaunchPermissionModifications) GoString() string {
	return s.String()
}

func (s *LaunchPermissionModifications) SetAdd(v []*LaunchPermission) *LaunchPermissionModifications {
	s.Add = v
	return s
}

func (s *LaunchPermissionModifications) SetRemove(v []*LaunchPermission) *LaunchPermissionModifications {
	s.Remove = v
	return s
}

type LaunchPermission struct {
	_ struct{} `type:"structure"`

	Group *string `locationName:"group" type:"string" enum:"PermissionGroup"`

	UserId *string `locationName:"userId" type:"string"`
}

func (s LaunchPermission) String() string {
	return awsutil.Prettify(s)
}

func (s LaunchPermission) GoString() string {
	return s.String()
}

func (s *LaunchPermission) SetGroup(v string) *LaunchPermission {
	s.Group = &v
	return s
}

func (s *LaunchPermission) SetUserId(v string) *LaunchPermission {
	s.UserId = &v
	return s
}

type DeleteTagsInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	Resources []*string `locationName:"resourceId" type:"list" required:"true"`

	Tags []*Tag `locationName:"tag" locationNameList:"item" type:"list"`
}

func (s DeleteTagsInput) String() string {
	return awsutil.Prettify(s)
}

func (s DeleteTagsInput) GoString() string {
	return s.String()
}

func (s *DeleteTagsInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeleteTagsInput"}
	if s.Resources == nil {
		invalidParams.Add(request.NewErrParamRequired("Resources"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

func (s *DeleteTagsInput) SetDryRun(v bool) *DeleteTagsInput {
	s.DryRun = &v
	return s
}

func (s *DeleteTagsInput) SetResources(v []*string) *DeleteTagsInput {
	s.Resources = v
	return s
}

func (s *DeleteTagsInput) SetTags(v []*Tag) *DeleteTagsInput {
	s.Tags = v
	return s
}

type DeleteTagsOutput struct {
	_ struct{} `type:"structure"`
}

func (s DeleteTagsOutput) String() string {
	return awsutil.Prettify(s)
}

func (s DeleteTagsOutput) GoString() string {
	return s.String()
}

type CreateTagsInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	Resources []*string `locationName:"ResourceId" type:"list" required:"true"`

	Tags []*Tag `locationName:"Tag" locationNameList:"item" type:"list" required:"true"`
}

func (s CreateTagsInput) String() string {
	return awsutil.Prettify(s)
}

func (s CreateTagsInput) GoString() string {
	return s.String()
}

func (s *CreateTagsInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "CreateTagsInput"}
	if s.Resources == nil {
		invalidParams.Add(request.NewErrParamRequired("Resources"))
	}
	if s.Tags == nil {
		invalidParams.Add(request.NewErrParamRequired("Tags"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

func (s *CreateTagsInput) SetDryRun(v bool) *CreateTagsInput {
	s.DryRun = &v
	return s
}

func (s *CreateTagsInput) SetResources(v []*string) *CreateTagsInput {
	s.Resources = v
	return s
}

func (s *CreateTagsInput) SetTags(v []*Tag) *CreateTagsInput {
	s.Tags = v
	return s
}

type CreateTagsOutput struct {
	_ struct{} `type:"structure"`
}

func (s CreateTagsOutput) String() string {
	return awsutil.Prettify(s)
}

func (s CreateTagsOutput) GoString() string {
	return s.String()
}

type DescribeTagsInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	MaxResults *int64 `locationName:"maxResults" type:"integer"`

	NextToken *string `locationName:"nextToken" type:"string"`
}

func (s DescribeTagsInput) String() string {
	return awsutil.Prettify(s)
}

func (s DescribeTagsInput) GoString() string {
	return s.String()
}

func (s *DescribeTagsInput) SetDryRun(v bool) *DescribeTagsInput {
	s.DryRun = &v
	return s
}

func (s *DescribeTagsInput) SetFilters(v []*Filter) *DescribeTagsInput {
	s.Filters = v
	return s
}

func (s *DescribeTagsInput) SetMaxResults(v int64) *DescribeTagsInput {
	s.MaxResults = &v
	return s
}

func (s *DescribeTagsInput) SetNextToken(v string) *DescribeTagsInput {
	s.NextToken = &v
	return s
}

type DescribeTagsOutput struct {
	_ struct{} `type:"structure"`

	NextToken *string `locationName:"nextToken" type:"string"`

	Tags []*TagDescription `locationName:"tagSet" locationNameList:"item" type:"list"`
}

func (s DescribeTagsOutput) String() string {
	return awsutil.Prettify(s)
}

func (s DescribeTagsOutput) GoString() string {
	return s.String()
}

func (s *DescribeTagsOutput) SetNextToken(v string) *DescribeTagsOutput {
	s.NextToken = &v
	return s
}

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

func (s TagDescription) String() string {
	return awsutil.Prettify(s)
}

func (s TagDescription) GoString() string {
	return s.String()
}

func (s *TagDescription) SetKey(v string) *TagDescription {
	s.Key = &v
	return s
}

func (s *TagDescription) SetResourceId(v string) *TagDescription {
	s.ResourceId = &v
	return s
}

func (s *TagDescription) SetResourceType(v string) *TagDescription {
	s.ResourceType = &v
	return s
}

func (s *TagDescription) SetValue(v string) *TagDescription {
	s.Value = &v
	return s
}

type TagSpecification struct {
	_ struct{} `type:"structure"`

	ResourceType *string `locationName:"resourceType" type:"string" enum:"ResourceType"`

	Tags []*Tag `locationName:"Tag" locationNameList:"item" type:"list"`
}

func (s TagSpecification) String() string {
	return awsutil.Prettify(s)
}

func (s TagSpecification) GoString() string {
	return s.String()
}

func (s *TagSpecification) SetResourceType(v string) *TagSpecification {
	s.ResourceType = &v
	return s
}

func (s *TagSpecification) SetTags(v []*Tag) *TagSpecification {
	s.Tags = v
	return s
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
func (s ImportKeyPairInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ImportKeyPairInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *ImportKeyPairInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "ImportKeyPairInput"}
	if s.KeyName == nil {
		invalidParams.Add(request.NewErrParamRequired("KeyName"))
	}
	if s.PublicKeyMaterial == nil {
		invalidParams.Add(request.NewErrParamRequired("PublicKeyMaterial"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDryRun sets the DryRun field's value.
func (s *ImportKeyPairInput) SetDryRun(v bool) *ImportKeyPairInput {
	s.DryRun = &v
	return s
}

// SetKeyName sets the KeyName field's value.
func (s *ImportKeyPairInput) SetKeyName(v string) *ImportKeyPairInput {
	s.KeyName = &v
	return s
}

// SetPublicKeyMaterial sets the PublicKeyMaterial field's value.
func (s *ImportKeyPairInput) SetPublicKeyMaterial(v []byte) *ImportKeyPairInput {
	s.PublicKeyMaterial = v
	return s
}

// Contains the output of ImportKeyPair.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/ImportKeyPairResult
type ImportKeyPairOutput struct {
	_ struct{} `type:"structure"`

	// The MD5 public key fingerprint as specified in section 4 of RFC 4716.
	KeyFingerprint *string `locationName:"keyFingerprint" type:"string"`

	// The key pair name you provided.
	KeyName *string `locationName:"keyName" type:"string"`
}

// String returns the string representation
func (s ImportKeyPairOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ImportKeyPairOutput) GoString() string {
	return s.String()
}

// SetKeyFingerprint sets the KeyFingerprint field's value.
func (s *ImportKeyPairOutput) SetKeyFingerprint(v string) *ImportKeyPairOutput {
	s.KeyFingerprint = &v
	return s
}

// SetKeyName sets the KeyName field's value.
func (s *ImportKeyPairOutput) SetKeyName(v string) *ImportKeyPairOutput {
	s.KeyName = &v
	return s
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

// String returns the string representation
func (s DescribeKeyPairsInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeKeyPairsInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeKeyPairsInput) SetDryRun(v bool) *DescribeKeyPairsInput {
	s.DryRun = &v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribeKeyPairsInput) SetFilters(v []*Filter) *DescribeKeyPairsInput {
	s.Filters = v
	return s
}

// SetKeyNames sets the KeyNames field's value.
func (s *DescribeKeyPairsInput) SetKeyNames(v []*string) *DescribeKeyPairsInput {
	s.KeyNames = v
	return s
}

// Contains the output of DescribeKeyPairs.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeKeyPairsResult
type DescribeKeyPairsOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more key pairs.
	KeyPairs []*KeyPairInfo `locationName:"keySet" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s DescribeKeyPairsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeKeyPairsOutput) GoString() string {
	return s.String()
}

// SetKeyPairs sets the KeyPairs field's value.
func (s *DescribeKeyPairsOutput) SetKeyPairs(v []*KeyPairInfo) *DescribeKeyPairsOutput {
	s.KeyPairs = v
	return s
}

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
func (s KeyPairInfo) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s KeyPairInfo) GoString() string {
	return s.String()
}

// SetKeyFingerprint sets the KeyFingerprint field's value.
func (s *KeyPairInfo) SetKeyFingerprint(v string) *KeyPairInfo {
	s.KeyFingerprint = &v
	return s
}

// SetKeyName sets the KeyName field's value.
func (s *KeyPairInfo) SetKeyName(v string) *KeyPairInfo {
	s.KeyName = &v
	return s
}

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

// String returns the string representation
func (s DeleteKeyPairInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteKeyPairInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteKeyPairInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeleteKeyPairInput"}
	if s.KeyName == nil {
		invalidParams.Add(request.NewErrParamRequired("KeyName"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDryRun sets the DryRun field's value.
func (s *DeleteKeyPairInput) SetDryRun(v bool) *DeleteKeyPairInput {
	s.DryRun = &v
	return s
}

// SetKeyName sets the KeyName field's value.
func (s *DeleteKeyPairInput) SetKeyName(v string) *DeleteKeyPairInput {
	s.KeyName = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteKeyPairOutput
type DeleteKeyPairOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeleteKeyPairOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteKeyPairOutput) GoString() string {
	return s.String()
}

// Contains the parameters for CreateKeyPair.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateKeyPairRequest
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

// String returns the string representation
func (s CreateKeyPairInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateKeyPairInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CreateKeyPairInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "CreateKeyPairInput"}
	if s.KeyName == nil {
		invalidParams.Add(request.NewErrParamRequired("KeyName"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDryRun sets the DryRun field's value.
func (s *CreateKeyPairInput) SetDryRun(v bool) *CreateKeyPairInput {
	s.DryRun = &v
	return s
}

// SetKeyName sets the KeyName field's value.
func (s *CreateKeyPairInput) SetKeyName(v string) *CreateKeyPairInput {
	s.KeyName = &v
	return s
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

// String returns the string representation
func (s CreateKeyPairOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateKeyPairOutput) GoString() string {
	return s.String()
}

// SetKeyFingerprint sets the KeyFingerprint field's value.
func (s *CreateKeyPairOutput) SetKeyFingerprint(v string) *CreateKeyPairOutput {
	s.KeyFingerprint = &v
	return s
}

// SetKeyMaterial sets the KeyMaterial field's value.
func (s *CreateKeyPairOutput) SetKeyMaterial(v string) *CreateKeyPairOutput {
	s.KeyMaterial = &v
	return s
}

// SetKeyName sets the KeyName field's value.
func (s *CreateKeyPairOutput) SetKeyName(v string) *CreateKeyPairOutput {
	s.KeyName = &v
	return s
}

type CreateSecurityGroupInput struct {
	_ struct{} `type:"structure"`

	Description *string `locationName:"GroupDescription" type:"string" required:"true"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	GroupName *string `type:"string" required:"true"`

	VpcId *string `type:"string"`
}

func (s CreateSecurityGroupInput) String() string {
	return awsutil.Prettify(s)
}

func (s CreateSecurityGroupInput) GoString() string {
	return s.String()
}

func (s *CreateSecurityGroupInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "CreateSecurityGroupInput"}
	if s.Description == nil {
		invalidParams.Add(request.NewErrParamRequired("Description"))
	}
	if s.GroupName == nil {
		invalidParams.Add(request.NewErrParamRequired("GroupName"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

func (s *CreateSecurityGroupInput) SetDescription(v string) *CreateSecurityGroupInput {
	s.Description = &v
	return s
}

func (s *CreateSecurityGroupInput) SetDryRun(v bool) *CreateSecurityGroupInput {
	s.DryRun = &v
	return s
}

func (s *CreateSecurityGroupInput) SetGroupName(v string) *CreateSecurityGroupInput {
	s.GroupName = &v
	return s
}

func (s *CreateSecurityGroupInput) SetVpcId(v string) *CreateSecurityGroupInput {
	s.VpcId = &v
	return s
}

type CreateSecurityGroupOutput struct {
	_ struct{} `type:"structure"`

	GroupId *string `locationName:"groupId" type:"string"`
}

func (s CreateSecurityGroupOutput) String() string {
	return awsutil.Prettify(s)
}

func (s CreateSecurityGroupOutput) GoString() string {
	return s.String()
}

func (s *CreateSecurityGroupOutput) SetGroupId(v string) *CreateSecurityGroupOutput {
	s.GroupId = &v
	return s
}

type DescribeSecurityGroupsInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	GroupIds []*string `locationName:"GroupId" locationNameList:"groupId" type:"list"`

	GroupNames []*string `locationName:"GroupName" locationNameList:"GroupName" type:"list"`
}

func (s DescribeSecurityGroupsInput) String() string {
	return awsutil.Prettify(s)
}

func (s DescribeSecurityGroupsInput) GoString() string {
	return s.String()
}

func (s *DescribeSecurityGroupsInput) SetDryRun(v bool) *DescribeSecurityGroupsInput {
	s.DryRun = &v
	return s
}

func (s *DescribeSecurityGroupsInput) SetFilters(v []*Filter) *DescribeSecurityGroupsInput {
	s.Filters = v
	return s
}

func (s *DescribeSecurityGroupsInput) SetGroupIds(v []*string) *DescribeSecurityGroupsInput {
	s.GroupIds = v
	return s
}

func (s *DescribeSecurityGroupsInput) SetGroupNames(v []*string) *DescribeSecurityGroupsInput {
	s.GroupNames = v
	return s
}

type DescribeSecurityGroupsOutput struct {
	_ struct{} `type:"structure"`

	SecurityGroups []*SecurityGroup `locationName:"securityGroupInfo" locationNameList:"item" type:"list"`
}

func (s DescribeSecurityGroupsOutput) String() string {
	return awsutil.Prettify(s)
}

func (s DescribeSecurityGroupsOutput) GoString() string {
	return s.String()
}

func (s *DescribeSecurityGroupsOutput) SetSecurityGroups(v []*SecurityGroup) *DescribeSecurityGroupsOutput {
	s.SecurityGroups = v
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

func (s SecurityGroup) String() string {
	return awsutil.Prettify(s)
}

func (s SecurityGroup) GoString() string {
	return s.String()
}

func (s *SecurityGroup) SetDescription(v string) *SecurityGroup {
	s.Description = &v
	return s
}

func (s *SecurityGroup) SetGroupId(v string) *SecurityGroup {
	s.GroupId = &v
	return s
}

func (s *SecurityGroup) SetGroupName(v string) *SecurityGroup {
	s.GroupName = &v
	return s
}

func (s *SecurityGroup) SetIpPermissions(v []*IpPermission) *SecurityGroup {
	s.IpPermissions = v
	return s
}

func (s *SecurityGroup) SetIpPermissionsEgress(v []*IpPermission) *SecurityGroup {
	s.IpPermissionsEgress = v
	return s
}

func (s *SecurityGroup) SetOwnerId(v string) *SecurityGroup {
	s.OwnerId = &v
	return s
}

func (s *SecurityGroup) SetTags(v []*Tag) *SecurityGroup {
	s.Tags = v
	return s
}

func (s *SecurityGroup) SetVpcId(v string) *SecurityGroup {
	s.VpcId = &v
	return s
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

func (s IpPermission) String() string {
	return awsutil.Prettify(s)
}

func (s IpPermission) GoString() string {
	return s.String()
}

func (s *IpPermission) SetFromPort(v int64) *IpPermission {
	s.FromPort = &v
	return s
}

func (s *IpPermission) SetIpProtocol(v string) *IpPermission {
	s.IpProtocol = &v
	return s
}

func (s *IpPermission) SetIpRanges(v []*IpRange) *IpPermission {
	s.IpRanges = v
	return s
}

func (s *IpPermission) SetIpv6Ranges(v []*Ipv6Range) *IpPermission {
	s.Ipv6Ranges = v
	return s
}

func (s *IpPermission) SetPrefixListIds(v []*PrefixListId) *IpPermission {
	s.PrefixListIds = v
	return s
}

func (s *IpPermission) SetToPort(v int64) *IpPermission {
	s.ToPort = &v
	return s
}

func (s *IpPermission) SetUserIdGroupPairs(v []*UserIdGroupPair) *IpPermission {
	s.UserIdGroupPairs = v
	return s
}

type IpRange struct {
	_ struct{} `type:"structure"`

	CidrIp *string `locationName:"cidrIp" type:"string"`
}

func (s IpRange) String() string {
	return awsutil.Prettify(s)
}

func (s IpRange) GoString() string {
	return s.String()
}

func (s *IpRange) SetCidrIp(v string) *IpRange {
	s.CidrIp = &v
	return s
}

type Ipv6Range struct {
	_ struct{} `type:"structure"`

	CidrIpv6 *string `locationName:"cidrIpv6" type:"string"`
}

func (s Ipv6Range) String() string {
	return awsutil.Prettify(s)
}

func (s Ipv6Range) GoString() string {
	return s.String()
}

func (s *Ipv6Range) SetCidrIpv6(v string) *Ipv6Range {
	s.CidrIpv6 = &v
	return s
}

type PrefixListId struct {
	_ struct{} `type:"structure"`

	PrefixListId *string `locationName:"prefixListId" type:"string"`
}

func (s PrefixListId) String() string {
	return awsutil.Prettify(s)
}

func (s PrefixListId) GoString() string {
	return s.String()
}

func (s *PrefixListId) SetPrefixListId(v string) *PrefixListId {
	s.PrefixListId = &v
	return s
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

func (s UserIdGroupPair) String() string {
	return awsutil.Prettify(s)
}

func (s UserIdGroupPair) GoString() string {
	return s.String()
}

func (s *UserIdGroupPair) SetGroupId(v string) *UserIdGroupPair {
	s.GroupId = &v
	return s
}

func (s *UserIdGroupPair) SetGroupName(v string) *UserIdGroupPair {
	s.GroupName = &v
	return s
}

func (s *UserIdGroupPair) SetPeeringStatus(v string) *UserIdGroupPair {
	s.PeeringStatus = &v
	return s
}

func (s *UserIdGroupPair) SetUserId(v string) *UserIdGroupPair {
	s.UserId = &v
	return s
}

func (s *UserIdGroupPair) SetVpcId(v string) *UserIdGroupPair {
	s.VpcId = &v
	return s
}

func (s *UserIdGroupPair) SetVpcPeeringConnectionId(v string) *UserIdGroupPair {
	s.VpcPeeringConnectionId = &v
	return s
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

func (s RevokeSecurityGroupEgressInput) String() string {
	return awsutil.Prettify(s)
}

func (s RevokeSecurityGroupEgressInput) GoString() string {
	return s.String()
}

func (s *RevokeSecurityGroupEgressInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "RevokeSecurityGroupEgressInput"}
	if s.GroupId == nil {
		invalidParams.Add(request.NewErrParamRequired("GroupId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

func (s *RevokeSecurityGroupEgressInput) SetCidrIp(v string) *RevokeSecurityGroupEgressInput {
	s.CidrIp = &v
	return s
}

func (s *RevokeSecurityGroupEgressInput) SetDryRun(v bool) *RevokeSecurityGroupEgressInput {
	s.DryRun = &v
	return s
}

func (s *RevokeSecurityGroupEgressInput) SetFromPort(v int64) *RevokeSecurityGroupEgressInput {
	s.FromPort = &v
	return s
}

func (s *RevokeSecurityGroupEgressInput) SetGroupId(v string) *RevokeSecurityGroupEgressInput {
	s.GroupId = &v
	return s
}

func (s *RevokeSecurityGroupEgressInput) SetIpPermissions(v []*IpPermission) *RevokeSecurityGroupEgressInput {
	s.IpPermissions = v
	return s
}

func (s *RevokeSecurityGroupEgressInput) SetIpProtocol(v string) *RevokeSecurityGroupEgressInput {
	s.IpProtocol = &v
	return s
}

func (s *RevokeSecurityGroupEgressInput) SetSourceSecurityGroupName(v string) *RevokeSecurityGroupEgressInput {
	s.SourceSecurityGroupName = &v
	return s
}

func (s *RevokeSecurityGroupEgressInput) SetSourceSecurityGroupOwnerId(v string) *RevokeSecurityGroupEgressInput {
	s.SourceSecurityGroupOwnerId = &v
	return s
}

func (s *RevokeSecurityGroupEgressInput) SetToPort(v int64) *RevokeSecurityGroupEgressInput {
	s.ToPort = &v
	return s
}

type RevokeSecurityGroupEgressOutput struct {
	_ struct{} `type:"structure"`
}

func (s RevokeSecurityGroupEgressOutput) String() string {
	return awsutil.Prettify(s)
}

func (s RevokeSecurityGroupEgressOutput) GoString() string {
	return s.String()
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

func (s RevokeSecurityGroupIngressInput) String() string {
	return awsutil.Prettify(s)
}

func (s RevokeSecurityGroupIngressInput) GoString() string {
	return s.String()
}

func (s *RevokeSecurityGroupIngressInput) SetCidrIp(v string) *RevokeSecurityGroupIngressInput {
	s.CidrIp = &v
	return s
}

func (s *RevokeSecurityGroupIngressInput) SetDryRun(v bool) *RevokeSecurityGroupIngressInput {
	s.DryRun = &v
	return s
}

func (s *RevokeSecurityGroupIngressInput) SetFromPort(v int64) *RevokeSecurityGroupIngressInput {
	s.FromPort = &v
	return s
}

func (s *RevokeSecurityGroupIngressInput) SetGroupId(v string) *RevokeSecurityGroupIngressInput {
	s.GroupId = &v
	return s
}

func (s *RevokeSecurityGroupIngressInput) SetGroupName(v string) *RevokeSecurityGroupIngressInput {
	s.GroupName = &v
	return s
}

func (s *RevokeSecurityGroupIngressInput) SetIpPermissions(v []*IpPermission) *RevokeSecurityGroupIngressInput {
	s.IpPermissions = v
	return s
}

func (s *RevokeSecurityGroupIngressInput) SetIpProtocol(v string) *RevokeSecurityGroupIngressInput {
	s.IpProtocol = &v
	return s
}

func (s *RevokeSecurityGroupIngressInput) SetSourceSecurityGroupName(v string) *RevokeSecurityGroupIngressInput {
	s.SourceSecurityGroupName = &v
	return s
}

func (s *RevokeSecurityGroupIngressInput) SetSourceSecurityGroupOwnerId(v string) *RevokeSecurityGroupIngressInput {
	s.SourceSecurityGroupOwnerId = &v
	return s
}

func (s *RevokeSecurityGroupIngressInput) SetToPort(v int64) *RevokeSecurityGroupIngressInput {
	s.ToPort = &v
	return s
}

type RevokeSecurityGroupIngressOutput struct {
	_ struct{} `type:"structure"`
}

func (s RevokeSecurityGroupIngressOutput) String() string {
	return awsutil.Prettify(s)
}

func (s RevokeSecurityGroupIngressOutput) GoString() string {
	return s.String()
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

func (s AuthorizeSecurityGroupEgressInput) String() string {
	return awsutil.Prettify(s)
}

func (s AuthorizeSecurityGroupEgressInput) GoString() string {
	return s.String()
}

func (s *AuthorizeSecurityGroupEgressInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "AuthorizeSecurityGroupEgressInput"}
	if s.GroupId == nil {
		invalidParams.Add(request.NewErrParamRequired("GroupId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

func (s *AuthorizeSecurityGroupEgressInput) SetCidrIp(v string) *AuthorizeSecurityGroupEgressInput {
	s.CidrIp = &v
	return s
}

func (s *AuthorizeSecurityGroupEgressInput) SetDryRun(v bool) *AuthorizeSecurityGroupEgressInput {
	s.DryRun = &v
	return s
}

func (s *AuthorizeSecurityGroupEgressInput) SetFromPort(v int64) *AuthorizeSecurityGroupEgressInput {
	s.FromPort = &v
	return s
}

func (s *AuthorizeSecurityGroupEgressInput) SetGroupId(v string) *AuthorizeSecurityGroupEgressInput {
	s.GroupId = &v
	return s
}

func (s *AuthorizeSecurityGroupEgressInput) SetIpPermissions(v []*IpPermission) *AuthorizeSecurityGroupEgressInput {
	s.IpPermissions = v
	return s
}

func (s *AuthorizeSecurityGroupEgressInput) SetIpProtocol(v string) *AuthorizeSecurityGroupEgressInput {
	s.IpProtocol = &v
	return s
}

func (s *AuthorizeSecurityGroupEgressInput) SetSourceSecurityGroupName(v string) *AuthorizeSecurityGroupEgressInput {
	s.SourceSecurityGroupName = &v
	return s
}

func (s *AuthorizeSecurityGroupEgressInput) SetSourceSecurityGroupOwnerId(v string) *AuthorizeSecurityGroupEgressInput {
	s.SourceSecurityGroupOwnerId = &v
	return s
}

func (s *AuthorizeSecurityGroupEgressInput) SetToPort(v int64) *AuthorizeSecurityGroupEgressInput {
	s.ToPort = &v
	return s
}

type AuthorizeSecurityGroupEgressOutput struct {
	_ struct{} `type:"structure"`
}

func (s AuthorizeSecurityGroupEgressOutput) String() string {
	return awsutil.Prettify(s)
}

func (s AuthorizeSecurityGroupEgressOutput) GoString() string {
	return s.String()
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

func (s AuthorizeSecurityGroupIngressInput) String() string {
	return awsutil.Prettify(s)
}

func (s AuthorizeSecurityGroupIngressInput) GoString() string {
	return s.String()
}

func (s *AuthorizeSecurityGroupIngressInput) SetCidrIp(v string) *AuthorizeSecurityGroupIngressInput {
	s.CidrIp = &v
	return s
}

func (s *AuthorizeSecurityGroupIngressInput) SetDryRun(v bool) *AuthorizeSecurityGroupIngressInput {
	s.DryRun = &v
	return s
}

func (s *AuthorizeSecurityGroupIngressInput) SetFromPort(v int64) *AuthorizeSecurityGroupIngressInput {
	s.FromPort = &v
	return s
}

func (s *AuthorizeSecurityGroupIngressInput) SetGroupId(v string) *AuthorizeSecurityGroupIngressInput {
	s.GroupId = &v
	return s
}

func (s *AuthorizeSecurityGroupIngressInput) SetGroupName(v string) *AuthorizeSecurityGroupIngressInput {
	s.GroupName = &v
	return s
}

func (s *AuthorizeSecurityGroupIngressInput) SetIpPermissions(v []*IpPermission) *AuthorizeSecurityGroupIngressInput {
	s.IpPermissions = v
	return s
}

func (s *AuthorizeSecurityGroupIngressInput) SetIpProtocol(v string) *AuthorizeSecurityGroupIngressInput {
	s.IpProtocol = &v
	return s
}

func (s *AuthorizeSecurityGroupIngressInput) SetSourceSecurityGroupName(v string) *AuthorizeSecurityGroupIngressInput {
	s.SourceSecurityGroupName = &v
	return s
}

func (s *AuthorizeSecurityGroupIngressInput) SetSourceSecurityGroupOwnerId(v string) *AuthorizeSecurityGroupIngressInput {
	s.SourceSecurityGroupOwnerId = &v
	return s
}

func (s *AuthorizeSecurityGroupIngressInput) SetToPort(v int64) *AuthorizeSecurityGroupIngressInput {
	s.ToPort = &v
	return s
}

type AuthorizeSecurityGroupIngressOutput struct {
	_ struct{} `type:"structure"`
}

func (s AuthorizeSecurityGroupIngressOutput) String() string {
	return awsutil.Prettify(s)
}

func (s AuthorizeSecurityGroupIngressOutput) GoString() string {
	return s.String()
}

type DeleteSecurityGroupInput struct {
	_ struct{} `type:"structure"`

	DryRun *bool `locationName:"dryRun" type:"boolean"`

	GroupId *string `type:"string"`

	GroupName *string `type:"string"`
}

func (s DeleteSecurityGroupInput) String() string {
	return awsutil.Prettify(s)
}

func (s DeleteSecurityGroupInput) GoString() string {
	return s.String()
}

func (s *DeleteSecurityGroupInput) SetDryRun(v bool) *DeleteSecurityGroupInput {
	s.DryRun = &v
	return s
}

func (s *DeleteSecurityGroupInput) SetGroupId(v string) *DeleteSecurityGroupInput {
	s.GroupId = &v
	return s
}

func (s *DeleteSecurityGroupInput) SetGroupName(v string) *DeleteSecurityGroupInput {
	s.GroupName = &v
	return s
}

type DeleteSecurityGroupOutput struct {
	_ struct{} `type:"structure"`
}

func (s DeleteSecurityGroupOutput) String() string {
	return awsutil.Prettify(s)
}

func (s DeleteSecurityGroupOutput) GoString() string {
	return s.String()
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

// String returns the string representation
func (s CreateVolumeInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateVolumeInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CreateVolumeInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "CreateVolumeInput"}
	if s.AvailabilityZone == nil {
		invalidParams.Add(request.NewErrParamRequired("AvailabilityZone"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAvailabilityZone sets the AvailabilityZone field's value.
func (s *CreateVolumeInput) SetAvailabilityZone(v string) *CreateVolumeInput {
	s.AvailabilityZone = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *CreateVolumeInput) SetDryRun(v bool) *CreateVolumeInput {
	s.DryRun = &v
	return s
}

// SetEncrypted sets the Encrypted field's value.
func (s *CreateVolumeInput) SetEncrypted(v bool) *CreateVolumeInput {
	s.Encrypted = &v
	return s
}

// SetIops sets the Iops field's value.
func (s *CreateVolumeInput) SetIops(v int64) *CreateVolumeInput {
	s.Iops = &v
	return s
}

// SetKmsKeyId sets the KmsKeyId field's value.
func (s *CreateVolumeInput) SetKmsKeyId(v string) *CreateVolumeInput {
	s.KmsKeyId = &v
	return s
}

// SetSize sets the Size field's value.
func (s *CreateVolumeInput) SetSize(v int64) *CreateVolumeInput {
	s.Size = &v
	return s
}

// SetSnapshotId sets the SnapshotId field's value.
func (s *CreateVolumeInput) SetSnapshotId(v string) *CreateVolumeInput {
	s.SnapshotId = &v
	return s
}

// SetTagSpecifications sets the TagSpecifications field's value.
func (s *CreateVolumeInput) SetTagSpecifications(v []*TagSpecification) *CreateVolumeInput {
	s.TagSpecifications = v
	return s
}

// SetVolumeType sets the VolumeType field's value.
func (s *CreateVolumeInput) SetVolumeType(v string) *CreateVolumeInput {
	s.VolumeType = &v
	return s
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

// String returns the string representation
func (s DeleteVolumeInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteVolumeInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteVolumeInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeleteVolumeInput"}
	if s.VolumeId == nil {
		invalidParams.Add(request.NewErrParamRequired("VolumeId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDryRun sets the DryRun field's value.
func (s *DeleteVolumeInput) SetDryRun(v bool) *DeleteVolumeInput {
	s.DryRun = &v
	return s
}

// SetVolumeId sets the VolumeId field's value.
func (s *DeleteVolumeInput) SetVolumeId(v string) *DeleteVolumeInput {
	s.VolumeId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteVolumeOutput
type DeleteVolumeOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeleteVolumeOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteVolumeOutput) GoString() string {
	return s.String()
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

// String returns the string representation
func (s DescribeVolumesInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeVolumesInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeVolumesInput) SetDryRun(v bool) *DescribeVolumesInput {
	s.DryRun = &v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribeVolumesInput) SetFilters(v []*Filter) *DescribeVolumesInput {
	s.Filters = v
	return s
}

// SetMaxResults sets the MaxResults field's value.
func (s *DescribeVolumesInput) SetMaxResults(v int64) *DescribeVolumesInput {
	s.MaxResults = &v
	return s
}

// SetNextToken sets the NextToken field's value.
func (s *DescribeVolumesInput) SetNextToken(v string) *DescribeVolumesInput {
	s.NextToken = &v
	return s
}

// SetVolumeIds sets the VolumeIds field's value.
func (s *DescribeVolumesInput) SetVolumeIds(v []*string) *DescribeVolumesInput {
	s.VolumeIds = v
	return s
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

	// The number of I/O operations per second (IOPS) that the volume supports.
	// For Provisioned IOPS SSD volumes, this represents the number of IOPS that
	// are provisioned for the volume. For General Purpose SSD volumes, this represents
	// the baseline performance of the volume and the rate at which the volume accumulates
	// I/O credits for bursting. For more information on General Purpose SSD baseline
	// performance, I/O credits, and bursting, see Amazon EBS Volume Types (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html)
	// in the Amazon Elastic Compute Cloud User Guide.
	//
	// Constraint: Range is 100-20000 IOPS for io1 volumes and 100-10000 IOPS for
	// gp2 volumes.
	//
	// Condition: This parameter is required for requests to create io1 volumes;
	// it is not used in requests to create gp2, st1, sc1, or standard volumes.
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

// Describes volume attachment details.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/VolumeAttachment
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

// Contains the parameters for AttachVolume.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/AttachVolumeRequest
type AttachVolumeInput struct {
	_ struct{} `type:"structure"`

	// The device name to expose to the instance (for example, /dev/sdh or xvdh).
	//
	// Device is a required field
	Device *string `type:"string" required:"true"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the instance.
	//
	// InstanceId is a required field
	InstanceId *string `type:"string" required:"true"`

	// The ID of the EBS volume. The volume and instance must be within the same
	// Availability Zone.
	//
	// VolumeId is a required field
	VolumeId *string `type:"string" required:"true"`
}

// String returns the string representation
func (s AttachVolumeInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AttachVolumeInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *AttachVolumeInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "AttachVolumeInput"}
	if s.Device == nil {
		invalidParams.Add(request.NewErrParamRequired("Device"))
	}
	if s.InstanceId == nil {
		invalidParams.Add(request.NewErrParamRequired("InstanceId"))
	}
	if s.VolumeId == nil {
		invalidParams.Add(request.NewErrParamRequired("VolumeId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDevice sets the Device field's value.
func (s *AttachVolumeInput) SetDevice(v string) *AttachVolumeInput {
	s.Device = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *AttachVolumeInput) SetDryRun(v bool) *AttachVolumeInput {
	s.DryRun = &v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *AttachVolumeInput) SetInstanceId(v string) *AttachVolumeInput {
	s.InstanceId = &v
	return s
}

// SetVolumeId sets the VolumeId field's value.
func (s *AttachVolumeInput) SetVolumeId(v string) *AttachVolumeInput {
	s.VolumeId = &v
	return s
}

type DetachVolumeInput struct {
	_ struct{} `type:"structure"`

	// The device name.
	Device *string `type:"string"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// Forces detachment if the previous detachment attempt did not occur cleanly
	// (for example, logging into an instance, unmounting the volume, and detaching
	// normally). This option can lead to data loss or a corrupted file system.
	// Use this option only as a last resort to detach a volume from a failed instance.
	// The instance won't have an opportunity to flush file system caches or file
	// system metadata. If you use this option, you must perform file system check
	// and repair procedures.
	Force *bool `type:"boolean"`

	// The ID of the instance.
	InstanceId *string `type:"string"`

	// The ID of the volume.
	//
	// VolumeId is a required field
	VolumeId *string `type:"string" required:"true"`
}

// String returns the string representation
func (s DetachVolumeInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DetachVolumeInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DetachVolumeInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DetachVolumeInput"}
	if s.VolumeId == nil {
		invalidParams.Add(request.NewErrParamRequired("VolumeId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetDevice sets the Device field's value.
func (s *DetachVolumeInput) SetDevice(v string) *DetachVolumeInput {
	s.Device = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *DetachVolumeInput) SetDryRun(v bool) *DetachVolumeInput {
	s.DryRun = &v
	return s
}

// SetForce sets the Force field's value.
func (s *DetachVolumeInput) SetForce(v bool) *DetachVolumeInput {
	s.Force = &v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *DetachVolumeInput) SetInstanceId(v string) *DetachVolumeInput {
	s.InstanceId = &v
	return s
}

// SetVolumeId sets the VolumeId field's value.
func (s *DetachVolumeInput) SetVolumeId(v string) *DetachVolumeInput {
	s.VolumeId = &v
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

//DescribeInternetGatewaysOutput Contains the output of DescribeInternetGateways.
type DescribeInternetGatewaysOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more Internet gateways.
	InternetGateways []*InternetGateway `locationName:"internetGatewaySet" locationNameList:"item" type:"list"`
	RequesterId      *string            `locationName:"requestId" type:"string"`
}

// Contains the parameters for DeleteInternetGateway.
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

type DeleteInternetGatewayOutput struct {
	_ struct{} `type:"structure"`
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
