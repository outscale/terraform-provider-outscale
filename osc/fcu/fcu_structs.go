package fcu

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/request"
)

const (
	// InstanceAttributeNameUserData is a InstanceAttributeName enum value
	InstanceAttributeNameUserData = "userData"
)

// DescribeInstancesInput Contains the parameters for DescribeInstances.
type DescribeInstancesInput struct {
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	InstanceIds []*string `locationName:"InstanceId" locationNameList:"InstanceId" type:"list"`

	MaxResults *int64 `locationName:"maxResults" type:"integer"`

	NextToken *string `locationName:"nextToken" type:"string"`
}

// String returns the string representation
func (s DescribeInstancesInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeInstancesInput) GoString() string {
	return s.String()
}

// SetFilters sets the Filters field's value.
func (s *DescribeInstancesInput) SetFilters(v []*Filter) *DescribeInstancesInput {
	s.Filters = v
	return s
}

// SetInstanceIds sets the InstanceIds field's value.
func (s *DescribeInstancesInput) SetInstanceIds(v []*string) *DescribeInstancesInput {
	s.InstanceIds = v
	return s
}

// SetMaxResults sets the MaxResults field's value.
func (s *DescribeInstancesInput) SetMaxResults(v int64) *DescribeInstancesInput {
	s.MaxResults = &v
	return s
}

// SetNextToken sets the NextToken field's value.
func (s *DescribeInstancesInput) SetNextToken(v string) *DescribeInstancesInput {
	s.NextToken = &v
	return s
}

// Filter can be used to match a set of resources by various criteria.
type Filter struct {
	Name *string `type:"string"`

	Values []*string `locationName:"Value" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s Filter) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Filter) GoString() string {
	return s.String()
}

// SetName sets the Name field's value.
func (s *Filter) SetName(v string) *Filter {
	s.Name = &v
	return s
}

// SetValues sets the Values field's value.
func (s *Filter) SetValues(v []*string) *Filter {
	s.Values = v
	return s
}

// DescribeInstancesOutput struct
type DescribeInstancesOutput struct {
	NextToken *string `locationName:"nextToken" type:"string"`

	Reservations []*Reservation `locationName:"reservationSet" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s DescribeInstancesOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeInstancesOutput) GoString() string {
	return s.String()
}

// SetNextToken sets the NextToken field's value.
func (s *DescribeInstancesOutput) SetNextToken(v string) *DescribeInstancesOutput {
	s.NextToken = &v
	return s
}

// SetReservations sets the Reservations field's value.
func (s *DescribeInstancesOutput) SetReservations(v []*Reservation) *DescribeInstancesOutput {
	s.Reservations = v
	return s
}

// Reservation struct
type Reservation struct {
	Groups []*GroupIdentifier `locationName:"groupSet" locationNameList:"item" type:"list"`

	Instances []*Instance `locationName:"instancesSet" locationNameList:"item" type:"list"`

	OwnerId *string `locationName:"ownerId" type:"string"`

	RequesterId *string `locationName:"requesterId" type:"string"`

	ReservationId *string `locationName:"reservationId" type:"string"`
}

// String returns the string representation
func (s Reservation) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Reservation) GoString() string {
	return s.String()
}

// SetGroups sets the Groups field's value.
func (s *Reservation) SetGroups(v []*GroupIdentifier) *Reservation {
	s.Groups = v
	return s
}

// SetInstances sets the Instances field's value.
func (s *Reservation) SetInstances(v []*Instance) *Reservation {
	s.Instances = v
	return s
}

// SetOwnerId sets the OwnerId field's value.
func (s *Reservation) SetOwnerId(v string) *Reservation {
	s.OwnerId = &v
	return s
}

// SetRequesterId sets the RequesterId field's value.
func (s *Reservation) SetRequesterId(v string) *Reservation {
	s.RequesterId = &v
	return s
}

// SetReservationId sets the ReservationId field's value.
func (s *Reservation) SetReservationId(v string) *Reservation {
	s.ReservationId = &v
	return s
}

// GroupIdentifier stuct
type GroupIdentifier struct {
	GroupId *string `locationName:"groupId" type:"string"`

	GroupName *string `locationName:"groupName" type:"string"`
}

// String returns the string representation
func (s GroupIdentifier) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s GroupIdentifier) GoString() string {
	return s.String()
}

// SetGroupId sets the GroupId field's value.
func (s *GroupIdentifier) SetGroupId(v string) *GroupIdentifier {
	s.GroupId = &v
	return s
}

// SetGroupName sets the GroupName field's value.
func (s *GroupIdentifier) SetGroupName(v string) *GroupIdentifier {
	s.GroupName = &v
	return s
}

