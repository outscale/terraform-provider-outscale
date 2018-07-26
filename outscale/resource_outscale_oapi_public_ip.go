package outscale

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
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
	conn := meta.(*OutscaleClient).OAPI

	domainOpt := resourceOutscaleOAPIPublicIPDomain(d)

	allocOpts := oapi.CreatePublicIpRequest{
		Placement: domainOpt,
	}

	fmt.Printf("[DEBUG] EIP create configuration: %#v", allocOpts)
	resp, err := conn.POST_CreatePublicIp(allocOpts)
	if err != nil {
		return fmt.Errorf("Error creating EIP: %s", err)
	}

	allocResp := resp.OK

	d.Set("placement", allocResp.Placement)

	fmt.Printf("[DEBUG] EIP Allocate: %#v", allocResp)
	if d.Get("placement").(string) == "vpc" {
		d.SetId(allocResp.ReservationId)
	} else {
		d.SetId(allocResp.PublicIp)
	}

	fmt.Printf("[INFO] EIP ID: %s (placement: %v)", d.Id(), allocResp.Placement)
	return resourceOutscaleOAPIPublicIPUpdate(d, meta)
}

func resourceOutscaleOAPIPublicIPRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	placement := resourceOutscaleOAPIPublicIPDomain(d)
	id := d.Id()

	req := oapi.ReadPublicIpsRequest{}

	//Not Used
	//filters := []oapi.Filters{}

	if placement == "vpc" {
		req.Filters.ReservationIds = []string{id}
	} else {
		req.Filters.PublicIps = []string{id}
	}

	var describeAddresses *oapi.ReadPublicIpsResponse
	resp, err := conn.POST_ReadPublicIps(req)
	describeAddresses = resp.OK

	if err != nil {
		if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving EIP: %s", err)
	}

	if len(describeAddresses.PublicIps) != 1 ||
		placement == "vpc" && describeAddresses.PublicIps[0].LinkId != id ||
		describeAddresses.PublicIps[0].PublicIp != id {
		if err != nil {
			return fmt.Errorf("Unable to find EIP: %#v", describeAddresses.PublicIps)
		}
	}

	address := describeAddresses.PublicIps[0]

	fmt.Printf("[DEBUG] EIP read configuration: %+v", address)

	if address.LinkId != "" {
		d.Set("link_id", address.LinkId)
	} else {
		d.Set("link_id", "")
	}
	if address.PublicIp != "" {
		d.Set("vm_id", address.VmId)
	} else {
		d.Set("vm_id", "")
	}
	if address.NicId != "" {
		d.Set("nic_id", address.NicId)
	} else {
		d.Set("nic_id", "")
	}
	if address.NicAccountId != "" {
		d.Set("nic_account_id", address.NicAccountId)
	} else {
		d.Set("nic_account_id", "")
	}
	d.Set("private_ip", address.PrivateIp)
	d.Set("public_ip", address.PublicIp)

	d.Set("placement", address.Placement)
	d.Set("reservation_id", address.ReservationId)

	if address.Placement == "vpc" && net.ParseIP(id) != nil {
		fmt.Printf("[DEBUG] Re-assigning EIP ID (%s) to it's Allocation ID (%s)", d.Id(), address.ReservationId)
		d.SetId(address.ReservationId)
	} else {
		d.SetId(address.PublicIp)
	}

	return nil
}

func resourceOutscaleOAPIPublicIPUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	placement := resourceOutscaleOAPIPublicIPDomain(d)

	vInstance, okInstance := d.GetOk("vm_id")
	vInterface, okInterface := d.GetOk("nic_id")

	if okInstance || okInterface {
		instanceID := vInstance.(string)
		networkInterfaceID := vInterface.(string)

		assocOpts := oapi.LinkPublicIpRequest{
			VmId:     instanceID,
			PublicIp: d.Id(),
		}

		if placement == "vpc" {
			var privateIPAddress string
			if v := d.Get("private_ip").(string); v != "" {
				privateIPAddress = v
			}
			assocOpts = oapi.LinkPublicIpRequest{
				NicId:         networkInterfaceID,
				VmId:          instanceID,
				ReservationId: d.Id(),
				PrivateIp:     privateIPAddress,
			}
		}

		err := resource.Retry(120*time.Second, func() *resource.RetryError {
			var err error
			_, err = conn.POST_LinkPublicIp(assocOpts)

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
	conn := meta.(*OutscaleClient).OAPI

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
			_, err = conn.POST_UnlinkPublicIp(oapi.UnlinkPublicIpRequest{
				LinkId: d.Get("link_id").(string),
			})
		case "standard":
			_, err = conn.POST_UnlinkPublicIp(oapi.UnlinkPublicIpRequest{
				PublicIp: d.Get("public_ip").(string),
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
			_, err = conn.POST_DeletePublicIp(oapi.DeletePublicIpRequest{
				ReservationId: d.Id(),
			})
		case "standard":
			fmt.Printf("[DEBUG] EIP release (destroy) address: %v", d.Id())
			_, err = conn.POST_DeletePublicIp(oapi.DeletePublicIpRequest{
				PublicIp: d.Id(),
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
