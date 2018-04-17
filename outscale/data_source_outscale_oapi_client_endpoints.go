package outscale

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOutscaleOAPICustomerGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPICustomerGatewaysRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"client_endpoint_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_endpoint": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bgp_asn": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"client_endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_set": tagsSchemaComputed(),
					},
				},
			},
		},
	}
}

func dataSourceOutscaleOAPICustomerGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeCustomerGatewaysInput{}

	filters, filtersOk := d.GetOk("filter")
	v, vOk := d.GetOk("client_endpoint_id")

	if filtersOk {
		req.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if vOk {
		var g []*string
		for _, s := range v.([]interface{}) {
			g = append(g, aws.String(s.(string)))
		}
		req.CustomerGatewayIds = g
	}

	var resp *fcu.DescribeCustomerGatewaysOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeCustomerGateways(req)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidCustomerGatewayID.NotFound") {
			d.SetId("")
			return nil
		} else {
			fmt.Printf("[ERROR] Error finding CustomerGateway: %s", err)
			return err
		}
	}

	if len(resp.CustomerGateways) == 0 {
		return fmt.Errorf("Unable to find Customer Gateways")
	}

	customerGateways := make([]map[string]interface{}, len(resp.CustomerGateways))

	for k, v := range resp.CustomerGateways {
		customerGateway := make(map[string]interface{})

		customerGateway["client_endpoint_id"] = *v.CustomerGatewayId
		customerGateway["public_ip"] = *v.IpAddress
		customerGateway["type"] = *v.Type
		customerGateway["tag_set"] = tagsToMap(v.Tags)

		if *v.BgpAsn != "" {
			val, err := strconv.ParseInt(*v.BgpAsn, 0, 0)
			if err != nil {
				return fmt.Errorf("error parsing bgp_asn: %s", err)
			}
			customerGateway["bgp_asn"] = int(val)
		}

		customerGateways[k] = customerGateway
	}

	d.Set("client_endpoint", customerGateways)
	d.Set("request_id", resp.RequestId)
	d.SetId(resource.UniqueId())

	return nil
}
