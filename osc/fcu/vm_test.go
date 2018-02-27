package fcu

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
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

func TestVM_RunInstance(t *testing.T) {
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
			<?xml version="1.0" encoding="UTF-8"?>
<RunInstancesResponse xmlns="http://ec2.amazonaws.com/doc/2014-06-15/"><requestId>b07d2bff-536a-4bc2-b528-da9d008638e0</requestId><reservationId>r-b4f82c1c</reservationId><ownerId>520679080430</ownerId><groupSet><item><groupId>sg-1385300f</groupId><groupName>default</groupName></item></groupSet><instancesSet><item><instanceId>i-0e8ea0a2</instanceId><imageId>ami-8a6a0120</imageId><instanceState><code>0</code><name>pending</name></instanceState><privateDnsName>ip-10-0-1-155.eu-west-2.compute.internal</privateDnsName><dnsName></dnsName><keyName>terraform-basic</keyName><amiLaunchIndex>0</amiLaunchIndex><productCodes/><instanceType>t2.micro</instanceType><launchTime>2018-02-22T20:48:32.524Z</launchTime><placement><availabilityZone>eu-west-2a</availabilityZone><groupName></groupName><tenancy>default</tenancy></placement><kernelId></kernelId><monitoring><state>disabled</state></monitoring><subnetId>subnet-861fbecc</subnetId><vpcId>vpc-e9d09d63</vpcId><privateIpAddress>10.0.1.155</privateIpAddress><sourceDestCheck>true</sourceDestCheck><groupSet><item><groupId>sg-1385300f</groupId><groupName>default</groupName></item></groupSet><architecture>x86_64</architecture><rootDeviceType>ebs</rootDeviceType><rootDeviceName>/dev/sda1</rootDeviceName><blockDeviceMapping><item><deviceName>/dev/sda1</deviceName><ebs><volumeId>vol-9454b3cc</volumeId><status>attaching</status><attachTime>2018-02-22T20:48:32.524Z</attachTime><deleteOnTermination>true</deleteOnTermination></ebs></item></blockDeviceMapping><virtualizationType>hvm</virtualizationType><clientToken></clientToken><hypervisor>xen</hypervisor><networkInterfaceSet><item><networkInterfaceId>eni-33a7d022</networkInterfaceId><subnetId>subnet-861fbecc</subnetId><vpcId>vpc-e9d09d63</vpcId><description>Primary network interface</description><ownerId>520679080430</ownerId><status>in-use</status><macAddress>aa:7f:a8:aa:94:33</macAddress><privateIpAddress>10.0.1.155</privateIpAddress><privateDnsName>ip-10-0-1-155.eu-west-2.compute.internal</privateDnsName><sourceDestCheck>true</sourceDestCheck><groupSet><item><groupId>sg-1385300f</groupId><groupName>default</groupName></item></groupSet><attachment><attachmentId>eni-attach-e23c25bf</attachmentId><deviceIndex>0</deviceIndex><status>attached</status><attachTime>2018-02-22T20:48:32.524Z</attachTime><deleteOnTermination>true</deleteOnTermination></attachment><privateIpAddressesSet><item><privateIpAddress>10.0.1.155</privateIpAddress><privateDnsName>ip-10-0-1-155.eu-west-2.compute.internal</privateDnsName><primary>true</primary></item></privateIpAddressesSet></item></networkInterfaceSet><ebsOptimized>false</ebsOptimized></item></instancesSet></RunInstancesResponse>
      `)
	})

	server, err := client.VM.RunInstance(input)
	if err != nil {
		t.Errorf("VM.RunInstance returned error: %v", err)
	}

	instanceID := *server.Instances[0].InstanceId
	expectedID := "i-0e8ea0a2"

	if instanceID != expectedID {
		t.Fatalf("Expected InstanceID:(%s), Got(%s)", instanceID, expectedID)
	}

}

func TestDescribe_Instance(t *testing.T) {
	setup()
	defer teardown()

	instanceID := "i-d470ce8f"

	input := DescribeInstancesInput{
		InstanceIds: []*string{&instanceID},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
		<DescribeInstancesResponse xmlns="http://ec2.amazonaws.com/doc/2014-06-15/"><requestId>16d09939-f72c-471c-99ae-6704d888197c</requestId><reservationSet><item><reservationId>r-3cedd89e</reservationId><ownerId>520679080430</ownerId><groupSet><item><groupId>sg-6ed31f3e</groupId><groupName>default</groupName></item></groupSet><instancesSet><item><instanceId>i-bee7ebf3</instanceId><imageId>ami-8a6a0120</imageId><instanceState><code>0</code><name>pending</name></instanceState><privateDnsName>ip-10-9-8-166.eu-west-2.compute.internal</privateDnsName><dnsName></dnsName><keyName></keyName><amiLaunchIndex>0</amiLaunchIndex><productCodes/><instanceType>m1.small</instanceType><launchTime>2018-02-08T01:46:45.269Z</launchTime><placement><availabilityZone>eu-west-2a</availabilityZone><groupName></groupName><tenancy>default</tenancy></placement><kernelId></kernelId><monitoring><state>disabled</state></monitoring><privateIpAddress>10.9.8.166</privateIpAddress><groupSet><item><groupId>sg-6ed31f3e</groupId><groupName>default</groupName></item></groupSet><architecture>x86_64</architecture><rootDeviceType>ebs</rootDeviceType><rootDeviceName>/dev/sda1</rootDeviceName><blockDeviceMapping><item><deviceName>/dev/sda1</deviceName><ebs><volumeId>vol-f65f0614</volumeId><status>attaching</status><attachTime>2018-02-08T01:46:45.269Z</attachTime><deleteOnTermination>true</deleteOnTermination></ebs></item></blockDeviceMapping><virtualizationType>hvm</virtualizationType><clientToken></clientToken><hypervisor>xen</hypervisor><networkInterfaceSet/><ebsOptimized>false</ebsOptimized></item></instancesSet></item></reservationSet></DescribeInstancesResponse>
      `)
	})

	desc, err := client.VM.DescribeInstances(&input)
	if err != nil {
		t.Errorf("VM.RunInstance returned error: %v", err)
	}

	expectedID := "i-bee7ebf3"
	outputInstanceID := *desc.Reservations[0].Instances[0].InstanceId

	if outputInstanceID != expectedID {
		t.Fatalf("Expected InstanceID:(%s), Got(%s)", outputInstanceID, expectedID)
	}

}

