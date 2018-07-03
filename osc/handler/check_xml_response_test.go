package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestUnmarshalXML(t *testing.T) {
	buf := bytes.NewReader([]byte("<OperationNameResponse><Str>myname</Str><FooNum>123</FooNum><FalseBool>false</FalseBool><TrueBool>true</TrueBool><Float>1.2</Float><Double>1.3</Double><Long>200</Long><Char>a</Char><RequestId>request-id</RequestId></OperationNameResponse>"))
	res := &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}
	var v struct{}
	if err := UnmarshalXML(&v, res, ""); err != nil {
		t.Fatalf("err: %s", err)
	}
}

type LBUResponse struct {
	CreateAppCookieStickinessPolicyResult *CreateAppCookieStickinessPolicyResult `type:"structure"`
	ResponseMetadata                      *ResponseMetadata                      `type:"structure"`
}

type CreateAppCookieStickinessPolicyResult struct {
	_ struct{} `type:"structure"`
}

//ResponseMatadata ...
type ResponseMetadata struct {
	RequestID *string `locationName:"RequestId" type:"string"`
}

func TestUnmarshalXMLRequestMetadata(t *testing.T) {
	buf := bytes.NewReader([]byte(`<CreateAppCookieStickinessPolicyResponse xmlns="http://elasticloadbalancing.amazonaws.com/doc/2012-06-01/"><CreateAppCookieStickinessPolicyResult/><ResponseMetadata><RequestId>fbb49983-45f5-4284-8c8f-52b7f180946a</RequestId></ResponseMetadata></CreateAppCookieStickinessPolicyResponse>`))
	res := &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	v := LBUResponse{}

	if err := UnmarshalXML(&v, res, "Action=CreateAppCookieStickinessPolicy&CookieName=MyOtherAppCookie&LoadBalancerName=tf-test-lb-tqi13&PolicyName=foo-policy&Version=2018-05-14"); err != nil {
		t.Fatalf("err: %s", err)
	}

	if v.ResponseMetadata == nil {
		t.Fatalf("Cannot unmarshal with the struct %+v", v)

		if v.ResponseMetadata.RequestID == nil {
			t.Fatalf("Cannot unmarshal ResponseMetadata %+v", v.ResponseMetadata)
		}
	}

	fmt.Printf("[Debug] %s\n", *v.ResponseMetadata.RequestID)

}
