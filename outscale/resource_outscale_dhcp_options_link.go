package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			ForceNew: true,
		},
		"vpc_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceOutscaleDHCPOptionLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	fmt.Printf(
		"[INFO] Creating DHCP Options Link: %s => %s",
		d.Get("vpc_id").(string),
		d.Get("dhcp_options_id").(string))

	optsID := aws.String(d.Get("dhcp_options_id").(string))
	vpcID := aws.String(d.Get("vpc_id").(string))

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.AssociateDhcpOptions(&fcu.AssociateDhcpOptionsInput{
			DhcpOptionsId: optsID,
			VpcId:         vpcID,
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Set the ID and return
	d.SetId(*optsID + "-" + *vpcID)
	fmt.Printf("[INFO] Association ID: %s", d.Id())

	return resourceOutscaleDHCPOptionLinkRead(d, meta)

}

func resourceOutscaleDHCPOptionLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	var resp *fcu.DescribeVpcsOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		DescribeVpcOpts := &fcu.DescribeVpcsInput{
			VpcIds: []*string{aws.String(d.Get("vpc_id").(string))},
		}
		resp, err = conn.VM.DescribeVpcs(DescribeVpcOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	vpc := resp.Vpcs[0]
	if *vpc.VpcId != d.Get("vpc_id") || *vpc.DhcpOptionsId != d.Get("dhcp_options_id") {
		fmt.Printf("[INFO] It seems the DHCP Options Link is gone. Deleting reference from Graph...")
		d.SetId("")
	}

	d.Set("request_id", resp.RequestId)

	return nil
}

func resourceOutscaleDHCPOptionLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	fmt.Printf("[INFO] Disassociating DHCP Options Set %s from VPC %s...", d.Get("dhcp_options_id"), d.Get("vpc_id"))

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.AssociateDhcpOptions(&fcu.AssociateDhcpOptionsInput{
			DhcpOptionsId: aws.String(d.Get("dhcp_options_id").(string)),
			VpcId:         aws.String(d.Get("vpc_id").(string)),
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId("")
	return nil

}

// VPCStateRefreshFunc ...
func VPCStateRefreshFunc(conn *fcu.Client, ID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		DescribeVpcOpts := &fcu.DescribeVpcsInput{
			VpcIds: []*string{aws.String(ID)},
		}
		resp, err := conn.VM.DescribeVpcs(DescribeVpcOpts)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidVpcID.NotFound") {
				resp = nil
			} else {
				log.Printf("Error on VPCStateRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil {
			return nil, "", nil
		}

		vpc := resp.Vpcs[0]
		return vpc, *vpc.State, nil
	}
}