func TestVM_GetPasswordData(t *testing.T) {
	setup()
	defer teardown()

	instanceID := "i-9c1b9711"

	input := GetPasswordDataInput{
		InstanceId: &instanceID,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
		<GetPasswordDataResponse xmlns="http://ec2.amazonaws.com/doc/2014-06-15/"><requestId>ce00bd1d-9f3f-4bfd-be6b-1b7b73454c20</requestId><instanceId>i-9c1b9711</instanceId><timestamp>2018-02-08T02:46:15.789Z</timestamp><passwordData></passwordData></GetPasswordDataResponse>
		`)
	})

	term, err := client.VM.GetPasswordData(&input)
	if err != nil {
		t.Errorf("VM.GetPasswordData returned error: %v", err)
	}

	expectedID := "i-9c1b9711"
	outputInstanceID := *term.InstanceId

	if outputInstanceID != expectedID {
		t.Fatalf("Expected InstanceID:(%s), Got(%s)", outputInstanceID, expectedID)
	}
}

func TestVM_ModifyInstanceKeyPair(t *testing.T) {
	setup()
	defer teardown()

	instanceID := "i-484e76e2"
	keypair := "testkey"

	input := ModifyInstanceKeyPairInput{
		InstanceId: &instanceID,
		KeyName:    &keypair,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, ``)
	})

	err := client.VM.ModifyInstanceKeyPair(&input)
	if err != nil {
		t.Errorf("VM.ModifyInstanceKeyPair returned error: %v", err)
	}
}

func TestVM_TerminateInstances(t *testing.T) {
	setup()
	defer teardown()

	instanceID := "i-484e76e2"

	input := TerminateInstancesInput{
		InstanceIds: []*string{&instanceID},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
		<TerminateInstancesResponse xmlns="http://ec2.amazonaws.com/doc/2014-06-15/"><requestId>f508de7e-fe4b-4572-a977-e74efb9f3b76</requestId><instancesSet><item><instanceId>i-484e76e2</instanceId><currentState><code>32</code><name>shutting-down</name></currentState><previousState><code>0</code><name>pending</name></previousState></item></instancesSet></TerminateInstancesResponse>
		`)
	})

	term, err := client.VM.TerminateInstances(&input)
	if err != nil {
		t.Errorf("VM.RunInstance returned error: %v", err)
	}

	expectedID := "i-484e76e2"
	outputInstanceID := *term.TerminatingInstances[0].InstanceId

	if outputInstanceID != expectedID {
		t.Fatalf("Expected InstanceID:(%s), Got(%s)", outputInstanceID, expectedID)
	}
}

