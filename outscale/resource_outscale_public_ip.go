package outscale

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscalePublicIP() *schema.Resource {
	return &schema.Resource{
		Create: resourcePublicIPCreate,
		Read:   resourcePublicIPRead,
		Update: resourcePublicIPUpdate,
		Delete: resourcePublicIPDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },

		// Timeouts: &schema.ResourceTimeout{
		// 	Create: schema.DefaultTimeout(10 * time.Minute),
		// 	Update: schema.DefaultTimeout(10 * time.Minute),
		// 	Delete: schema.DefaultTimeout(10 * time.Minute),
		// },

		Schema: getPublicIPSchema(),
	}
}

func resourcePublicIPCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	// Build the creation struct
	runOpts := &fcu.AllocateAddressInput{}

	domain, ok := d.GetOk("domain")
	if ok {
		runOpts.Domain = domain.(*string)
	}

	allocResp, err := conn.VM.AllocateAddress(runOpts)
	if err != nil {
		return fmt.Errorf("Error allocating address: %s", err)
	}

	d.Set("domain", allocResp.Domain)
	d.SetId(*allocResp.PublicIp)
	return resourcePublicIPRead(d, meta)
}

func resourcePublicIPRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	domain := resourceOutscaleDomain(d)
	id := d.Id()

	req := &fcu.DescribeAddressesInput{}

	if domain == "vpc" {
		req.AllocationIds = []*string{aws.String(id)}
	} else {
		req.PublicIps = []*string{aws.String(id)}
	}

	log.Printf(
		"[DEBUG] Public IP describe configuration: %s (domain: %s)",
		req, domain)

	describeAddresses, err := conn.VM.DescribeAddressesRequest(req)
	if err != nil {
		if ec2err, ok := err.(awserr.Error); ok && (ec2err.Code() == "InvalidAllocationID.NotFound" || ec2err.Code() == "InvalidAddress.NotFound") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Public IP: %s", err)
	}

	// Verify Outscale returned our PublicIP
	if len(describeAddresses.Addresses) != 1 ||
		domain == "vpc" && *describeAddresses.Addresses[0].AllocationId != id ||
		*describeAddresses.Addresses[0].PublicIp != id {
		if err != nil {
			return fmt.Errorf("Unable to find Public IP: %#v", describeAddresses.Addresses)
		}
	}

	address := describeAddresses.Addresses[0]

	if address.AllocationId != nil {
		d.Set("allocation_id", address.AllocationId)
	}
	if address.AssociationId != nil {
		d.Set("association_id", address.AssociationId)
	}
	if address.Domain != nil {
		d.Set("domain", address.Domain)
	}
	if address.InstanceId != nil {
		d.Set("instance_id", address.InstanceId)
	}
	if address.NetworkInterfaceId != nil {
		d.Set("network_interface_id", address.NetworkInterfaceId)
	}
	if address.NetworkInterfaceOwnerId != nil {
		d.Set("network_interface_owner_id", address.NetworkInterfaceOwnerId)
	}
	if address.PrivateIpAddress != nil {
		d.Set("private_ip_address", address.PrivateIpAddress)
	}
	if address.PublicIp != nil {
		d.Set("public_ip", address.PublicIp)
	}

	// On import (domain never set, which it must've been if we created),
	// set the 'vpc' attribute depending on if we're in a VPC.
	if address.Domain != nil {
		d.Set("vpc", *address.Domain == "vpc")
	}

	d.Set("domain", address.Domain)

	// Force ID to be an Allocation ID if we're on a VPC
	// This allows users to import the PublicIP based on the IP if they are in a VPC
	if *address.Domain == "vpc" && net.ParseIP(id) != nil {
		log.Printf("[DEBUG] Re-assigning Public IP ID (%s) to it's Allocation ID (%s)", d.Id(), *address.AllocationId)
		d.SetId(*address.AllocationId)
	}

	return nil
}

func getPublicIPSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"allocation_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"association_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"domain": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"instance_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"network_interface_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"network_interface_owner_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"private_ip_address": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"public_ip": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func resourceOutscaleDomain(d *schema.ResourceData) string {
	if v, ok := d.GetOk("domain"); ok {
		return v.(string)
	} else if strings.Contains(d.Id(), "allocation") {
		// We have to do this for backwards compatibility since TF 0.1
		// didn't have the "domain" computed attribute.
		return "vpc"
	}

	return "standard"
}

func resourcePublicIPDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePublicIPUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}
