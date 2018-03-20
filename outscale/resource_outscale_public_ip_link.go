package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscalePublicIPLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscalePublicIPLinkCreate,
		Read:   resourceOutscalePublicIPLinkRead,
		Delete: resourceOutscalePublicIPLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: getPublicIPLinkSchema(),
	}
}

func resourceOutscalePublicIPLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	request := &fcu.AssociateAddressInput{}

	if v, ok := d.GetOk("allocation_id"); ok {
		fmt.Println(v.(string))
		request.AllocationId = aws.String(v.(string))
	}
	if v, ok := d.GetOk("allow_reassociation"); ok {
		request.AllowReassociation = aws.Bool(v.(bool))
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = aws.String(v.(string))
	}
	if v, ok := d.GetOk("network_interface_id"); ok {
		request.NetworkInterfaceId = aws.String(v.(string))
	}
	if v, ok := d.GetOk("private_ip_address"); ok {
		request.PrivateIpAddress = aws.String(v.(string))
	}
	if v, ok := d.GetOk("public_ip"); ok {
		request.PublicIp = aws.String(v.(string))
	}

	log.Printf("[DEBUG] EIP association configuration: %#v", request)

	var resp *fcu.AssociateAddressOutput
	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {

		resp, err = conn.VM.AssociateAddress(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})

	if err != nil {
		fmt.Printf("[WARN] ERROR resourceOutscalePublicIPLinkCreate (%s)", err)
		return err
	}

	if resp != nil && resp.AssociationId != nil && len(*resp.AssociationId) > 0 {
		d.SetId(*resp.AssociationId)
	} else {
		d.SetId(*request.PublicIp)
	}

	return resourceOutscalePublicIPLinkRead(d, meta)
}

func resourceOutscalePublicIPLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	id := d.Id()
	var request *fcu.DescribeAddressesInput

	if strings.Contains(id, "eipassoc") {
		request = &fcu.DescribeAddressesInput{
			Filters: []*fcu.Filter{
				&fcu.Filter{
					Name:   aws.String("association-id"),
					Values: []*string{aws.String(id)},
				},
			},
		}
	} else {
		request = &fcu.DescribeAddressesInput{
			Filters: []*fcu.Filter{
				&fcu.Filter{
					Name:   aws.String("public-ip"),
					Values: []*string{aws.String(id)},
				},
			},
		}
	}

	var response *fcu.DescribeAddressesOutput
	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {

		response, err = conn.VM.DescribeAddressesRequest(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})

	fmt.Printf("[WARN] ERROR resourceOutscalePublicIPLinkRead (%s)", err)

	if err != nil {
		return fmt.Errorf("Error reading Outscale VM Public IP %s: %#v", d.Get("allocation_id").(string), err)
	}

	if response.Addresses == nil || len(response.Addresses) == 0 {
		log.Printf("[INFO] EIP Association ID Not Found. Refreshing from state")
		d.SetId("")
		return nil
	}

	d.Set("request_id", *response.RequestId)

	return readOutscalePublicIPLink(d, response.Addresses[0])
}

func resourceOutscalePublicIPLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	assocId := d.Get("association_id")

	opts := &fcu.DisassociateAddressInput{
		AssociationId: aws.String(assocId.(string)),
	}

	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {

		_, err = conn.VM.DisassociateAddress(opts)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})

	fmt.Printf("[WARN] ERROR resourceOutscalePublicIPLinkDelete (%s)", err)

	if err != nil {
		return fmt.Errorf("Error deleting Elastic IP association: %s", err)
	}

	return nil
}

func readOutscalePublicIPLink(d *schema.ResourceData, address *fcu.Address) error {
	if err := d.Set("allocation_id", address.AllocationId); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink1 (%s)", err)

		return err
	}
	if err := d.Set("instance_id", address.InstanceId); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink2 (%s)", err)

		return err
	}
	if err := d.Set("network_interface_id", address.NetworkInterfaceId); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink3 (%s)", err)

		return err
	}
	if err := d.Set("private_ip_address", address.PrivateIpAddress); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink4 (%s)", err)

		return err
	}
	if err := d.Set("public_ip", address.PublicIp); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink (%s)", err)

		return err
	}

	if err := d.Set("association_id", address.AssociationId); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink (%s)", err)

		return err
	}

	if err := d.Set("allow_reassociation", address.AllowReassociation); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink (%s)", err)

		return err
	}

	return nil
}

func getPublicIPLinkSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allocation_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"association_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
		"domain": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
		"allow_reassociation": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"instance_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"network_interface_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"network_interface_owner_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},

		"private_ip_address": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"public_ip": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		"request_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
