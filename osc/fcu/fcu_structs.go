package fcu

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/awsutil"
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

// Filter can be used to match a set of resources by various criteria.
type Filter struct {
	Name *string `type:"string"`

	Values []*string `locationName:"Value" locationNameList:"item" type:"list"`
}

// DescribeInstancesOutput struct
type DescribeInstancesOutput struct {
	_ struct{} `type:"structure"`

	NextToken *string `locationName:"nextToken" type:"string"`

	OwnerId *string `locationName:"ownerId" locationNameList:"item" type:"string"`

	RequesterId *string `locationName:"requesterId" locationNameList:"item" type:"string"`

	ReservationId *string `locationName:"reservationId" locationNameList:"item" type:"string"`

	Reservations []*Reservation `locationName:"reservationSet" locationNameList:"item" type:"list"`

	GroupSet []*GroupIdentifier `locationName:"groupSet" locationNameList:"item" type:"list"`
}

type GroupIdentifier struct {
	_ struct{} `type:"structure"`

	// The ID of the security group.
	GroupId *string `locationName:"groupId" type:"string"`

	// The name of the security group.
	GroupName *string `locationName:"groupName" type:"string"`
}

// Describes a reservation.
// See also, https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/Reservation
type Reservation struct {
	_ struct{} `type:"structure"`

	// [EC2-Classic only] One or more security groups.
	Groups []*GroupIdentifier `locationName:"groupSet" locationNameList:"item" type:"list"`

	// One or more instances.
	Instances []*Instance `locationName:"instancesSet" locationNameList:"item" type:"list"`

	// The ID of the AWS account that owns the reservation.
	OwnerId *string `locationName:"ownerId" type:"string"`

	// The ID of the requester that launched the instances on your behalf (for example,
	// AWS Management Console or Auto Scaling).
	RequesterId *string `locationName:"requesterId" type:"string"`

	// The ID of the reservation.
	ReservationId *string `locationName:"reservationId" type:"string"`
}

// Instance struct
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