// Instance struct
type Instance struct {
	AmiLaunchIndex *int64 `locationName:"amiLaunchIndex" type:"integer"`

	Architecture *string `locationName:"architecture" type:"string" enum:"ArchitectureValues"`

	BlockDeviceMappings []*InstanceBlockDeviceMapping `locationName:"blockDeviceMapping" locationNameList:"item" type:"list"`

	ClientToken *string `locationName:"clientToken" type:"string"`

	DnsName *string `type:"string"`

	EbsOptimized *bool `locationName:"ebsOptimized" type:"boolean"`

	GroupSet []*GroupIdentifier `locationName:"groupSet" locationNameList:"item" type:"list"`

	Hypervisor *string `locationName:"hypervisor" type:"string" enum:"HypervisorType"`

	IamInstanceProfile *IamInstanceProfile `locationName:"iamInstanceProfile" type:"structure"`

	ImageId *string `locationName:"imageId" type:"string"`

	InstanceId *string `locationName:"instanceId" type:"string"`

	InstanceLifecycle *string `locationName:"instanceLifecycle" type:"string" enum:"InstanceLifecycleType"`

	InstanceState *InstanceState `locationName:"instanceState" type:"structure"`

	InstanceType *string `locationName:"instanceType" type:"string" enum:"InstanceType"`

	IpAddress *string `type:"string"`

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

	Reason *string `locationName:"ramdiskId" type:"string"`

	RootDeviceName *string `locationName:"rootDeviceName" type:"string"`

	RootDeviceType *string `locationName:"rootDeviceType" type:"string" enum:"DeviceType"`

	SourceDestCheck *bool `locationName:"sourceDestCheck" type:"boolean"`

	SpotInstanceRequestId *string `locationName:"spotInstanceRequestId" type:"string"`

	SriovNetSupport *string `locationName:"sriovNetSupport" type:"string"`

	StateReason *StateReason `locationName:"stateReason" type:"structure"`

	SubnetId *string `locationName:"subnetId" type:"string"`

	Tags []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`

	VirtualizationType *string `locationName:"virtualizationType" type:"string" enum:"VirtualizationType"`

	VpcId *string `locationName:"vpcId" type:"string"`
}

// String returns the string representation
func (s Instance) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Instance) GoString() string {
	return s.String()
}

// SetAmiLaunchIndex sets the AmiLaunchIndex field's value.
func (s *Instance) SetAmiLaunchIndex(v int64) *Instance {
	s.AmiLaunchIndex = &v
	return s
}

// SetArchitecture sets the Architecture field's value.
func (s *Instance) SetArchitecture(v string) *Instance {
	s.Architecture = &v
	return s
}

// SetBlockDeviceMappings sets the BlockDeviceMappings field's value.
func (s *Instance) SetBlockDeviceMappings(v []*InstanceBlockDeviceMapping) *Instance {
	s.BlockDeviceMappings = v
	return s
}

// SetClientToken sets the ClientToken field's value.
func (s *Instance) SetClientToken(v string) *Instance {
	s.ClientToken = &v
	return s
}

// SetEbsOptimized sets the EbsOptimized field's value.
func (s *Instance) SetEbsOptimized(v bool) *Instance {
	s.EbsOptimized = &v
	return s
}

// SetHypervisor sets the Hypervisor field's value.
func (s *Instance) SetHypervisor(v string) *Instance {
	s.Hypervisor = &v
	return s
}

// SetIamInstanceProfile sets the IamInstanceProfile field's value.
func (s *Instance) SetIamInstanceProfile(v *IamInstanceProfile) *Instance {
	s.IamInstanceProfile = v
	return s
}

// SetImageId sets the ImageId field's value.
func (s *Instance) SetImageId(v string) *Instance {
	s.ImageId = &v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *Instance) SetInstanceId(v string) *Instance {
	s.InstanceId = &v
	return s
}

// SetInstanceLifecycle sets the InstanceLifecycle field's value.
func (s *Instance) SetInstanceLifecycle(v string) *Instance {
	s.InstanceLifecycle = &v
	return s
}

// SetInstanceType sets the InstanceType field's value.
func (s *Instance) SetInstanceType(v string) *Instance {
	s.InstanceType = &v
	return s
}

// SetKernelId sets the KernelId field's value.
func (s *Instance) SetKernelId(v string) *Instance {
	s.KernelId = &v
	return s
}

// SetKeyName sets the KeyName field's value.
func (s *Instance) SetKeyName(v string) *Instance {
	s.KeyName = &v
	return s
}

// SetMonitoring sets the Monitoring field's value.
func (s *Instance) SetMonitoring(v *Monitoring) *Instance {
	s.Monitoring = v
	return s
}

// SetNetworkInterfaces sets the NetworkInterfaces field's value.
func (s *Instance) SetNetworkInterfaces(v []*InstanceNetworkInterface) *Instance {
	s.NetworkInterfaces = v
	return s
}

// SetPlacement sets the Placement field's value.
func (s *Instance) SetPlacement(v *Placement) *Instance {
	s.Placement = v
	return s
}

// SetPlatform sets the Platform field's value.
func (s *Instance) SetPlatform(v string) *Instance {
	s.Platform = &v
	return s
}

// SetPrivateDnsName sets the PrivateDnsName field's value.
func (s *Instance) SetPrivateDnsName(v string) *Instance {
	s.PrivateDnsName = &v
	return s
}

// SetPrivateIpAddress sets the PrivateIpAddress field's value.
func (s *Instance) SetPrivateIpAddress(v string) *Instance {
	s.PrivateIpAddress = &v
	return s
}

// SetProductCodes sets the ProductCodes field's value.
func (s *Instance) SetProductCodes(v []*ProductCode) *Instance {
	s.ProductCodes = v
	return s
}

// SetRamdiskId sets the RamdiskId field's value.
func (s *Instance) SetRamdiskId(v string) *Instance {
	s.RamdiskId = &v
	return s
}

// SetRootDeviceName sets the RootDeviceName field's value.
func (s *Instance) SetRootDeviceName(v string) *Instance {
	s.RootDeviceName = &v
	return s
}

// SetRootDeviceType sets the RootDeviceType field's value.
func (s *Instance) SetRootDeviceType(v string) *Instance {
	s.RootDeviceType = &v
	return s
}

// SetSourceDestCheck sets the SourceDestCheck field's value.
func (s *Instance) SetSourceDestCheck(v bool) *Instance {
	s.SourceDestCheck = &v
	return s
}

// SetSpotInstanceRequestId sets the SpotInstanceRequestId field's value.
func (s *Instance) SetSpotInstanceRequestId(v string) *Instance {
	s.SpotInstanceRequestId = &v
	return s
}

// SetSriovNetSupport sets the SriovNetSupport field's value.
func (s *Instance) SetSriovNetSupport(v string) *Instance {
	s.SriovNetSupport = &v
	return s
}

// SetStateReason sets the StateReason field's value.
func (s *Instance) SetStateReason(v *StateReason) *Instance {
	s.StateReason = v
	return s
}

// SetSubnetId sets the SubnetId field's value.
func (s *Instance) SetSubnetId(v string) *Instance {
	s.SubnetId = &v
	return s
}

// SetTags sets the Tags field's value.
func (s *Instance) SetTags(v []*Tag) *Instance {
	s.Tags = v
	return s
}

// SetVirtualizationType sets the VirtualizationType field's value.
func (s *Instance) SetVirtualizationType(v string) *Instance {
	s.VirtualizationType = &v
	return s
}

// SetVpcId sets the VpcId field's value.
func (s *Instance) SetVpcId(v string) *Instance {
	s.VpcId = &v
	return s
}

// InstanceBlockDeviceMapping struct
type InstanceBlockDeviceMapping struct {
	DeviceName *string `locationName:"deviceName" type:"string"`

	Ebs *EbsInstanceBlockDevice `locationName:"ebs" type:"structure"`
}

// String returns the string representation
func (s InstanceBlockDeviceMapping) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstanceBlockDeviceMapping) GoString() string {
	return s.String()
}

// SetDeviceName sets the DeviceName field's value.
func (s *InstanceBlockDeviceMapping) SetDeviceName(v string) *InstanceBlockDeviceMapping {
	s.DeviceName = &v
	return s
}

// SetEbs sets the Ebs field's value.
func (s *InstanceBlockDeviceMapping) SetEbs(v *EbsInstanceBlockDevice) *InstanceBlockDeviceMapping {
	s.Ebs = v
	return s
}

// InstanceBlockDeviceMappingSpecification struct
type InstanceBlockDeviceMappingSpecification struct {
	_ struct{} `type:"structure"`

	// The device name exposed to the instance (for example, /dev/sdh or xvdh).
	DeviceName *string `locationName:"deviceName" type:"string"`

	// Parameters used to automatically set up EBS volumes when the instance is
	// launched.
	Ebs *EbsInstanceBlockDeviceSpecification `locationName:"ebs" type:"structure"`

	// suppress the specified device included in the block device mapping.
	NoDevice *string `locationName:"noDevice" type:"string"`

	// The virtual device name.
	VirtualName *string `locationName:"virtualName" type:"string"`
}

// String returns the string representation
func (s InstanceBlockDeviceMappingSpecification) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstanceBlockDeviceMappingSpecification) GoString() string {
	return s.String()
}

// SetDeviceName sets the DeviceName field's value.
func (s *InstanceBlockDeviceMappingSpecification) SetDeviceName(v string) *InstanceBlockDeviceMappingSpecification {
	s.DeviceName = &v
	return s
}

// SetEbs sets the Ebs field's value.
func (s *InstanceBlockDeviceMappingSpecification) SetEbs(v *EbsInstanceBlockDeviceSpecification) *InstanceBlockDeviceMappingSpecification {
	s.Ebs = v
	return s
}

// SetNoDevice sets the NoDevice field's value.
func (s *InstanceBlockDeviceMappingSpecification) SetNoDevice(v string) *InstanceBlockDeviceMappingSpecification {
	s.NoDevice = &v
	return s
}

// SetVirtualName sets the VirtualName field's value.
func (s *InstanceBlockDeviceMappingSpecification) SetVirtualName(v string) *InstanceBlockDeviceMappingSpecification {
	s.VirtualName = &v
	return s
}

// InstanceCapacity struct
type InstanceCapacity struct {
	_ struct{} `type:"structure"`

	// The number of instances that can still be launched onto the Dedicated Host.
	AvailableCapacity *int64 `locationName:"availableCapacity" type:"integer"`

	// The instance type size supported by the Dedicated Host.
	InstanceType *string `locationName:"instanceType" type:"string"`

	// The total number of instances that can be launched onto the Dedicated Host.
	TotalCapacity *int64 `locationName:"totalCapacity" type:"integer"`
}

// String returns the string representation
func (s InstanceCapacity) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstanceCapacity) GoString() string {
	return s.String()
}

// SetAvailableCapacity sets the AvailableCapacity field's value.
func (s *InstanceCapacity) SetAvailableCapacity(v int64) *InstanceCapacity {
	s.AvailableCapacity = &v
	return s
}

// SetInstanceType sets the InstanceType field's value.
func (s *InstanceCapacity) SetInstanceType(v string) *InstanceCapacity {
	s.InstanceType = &v
	return s
}

// SetTotalCapacity sets the TotalCapacity field's value.
func (s *InstanceCapacity) SetTotalCapacity(v int64) *InstanceCapacity {
	s.TotalCapacity = &v
	return s
}

// InstanceCount struct
type InstanceCount struct {
	_ struct{} `type:"structure"`

	// The number of listed Reserved Instances in the state specified by the state.
	InstanceCount *int64 `locationName:"instanceCount" type:"integer"`

	// The states of the listed Reserved Instances.
	State *string `locationName:"state" type:"string" enum:"ListingState"`
}

// String returns the string representation
func (s InstanceCount) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstanceCount) GoString() string {
	return s.String()
}

// SetInstanceCount sets the InstanceCount field's value.
func (s *InstanceCount) SetInstanceCount(v int64) *InstanceCount {
	s.InstanceCount = &v
	return s
}

// SetState sets the State field's value.
func (s *InstanceCount) SetState(v string) *InstanceCount {
	s.State = &v
	return s
}

// InstanceExportDetails struct
type InstanceExportDetails struct {
	_ struct{} `type:"structure"`

	// The ID of the resource being exported.
	InstanceId *string `locationName:"instanceId" type:"string"`

	// The target virtualization environment.
	TargetEnvironment *string `locationName:"targetEnvironment" type:"string" enum:"ExportEnvironment"`
}

// String returns the string representation
func (s InstanceExportDetails) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstanceExportDetails) GoString() string {
	return s.String()
}

// SetInstanceId sets the InstanceId field's value.
func (s *InstanceExportDetails) SetInstanceId(v string) *InstanceExportDetails {
	s.InstanceId = &v
	return s
}

// SetTargetEnvironment sets the TargetEnvironment field's value.
func (s *InstanceExportDetails) SetTargetEnvironment(v string) *InstanceExportDetails {
	s.TargetEnvironment = &v
	return s
}

// InstanceMonitoring struct
type InstanceMonitoring struct {
	_ struct{} `type:"structure"`

	// The ID of the instance.
	InstanceId *string `locationName:"instanceId" type:"string"`

	// The monitoring for the instance.
	Monitoring *Monitoring `locationName:"monitoring" type:"structure"`
}

// String returns the string representation
func (s InstanceMonitoring) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstanceMonitoring) GoString() string {
	return s.String()
}

// SetInstanceId sets the InstanceId field's value.
func (s *InstanceMonitoring) SetInstanceId(v string) *InstanceMonitoring {
	s.InstanceId = &v
	return s
}

// SetMonitoring sets the Monitoring field's value.
func (s *InstanceMonitoring) SetMonitoring(v *Monitoring) *InstanceMonitoring {
	s.Monitoring = v
	return s
}

// InstanceNetworkInterface struct
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

	SourceDestCheck *bool `locationName:"sourceDestCheck" type:"boolean"`

	Status *string `locationName:"status" type:"string" enum:"NetworkInterfaceStatus"`

	SubnetId *string `locationName:"subnetId" type:"string"`

	VpcId *string `locationName:"vpcId" type:"string"`
}

// String returns the string representation
func (s InstanceNetworkInterface) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstanceNetworkInterface) GoString() string {
	return s.String()
}

// SetAssociation sets the Association field's value.
func (s *InstanceNetworkInterface) SetAssociation(v *InstanceNetworkInterfaceAssociation) *InstanceNetworkInterface {
	s.Association = v
	return s
}

// SetAttachment sets the Attachment field's value.
func (s *InstanceNetworkInterface) SetAttachment(v *InstanceNetworkInterfaceAttachment) *InstanceNetworkInterface {
	s.Attachment = v
	return s
}

// SetDescription sets the Description field's value.
func (s *InstanceNetworkInterface) SetDescription(v string) *InstanceNetworkInterface {
	s.Description = &v
	return s
}

// SetGroups sets the Groups field's value.
func (s *InstanceNetworkInterface) SetGroups(v []*GroupIdentifier) *InstanceNetworkInterface {
	s.Groups = v
	return s
}

// SetMacAddress sets the MacAddress field's value.
func (s *InstanceNetworkInterface) SetMacAddress(v string) *InstanceNetworkInterface {
	s.MacAddress = &v
	return s
}

// SetNetworkInterfaceId sets the NetworkInterfaceId field's value.
func (s *InstanceNetworkInterface) SetNetworkInterfaceId(v string) *InstanceNetworkInterface {
	s.NetworkInterfaceId = &v
	return s
}

// SetOwnerId sets the OwnerId field's value.
func (s *InstanceNetworkInterface) SetOwnerId(v string) *InstanceNetworkInterface {
	s.OwnerId = &v
	return s
}

// SetPrivateDnsName sets the PrivateDnsName field's value.
func (s *InstanceNetworkInterface) SetPrivateDnsName(v string) *InstanceNetworkInterface {
	s.PrivateDnsName = &v
	return s
}

// SetPrivateIpAddress sets the PrivateIpAddress field's value.
func (s *InstanceNetworkInterface) SetPrivateIpAddress(v string) *InstanceNetworkInterface {
	s.PrivateIpAddress = &v
	return s
}

// SetPrivateIpAddresses sets the PrivateIpAddresses field's value.
func (s *InstanceNetworkInterface) SetPrivateIpAddresses(v []*InstancePrivateIpAddress) *InstanceNetworkInterface {
	s.PrivateIpAddresses = v
	return s
}

// SetSourceDestCheck sets the SourceDestCheck field's value.
func (s *InstanceNetworkInterface) SetSourceDestCheck(v bool) *InstanceNetworkInterface {
	s.SourceDestCheck = &v
	return s
}

// SetStatus sets the Status field's value.
func (s *InstanceNetworkInterface) SetStatus(v string) *InstanceNetworkInterface {
	s.Status = &v
	return s
}

// SetSubnetId sets the SubnetId field's value.
func (s *InstanceNetworkInterface) SetSubnetId(v string) *InstanceNetworkInterface {
	s.SubnetId = &v
	return s
}

// SetVpcId sets the VpcId field's value.
func (s *InstanceNetworkInterface) SetVpcId(v string) *InstanceNetworkInterface {
	s.VpcId = &v
	return s
}

// InstanceNetworkInterfaceAssociation struct
type InstanceNetworkInterfaceAssociation struct {
	IpOwnerId *string `locationName:"ipOwnerId" type:"string"`

	PublicDnsName *string `locationName:"publicDnsName" type:"string"`

	PublicIp *string `locationName:"publicIp" type:"string"`
}

// String returns the string representation
func (s InstanceNetworkInterfaceAssociation) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstanceNetworkInterfaceAssociation) GoString() string {
	return s.String()
}

// SetIpOwnerId sets the IpOwnerId field's value.
func (s *InstanceNetworkInterfaceAssociation) SetIpOwnerId(v string) *InstanceNetworkInterfaceAssociation {
	s.IpOwnerId = &v
	return s
}

// SetPublicDnsName sets the PublicDnsName field's value.
func (s *InstanceNetworkInterfaceAssociation) SetPublicDnsName(v string) *InstanceNetworkInterfaceAssociation {
	s.PublicDnsName = &v
	return s
}

// SetPublicIp sets the PublicIp field's value.
func (s *InstanceNetworkInterfaceAssociation) SetPublicIp(v string) *InstanceNetworkInterfaceAssociation {
	s.PublicIp = &v
	return s
}

// InstanceNetworkInterfaceAttachment struct
type InstanceNetworkInterfaceAttachment struct {
	AttachmentId *string `locationName:"attachmentId" type:"string"`

	DeleteOnTermination *bool `locationName:"deleteOnTermination" type:"boolean"`

	DeviceIndex *int64 `locationName:"deviceIndex" type:"integer"`

	Status *string `locationName:"status" type:"string" enum:"AttachmentStatus"`
}

// String returns the string representation
func (s InstanceNetworkInterfaceAttachment) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstanceNetworkInterfaceAttachment) GoString() string {
	return s.String()
}

// SetAttachmentId sets the AttachmentId field's value.
func (s *InstanceNetworkInterfaceAttachment) SetAttachmentId(v string) *InstanceNetworkInterfaceAttachment {
	s.AttachmentId = &v
	return s
}

// SetDeleteOnTermination sets the DeleteOnTermination field's value.
func (s *InstanceNetworkInterfaceAttachment) SetDeleteOnTermination(v bool) *InstanceNetworkInterfaceAttachment {
	s.DeleteOnTermination = &v
	return s
}

// SetDeviceIndex sets the DeviceIndex field's value.
func (s *InstanceNetworkInterfaceAttachment) SetDeviceIndex(v int64) *InstanceNetworkInterfaceAttachment {
	s.DeviceIndex = &v
	return s
}

// SetStatus sets the Status field's value.
func (s *InstanceNetworkInterfaceAttachment) SetStatus(v string) *InstanceNetworkInterfaceAttachment {
	s.Status = &v
	return s
}

// InstanceNetworkInterfaceSpecification struct
type InstanceNetworkInterfaceSpecification struct {
	_ struct{} `type:"structure"`

	// Indicates whether to assign a public IPv4 address to an instance you launch
	// in a VPC. The public IP address can only be assigned to a network interface
	// for eth0, and can only be assigned to a new network interface, not an existing
	// one. You cannot specify more than one network interface in the request. If
	// launching into a default subnet, the default value is true.
	AssociatePublicIpAddress *bool `locationName:"associatePublicIpAddress" type:"boolean"`

	// If set to true, the interface is deleted when the instance is terminated.
	// You can specify true only if creating a new network interface when launching
	// an instance.
	DeleteOnTermination *bool `locationName:"deleteOnTermination" type:"boolean"`

	// The description of the network interface. Applies only if creating a network
	// interface when launching an instance.
	Description *string `locationName:"description" type:"string"`

	// The index of the device on the instance for the network interface attachment.
	// If you are specifying a network interface in a RunInstances request, you
	// must provide the device index.
	DeviceIndex *int64 `locationName:"deviceIndex" type:"integer"`

	// The IDs of the security groups for the network interface. Applies only if
	// creating a network interface when launching an instance.
	Groups []*string `locationName:"SecurityGroupId" locationNameList:"SecurityGroupId" type:"list"`

	// A number of IPv6 addresses to assign to the network interface. Amazon EC2
	// chooses the IPv6 addresses from the range of the subnet. You cannot specify
	// this option and the option to assign specific IPv6 addresses in the same
	// request. You can specify this option if you've specified a minimum number
	// of instances to launch.
	Ipv6AddressCount *int64 `locationName:"ipv6AddressCount" type:"integer"`

	// The ID of the network interface.
	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string"`

	// The private IPv4 address of the network interface. Applies only if creating
	// a network interface when launching an instance. You cannot specify this option
	// if you're launching more than one instance in a RunInstances request.
	PrivateIpAddress *string `locationName:"privateIpAddress" type:"string"`

	// One or more private IPv4 addresses to assign to the network interface. Only
	// one private IPv4 address can be designated as primary. You cannot specify
	// this option if you're launching more than one instance in a RunInstances
	// request.
	PrivateIpAddresses []*PrivateIpAddressSpecification `locationName:"privateIpAddressesSet" queryName:"PrivateIpAddresses" locationNameList:"item" type:"list"`

	// The number of secondary private IPv4 addresses. You can't specify this option
	// and specify more than one private IP address using the private IP addresses
	// option. You cannot specify this option if you're launching more than one
	// instance in a RunInstances request.
	SecondaryPrivateIpAddressCount *int64 `locationName:"secondaryPrivateIpAddressCount" type:"integer"`

	// The ID of the subnet associated with the network string. Applies only if
	// creating a network interface when launching an instance.
	SubnetId *string `locationName:"subnetId" type:"string"`
}

