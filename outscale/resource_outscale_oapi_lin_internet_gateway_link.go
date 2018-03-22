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

func resourceOutscaleOAPILinInternetGatewayLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILinInternetGatewayLinkCreate,
		Read:   resourceOutscaleOAPILinInternetGatewayLinkRead,
		Delete: resourceOutscaleOAPILinInternetGatewayLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: getOAPILinInternetGatewayLinkSchema(),
	}
}

func resourceOutscaleOAPILinInternetGatewayLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	vpcID := d.Get("lin_id").(string)
	igID := d.Get("lin_internet_gateway_id").(string)

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

	return resourceOutscaleOAPILinInternetGatewayLinkRead(d, meta)
}

func resourceOutscaleOAPILinInternetGatewayLinkRead(d *schema.ResourceData, meta interface{}) error {
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
		d.Set("lin_to_lin_internet_gateway_link", attchs)
	}

	if err := d.Set("tag", dataSourceTags(resp.InternetGateways[0].Tags)); err != nil {
		return err
	}

	d.Set("lin_internet_gateway_id", resp.InternetGateways[0].InternetGatewayId)
	d.Set("request_id", resp.RequesterId)

	return nil
}

func resourceOutscaleOAPILinInternetGatewayLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	vpcID := d.Get("lin_id").(string)
	igID := d.Get("lin_internet_gateway_id").(string)

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

func getOAPILinInternetGatewayLinkSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Arguments
		"lin_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"lin_internet_gateway_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		// Attributes
		"lin_to_lin_internet_gateway_link": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"state": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"lin_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"tag": dataSourceTagsSchema(),
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
