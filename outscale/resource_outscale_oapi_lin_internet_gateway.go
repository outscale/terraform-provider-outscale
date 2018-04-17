package outscale

import (
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleOAPILinInternetGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILinInternetGatewayCreate,
		Read:   resourceOutscaleOAPILinInternetGatewayRead,
		Delete: resourceOutscaleOAPILinInternetGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: getOAPILinInternetGatewaySchema(),
	}
}

func resourceOutscaleOAPILinInternetGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	log.Println("[DEBUG] Creating LIN Internet Gateway")
	r, err := conn.VM.CreateInternetGateway(nil)
	if err != nil {
		log.Printf("[DEBUG] Error creating LIN Internet Gateway %s", err)

		return err
	}

	d.SetId(*r.InternetGateway.InternetGatewayId)

	return resourceOutscaleOAPILinInternetGatewayRead(d, meta)
}

func resourceOutscaleOAPILinInternetGatewayRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	id := d.Id()

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
	}

	log.Printf("[DEBUG] Setting LIN Internet Gateway id (%s)", err)

	d.Set("request_id", resp.RequestId)
	d.Set("lin_internet_gateway_id", resp.InternetGateways[0].InternetGatewayId)

	err = d.Set("lin_to_lin_internet_gateway_link", flattenOAPIInternetAttachements(resp.InternetGateways[0].Attachments))
	if err != nil {
		return err
	}

	if err := d.Set("tag", dataSourceTags(resp.InternetGateways[0].Tags)); err != nil {
		return err
	}

	return nil
}

func resourceOutscaleOAPILinInternetGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	id := d.Id()
	log.Printf("[DEBUG] Deleting LIN Internet Gateway id (%s)", id)

	req := &fcu.DeleteInternetGatewayInput{
		InternetGatewayId: &id,
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.VM.DeleteInternetGateway(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		log.Printf("[DEBUG] Error deleting LIN Internet Gateway id (%s)", err)
		return err
	}

	return nil
}

func getOAPILinInternetGatewaySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"lin_internet_gateway_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tag": dataSourceTagsSchema(),
	}
}

func flattenOAPIInternetAttachements(attachements []*fcu.InternetGatewayAttachment) []map[string]interface{} {
	res := make([]map[string]interface{}, len(attachements))

	for i, a := range attachements {
		res[i]["state"] = a.State
		res[i]["lin_id"] = a.VpcId
	}

	return res
}
