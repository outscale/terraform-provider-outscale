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

func dataSourceOutscaleCustomerGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleCustomerGatewayRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
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
				Optional: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag_set": tagsSchemaComputed(),
		},
	}
}

func dataSourceOutscaleCustomerGatewayRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeCustomerGatewaysInput{}

	filters, filtersOk := d.GetOk("filter")
	v, vOk := d.GetOk("customer_gateway_id")

	if filtersOk == false && vOk == false {
		return fmt.Errorf("One of filters, or customer_gateway_id must be assigned")
	}

	if filtersOk {
		req.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if vOk {
		req.CustomerGatewayIds = []*string{aws.String(v.(string))}
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
		return fmt.Errorf("Unable to find Customer Gateway")
	}

	if len(resp.CustomerGateways) > 1 {
		return fmt.Errorf("multiple results returned, please use a more specific criteria in your query")
	}

	customerGateway := resp.CustomerGateways[0]
	d.SetId(*customerGateway.CustomerGatewayId)
	d.Set("ip_address", customerGateway.IpAddress)
	d.Set("type", customerGateway.Type)
	d.Set("tag_set", tagsToMap(customerGateway.Tags))

	if *customerGateway.BgpAsn != "" {
		val, err := strconv.ParseInt(*customerGateway.BgpAsn, 0, 0)
		if err != nil {
			return fmt.Errorf("error parsing bgp_asn: %s", err)
		}

		d.Set("bgp_asn", int(val))
	}
	d.Set("request_id", resp.RequesterId)

	return nil
}
