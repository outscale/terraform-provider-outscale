package lbu

import (
	"time"
)

// CreateLoadBalancerInput ...
type CreateLoadBalancerInput struct {
	_ struct{} `type:"structure"`

	AvailabilityZones []*string `type:"list"`

	Listeners []*Listener `type:"list"`

	LoadBalancerName *string `type:"string" required:"true"`

	Scheme *string `type:"string"`

	SecurityGroups []*string `type:"list"`

	Subnets []*string `type:"list"`

	Tags []*Tag `type:"list"`
}

// Tag ...
type Tag struct {
	_ struct{} `type:"structure"`

	Key *string `min:"1" type:"string" required:"true"`

	Value *string `type:"string"`
}

// CreateLoadBalancerListenersInput ...
type CreateLoadBalancerListenersInput struct {
	_ struct{} `type:"structure"`

	Listeners []*Listener `type:"list"`

	LoadBalancerName *string `type:"string" required:"true"`
}

// CreateLoadBalancerListenersOutput ...
type CreateLoadBalancerListenersOutput struct {
	_ struct{} `type:"structure"`
}

// CreateLoadBalancerOutput ...
type CreateLoadBalancerOutput struct {
	_ struct{} `type:"structure"`

	DNSName *string ` type:"string"`
}

// Listener ...
type Listener struct {
	_ struct{} `type:"structure"`

	InstancePort *int64 `min:"1" type:"integer" required:"true"`

	InstanceProtocol *string `type:"string"`

	LoadBalancerPort *int64 `type:"integer" required:"true"`

	Protocol *string `type:"string" required:"true"`

	SSLCertificateId *string `type:"string"`
}

// DescribeLoadBalancersInput ...
type DescribeLoadBalancersInput struct {
	_ struct{} `type:"structure"`

	LoadBalancerNames []*string `type:"list"`

	Marker *string `type:"string"`

	PageSize *int64 `min:"1" type:"integer"`
}

// DescribeLoadBalancersOutput ...
type DescribeLoadBalancersOutput struct {
	_ struct{} `type:"structure"`

	LoadBalancerDescriptions []*LoadBalancerDescription `type:"list"`

	NextMarker *string `type:"string"`

	ResponseMetadata *RequestID `type:"structre"`
}

// RequestID ...
type RequestID struct {
	RequestID *string `type:"string"`
}

// LoadBalancerDescription ...
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

// BackendServerDescription ...
type BackendServerDescription struct {
	_ struct{} `type:"structure"`

	// The port on which the EC2 instance is listening.
	InstancePort *int64 `min:"1" type:"integer"`

	// The names of the policies enabled for the EC2 instance.
	PolicyNames []*string `type:"list"`
}

// HealthCheck ...
type HealthCheck struct {
	_ struct{} `type:"structure"`

	HealthyThreshold *int64 `min:"2" type:"integer" required:"true"`

	Interval *int64 `min:"5" type:"integer" required:"true"`

	Target *string `type:"string" required:"true"`

	Timeout *int64 `min:"2" type:"integer" required:"true"`

	UnhealthyThreshold *int64 `min:"2" type:"integer" required:"true"`
}

// Instance ...
type Instance struct {
	_ struct{} `type:"structure"`

	InstanceId *string `type:"string"`
}

// ListenerDescription ...
type ListenerDescription struct {
	_ struct{} `type:"structure"`

	Listener *Listener `type:"structure"`

	PolicyNames []*string `type:"list"`
}

// SourceSecurityGroup ...
type SourceSecurityGroup struct {
	_ struct{} `type:"structure"`

	GroupName *string `type:"string"`

	OwnerAlias *string `type:"string"`
}

