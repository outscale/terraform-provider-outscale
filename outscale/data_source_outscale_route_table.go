package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/spf13/cast"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleRouteTable() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleRouteTableRead,

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
			"tags": dataSourceTagsSchema(),
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

func setOSCAPIRoutes(rt []oscgo.Route) []map[string]interface{} {
	route := make([]map[string]interface{}, len(rt))
	if len(rt) > 0 {
		for k, r := range rt {
			m := make(map[string]interface{})

			if r.GetNatServiceId() != "" {
				m["nat_service_id"] = r.GetNatServiceId()
			}

			if r.GetCreationMethod() != "" {
				m["creation_method"] = r.GetCreationMethod()
			}
			if r.GetDestinationIpRange() != "" {
				m["destination_ip_range"] = r.GetDestinationIpRange()
			}
			if r.GetDestinationServiceId() != "" {
				m["destination_service_id"] = r.GetDestinationServiceId()
			}
			if r.GetGatewayId() != "" {
				m["gateway_id"] = r.GetGatewayId()
			}
			if r.GetNetAccessPointId() != "" {
				m["net_access_point_id"] = r.GetNetAccessPointId()
			}
			if r.GetNetPeeringId() != "" {
				m["net_peering_id"] = r.GetNetPeeringId()
			}
			if r.GetVmId() != "" {
				m["vm_id"] = r.GetVmId()
			}
			if r.GetNicId() != "" {
				m["nic_id"] = r.GetNicId()
			}
			if r.GetState() != "" {
				m["state"] = r.GetState()
			}
			if r.GetVmAccountId() != "" {
				m["vm_account_id"] = r.GetVmAccountId()
			}
			route[k] = m
		}
	}

	return route
}

func setOSCAPILinkRouteTables(rt []oscgo.LinkRouteTable) []map[string]interface{} {
	linkRouteTables := make([]map[string]interface{}, len(rt))
	log.Printf("[DEBUG] LinkRouteTable: %#v", rt)
	if len(rt) > 0 {
		for k, r := range rt {
			m := make(map[string]interface{})
			if r.GetMain() {
				m["main"] = r.GetMain()
			}
			if r.GetRouteTableId() != "" {
				m["route_table_id"] = r.GetRouteTableId()
			}
			if r.GetLinkRouteTableId() != "" {
				m["link_route_table_id"] = r.GetLinkRouteTableId()
			}
			if r.GetSubnetId() != "" {
				m["subnet_id"] = r.GetSubnetId()
			}
			linkRouteTables[k] = m
		}
	}

	return linkRouteTables
}

func setOSCAPIPropagatingVirtualGateways(vg []oscgo.RoutePropagatingVirtualGateway) (propagatingVGWs []map[string]interface{}) {
	propagatingVGWs = make([]map[string]interface{}, len(vg))

	if len(vg) > 0 {
		for k, vgw := range vg {
			m := make(map[string]interface{})
			if vgw.GetVirtualGatewayId() != "" {
				m["virtual_gateway_id"] = vgw.GetVirtualGatewayId()
			}
			propagatingVGWs[k] = m
		}
	}
	return propagatingVGWs
}

func DataSourceOutscaleRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	routeTableID, routeTableIDOk := d.GetOk("route_table_id")
	filter, filterOk := d.GetOk("filter")

	if !filterOk && !routeTableIDOk {
		return fmt.Errorf("One of route_table_id or filters must be assigned")
	}

	params := oscgo.ReadRouteTablesRequest{}
	if routeTableIDOk {
		params.Filters = &oscgo.FiltersRouteTable{
			RouteTableIds: &[]string{routeTableID.(string)},
		}
	}

	var err error
	if filterOk {
		params.Filters, err = buildOutscaleDataSourceRouteTableFilters(filter.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp oscgo.ReadRouteTablesResponse
	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.RouteTableApi.ReadRouteTables(context.Background()).ReadRouteTablesRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	numRouteTables := len(resp.GetRouteTables())
	if numRouteTables <= 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}
	if numRouteTables > 1 {
		return fmt.Errorf("Multiple Route Table matched; use additional constraints to reduce matches to a single Route Table")
	}

	rt := resp.GetRouteTables()[0]
	if err :=
		d.Set("route_propagating_virtual_gateways", setOSCAPIPropagatingVirtualGateways(rt.GetRoutePropagatingVirtualGateways())); err != nil {
		return err
	}
	if err := d.Set("route_table_id", rt.GetRouteTableId()); err != nil {
		return err
	}
	if err := d.Set("net_id", rt.GetNetId()); err != nil {
		return err
	}
	if err := d.Set("tags", tagsOSCAPIToMap(rt.GetTags())); err != nil {
		return err
	}
	if err := d.Set("routes", setOSCAPIRoutes(rt.GetRoutes())); err != nil {
		return err
	}

	d.SetId(rt.GetRouteTableId())

	return d.Set("link_route_tables", setOSCAPILinkRouteTables(rt.GetLinkRouteTables()))
}

func buildOutscaleDataSourceRouteTableFilters(set *schema.Set) (*oscgo.FiltersRouteTable, error) {
	var filters oscgo.FiltersRouteTable
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}
		switch name := m["name"].(string); name {
		case "route_table_ids":
			filters.SetRouteTableIds(filterValues)
		case "link_route_table_link_route_table_ids":
			filters.SetLinkRouteTableLinkRouteTableIds(filterValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		case "tags":
			filters.SetTags(filterValues)
		case "link_route_table_ids":
			filters.SetLinkRouteTableIds(filterValues)
		case "link_route_table_main":
			filters.SetLinkRouteTableMain(cast.ToBool(filterValues[0]))
		case "link_subnet_ids":
			filters.SetLinkSubnetIds(filterValues)
		case "net_ids":
			filters.SetNetIds(filterValues)
		case "route_creation_methods":
			filters.SetRouteCreationMethods(filterValues)
		case "route_destination_ip_ranges":
			filters.SetRouteDestinationIpRanges(filterValues)
		case "route_destination_service_ids":
			filters.SetRouteDestinationServiceIds(filterValues)
		case "route_gateway_ids":
			filters.SetRouteGatewayIds(filterValues)
		case "route_nat_service_ids":
			filters.SetRouteNatServiceIds(filterValues)
		case "route_net_peering_ids":
			filters.SetRouteNetPeeringIds(filterValues)
		case "route_states":
			filters.SetRouteStates(filterValues)
		case "route_vm_ids":
			filters.SetRouteVmIds(filterValues)
		default:
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}
