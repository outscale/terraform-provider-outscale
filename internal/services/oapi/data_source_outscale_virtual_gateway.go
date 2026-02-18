package oapi

import (
	"context"
	"fmt"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleVirtualGateway() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleVirtualGatewayRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"virtual_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"clientection_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"net_to_virtual_gateway_links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": TagsSchemaComputedSDK(),
		},
	}
}

func DataSourceOutscaleVirtualGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.OutscaleClient).OSC
	

	filters, filtersOk := d.GetOk("filter")
	virtualId, vpnOk := d.GetOk("virtual_gateway_id")

	if !filtersOk && !vpnOk {
		return fmt.Errorf("one of virtual_gateway_id or filter must be assigned")
	}

	params := osc.ReadVirtualGatewaysRequest{}

	if vpnOk {
		params.SetFilters(osc.FiltersVirtualGateway{VirtualGatewayIds: &[]string{virtualId.(string)}})
	}

	var err error
	if filtersOk {
		params.Filters, err = buildOutscaleAPIVirtualGatewayFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp osc.ReadVirtualGatewaysResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := client.VirtualGatewayApi.ReadVirtualGateways(ctx).ReadVirtualGatewaysRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}
	if resp.GetVirtualGateways() == nil || len(resp.GetVirtualGateways()) == 0 {
		return fmt.Errorf("no matching virtual gateway found: %#v", params)
	}
	if len(resp.GetVirtualGateways()) > 1 {
		return ErrMultipleResults
	}

	vgw := resp.GetVirtualGateways()[0]

	d.SetId(vgw.GetVirtualGatewayId())
	vs := make([]map[string]interface{}, len(vgw.GetNetToVirtualGatewayLinks()))

	for k, v := range vgw.GetNetToVirtualGatewayLinks() {
		vp := make(map[string]interface{})

		vp["state"] = v.GetState()
		vp["net_id"] = v.GetNetId()

		vs[k] = vp
	}

	d.Set("net_to_virtual_gateway_links", vs)
	d.Set("state", ptr.From(vgw.State))
	d.Set("clientection_type", vgw.clientectionType)
	d.Set("tags", FlattenOAPITagsSDK(vgw.Tags))

	return nil
}

func buildOutscaleAPIVirtualGatewayFilters(set *schema.Set) (*osc.FiltersVirtualGateway, error) {
	var filters osc.FiltersVirtualGateway
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}
		switch name := m["name"].(string); name {
		// case "available_ips_counts":
		// 	filters.AvailableIpsCounts = filterValues
		case "tags":
			filters.SetTags(filterValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		case "states":
			filters.SetStates(filterValues)
		case "clientection_types":
			filters.SetclientectionTypes(filterValues)
		case "link_net_ids":
			filters.SetLinkNetIds(filterValues)
		case "link_states":
			filters.SetLinkStates(filterValues)
		case "virtual_gateway_ids":
			filters.SetVirtualGatewayIds(filterValues)
		default:
			return nil, utils.UnknownDataSourceFilterError(ctx, name)
		}
	}
	return &filters, nil
}
