package outscale

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/aws-sdk-go/service/ec2"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleDHCPOptionLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleDHCPOptionLinkCreate,
		Read:   resourceOutscaleDHCPOptionLinkRead,
		Delete: resourceOutscaleDHCPOptionLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: getDHCPOptionLinkSchema(),
	}
}

func getDHCPOptionLinkSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dhcp_options_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"vpc_id": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func resourceOutscaleDHCPOptionLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	log.Printf(
		"[INFO] Creating DHCP Options Link: %s => %s",
		d.Get("vpc_id").(string),
		d.Get("dhcp_options_id").(string))

	optsID := aws.String(d.Get("dhcp_options_id").(string))
	vpcID := aws.String(d.Get("vpc_id").(string))

	if _, err := conn.VM.AssociateDhcpOptions(&fcu.AssociateDhcpOptionsInput{
		DhcpOptionsId: optsID,
		VpcId:         vpcID,
	}); err != nil {
		return err
	}

	// Set the ID and return
	d.SetId(*optsID + "-" + *vpcID)
	log.Printf("[INFO] Association ID: %s", d.Id())

	return nil

}

func resourceOutscaleDHCPOptionLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	// Get the VPC that this association belongs to
	vpcRaw, _, err := VPCStateRefreshFunc(conn, d.Get("vpc_id").(string))()

	if err != nil {
		return err
	}

	if vpcRaw == nil {
		return nil
	}

	vpc := vpcRaw.(*ec2.Vpc)
	if *vpc.VpcId != d.Get("vpc_id") || *vpc.DhcpOptionsId != d.Get("dhcp_options_id") {
		log.Printf("[INFO] It seems the DHCP Options Link is gone. Deleting reference from Graph...")
		d.SetId("")
	}

	return nil
}

func resourceOutscaleDHCPOptionLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	log.Printf("[INFO] Disassociating DHCP Options Set %s from VPC %s...", d.Get("dhcp_options_id"), d.Get("vpc_id"))
	if _, err := conn.VM.AssociateDhcpOptions(&fcu.AssociateDhcpOptionsInput{
		DhcpOptionsId: aws.String("default"),
		VpcId:         aws.String(d.Get("vpc_id").(string)),
	}); err != nil {
		return err
	}

	d.SetId("")
	return nil

}

// DHCP Options Asociations cannot be updated.
func resourceAwsVpcDhcpOptionsAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceOutscaleDHCPOptionLinkCreate(d, meta)
}

func VPCStateRefreshFunc(conn *fcu.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		DescribeVpcOpts := &fcu.DescribeVpcsInput{
			VpcIds: []*string{aws.String(id)},
		}
		resp, err := conn.VM.DescribeVpcs(DescribeVpcOpts)
		if err != nil {
			if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidVpcID.NotFound" {
				resp = nil
			} else {
				log.Printf("Error on VPCStateRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil {
			// Sometimes AWS just has consistency issues and doesn't see
			// our instance yet. Return an empty state.
			return nil, "", nil
		}

		vpc := resp.Vpcs[0]
		return vpc, *vpc.State, nil
	}
}
