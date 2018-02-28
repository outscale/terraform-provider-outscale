package fcu

import (
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

	RequesterId *string `locationName:"requestId" type:"string"`

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
	RequesterId *string `locationName:"requestId" type:"string"`

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

	SourceDestCheck *bool `locationName:"sourceDestCheck" type:"bool"`

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

	SecurityGroupIds []*string `locationName:"SecurityGroupId" locationNameList:"SecurityGroupId" type:"list"`

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
	InstanceInitiatedShutdownBehavior *string `locationName:"instanceInitiatedShutdownBehavior" type:"string" enum:"ShutdownBehavior"`

	// The instance type. For more information, see Instance Types (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html)
	// in the Amazon Elastic Compute Cloud User Guide.
	//
	// Default: m1.small
	InstanceType *string `type:"string" enum:"InstanceType"`

	InstanceName *string `type:"string" enum:"InstanceName"`

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

	PrivateIPAddresses *string `locationName:"privateIpAddresses" type:"string"`

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

	TagSpecifications []*TagSpecification `locationName:"TagSpecification" locationNameList:"item" type:"list"`

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

	// Set to vpc to allocate the address for use with instances in a VPC.
	//
	// Default: The address is for use with instances in EC2-Classic.
	Domain *string `type:"string" enum:"DomainType"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`
}

// String returns the string representation
func (s AllocateAddressInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AllocateAddressInput) GoString() string {
	return s.String()
}

// SetDomain sets the Domain field's value.
func (s *AllocateAddressInput) SetDomain(v string) *AllocateAddressInput {
	s.Domain = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *AllocateAddressInput) SetDryRun(v bool) *AllocateAddressInput {
	s.DryRun = &v
	return s
}

// Contains the output of AllocateAddress.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/AllocateAddressResult
type AllocateAddressOutput struct {
	_ struct{} `type:"structure"`

	// [EC2-VPC] The ID that AWS assigns to represent the allocation of the Elastic
	// IP address for use with instances in a VPC.
	AllocationId *string `locationName:"allocationId" type:"string"`

	// Indicates whether this Elastic IP address is for use with instances in EC2-Classic
	// (standard) or instances in a VPC (vpc).
	Domain *string `locationName:"domain" type:"string" enum:"DomainType"`

	// The Elastic IP address.
	PublicIp *string `locationName:"publicIp" type:"string"`
}

// String returns the string representation
func (s AllocateAddressOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AllocateAddressOutput) GoString() string {
	return s.String()
}

// SetAllocationId sets the AllocationId field's value.
func (s *AllocateAddressOutput) SetAllocationId(v string) *AllocateAddressOutput {
	s.AllocationId = &v
	return s
}

// SetDomain sets the Domain field's value.
func (s *AllocateAddressOutput) SetDomain(v string) *AllocateAddressOutput {
	s.Domain = &v
	return s
}

// SetPublicIp sets the PublicIp field's value.
func (s *AllocateAddressOutput) SetPublicIp(v string) *AllocateAddressOutput {
	s.PublicIp = &v
	return s
}

type DescribeAddressesInput struct {
	_ struct{} `type:"structure"`

	// [EC2-VPC] One or more allocation IDs.
	//
	// Default: Describes all your Elastic IP addresses.
	AllocationIds []*string `locationName:"AllocationId" locationNameList:"AllocationId" type:"list"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more filters. Filter names and values are case-sensitive.
	//
	//    * allocation-id - [EC2-VPC] The allocation ID for the address.
	//
	//    * association-id - [EC2-VPC] The association ID for the address.
	//
	//    * domain - Indicates whether the address is for use in EC2-Classic (standard)
	//    or in a VPC (vpc).
	//
	//    * instance-id - The ID of the instance the address is associated with,
	//    if any.
	//
	//    * network-interface-id - [EC2-VPC] The ID of the network interface that
	//    the address is associated with, if any.
	//
	//    * network-interface-owner-id - The AWS account ID of the owner.
	//
	//    * private-ip-address - [EC2-VPC] The private IP address associated with
	//    the Elastic IP address.
	//
	//    * public-ip - The Elastic IP address.
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	//
	// Default: Describes all your Elastic IP addresses.
	PublicIps []*string `locationName:"PublicIp" locationNameList:"PublicIp" type:"list"`
}

// String returns the string representation
func (s DescribeAddressesInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeAddressesInput) GoString() string {
	return s.String()
}

// SetAllocationIds sets the AllocationIds field's value.
func (s *DescribeAddressesInput) SetAllocationIds(v []*string) *DescribeAddressesInput {
	s.AllocationIds = v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeAddressesInput) SetDryRun(v bool) *DescribeAddressesInput {
	s.DryRun = &v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribeAddressesInput) SetFilters(v []*Filter) *DescribeAddressesInput {
	s.Filters = v
	return s
}

// SetPublicIps sets the PublicIps field's value.
func (s *DescribeAddressesInput) SetPublicIps(v []*string) *DescribeAddressesInput {
	s.PublicIps = v
	return s
}

// Contains the output of DescribeAddresses
type DescribeAddressesOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more Elastic IP addresses.
	Addresses []*Address `locationName:"addressesSet" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s DescribeAddressesOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeAddressesOutput) GoString() string {
	return s.String()
}

// SetAddresses sets the Addresses field's value.
func (s *DescribeAddressesOutput) SetAddresses(v []*Address) *DescribeAddressesOutput {
	s.Addresses = v
	return s
}

