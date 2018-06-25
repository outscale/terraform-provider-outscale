package lbu

import (
	"time"
)

// CreateLoadBalancerInput ...
type CreateLoadBalancerInput struct {
	_                 struct{}    `type:"structure"`
	AvailabilityZones []*string   `type:"list"`
	Listeners         []*Listener `type:"list"`
	LoadBalancerName  *string     `type:"string" required:"true"`
	Scheme            *string     `type:"string"`
	SecurityGroups    []*string   `type:"list"`
	Subnets           []*string   `type:"list"`
	Tags              []*Tag      `type:"list"`
}

// Tag ...
type Tag struct {
	_     struct{} `type:"structure"`
	Key   *string  `min:"1" type:"string" required:"true"`
	Value *string  `type:"string"`
}

// CreateLoadBalancerListenersInput ...
type CreateLoadBalancerListenersInput struct {
	_                struct{}    `type:"structure"`
	Listeners        []*Listener `type:"list"`
	LoadBalancerName *string     `type:"string" required:"true"`
}

// CreateLoadBalancerListenersOutput ...
type CreateLoadBalancerListenersOutput struct {
	CreateLoadBalancerListenersResult *CreateLoadBalancerListenersResult `type:"structure"`
	ResponseMetadata                  *ResponseMetadata                  `type:"structure"`
}

//CreateLoadBalancerListenersResult ...
type CreateLoadBalancerListenersResult struct {
	_ struct{} `type:"structure"`
}

// CreateLoadBalancerOutput ...
type CreateLoadBalancerOutput struct {
	_       struct{} `type:"structure"`
	DNSName *string  ` type:"string"`
}

// Listener ...
type Listener struct {
	_                struct{} `type:"structure"`
	InstancePort     *int64   `min:"1" type:"integer" required:"true"`
	InstanceProtocol *string  `type:"string"`
	LoadBalancerPort *int64   `type:"integer" required:"true"`
	Protocol         *string  `type:"string" required:"true"`
	SSLCertificateId *string  `type:"string"`
}

// DescribeLoadBalancersInput ...
type DescribeLoadBalancersInput struct {
	_                 struct{}  `type:"structure"`
	LoadBalancerNames []*string `type:"list"`
	Marker            *string   `type:"string"`
	PageSize          *int64    `min:"1" type:"integer"`
}

// DescribeLoadBalancersOutput ...
type DescribeLoadBalancersOutput struct {
	DescribeLoadBalancersResult *DescribeLoadBalancersResult `type:"structure"`
	ResponseMetadata            *ResponseMetadata            `type:"structure"`
}

// DescribeLoadBalancersResult ...
type DescribeLoadBalancersResult struct {
	LoadBalancerDescriptions []*LoadBalancerDescription `type:"list"`
	NextMarker               *string                    `type:"string"`
}

// LoadBalancerDescription ...
type LoadBalancerDescription struct {
	_                         struct{}                    `type:"structure"`
	AvailabilityZones         []*string                   `type:"list"`                                // The Availability Zones for the load balancer.
	BackendServerDescriptions []*BackendServerDescription `type:"list"`                                // Information about your EC2 instances.
	CanonicalHostedZoneName   *string                     `type:"string"`                              // The DNS name of the load balancer. For more information, see Configure a Custom Domain Name (http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/using-domain-names-with-elb.html) in the Classic Load Balancers Guide.
	CanonicalHostedZoneNameID *string                     `type:"string"`                              // The ID of the Amazon Route 53 hosted zone for the load balancer.
	CreatedTime               *time.Time                  `type:"timestamp" timestampFormat:"iso8601"` // The date and time the load balancer was created.
	DNSName                   *string                     `type:"string"`                              // The DNS name of the load balancer.
	HealthCheck               *HealthCheck                `type:"structure"`                           // Information about the health checks conducted on the load balancer.
	Instances                 []*Instance                 `type:"list"`                                // The IDs of the instances for the load balancer.
	ListenerDescriptions      []*ListenerDescription      `type:"list"`                                // The listeners for the load balancer.
	LoadBalancerName          *string                     `type:"string"`                              // The name of the load balancer.
	Policies                  *Policies                   `type:"structure"`                           // The policies defined for the load balancer.
	Scheme                    *string                     `type:"string"`                              // The type of load balancer. Valid only for load balancers in a VPC.
	SecurityGroups            []*string                   `type:"list"`                                // The security groups for the load balancer. Valid only for load balancers in a VPC.
	SourceSecurityGroup       *SourceSecurityGroup        `type:"structure"`                           // The security group for the load balancer, which you can use as part of your inbound rules for your registered instances. To only allow traffic from load balancers, add a security group rule that specifies this source security group as the inbound source.
	Subnets                   []*string                   `type:"list"`                                // The IDs of the subnets for the load balancer.
	VPCId                     *string                     `type:"string"`                              // The ID of the VPC for the load balancer.
}