func TestVM_ModifyInstanceAttribute(t *testing.T) {
	setup()
	defer teardown()

	instanceID := "i-d742ed97"

	input := ModifyInstanceAttributeInput{
		InstanceId: aws.String(instanceID),
		DisableApiTermination: &AttributeBooleanValue{
			Value: aws.Bool(false),
		},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
		<ModifyInstanceAttributeResponse
	xmlns="http://ec2.amazonaws.com/doc/2014-06-15/">
	<requestId>f508de7e-fe4b-4572-a977-e74efb9f3b76</requestId>
	<return>true</return>
</ModifyInstanceAttributeResponse>
		`)
	})

	_, err := client.VM.ModifyInstanceAttribute(&input)
	if err != nil {
		t.Errorf("VM.ModifyInstanceAttribute returned error: %v", err)
	}
}

func TestVM_StopInstances(t *testing.T) {
	setup()
	defer teardown()

	instanceID := "i-d742ed97"

	input := StopInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
		<StoppingInstancesResponse
	xmlns="http://ec2.amazonaws.com/doc/2014-06-15/">
	<requestId>f508de7e-fe4b-4572-a977-e74efb9f3b76</requestId>
	<stoppingInstances>
		<item>
			<instanceId>i-d742ed97</instanceId>
			<currentState>
				<code>64</code>
				<name>stopping</name>
			</currentState>
			<previousState>
				<code>16</code>
				<name>running</name>
			</previousState>
		</item>
	</stoppingInstances>
</StoppingInstancesResponse>
		`)
	})

	_, err := client.VM.StopInstances(&input)
	if err != nil {
		t.Errorf("VM.StopInstances returned error: %v", err)
	}
}

func TestVM_StartInstances(t *testing.T) {
	setup()
	defer teardown()

	instanceID := "i-d742ed97"

	input := StartInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
		<StartingInstancesResponse
	xmlns="http://ec2.amazonaws.com/doc/2014-06-15/">
	<requestId>f508de7e-fe4b-4572-a977-e74efb9f3b76</requestId>
	<startingInstances>
		<item>
			<instanceId>i-d742ed97</instanceId>
			<currentState>
				<code>0</code>
				<name>pending</name>
			</currentState>
			<previousState>
				<code>80</code>
				<name>pending</name>
			</previousState>
		</item>
	</startingInstances>
</StartingInstancesResponse>
		`)
	})

	_, err := client.VM.StartInstances(&input)
	if err != nil {
		t.Errorf("VM.StartInstances returned error: %v", err)
	}
}

func TestVM_GetOwnerId(t *testing.T) {
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

	ownerID := *server.OwnerId
	expectedOwnerID := "520679080430"

	if ownerID != expectedOwnerID {
		t.Fatalf("Expected OwnerID:(%s), Got(%s)", ownerID, expectedOwnerID)
	}
}

func TestVM_GetRequesterID(t *testing.T) {
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

	requesterID := *server.RequesterId
	expectedrequesterID := "193ddebf-63d4-466d-9fe1-d5b74b9962f0"

	fmt.Println(requesterID, expectedrequesterID)

	if requesterID != expectedrequesterID {
		t.Fatalf("Expected OwnerID:(%s), Got(%s)", requesterID, expectedrequesterID)
	}
}
func TestVM_GetReservationID(t *testing.T) {
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

	reservationID := *server.ReservationId
	expectedReservationID := "r-071eb05d"

	fmt.Println("OUTPUT =>", reservationID, expectedReservationID)

	if reservationID != expectedReservationID {
		t.Fatalf("Expected OwnerID:(%s), Got(%s)", reservationID, expectedReservationID)
	}
}
