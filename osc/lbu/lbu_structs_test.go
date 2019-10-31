package lbu

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/terraform-providers/terraform-provider-outscale/osc/handler"
)

func TestUnmarshalXMLDescribeLoadBalancers(t *testing.T) {
	buf := bytes.NewReader([]byte(`<?xml version="1.0" encoding="UTF-8"?>
		<DescribeLoadBalancersResponse xmlns="http://elasticloadbalancing.amazonaws.com/doc/2012-06-01/"><DescribeLoadBalancersResult><LoadBalancerDescriptions><member><LoadBalancerName>tf-test-lb-0i8z4</LoadBalancerName><DNSName>tf-test-lb-0i8z4-614518246.eu-west-2.lbu.outscale.com</DNSName><ListenerDescriptions><member><Listener><Protocol>HTTP</Protocol><LoadBalancerPort>80</LoadBalancerPort><InstanceProtocol>HTTP</InstanceProtocol><InstancePort>8000</InstancePort></Listener><PolicyNames/></member></ListenerDescriptions><Policies><AppCookieStickinessPolicies><member><PolicyName>foo-policy</PolicyName><CookieName>MyAppCookie</CookieName></member></AppCookieStickinessPolicies><LBCookieStickinessPolicies/><OtherPolicies/></Policies><AvailabilityZones><member>eu-west-2a</member></AvailabilityZones><Instances/><HealthCheck><Target>TCP:8000</Target><Interval>30</Interval><Timeout>5</Timeout><UnhealthyThreshold>2</UnhealthyThreshold><HealthyThreshold>10</HealthyThreshold></HealthCheck><SourceSecurityGroup><OwnerAlias>outscale-elb</OwnerAlias><GroupName>outscale-elb-sg</GroupName></SourceSecurityGroup><CreatedTime>2018-06-22T20:36:52.106Z</CreatedTime><Scheme>internet-facing</Scheme></member></LoadBalancerDescriptions></DescribeLoadBalancersResult><ResponseMetadata><RequestId>92da855e-5073-4fd9-944f-da370a64d814</RequestId></ResponseMetadata></DescribeLoadBalancersResponse>`))
	res := &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	v := DescribeLoadBalancersOutput{}

	if err := handler.UnmarshalXML(&v, res, "Action=CreateAppCookieStickinessPolicy&CookieName=MyOtherAppCookie&LoadBalancerName=tf-test-lb-tqi13&PolicyName=foo-policy&Version=2018-05-14"); err != nil {
		t.Fatalf("err: %s", err)
	}

	if v.ResponseMetadata == nil {
		t.Fatalf("Cannot unmarshal with the struct %+v", v)

		if v.ResponseMetadata.RequestID == nil {
			t.Fatalf("Cannot unmarshal ResponseMetadata %+v", v.ResponseMetadata)
		}
	}

	log.Printf("[DEBUG] %s\n", *v.ResponseMetadata.RequestID)

}
