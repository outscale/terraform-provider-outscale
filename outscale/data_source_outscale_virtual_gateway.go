package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOutscaleOAPIVirtualGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIVirtualGatewayRead,

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
			"connection_type": {
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
			"tags": dataSourceTagsSchema(),
		},
	}
}

func dataSourceOutscaleOAPIVirtualGatewayRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	virtualId, vpnOk := d.GetOk("virtual_gateway_id")

	if !filtersOk && !vpnOk {
		return fmt.Errorf("One of virtual_gateway_id or filter must be assigned")
	}

	params := oscgo.ReadVirtualGatewaysRequest{}

	if vpnOk {
		params.SetFilters(oscgo.FiltersVirtualGateway{VirtualGatewayIds: &[]string{virtualId.(string)}})
	}

	if filtersOk {
		params.SetFilters(buildOutscaleAPIVirtualGatewayFilters(filters.(*schema.Set)))
	}

	var resp oscgo.ReadVirtualGatewaysResponse
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.VirtualGatewayApi.ReadVirtualGateways(context.Background()).ReadVirtualGatewaysRequest(params).Execute()
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
		return fmt.Errorf("multiple virtual gateways matched; use additional constraints to reduce matches to a single virtual gateway")
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
	d.Set("state", aws.StringValue(vgw.State))
	d.Set("connection_type", vgw.ConnectionType)
	d.Set("tags", tagsOSCAPIToMap(vgw.GetTags()))

	return nil
}

func buildOutscaleAPIVirtualGatewayFilters(set *schema.Set) oscgo.FiltersVirtualGateway {
	var filters oscgo.FiltersVirtualGateway
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
		case "connection_types":
			filters.SetConnectionTypes(filterValues)
		case "link_net_ids":
			filters.SetLinkNetIds(filterValues)
		case "link_states":
			filters.SetLinkStates(filterValues)
		case "virtual_gateway_ids":
			filters.SetVirtualGatewayIds(filterValues)
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
