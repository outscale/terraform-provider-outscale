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

func DataSourceOutscaleVpcs() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleVpcsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"net_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_range": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"dhcp_options_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"net_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tenancy": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": TagsSchemaComputedSDK(),
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

func DataSourceOutscaleVpcsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	var err error
	req := osc.ReadNetsRequest{}

	filters, filtersOk := d.GetOk("filter")
	netIds, netIdsOk := d.GetOk("net_id")

	if !filtersOk && !netIdsOk {
		return diag.Errorf("filters or net_id(s) must be provided")
	}

	if filtersOk {
		req.Filters, err = buildOutscaleDataSourceNetFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if netIdsOk {
		ids := make([]string, len(netIds.([]interface{})))

		for k, v := range netIds.([]interface{}) {
			ids[k] = v.(string)
		}
		var filters osc.FiltersNet
		filters.NetIds = &ids
		req.Filters = &filters

	}

	resp, err := client.ReadNets(ctx, req, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.Nets == nil || len(*resp.Nets) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	d.SetId(id.UniqueId())

	nets := make([]map[string]interface{}, len(*resp.Nets))

	for i, v := range *resp.Nets {
		net := make(map[string]interface{})

		net["net_id"] = v.NetId
		net["ip_range"] = v.IpRange
		net["dhcp_options_set_id"] = v.DhcpOptionsSetId
		net["tenancy"] = v.Tenancy
		net["state"] = v.State
		if v.Tags != nil {
			net["tags"] = FlattenOAPITagsSDK(v.Tags)
		}

		nets[i] = net
	}

	if err := d.Set("nets", nets); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
