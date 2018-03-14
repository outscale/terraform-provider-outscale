package fcu

import (
	"fmt"
	"net/http"
	"testing"
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
