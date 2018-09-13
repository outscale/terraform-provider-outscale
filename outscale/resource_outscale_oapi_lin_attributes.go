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

func resourceOutscaleOAPILinAttributes() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILinAttrCreate,
		Read:   resourceOutscaleOAPILinAttrRead,
		Update: resourceOutscaleOAPILinAttrUpdate,
		Delete: resourceOutscaleOAPILinAttrDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"dns_hostnames_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"dns_support_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"net_id": {
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

func resourceOutscaleOAPILinAttrCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.ModifyVpcAttributeInput{}

	req.VpcId = aws.String(d.Get("net_id").(string))

	if c, ok := d.GetOk("dns_hostnames_enabled"); ok {
		req.EnableDnsHostnames = &fcu.AttributeBooleanValue{Value: aws.Bool(c.(bool))}
	}
	if c, ok := d.GetOk("dns_support_enabled"); ok {
		req.EnableDnsHostnames = &fcu.AttributeBooleanValue{Value: aws.Bool(c.(bool))}
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.VM.ModifyVpcAttribute(req)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[DEBUG] Error creating lin (%s)", err)
		return err
	}

	d.SetId(resource.UniqueId())

	return resourceOutscaleOAPILinAttrRead(d, meta)
}

func resourceOutscaleOAPILinAttrUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.ModifyVpcAttributeInput{}

	if d.HasChange("net_id") && !d.IsNewResource() {
		req.VpcId = aws.String(d.Get("net_id").(string))
	}
	if d.HasChange("dns_hostnames_enabled") && !d.IsNewResource() {
		req.EnableDnsHostnames = &fcu.AttributeBooleanValue{Value: aws.Bool(d.Get("dns_hostnames_enabled").(bool))}
	}
	if d.HasChange("dns_support_enabled") && !d.IsNewResource() {
		req.EnableDnsHostnames = &fcu.AttributeBooleanValue{Value: aws.Bool(d.Get("dns_support_enabled").(bool))}
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.VM.ModifyVpcAttribute(req)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[DEBUG] Error creating lin (%s)", err)
		return err
	}

	return resourceOutscaleOAPILinAttrRead(d, meta)
}

func resourceOutscaleOAPILinAttrRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeVpcAttributeInput{
		Attribute: aws.String(d.Get("attribute").(string)),
		VpcId:     aws.String(d.Get("net_id").(string)),
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

	d.Set("net_id", resp.VpcId)
	if resp.EnableDnsHostnames != nil {
		d.Set("dns_hostnames_enabled", *resp.EnableDnsHostnames.Value)
	}
	if resp.EnableDnsSupport != nil {
		d.Set("dns_support_enabled", *resp.EnableDnsSupport.Value)
	}

	return nil
}

func resourceOutscaleOAPILinAttrDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	return nil
}
