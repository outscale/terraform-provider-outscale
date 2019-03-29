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

func dataSourceOutscaleCustomerGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleCustomerGatewaysRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"customer_gateway_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_gateway_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bgp_asn": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"customer_gateway_id": {
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

func dataSourceOutscaleCustomerGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeCustomerGatewaysInput{}

	filters, filtersOk := d.GetOk("filter")
	v, vOk := d.GetOk("customer_gateway_id")

	if filtersOk == false && vOk == false {
		return fmt.Errorf("One of filters, or customer_gateway_id(s) must be assigned")
	}

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
		}
		fmt.Printf("[ERROR] Error finding CustomerGateway: %s", err)
		return err
	}

	if len(resp.CustomerGateways) == 0 {
		return fmt.Errorf("Unable to find Customer Gateways")
	}

	customerGateways := make([]map[string]interface{}, len(resp.CustomerGateways))

	for k, v := range resp.CustomerGateways {
		customerGateway := make(map[string]interface{})

		customerGateway["customer_gateway_id"] = *v.CustomerGatewayId
		customerGateway["ip_address"] = *v.IpAddress
		customerGateway["type"] = *v.Type
		customerGateway["tag_set"] = tagsToMap(v.Tags)
		customerGateway["state"] = *v.State

		if *v.BgpAsn != "" {
			val, err := strconv.ParseInt(*v.BgpAsn, 0, 0)
			if err != nil {
				return fmt.Errorf("error parsing bgp_asn: %s", err)
			}
			customerGateway["bgp_asn"] = int(val)
		}

		customerGateways[k] = customerGateway
	}

	d.Set("customer_gateway_set", customerGateways)
	d.Set("request_id", resp.RequestId)
	d.SetId(resource.UniqueId())

	return nil
}