// String returns the string representation
func (s InstanceNetworkInterfaceSpecification) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstanceNetworkInterfaceSpecification) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *InstanceNetworkInterfaceSpecification) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "InstanceNetworkInterfaceSpecification"}
	if s.PrivateIpAddresses != nil {
		for i, v := range s.PrivateIpAddresses {
			if v == nil {
				continue
			}
			if err := v.Validate(); err != nil {
				invalidParams.AddNested(fmt.Sprintf("%s[%v]", "PrivateIpAddresses", i), err.(request.ErrInvalidParams))
			}
		}
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAssociatePublicIpAddress sets the AssociatePublicIpAddress field's value.
func (s *InstanceNetworkInterfaceSpecification) SetAssociatePublicIpAddress(v bool) *InstanceNetworkInterfaceSpecification {
	s.AssociatePublicIpAddress = &v
	return s
}

// SetDeleteOnTermination sets the DeleteOnTermination field's value.
func (s *InstanceNetworkInterfaceSpecification) SetDeleteOnTermination(v bool) *InstanceNetworkInterfaceSpecification {
	s.DeleteOnTermination = &v
	return s
}

// SetDescription sets the Description field's value.
func (s *InstanceNetworkInterfaceSpecification) SetDescription(v string) *InstanceNetworkInterfaceSpecification {
	s.Description = &v
	return s
}

// SetDeviceIndex sets the DeviceIndex field's value.
func (s *InstanceNetworkInterfaceSpecification) SetDeviceIndex(v int64) *InstanceNetworkInterfaceSpecification {
	s.DeviceIndex = &v
	return s
}

// SetGroups sets the Groups field's value.
func (s *InstanceNetworkInterfaceSpecification) SetGroups(v []*string) *InstanceNetworkInterfaceSpecification {
	s.Groups = v
	return s
}

// SetIpv6AddressCount sets the Ipv6AddressCount field's value.
func (s *InstanceNetworkInterfaceSpecification) SetIpv6AddressCount(v int64) *InstanceNetworkInterfaceSpecification {
	s.Ipv6AddressCount = &v
	return s
}

// SetNetworkInterfaceId sets the NetworkInterfaceId field's value.
func (s *InstanceNetworkInterfaceSpecification) SetNetworkInterfaceId(v string) *InstanceNetworkInterfaceSpecification {
	s.NetworkInterfaceId = &v
	return s
}

// SetPrivateIpAddress sets the PrivateIpAddress field's value.
func (s *InstanceNetworkInterfaceSpecification) SetPrivateIpAddress(v string) *InstanceNetworkInterfaceSpecification {
	s.PrivateIpAddress = &v
	return s
}

// SetPrivateIpAddresses sets the PrivateIpAddresses field's value.
func (s *InstanceNetworkInterfaceSpecification) SetPrivateIpAddresses(v []*PrivateIpAddressSpecification) *InstanceNetworkInterfaceSpecification {
	s.PrivateIpAddresses = v
	return s
}

// SetSecondaryPrivateIpAddressCount sets the SecondaryPrivateIpAddressCount field's value.
func (s *InstanceNetworkInterfaceSpecification) SetSecondaryPrivateIpAddressCount(v int64) *InstanceNetworkInterfaceSpecification {
	s.SecondaryPrivateIpAddressCount = &v
	return s
}

// SetSubnetId sets the SubnetId field's value.
func (s *InstanceNetworkInterfaceSpecification) SetSubnetId(v string) *InstanceNetworkInterfaceSpecification {
	s.SubnetId = &v
	return s
}

// InstancePrivateIpAddress struct
type InstancePrivateIpAddress struct {
	Association *InstanceNetworkInterfaceAssociation `locationName:"association" type:"structure"`

	Primary *bool `locationName:"primary" type:"boolean"`

	PrivateDnsName *string `locationName:"privateDnsName" type:"string"`

	PrivateIpAddress *string `locationName:"privateIpAddress" type:"string"`
}

// String returns the string representation
func (s InstancePrivateIpAddress) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstancePrivateIpAddress) GoString() string {
	return s.String()
}