// Policies ...
type Policies struct {
	_ struct{} `type:"structure"`

	// The stickiness policies created using CreateAppCookieStickinessPolicy.
	AppCookieStickinessPolicies []*AppCookieStickinessPolicy `type:"list"`

	// The stickiness policies created using CreateLBCookieStickinessPolicy.
	LBCookieStickinessPolicies []*LBCookieStickinessPolicy `type:"list"`

	// The policies other than the stickiness policies.
	OtherPolicies []*string `type:"list"`
}

// AppCookieStickinessPolicy ...
type AppCookieStickinessPolicy struct {
	_ struct{} `type:"structure"`

	// The name of the application cookie used for stickiness.
	CookieName *string `type:"string"`

	// The mnemonic name for the policy being created. The name must be unique within
	// a set of policies for this load balancer.
	PolicyName *string `type:"string"`
}

// LBCookieStickinessPolicy ...
type LBCookieStickinessPolicy struct {
	_ struct{} `type:"structure"`

	// The time period, in seconds, after which the cookie should be considered
	// stale. If this parameter is not specified, the stickiness session lasts for
	// the duration of the browser session.
	CookieExpirationPeriod *int64 `type:"long"`

	// The name of the policy. This name must be unique within the set of policies
	// for this load balancer.
	PolicyName *string `type:"string"`
}

// DescribeLoadBalancerAttributesInput ...
type DescribeLoadBalancerAttributesInput struct {
	_ struct{} `type:"structure"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`
}

// DescribeLoadBalancerAttributesOutput ...
type DescribeLoadBalancerAttributesOutput struct {
	_ struct{} `type:"structure"`

	// Information about the load balancer attributes.
	LoadBalancerAttributes *LoadBalancerAttributes `type:"structure"`
}

// LoadBalancerAttributes ...
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

// AdditionalAttribute ...
type AdditionalAttribute struct {
	_ struct{} `type:"structure"`

	Key *string ` type:"string"`

	Value *string ` type:"string"`
}

// ConnectionDraining ...
type ConnectionDraining struct {
	_ struct{} `type:"structure"`

	Enabled *bool ` type:"boolean" required:"true"`

	Timeout *int64 ` type:"integer"`
}

// ConnectionSettings ...
type ConnectionSettings struct {
	_ struct{} `type:"structure"`

	IdleTimeout *int64 ` min:"1" type:"integer" required:"true"`
}

// CrossZoneLoadBalancing ...
type CrossZoneLoadBalancing struct {
	_ struct{} `type:"structure"`

	Enabled *bool ` type:"boolean" required:"true"`
}

// AccessLog ...
type AccessLog struct {
	_ struct{} `type:"structure"`

	EmitInterval *int64 ` type:"integer"`

	Enabled *bool ` type:"boolean" required:"true"`

	S3BucketName *string ` type:"string"`

	S3BucketPrefix *string ` type:"string"`
}

// DeleteLoadBalancerListenersInput ...
type DeleteLoadBalancerListenersInput struct {
	_ struct{} `type:"structure"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`

	// The client port numbers of the listeners.
	//
	// LoadBalancerPorts is a required field
	LoadBalancerPorts []*int64 `type:"list" required:"true"`
}

// DeleteLoadBalancerListenersOutput ...
type DeleteLoadBalancerListenersOutput struct {
	_ struct{} `type:"structure"`
}

// ConfigureHealthCheckInput ...
type ConfigureHealthCheckInput struct {
	_ struct{} `type:"structure"`

	// The configuration information.
	//
	// HealthCheck is a required field
	HealthCheck *HealthCheck `type:"structure" required:"true"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`
}

// ConfigureHealthCheckOutput ...
type ConfigureHealthCheckOutput struct {
	_ struct{} `type:"structure"`

	// The updated health check.
	HealthCheck *HealthCheck `type:"structure"`
}

// ApplySecurityGroupsToLoadBalancerInput ...
type ApplySecurityGroupsToLoadBalancerInput struct {
	_ struct{} `type:"structure"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`

	// The IDs of the security groups to associate with the load balancer. Note
	// that you cannot specify the name of the security group.
	//
	// SecurityGroups is a required field
	SecurityGroups []*string `type:"list" required:"true"`
}