// BackendServerDescription ...
type BackendServerDescription struct {
	_            struct{}  `type:"structure"`
	InstancePort *int64    `min:"1" type:"integer"` // The port on which the EC2 instance is listening.
	PolicyNames  []*string `type:"list"`            // The names of the policies enabled for the EC2 instance.
}

// HealthCheck ...
type HealthCheck struct {
	_                  struct{} `type:"structure"`
	HealthyThreshold   *int64   `min:"2" type:"integer" required:"true"`
	Interval           *int64   `min:"5" type:"integer" required:"true"`
	Target             *string  `type:"string" required:"true"`
	Timeout            *int64   `min:"2" type:"integer" required:"true"`
	UnhealthyThreshold *int64   `min:"2" type:"integer" required:"true"`
}

// Instance ...
type Instance struct {
	_          struct{} `type:"structure"`
	InstanceId *string  `type:"string"`
}

// ListenerDescription ...
type ListenerDescription struct {
	_           struct{}  `type:"structure"`
	Listener    *Listener `type:"structure"`
	PolicyNames []*string `type:"list"`
}

// SourceSecurityGroup ...
type SourceSecurityGroup struct {
	_          struct{} `type:"structure"`
	GroupName  *string  `type:"string"`
	OwnerAlias *string  `type:"string"`
}

// Policies ...
type Policies struct {
	_                           struct{}                     `type:"structure"`
	AppCookieStickinessPolicies []*AppCookieStickinessPolicy `type:"list"` // The stickiness policies created using CreateAppCookieStickinessPolicy.
	LBCookieStickinessPolicies  []*LBCookieStickinessPolicy  `type:"list"` // The stickiness policies created using CreateLBCookieStickinessPolicy.
	OtherPolicies               []*string                    `type:"list"` // The policies other than the stickiness policies.
}

// AppCookieStickinessPolicy ...
type AppCookieStickinessPolicy struct {
	_          struct{} `type:"structure"`
	CookieName *string  `type:"string"` // The name of the application cookie used for stickiness.
	PolicyName *string  `type:"string"` // The mnemonic name for the policy being created. The name must be unique within a set of policies for this load balancer.
}

// LBCookieStickinessPolicy ...
type LBCookieStickinessPolicy struct {
	_                      struct{} `type:"structure"`
	CookieExpirationPeriod *int64   `type:"long"`   // The time period, in seconds, after which the cookie should be considered stale. If this parameter is not specified, the stickiness session lasts for the duration of the browser session.
	PolicyName             *string  `type:"string"` // The name of the policy. This name must be unique within the set of policies for this load balancer.
}

// DescribeLoadBalancerAttributesInput ...
type DescribeLoadBalancerAttributesInput struct {
	_                struct{} `type:"structure"`
	LoadBalancerName *string  `type:"string" required:"true"` //The name of the load balancer.
}

// DescribeLoadBalancerAttributesOutput ...
type DescribeLoadBalancerAttributesOutput struct {
	DescribeLoadBalancerAttributesResult *DescribeLoadBalancerAttributesResult `type:"structure"`
	ResponseMetadata                     *ResponseMetadata                     `type:"structure"`
}

