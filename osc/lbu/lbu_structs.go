package lbu

import (
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/common"
)

// CreateAvailabilityZonesMember ...
type CreateAvailabilityZonesMember struct {
	Member []*string `locationName:"member" locationNameList:"member" type:"list"`
}

// CreateListenersMember ...
type CreateListenersMember struct {
	Member []*Listener `locationName:"member" locationNameList:"member" type:"list"`
}

// SecurityGroupsMember ...
type SecurityGroupsMember struct {
	Member []*string `locationName:"member" locationNameList:"member" type:"list"`
}

// SubnetsMember ...
type SubnetsMember struct {
	Member []*string `locationName:"member" locationNameList:"member" type:"list"`
}

// TagsMember ...
type TagsMember struct {
	Member []*common.Tag `locationName:"member" locationNameList:"member" type:"list"`
}

// CreateLoadBalancerInput ...
type CreateLoadBalancerInput struct {
	_ struct{} `type:"structure"`

	AvailabilityZones *CreateAvailabilityZonesMember `locationName:"availabilityZones"`

	Listeners *CreateListenersMember `locationName:"listeners" required:"true"`

	LoadBalancerName *string `locationName:"loadBalancerName" type:"string" required:"true"`

	Scheme *string `locationName:"scheme" type:"string"`

	SecurityGroups *SecurityGroupsMember `locationName:"securityGroups"`

	Subnets *SubnetsMember `locationName:"subnets"`

	Tags *TagsMember `locationName:"tags" min:"1"`
}

// CreateLoadBalancerListenersInput ...
type CreateLoadBalancerListenersInput struct {
	_ struct{} `type:"structure"`

	Listeners *CreateListenersMember `locationName:"listeners"`

	LoadBalancerName *string `locationName:"loadBalancerName" type:"string" required:"true"`
}

// CreateLoadBalancerListenersOutput ...
type CreateLoadBalancerListenersOutput struct {
	_ struct{} `type:"structure"`
}

// CreateLoadBalancerOutput ...
type CreateLoadBalancerOutput struct {
	_ struct{} `type:"structure"`

	DNSName *string `locationName:"dnsName" type:"string"`
}

// Listener ...
type Listener struct {
	_ struct{} `type:"structure"`

	InstancePort *int64 `locationName:"instancePort" min:"1" type:"integer" required:"true"`

	InstanceProtocol *string `locationName:"instanceProtocol" type:"string"`

	LoadBalancerPort *int64 `locationName:"loadBalancerPort" type:"integer" required:"true"`

	Protocol *string `locationName:"protocol" type:"string" required:"true"`

	SSLCertificateID *string `locationName:"sslCertificateId" type:"string"`
}

// LoadBalancerNamesMember ...
type LoadBalancerNamesMember struct {
	Member []*string `locationName:"member" locationNameList:"member" type:"list"`
}

// DescribeLoadBalancersInput ...
type DescribeLoadBalancersInput struct {
	_ struct{} `type:"structure"`

	LoadBalancerNames *LoadBalancerNamesMember `locationName:"loadBalancerNames"`

	Marker *string `locationName:"marker" type:"string"`

	PageSize *int64 `locationName:"pageSize" min:"1" type:"integer"`
}

// DescribeLoadBalancersResult ...
type DescribeLoadBalancersResult struct {
	_ struct{} `type:"structure"`

	LoadBalancerDescriptions []*LoadBalancerDescription `locationName:"loadBalancerDescriptions" locationNameList:"member" type:"list"`
}

// DescribeLoadBalancersOutput ...
type DescribeLoadBalancersOutput struct {
	_ struct{} `type:"structure"`

	DescribeLoadBalancersResult *DescribeLoadBalancersResult `locationName:"describeLoadBalancersResult" type:"structre"`

	ResponseMetadata *RequestID `locationName:"responseMetadata" type:"structre"`
}

// RequestID ...
type RequestID struct {
	RequestID *string `locationName:"requestId" type:"string"`
}

