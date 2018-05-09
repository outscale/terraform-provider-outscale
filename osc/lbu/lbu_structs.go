package lbu

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/terraform-providers/terraform-provider-outscale/osc/common"
)

type CreateLoadBalancerInput struct {
	_ struct{} `type:"structure"`

	// One or more Availability Zones from the same region as the load balancer.
	//
	// You must specify at least one Availability Zone.
	//
	// You can add more Availability Zones after you create the load balancer using
	// EnableAvailabilityZonesForLoadBalancer.
	AvailabilityZones []*string `type:"list"`

	// The listeners.
	//
	// For more information, see Listeners for Your Classic Load Balancer (http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/elb-listener-config.html)
	// in the Classic Load Balancers Guide.
	//
	// Listeners is a required field
	Listeners []*Listener `type:"list" required:"true"`

	// The name of the load balancer.
	//
	// This name must be unique within your set of load balancers for the region,
	// must have a maximum of 32 characters, must contain only alphanumeric characters
	// or hyphens, and cannot begin or end with a hyphen.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`

	// The type of a load balancer. Valid only for load balancers in a VPC.
	//
	// By default, Elastic Load Balancing creates an Internet-facing load balancer
	// with a DNS name that resolves to public IP addresses. For more information
	// about Internet-facing and Internal load balancers, see Load Balancer Scheme
	// (http://docs.aws.amazon.com/elasticloadbalancing/latest/userguide/how-elastic-load-balancing-works.html#load-balancer-scheme)
	// in the Elastic Load Balancing User Guide.
	//
	// Specify internal to create a load balancer with a DNS name that resolves
	// to private IP addresses.
	Scheme *string `type:"string"`

	// The IDs of the security groups to assign to the load balancer.
	SecurityGroups []*string `type:"list"`

	// The IDs of the subnets in your VPC to attach to the load balancer. Specify
	// one subnet per Availability Zone specified in AvailabilityZones.
	Subnets []*string `type:"list"`

	// A list of tags to assign to the load balancer.
	//
	// For more information about tagging your load balancer, see Tag Your Classic
	// Load Balancer (http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/add-remove-tags.html)
	// in the Classic Load Balancers Guide.
	Tags []*common.Tag `min:"1" type:"list"`
}

// String returns the string representation
func (s CreateLoadBalancerInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateLoadBalancerInput) GoString() string {
	return s.String()
}

// SetAvailabilityZones sets the AvailabilityZones field's value.
func (s *CreateLoadBalancerInput) SetAvailabilityZones(v []*string) *CreateLoadBalancerInput {
	s.AvailabilityZones = v
	return s
}

// SetListeners sets the Listeners field's value.
func (s *CreateLoadBalancerInput) SetListeners(v []*Listener) *CreateLoadBalancerInput {
	s.Listeners = v
	return s
}

// SetLoadBalancerName sets the LoadBalancerName field's value.
func (s *CreateLoadBalancerInput) SetLoadBalancerName(v string) *CreateLoadBalancerInput {
	s.LoadBalancerName = &v
	return s
}

// SetScheme sets the Scheme field's value.
func (s *CreateLoadBalancerInput) SetScheme(v string) *CreateLoadBalancerInput {
	s.Scheme = &v
	return s
}

// SetSecurityGroups sets the SecurityGroups field's value.
func (s *CreateLoadBalancerInput) SetSecurityGroups(v []*string) *CreateLoadBalancerInput {
	s.SecurityGroups = v
	return s
}

// SetSubnets sets the Subnets field's value.
func (s *CreateLoadBalancerInput) SetSubnets(v []*string) *CreateLoadBalancerInput {
	s.Subnets = v
	return s
}

// SetTags sets the Tags field's value.
func (s *CreateLoadBalancerInput) SetTags(v []*common.Tag) *CreateLoadBalancerInput {
	s.Tags = v
	return s
}

// Contains the parameters for CreateLoadBalancerListeners.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/elasticloadbalancing-2012-06-01/CreateLoadBalancerListenerInput
type CreateLoadBalancerListenersInput struct {
	_ struct{} `type:"structure"`

	// The listeners.
	//
	// Listeners is a required field
	Listeners []*Listener `type:"list" required:"true"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`
}

// String returns the string representation
func (s CreateLoadBalancerListenersInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateLoadBalancerListenersInput) GoString() string {
	return s.String()
}

