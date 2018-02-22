package outscale

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscalePublicIP() *schema.Resource {
	return &schema.Resource{
		Create: resourcePublicIPCreate,
		Read:   resourcePublicIPRead,
		Delete: resourcePublicIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: getPublicIPSchema(),
	}
}

func resourcePublicIPCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	domain := resourceOutscaleDomain(d)

	allocOpts := &fcu.AllocateAddressInput{
		Domain: aws.String(domain),
	}

	log.Printf("[DEBUG] EIP create configuration: %#v", allocOpts)
	allocResp, err := conn.VM.AllocateAddress(allocOpts)
	if err != nil {
		return fmt.Errorf("Error creating EIP: %s", err)
	}

	d.Set("domain", allocResp.Domain)

	log.Printf("[DEBUG] EIP Allocate: %#v", allocResp)
	if d.Get("domain").(string) == "vpc" {
		d.SetId(*allocResp.AllocationId)
	} else {
		d.SetId(*allocResp.PublicIp)
	}

	log.Printf("[INFO] EIP ID: %s (domain: %v)", d.Id(), *allocResp.Domain)
	return resourceOutscaleUpdate(d, meta)
}

func resourceOutscaleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	domain := resourceOutscaleDomain(d)

	// Associate to instance or interface if specified
	v_instance, ok_instance := d.GetOk("instance")
	v_interface, ok_interface := d.GetOk("network_interface")

	if ok_instance || ok_interface {
		instanceId := v_instance.(string)
		networkInterfaceId := v_interface.(string)

		assocOpts := &fcu.AssociateAddressInput{
			InstanceId: aws.String(instanceId),
			PublicIp:   aws.String(d.Id()),
		}

		// more unique ID conditionals
		if domain == "vpc" {
			var privateIpAddress *string
			if v := d.Get("associate_with_private_ip").(string); v != "" {
				privateIpAddress = aws.String(v)
			}
			assocOpts = &fcu.AssociateAddressInput{
				NetworkInterfaceId: aws.String(networkInterfaceId),
				InstanceId:         aws.String(instanceId),
				AllocationId:       aws.String(d.Id()),
				PrivateIpAddress:   privateIpAddress,
			}
		}

		log.Printf("[DEBUG] EIP associate configuration: %s (domain: %s)", assocOpts, domain)

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err := conn.VM.AssociateAddress(assocOpts)
			if err != nil {
				if awsErr, ok := err.(awserr.Error); ok {
					if awsErr.Code() == "InvalidAllocationID.NotFound" {
						return resource.RetryableError(awsErr)
					}
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			// Prevent saving instance if association failed
			// e.g. missing internet gateway in VPC
			d.Set("instance", "")
			d.Set("network_interface", "")
			return fmt.Errorf("Failure associating EIP: %s", err)
		}
	}

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
			Computed: true,
		},
		"association_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"domain": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"instance_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"network_interface_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"network_interface_owner_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"private_ip_address": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"public_ip": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
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
	conn := meta.(*OutscaleClient).FCU

	if err := resourcePublicIPRead(d, meta); err != nil {
		return err
	}
	if d.Id() == "" {
		// This might happen from the read
		return nil
	}

	v_instance, ok_instance := d.GetOk("instance")
	v_association_id, ok_association_id := d.GetOk("association_id")

	// If we are attached to an instance or interface, detach first.
	if (ok_instance && v_instance.(string) != "") || ok_association_id && v_association_id.(string) != "" {
		log.Printf("[DEBUG] Disassociating EIP: %s", d.Id())
		var err error
		switch resourceOutscaleDomain(d) {
		case "vpc":
			_, err = conn.VM.DisassociateAddress(&fcu.DisassociateAddressInput{
				AssociationId: aws.String(d.Get("association_id").(string)),
			})
		case "standard":
			_, err = conn.VM.DisassociateAddress(&fcu.DisassociateAddressInput{
				PublicIp: aws.String(d.Get("public_ip").(string)),
			})
		}

		if err != nil {
			// First check if the association ID is not found. If this
			// is the case, then it was already disassociated somehow,
			// and that is okay. The most commmon reason for this is that
			// the instance or ENI it was attached it was destroyed.
			if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidAssociationID.NotFound" {
				err = nil
			}
		}

		if err != nil {
			return err
		}
	}

	domain := resourceOutscaleDomain(d)
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		var err error
		switch domain {
		case "vpc":
			log.Printf(
				"[DEBUG] EIP release (destroy) address allocation: %v",
				d.Id())
			_, err = conn.VM.ReleaseAddress(&fcu.ReleaseAddressInput{
				AllocationId: aws.String(d.Id()),
			})
		case "standard":
			log.Printf("[DEBUG] EIP release (destroy) address: %v", d.Id())
			_, err = conn.VM.ReleaseAddress(&fcu.ReleaseAddressInput{
				PublicIp: aws.String(d.Id()),
			})
		}

		if err == nil {
			return nil
		}
		if _, ok := err.(awserr.Error); !ok {
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(err)
	})
}
