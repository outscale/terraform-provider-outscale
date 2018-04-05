package outscale

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleLinInternetGatewayLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleLinInternetGatewayLinkCreate,
		Read:   resourceOutscaleLinInternetGatewayLinkRead,
		Delete: resourceOutscaleLinInternetGatewayLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: getLinInternetGatewayLinkSchema(),
	}
}

func resourceOutscaleLinInternetGatewayLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	vpcID := d.Get("vpc_id").(string)
	igID := d.Get("internet_gateway_id").(string)

	req := &fcu.AttachInternetGatewayInput{
		VpcId:             aws.String(vpcID),
		InternetGatewayId: aws.String(igID),
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.VM.AttachInternetGateway(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		log.Printf("[DEBUG] Error linking Internet Gateway id (%s)", err)
	}

	d.SetId(igID)

	return resourceOutscaleLinInternetGatewayLinkRead(d, meta)
}

func resourceOutscaleLinInternetGatewayLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	id := d.Id()

	d.SetId(id)

	log.Printf("[DEBUG] Reading LIN Internet Gateway id (%s)", id)

	req := &fcu.DescribeInternetGatewaysInput{
		InternetGatewayIds: []*string{aws.String(id)},
	}

	var resp *fcu.DescribeInternetGatewaysOutput
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.VM.DescribeInternetGateways(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		log.Printf("[DEBUG] Error reading LIN Internet Gateway id (%s)", err)
		return err
	}

	if resp == nil {
		d.SetId("")
		return errors.New("Got a nil response for Internet Gateway Link")
	}

	if resp.InternetGateways == nil {
		return errors.New("Failed to retrieve attachments Internet Gateway Link")
	}

	if len(resp.InternetGateways) > 0 {
		attchs := flattenInternetAttachements(resp.InternetGateways[0].Attachments)
		d.Set("attachment_set", attchs)
	}

	if err := d.Set("tag_set", dataSourceTags(resp.InternetGateways[0].Tags)); err != nil {
		return err
	}

	d.Set("internet_gateway_id", resp.InternetGateways[0].InternetGatewayId)
	d.Set("request_id", resp.RequesterId)

	return nil
}

func resourceOutscaleLinInternetGatewayLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	vpcID := d.Get("vpc_id").(string)
	igID := d.Get("internet_gateway_id").(string)

	req := &fcu.DetachInternetGatewayInput{
		VpcId:             aws.String(vpcID),
		InternetGatewayId: aws.String(igID),
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.VM.DetachInternetGateway(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		log.Printf("[DEBUG] Error dettaching Internet Gateway id (%s)", err)
	}

	d.SetId("")

	return nil
}

func getLinInternetGatewayLinkSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Arguments
		"vpc_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"internet_gateway_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		// Attributes
		"attachment_set": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"state": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"vpc_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"tag_set": dataSourceTagsSchema(),
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