// DescribeLoadBalancerAttributesResult ...
type DescribeLoadBalancerAttributesResult struct {
	_                      struct{}                `type:"structure"`
	LoadBalancerAttributes *LoadBalancerAttributes `type:"structure"` // Information about the load balancer attributes.
}

// LoadBalancerAttributes ...
type LoadBalancerAttributes struct {
	_                      struct{}                `type:"structure"`
	AccessLog              *AccessLog              `type:"structure"` //If enabled, the load balancer captures detailed information of all requests and delivers the information to the Amazon S3 bucket that you specify.
	AdditionalAttributes   []*AdditionalAttribute  `type:"list"`      // This parameter is reserved.
	ConnectionDraining     *ConnectionDraining     `type:"structure"` // If enabled, the load balancer allows existing requests to complete before the load balancer shifts traffic away from a deregistered or unhealthy instance.
	ConnectionSettings     *ConnectionSettings     `type:"structure"` // If enabled, the load balancer allows the connections to remain idle (no data is sent over the connection) for the specified duration.
	CrossZoneLoadBalancing *CrossZoneLoadBalancing `type:"structure"` // If enabled, the load balancer routes the request traffic evenly across all instances regardless of the Availability Zones.
}

// AdditionalAttribute ...
type AdditionalAttribute struct {
	_     struct{} `type:"structure"`
	Key   *string  ` type:"string"`
	Value *string  ` type:"string"`
}

// ConnectionDraining ...
type ConnectionDraining struct {
	_       struct{} `type:"structure"`
	Enabled *bool    ` type:"boolean" required:"true"`
	Timeout *int64   ` type:"integer"`
}

// ConnectionSettings ...
type ConnectionSettings struct {
	_           struct{} `type:"structure"`
	IdleTimeout *int64   ` min:"1" type:"integer" required:"true"`
}

// CrossZoneLoadBalancing ...
type CrossZoneLoadBalancing struct {
	_       struct{} `type:"structure"`
	Enabled *bool    ` type:"boolean" required:"true"`
}

// AccessLog ...
type AccessLog struct {
	_              struct{} `type:"structure"`
	EmitInterval   *int64   ` type:"integer"`
	Enabled        *bool    ` type:"boolean" required:"true"`
	S3BucketName   *string  ` type:"string"`
	S3BucketPrefix *string  ` type:"string"`
}

// DeleteLoadBalancerListenersInput ...
type DeleteLoadBalancerListenersInput struct {
	_                 struct{} `type:"structure"`
	LoadBalancerName  *string  `type:"string" required:"true"` //The name of the load balancer.
	LoadBalancerPorts []*int64 `type:"list" required:"true"`   // The client port numbers of the listeners.
}

// DeleteLoadBalancerListenersOutput ...
type DeleteLoadBalancerListenersOutput struct {
	DeleteLoadBalancerListenersResult *DeleteLoadBalancerListenersResult `type:"structure"`
	ResponseMetadata                  *ResponseMetadata                  `type:"structure"`
}

// DeleteLoadBalancerListenersResult ...
type DeleteLoadBalancerListenersResult struct {
	_ struct{} `type:"structure"`
}

// ConfigureHealthCheckInput ...
type ConfigureHealthCheckInput struct {
	_                struct{}     `type:"structure"`
	HealthCheck      *HealthCheck `type:"structure" required:"true"` //The configuration information. HealthCheck is a required field
	LoadBalancerName *string      `type:"string" required:"true"`    // The name of the load balancer.
}

// ConfigureHealthCheckOutput ...
type ConfigureHealthCheckOutput struct {
	ConfigureHealthCheckResult *ConfigureHealthCheckResult `type:"structure"` // The updated health check.
	ResponseMetadata           *ResponseMetadata           `type:"structure"`
}

// ConfigureHealthCheckResult ...
type ConfigureHealthCheckResult struct {
	_           struct{}     `type:"structure"`
	HealthCheck *HealthCheck `type:"structure"` // The updated health check.
}