// SetAssociation sets the Association field's value.
func (s *InstancePrivateIpAddress) SetAssociation(v *InstanceNetworkInterfaceAssociation) *InstancePrivateIpAddress {
	s.Association = v
	return s
}

// SetPrimary sets the Primary field's value.
func (s *InstancePrivateIpAddress) SetPrimary(v bool) *InstancePrivateIpAddress {
	s.Primary = &v
	return s
}

// SetPrivateDnsName sets the PrivateDnsName field's value.
func (s *InstancePrivateIpAddress) SetPrivateDnsName(v string) *InstancePrivateIpAddress {
	s.PrivateDnsName = &v
	return s
}

// SetPrivateIpAddress sets the PrivateIpAddress field's value.
func (s *InstancePrivateIpAddress) SetPrivateIpAddress(v string) *InstancePrivateIpAddress {
	s.PrivateIpAddress = &v
	return s
}

// InstanceState struct
type InstanceState struct {
	Code *int64 `locationName:"code" type:"integer"`

	Name *string `locationName:"name" type:"string" enum:"InstanceStateName"`
}

// String returns the string representation
func (s InstanceState) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstanceState) GoString() string {
	return s.String()
}

// SetCode sets the Code field's value.
func (s *InstanceState) SetCode(v int64) *InstanceState {
	s.Code = &v
	return s
}

