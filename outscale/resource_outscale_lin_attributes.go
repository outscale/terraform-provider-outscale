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

func resourceOutscaleLinAttributes() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleLinAttrCreate,
		Read:   resourceOutscaleLinAttrRead,
		Update: resourceOutscaleLinAttrUpdate,
		Delete: resourceOutscaleLinAttrDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"enable_dns_hostnames": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"enable_dns_support": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"attribute": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceOutscaleLinAttrCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.ModifyVpcAttributeInput{}

	req.VpcId = aws.String(d.Get("vpc_id").(string))

	if c, ok := d.GetOk("enable_dns_hostnames"); ok {
		req.EnableDnsHostnames = &fcu.AttributeBooleanValue{Value: aws.Bool(c.(bool))}
	}
	if c, ok := d.GetOk("enable_dns_support"); ok {
		req.EnableDnsHostnames = &fcu.AttributeBooleanValue{Value: aws.Bool(c.(bool))}
	}

	var resp *fcu.ModifyVpcAttributeOutput
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.VM.ModifyVpcAttribute(req)

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

	d.SetId(resource.UniqueId())

	return resourceOutscaleLinAttrRead(d, meta)
}

func resourceOutscaleLinAttrUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.ModifyVpcAttributeInput{}

	if d.HasChange("vpc_id") && !d.IsNewResource() {
		req.VpcId = aws.String(d.Get("vpc_id").(string))
	}
	if d.HasChange("enable_dns_hostnames") && !d.IsNewResource() {
		req.EnableDnsHostnames = &fcu.AttributeBooleanValue{Value: aws.Bool(d.Get("enable_dns_hostnames").(bool))}
	}
	if d.HasChange("enable_dns_support") && !d.IsNewResource() {
		req.EnableDnsHostnames = &fcu.AttributeBooleanValue{Value: aws.Bool(d.Get("enable_dns_support").(bool))}
	}

	var resp *fcu.ModifyVpcAttributeOutput
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.VM.ModifyVpcAttribute(req)

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

	return resourceOutscaleLinAttrRead(d, meta)
}

func resourceOutscaleLinAttrRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeVpcAttributeInput{
		Attribute: aws.String(d.Get("attribute").(string)),
		VpcId:     aws.String(d.Get("vpc_id").(string)),
	}

	var resp *fcu.DescribeVpcAttributeOutput
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.VM.DescribeVpcAttribute(req)

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

	d.Set("vpc_id", resp.VpcId)
	d.Set("enable_dns_hostnames", resp.EnableDnsHostnames)
	d.Set("enable_dns_support", resp.EnableDnsSupport)

	return nil
}

func resourceOutscaleLinAttrDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	return nil
}
