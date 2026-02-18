package oapi

import (
	"fmt"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleVPNConnection() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleVPNConnectionRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vpn_clientection_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtual_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"clientection_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"static_routes_only": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"client_gateway_configuration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_ip_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tags": TagsSchemaComputedSDK(),
			"vgw_telemetries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accepted_route_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"last_state_change_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"outside_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_description": {
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
		},
	}
}

func DataSourceOutscaleVPNclientectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	vpnclientectionID, vpnclientectionOk := d.GetOk("vpn_clientection_id")

	if !filtersOk && !vpnclientectionOk {
		return fmt.Errorf("one of filters, or vpn_clientection_id must be assigned")
	}

	params := osc.ReadVpnclientectionsRequest{}

	if vpnclientectionOk {
		params.Filters = &osc.FiltersVpnclientection{
			VpnclientectionIds: &[]string{vpnclientectionID.(string)},
		}
	}

	var err error
	if filtersOk {
		params.Filters, err = buildOutscaleDataSourceVPNclientectionFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp osc.ReadVpnclientectionsResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := client.VpnclientectionApi.ReadVpnclientections(ctx).ReadVpnclientectionsRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	if len(resp.GetVpnclientections()) == 0 {
		return ErrNoResults
	}
	if len(resp.GetVpnclientections()) > 1 {
		return ErrMultipleResults
	}

	vpnclientection := resp.GetVpnclientections()[0]

	if err := d.Set("client_gateway_id", vpnclientection.GetClientGatewayId()); err != nil {
		return err
	}
	if err := d.Set("virtual_gateway_id", vpnclientection.GetVirtualGatewayId()); err != nil {
		return err
	}
	if err := d.Set("clientection_type", vpnclientection.GetclientectionType()); err != nil {
		return err
	}
	if err := d.Set("static_routes_only", vpnclientection.GetStaticRoutesOnly()); err != nil {
		return err
	}
	if err := d.Set("client_gateway_configuration", vpnclientection.GetClientGatewayConfiguration()); err != nil {
		return err
	}
	if err := d.Set("vpn_clientection_id", vpnclientection.GetVpnclientectionId()); err != nil {
		return err
	}
	if err := d.Set("state", vpnclientection.GetState()); err != nil {
		return err
	}
	if err := d.Set("routes", flattenVPNclientection(vpnclientection.GetRoutes())); err != nil {
		return err
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(vpnclientection.Tags)); err != nil {
		return err
	}
	if err := d.Set("vgw_telemetries", flattenVgwTelemetries(vpnclientection.GetVgwTelemetries())); err != nil {
		return err
	}
	d.SetId(vpnclientection.GetVpnclientectionId())

	return nil
}

func buildOutscaleDataSourceVPNclientectionFilters(set *schema.Set) (*osc.FiltersVpnclientection, error) {
	var filters osc.FiltersVpnclientection
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		var filteBgpAsnsValues []int32
		for _, e := range m["values"].([]interface{}) {
			filteBgpAsnsValues = append(filteBgpAsnsValues, cast.ToInt32(e))
		}

		switch name := m["name"].(string); name {
		case "vpn_clientection_ids":
			filters.SetVpnclientectionIds(filterValues)
		case "virtual_gateway_ids":
			filters.SetVirtualGatewayIds(filterValues)
		case "client_gateway_ids":
			filters.SetClientGatewayIds(filterValues)
		case "clientection_types":
			filters.SetclientectionTypes(filterValues)
		case "route_destination_ip_ranges":
			filters.SetRouteDestinationIpRanges(filterValues)
		case "states":
			filters.SetStates(filterValues)
		case "static_routes_only":
			filters.SetStaticRoutesOnly(cast.ToBool(filterValues[0]))
		case "bgp_asns":
			filters.SetBgpAsns(filteBgpAsnsValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		case "tags":
			filters.SetTags(filterValues)
		default:
			return nil, utils.UnknownDataSourceFilterError(ctx, name)
		}
	}
	return &filters, nil
}
