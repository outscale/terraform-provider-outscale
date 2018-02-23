package outscale

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscalePublicIPLink() *schema.Resource {
	return &schema.Resource{
		Create: resourcePublicIPLinkCreate,
		Read:   resourcePublicIPLinkRead,
		Delete: resourcePublicIPLinkDelete,
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

func resourcePublicIPLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	request := &fcu.AssociateAddressInput{}

	if v, ok := d.GetOk("allocation_id"); ok {
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

	fmt.Printf("[DEBUG] EIP association configuration: %#v", request)

	resp, err := conn.VM.AssociateAddress(request)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			return fmt.Errorf("[WARN] Error attaching EIP, message: \"%s\", code: \"%s\"",
				awsErr.Message(), awsErr.Code())
		}
		return err
	}

	fmt.Printf("\n [DEBUG] resourcePublicIPLinkCreate Error 3: %v", err)

	if resp.AssociationId != nil {
		d.SetId(*resp.AssociationId)
	} else {
		d.SetId("")
	}

	return resourcePublicIPLinkRead(d, meta)

}

func resourcePublicIPLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	request := &fcu.DescribeAddressesInput{
		Filters: []*fcu.Filter{
			&fcu.Filter{
				Name:   aws.String("association-id"),
				Values: []*string{aws.String(d.Id())},
			},
		},
	}

	response, err := conn.VM.DescribeAddressesRequest(request)

	fmt.Printf("\n [DEBUG] resourcePublicIPLinkRead Error 3: %v", err)

	if err != nil {
		return fmt.Errorf("Error reading Outscale VM Public IP %s: %#v", d.Get("allocation_id").(string), err)
	}

	if response.Addresses == nil || len(response.Addresses) == 0 {
		fmt.Printf("[INFO] EIP Association ID Not Found. Refreshing from state")
		d.SetId("")
		return nil
	}

	return readOutscalePublicIPAssociation(d, response.Addresses[0])
}

func resourcePublicIPLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	fmt.Printf("\n [DEBUG] ID => %s", d.Id())

	opts := &fcu.DisassociateAddressInput{
		AssociationId: aws.String(d.Id()),
	}

	_, err := conn.VM.DisassociateAddress(opts)

	fmt.Printf("\n [DEBUG] resourcePublicIPLinkDelete Error 2: %v", err)

	if err != nil {
		return fmt.Errorf("Error deleting Public IP association: %s", err)
	}

	return nil
}

func readOutscalePublicIPAssociation(d *schema.ResourceData, address *fcu.Address) error {
	if err := d.Set("allocation_id", address.AllocationId); err != nil {
		fmt.Printf("\n [DEBUG] readOutscalePublicIPAssociation Error 2: %v", err)

		return err
	}
	if err := d.Set("instance_id", address.InstanceId); err != nil {
		fmt.Printf("\n [DEBUG] readOutscalePublicIPAssociation Error 2: %v", err)

		return err
	}
	if err := d.Set("network_interface_id", address.NetworkInterfaceId); err != nil {
		fmt.Printf("\n [DEBUG] readOutscalePublicIPAssociation Error 2: %v", err)

		return err
	}
	if err := d.Set("private_ip_address", address.PrivateIpAddress); err != nil {
		fmt.Printf("\n [DEBUG] readOutscalePublicIPAssociation Error 2: %v", err)

		return err
	}
	if err := d.Set("public_ip", address.PublicIp); err != nil {
		fmt.Printf("\n [DEBUG] readOutscalePublicIPAssociation Error 2: %v", err)

		return err
	}
	if err := d.Set("association_id", address.AssociationId); err != nil {
		fmt.Printf("\n [DEBUG] readOutscalePublicIPAssociation Error 2: %v", err)

		return err
	}
	if err := d.Set("domain", address.Domain); err != nil {
		fmt.Printf("\n [DEBUG] readOutscalePublicIPAssociation Error 2: %v", err)

		return err
	}
	if err := d.Set("network_interface_owner_id", address.NetworkInterfaceOwnerId); err != nil {
		fmt.Printf("\n [DEBUG] readOutscalePublicIPAssociation Error 2: %v", err)

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
		"request_id": &schema.Schema{
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
	}
}