// SetName sets the Name field's value.
func (s *InstanceState) SetName(v string) *InstanceState {
	s.Name = &v
	return s
}

// InstanceStateChange struct
type InstanceStateChange struct {
	_ struct{} `type:"structure"`

	// The current state of the instance.
	CurrentState *InstanceState `locationName:"currentState" type:"structure"`

	// The ID of the instance.
	InstanceId *string `locationName:"instanceId" type:"string"`

	// The previous state of the instance.
	PreviousState *InstanceState `locationName:"previousState" type:"structure"`
}

// String returns the string representation
func (s InstanceStateChange) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstanceStateChange) GoString() string {
	return s.String()
}

// SetCurrentState sets the CurrentState field's value.
func (s *InstanceStateChange) SetCurrentState(v *InstanceState) *InstanceStateChange {
	s.CurrentState = v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *InstanceStateChange) SetInstanceId(v string) *InstanceStateChange {
	s.InstanceId = &v
	return s
}

// SetPreviousState sets the PreviousState field's value.
func (s *InstanceStateChange) SetPreviousState(v *InstanceState) *InstanceStateChange {
	s.PreviousState = v
	return s
}

// Describes the status of an instance.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/InstanceStatus
type InstanceStatus struct {
	_ struct{} `type:"structure"`

	// The Availability Zone of the instance.
	AvailabilityZone *string `locationName:"availabilityZone" type:"string"`

	// Any scheduled events associated with the instance.
	Events []*InstanceStatusEvent `locationName:"eventsSet" locationNameList:"item" type:"list"`

	// The ID of the instance.
	InstanceId *string `locationName:"instanceId" type:"string"`

	// The intended state of the instance. DescribeInstanceStatus requires that
	// an instance be in the running state.
	InstanceState *InstanceState `locationName:"instanceState" type:"structure"`

	// Reports impaired functionality that stems from issues internal to the instance,
	// such as impaired reachability.
	InstanceStatus *InstanceStatusSummary `locationName:"instanceStatus" type:"structure"`

	// Reports impaired functionality that stems from issues related to the systems
	// that support an instance, such as hardware failures and network connectivity
	// problems.
	SystemStatus *InstanceStatusSummary `locationName:"systemStatus" type:"structure"`
}