// ApplySecurityGroupsToLoadBalancerOutput ...
type ApplySecurityGroupsToLoadBalancerOutput struct {
	_ struct{} `type:"structure"`

	// The IDs of the security groups associated with the load balancer.
	SecurityGroups []*string `type:"list"`
}

// EnableAvailabilityZonesForLoadBalancerInput ...
type EnableAvailabilityZonesForLoadBalancerInput struct {
	_ struct{} `type:"structure"`

	// The Availability Zones. These must be in the same region as the load balancer.
	//
	// AvailabilityZones is a required field
	AvailabilityZones []*string `type:"list" required:"true"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`
}

// EnableAvailabilityZonesForLoadBalancerOutput ...
type EnableAvailabilityZonesForLoadBalancerOutput struct {
	_ struct{} `type:"structure"`

	// The updated list of Availability Zones for the load balancer.
	AvailabilityZones []*string `type:"list"`
}

// DisableAvailabilityZonesForLoadBalancerInput ...
type DisableAvailabilityZonesForLoadBalancerInput struct {
	_ struct{} `type:"structure"`

	// The Availability Zones.
	//
	// AvailabilityZones is a required field
	AvailabilityZones []*string `type:"list" required:"true"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`
}

// DisableAvailabilityZonesForLoadBalancerOutput ...
type DisableAvailabilityZonesForLoadBalancerOutput struct {
	_ struct{} `type:"structure"`

	// The remaining Availability Zones for the load balancer.
	AvailabilityZones []*string `type:"list"`
}

// AttachLoadBalancerToSubnetsInput ...
type AttachLoadBalancerToSubnetsInput struct {
	_ struct{} `type:"structure"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`

	// The IDs of the subnets to add. You can add only one subnet per Availability
	// Zone.
	//
	// Subnets is a required field
	Subnets []*string `type:"list" required:"true"`
}

// AttachLoadBalancerToSubnetsOutput ...
type AttachLoadBalancerToSubnetsOutput struct {
	_ struct{} `type:"structure"`

	// The IDs of the subnets attached to the load balancer.
	Subnets []*string `type:"list"`
}

// DeleteLoadBalancerInput ...
type DeleteLoadBalancerInput struct {
	_ struct{} `type:"structure"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`
}

// DeleteLoadBalancerOutput ...
type DeleteLoadBalancerOutput struct {
	_ struct{} `type:"structure"`
}

// RegisterInstancesWithLoadBalancerInput ...
type RegisterInstancesWithLoadBalancerInput struct {
	_ struct{} `type:"structure"`

	// The IDs of the instances.
	//
	// Instances is a required field
	Instances []*Instance `type:"list" required:"true"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`
}

// RegisterInstancesWithLoadBalancerOutput ...
type RegisterInstancesWithLoadBalancerOutput struct {
	_ struct{} `type:"structure"`

	// The updated list of instances for the load balancer.
	Instances []*Instance `type:"list"`
}

// DeregisterInstancesFromLoadBalancerInput ...
type DeregisterInstancesFromLoadBalancerInput struct {
	_ struct{} `type:"structure"`

	// The IDs of the instances.
	//
	// Instances is a required field
	Instances []*Instance `type:"list" required:"true"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`
}

// DeregisterInstancesFromLoadBalancerOutput ...
type DeregisterInstancesFromLoadBalancerOutput struct {
	_ struct{} `type:"structure"`

	// The remaining instances registered with the load balancer.
	Instances []*Instance `type:"list"`
}

// DetachLoadBalancerFromSubnetsInput ...
type DetachLoadBalancerFromSubnetsInput struct {
	_ struct{} `type:"structure"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`

	// The IDs of the subnets.
	//
	// Subnets is a required field
	Subnets []*string `type:"list" required:"true"`
}

