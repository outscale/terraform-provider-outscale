ckage outscale

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

	log.Println("[DEBUG] Creating Subnet")
	r, err := conn.VM.CreateSubNet(nil)
	if err != nil {
		log.Printf("[DEBUG] Subnet %s", err)

		return err
	}

	d.SetId(*r.Subnet.SubnetId)
	d.Set("availability_zone", *r.Subnet.AvailabilityZone)
	d.Set("cidr_block", *r.Subnet.CidrBlock)
	d.Set("vpc_id", *r.Subnet.VpcId)

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

	d.Set("subnet_id", resp.Subnets)

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
		"availability_zone": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		"available_ip_address_count": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		"cidr_block": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			Computed: true,
			ForceNew: true,
		},
		"state": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"subnet_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		"request_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"vpc_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			Computed: true,
			ForceNew: true,
		},
		"tags": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"value": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
			Computed: true,
		},
		"tag": tagsSchema(),
	}
}