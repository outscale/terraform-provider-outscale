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

func dataSourceOutscaleOAPINatService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPINatServiceRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"nat_service_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Attributes
			"public_ip": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reservation_id": {
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
			"lin_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPINatServiceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	natGatewayID, natGatewayIDOK := d.GetOk("nat_service_id")

	if filtersOk == false && natGatewayIDOK == false {
		return fmt.Errorf("filters, or owner must be assigned, or nat_service_id must be provided")
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
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	if len(res.NatGateways) > 1 {
		return fmt.Errorf("Your query returned more than one result. Please try a more " +
			"specific search criteria.")
	}

	return ngOAPIDescriptionAttributes(d, res.NatGateways[0])
}

// populate the numerous fields that the image description returns.
func ngOAPIDescriptionAttributes(d *schema.ResourceData, ng *fcu.NatGateway) error {

	d.SetId(*ng.NatGatewayId)
	d.Set("nat_service_id", *ng.NatGatewayId)

	if ng.State != nil {
		d.Set("state", *ng.State)
	}
	if ng.SubnetId != nil {
		d.Set("subnet_id", *ng.SubnetId)
	}
	if ng.VpcId != nil {
		d.Set("lin_id", *ng.VpcId)
	}

	if ng.NatGatewayAddresses != nil {
		addresses := make([]map[string]interface{}, len(ng.NatGatewayAddresses))

		for k, v := range ng.NatGatewayAddresses {
			address := make(map[string]interface{})
			if v.AllocationId != nil {
				address["reservation_id"] = *v.AllocationId
			}
			if v.PublicIp != nil {
				address["public_ip"] = *v.PublicIp
			}
			addresses[k] = address
		}
		if err := d.Set("public_ip", addresses); err != nil {
			return err
		}
	}
	return nil
}