// ApplySecurityGroupsToLoadBalancerInput ...
type ApplySecurityGroupsToLoadBalancerInput struct {
	_                struct{}  `type:"structure"`
	LoadBalancerName *string   `type:"string" required:"true"` // The name of the load balancer.
	SecurityGroups   []*string `type:"list" required:"true"`   // The IDs of the security groups to associate with the load balancer. Note that you cannot specify the name of the security group.
}

// ApplySecurityGroupsToLoadBalancerOutput ...
type ApplySecurityGroupsToLoadBalancerOutput struct {
	ApplySecurityGroupsToLoadBalancerResult *ApplySecurityGroupsToLoadBalancerResult `type:"structure"`
	ResponseMetadata                        *ResponseMetadata                        `type:"structure"`
}

// ApplySecurityGroupsToLoadBalancerResult ...
type ApplySecurityGroupsToLoadBalancerResult struct {
	_              struct{}  `type:"structure"`
	SecurityGroups []*string `type:"list"` // The IDs of the security groups associated with the load balancer.
}

// EnableAvailabilityZonesForLoadBalancerInput ...
type EnableAvailabilityZonesForLoadBalancerInput struct {
	_                 struct{}  `type:"structure"`
	AvailabilityZones []*string `type:"list" required:"true"`   //The Availability Zones. These must be in the same region as the load balancer.
	LoadBalancerName  *string   `type:"string" required:"true"` //The name of the load balancer.
}

// EnableAvailabilityZonesForLoadBalancerOutput ...
type EnableAvailabilityZonesForLoadBalancerOutput struct {
	EnableAvailabilityZonesForLoadBalancerResult *EnableAvailabilityZonesForLoadBalancerResult `type:"structure"`
	ResponseMetadata                             *ResponseMetadata                             `type:"structure"`
}

// EnableAvailabilityZonesForLoadBalancerResult ...
type EnableAvailabilityZonesForLoadBalancerResult struct {
	_                 struct{}  `type:"structure"`
	AvailabilityZones []*string `type:"list"` // The updated list of Availability Zones for the load balancer.
}

// DisableAvailabilityZonesForLoadBalancerInput ...
type DisableAvailabilityZonesForLoadBalancerInput struct {
	_                 struct{}  `type:"structure"`
	AvailabilityZones []*string `type:"list" required:"true"`   // The Availability Zones.
	LoadBalancerName  *string   `type:"string" required:"true"` // The name of the load balancer.
}

// DisableAvailabilityZonesForLoadBalancerOutput ...
type DisableAvailabilityZonesForLoadBalancerOutput struct {
	DisableAvailabilityZonesForLoadBalancerResult *DisableAvailabilityZonesForLoadBalancerResult `type:"structure"`
	ResponseMetadata                              *ResponseMetadata                              `type:"structure"`
}

// DisableAvailabilityZonesForLoadBalancerResult ...
type DisableAvailabilityZonesForLoadBalancerResult struct {
	_                 struct{}  `type:"structure"`
	AvailabilityZones []*string `type:"list"` // The remaining Availability Zones for the load balancer.
}

// AttachLoadBalancerToSubnetsInput ...
type AttachLoadBalancerToSubnetsInput struct {
	_                struct{}  `type:"structure"`
	LoadBalancerName *string   `type:"string" required:"true"` //The name of the load balancer.
	Subnets          []*string `type:"list" required:"true"`   //// The IDs of the subnets to add. You can add only one subnet per Availability Zone.
}

// AttachLoadBalancerToSubnetsOutput ...
type AttachLoadBalancerToSubnetsOutput struct {
	AttachLoadBalancerToSubnetsResult *AttachLoadBalancerToSubnetsResult `type:"structure"`
	ResponseMetadata                  *ResponseMetadata                  `type:"structure"`
}

// AttachLoadBalancerToSubnetsResult ...
type AttachLoadBalancerToSubnetsResult struct {
	_       struct{}  `type:"structure"`
	Subnets []*string `type:"list"` // The IDs of the subnets attached to the load balancer.
}

