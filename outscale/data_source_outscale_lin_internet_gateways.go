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

func datasourceOutscaleLinInternetGateways() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOutscaleLinInternetGatewaysRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"internet_gateway_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"internet_gateway_set": {
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

func datasourceOutscaleLinInternetGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	internetID, insternetIDOk := d.GetOk("internet_gateway_id")

	if filtersOk == false && insternetIDOk == false {
		return fmt.Errorf("One of filters, or internet_gateway_id must be assigned")
	}

	// Build up search parameters
	params := &fcu.DescribeInternetGatewaysInput{}
	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if insternetIDOk {
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
		log.Printf("[DEBUG] Error reading LIN Internet Gateways id (%s)", err)
	}

	log.Printf("[DEBUG] Setting LIN Internet Gateways id (%s)", err)

	d.Set("request_id", resp.RequesterId)
	d.SetId(resource.UniqueId())

	return internetGatewaysDescriptionAttributes(d, resp.InternetGateways)
}

func flattenInternetGwsAttachements(attachements []*fcu.InternetGatewayAttachment) []map[string]interface{} {
	res := make([]map[string]interface{}, len(attachements))

	for i, a := range attachements {
		res[i]["state"] = a.State
		res[i]["vpc_id"] = a.VpcId
	}

	return res
}

func internetGatewaysDescriptionAttributes(d *schema.ResourceData, internetGateways []*fcu.InternetGateway) error {

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
					at["vpc_id"] = *n.VpcId
				}
				a[m] = at
			}
			im["attachment_set"] = a
		}
		if v.InternetGatewayId != nil {
			im["internet_gateway_id"] = *v.InternetGatewayId
		}
		if v.Tags != nil {
			im["tag_set"] = tagsToMap(v.Tags)
		}
		i[k] = im
	}

	return d.Set("internet_gateway_set", i)
}
