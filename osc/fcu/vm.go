package fcu

import (
	"context"
	"net/http"

	"github.com/terraform-providers/terraform-provider-outscale/osc"
)

//VMOperations defines all the operations needed for FCU VMs
type VMOperations struct {
	client *osc.Client
}

//VMService all the necessary actions for them VM service
type VMService interface {
	RunInstance(input *RunInstancesInput) (*Reservation, error)
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
	DisableAPITermination *bool `locationName:"disableApiTermination" type:"boolean"`

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

	// The name of the key pair. You can create a key pair using CreateKeyPair or
	// ImportKeyPair.
	//
	// If you do not specify a key pair, you can't connect to the instance unless
	// you choose an AMI that is configured to allow users another way to log in.
	KeyName *string `type:"string"`

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
}

//BlockDeviceMapping input to specify the mapping
type BlockDeviceMapping struct {
	_ struct{} `type:"structure"`

	// The device name (for example, /dev/sdh or xvdh).
	DeviceName *string `locationName:"deviceName" type:"string"`

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

const opRunInstances = "RunInstances"

func (v VMOperations) RunInstance(input *RunInstancesInput) (*Reservation, error) {
	req, err := v.client.NewRequest(context.Background(), opRunInstances, http.MethodGet, "/", input)
	if err != nil {
		return nil, err
	}

	output := &Reservation{}

	err = v.client.Do(context.Background(), req, output)
	if err != nil {
		return nil, err
	}

	return nil, err
}