// DeleteLoadBalancerInput ...
type DeleteLoadBalancerInput struct {
	_                struct{} `type:"structure"`
	LoadBalancerName *string  `type:"string" required:"true"` //The name of the load balancer.
}

// DeleteLoadBalancerOutput ...
type DeleteLoadBalancerOutput struct {
	DeleteLoadBalancerResult *DeleteLoadBalancerResult `type:"structure"`
	ResponseMetadata         *ResponseMetadata         `type:"structure"`
}

// DeleteLoadBalancerResult ...
type DeleteLoadBalancerResult struct {
	_ struct{} `type:"structure"`
}

// RegisterInstancesWithLoadBalancerInput ...
type RegisterInstancesWithLoadBalancerInput struct {
	_                struct{}    `type:"structure"`
	Instances        []*Instance `type:"list" required:"true"`   //The IDs of the instances.
	LoadBalancerName *string     `type:"string" required:"true"` // The name of the load balancer.
}

// RegisterInstancesWithLoadBalancerOutput ...
type RegisterInstancesWithLoadBalancerOutput struct {
	RegisterInstancesWithLoadBalancerResult *RegisterInstancesWithLoadBalancerResult `type:"structure"`
	ResponseMetadata                        *ResponseMetadata                        `type:"structure"`
}

// RegisterInstancesWithLoadBalancerResult ...
type RegisterInstancesWithLoadBalancerResult struct {
	_         struct{}    `type:"structure"`
	Instances []*Instance `type:"list"` // The updated list of instances for the load balancer.
}

// DeregisterInstancesFromLoadBalancerInput ...
type DeregisterInstancesFromLoadBalancerInput struct {
	_                struct{}    `type:"structure"`
	Instances        []*Instance `type:"list" required:"true"`   // The IDs of the instances.
	LoadBalancerName *string     `type:"string" required:"true"` // The name of the load balancer.
}

// DeregisterInstancesFromLoadBalancerOutput ...
type DeregisterInstancesFromLoadBalancerOutput struct {
	DeregisterInstancesFromLoadBalancerResult *DeregisterInstancesFromLoadBalancerResult `type:"structure"`
	ResponseMetadata                          *ResponseMetadata                          `type:"structure"`
}

// DeregisterInstancesFromLoadBalancerResult ...
type DeregisterInstancesFromLoadBalancerResult struct {
	_         struct{}    `type:"structure"`
	Instances []*Instance `type:"list"` // The remaining instances registered with the load balancer.
}

// DetachLoadBalancerFromSubnetsInput ...
type DetachLoadBalancerFromSubnetsInput struct {
	_                struct{}  `type:"structure"`
	LoadBalancerName *string   `type:"string" required:"true"` // The name of the load balancer.
	Subnets          []*string `type:"list" required:"true"`   // The IDs of the subnets.
}

// DetachLoadBalancerFromSubnetsOutput ...
type DetachLoadBalancerFromSubnetsOutput struct {
	DetachLoadBalancerFromSubnetsResult *DetachLoadBalancerFromSubnetsResult `type:"structure"`
	ResponseMetadata                    *ResponseMetadata                    `type:"structure"`
}

// DetachLoadBalancerFromSubnetsResult ...
type DetachLoadBalancerFromSubnetsResult struct {
	_       struct{}  `type:"structure"`
	Subnets []*string `type:"list"` // The IDs of the remaining subnets for the load balancer.
}

// CreateLBCookieStickinessPolicyInput ...
type CreateLBCookieStickinessPolicyInput struct {
	_                      struct{} `type:"structure"`
	CookieExpirationPeriod *int64   `type:"long"`                   // The time period, in seconds, after which the cookie should be considered stale. If you do not specify this parameter, the default value is 0, which indicates that the sticky session should last for the duration of the browser session.
	LoadBalancerName       *string  `type:"string" required:"true"` // The name of the load balancer.
	PolicyName             *string  `type:"string" required:"true"` // The name of the policy being created. Policy names must consist of alphanumeric characters and dashes (-). This name must be unique within the set of policies for this load balancer.

}