// Describes an Elastic IP address.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/Address
type Address struct {
	_ struct{} `type:"structure"`

	// The ID representing the allocation of the address for use with EC2-VPC.
	AllocationId *string `locationName:"allocationId" type:"string"`

	// The ID representing the association of the address with an instance in a
	// VPC.
	AssociationId *string `locationName:"associationId" type:"string"`

	AllowReassociation *bool `locationName:"allowReassociation" type:"bool"`

	// Indicates whether this Elastic IP address is for use with instances in EC2-Classic
	// (standard) or instances in a VPC (vpc).
	Domain *string `locationName:"domain" type:"string" enum:"DomainType"`

	// The ID of the instance that the address is associated with (if any).
	InstanceId *string `locationName:"instanceId" type:"string"`

	// The ID of the network interface.
	NetworkInterfaceId *string `locationName:"networkInterfaceId" type:"string"`

	// The ID of the AWS account that owns the network interface.
	NetworkInterfaceOwnerId *string `locationName:"networkInterfaceOwnerId" type:"string"`

	// The private IP address associated with the Elastic IP address.
	PrivateIpAddress *string `locationName:"privateIpAddress" type:"string"`

	// The Elastic IP address.
	PublicIp *string `locationName:"publicIp" type:"string"`
}

// String returns the string representation
func (s Address) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Address) GoString() string {
	return s.String()
}

// SetAllocationId sets the AllocationId field's value.
func (s *Address) SetAllocationId(v string) *Address {
	s.AllocationId = &v
	return s
}

// SetAllowReassociation sets the AllowReassociation field's value.
func (s *Address) SetAllowReassociation(v bool) *Address {
	s.AllowReassociation = &v
	return s
}

// SetAssociationId sets the AssociationId field's value.
func (s *Address) SetAssociationId(v string) *Address {
	s.AssociationId = &v
	return s
}

// SetDomain sets the Domain field's value.
func (s *Address) SetDomain(v string) *Address {
	s.Domain = &v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *Address) SetInstanceId(v string) *Address {
	s.InstanceId = &v
	return s
}

// SetNetworkInterfaceId sets the NetworkInterfaceId field's value.
func (s *Address) SetNetworkInterfaceId(v string) *Address {
	s.NetworkInterfaceId = &v
	return s
}

// SetNetworkInterfaceOwnerId sets the NetworkInterfaceOwnerId field's value.
func (s *Address) SetNetworkInterfaceOwnerId(v string) *Address {
	s.NetworkInterfaceOwnerId = &v
	return s
}

// SetPrivateIpAddress sets the PrivateIpAddress field's value.
func (s *Address) SetPrivateIpAddress(v string) *Address {
	s.PrivateIpAddress = &v
	return s
}

// SetPublicIp sets the PublicIp field's value.
func (s *Address) SetPublicIp(v string) *Address {
	s.PublicIp = &v
	return s
}

// Contains the parameters for ModifyInstanceAttribute.
type ModifyInstanceAttributeInput struct {
	_ struct{} `type:"structure"`
	// The name of the attribute.
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

// String returns the string representation
func (s ModifyInstanceAttributeInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ModifyInstanceAttributeInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
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

// SetAttribute sets the Attribute field's value.
func (s *ModifyInstanceAttributeInput) SetAttribute(v string) *ModifyInstanceAttributeInput {
	s.Attribute = &v
	return s
}

// SetBlockDeviceMappings sets the BlockDeviceMappings field's value.
func (s *ModifyInstanceAttributeInput) SetBlockDeviceMappings(v []*InstanceBlockDeviceMappingSpecification) *ModifyInstanceAttributeInput {
	s.BlockDeviceMappings = v
	return s
}

// SetDisableApiTermination sets the DisableApiTermination field's value.
func (s *ModifyInstanceAttributeInput) SetDisableApiTermination(v *AttributeBooleanValue) *ModifyInstanceAttributeInput {
	s.DisableApiTermination = v
	return s
}

// SetEbsOptimized sets the EbsOptimized field's value.
func (s *ModifyInstanceAttributeInput) SetEbsOptimized(v *AttributeBooleanValue) *ModifyInstanceAttributeInput {
	s.EbsOptimized = v
	return s
}

// SetGroups sets the Groups field's value.
func (s *ModifyInstanceAttributeInput) SetGroups(v []*string) *ModifyInstanceAttributeInput {
	s.Groups = v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *ModifyInstanceAttributeInput) SetInstanceId(v string) *ModifyInstanceAttributeInput {
	s.InstanceId = &v
	return s
}

// SetInstanceInitiatedShutdownBehavior sets the InstanceInitiatedShutdownBehavior field's value.
func (s *ModifyInstanceAttributeInput) SetInstanceInitiatedShutdownBehavior(v *AttributeValue) *ModifyInstanceAttributeInput {
	s.InstanceInitiatedShutdownBehavior = v
	return s
}

// SetInstanceType sets the InstanceType field's value.
func (s *ModifyInstanceAttributeInput) SetInstanceType(v *AttributeValue) *ModifyInstanceAttributeInput {
	s.InstanceType = v
	return s
}

// SetSourceDestCheck sets the SourceDestCheck field's value.
func (s *ModifyInstanceAttributeInput) SetSourceDestCheck(v *AttributeBooleanValue) *ModifyInstanceAttributeInput {
	s.SourceDestCheck = v
	return s
}

// SetUserData sets the UserData field's value.
func (s *ModifyInstanceAttributeInput) SetUserData(v *BlobAttributeValue) *ModifyInstanceAttributeInput {
	s.UserData = v
	return s
}

// SetValue sets the Value field's value.
func (s *ModifyInstanceAttributeInput) SetValue(v string) *ModifyInstanceAttributeInput {
	s.Value = &v
	return s
}

type BlobAttributeValue struct {
	_ struct{} `type:"structure"`

	// Value is automatically base64 encoded/decoded by the SDK.
	Value []byte `locationName:"value" type:"blob"`
}

// String returns the string representation
func (s BlobAttributeValue) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s BlobAttributeValue) GoString() string {
	return s.String()
}

// SetValue sets the Value field's value.
func (s *BlobAttributeValue) SetValue(v []byte) *BlobAttributeValue {
	s.Value = v
	return s
}

type StopInstancesInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// Forces the instances to stop. The instances do not have an opportunity to
	// flush file system caches or file system metadata. If you use this option,
	// you must perform file system check and repair procedures. This option is
	// not recommended for Windows instances.
	//
	// Default: false
	Force *bool `locationName:"force" type:"boolean"`

	// One or more instance IDs.
	//
	// InstanceIds is a required field
	InstanceIds []*string `locationName:"InstanceId" locationNameList:"InstanceId" type:"list" required:"true"`
}

