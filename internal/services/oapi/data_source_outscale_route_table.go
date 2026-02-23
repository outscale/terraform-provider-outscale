package oapi

import (
	"context"
	"log"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleRouteTable() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleRouteTableRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"route_table_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": TagsSchemaComputedSDK(),
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_ip_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_access_point_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_peering_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nic_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"link_route_tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_table_to_subnet_link_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"link_route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"main": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"route_propagating_virtual_gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func setOSCAPIRoutes(rt []osc.Route) []map[string]interface{} {
	route := make([]map[string]interface{}, len(rt))
	if len(rt) > 0 {
		for k, r := range rt {
			m := make(map[string]interface{})

			if ptr.From(r.NatServiceId) != "" {
				m["nat_service_id"] = r.NatServiceId
			}

			if r.CreationMethod != "" {
				m["creation_method"] = r.CreationMethod
			}
			if r.DestinationIpRange != "" {
				m["destination_ip_range"] = r.DestinationIpRange
			}
			if ptr.From(r.DestinationServiceId) != "" {
				m["destination_service_id"] = r.DestinationServiceId
			}
			if ptr.From(r.GatewayId) != "" {
				m["gateway_id"] = r.GatewayId
			}
			if ptr.From(r.NetAccessPointId) != "" {
				m["net_access_point_id"] = r.NetAccessPointId
			}
			if ptr.From(r.NetPeeringId) != "" {
				m["net_peering_id"] = r.NetPeeringId
			}
			if ptr.From(r.VmId) != "" {
				m["vm_id"] = r.VmId
			}
			if ptr.From(r.NicId) != "" {
				m["nic_id"] = r.NicId
			}
			if r.State != "" {
				m["state"] = r.State
			}
			if ptr.From(r.VmAccountId) != "" {
				m["vm_account_id"] = r.VmAccountId
			}
			route[k] = m
		}
	}

	return route
}

func setOSCAPILinkRouteTables(rt []osc.LinkRouteTable) []map[string]interface{} {
	linkRouteTables := make([]map[string]interface{}, len(rt))
	log.Printf("[DEBUG] LinkRouteTable: %#v", rt)
	if len(rt) > 0 {
		for k, r := range rt {
			m := make(map[string]interface{})
			if r.Main {
				m["main"] = r.Main
			}
			if r.RouteTableId != "" {
				m["route_table_id"] = r.RouteTableId
			}
			if r.LinkRouteTableId != "" {
				m["link_route_table_id"] = r.LinkRouteTableId
			}
			if r.SubnetId != "" {
				m["subnet_id"] = r.SubnetId
			}
			linkRouteTables[k] = m
		}
	}

	return linkRouteTables
}

func setOSCAPIPropagatingVirtualGateways(vg []osc.RoutePropagatingVirtualGateway) (propagatingVGWs []map[string]interface{}) {
	propagatingVGWs = make([]map[string]interface{}, len(vg))

	if len(vg) > 0 {
		for k, vgw := range vg {
			m := make(map[string]interface{})
			if ptr.From(vgw.VirtualGatewayId) != "" {
				m["virtual_gateway_id"] = vgw.VirtualGatewayId
			}
			propagatingVGWs[k] = m
		}
	}
	return propagatingVGWs
}

func DataSourceOutscaleRouteTableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	routeTableID, routeTableIDOk := d.GetOk("route_table_id")
	filter, filterOk := d.GetOk("filter")

	if !filterOk && !routeTableIDOk {
		return diag.Errorf("one of route_table_id or filters must be assigned")
	}

	params := osc.ReadRouteTablesRequest{}
	if routeTableIDOk {
		params.Filters = &osc.FiltersRouteTable{
			RouteTableIds: &[]string{routeTableID.(string)},
		}
	}

	var err error
	if filterOk {
		params.Filters, err = buildOutscaleDataSourceRouteTableFilters(filter.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadRouteTables(ctx, params, options.WithRetryTimeout(60*time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	numRouteTables := len(ptr.From(resp.RouteTables))
	if numRouteTables == 0 {
		return diag.FromErr(ErrNoResults)
	}
	if numRouteTables > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	rt := ptr.From(resp.RouteTables)[0]
	if err := d.Set("route_propagating_virtual_gateways", setOSCAPIPropagatingVirtualGateways(rt.RoutePropagatingVirtualGateways)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("route_table_id", rt.RouteTableId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("net_id", rt.NetId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(rt.Tags)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("routes", setOSCAPIRoutes(rt.Routes)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rt.RouteTableId)

	return diag.FromErr(d.Set("link_route_tables", setOSCAPILinkRouteTables(rt.LinkRouteTables)))
}

func buildOutscaleDataSourceRouteTableFilters(set *schema.Set) (*osc.FiltersRouteTable, error) {
	var filters osc.FiltersRouteTable
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}
		switch name := m["name"].(string); name {
		case "route_table_ids":
			filters.RouteTableIds = &filterValues
		case "link_route_table_link_route_table_ids":
			filters.LinkRouteTableLinkRouteTableIds = &filterValues
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		case "tags":
			filters.Tags = &filterValues
		case "link_route_table_ids":
			filters.LinkRouteTableIds = &filterValues
		case "link_route_table_main":
			filters.LinkRouteTableMain = new(cast.ToBool(filterValues[0]))
		case "link_subnet_ids":
			filters.LinkSubnetIds = &filterValues
		case "net_ids":
			filters.NetIds = &filterValues
		case "route_creation_methods":
			filters.RouteCreationMethods = &filterValues
		case "route_destination_ip_ranges":
			filters.RouteDestinationIpRanges = &filterValues
		case "route_destination_service_ids":
			filters.RouteDestinationServiceIds = &filterValues
		case "route_gateway_ids":
			filters.RouteGatewayIds = &filterValues
		case "route_nat_service_ids":
			filters.RouteNatServiceIds = &filterValues
		case "route_net_peering_ids":
			filters.RouteNetPeeringIds = &filterValues
		case "route_states":
			filters.RouteStates = &filterValues
		case "route_vm_ids":
			filters.RouteVmIds = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