// SetListeners sets the Listeners field's value.
func (s *CreateLoadBalancerListenersInput) SetListeners(v []*Listener) *CreateLoadBalancerListenersInput {
	s.Listeners = v
	return s
}

// SetLoadBalancerName sets the LoadBalancerName field's value.
func (s *CreateLoadBalancerListenersInput) SetLoadBalancerName(v string) *CreateLoadBalancerListenersInput {
	s.LoadBalancerName = &v
	return s
}

// Contains the parameters for CreateLoadBalancerListener.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/elasticloadbalancing-2012-06-01/CreateLoadBalancerListenerOutput
type CreateLoadBalancerListenersOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s CreateLoadBalancerListenersOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateLoadBalancerListenersOutput) GoString() string {
	return s.String()
}

// Contains the output for CreateLoadBalancer.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/elasticloadbalancing-2012-06-01/CreateAccessPointOutput
type CreateLoadBalancerOutput struct {
	_ struct{} `type:"structure"`

	// The DNS name of the load balancer.
	DNSName *string `type:"string"`
}

// String returns the string representation
func (s CreateLoadBalancerOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateLoadBalancerOutput) GoString() string {
	return s.String()
}

// SetDNSName sets the DNSName field's value.
func (s *CreateLoadBalancerOutput) SetDNSName(v string) *CreateLoadBalancerOutput {
	s.DNSName = &v
	return s
}

type Listener struct {
	_ struct{} `type:"structure"`

	// The port on which the instance is listening.
	//
	// InstancePort is a required field
	InstancePort *int64 `min:"1" type:"integer" required:"true"`

	// The protocol to use for routing traffic to instances: HTTP, HTTPS, TCP, or
	// SSL.
	//
	// If the front-end protocol is HTTP, HTTPS, TCP, or SSL, InstanceProtocol must
	// be at the same protocol.
	//
	// If there is another listener with the same InstancePort whose InstanceProtocol
	// is secure, (HTTPS or SSL), the listener's InstanceProtocol must also be secure.
	//
	// If there is another listener with the same InstancePort whose InstanceProtocol
	// is HTTP or TCP, the listener's InstanceProtocol must be HTTP or TCP.
	InstanceProtocol *string `type:"string"`

	// The port on which the load balancer is listening. On EC2-VPC, you can specify
	// any port from the range 1-65535. On EC2-Classic, you can specify any port
	// from the following list: 25, 80, 443, 465, 587, 1024-65535.
	//
	// LoadBalancerPort is a required field
	LoadBalancerPort *int64 `type:"integer" required:"true"`

	// The load balancer transport protocol to use for routing: HTTP, HTTPS, TCP,
	// or SSL.
	//
	// Protocol is a required field
	Protocol *string `type:"string" required:"true"`

	// The Amazon Resource Name (ARN) of the server certificate.
	SSLCertificateId *string `type:"string"`
}

// String returns the string representation
func (s Listener) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Listener) GoString() string {
	return s.String()
}

// SetInstancePort sets the InstancePort field's value.
func (s *Listener) SetInstancePort(v int64) *Listener {
	s.InstancePort = &v
	return s
}

// SetInstanceProtocol sets the InstanceProtocol field's value.
func (s *Listener) SetInstanceProtocol(v string) *Listener {
	s.InstanceProtocol = &v
	return s
}

// SetLoadBalancerPort sets the LoadBalancerPort field's value.
func (s *Listener) SetLoadBalancerPort(v int64) *Listener {
	s.LoadBalancerPort = &v
	return s
}

// SetProtocol sets the Protocol field's value.
func (s *Listener) SetProtocol(v string) *Listener {
	s.Protocol = &v
	return s
}

// SetSSLCertificateId sets the SSLCertificateId field's value.
func (s *Listener) SetSSLCertificateId(v string) *Listener {
	s.SSLCertificateId = &v
	return s
}

type DescribeLoadBalancersInput struct {
	_ struct{} `type:"structure"`

	// The names of the load balancers.
	LoadBalancerNames []*string `type:"list"`

	// The marker for the next set of results. (You received this marker from a
	// previous call.)
	Marker *string `type:"string"`

	// The maximum number of results to return with this call (a number from 1 to
	// 400). The default is 400.
	PageSize *int64 `min:"1" type:"integer"`
}

type DescribeLoadBalancersOutput struct {
	_ struct{} `type:"structure"`

	// Information about the load balancers.
	LoadBalancerDescriptions []*LoadBalancerDescription `type:"list"`

	// The marker to use when requesting the next set of results. If there are no
	// additional results, the string is empty.
	NextMarker *string `type:"string"`
}

