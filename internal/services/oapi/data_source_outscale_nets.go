package oapi

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleVpcs() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleVpcsRead,

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

func DataSourceOutscaleVpcsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	var err error
	req := oscgo.ReadNetsRequest{}

	filters, filtersOk := d.GetOk("filter")
	netIds, netIdsOk := d.GetOk("net_id")

	if !filtersOk && !netIdsOk {
		return fmt.Errorf("filters or net_id(s) must be provided")
	}

	if filtersOk {
		req.Filters, err = buildOutscaleDataSourceNetFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	if netIdsOk {
		ids := make([]string, len(netIds.([]interface{})))

		for k, v := range netIds.([]interface{}) {
			ids[k] = v.(string)
		}
		var filters oscgo.FiltersNet
		filters.SetNetIds(ids)
		req.SetFilters(filters)

	}

	var resp oscgo.ReadNetsResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.NetApi.ReadNets(context.Background()).ReadNetsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}
	if len(resp.GetNets()) == 0 {
		return fmt.Errorf("no matching VPC found")
	}

	d.SetId(id.UniqueId())

	nets := make([]map[string]interface{}, len(resp.GetNets()))

	for i, v := range resp.GetNets() {
		net := make(map[string]interface{})

		net["net_id"] = v.GetNetId()
		net["ip_range"] = v.GetIpRange()
		net["dhcp_options_set_id"] = v.GetDhcpOptionsSetId()
		net["tenancy"] = v.GetTenancy()
		net["state"] = v.GetState()
		if v.Tags != nil {
			net["tags"] = FlattenOAPITagsSDK(v.GetTags())
		}

		nets[i] = net
	}

	if err := d.Set("nets", nets); err != nil {
		return err
	}

	return nil
}