// DetachLoadBalancerFromSubnetsOutput ...
type DetachLoadBalancerFromSubnetsOutput struct {
	_ struct{} `type:"structure"`

	// The IDs of the remaining subnets for the load balancer.
	Subnets []*string `type:"list"`
}

// CreateLBCookieStickinessPolicyInput ...
type CreateLBCookieStickinessPolicyInput struct {
	_ struct{} `type:"structure"`

	// The time period, in seconds, after which the cookie should be considered
	// stale. If you do not specify this parameter, the default value is 0, which
	// indicates that the sticky session should last for the duration of the browser
	// session.
	CookieExpirationPeriod *int64 `type:"long"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`

	// The name of the policy being created. Policy names must consist of alphanumeric
	// characters and dashes (-). This name must be unique within the set of policies
	// for this load balancer.
	//
	// PolicyName is a required field
	PolicyName *string `type:"string" required:"true"`
}

// CreateLBCookieStickinessPolicyOutput ...
type CreateLBCookieStickinessPolicyOutput struct {
	_ struct{} `type:"structure"`
}

// CreateAppCookieStickinessPolicyInput ...
type CreateAppCookieStickinessPolicyInput struct {
	_ struct{} `type:"structure"`

	// The name of the application cookie used for stickiness.
	//
	// CookieName is a required field
	CookieName *string `type:"string" required:"true"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`

	// The name of the policy being created. Policy names must consist of alphanumeric
	// characters and dashes (-). This name must be unique within the set of policies
	// for this load balancer.
	//
	// PolicyName is a required field
	PolicyName *string `type:"string" required:"true"`
}

// CreateAppCookieStickinessPolicyOutput ...
type CreateAppCookieStickinessPolicyOutput struct {
	_ struct{} `type:"structure"`
}

// SetLoadBalancerPoliciesOfListenerInput ...
type SetLoadBalancerPoliciesOfListenerInput struct {
	_ struct{} `type:"structure"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`

	// The external port of the load balancer.
	//
	// LoadBalancerPort is a required field
	LoadBalancerPort *int64 `type:"integer" required:"true"`

	// The names of the policies. This list must include all policies to be enabled.
	// If you omit a policy that is currently enabled, it is disabled. If the list
	// is empty, all current policies are disabled.
	//
	// PolicyNames is a required field
	PolicyNames []*string `type:"list" required:"true"`
}

// SetLoadBalancerPoliciesOfListenerOutput ...
type SetLoadBalancerPoliciesOfListenerOutput struct {
	_ struct{} `type:"structure"`
}

// DescribeLoadBalancerPoliciesInput ...
type DescribeLoadBalancerPoliciesInput struct {
	_ struct{} `type:"structure"`

	// The name of the load balancer.
	LoadBalancerName *string `type:"string"`

	// The names of the policies.
	PolicyNames []*string `type:"list"`
}

// DescribeLoadBalancerPoliciesOutput ...
type DescribeLoadBalancerPoliciesOutput struct {
	_ struct{} `type:"structure"`

	// Information about the policies.
	PolicyDescriptions []*PolicyDescription `type:"list"`
}

// PolicyDescription ...
type PolicyDescription struct {
	_ struct{} `type:"structure"`

	// The policy attributes.
	PolicyAttributeDescriptions []*PolicyAttributeDescription `type:"list"`

	// The name of the policy.
	PolicyName *string `type:"string"`

	// The name of the policy type.
	PolicyTypeName *string `type:"string"`
}

// PolicyAttributeDescription ...
type PolicyAttributeDescription struct {
	_ struct{} `type:"structure"`

	// The name of the attribute.
	AttributeName *string `type:"string"`

	// The value of the attribute.
	AttributeValue *string `type:"string"`
}

// DeleteLoadBalancerPolicyInput ...
type DeleteLoadBalancerPolicyInput struct {
	_ struct{} `type:"structure"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`

	// The name of the policy.
	//
	// PolicyName is a required field
	PolicyName *string `type:"string" required:"true"`
}