// CreateLBCookieStickinessPolicyOutput ...
type CreateLBCookieStickinessPolicyOutput struct {
	CreateLBCookieStickinessPolicyResult *CreateLBCookieStickinessPolicyResult `type:"structure"`
	ResponseMetadata                     *ResponseMetadata                     `type:"structure"`
}

// CreateLBCookieStickinessPolicyResult ...
type CreateLBCookieStickinessPolicyResult struct {
	_ struct{} `type:"structure"`
}

// CreateAppCookieStickinessPolicyInput ...
type CreateAppCookieStickinessPolicyInput struct {
	_                struct{} `type:"structure"`
	CookieName       *string  `type:"string" required:"true"` // The name of the application cookie used for stickiness.
	LoadBalancerName *string  `type:"string" required:"true"` // The name of the load balancer.
	PolicyName       *string  `type:"string" required:"true"` // The name of the policy being created. Policy names must consist of alphanumeric characters and dashes (-). This name must be unique within the set of policies for this load balancer.
}

// CreateAppCookieStickinessPolicyOutput ...
type CreateAppCookieStickinessPolicyOutput struct {
	CreateAppCookieStickinessPolicyResult *CreateAppCookieStickinessPolicyResult `type:"structure"`
	ResponseMatadata                      *ResponseMetadata                      `type:"structure"`
}

//CreateAppCookieStickinessPolicyResult inner result
type CreateAppCookieStickinessPolicyResult struct {
	_ struct{} `type:"structure"`
}

//ResponseMetadata ...
type ResponseMetadata struct {
	RequestID *string `locationName:"RequestId" type:"string"`
}

// SetLoadBalancerPoliciesOfListenerInput ...
type SetLoadBalancerPoliciesOfListenerInput struct {
	_                struct{}  `type:"structure"`
	LoadBalancerName *string   `type:"string" required:"true"`  // The name of the load balancer.
	LoadBalancerPort *int64    `type:"integer" required:"true"` // The external port of the load balancer.
	PolicyNames      []*string `type:"list" required:"true"`    // The names of the policies. This list must include all policies to be enabled.
}

// SetLoadBalancerPoliciesOfListenerOutput ...
type SetLoadBalancerPoliciesOfListenerOutput struct {
	SetLoadBalancerPoliciesOfListenerResult SetLoadBalancerPoliciesOfListenerResult `type:"structure"`
	ResponseMatadata                        *ResponseMetadata                       `type:"structure"`
}

// SetLoadBalancerPoliciesOfListenerResult ...
type SetLoadBalancerPoliciesOfListenerResult struct {
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
	ModifyLoadBalancerAttributesResult *ModifyLoadBalancerAttributesResult `type:"structure"`
	ResponseMetadata                   *ResponseMetadata                   `type:"structure"`
}

// ModifyLoadBalancerAttributesResult ...
type ModifyLoadBalancerAttributesResult struct {
	_                      struct{}                `type:"structure"`
	LoadBalancerAttributes *LoadBalancerAttributes `type:"structure"` // The attributes for a load balancer.
	LoadBalancerName       *string                 `type:"string"`    // The name of the load balancer.
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

// SetLoadBalancerListenerSSLCertificateInput Contains the parameters for SetLoadBalancerListenerSSLCertificate.
type SetLoadBalancerListenerSSLCertificateInput struct {
	_                struct{} `type:"structure"`
	LoadBalancerName *string  `type:"string" required:"true"`  // The name of the load balancer. LoadBalancerName is a required field
	LoadBalancerPort *int64   `type:"integer" required:"true"` // The port that uses the specified SSL certificate. LoadBalancerPort is a required field
	SSLCertificateId *string  `type:"string" required:"true"`  // The Outscale Resource Name (ORN) of the SSL certificate. SSLCertificateId is a required field
}

// SetLoadBalancerListenerSSLCertificateOutput Contains the output of SetLoadBalancerListenerSSLCertificate.
type SetLoadBalancerListenerSSLCertificateOutput struct {
	_                struct{}          `type:"structure"`
	ResponseMetadata *ResponseMetadata `type:"structure"`
}
