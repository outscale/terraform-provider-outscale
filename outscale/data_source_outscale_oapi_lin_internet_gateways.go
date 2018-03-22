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

func datasourceOutscaleOAPILinInternetGateways() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOutscaleLinInternetGatewaysRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			"lin_to_lin_internet_gateway_link": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"attachement_set": {
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
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceOutscaleOAPILinInternetGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	internetID, internetIDOk := d.GetOk("lin_internet_gateway_id")

	if filtersOk == false && internetIDOk == false {
		return fmt.Errorf("One of filters, or instance_id must be assigned")
	}

	// Build up search parameters
	params := &fcu.DescribeInternetGatewaysInput{}
	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if internetIDOk {
		i := internetID.([]string)
		in := make([]*string, len(i))
		for k, v := range i {
			in[k] = aws.String(v)
		}
		params.InternetGatewayIds = in
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
		log.Printf("[DEBUG] Error reading OAPI LIN Internet Gateways id (%s)", err)
	}

	log.Printf("[DEBUG] Setting OAPI LIN Internet Gateways id (%s)", err)

	d.Set("request_id", resp.RequesterId)
	d.SetId(resource.UniqueId())

	return internetGatewaysOAPIDescriptionAttributes(d, resp.InternetGateways)
}

func flattenOAPIInternetGwsAttachements(attachements []*fcu.InternetGatewayAttachment) []map[string]interface{} {
	res := make([]map[string]interface{}, len(attachements))

	for i, a := range attachements {
		res[i]["state"] = a.State
		res[i]["lin_id"] = a.VpcId
	}

	return res
}

func internetGatewaysOAPIDescriptionAttributes(d *schema.ResourceData, internetGateways []*fcu.InternetGateway) error {

	i := make([]map[string]interface{}, len(internetGateways))

	for k, v := range internetGateways {
		im := make(map[string]interface{})

		if v.Attachments != nil {
			a := make([]map[string]interface{}, len(v.Attachments))
			for m, n := range v.Attachments {
				at := make(map[string]interface{})
				if n.State != nil {
					at["state"] = *n.State
				}
				if n.VpcId != nil {
					at["lin_id"] = *n.VpcId
				}
				a[m] = at
			}
			im["attachment_set"] = a
		}
		if v.InternetGatewayId != nil {
			im["lin_internet_gateway_id"] = *v.InternetGatewayId
		}
		if v.Tags != nil {
			im["tag_set"] = tagsToMap(v.Tags)
		}
		i[k] = im
	}

	return d.Set("lin_to_lin_internet_gateway_link", i)
}
