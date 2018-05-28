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

func resourceOutscaleOAPISubNet() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPISubNetCreate,
		Read:   resourceOutscaleOAPISubNetRead,
		Delete: resourceOutscaleOAPISubNetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: getOAPISubNetSchema(),
	}
}

//Create SubNet
func resourceOutscaleOAPISubNetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	req := &fcu.CreateSubnetInput{
		CidrBlock: aws.String(d.Get("ip_range").(string)),
		VpcId:     aws.String(d.Get("lin_id").(string)),
	}
	if a, aok := d.GetOk("sub_region_name"); aok {
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

	return resourceOutscaleOAPISubNetRead(d, meta)
}

//Read SubNet

func resourceOutscaleOAPISubNetRead(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("sub_region_name", resp.Subnets[0].AvailabilityZone)
	d.Set("ip_range", resp.Subnets[0].CidrBlock)
	d.Set("lin_id", resp.Subnets[0].VpcId)

	return d.Set("tag", tagsToMap(resp.Subnets[0].Tags))
}

func resourceOutscaleOAPISubNetDelete(d *schema.ResourceData, meta interface{}) error {
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

func readOutscaleOAPISubNet(d *schema.ResourceData, subnet *fcu.Subnet) error {
	if err := d.Set("sub_region_name", subnet.AvailabilityZone); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet1 (%s)", err)

		return err
	}
	if err := d.Set("available_ips_count", subnet.AvailableIpAddressCount); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet2 (%s)", err)

		return err
	}
	if err := d.Set("ip_range", subnet.CidrBlock); err != nil {
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

	if err := d.Set("lin_id", subnet.VpcId); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet6 (%s)", err)

		return err
	}

	return nil
}

func getOAPISubNetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		//This is attribute part for schema SubNet
		"lin_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"ip_range": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"sub_region_name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		//This is arguments part for schema SubNet
		"available_ips_count": &schema.Schema{
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
		"tag": dataSourceTagsSchema(),
	}
}