// String returns the string representation
func (s InstanceStatus) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstanceStatus) GoString() string {
	return s.String()
}

// SetAvailabilityZone sets the AvailabilityZone field's value.
func (s *InstanceStatus) SetAvailabilityZone(v string) *InstanceStatus {
	s.AvailabilityZone = &v
	return s
}

// SetEvents sets the Events field's value.
func (s *InstanceStatus) SetEvents(v []*InstanceStatusEvent) *InstanceStatus {
	s.Events = v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *InstanceStatus) SetInstanceId(v string) *InstanceStatus {
	s.InstanceId = &v
	return s
}

// SetInstanceState sets the InstanceState field's value.
func (s *InstanceStatus) SetInstanceState(v *InstanceState) *InstanceStatus {
	s.InstanceState = v
	return s
}

// SetInstanceStatus sets the InstanceStatus field's value.
func (s *InstanceStatus) SetInstanceStatus(v *InstanceStatusSummary) *InstanceStatus {
	s.InstanceStatus = v
	return s
}

// SetSystemStatus sets the SystemStatus field's value.
func (s *InstanceStatus) SetSystemStatus(v *InstanceStatusSummary) *InstanceStatus {
	s.SystemStatus = v
	return s
}

// Describes the instance status.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/InstanceStatusDetails
type InstanceStatusDetails struct {
	_ struct{} `type:"structure"`

	// The time when a status check failed. For an instance that was launched and
	// impaired, this is the time when the instance was launched.
	ImpairedSince *time.Time `locationName:"impairedSince" type:"timestamp" timestampFormat:"iso8601"`

	// The type of instance status.
	Name *string `locationName:"name" type:"string" enum:"StatusName"`

	// The status.
	Status *string `locationName:"status" type:"string" enum:"StatusType"`
}

// String returns the string representation
func (s InstanceStatusDetails) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstanceStatusDetails) GoString() string {
	return s.String()
}

// SetImpairedSince sets the ImpairedSince field's value.
func (s *InstanceStatusDetails) SetImpairedSince(v time.Time) *InstanceStatusDetails {
	s.ImpairedSince = &v
	return s
}

// SetName sets the Name field's value.
func (s *InstanceStatusDetails) SetName(v string) *InstanceStatusDetails {
	s.Name = &v
	return s
}

// SetStatus sets the Status field's value.
func (s *InstanceStatusDetails) SetStatus(v string) *InstanceStatusDetails {
	s.Status = &v
	return s
}

// Describes a scheduled event for an instance.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/InstanceStatusEvent
type InstanceStatusEvent struct {
	_ struct{} `type:"structure"`

	// The event code.
	Code *string `locationName:"code" type:"string" enum:"EventCode"`

	// A description of the event.
	//
	// After a scheduled event is completed, it can still be described for up to
	// a week. If the event has been completed, this description starts with the
	// following text: [Completed].
	Description *string `locationName:"description" type:"string"`

	// The latest scheduled end time for the event.
	NotAfter *time.Time `locationName:"notAfter" type:"timestamp" timestampFormat:"iso8601"`

	// The earliest scheduled start time for the event.
	NotBefore *time.Time `locationName:"notBefore" type:"timestamp" timestampFormat:"iso8601"`
}

// String returns the string representation
func (s InstanceStatusEvent) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstanceStatusEvent) GoString() string {
	return s.String()
}

// SetCode sets the Code field's value.
func (s *InstanceStatusEvent) SetCode(v string) *InstanceStatusEvent {
	s.Code = &v
	return s
}

// SetDescription sets the Description field's value.
func (s *InstanceStatusEvent) SetDescription(v string) *InstanceStatusEvent {
	s.Description = &v
	return s
}

// SetNotAfter sets the NotAfter field's value.
func (s *InstanceStatusEvent) SetNotAfter(v time.Time) *InstanceStatusEvent {
	s.NotAfter = &v
	return s
}

// SetNotBefore sets the NotBefore field's value.
func (s *InstanceStatusEvent) SetNotBefore(v time.Time) *InstanceStatusEvent {
	s.NotBefore = &v
	return s
}

// Describes the status of an instance.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/InstanceStatusSummary
type InstanceStatusSummary struct {
	_ struct{} `type:"structure"`

	// The system instance health or application instance health.
	Details []*InstanceStatusDetails `locationName:"details" locationNameList:"item" type:"list"`

	// The status.
	Status *string `locationName:"status" type:"string" enum:"SummaryStatus"`
}

// String returns the string representation
func (s InstanceStatusSummary) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InstanceStatusSummary) GoString() string {
	return s.String()
}

// SetDetails sets the Details field's value.
func (s *InstanceStatusSummary) SetDetails(v []*InstanceStatusDetails) *InstanceStatusSummary {
	s.Details = v
	return s
}

// SetStatus sets the Status field's value.
func (s *InstanceStatusSummary) SetStatus(v string) *InstanceStatusSummary {
	s.Status = &v
	return s
}

// EbsInstanceBlockDevice struct
type EbsInstanceBlockDevice struct {
	AttachTime *time.Time `locationName:"attachTime" type:"timestamp" timestampFormat:"iso8601"`

	DeleteOnTermination *bool `locationName:"deleteOnTermination" type:"boolean"`

	Status *string `locationName:"status" type:"string" enum:"AttachmentStatus"`

	VolumeId *string `locationName:"volumeId" type:"string"`
}

// String returns the string representation
func (s EbsInstanceBlockDevice) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s EbsInstanceBlockDevice) GoString() string {
	return s.String()
}

// SetAttachTime sets the AttachTime field's value.
func (s *EbsInstanceBlockDevice) SetAttachTime(v time.Time) *EbsInstanceBlockDevice {
	s.AttachTime = &v
	return s
}

// SetDeleteOnTermination sets the DeleteOnTermination field's value.
func (s *EbsInstanceBlockDevice) SetDeleteOnTermination(v bool) *EbsInstanceBlockDevice {
	s.DeleteOnTermination = &v
	return s
}

// SetStatus sets the Status field's value.
func (s *EbsInstanceBlockDevice) SetStatus(v string) *EbsInstanceBlockDevice {
	s.Status = &v
	return s
}

// SetVolumeId sets the VolumeId field's value.
func (s *EbsInstanceBlockDevice) SetVolumeId(v string) *EbsInstanceBlockDevice {
	s.VolumeId = &v
	return s
}

// Describes information used to set up an EBS volume specified in a block device
// mapping.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/EbsInstanceBlockDeviceSpecification
type EbsInstanceBlockDeviceSpecification struct {
	_ struct{} `type:"structure"`

	// Indicates whether the volume is deleted on instance termination.
	DeleteOnTermination *bool `locationName:"deleteOnTermination" type:"boolean"`

	// The ID of the EBS volume.
	VolumeId *string `locationName:"volumeId" type:"string"`
}

// String returns the string representation
func (s EbsInstanceBlockDeviceSpecification) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s EbsInstanceBlockDeviceSpecification) GoString() string {
	return s.String()
}

