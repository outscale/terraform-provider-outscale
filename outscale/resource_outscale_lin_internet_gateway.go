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

func resourceOutscaleLinInternetGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleLinInternetGatewayCreate,
		Read:   resourceOutscaleLinInternetGatewayRead,
		Delete: resourceOutscaleLinInternetGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: getLinInternetGatewaySchema(),
	}
}

func resourceOutscaleLinInternetGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	log.Println("[DEBUG] Creating LIN Internet Gateway")
	r, err := conn.VM.CreateInternetGateway(nil)
	if err != nil {
		log.Printf("[DEBUG] Error creating LIN Internet Gateway %s", err)

		return err
	}

	d.SetId(*r.InternetGateway.InternetGatewayId)

	return resourceOutscaleLinInternetGatewayRead(d, meta)
}

func resourceOutscaleLinInternetGatewayRead(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}

	d.SetId(*resp.InternetGateways[0].InternetGatewayId)
	d.Set("request_id", resp.RequesterId)
	d.Set("internet_gateway_id", resp.InternetGateways[0].InternetGatewayId)

	if err := d.Set("attachement_set", flattenInternetAttachements(resp.InternetGateways[0].Attachments)); err != nil {
		return err
	}

	return d.Set("tag_set", tagsToMap(resp.InternetGateways[0].Tags))
}

func resourceOutscaleLinInternetGatewayDelete(d *schema.ResourceData, meta interface{}) error {
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

func getLinInternetGatewaySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"attachement_set": {
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
		"internet_gateway_id": {
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

func flattenInternetAttachements(attachements []*fcu.InternetGatewayAttachment) []map[string]interface{} {
	res := make([]map[string]interface{}, len(attachements))

	for i, a := range attachements {
		res[i] = map[string]interface{}{
			"state":  *a.State,
			"vpc_id": *a.VpcId,
		}
	}

	return res
}
