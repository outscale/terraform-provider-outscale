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

func resourceOutscaleSubNet() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleSubNetCreate,
		Read:   resourceOutscaleSubNetRead,
		Delete: resourceOutscaleSubNetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: getSubNetSchema(),
	}
}

//Create SubNet
func resourceOutscaleSubNetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	req := &fcu.CreateSubnetInput{
		CidrBlock: aws.String(d.Get("cidr_block").(string)),
		VpcId:     aws.String(d.Get("vpc_id").(string)),
	}
	if a, aok := d.GetOk("availability_zone"); aok {
		req.AvailabilityZone = aws.String(a.(string))
	}
	if a, aok := d.GetOk("dry_run"); aok {
		req.DryRun = aws.Bool(a.(bool))
	}
	var res *fcu.CreateSubnetOutput
	var err error
	err = resource.Retry(40*time.Minute, func() *resource.RetryError {
		res, err = conn.VM.CreateSubNet(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(err)
	})

	d.SetId(*res.Subnet.SubnetId)

	return resourceOutscaleSubNetRead(d, meta)
}

//Read SubNet

func resourceOutscaleSubNetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	id := d.Id()

	log.Printf("[DEBUG] Reading Subnet(%s)", id)

	req := &fcu.DescribeSubnetsInput{
		SubnetIds: []*string{aws.String(id)},
	}

	var resp *fcu.DescribeSubnetsOutput
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.VM.DescribeSubNet(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		log.Printf("[DEBUG] Error reading Subnet (%s)", err)
	}

	log.Printf("[DEBUG] Setting Subnet (%s)", err)

	d.Set("subnet_id", resp.Subnets[0].SubnetId)
	d.Set("availability_zone", resp.Subnets[0].AvailabilityZone)
	d.Set("cidr_block", resp.Subnets[0].CidrBlock)
	d.Set("vpc_id", resp.Subnets[0].VpcId)

	if err := d.Set("tag_set", dataSourceTags(resp.Subnets[0].Tags)); err != nil {
		return err
	}

	return nil
}

func resourceOutscaleSubNetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	id := d.Id()
	log.Printf("[DEBUG] Deleting Subnet (%s)", id)

	req := &fcu.DeleteSubnetInput{
		SubnetId: &id,
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.VM.DeleteSubNet(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		log.Printf("[DEBUG] Error deleting Subnet(%s)", err)
		return err
	}

	return nil
}

func readOutscaleSubNet(d *schema.ResourceData, subnet *fcu.Subnet) error {
	if err := d.Set("availability_zone", subnet.AvailabilityZone); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet1 (%s)", err)

		return err
	}
	if err := d.Set("available_ip_address_count", subnet.AvailableIpAddressCount); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet2 (%s)", err)

		return err
	}
	if err := d.Set("cidr_block", subnet.CidrBlock); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet3 (%s)", err)

		return err
	}
	if err := d.Set("state", subnet.State); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet4 (%s)", err)

		return err
	}
	if err := d.Set("subnet_id", subnet.SubnetId); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet5 (%s)", err)

		return err
	}

	if err := d.Set("vpc_id", subnet.VpcId); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet6 (%s)", err)

		return err
	}

	return nil
}

func getSubNetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		//This is attribute part for schema SubNet
		"vpc_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"cidr_block": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"availability_zone": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		//This is arguments part for schema SubNet
		"available_ip_address_count": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},

		"state": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"subnet_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"request_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"tag_set": dataSourceTagsSchema(),
	}
}
