package oapi

import (
	"context"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleSubnets() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleSubnetsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"subnet_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subregion_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tags": TagsSchemaComputedSDK(),

						"net_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available_ips_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"map_public_ip_on_launch": {
							Type:     schema.TypeBool,
							Computed: true,
						},
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

func DataSourceOutscaleSubnetsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadSubnetsRequest{}

	if id := d.Get("subnet_ids"); id != "" {
		var ids []string
		for _, v := range id.([]interface{}) {
			ids = append(ids, v.(string))
		}
		req.Filters = &osc.FiltersSubnet{SubnetIds: &ids}
	}

	filters, filtersOk := d.GetOk("filter")

	var err error
	if filtersOk {
		req.Filters, err = buildOutscaleSubnetDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadSubnets(ctx, req, options.WithRetryTimeout(120*time.Second))

	var errString string

	if err != nil {
		errString = err.Error()

		return diag.Errorf("error reading subnet (%s)", errString)
	}

	if resp.Subnets == nil || len(*resp.Subnets) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	subnets := make([]map[string]interface{}, len(*resp.Subnets))

	for k, v := range *resp.Subnets {
		subnet := make(map[string]interface{})

		if v.SubregionName != "" {
			subnet["subregion_name"] = v.SubregionName
		}
		// if v.AvailableIpsCount != 0 {
		subnet["available_ips_count"] = v.AvailableIpsCount
		//}
		if v.IpRange != "" {
			subnet["ip_range"] = v.IpRange
		}
		if v.State != "" {
			subnet["state"] = v.State
		}
		if v.SubnetId != "" {
			subnet["subnet_id"] = v.SubnetId
		}
		if v.Tags != nil {
			subnet["tags"] = FlattenOAPITagsSDK(v.Tags)
		}
		if v.NetId != "" {
			subnet["net_id"] = v.NetId
		}
		subnet["map_public_ip_on_launch"] = v.MapPublicIpOnLaunch

		subnets[k] = subnet
	}

	if err := d.Set("subnets", subnets); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id.UniqueId())

	return nil
}
