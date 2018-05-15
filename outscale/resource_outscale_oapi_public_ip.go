package outscale

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleOAPIPublicIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIPublicIPCreate,
		Read:   resourceOutscaleOAPIPublicIPRead,
		Delete: resourceOutscaleOAPIPublicIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"reservation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"link_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"placement": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"vm_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nic_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nic_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPIPublicIPCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	domainOpt := resourceOutscaleOAPIPublicIPDomain(d)

	allocOpts := &fcu.AllocateAddressInput{
		Domain: aws.String(domainOpt),
	}

	fmt.Printf("[DEBUG] EIP create configuration: %#v", allocOpts)
	allocResp, err := conn.VM.AllocateAddress(allocOpts)
	if err != nil {
		return fmt.Errorf("Error creating EIP: %s", err)
	}

	d.Set("placement", allocResp.Domain)

	fmt.Printf("[DEBUG] EIP Allocate: %#v", allocResp)
	if d.Get("placement").(string) == "vpc" {
		d.SetId(*allocResp.AllocationId)
	} else {
		d.SetId(*allocResp.PublicIp)
	}

	fmt.Printf("[INFO] EIP ID: %s (placement: %v)", d.Id(), *allocResp.Domain)
	return resourceOutscaleOAPIPublicIPUpdate(d, meta)
}

func resourceOutscaleOAPIPublicIPRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	placement := resourceOutscaleOAPIPublicIPDomain(d)
	id := d.Id()

	req := &fcu.DescribeAddressesInput{}

	if placement == "vpc" {
		req.AllocationIds = []*string{aws.String(id)}
	} else {
		req.PublicIps = []*string{aws.String(id)}
	}

	var describeAddresses *fcu.DescribeAddressesOutput
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		var err error
		describeAddresses, err = conn.VM.DescribeAddressesRequest(req)

		return resource.RetryableError(err)
	})

	if err != nil {
		if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving EIP: %s", err)
	}

	if len(describeAddresses.Addresses) != 1 ||
		placement == "vpc" && *describeAddresses.Addresses[0].AllocationId != id ||
		*describeAddresses.Addresses[0].PublicIp != id {
		if err != nil {
			return fmt.Errorf("Unable to find EIP: %#v", describeAddresses.Addresses)
		}
	}

	address := describeAddresses.Addresses[0]

	fmt.Printf("[DEBUG] EIP read configuration: %+v", *address)

	if address.AssociationId != nil {
		d.Set("link_id", address.AssociationId)
	} else {
		d.Set("link_id", "")
	}
	if address.InstanceId != nil {
		d.Set("vm_id", address.InstanceId)
	} else {
		d.Set("vm_id", "")
	}
	if address.NetworkInterfaceId != nil {
		d.Set("nic_id", address.NetworkInterfaceId)
	} else {
		d.Set("nic_id", "")
	}
	if address.NetworkInterfaceOwnerId != nil {
		d.Set("nic_account_id", address.NetworkInterfaceOwnerId)
	} else {
		d.Set("nic_account_id", "")
	}
	d.Set("private_ip", address.PrivateIpAddress)
	d.Set("public_ip", address.PublicIp)

	d.Set("placement", address.Domain)
	d.Set("reservation_id", address.AllocationId)

	if *address.Domain == "vpc" && net.ParseIP(id) != nil {
		fmt.Printf("[DEBUG] Re-assigning EIP ID (%s) to it's Allocation ID (%s)", d.Id(), *address.AllocationId)
		d.SetId(*address.AllocationId)
	} else {
		d.SetId(*address.PublicIp)
	}

	return nil
}

func resourceOutscaleOAPIPublicIPUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	placement := resourceOutscaleOAPIPublicIPDomain(d)

	vInstance, okInstance := d.GetOk("vm_id")
	vInterface, okInterface := d.GetOk("nic_id")

	if okInstance || okInterface {
		instanceID := vInstance.(string)
		networkInterfaceID := vInterface.(string)

		assocOpts := &fcu.AssociateAddressInput{
			InstanceId: aws.String(instanceID),
			PublicIp:   aws.String(d.Id()),
		}

		if placement == "vpc" {
			var privateIPAddress *string
			if v := d.Get("private_ip").(string); v != "" {
				privateIPAddress = aws.String(v)
			}
			assocOpts = &fcu.AssociateAddressInput{
				NetworkInterfaceId: aws.String(networkInterfaceID),
				InstanceId:         aws.String(instanceID),
				AllocationId:       aws.String(d.Id()),
				PrivateIpAddress:   privateIPAddress,
			}
		}

		err := resource.Retry(120*time.Second, func() *resource.RetryError {
			var err error
			_, err = conn.VM.AssociateAddress(assocOpts)

			if err != nil {
				if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err != nil {
			d.Set("vm_id", "")
			d.Set("nic_id", "")
			return fmt.Errorf("Failure associating EIP: %s", err)
		}
	}

	return resourceOutscaleOAPIPublicIPRead(d, meta)
}

func resourceOutscaleOAPIPublicIPDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	if err := resourceOutscaleOAPIPublicIPRead(d, meta); err != nil {
		return err
	}
	if d.Id() == "" {
		return nil
	}

	vInstance, okInstance := d.GetOk("vm_id")
	vAssociationID, okAssociationID := d.GetOk("link_id")

	if (okInstance && vInstance.(string) != "") || (okAssociationID && vAssociationID.(string) != "") {
		fmt.Printf("[DEBUG] Disassociating EIP: %s", d.Id())
		var err error
		switch resourceOutscaleOAPIPublicIPDomain(d) {
		case "vpc":
			_, err = conn.VM.DisassociateAddress(&fcu.DisassociateAddressInput{
				AssociationId: aws.String(d.Get("link_id").(string)),
			})
		case "standard":
			_, err = conn.VM.DisassociateAddress(&fcu.DisassociateAddressInput{
				PublicIp: aws.String(d.Get("public_ip").(string)),
			})
		}

		if err != nil {
			if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
				return nil
			}
			return err
		}
	}

	placement := resourceOutscaleOAPIPublicIPDomain(d)
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		var err error
		switch placement {
		case "vpc":
			fmt.Printf(
				"[DEBUG] EIP release (destroy) address allocation: %v",
				d.Id())
			_, err = conn.VM.ReleaseAddress(&fcu.ReleaseAddressInput{
				AllocationId: aws.String(d.Id()),
			})
		case "standard":
			fmt.Printf("[DEBUG] EIP release (destroy) address: %v", d.Id())
			_, err = conn.VM.ReleaseAddress(&fcu.ReleaseAddressInput{
				PublicIp: aws.String(d.Id()),
			})
		}

		if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
			return nil
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

func resourceOutscaleOAPIPublicIPDomain(d *schema.ResourceData) string {
	if v, ok := d.GetOk("placement"); ok {
		return v.(string)
	} else if strings.Contains(d.Id(), "eipalloc") {
		return "vpc"
	}

	return "standard"
}
