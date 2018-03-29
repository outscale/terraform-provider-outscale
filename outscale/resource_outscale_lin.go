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

func resourceOutscaleLin() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleLinCreate,
		Read:   resourceOutscaleLinRead,
		Delete: resourceOutscaleLinDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: getLinSchema(),
	}
}

func resourceOutscaleLinCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.CreateVpcInput{}

	req.CidrBlock = aws.String(d.Get("cidr_block").(string))

	if c, ok := d.GetOk("instance_tenancy"); ok {
		cidr := c.(string)
		if cidr == "default" || cidr == "dedicated" {
			req.InstanceTenancy = aws.String(cidr)
		} else {
			return fmt.Errorf("cidr_block option not supported %s", cidr)
		}
	}

	var resp *fcu.CreateVpcOutput
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.VM.CreateVpc(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		log.Printf("[DEBUG] Error creating lin (%s)", err)
		return err
	}

	if resp == nil {
		return fmt.Errorf("Cannot create the vpc, empty response")
	}

	d.SetId(*resp.Vpc.VpcId)

	return resourceOutscaleLinRead(d, meta)
}

func resourceOutscaleLinRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	id := d.Id()

	req := &fcu.DescribeVpcsInput{
		VpcIds: []*string{aws.String(id)},
	}

	var resp *fcu.DescribeVpcsOutput
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.VM.DescribeVpcs(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		log.Printf("[DEBUG] Error reading lin (%s)", err)
	}

	if resp == nil {
		d.SetId("")
		return fmt.Errorf("Lin not found")
	}

	if len(resp.Vpcs) == 0 {
		d.SetId("")
		return fmt.Errorf("Lin not found")
	}

	d.Set("cidr_block", resp.Vpcs[0].CidrBlock)
	d.Set("instance_tenancy", resp.Vpcs[0].InstanceTenancy)
	d.Set("dhcp_options_id", resp.Vpcs[0].DhcpOptionsId)
	d.Set("request_id", resp.RequesterId)
	d.Set("state", resp.Vpcs[0].State)
	d.Set("vpc_id", resp.Vpcs[0].VpcId)

	if err := d.Set("tag_set", dataSourceTags(resp.Vpcs[0].Tags)); err != nil {
		return err
	}

	return nil
}

func resourceOutscaleLinDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	id := d.Id()

	req := &fcu.DeleteVpcInput{
		VpcId: &id,
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.VM.DeleteVpc(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func getLinSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cidr_block": {
			Type:     schema.TypeString,
			ForceNew: true,
			Required: true,
		},
		"instance_tenancy": {
			Type:     schema.TypeString,
			ForceNew: true,
			Computed: true,
			Optional: true,
		},

		// Attributes
		"dhcp_options_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vpc_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tag_set": dataSourceTagsSchema(),
	}
}