// DeleteLoadBalancerPolicyOutput ...
type DeleteLoadBalancerPolicyOutput struct {
	_ struct{} `type:"structure"`
}

// ModifyLoadBalancerAttributesInput ...
type ModifyLoadBalancerAttributesInput struct {
	_ struct{} `type:"structure"`

	// The attributes of the load balancer.
	//
	// LoadBalancerAttributes is a required field
	LoadBalancerAttributes *LoadBalancerAttributes `type:"structure" required:"true"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`
}

// ModifyLoadBalancerAttributesOutput ...
type ModifyLoadBalancerAttributesOutput struct {
	_ struct{} `type:"structure"`

	// The attributes for a load balancer.
	LoadBalancerAttributes *LoadBalancerAttributes `type:"structure"`

	// The name of the load balancer.
	LoadBalancerName *string `type:"string"`
}

// AddTagsInput ...
type AddTagsInput struct {
	_ struct{} `type:"structure"`

	// The name of the load balancer. You can specify one load balancer only.
	//
	// LoadBalancerNames is a required field
	LoadBalancerNames []*string `type:"list"`

	// The tags.
	//
	// Tags is a required field
	Tags []*Tag `type:"list"`
}

// AddTagsOutput ...
type AddTagsOutput struct {
	_ struct{} `type:"structure"`
}

// DescribeTagsInput ...
type DescribeTagsInput struct {
	_ struct{} `type:"structure"`

	// The names of the load balancers.
	//
	// LoadBalancerNames is a required field
	LoadBalancerNames []*string `min:"1" type:"list" required:"true"`
}

// DescribeTagsOutput ...
type DescribeTagsOutput struct {
	_ struct{} `type:"structure"`

	// Information about the tags.
	TagDescriptions []*TagDescription `type:"list"`
}

// TagDescription ...
type TagDescription struct {
	_ struct{} `type:"structure"`

	// The name of the load balancer.
	LoadBalancerName *string `type:"string"`

	// The tags.
	Tags []*Tag `min:"1" type:"list"`
}

// RemoveTagsInput ...
type RemoveTagsInput struct {
	_ struct{} `type:"structure"`

	// The name of the load balancer. You can specify a maximum of one load balancer
	// name.
	//
	// LoadBalancerNames is a required field
	LoadBalancerNames []*string `type:"list" required:"true"`

	// The list of tag keys to remove.
	//
	// Tags is a required field
	Tags []*TagKeyOnly `min:"1" type:"list" required:"true"`
}

// RemoveTagsOutput ...
type RemoveTagsOutput struct {
	_ struct{} `type:"structure"`
}

// TagKeyOnly ...
type TagKeyOnly struct {
	_ struct{} `type:"structure"`

	// The name of the key.
	Key *string `min:"1" type:"string"`
}

// DescribeInstanceHealthInput ...
type DescribeInstanceHealthInput struct {
	_ struct{} `type:"structure"`

	// The IDs of the instances.
	Instances []*Instance `type:"list"`

	// The name of the load balancer.
	//
	// LoadBalancerName is a required field
	LoadBalancerName *string `type:"string" required:"true"`
}

// DescribeInstanceHealthOutput ...
type DescribeInstanceHealthOutput struct {
	_ struct{} `type:"structure"`

	// Information about the health of the instances.
	InstanceStates []*InstanceState `type:"list"`
}

// InstanceState ...
type InstanceState struct {
	_ struct{} `type:"structure"`

	Description *string `type:"string"`

	// The ID of the instance.
	InstanceId *string `type:"string"`

	// Information about the cause of OutOfService instances. Specifically, whether
	// the cause is Elastic Load Balancing or the instance.
	//
	// Valid values: ELB | Instance | N/A
	ReasonCode *string `type:"string"`

	// The current state of the instance.
	//
	// Valid values: InService | OutOfService | Unknown
	State *string `type:"string"`
}
