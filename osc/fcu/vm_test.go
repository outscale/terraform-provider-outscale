package fcu

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/terraform-providers/terraform-provider-outscale/osc"
)

var (
	mux *http.ServeMux

	ctx = context.TODO()

	client *Client

	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	config := osc.Config{
		Credentials: &osc.Credentials{
			AccessKey: "AKID",
			SecretKey: "SecretKey",
			Region:    "region",
		},
	}

	client, _ = NewFCUClient(config)

	u, _ := url.Parse(server.URL)
	client.client.Config.BaseURL = u

}

func teardown() {
	server.Close()
}

func TestServers_Create(t *testing.T) {
	setup()
	defer teardown()

	var maxC int64
	imageID := "ami-8a6a0120"
	maxC = 1

	input := &RunInstancesInput{
		ImageId:  &imageID,
		MaxCount: &maxC,
		MinCount: &maxC,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
      <RunInstancesResponse xmlns="http://ec2.amazonaws.com/doc/2014-06-15/"><requestId>193ddebf-63d4-466d-9fe1-d5b74b9962f0</requestId><reservationId>r-071eb05d</reservationId><ownerId>520679080430</ownerId><groupSet><item><groupId>sg-6ed31f3e</groupId><groupName>default</groupName></item></groupSet><instancesSet><item><instanceId>i-d470ce8f</instanceId><imageId>ami-8a6a0120</imageId><instanceState><code>0</code><name>pending</name></instanceState><privateDnsName>ip-10-9-10-212.eu-west-2.compute.internal</privateDnsName><dnsName></dnsName><keyName></keyName><amiLaunchIndex>0</amiLaunchIndex><productCodes/><instanceType>m1.small</instanceType><launchTime>2018-02-08T00:51:38.866Z</launchTime><placement><availabilityZone>eu-west-2a</availabilityZone><groupName></groupName><tenancy>default</tenancy></placement><kernelId></kernelId><monitoring><state>disabled</state></monitoring><privateIpAddress>10.9.10.212</privateIpAddress><groupSet><item><groupId>sg-6ed31f3e</groupId><groupName>default</groupName></item></groupSet><architecture>x86_64</architecture><rootDeviceType>ebs</rootDeviceType><rootDeviceName>/dev/sda1</rootDeviceName><blockDeviceMapping><item><deviceName>/dev/sda1</deviceName><ebs><volumeId>vol-ee2f2a14</volumeId><status>attaching</status><attachTime>2018-02-08T00:51:38.866Z</attachTime><deleteOnTermination>true</deleteOnTermination></ebs></item></blockDeviceMapping><virtualizationType>hvm</virtualizationType><clientToken></clientToken><hypervisor>xen</hypervisor><networkInterfaceSet/><ebsOptimized>false</ebsOptimized></item></instancesSet></RunInstancesResponse>
      `)
	})

	server, err := client.VM.RunInstance(input)
	if err != nil {
		t.Errorf("VM.RunInstance returned error: %v", err)
	}

	instanceID := *server.Instances[0].InstanceId
	expectedID := "i-d470ce8f"

	if instanceID != expectedID {
		t.Fatalf("Expected InstanceID:(%s), Got(%s)", instanceID, expectedID)
	}

}