// SetDeleteOnTermination sets the DeleteOnTermination field's value.
func (s *EbsInstanceBlockDeviceSpecification) SetDeleteOnTermination(v bool) *EbsInstanceBlockDeviceSpecification {
	s.DeleteOnTermination = &v
	return s
}

// SetVolumeId sets the VolumeId field's value.
func (s *EbsInstanceBlockDeviceSpecification) SetVolumeId(v string) *EbsInstanceBlockDeviceSpecification {
	s.VolumeId = &v
	return s
}

// IamInstanceProfile struct
type IamInstanceProfile struct {
	Arn *string `locationName:"arn" type:"string"`

	Id *string `locationName:"id" type:"string"`
}

// String returns the string representation
func (s IamInstanceProfile) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s IamInstanceProfile) GoString() string {
	return s.String()
}

// SetArn sets the Arn field's value.
func (s *IamInstanceProfile) SetArn(v string) *IamInstanceProfile {
	s.Arn = &v
	return s
}

// SetId sets the Id field's value.
func (s *IamInstanceProfile) SetId(v string) *IamInstanceProfile {
	s.Id = &v
	return s
}

// Describes the monitoring of an instance.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/Monitoring
type Monitoring struct {
	_ struct{} `type:"structure"`

	// Indicates whether detailed monitoring is enabled. Otherwise, basic monitoring
	// is enabled.
	State *string `locationName:"state" type:"string" enum:"MonitoringState"`
}

// String returns the string representation
func (s Monitoring) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Monitoring) GoString() string {
	return s.String()
}

// SetState sets the State field's value.
func (s *Monitoring) SetState(v string) *Monitoring {
	s.State = &v
	return s
}

// Placement struct
type Placement struct {
	Affinity *string `locationName:"affinity" type:"string"`

	AvailabilityZone *string `locationName:"availabilityZone" type:"string"`

	GroupName *string `locationName:"groupName" type:"string"`

	HostId *string `locationName:"hostId" type:"string"`

	Tenancy *string `locationName:"tenancy" type:"string" enum:"Tenancy"`
}

// String returns the string representation
func (s Placement) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Placement) GoString() string {
	return s.String()
}

// SetAffinity sets the Affinity field's value.
func (s *Placement) SetAffinity(v string) *Placement {
	s.Affinity = &v
	return s
}

// SetAvailabilityZone sets the AvailabilityZone field's value.
func (s *Placement) SetAvailabilityZone(v string) *Placement {
	s.AvailabilityZone = &v
	return s
}

// SetGroupName sets the GroupName field's value.
func (s *Placement) SetGroupName(v string) *Placement {
	s.GroupName = &v
	return s
}

// SetHostId sets the HostId field's value.
func (s *Placement) SetHostId(v string) *Placement {
	s.HostId = &v
	return s
}

// SetTenancy sets the Tenancy field's value.
func (s *Placement) SetTenancy(v string) *Placement {
	s.Tenancy = &v
	return s
}

// ProductCode struct.
type ProductCode struct {
	ProductCode *string `locationName:"productCode" type:"string"`

	Type *string `locationName:"type" type:"string" enum:"ProductCodeValues"`
}

// String returns the string representation
func (s ProductCode) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ProductCode) GoString() string {
	return s.String()
}

// SetProductCode sets the ProductCode field's value.
func (s *ProductCode) SetProductCode(v string) *ProductCode {
	s.ProductCode = &v
	return s
}

// SetType sets the Type field's value.
func (s *ProductCode) SetType(v string) *ProductCode {
	s.Type = &v
	return s
}

// StateReason struct
type StateReason struct {
	Code    *string `locationName:"code" type:"string"`
	Message *string `locationName:"message" type:"string"`
}

// String returns the string representation
func (s StateReason) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s StateReason) GoString() string {
	return s.String()
}

// SetCode sets the Code field's value.
func (s *StateReason) SetCode(v string) *StateReason {
	s.Code = &v
	return s
}

// SetMessage sets the Message field's value.
func (s *StateReason) SetMessage(v string) *StateReason {
	s.Message = &v
	return s
}

// Describes a tag.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/Tag
type Tag struct {
	_ struct{} `type:"structure"`

	// The key of the tag.
	//
	// Constraints: Tag keys are case-sensitive and accept a maximum of 127 Unicode
	// characters. May not begin with aws:
	Key *string `locationName:"key" type:"string"`

	// The value of the tag.
	//
	// Constraints: Tag values are case-sensitive and accept a maximum of 255 Unicode
	// characters.
	Value *string `locationName:"value" type:"string"`
}

// String returns the string representation
func (s Tag) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Tag) GoString() string {
	return s.String()
}

// SetKey sets the Key field's value.
func (s *Tag) SetKey(v string) *Tag {
	s.Key = &v
	return s
}

// SetValue sets the Value field's value.
func (s *Tag) SetValue(v string) *Tag {
	s.Value = &v
	return s
}