// LoadBalancerDescription ...
type LoadBalancerDescription struct {
	_ struct{} `type:"structure"`

	AvailabilityZones []*string `locationName:"availabilityZones" locationNameList:"member" type:"list"`

	BackendServerDescriptions []*BackendServerDescription `locationName:"backendServerDescriptions" locationNameList:"member" type:"list"`

	CanonicalHostedZoneName *string `locationName:"canonicalHostedZoneName" type:"string"`

	CanonicalHostedZoneNameID *string `locationName:"canonicalHostedZoneNameID" type:"string"`

	CreatedTime *time.Time `locationName:"createdTime" type:"timestamp" timestampFormat:"iso8601"`

	DNSName *string `locationName:"dnsName" type:"string"`

	HealthCheck *HealthCheck `locationName:"healthCheck" type:"structure"`

	Instances []*Instance `locationName:"instances" locationNameList:"member" type:"list"`

	ListenerDescriptions []*ListenerDescription `locationName:"listenerDescriptions" locationNameList:"member" type:"list"`

	LoadBalancerName *string `locationName:"loadBalancerName" type:"string"`

	Policies *Policies `locationName:"policies" type:"structure"`

	Scheme *string `locationName:"scheme" type:"string"`

	SecurityGroups []*string `locationName:"securityGroups" locationNameList:"member" type:"list"`

	SourceSecurityGroup *SourceSecurityGroup `locationName:"sourceSecurityGroup" type:"structure"`

	Subnets []*string `locationName:"subnets" locationNameList:"member" type:"list"`

	VPCId *string `locationName:"vpcId" type:"string"`
}

// BackendServerDescription ...
type BackendServerDescription struct {
	_ struct{} `type:"structure"`

	InstancePort *int64 `locationName:"instancePort" min:"1" type:"integer"`

	PolicyNames []*string `locationName:"policyNames" locationNameList:"member" type:"list"`
}

// HealthCheck ...
type HealthCheck struct {
	_ struct{} `type:"structure"`

	HealthyThreshold *int64 `locationName:"healthyThreshold" min:"2" type:"integer" required:"true"`

	Interval *int64 `locationName:"interval" min:"5" type:"integer" required:"true"`

	Target *string `locationName:"target" type:"string" required:"true"`

	Timeout *int64 `locationName:"timeout" min:"2" type:"integer" required:"true"`

	UnhealthyThreshold *int64 `locationName:"unhealthyThreshold" min:"2" type:"integer" required:"true"`
}

// Instance ...
type Instance struct {
	_ struct{} `type:"structure"`

	InstanceID *string `locationName:"instanceId" type:"string"`
}

// ListenerDescription ...
type ListenerDescription struct {
	_ struct{} `type:"structure"`

	Listener *Listener `locationName:"listener" type:"structure"`

	PolicyNames []*string `locationName:"policyNames" locationNameList:"member" type:"list"`
}

// SourceSecurityGroup ...
type SourceSecurityGroup struct {
	_ struct{} `type:"structure"`

	GroupName *string `locationName:"groupName" type:"string"`

	OwnerAlias *string `locationName:"ownerAlias" type:"string"`
}

// Policies ...
type Policies struct {
	_ struct{} `type:"structure"`

	AppCookieStickinessPolicies []*AppCookieStickinessPolicy `locationName:"appCookieStickinessPolicies" locationNameList:"member" type:"list"`

	LBCookieStickinessPolicies []*LBCookieStickinessPolicy `locationName:"lbCookieStickinessPolicies" locationNameList:"member" type:"list"`

	OtherPolicies []*string `locationName:"otherPolicies" locationNameList:"member" type:"list"`
}

// AppCookieStickinessPolicy ...
type AppCookieStickinessPolicy struct {
	_ struct{} `type:"structure"`

	CookieName *string `locationName:"cookieName" type:"string"`

	PolicyName *string `locationName:"policyName" type:"string"`
}

// LBCookieStickinessPolicy ...
type LBCookieStickinessPolicy struct {
	_ struct{} `type:"structure"`

	CookieExpirationPeriod *int64 `locationName:"cookieExpirationPeriod" type:"long"`

	PolicyName *string `locationName:"policyName" type:"string"`
}

// DescribeLoadBalancerAttributesInput ...
type DescribeLoadBalancerAttributesInput struct {
	_ struct{} `type:"structure"`

	LoadBalancerName *string `locationName:"loadBalancerName" type:"string" required:"true"`
}

