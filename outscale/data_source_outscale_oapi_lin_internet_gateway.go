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

func datasourceOutscaleOAPILinInternetGateway() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOutscaleLinInternetGatewayRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
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
			"tag_set": dataSourceTagsSchema(),
		},
	}
}

func datasourceOutscaleOAPILinInternetGatewayRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	internetID, insternetIDOk := d.GetOk("lin_internet_gateway_id")

	if filtersOk == false && insternetIDOk == false {
		return fmt.Errorf("One of filters, or instance_id must be assigned")
	}

	// Build up search parameters
	params := &fcu.DescribeInternetGatewaysInput{}
	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if insternetIDOk {
		params.InternetGatewayIds = []*string{aws.String(internetID.(string))}
	}

	var resp *fcu.DescribeInternetGatewaysOutput
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.VM.DescribeInternetGateways(params)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		log.Printf("[DEBUG] Error reading OAPI LIN Internet Gateway id (%s)", err)
	}

	log.Printf("[DEBUG] Setting OAPI LIN Internet Gateway id (%s)", err)

	d.Set("request_id", resp.RequestId)
	d.Set("lin_internet_gateway_id", resp.InternetGateways[0].InternetGatewayId)
	d.Set("tag_set", tagsToMap(resp.InternetGateways[0].Tags))

	err = d.Set("lin_to_lin_internet_gateway_link", flattenOAPIInternetGwAttachements(resp.InternetGateways[0].Attachments))
	if err != nil {
		return err
	}

	return d.Set("tag_set", tagsToMap(resp.InternetGateways[0].Tags))
}

func flattenOAPIInternetGwAttachements(attachements []*fcu.InternetGatewayAttachment) []map[string]interface{} {
	res := make([]map[string]interface{}, len(attachements))

	for i, a := range attachements {
		res[i]["state"] = a.State
		res[i]["lin_id"] = a.VpcId
	}

	return res
}