// String returns the string representation
func (s StopInstancesInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s StopInstancesInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
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

// SetDryRun sets the DryRun field's value.
func (s *StopInstancesInput) SetDryRun(v bool) *StopInstancesInput {
	s.DryRun = &v
	return s
}

// SetForce sets the Force field's value.
func (s *StopInstancesInput) SetForce(v bool) *StopInstancesInput {
	s.Force = &v
	return s
}

// SetInstanceIds sets the InstanceIds field's value.
func (s *StopInstancesInput) SetInstanceIds(v []*string) *StopInstancesInput {
	s.InstanceIds = v
	return s
}

// Contains the output of StopInstances.
type StopInstancesOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more stopped instances.
	StoppingInstances []*InstanceStateChange `locationName:"instancesSet" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s StopInstancesOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s StopInstancesOutput) GoString() string {
	return s.String()
}

// SetStoppingInstances sets the StoppingInstances field's value.
func (s *StopInstancesOutput) SetStoppingInstances(v []*InstanceStateChange) *StopInstancesOutput {
	s.StoppingInstances = v
	return s
}

type ModifyInstanceAttributeOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s ModifyInstanceAttributeOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ModifyInstanceAttributeOutput) GoString() string {
	return s.String()
}

type StartInstancesInput struct {
	_ struct{} `type:"structure"`

	// Reserved.
	AdditionalInfo *string `locationName:"additionalInfo" type:"string"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more instance IDs.
	//
	// InstanceIds is a required field
	InstanceIds []*string `locationName:"InstanceId" locationNameList:"InstanceId" type:"list" required:"true"`
}

// String returns the string representation
func (s StartInstancesInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s StartInstancesInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
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

// SetAdditionalInfo sets the AdditionalInfo field's value.
func (s *StartInstancesInput) SetAdditionalInfo(v string) *StartInstancesInput {
	s.AdditionalInfo = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *StartInstancesInput) SetDryRun(v bool) *StartInstancesInput {
	s.DryRun = &v
	return s
}

// SetInstanceIds sets the InstanceIds field's value.
func (s *StartInstancesInput) SetInstanceIds(v []*string) *StartInstancesInput {
	s.InstanceIds = v
	return s
}

type StartInstancesOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more started instances.
	StartingInstances []*InstanceStateChange `locationName:"instancesSet" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s StartInstancesOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s StartInstancesOutput) GoString() string {
	return s.String()
}

// SetStartingInstances sets the StartingInstances field's value.
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

// String returns the string representation
func (s AssociateAddressInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AssociateAddressInput) GoString() string {
	return s.String()
}

// SetAllocationId sets the AllocationId field's value.
func (s *AssociateAddressInput) SetAllocationId(v string) *AssociateAddressInput {
	s.AllocationId = &v
	return s
}

// SetAllowReassociation sets the AllowReassociation field's value.
func (s *AssociateAddressInput) SetAllowReassociation(v bool) *AssociateAddressInput {
	s.AllowReassociation = &v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *AssociateAddressInput) SetInstanceId(v string) *AssociateAddressInput {
	s.InstanceId = &v
	return s
}

// SetNetworkInterfaceId sets the NetworkInterfaceId field's value.
func (s *AssociateAddressInput) SetNetworkInterfaceId(v string) *AssociateAddressInput {
	s.NetworkInterfaceId = &v
	return s
}

// SetPrivateIpAddress sets the PrivateIpAddress field's value.
func (s *AssociateAddressInput) SetPrivateIpAddress(v string) *AssociateAddressInput {
	s.PrivateIpAddress = &v
	return s
}

// SetPublicIp sets the PublicIp field's value.
func (s *AssociateAddressInput) SetPublicIp(v string) *AssociateAddressInput {
	s.PublicIp = &v
	return s
}

// Contains the output of AssociateAddress.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/AssociateAddressResult
type AssociateAddressOutput struct {
	_ struct{} `type:"structure"`

	// [EC2-VPC] The ID that represents the association of the Elastic IP address
	// with an instance.
	AssociationId *string `locationName:"associationId" type:"string"`

	RequestId *string `locationName:"requestId" type:"string"`
}

// String returns the string representation
func (s AssociateAddressOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AssociateAddressOutput) GoString() string {
	return s.String()
}

// SetAssociationId sets the AssociationId field's value.
func (s *AssociateAddressOutput) SetAssociationId(v string) *AssociateAddressOutput {
	s.AssociationId = &v
	return s
}

// SetRequestId sets the AssociationId field's value.
func (s *AssociateAddressOutput) SetRequestId(v string) *AssociateAddressOutput {
	s.RequestId = &v
	return s
}

type DisassociateAddressInput struct {
	_ struct{} `type:"structure"`

	AssociationId *string `type:"string"`

	PublicIp *string `type:"string"`
}

// String returns the string representation
func (s DisassociateAddressInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DisassociateAddressInput) GoString() string {
	return s.String()
}

// SetAssociationId sets the AssociationId field's value.
func (s *DisassociateAddressInput) SetAssociationId(v string) *DisassociateAddressInput {
	s.AssociationId = &v
	return s
}

// SetPublicIp sets the PublicIp field's value.
func (s *DisassociateAddressInput) SetPublicIp(v string) *DisassociateAddressInput {
	s.PublicIp = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DisassociateAddressOutput
type DisassociateAddressOutput struct {
	_ struct{} `type:"structure"`

	RequestId *string `locationName:"requestId" type:"string"`
	Return    *bool   `locationName:"return" type:"boolean"`
}

// String returns the string representation
func (s DisassociateAddressOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DisassociateAddressOutput) GoString() string {
	return s.String()
}

// SetAssociationId sets the AssociationId field's value.
func (s *DisassociateAddressOutput) SetReturn(v bool) *DisassociateAddressOutput {
	s.Return = &v
	return s
}

// SetRequestId sets the AssociationId field's value.
func (s *DisassociateAddressOutput) SetRequestId(v string) *DisassociateAddressOutput {
	s.RequestId = &v
	return s
}

type ReleaseAddressInput struct {
	_ struct{} `type:"structure"`

	// [EC2-VPC] The allocation ID. Required for EC2-VPC.
	AllocationId *string `type:"string"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// [EC2-Classic] The Elastic IP address. Required for EC2-Classic.
	PublicIp *string `type:"string"`
}