// Describes a secondary private IPv4 address for a network interface.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/PrivateIpAddressSpecification
type PrivateIpAddressSpecification struct {
	_ struct{} `type:"structure"`

	// Indicates whether the private IPv4 address is the primary private IPv4 address.
	// Only one IPv4 address can be designated as primary.
	Primary *bool `locationName:"primary" type:"boolean"`

	// The private IPv4 addresses.
	//
	// PrivateIpAddress is a required field
	PrivateIpAddress *string `locationName:"privateIpAddress" type:"string" required:"true"`
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

// Contains the parameters for DescribeInstanceAttribute.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeInstanceAttributeRequest
type DescribeInstanceAttributeInput struct {
	_ struct{} `type:"structure"`

	// The instance attribute.
	//
	// Note: The enaSupport attribute is not supported at this time.
	//
	// Attribute is a required field
	Attribute *string `locationName:"attribute" type:"string" required:"true" enum:"InstanceAttributeName"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the instance.
	//
	// InstanceId is a required field
	InstanceId *string `locationName:"instanceId" type:"string" required:"true"`
}

// String returns the string representation
func (s DescribeInstanceAttributeInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeInstanceAttributeInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DescribeInstanceAttributeInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DescribeInstanceAttributeInput"}
	if s.Attribute == nil {
		invalidParams.Add(request.NewErrParamRequired("Attribute"))
	}
	if s.InstanceId == nil {
		invalidParams.Add(request.NewErrParamRequired("InstanceId"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAttribute sets the Attribute field's value.
func (s *DescribeInstanceAttributeInput) SetAttribute(v string) *DescribeInstanceAttributeInput {
	s.Attribute = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeInstanceAttributeInput) SetDryRun(v bool) *DescribeInstanceAttributeInput {
	s.DryRun = &v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *DescribeInstanceAttributeInput) SetInstanceId(v string) *DescribeInstanceAttributeInput {
	s.InstanceId = &v
	return s
}

// Describes an instance attribute.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/InstanceAttribute
type DescribeInstanceAttributeOutput struct {
	_ struct{} `type:"structure"`

	// The block device mapping of the instance.
	BlockDeviceMappings []*InstanceBlockDeviceMapping `locationName:"blockDeviceMapping" locationNameList:"item" type:"list"`

	// If the value is true, you can't terminate the instance through the Amazon
	// EC2 console, CLI, or API; otherwise, you can.
	DisableApiTermination *AttributeBooleanValue `locationName:"disableApiTermination" type:"structure"`

	// Indicates whether the instance is optimized for EBS I/O.
	EbsOptimized *AttributeBooleanValue `locationName:"ebsOptimized" type:"structure"`

	// Indicates whether enhanced networking with ENA is enabled.
	EnaSupport *AttributeBooleanValue `locationName:"enaSupport" type:"structure"`

	// The security groups associated with the instance.
	Groups []*GroupIdentifier `locationName:"groupSet" locationNameList:"item" type:"list"`

	// The ID of the instance.
	InstanceId *string `locationName:"instanceId" type:"string"`

	// Indicates whether an instance stops or terminates when you initiate shutdown
	// from the instance (using the operating system command for system shutdown).
	InstanceInitiatedShutdownBehavior *AttributeValue `locationName:"instanceInitiatedShutdownBehavior" type:"structure"`

	// The instance type.
	InstanceType *AttributeValue `locationName:"instanceType" type:"structure"`

	// The kernel ID.
	KernelId *AttributeValue `locationName:"kernel" type:"structure"`

	// A list of product codes.
	ProductCodes []*ProductCode `locationName:"productCodes" locationNameList:"item" type:"list"`

	// The RAM disk ID.
	RamdiskId *AttributeValue `locationName:"ramdisk" type:"structure"`

	// The name of the root device (for example, /dev/sda1 or /dev/xvda).
	RootDeviceName *AttributeValue `locationName:"rootDeviceName" type:"structure"`

	// Indicates whether source/destination checking is enabled. A value of true
	// means checking is enabled, and false means checking is disabled. This value
	// must be false for a NAT instance to perform NAT.
	SourceDestCheck *AttributeBooleanValue `locationName:"sourceDestCheck" type:"structure"`

	// Indicates whether enhanced networking with the Intel 82599 Virtual Function
	// interface is enabled.
	SriovNetSupport *AttributeValue `locationName:"sriovNetSupport" type:"structure"`

	// The user data.
	UserData *AttributeValue `locationName:"userData" type:"structure"`
}

// String returns the string representation
func (s DescribeInstanceAttributeOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeInstanceAttributeOutput) GoString() string {
	return s.String()
}

// SetBlockDeviceMappings sets the BlockDeviceMappings field's value.
func (s *DescribeInstanceAttributeOutput) SetBlockDeviceMappings(v []*InstanceBlockDeviceMapping) *DescribeInstanceAttributeOutput {
	s.BlockDeviceMappings = v
	return s
}

// SetDisableApiTermination sets the DisableApiTermination field's value.
func (s *DescribeInstanceAttributeOutput) SetDisableApiTermination(v *AttributeBooleanValue) *DescribeInstanceAttributeOutput {
	s.DisableApiTermination = v
	return s
}

// SetEbsOptimized sets the EbsOptimized field's value.
func (s *DescribeInstanceAttributeOutput) SetEbsOptimized(v *AttributeBooleanValue) *DescribeInstanceAttributeOutput {
	s.EbsOptimized = v
	return s
}

// SetEnaSupport sets the EnaSupport field's value.
func (s *DescribeInstanceAttributeOutput) SetEnaSupport(v *AttributeBooleanValue) *DescribeInstanceAttributeOutput {
	s.EnaSupport = v
	return s
}

// SetGroups sets the Groups field's value.
func (s *DescribeInstanceAttributeOutput) SetGroups(v []*GroupIdentifier) *DescribeInstanceAttributeOutput {
	s.Groups = v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *DescribeInstanceAttributeOutput) SetInstanceId(v string) *DescribeInstanceAttributeOutput {
	s.InstanceId = &v
	return s
}

// SetInstanceInitiatedShutdownBehavior sets the InstanceInitiatedShutdownBehavior field's value.
func (s *DescribeInstanceAttributeOutput) SetInstanceInitiatedShutdownBehavior(v *AttributeValue) *DescribeInstanceAttributeOutput {
	s.InstanceInitiatedShutdownBehavior = v
	return s
}

// SetInstanceType sets the InstanceType field's value.
func (s *DescribeInstanceAttributeOutput) SetInstanceType(v *AttributeValue) *DescribeInstanceAttributeOutput {
	s.InstanceType = v
	return s
}

// SetKernelId sets the KernelId field's value.
func (s *DescribeInstanceAttributeOutput) SetKernelId(v *AttributeValue) *DescribeInstanceAttributeOutput {
	s.KernelId = v
	return s
}

// SetProductCodes sets the ProductCodes field's value.
func (s *DescribeInstanceAttributeOutput) SetProductCodes(v []*ProductCode) *DescribeInstanceAttributeOutput {
	s.ProductCodes = v
	return s
}

// SetRamdiskId sets the RamdiskId field's value.
func (s *DescribeInstanceAttributeOutput) SetRamdiskId(v *AttributeValue) *DescribeInstanceAttributeOutput {
	s.RamdiskId = v
	return s
}

// SetRootDeviceName sets the RootDeviceName field's value.
func (s *DescribeInstanceAttributeOutput) SetRootDeviceName(v *AttributeValue) *DescribeInstanceAttributeOutput {
	s.RootDeviceName = v
	return s
}

// SetSourceDestCheck sets the SourceDestCheck field's value.
func (s *DescribeInstanceAttributeOutput) SetSourceDestCheck(v *AttributeBooleanValue) *DescribeInstanceAttributeOutput {
	s.SourceDestCheck = v
	return s
}

// SetSriovNetSupport sets the SriovNetSupport field's value.
func (s *DescribeInstanceAttributeOutput) SetSriovNetSupport(v *AttributeValue) *DescribeInstanceAttributeOutput {
	s.SriovNetSupport = v
	return s
}

// SetUserData sets the UserData field's value.
func (s *DescribeInstanceAttributeOutput) SetUserData(v *AttributeValue) *DescribeInstanceAttributeOutput {
	s.UserData = v
	return s
}

// Describes a value for a resource attribute that is a Boolean value.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/AttributeBooleanValue
type AttributeBooleanValue struct {
	_ struct{} `type:"structure"`

	// The attribute value. The valid values are true or false.
	Value *bool `locationName:"value" type:"boolean"`
}

// String returns the string representation
func (s AttributeBooleanValue) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AttributeBooleanValue) GoString() string {
	return s.String()
}

// SetValue sets the Value field's value.
func (s *AttributeBooleanValue) SetValue(v bool) *AttributeBooleanValue {
	s.Value = &v
	return s
}

// Describes a value for a resource attribute that is a String.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/AttributeValue
type AttributeValue struct {
	_ struct{} `type:"structure"`

	// The attribute value. Note that the value is case-sensitive.
	Value *string `locationName:"value" type:"string"`
}

// String returns the string representation
func (s AttributeValue) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AttributeValue) GoString() string {
	return s.String()
}

// SetValue sets the Value field's value.
func (s *AttributeValue) SetValue(v string) *AttributeValue {
	s.Value = &v
	return s
}