type LoadBalancerDescription struct {
	_ struct{} `type:"structure"`

	// The Availability Zones for the load balancer.
	AvailabilityZones []*string `type:"list"`

	// Information about your EC2 instances.
	BackendServerDescriptions []*BackendServerDescription `type:"list"`

	// The DNS name of the load balancer.
	//
	// For more information, see Configure a Custom Domain Name (http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/using-domain-names-with-elb.html)
	// in the Classic Load Balancers Guide.
	CanonicalHostedZoneName *string `type:"string"`

	// The ID of the Amazon Route 53 hosted zone for the load balancer.
	CanonicalHostedZoneNameID *string `type:"string"`

	// The date and time the load balancer was created.
	CreatedTime *time.Time `type:"timestamp" timestampFormat:"iso8601"`

	// The DNS name of the load balancer.
	DNSName *string `type:"string"`

	// Information about the health checks conducted on the load balancer.
	HealthCheck *HealthCheck `type:"structure"`

	// The IDs of the instances for the load balancer.
	Instances []*Instance `type:"list"`

	// The listeners for the load balancer.
	ListenerDescriptions []*ListenerDescription `type:"list"`

	// The name of the load balancer.
	LoadBalancerName *string `type:"string"`

	// The policies defined for the load balancer.
	Policies *Policies `type:"structure"`

	// The type of load balancer. Valid only for load balancers in a VPC.
	//
	// If Scheme is internet-facing, the load balancer has a public DNS name that
	// resolves to a public IP address.
	//
	// If Scheme is internal, the load balancer has a public DNS name that resolves
	// to a private IP address.
	Scheme *string `type:"string"`

	// The security groups for the load balancer. Valid only for load balancers
	// in a VPC.
	SecurityGroups []*string `type:"list"`

	// The security group for the load balancer, which you can use as part of your
	// inbound rules for your registered instances. To only allow traffic from load
	// balancers, add a security group rule that specifies this source security
	// group as the inbound source.
	SourceSecurityGroup *SourceSecurityGroup `type:"structure"`

	// The IDs of the subnets for the load balancer.
	Subnets []*string `type:"list"`

	// The ID of the VPC for the load balancer.
	VPCId *string `type:"string"`
}

type BackendServerDescription struct {
	_ struct{} `type:"structure"`

	// The port on which the EC2 instance is listening.
	InstancePort *int64 `min:"1" type:"integer"`

	// The names of the policies enabled for the EC2 instance.
	PolicyNames []*string `type:"list"`
}

type HealthCheck struct {
	_ struct{} `type:"structure"`

	HealthyThreshold *int64 `min:"2" type:"integer" required:"true"`

	Interval *int64 `min:"5" type:"integer" required:"true"`

	Target *string `type:"string" required:"true"`

	Timeout *int64 `min:"2" type:"integer" required:"true"`

	UnhealthyThreshold *int64 `min:"2" type:"integer" required:"true"`
}

type Instance struct {
	_ struct{} `type:"structure"`

	// The instance ID.
	InstanceId *string `type:"string"`
}
type ListenerDescription struct {
	_ struct{} `type:"structure"`

	Listener *Listener `type:"structure"`

	PolicyNames []*string `type:"list"`
}

type SourceSecurityGroup struct {
	_ struct{} `type:"structure"`

	GroupName *string `type:"string"`

	OwnerAlias *string `type:"string"`
}

type Policies struct {
	_ struct{} `type:"structure"`

	AppCookieStickinessPolicies []*AppCookieStickinessPolicy `type:"list"`

	LBCookieStickinessPolicies []*LBCookieStickinessPolicy `type:"list"`

	OtherPolicies []*string `type:"list"`
}

type AppCookieStickinessPolicy struct {
	_ struct{} `type:"structure"`

	CookieName *string `type:"string"`

	PolicyName *string `type:"string"`
}

type LBCookieStickinessPolicy struct {
	_ struct{} `type:"structure"`

	CookieExpirationPeriod *int64 `type:"long"`

	PolicyName *string `type:"string"`
}

type DescribeLoadBalancerAttributesInput struct {
	_ struct{} `type:"structure"`

	LoadBalancerName *string `type:"string" required:"true"`
}

// String returns the string representation
func (s DescribeLoadBalancerAttributesInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeLoadBalancerAttributesInput) GoString() string {
	return s.String()
}

