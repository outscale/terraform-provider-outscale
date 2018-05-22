package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleNatService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleNatServiceRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Attributes
			"nat_gateway_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allocation_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleNatServiceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	natGatewayID, natGatewayIDOK := d.GetOk("nat_gateway_id")

	if filtersOk == false && natGatewayIDOK == false {
		return fmt.Errorf("filters, or owner must be assigned, or nat_gateway_id must be provided")
	}

	params := &fcu.DescribeNatGatewaysInput{}
	if filtersOk {
		params.Filter = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if natGatewayIDOK {
		params.NatGatewayIds = []*string{aws.String(natGatewayID.(string))}
	}

	var err error
	var res *fcu.DescribeNatGatewaysOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error

		res, err = conn.VM.DescribeNatGateways(params)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if len(res.NatGateways) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	if len(res.NatGateways) > 1 {
		return fmt.Errorf("your query returned more than one result, please try a more " +
			"specific search criteria")
	}

	d.Set("request_id", res.RequestId)

	return ngDescriptionAttributes(d, res.NatGateways[0])
}

// populate the numerous fields that the image description returns.
func ngDescriptionAttributes(d *schema.ResourceData, ng *fcu.NatGateway) error {

	d.SetId(*ng.NatGatewayId)
	d.Set("nat_gateway_id", *ng.NatGatewayId)
	d.Set("state", aws.StringValue(ng.State))
	d.Set("subnet_id", aws.StringValue(ng.SubnetId))
	d.Set("vpc_id", aws.StringValue(ng.VpcId))

	addresses := make([]map[string]interface{}, len(ng.NatGatewayAddresses))
	if ng.NatGatewayAddresses != nil {
		for k, v := range ng.NatGatewayAddresses {
			address := make(map[string]interface{})
			if v.AllocationId != nil {
				address["allocation_id"] = *v.AllocationId
			}
			if v.PublicIp != nil {
				address["public_ip"] = *v.PublicIp
			}
			addresses[k] = address
		}
	}

	return d.Set("nat_gateway_address", addresses)
}