// String returns the string representation
func (s ReleaseAddressInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ReleaseAddressInput) GoString() string {
	return s.String()
}

// SetAllocationId sets the AllocationId field's value.
func (s *ReleaseAddressInput) SetAllocationId(v string) *ReleaseAddressInput {
	s.AllocationId = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *ReleaseAddressInput) SetDryRun(v bool) *ReleaseAddressInput {
	s.DryRun = &v
	return s
}

// SetPublicIp sets the PublicIp field's value.
func (s *ReleaseAddressInput) SetPublicIp(v string) *ReleaseAddressInput {
	s.PublicIp = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/ReleaseAddressOutput
type ReleaseAddressOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s ReleaseAddressOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ReleaseAddressOutput) GoString() string {
	return s.String()
}

type RegisterImageInput struct {
	_ struct{} `type:"structure"`

	// The architecture of the AMI.
	//
	// Default: For Amazon EBS-backed AMIs, i386. For instance store-backed AMIs,
	// the architecture specified in the manifest file.
	Architecture *string `locationName:"architecture" type:"string" enum:"ArchitectureValues"`

	// The billing product codes. Your account must be authorized to specify billing
	// product codes. Otherwise, you can use the AWS Marketplace to bill for the
	// use of an AMI.
	BillingProducts []*string `locationName:"BillingProduct" locationNameList:"item" type:"list"`

	// One or more block device mapping entries.
	BlockDeviceMappings []*BlockDeviceMapping `locationName:"BlockDeviceMapping" locationNameList:"BlockDeviceMapping" type:"list"`

	// A description for your AMI.
	Description *string `locationName:"description" type:"string"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// Set to true to enable enhanced networking with ENA for the AMI and any instances
	// that you launch from the AMI.
	//
	// This option is supported only for HVM AMIs. Specifying this option with a
	// PV AMI can make instances launched from the AMI unreachable.
	EnaSupport *bool `locationName:"enaSupport" type:"boolean"`

	// The full path to your AMI manifest in Amazon S3 storage.
	ImageLocation *string `type:"string"`

	InstanceId *string `type:"string"`

	NoReboot *bool `type:"boolean"`

	// The ID of the kernel.
	KernelId *string `locationName:"kernelId" type:"string"`

	// A name for your AMI.
	//
	// Constraints: 3-128 alphanumeric characters, parentheses (()), square brackets
	// ([]), spaces ( ), periods (.), slashes (/), dashes (-), single quotes ('),
	// at-signs (@), or underscores(_)
	//
	// Name is a required field
	Name *string `locationName:"name" type:"string" required:"true"`

	// The ID of the RAM disk.
	RamdiskId *string `locationName:"ramdiskId" type:"string"`

	// The name of the root device (for example, /dev/sda1, or /dev/xvda).
	RootDeviceName *string `locationName:"rootDeviceName" type:"string"`

	// Set to simple to enable enhanced networking with the Intel 82599 Virtual
	// Function interface for the AMI and any instances that you launch from the
	// AMI.
	//
	// There is no way to disable sriovNetSupport at this time.
	//
	// This option is supported only for HVM AMIs. Specifying this option with a
	// PV AMI can make instances launched from the AMI unreachable.
	SriovNetSupport *string `locationName:"sriovNetSupport" type:"string"`

	// The type of virtualization.
	//
	// Default: paravirtual
	VirtualizationType *string `locationName:"virtualizationType" type:"string"`
}

// String returns the string representation
func (s RegisterImageInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s RegisterImageInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
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

// SetArchitecture sets the Architecture field's value.
func (s *RegisterImageInput) SetArchitecture(v string) *RegisterImageInput {
	s.Architecture = &v
	return s
}

// SetBillingProducts sets the BillingProducts field's value.
func (s *RegisterImageInput) SetBillingProducts(v []*string) *RegisterImageInput {
	s.BillingProducts = v
	return s
}

// SetBlockDeviceMappings sets the BlockDeviceMappings field's value.
func (s *RegisterImageInput) SetBlockDeviceMappings(v []*BlockDeviceMapping) *RegisterImageInput {
	s.BlockDeviceMappings = v
	return s
}

// SetDescription sets the Description field's value.
func (s *RegisterImageInput) SetDescription(v string) *RegisterImageInput {
	s.Description = &v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *RegisterImageInput) SetDryRun(v bool) *RegisterImageInput {
	s.DryRun = &v
	return s
}

// SetEnaSupport sets the EnaSupport field's value.
func (s *RegisterImageInput) SetEnaSupport(v bool) *RegisterImageInput {
	s.EnaSupport = &v
	return s
}

// SetImageLocation sets the ImageLocation field's value.
func (s *RegisterImageInput) SetImageLocation(v string) *RegisterImageInput {
	s.ImageLocation = &v
	return s
}

// SetKernelId sets the KernelId field's value.
func (s *RegisterImageInput) SetKernelId(v string) *RegisterImageInput {
	s.KernelId = &v
	return s
}

// SetName sets the Name field's value.
func (s *RegisterImageInput) SetName(v string) *RegisterImageInput {
	s.Name = &v
	return s
}

// SetRamdiskId sets the RamdiskId field's value.
func (s *RegisterImageInput) SetRamdiskId(v string) *RegisterImageInput {
	s.RamdiskId = &v
	return s
}

// SetRootDeviceName sets the RootDeviceName field's value.
func (s *RegisterImageInput) SetRootDeviceName(v string) *RegisterImageInput {
	s.RootDeviceName = &v
	return s
}

// SetSriovNetSupport sets the SriovNetSupport field's value.
func (s *RegisterImageInput) SetSriovNetSupport(v string) *RegisterImageInput {
	s.SriovNetSupport = &v
	return s
}

// SetVirtualizationType sets the VirtualizationType field's value.
func (s *RegisterImageInput) SetVirtualizationType(v string) *RegisterImageInput {
	s.VirtualizationType = &v
	return s
}

// Contains the output of RegisterImage.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/RegisterImageResult
type RegisterImageOutput struct {
	_ struct{} `type:"structure"`

	// The ID of the newly registered AMI.
	ImageId *string `locationName:"imageId" type:"string"`
}

// String returns the string representation
func (s RegisterImageOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s RegisterImageOutput) GoString() string {
	return s.String()
}

// SetImageId sets the ImageId field's value.
func (s *RegisterImageOutput) SetImageId(v string) *RegisterImageOutput {
	s.ImageId = &v
	return s
}

type DeregisterImageInput struct {
	_ struct{} `type:"structure"`

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
func (s DeregisterImageInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeregisterImageInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
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

// SetDryRun sets the DryRun field's value.
func (s *DeregisterImageInput) SetDryRun(v bool) *DeregisterImageInput {
	s.DryRun = &v
	return s
}

// SetImageId sets the ImageId field's value.
func (s *DeregisterImageInput) SetImageId(v string) *DeregisterImageInput {
	s.ImageId = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeregisterImageOutput
type DeregisterImageOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeregisterImageOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeregisterImageOutput) GoString() string {
	return s.String()
}

type Image struct {
	_ struct{} `type:"structure"`

	// The architecture of the image.
	Architecture *string `locationName:"architecture" type:"string" enum:"ArchitectureValues"`

	// Any block device mapping entries.
	BlockDeviceMappings []*BlockDeviceMapping `locationName:"blockDeviceMapping" locationNameList:"item" type:"list"`

	// The date and time the image was created.
	CreationDate *string `locationName:"creationDate" type:"string"`

	// The description of the AMI that was provided during image creation.
	Description *string `locationName:"description" type:"string"`

	// Specifies whether enhanced networking with ENA is enabled.
	EnaSupport *bool `locationName:"enaSupport" type:"boolean"`

	// The hypervisor type of the image.
	Hypervisor *string `locationName:"hypervisor" type:"string" enum:"HypervisorType"`

	// The ID of the AMI.
	ImageId *string `locationName:"imageId" type:"string"`

	// The location of the AMI.
	ImageLocation *string `locationName:"imageLocation" type:"string"`

	// The AWS account alias (for example, amazon, self) or the AWS account ID of
	// the AMI owner.
	ImageOwnerAlias *string `locationName:"imageOwnerAlias" type:"string"`

	// The type of image.
	ImageType *string `locationName:"imageType" type:"string" enum:"ImageTypeValues"`

	// The kernel associated with the image, if any. Only applicable for machine
	// images.
	KernelId *string `locationName:"kernelId" type:"string"`

	// The name of the AMI that was provided during image creation.
	Name *string `locationName:"name" type:"string"`

	// The AWS account ID of the image owner.
	OwnerId *string `locationName:"imageOwnerId" type:"string"`

	// The value is Windows for Windows AMIs; otherwise blank.
	Platform *string `locationName:"platform" type:"string" enum:"PlatformValues"`

	// Any product codes associated with the AMI.
	ProductCodes []*ProductCode `locationName:"productCodes" locationNameList:"item" type:"list"`

	// Indicates whether the image has public launch permissions. The value is true
	// if this image has public launch permissions or false if it has only implicit
	// and explicit launch permissions.
	Public *bool `locationName:"isPublic" type:"boolean"`

	// The RAM disk associated with the image, if any. Only applicable for machine
	// images.
	RamdiskId *string `locationName:"ramdiskId" type:"string"`

	// The device name of the root device (for example, /dev/sda1 or /dev/xvda).
	RootDeviceName *string `locationName:"rootDeviceName" type:"string"`

	// The type of root device used by the AMI. The AMI can use an EBS volume or
	// an instance store volume.
	RootDeviceType *string `locationName:"rootDeviceType" type:"string" enum:"DeviceType"`

	// Specifies whether enhanced networking with the Intel 82599 Virtual Function
	// interface is enabled.
	SriovNetSupport *string `locationName:"sriovNetSupport" type:"string"`

	// The current state of the AMI. If the state is available, the image is successfully
	// registered and can be used to launch an instance.
	State *string `locationName:"imageState" type:"string" enum:"ImageState"`

	// The reason for the state change.
	StateReason *StateReason `locationName:"stateReason" type:"structure"`

	// Any tags assigned to the image.
	Tags []*Tag `locationName:"tagSet" locationNameList:"item" type:"list"`

	// The type of virtualization of the AMI.
	VirtualizationType *string `locationName:"virtualizationType" type:"string" enum:"VirtualizationType"`
}

// String returns the string representation
func (s Image) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Image) GoString() string {
	return s.String()
}

// SetArchitecture sets the Architecture field's value.
func (s *Image) SetArchitecture(v string) *Image {
	s.Architecture = &v
	return s
}

type DescribeImagesInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// Scopes the images by users with explicit launch permissions. Specify an AWS
	// account ID, self (the sender of the request), or all (public AMIs).
	ExecutableUsers []*string `locationName:"ExecutableBy" locationNameList:"ExecutableBy" type:"list"`

	// One or more filters.
	//
	//    * architecture - The image architecture (i386 | x86_64).
	//
	//    * block-device-mapping.delete-on-termination - A Boolean value that indicates
	//    whether the Amazon EBS volume is deleted on instance termination.
	//
	//    * block-device-mapping.device-name - The device name for the EBS volume
	//    (for example, /dev/sdh).
	//
	//    * block-device-mapping.snapshot-id - The ID of the snapshot used for the
	//    EBS volume.
	//
	//    * block-device-mapping.volume-size - The volume size of the EBS volume,
	//    in GiB.
	//
	//    * block-device-mapping.volume-type - The volume type of the EBS volume
	//    (gp2 | io1 | st1 | sc1 | standard).
	//
	//    * description - The description of the image (provided during image creation).
	//
	//    * ena-support - A Boolean that indicates whether enhanced networking with
	//    ENA is enabled.
	//
	//    * hypervisor - The hypervisor type (ovm | xen).
	//
	//    * image-id - The ID of the image.
	//
	//    * image-type - The image type (machine | kernel | ramdisk).
	//
	//    * is-public - A Boolean that indicates whether the image is public.
	//
	//    * kernel-id - The kernel ID.
	//
	//    * manifest-location - The location of the image manifest.
	//
	//    * name - The name of the AMI (provided during image creation).
	//
	//    * owner-alias - String value from an Amazon-maintained list (amazon |
	//    aws-marketplace | microsoft) of snapshot owners. Not to be confused with
	//    the user-configured AWS account alias, which is set from the IAM console.
	//
	//    * owner-id - The AWS account ID of the image owner.
	//
	//    * platform - The platform. To only list Windows-based AMIs, use windows.
	//
	//    * product-code - The product code.
	//
	//    * product-code.type - The type of the product code (devpay | marketplace).
	//
	//    * ramdisk-id - The RAM disk ID.
	//
	//    * root-device-name - The name of the root device volume (for example,
	//    /dev/sda1).
	//
	//    * root-device-type - The type of the root device volume (ebs | instance-store).
	//
	//    * state - The state of the image (available | pending | failed).
	//
	//    * state-reason-code - The reason code for the state change.
	//
	//    * state-reason-message - The message for the state change.
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
	//    * virtualization-type - The virtualization type (paravirtual | hvm).
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// One or more image IDs.
	//
	// Default: Describes all images available to you.
	ImageIds []*string `locationName:"ImageId" locationNameList:"ImageId" type:"list"`

	// Filters the images by the owner. Specify an AWS account ID, self (owner is
	// the sender of the request), or an AWS owner alias (valid values are amazon
	// | aws-marketplace | microsoft). Omitting this option returns all images for
	// which you have launch permissions, regardless of ownership.
	Owners []*string `locationName:"Owner" locationNameList:"Owner" type:"list"`
}

// String returns the string representation
func (s DescribeImagesInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeImagesInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeImagesInput) SetDryRun(v bool) *DescribeImagesInput {
	s.DryRun = &v
	return s
}

// SetExecutableUsers sets the ExecutableUsers field's value.
func (s *DescribeImagesInput) SetExecutableUsers(v []*string) *DescribeImagesInput {
	s.ExecutableUsers = v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribeImagesInput) SetFilters(v []*Filter) *DescribeImagesInput {
	s.Filters = v
	return s
}

// SetImageIds sets the ImageIds field's value.
func (s *DescribeImagesInput) SetImageIds(v []*string) *DescribeImagesInput {
	s.ImageIds = v
	return s
}

// SetOwners sets the Owners field's value.
func (s *DescribeImagesInput) SetOwners(v []*string) *DescribeImagesInput {
	s.Owners = v
	return s
}

// Contains the output of DescribeImages.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeImagesResult
type DescribeImagesOutput struct {
	_ struct{} `type:"structure"`

	// Information about one or more images.
	Images []*Image `locationName:"imagesSet" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s DescribeImagesOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeImagesOutput) GoString() string {
	return s.String()
}

