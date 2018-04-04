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

func dataSourceOutscaleOAPICustomerGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPICustomerGatewayRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
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
				Optional: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag": tagsSchemaComputed(),
		},
	}
}

func dataSourceOutscaleOAPICustomerGatewayRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeCustomerGatewaysInput{}

	filters, filtersOk := d.GetOk("filter")
	v, vOk := d.GetOk("client_endpoint_id")

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
	d.Set("public_ip", customerGateway.IpAddress)
	d.Set("type", customerGateway.Type)
	d.Set("tag", tagsToMap(customerGateway.Tags))

	if *customerGateway.BgpAsn != "" {
		val, err := strconv.ParseInt(*customerGateway.BgpAsn, 0, 0)
		if err != nil {
			return fmt.Errorf("error parsing bgp_asn: %s", err)
		}

		d.Set("bgp_asn", int(val))
	}

	return nil
}
