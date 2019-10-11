package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func dataSourceOutscaleOAPIRouteTable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIRouteTableRead,

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
			"tags": tagsListOAPISchema(),
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_ip_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_prefix_list_id": {
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
						"nat_service_id": {
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

func dataSourceOutscaleOAPIRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI
	routeTableID, routeTableIDOk := d.GetOk("route_table_id")
	filter, filterOk := d.GetOk("filter")

	if !filterOk && !routeTableIDOk {
		return fmt.Errorf("One of route_table_id or filters must be assigned")
	}

	params := &oapi.ReadRouteTablesRequest{}
	if routeTableIDOk {
		params.Filters = oapi.FiltersRouteTable{
			RouteTableIds: []string{routeTableID.(string)},
		}
	}

	if filterOk {
		params.Filters = buildOutscaleOAPIDataSourceRouteTableFilters(filter.(*schema.Set))
	}

	var resp *oapi.POST_ReadRouteTablesResponses
	var err error
	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		resp, err = conn.POST_ReadRouteTables(*params)
		if err != nil && strings.Contains(err.Error(), "RequestLimitExceeded") {
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	numRouteTables := len(resp.OK.RouteTables)
	if numRouteTables <= 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}
	if numRouteTables > 1 {
		return fmt.Errorf("Multiple Route Table matched; use additional constraints to reduce matches to a single Route Table")
	}

	rt := resp.OK.RouteTables[0]

	d.Set("route_propagating_virtual_gateways", setOAPIPropagatingVirtualGateways(rt.RoutePropagatingVirtualGateways))
	d.SetId(rt.RouteTableId)
	d.Set("route_table_id", rt.RouteTableId)
	d.Set("net_id", rt.NetId)
	d.Set("tags", tagsOAPIToMap(rt.Tags))
	d.Set("request_id", resp.OK.ResponseContext.RequestId)

	if err := d.Set("routes", setOAPIRoutes(rt.Routes)); err != nil {
		return err
	}

	return d.Set("link_route_tables", setOAPILinkRouteTables(rt.LinkRouteTables))
}

func buildOutscaleOAPIDataSourceRouteTableFilters(set *schema.Set) oapi.FiltersRouteTable {
	var filters oapi.FiltersRouteTable
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}
		switch name := m["name"].(string); name {
		case "route_table_ids":
			filters.RouteTableIds = filterValues
		case "link_route_table_ids":
			filters.LinkRouteTableLinkRouteTableIds = filterValues

		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