// SetImages sets the Images field's value.
func (s *DescribeImagesOutput) SetImages(v []*Image) *DescribeImagesOutput {
	s.Images = v
	return s
}

type ModifyImageAttributeInput struct {
	_ struct{} `type:"structure"`

	// The name of the attribute to modify.
	Attribute *string `type:"string"`

	// A description for the AMI.
	Description *AttributeValue `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the AMI.
	//
	// ImageId is a required field
	ImageId *string `type:"string" required:"true"`

	// A launch permission modification.
	LaunchPermission *LaunchPermissionModifications `type:"structure"`

	// The operation type.
	OperationType *string `type:"string" enum:"OperationType"`

	// One or more product codes. After you add a product code to an AMI, it can't
	// be removed. This is only valid when modifying the productCodes attribute.
	ProductCodes []*string `locationName:"ProductCode" locationNameList:"ProductCode" type:"list"`

	// One or more user groups. This is only valid when modifying the launchPermission
	// attribute.
	UserGroups []*string `locationName:"UserGroup" locationNameList:"UserGroup" type:"list"`

	// One or more AWS account IDs. This is only valid when modifying the launchPermission
	// attribute.
	UserIds []*string `locationName:"UserId" locationNameList:"UserId" type:"list"`

	// The value of the attribute being modified. This is only valid when modifying
	// the description attribute.
	Value *string `type:"string"`
}

// String returns the string representation
func (s ModifyImageAttributeInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ModifyImageAttributeInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
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

// SetAttribute sets the Attribute field's value.
func (s *ModifyImageAttributeInput) SetAttribute(v string) *ModifyImageAttributeInput {
	s.Attribute = &v
	return s
}

// SetDescription sets the Description field's value.
func (s *ModifyImageAttributeInput) SetDescription(v *AttributeValue) *ModifyImageAttributeInput {
	s.Description = v
	return s
}

// SetDryRun sets the DryRun field's value.
func (s *ModifyImageAttributeInput) SetDryRun(v bool) *ModifyImageAttributeInput {
	s.DryRun = &v
	return s
}

// SetImageId sets the ImageId field's value.
func (s *ModifyImageAttributeInput) SetImageId(v string) *ModifyImageAttributeInput {
	s.ImageId = &v
	return s
}

// SetLaunchPermission sets the LaunchPermission field's value.
func (s *ModifyImageAttributeInput) SetLaunchPermission(v *LaunchPermissionModifications) *ModifyImageAttributeInput {
	s.LaunchPermission = v
	return s
}

// SetOperationType sets the OperationType field's value.
func (s *ModifyImageAttributeInput) SetOperationType(v string) *ModifyImageAttributeInput {
	s.OperationType = &v
	return s
}

// SetProductCodes sets the ProductCodes field's value.
func (s *ModifyImageAttributeInput) SetProductCodes(v []*string) *ModifyImageAttributeInput {
	s.ProductCodes = v
	return s
}

// SetUserGroups sets the UserGroups field's value.
func (s *ModifyImageAttributeInput) SetUserGroups(v []*string) *ModifyImageAttributeInput {
	s.UserGroups = v
	return s
}

// SetUserIds sets the UserIds field's value.
func (s *ModifyImageAttributeInput) SetUserIds(v []*string) *ModifyImageAttributeInput {
	s.UserIds = v
	return s
}

// SetValue sets the Value field's value.
func (s *ModifyImageAttributeInput) SetValue(v string) *ModifyImageAttributeInput {
	s.Value = &v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/ModifyImageAttributeOutput
type ModifyImageAttributeOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s ModifyImageAttributeOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ModifyImageAttributeOutput) GoString() string {
	return s.String()
}

type LaunchPermissionModifications struct {
	_ struct{} `type:"structure"`

	// The AWS account ID to add to the list of launch permissions for the AMI.
	Add []*LaunchPermission `locationNameList:"item" type:"list"`

	// The AWS account ID to remove from the list of launch permissions for the
	// AMI.
	Remove []*LaunchPermission `locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s LaunchPermissionModifications) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s LaunchPermissionModifications) GoString() string {
	return s.String()
}

