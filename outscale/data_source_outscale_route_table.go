package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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

func dataSourceOutscaleOAPIRouteTableRead(d *schema.ResourceData, meta interface{}) error {
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

	if filterOk {
		params.Filters = buildOutscaleOAPIDataSourceRouteTableFilters(filter.(*schema.Set))
	}

	var resp oscgo.ReadRouteTablesResponse
	var err error
	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		resp, _, err = conn.RouteTableApi.ReadRouteTables(context.Background(), &oscgo.ReadRouteTablesOpts{ReadRouteTablesRequest: optional.NewInterface(params)})
		if err != nil && strings.Contains(err.Error(), "RequestLimitExceeded") {
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	numRouteTables := len(resp.GetRouteTables())
	if numRouteTables <= 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}
	if numRouteTables > 1 {
		return fmt.Errorf("Multiple Route Table matched; use additional constraints to reduce matches to a single Route Table")
	}

	rt := resp.GetRouteTables()[0]

	d.Set("route_propagating_virtual_gateways", setOSCAPIPropagatingVirtualGateways(rt.GetRoutePropagatingVirtualGateways()))
	d.SetId(rt.GetRouteTableId())
	d.Set("route_table_id", rt.GetRouteTableId())
	d.Set("net_id", rt.GetNetId())
	d.Set("tags", tagsOSCAPIToMap(rt.GetTags()))
	d.Set("request_id", resp.ResponseContext.GetRequestId())

	if err := d.Set("routes", setOSCAPIRoutes(rt.GetRoutes())); err != nil {
		return err
	}

	return d.Set("link_route_tables", setOSCAPILinkRouteTables(rt.GetLinkRouteTables()))
}

func buildOutscaleOAPIDataSourceRouteTableFilters(set *schema.Set) *oscgo.FiltersRouteTable {
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
		case "link_route_table_ids":
			filters.SetLinkRouteTableLinkRouteTableIds(filterValues)

		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return &filters
}