// InstanceBlockDeviceMapping struct
type InstanceBlockDeviceMapping struct {
	DeviceName *string `locationName:"deviceName" type:"string"`

	Ebs *EbsInstanceBlockDevice `locationName:"ebs" type:"structure"`
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

// InstanceCount struct
type InstanceCount struct {
	_ struct{} `type:"structure"`

	// The number of listed Reserved Instances in the state specified by the state.
	InstanceCount *int64 `locationName:"instanceCount" type:"integer"`

	// The states of the listed Reserved Instances.
	State *string `locationName:"state" type:"string" enum:"ListingState"`
}

// InstanceExportDetails struct
type InstanceExportDetails struct {
	_ struct{} `type:"structure"`

	// The ID of the resource being exported.
	InstanceId *string `locationName:"instanceId" type:"string"`

	// The target virtualization environment.
	TargetEnvironment *string `locationName:"targetEnvironment" type:"string" enum:"ExportEnvironment"`
}

// InstanceMonitoring struct
type InstanceMonitoring struct {
	_ struct{} `type:"structure"`

	// The ID of the instance.
	InstanceId *string `locationName:"instanceId" type:"string"`

	// The monitoring for the instance.
	Monitoring *Monitoring `locationName:"monitoring" type:"structure"`
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

// InstanceNetworkInterfaceAssociation struct
type InstanceNetworkInterfaceAssociation struct {
	IpOwnerId *string `locationName:"ipOwnerId" type:"string"`

	PublicDnsName *string `locationName:"publicDnsName" type:"string"`

	PublicIp *string `locationName:"publicIp" type:"string"`
}

// InstanceNetworkInterfaceAttachment struct
type InstanceNetworkInterfaceAttachment struct {
	AttachmentId *string `locationName:"attachmentId" type:"string"`

	DeleteOnTermination *bool `locationName:"deleteOnTermination" type:"boolean"`

	DeviceIndex *int64 `locationName:"deviceIndex" type:"integer"`

	Status *string `locationName:"status" type:"string" enum:"AttachmentStatus"`
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

// InstancePrivateIpAddress struct
type InstancePrivateIpAddress struct {
	Association *InstanceNetworkInterfaceAssociation `locationName:"association" type:"structure"`

	Primary *bool `locationName:"primary" type:"boolean"`

	PrivateDnsName *string `locationName:"privateDnsName" type:"string"`

	PrivateIpAddress *string `locationName:"privateIpAddress" type:"string"`
}

// InstanceState struct
type InstanceState struct {
	Code *int64 `locationName:"code" type:"integer"`

	Name *string `locationName:"name" type:"string" enum:"InstanceStateName"`
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

// IamInstanceProfile struct
type IamInstanceProfile struct {
	Arn *string `locationName:"arn" type:"string"`

	Id *string `locationName:"id" type:"string"`
}

// Describes the monitoring of an instance.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/Monitoring
type Monitoring struct {
	_ struct{} `type:"structure"`

	// Indicates whether detailed monitoring is enabled. Otherwise, basic monitoring
	// is enabled.
	State *string `locationName:"state" type:"string" enum:"MonitoringState"`
}

// Placement struct
type Placement struct {
	Affinity *string `locationName:"affinity" type:"string"`

	AvailabilityZone *string `locationName:"availabilityZone" type:"string"`

	GroupName *string `locationName:"groupName" type:"string"`

	HostId *string `locationName:"hostId" type:"string"`

	Tenancy *string `locationName:"tenancy" type:"string" enum:"Tenancy"`
}

// ProductCode struct.
type ProductCode struct {
	ProductCode *string `locationName:"productCode" type:"string"`

	Type *string `locationName:"type" type:"string" enum:"ProductCodeValues"`
}

// StateReason struct
type StateReason struct {
	Code    *string `locationName:"code" type:"string"`
	Message *string `locationName:"message" type:"string"`
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

// Describes a value for a resource attribute that is a Boolean value.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/AttributeBooleanValue
type AttributeBooleanValue struct {
	_ struct{} `type:"structure"`

	// The attribute value. The valid values are true or false.
	Value *bool `locationName:"value" type:"boolean"`
}

// Describes a value for a resource attribute that is a String.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/AttributeValue
type AttributeValue struct {
	_ struct{} `type:"structure"`

	// The attribute value. Note that the value is case-sensitive.
	Value *string `locationName:"value" type:"string"`
}

//RunInstancesInput is the specification to run the an instance
type RunInstancesInput struct {
	_ struct{} `type:"structure"`

	// One or more block device mapping entries. You can't specify both a snapshot
	// Id and an encryption value. This is because only blank volumes can be encrypted
	// on creation. If a snapshot is the basis for a volume, it is not blank and
	// its encryption status is used for the volume encryption status.
	BlockDeviceMappings []*BlockDeviceMapping `locationName:"BlockDeviceMapping" locationNameList:"BlockDeviceMapping" type:"list"`

	// Unique, case-sensitive identifier you provide to ensure the idempotency of
	// the request. For more information, see Ensuring Idempotency (http://docs.aws.amazon.com/AWSEC2/latest/APIReference/Run_Instance_Idempotency.html).
	//
	// Constraints: Maximum 64 ASCII characters
	ClientToken *string `locationName:"clientToken" type:"string"`

	// If you set this parameter to true, you can't terminate the instance using
	// the Amazon EC2 console, CLI, or API; otherwise, you can. To change this attribute
	// to false after launch, use ModifyInstanceAttribute. Alternatively, if you
	// set InstanceInitiatedShutdownBehavior to terminate, you can terminate the
	// instance by running the shutdown command from the instance.
	//
	// Default: false
	DisableApiTermination *bool `locationName:"disableApiTermination" type:"boolean"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// Indicates whether the instance is optimized for Amazon EBS I/O. This optimization
	// provides dedicated throughput to Amazon EBS and an optimized configuration
	// stack to provide optimal Amazon EBS I/O performance. This optimization isn't
	// available with all instance types. Additional usage charges apply when using
	// an EBS-optimized instance.
	//
	// Default: false
	EbsOptimized *bool `locationName:"ebsOptimized" type:"boolean"`

	// The Id of the AMI, which you can get by calling DescribeImages. An AMI is
	// required to launch an instance and must be specified here or in a launch
	// template.
	ImageId *string `type:"string"`

	// Indicates whether an instance stops or terminates when you initiate shutdown
	// from the instance (using the operating system command for system shutdown).
	//
	// Default: stop
	InstanceInitiatedShutdownBehavior *bool `locationName:"instanceInitiatedShutdownBehavior" type:"bool" enum:"ShutdownBehavior"`

	// The instance type. For more information, see Instance Types (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html)
	// in the Amazon Elastic Compute Cloud User Guide.
	//
	// Default: m1.small
	InstanceType *string `type:"string" enum:"InstanceType"`

	// The name of the key pair. You can create a key pair using CreateKeyPair or
	// ImportKeyPair.
	//
	// If you do not specify a key pair, you can't connect to the instance unless
	// you choose an AMI that is configured to allow users another way to log in.
	KeyName *string `locationName:"keyName" type:"string"`

	// The maximum number of instances to launch. If you specify more instances
	// than Amazon EC2 can launch in the target Availability Zone, Amazon EC2 launches
	// the largest possible number of instances above MinCount.
	//
	// Constraints: Between 1 and the maximum number you're allowed for the specified
	// instance type. For more information about the default limits, and how to
	// request an increase, see How many instances can I run in Amazon EC2 (http://aws.amazon.com/ec2/faqs/#How_many_instances_can_I_run_in_Amazon_EC2)
	// in the Amazon EC2 FAQ.
	//
	// MaxCount is a required field
	MaxCount *int64 `type:"integer" required:"true"`

	// The minimum number of instances to launch. If you specify a minimum that
	// is more instances than Amazon EC2 can launch in the target Availability Zone,
	// Amazon EC2 launches no instances.
	//
	// Constraints: Between 1 and the maximum number you're allowed for the specified
	// instance type. For more information about the default limits, and how to
	// request an increase, see How many instances can I run in Amazon EC2 (http://aws.amazon.com/ec2/faqs/#How_many_instances_can_I_run_in_Amazon_EC2)
	// in the Amazon EC2 General FAQ.
	//
	// MinCount is a required field
	MinCount *int64 `type:"integer" required:"true"`

	// One or more network interfaces.
	NetworkInterfaces []*InstanceNetworkInterfaceSpecification `locationName:"networkInterface" locationNameList:"item" type:"list"`

	// The placement for the instance.
	Placement *Placement `type:"structure"`

	// [EC2-VPC] The primary IPv4 address. You must specify a value from the IPv4
	// address range of the subnet.
	//
	// Only one private IP address can be designated as primary. You can't specify
	// this option if you've specified the option to designate a private IP address
	// as the primary IP address in a network interface specification. You cannot
	// specify this option if you're launching more than one instance in the request.
	PrivateIPAddress *string `locationName:"privateIpAddress" type:"string"`

	// The Id of the RAM disk.
	//
	// We recommend that you use PV-GRUB instead of kernels and RAM disks. For more
	// information, see  PV-GRUB (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/UserProvidedkernels.html)
	// in the Amazon Elastic Compute Cloud User Guide.
	RamdiskId *string `type:"string"`

	// One or more security group Ids. You can create a security group using CreateSecurityGroup.
	//
	// Default: Amazon EC2 uses the default security group.
	SecurityGroupIds []*string `locationName:"SecurityGroupId" locationNameList:"SecurityGroupId" type:"list"`

	// [EC2-Classic, default VPC] One or more security group names. For a nondefault
	// VPC, you must use security group Ids instead.
	//
	// Default: Amazon EC2 uses the default security group.
	SecurityGroups []*string `locationName:"SecurityGroup" locationNameList:"SecurityGroup" type:"list"`

	// [EC2-VPC] The Id of the subnet to launch the instance into.
	SubnetId *string `type:"string"`

	// The user data to make available to the instance. For more information, see
	// Running Commands on Your Linux Instance at Launch (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html)
	// (Linux) and Adding User Data (http://docs.aws.amazon.com/AWSEC2/latest/WindowsGuide/ec2-instance-metadata.html#instancedata-add-user-data)
	// (Windows). If you are using a command line tool, base64-encoding is performed
	// for you, and you can load the text from a file. Otherwise, you must provide
	// base64-encoded text.
	UserData *string `type:"string"`

	OwnerId *string `type:"string"`

	RequesterId *string `type:"string"`

	ReservationId *string `type:"string"`

	PasswordData *string `type:"string"`
}

//BlockDeviceMapping input to specify the mapping
type BlockDeviceMapping struct {
	_ struct{} `type:"structure"`

	// The device name (for example, /dev/sdh or xvdh).
	DeviceName *string `locationName:"deviceName" type:"string"`

	Ebs *EbsBlockDevice `locationName:"ebs" type:"structure"`

	// Suppresses the specified device included in the block device mapping of the
	// AMI.
	NoDevice *string `locationName:"noDevice" type:"string"`

	// The virtual device name (ephemeralN). Instance store volumes are numbered
	// starting from 0. An instance type with 2 available instance store volumes
	// can specify mappings for ephemeral0 and ephemeral1.The number of available
	// instance store volumes depends on the instance type. After you connect to
	// the instance, you must mount the volume.
	//
	// Constraints: For M3 instances, you must specify instance store volumes in
	// the block device mapping for the instance. When you launch an M3 instance,
	// we ignore any instance store volumes specified in the block device mapping
	// for the AMI.
	VirtualName *string `locationName:"virtualName" type:"string"`
}

//PrivateIPAddressSpecification ...
type PrivateIPAddressSpecification struct {
	_ struct{} `type:"structure"`

	// Indicates whether the private IPv4 address is the primary private IPv4 address.
	// Only one IPv4 address can be designated as primary.
	Primary *bool `locationName:"primary" type:"boolean"`

	// The private IPv4 addresses.
	//
	// PrivateIpAddress is a required field
	PrivateIPAddress *string `locationName:"privateIpAddress" type:"string" required:"true"`
}

type ModifyInstanceKeyPairInput struct {
	_ struct{} `type:"structure"`

	// The ID of the Windows instance.
	InstanceId *string `locationName:"instanceId" type:"string"`

	KeyName *string `locationName:"keyName" type:"string"`
}

type EbsBlockDevice struct {
	_ struct{} `type:"structure"`

	// Indicates whether the EBS volume is deleted on instance termination.
	DeleteOnTermination *bool `locationName:"deleteOnTermination" type:"boolean"`

	// Indicates whether the EBS volume is encrypted. Encrypted volumes can only
	// be attached to instances that support Amazon EBS encryption. If you are creating
	// a volume from a snapshot, you can't specify an encryption value. This is
	// because only blank volumes can be encrypted on creation.
	Encrypted *bool `locationName:"encrypted" type:"boolean"`

	// The number of I/O operations per second (IOPS) that the volume supports.
	// For io1, this represents the number of IOPS that are provisioned for the
	// volume. For gp2, this represents the baseline performance of the volume and
	// the rate at which the volume accumulates I/O credits for bursting. For more
	// information about General Purpose SSD baseline performance, I/O credits,
	// and bursting, see Amazon EBS Volume Types (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html)
	// in the Amazon Elastic Compute Cloud User Guide.
	//
	// Constraint: Range is 100-20000 IOPS for io1 volumes and 100-10000 IOPS for
	// gp2 volumes.
	//
	// Condition: This parameter is required for requests to create io1 volumes;
	// it is not used in requests to create gp2, st1, sc1, or standard volumes.
	Iops *int64 `locationName:"iops" type:"integer"`

	// ID for a user-managed CMK under which the EBS volume is encrypted.
	//
	// Note: This parameter is only supported on BlockDeviceMapping objects called
	// by RunInstances (http://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_RunInstances.html),
	// RequestSpotFleet (http://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_RequestSpotFleet.html),
	// and RequestSpotInstances (http://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_RequestSpotInstances.html).
	KmsKeyId *string `type:"string"`

	// The ID of the snapshot.
	SnapshotId *string `locationName:"snapshotId" type:"string"`

	// The size of the volume, in GiB.
	//
	// Constraints: 1-16384 for General Purpose SSD (gp2), 4-16384 for Provisioned
	// IOPS SSD (io1), 500-16384 for Throughput Optimized HDD (st1), 500-16384 for
	// Cold HDD (sc1), and 1-1024 for Magnetic (standard) volumes. If you specify
	// a snapshot, the volume size must be equal to or larger than the snapshot
	// size.
	//
	// Default: If you're creating the volume from a snapshot and don't specify
	// a volume size, the default is the snapshot size.
	VolumeSize *int64 `locationName:"volumeSize" type:"integer"`

	// The volume type: gp2, io1, st1, sc1, or standard.
	//
	// Default: standard
	VolumeType *string `locationName:"volumeType" type:"string" enum:"VolumeType"`
}

type GetPasswordDataInput struct {
	_ struct{} `type:"structure"`

	// The ID of the Windows instance.
	//
	// InstanceId is a required field
	InstanceId *string `type:"string" required:"true"`
}

type GetPasswordDataOutput struct {
	_ struct{} `type:"structure"`

	// The ID of the Windows instance.
	InstanceId *string `locationName:"instanceId" type:"string"`

	// The password of the instance. Returns an empty string if the password is
	// not available.
	PasswordData *string `locationName:"passwordData" type:"string"`

	// The time the data was last updated.
	Timestamp *time.Time `locationName:"timestamp" type:"timestamp" timestampFormat:"iso8601"`
}

type TerminateInstancesInput struct {
	// _ struct{} `type:"structure"`

	InstanceIds []*string `locationName:"InstanceId" locationNameList:"InstanceId" type:"list" required:"true"`
}

type TerminateInstancesOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more terminated instances.
	TerminatingInstances []*InstanceStateChange `locationName:"instancesSet" locationNameList:"item" type:"list"`
}