// SetAdd sets the Add field's value.
func (s *LaunchPermissionModifications) SetAdd(v []*LaunchPermission) *LaunchPermissionModifications {
	s.Add = v
	return s
}

// SetRemove sets the Remove field's value.
func (s *LaunchPermissionModifications) SetRemove(v []*LaunchPermission) *LaunchPermissionModifications {
	s.Remove = v
	return s
}

type LaunchPermission struct {
	_ struct{} `type:"structure"`

	// The name of the group.
	Group *string `locationName:"group" type:"string" enum:"PermissionGroup"`

	// The AWS account ID.
	UserId *string `locationName:"userId" type:"string"`
}

// String returns the string representation
func (s LaunchPermission) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s LaunchPermission) GoString() string {
	return s.String()
}

// SetGroup sets the Group field's value.
func (s *LaunchPermission) SetGroup(v string) *LaunchPermission {
	s.Group = &v
	return s
}

// SetUserId sets the UserId field's value.
func (s *LaunchPermission) SetUserId(v string) *LaunchPermission {
	s.UserId = &v
	return s
}

type DeleteTagsInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The ID of the resource. For example, ami-1a2b3c4d. You can specify more than
	// one resource ID.
	//
	// Resources is a required field
	Resources []*string `locationName:"resourceId" type:"list" required:"true"`

	// One or more tags to delete. If you omit the value parameter, we delete the
	// tag regardless of its value. If you specify this parameter with an empty
	// string as the value, we delete the key only if its value is an empty string.
	Tags []*Tag `locationName:"tag" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s DeleteTagsInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteTagsInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
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

// SetDryRun sets the DryRun field's value.
func (s *DeleteTagsInput) SetDryRun(v bool) *DeleteTagsInput {
	s.DryRun = &v
	return s
}

// SetResources sets the Resources field's value.
func (s *DeleteTagsInput) SetResources(v []*string) *DeleteTagsInput {
	s.Resources = v
	return s
}

// SetTags sets the Tags field's value.
func (s *DeleteTagsInput) SetTags(v []*Tag) *DeleteTagsInput {
	s.Tags = v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DeleteTagsOutput
type DeleteTagsOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeleteTagsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteTagsOutput) GoString() string {
	return s.String()
}

type CreateTagsInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The IDs of one or more resources to tag. For example, ami-1a2b3c4d.
	//
	// Resources is a required field
	Resources []*string `locationName:"ResourceId" type:"list" required:"true"`

	// One or more tags. The value parameter is required, but if you don't want
	// the tag to have a value, specify the parameter with no value, and we set
	// the value to an empty string.
	//
	// Tags is a required field
	Tags []*Tag `locationName:"Tag" locationNameList:"item" type:"list" required:"true"`
}

// String returns the string representation
func (s CreateTagsInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateTagsInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
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

// SetDryRun sets the DryRun field's value.
func (s *CreateTagsInput) SetDryRun(v bool) *CreateTagsInput {
	s.DryRun = &v
	return s
}

// SetResources sets the Resources field's value.
func (s *CreateTagsInput) SetResources(v []*string) *CreateTagsInput {
	s.Resources = v
	return s
}

// SetTags sets the Tags field's value.
func (s *CreateTagsInput) SetTags(v []*Tag) *CreateTagsInput {
	s.Tags = v
	return s
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateTagsOutput
type CreateTagsOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s CreateTagsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateTagsOutput) GoString() string {
	return s.String()
}

type DescribeTagsInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more filters.
	//
	//    * key - The tag key.
	//
	//    * resource-id - The resource ID.
	//
	//    * resource-type - The resource type (customer-gateway | dhcp-options |
	//    image | instance | internet-gateway | network-acl | network-interface
	//    | reserved-instances | route-table | security-group | snapshot | spot-instances-request
	//    | subnet | volume | vpc | vpn-connection | vpn-gateway).
	//
	//    * value - The tag value.
	Filters []*Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// The maximum number of results to return in a single call. This value can
	// be between 5 and 1000. To retrieve the remaining results, make another call
	// with the returned NextToken value.
	MaxResults *int64 `locationName:"maxResults" type:"integer"`

	// The token to retrieve the next page of results.
	NextToken *string `locationName:"nextToken" type:"string"`
}

// String returns the string representation
func (s DescribeTagsInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeTagsInput) GoString() string {
	return s.String()
}

// SetDryRun sets the DryRun field's value.
func (s *DescribeTagsInput) SetDryRun(v bool) *DescribeTagsInput {
	s.DryRun = &v
	return s
}

// SetFilters sets the Filters field's value.
func (s *DescribeTagsInput) SetFilters(v []*Filter) *DescribeTagsInput {
	s.Filters = v
	return s
}

// SetMaxResults sets the MaxResults field's value.
func (s *DescribeTagsInput) SetMaxResults(v int64) *DescribeTagsInput {
	s.MaxResults = &v
	return s
}

// SetNextToken sets the NextToken field's value.
func (s *DescribeTagsInput) SetNextToken(v string) *DescribeTagsInput {
	s.NextToken = &v
	return s
}

// Contains the output of DescribeTags.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeTagsResult
type DescribeTagsOutput struct {
	_ struct{} `type:"structure"`

	// The token to use to retrieve the next page of results. This value is null
	// when there are no more results to return..
	NextToken *string `locationName:"nextToken" type:"string"`

	// A list of tags.
	Tags []*TagDescription `locationName:"tagSet" locationNameList:"item" type:"list"`
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

// SetTags sets the Tags field's value.
func (s *DescribeTagsOutput) SetTags(v []*TagDescription) *DescribeTagsOutput {
	s.Tags = v
	return s
}

type TagDescription struct {
	_ struct{} `type:"structure"`

	// The tag key.
	Key *string `locationName:"key" type:"string"`

	// The ID of the resource. For example, ami-1a2b3c4d.
	ResourceId *string `locationName:"resourceId" type:"string"`

	// The resource type.
	ResourceType *string `locationName:"resourceType" type:"string" enum:"ResourceType"`

	// The tag value.
	Value *string `locationName:"value" type:"string"`
}

// String returns the string representation
func (s TagDescription) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s TagDescription) GoString() string {
	return s.String()
}

// SetKey sets the Key field's value.
func (s *TagDescription) SetKey(v string) *TagDescription {
	s.Key = &v
	return s
}

// SetResourceId sets the ResourceId field's value.
func (s *TagDescription) SetResourceId(v string) *TagDescription {
	s.ResourceId = &v
	return s
}

// SetResourceType sets the ResourceType field's value.
func (s *TagDescription) SetResourceType(v string) *TagDescription {
	s.ResourceType = &v
	return s
}

// SetValue sets the Value field's value.
func (s *TagDescription) SetValue(v string) *TagDescription {
	s.Value = &v
	return s
}

// The tags to apply to a resource when the resource is being created.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/TagSpecification
type TagSpecification struct {
	_ struct{} `type:"structure"`

	// The type of resource to tag. Currently, the resource types that support tagging
	// on creation are instance and volume.
	ResourceType *string `locationName:"resourceType" type:"string" enum:"ResourceType"`

	// The tags to apply to the resource.
	Tags []*Tag `locationName:"Tag" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s TagSpecification) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s TagSpecification) GoString() string {
	return s.String()
}

// SetResourceType sets the ResourceType field's value.
func (s *TagSpecification) SetResourceType(v string) *TagSpecification {
	s.ResourceType = &v
	return s
}

// SetTags sets the Tags field's value.
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