// SetLoadBalancerName sets the LoadBalancerName field's value.
func (s *DescribeLoadBalancerAttributesInput) SetLoadBalancerName(v string) *DescribeLoadBalancerAttributesInput {
	s.LoadBalancerName = &v
	return s
}

type DescribeLoadBalancerAttributesOutput struct {
	_ struct{} `type:"structure"`

	// Information about the load balancer attributes.
	LoadBalancerAttributes *LoadBalancerAttributes `type:"structure"`
}

// String returns the string representation
func (s DescribeLoadBalancerAttributesOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeLoadBalancerAttributesOutput) GoString() string {
	return s.String()
}

// SetLoadBalancerAttributes sets the LoadBalancerAttributes field's value.
func (s *DescribeLoadBalancerAttributesOutput) SetLoadBalancerAttributes(v *LoadBalancerAttributes) *DescribeLoadBalancerAttributesOutput {
	s.LoadBalancerAttributes = v
	return s
}

type LoadBalancerAttributes struct {
	_ struct{} `type:"structure"`

	// If enabled, the load balancer captures detailed information of all requests
	// and delivers the information to the Amazon S3 bucket that you specify.
	//
	// For more information, see Enable Access Logs (http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/enable-access-logs.html)
	// in the Classic Load Balancers Guide.
	AccessLog *AccessLog `type:"structure"`

	// This parameter is reserved.
	AdditionalAttributes []*AdditionalAttribute `type:"list"`

	// If enabled, the load balancer allows existing requests to complete before
	// the load balancer shifts traffic away from a deregistered or unhealthy instance.
	//
	// For more information, see Configure Connection Draining (http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/config-conn-drain.html)
	// in the Classic Load Balancers Guide.
	ConnectionDraining *ConnectionDraining `type:"structure"`

	// If enabled, the load balancer allows the connections to remain idle (no data
	// is sent over the connection) for the specified duration.
	//
	// By default, Elastic Load Balancing maintains a 60-second idle connection
	// timeout for both front-end and back-end connections of your load balancer.
	// For more information, see Configure Idle Connection Timeout (http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/config-idle-timeout.html)
	// in the Classic Load Balancers Guide.
	ConnectionSettings *ConnectionSettings `type:"structure"`

	// If enabled, the load balancer routes the request traffic evenly across all
	// instances regardless of the Availability Zones.
	//
	// For more information, see Configure Cross-Zone Load Balancing (http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/enable-disable-crosszone-lb.html)
	// in the Classic Load Balancers Guide.
	CrossZoneLoadBalancing *CrossZoneLoadBalancing `type:"structure"`
}

type AdditionalAttribute struct {
	_ struct{} `type:"structure"`

	// This parameter is reserved.
	Key *string `type:"string"`

	// This parameter is reserved.
	Value *string `type:"string"`
}

type ConnectionDraining struct {
	_ struct{} `type:"structure"`

	// Specifies whether connection draining is enabled for the load balancer.
	//
	// Enabled is a required field
	Enabled *bool `type:"boolean" required:"true"`

	// The maximum time, in seconds, to keep the existing connections open before
	// deregistering the instances.
	Timeout *int64 `type:"integer"`
}

type ConnectionSettings struct {
	_ struct{} `type:"structure"`

	// The time, in seconds, that the connection is allowed to be idle (no data
	// has been sent over the connection) before it is closed by the load balancer.
	//
	// IdleTimeout is a required field
	IdleTimeout *int64 `min:"1" type:"integer" required:"true"`
}

type CrossZoneLoadBalancing struct {
	_ struct{} `type:"structure"`

	// Specifies whether cross-zone load balancing is enabled for the load balancer.
	//
	// Enabled is a required field
	Enabled *bool `type:"boolean" required:"true"`
}

type AccessLog struct {
	_ struct{} `type:"structure"`

	// The interval for publishing the access logs. You can specify an interval
	// of either 5 minutes or 60 minutes.
	//
	// Default: 60 minutes
	EmitInterval *int64 `type:"integer"`

	// Specifies whether access logs are enabled for the load balancer.
	//
	// Enabled is a required field
	Enabled *bool `type:"boolean" required:"true"`

	// The name of the Amazon S3 bucket where the access logs are stored.
	S3BucketName *string `type:"string"`

	// The logical hierarchy you created for your Amazon S3 bucket, for example
	// my-bucket-prefix/prod. If the prefix is not provided, the log is placed at
	// the root level of the bucket.
	S3BucketPrefix *string `type:"string"`
}
