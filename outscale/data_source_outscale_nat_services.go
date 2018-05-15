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

func dataSourceOutscaleNatServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleNatServicesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"nat_gateway_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// Attributes
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nat_gateway": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"nat_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
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
		},
	}
}

func dataSourceOutscaleNatServicesRead(d *schema.ResourceData, meta interface{}) error {
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
		ids := make([]*string, len(natGatewayID.([]interface{})))

		for k, v := range natGatewayID.([]interface{}) {
			ids[k] = aws.String(v.(string))
		}

		params.NatGatewayIds = ids
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

	d.Set("request_id", res.RequestId)

	return ngsDescriptionAttributes(d, res.NatGateways)
}

// populate the numerous fields that the image description returns.
func ngsDescriptionAttributes(d *schema.ResourceData, ngs []*fcu.NatGateway) error {

	d.SetId(resource.UniqueId())

	addngs := make([]map[string]interface{}, len(ngs))

	for k, v := range ngs {
		addng := make(map[string]interface{})

		if v.NatGatewayAddresses != nil {
			ngas := make([]interface{}, len(v.NatGatewayAddresses))

			for i, w := range v.NatGatewayAddresses {
				nga := make(map[string]interface{})
				if w.AllocationId != nil {
					nga["allocation_id"] = *w.AllocationId
				}
				if w.PublicIp != nil {
					nga["public_ip"] = *w.PublicIp
				}
				ngas[i] = nga
			}
			addng["nat_gateway_address"] = ngas
		}

		if v.NatGatewayId != nil {
			addng["nat_gateway_id"] = *v.NatGatewayId
		}
		if v.State != nil {
			addng["state"] = *v.State
		}
		if v.SubnetId != nil {
			addng["subnet_id"] = *v.SubnetId
		}
		if v.VpcId != nil {
			addng["vpc_id"] = *v.VpcId
		}

		addngs[k] = addng
	}

	return d.Set("nat_gateway", addngs)
}
