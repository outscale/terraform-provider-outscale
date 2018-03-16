package fcu

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
)

func TestVM_CreateInternetGateaway(t *testing.T) {
	setup()
	defer teardown()

	input := CreateInternetGatewayInput{}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<CreateInternetGatewayResponse xmlns="http://ec2.amazonaws.com/doc/2014-06-15/"><requestId>6bc7d085-32c3-4b87-b07b-04ccc5c25f19</requestId><internetGateway><internetGatewayId>igw-349c9be7</internetGatewayId><attachmentSet/><tagSet/></internetGateway></CreateInternetGatewayResponse>`)
	})

	desc, err := client.VM.CreateInternetGateway(&input)
	if err != nil {
		t.Errorf("VM.CreateInternetGateway returned error: %v", err)
	}

	expectedID := "igw-349c9be7"
	outputInstanceID := *desc.InternetGateway.InternetGatewayId

	if outputInstanceID != expectedID {
		t.Fatalf("Expected InternetGatewayID:(%s), Got(%s)", outputInstanceID, expectedID)
	}

}

func TestVM_DescribeInternetGateaways(t *testing.T) {
	setup()
	defer teardown()

	expectedID := "igw-251475c9"

	input := DescribeInternetGatewaysInput{
		InternetGatewayIds: []*string{&expectedID},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<DescribeInternetGatewaysResponse xmlns="http://ec2.amazonaws.com/doc/2014-06-15/"><requestId>89135ead-4cab-4e11-b842-7b70d0d5f583</requestId><internetGatewaySet><item><internetGatewayId>igw-251475c9</internetGatewayId><attachmentSet><item><vpcId>vpc-e9d09d63</vpcId><state>available</state></item></attachmentSet><tagSet><item><key>Name</key><value>Default</value></item></tagSet></item></internetGatewaySet></DescribeInternetGatewaysResponse>`)
	})

	desc, err := client.VM.DescribeInternetGateways(&input)
	if err != nil {
		t.Errorf("VM.DescribeInternetGateways returned error: %v", err)
	}

	outputInstanceID := *desc.InternetGateways[0].InternetGatewayId

	if outputInstanceID != expectedID {
		t.Fatalf("Expected InstanceID:(%s), Got(%s)", outputInstanceID, expectedID)
	}

}

func TestVM_DeleteInternetGateway(t *testing.T) {
	setup()
	defer teardown()

	expectedID := "igw-349c9be7"

	input := DeleteInternetGatewayInput{
		InternetGatewayId: &expectedID,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, ``)
	})

	_, err := client.VM.DeleteInternetGateway(&input)
	if err != nil {
		t.Errorf("VM.DeleteInternetGateway returned error: %v", err)
	}

}

func TestVM_CreateVpc(t *testing.T) {
	setup()
	defer teardown()

	expectedID := "vpc-53769ad9"

	input := CreateVpcInput{
		CidrBlock: aws.String("10.0.0.0/16"),
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<CreateVpcResponse xmlns="http://ec2.amazonaws.com/doc/2014-06-15/"><requestId>52faf8ea-65ea-46cb-896c-df33cd06e8fc</requestId><vpc><vpcId>vpc-53769ad9</vpcId><state>available</state><cidrBlock>10.0.0.0/16</cidrBlock><dhcpOptionsId>dopt-1ea5389e</dhcpOptionsId><tagSet/><instanceTenancy>default</instanceTenancy><isDefault>false</isDefault></vpc></CreateVpcResponse>`)
	})

	desc, err := client.VM.CreateVpc(&input)
	if err != nil {
		t.Errorf("VM.CreateVpc returned error: %v", err)
	}

	outputVpcID := *desc.Vpc.VpcId

	if outputVpcID != expectedID {
		t.Fatalf("Expected VpcId:(%s), Got(%s)", outputVpcID, expectedID)
	}
	expectedState := "available"
	state := *desc.Vpc.State
	if expectedState != state {
		t.Fatalf("Expected state:(%s), Got(%s)", state, expectedState)
	}

	expectedCIDR := "10.0.0.0/16"
	cidr := *desc.Vpc.CidrBlock
	if expectedCIDR != expectedCIDR {
		t.Fatalf("Expected cidr:(%s), Got(%s)", cidr, expectedCIDR)
	}

}

func TestVM_DescribeVpcs(t *testing.T) {
	setup()
	defer teardown()

	expectedID := "vpc-53769ad9"

	input := DescribeVpcsInput{
		VpcIds: []*string{&expectedID},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<DescribeVpcsResponse xmlns="http://ec2.amazonaws.com/doc/2014-06-15/"><requestId>1ab37b1d-67fb-4edb-aeea-1cdb02e1c232</requestId><vpcSet><item><vpcId>vpc-53769ad9</vpcId><state>available</state><cidrBlock>10.0.0.0/16</cidrBlock><dhcpOptionsId>dopt-1ea5389e</dhcpOptionsId><tagSet/><instanceTenancy>default</instanceTenancy><isDefault>false</isDefault></item></vpcSet></DescribeVpcsResponse>`)
	})

	desc, err := client.VM.DescribeVpcs(&input)
	if err != nil {
		t.Errorf("VM.DescribeVpcs returned error: %v", err)
	}

	outputVpcID := *desc.Vpcs[0].VpcId

	if outputVpcID != expectedID {
		t.Fatalf("Expected VPCID:(%s), Got(%s)", outputVpcID, expectedID)
	}
	expectedState := "available"
	state := *desc.Vpcs[0].State
	if expectedState != state {
		t.Fatalf("Expected state:(%s), Got(%s)", state, expectedState)
	}

	expectedCIDR := "10.0.0.0/16"
	cidr := *desc.Vpcs[0].CidrBlock
	if expectedCIDR != cidr {
		t.Fatalf("Expected cidr:(%s), Got(%s)", cidr, expectedCIDR)
	}

}

func TestVM_DeleteVpc(t *testing.T) {
	setup()
	defer teardown()

	expectedID := "vpc-53769ad9"

	input := DeleteVpcInput{
		VpcId: &expectedID,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, ``)
	})

	_, err := client.VM.DeleteVpc(&input)
	if err != nil {
		t.Errorf("VM.DeleteVpc returned error: %v", err)
	}

}