// DescribeLoadBalancerAttributesOutput ...
type DescribeLoadBalancerAttributesOutput struct {
	_ struct{} `type:"structure"`

	LoadBalancerAttributes *LoadBalancerAttributes `locationName:"loadBalancerAttributes" type:"structure"`

	RequestID *string `locationName:"requestId" type:"string"`
}

// LoadBalancerAttributes ...
type LoadBalancerAttributes struct {
	_ struct{} `type:"structure"`

	AccessLog *AccessLog `locationName:"accessLog" type:"structure"`

	AdditionalAttributes []*AdditionalAttribute `locationName:"additionalAttributes" locationNameList:"member" type:"list"`

	ConnectionDraining *ConnectionDraining `locationName:"connectionDraining" type:"structure"`

	ConnectionSettings *ConnectionSettings `locationName:"connectionSettings" type:"structure"`

	CrossZoneLoadBalancing *CrossZoneLoadBalancing `locationName:"crossZoneLoadBalancing" type:"structure"`
}

// AdditionalAttribute ...
type AdditionalAttribute struct {
	_ struct{} `type:"structure"`

	Key *string `locationName:"key" type:"string"`

	Value *string `locationName:"value" type:"string"`
}

// ConnectionDraining ...
type ConnectionDraining struct {
	_ struct{} `type:"structure"`

	Enabled *bool `locationName:"enabled" type:"boolean" required:"true"`

	Timeout *int64 `locationName:"timeout" type:"integer"`
}

// ConnectionSettings ...
type ConnectionSettings struct {
	_ struct{} `type:"structure"`

	IdleTimeout *int64 `locationName:"idleTimeout" min:"1" type:"integer" required:"true"`
}

// CrossZoneLoadBalancing ...
type CrossZoneLoadBalancing struct {
	_ struct{} `type:"structure"`

	Enabled *bool `locationName:"enabled" type:"boolean" required:"true"`
}

// AccessLog ...
type AccessLog struct {
	_ struct{} `type:"structure"`

	EmitInterval *int64 `locationName:"emitInterval" type:"integer"`

	Enabled *bool `locationName:"enabled" type:"boolean" required:"true"`

	S3BucketName *string `locationName:"s3BucketName" type:"string"`

	S3BucketPrefix *string `locationName:"s3BucketPrefix" type:"string"`
}

// DeleteLoadBalancerListenersInput ...
type DeleteLoadBalancerListenersInput struct {
	_ struct{} `type:"structure"`

	LoadBalancerName *string `locationName:"loadBalancerName" type:"string" required:"true"`

	LoadBalancerPorts []*int64 `locationName:"loadBalancerPorts" locationNameList:"member" type:"list" required:"true"`
}

// DeleteLoadBalancerListenersOutput ...
type DeleteLoadBalancerListenersOutput struct {
	_ struct{} `type:"structure"`
}

// ConfigureHealthCheckInput ...
type ConfigureHealthCheckInput struct {
	_ struct{} `type:"structure"`

	HealthCheck *HealthCheck `locationName:"healthCheck" type:"structure" required:"true"`

	LoadBalancerName *string `locationName:"loadBalancerName" type:"string" required:"true"`
}

// ConfigureHealthCheckOutput ...
type ConfigureHealthCheckOutput struct {
	_ struct{} `type:"structure"`

	HealthCheck *HealthCheck `locationName:"healthCheck" type:"structure"`
}

// ApplySecurityGroupsToLoadBalancerInput ...
type ApplySecurityGroupsToLoadBalancerInput struct {
	_ struct{} `type:"structure"`

	LoadBalancerName *string `locationName:"loadBalancerName" type:"string" required:"true"`

	SecurityGroups []*string `locationName:"securityGroups" locationNameList:"member" type:"list" required:"true"`
}

// ApplySecurityGroupsToLoadBalancerOutput ...
type ApplySecurityGroupsToLoadBalancerOutput struct {
	_ struct{} `type:"structure"`

	SecurityGroups []*string `locationName:"securityGroups" locationNameList:"member" type:"list"`
}

// EnableAvailabilityZonesForLoadBalancerInput ...
type EnableAvailabilityZonesForLoadBalancerInput struct {
	_ struct{} `type:"structure"`

	AvailabilityZones []*string `locationName:"availabilityZones" locationNameList:"member" type:"list" required:"true"`

	LoadBalancerName *string `locationName:"loadBalancerName" type:"string" required:"true"`
}

// EnableAvailabilityZonesForLoadBalancerOutput ...
type EnableAvailabilityZonesForLoadBalancerOutput struct {
	_ struct{} `type:"structure"`

	AvailabilityZones []*string `locationName:"availabilityZones" locationNameList:"member" type:"list"`
}

// DisableAvailabilityZonesForLoadBalancerInput ...
type DisableAvailabilityZonesForLoadBalancerInput struct {
	_ struct{} `type:"structure"`

	AvailabilityZones []*string `locationName:"availabilityZones" locationNameList:"member" type:"list" required:"true"`

	LoadBalancerName *string `locationName:"loadBalancerName" type:"string" required:"true"`
}

// DisableAvailabilityZonesForLoadBalancerOutput ...
type DisableAvailabilityZonesForLoadBalancerOutput struct {
	_ struct{} `type:"structure"`

	AvailabilityZones []*string `locationName:"availabilityZones" locationNameList:"member" type:"list"`
}

// AttachLoadBalancerToSubnetsInput ...
type AttachLoadBalancerToSubnetsInput struct {
	_ struct{} `type:"structure"`

	LoadBalancerName *string `locationName:"loadBalancerName" type:"string" required:"true"`

	Subnets []*string `locationName:"subnets" locationNameList:"member" type:"list" required:"true"`
}

// AttachLoadBalancerToSubnetsOutput ...
type AttachLoadBalancerToSubnetsOutput struct {
	_ struct{} `type:"structure"`

	Subnets []*string `locationName:"subnets" locationNameList:"member" type:"list"`
}

// DeleteLoadBalancerInput ...
type DeleteLoadBalancerInput struct {
	_ struct{} `type:"structure"`

	LoadBalancerName *string `locationName:"loadBalancerName" type:"string" required:"true"`
}

// DeleteLoadBalancerOutput ...
type DeleteLoadBalancerOutput struct {
	_ struct{} `type:"structure"`
}

// RegisterInstancesWithLoadBalancerInput ...
type RegisterInstancesWithLoadBalancerInput struct {
	_ struct{} `type:"structure"`

	Instances []*Instance `locationName:"instances" locationNameList:"member" type:"list" required:"true"`

	LoadBalancerName *string `locationName:"loadBalancerName" type:"string" required:"true"`
}

// RegisterInstancesWithLoadBalancerOutput ...
type RegisterInstancesWithLoadBalancerOutput struct {
	_ struct{} `type:"structure"`

	Instances []*Instance `locationName:"instances" locationNameList:"member" type:"list"`
}

// DeregisterInstancesFromLoadBalancerInput ...
type DeregisterInstancesFromLoadBalancerInput struct {
	_ struct{} `type:"structure"`

	Instances []*Instance `locationName:"instances" locationNameList:"member" type:"list" required:"true"`

	LoadBalancerName *string `locationName:"loadBalancerName" type:"string" required:"true"`
}

// DeregisterInstancesFromLoadBalancerOutput ...
type DeregisterInstancesFromLoadBalancerOutput struct {
	_ struct{} `type:"structure"`

	Instances []*Instance `locationName:"instances" locationNameList:"member" type:"list"`
}

// DetachLoadBalancerFromSubnetsInput ...
type DetachLoadBalancerFromSubnetsInput struct {
	_ struct{} `type:"structure"`

	LoadBalancerName *string `locationName:"loadBalancerName" type:"string" required:"true"`

	Subnets []*string `locationName:"subnets" locationNameList:"member" type:"list" required:"true"`
}

// DetachLoadBalancerFromSubnetsOutput ...
type DetachLoadBalancerFromSubnetsOutput struct {
	_ struct{} `type:"structure"`

	Subnets []*string `locationName:"subnets" locationNameList:"member" type:"list"`
}
