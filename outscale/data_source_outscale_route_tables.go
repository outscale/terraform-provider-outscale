package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/outscale/osc-go/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleOAPIRouteTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIRouteTablesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"route_table_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"net_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": tagsListOAPISchema(),
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
									"net_peering_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"net_access_point_id": {
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
									"main": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"route_table_to_subnet_link_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"link_route_table_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"route_table_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subnet_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceOutscaleOAPIRouteTablesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI
	rtbID, rtbOk := d.GetOk("route_table_id")
	filter, filterOk := d.GetOk("filter")
	if !filterOk && !rtbOk {
		return fmt.Errorf("One of route_table_id or filters must be assigned")
	}

	params := &oapi.ReadRouteTablesRequest{
		Filters: oapi.FiltersRouteTable{},
	}

	if rtbOk {
		i := rtbID.([]interface{})
		in := make([]string, len(i))
		for k, v := range i {
			in[k] = v.(string)
		}
		params.Filters.RouteTableIds = in
	}

	if filterOk {
		params.Filters = buildOutscaleOAPIDataSourceRouteTableFilters(filter.(*schema.Set))
	}

	var resp *oapi.POST_ReadRouteTablesResponses
	var err error
	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		resp, err = conn.POST_ReadRouteTables(*params)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})

	var errString string
	if err != nil || resp.OK == nil {
		if err != nil {
			errString = err.Error()
		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
		}
		return fmt.Errorf("[DEBUG] Error reading Internet Services (%s)", errString)
	}

	if err != nil {
		return err
	}

	rt := resp.OK.RouteTables
	if resp == nil || len(rt) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	routeTables := make([]map[string]interface{}, len(rt))

	for k, v := range rt {
		routeTable := make(map[string]interface{})
		routeTable["route_propagating_virtual_gateways"] = setOAPIPropagatingVirtualGateways(v.RoutePropagatingVirtualGateways)
		routeTable["route_table_id"] = v.RouteTableId
		routeTable["net_id"] = v.NetId
		routeTable["tags"] = tagsOAPIToMap(v.Tags)
		routeTable["routes"] = setOAPIRoutes(v.Routes)
		routeTable["link_route_tables"] = setOAPILinkRouteTables(v.LinkRouteTables)
		routeTables[k] = routeTable
	}

	d.SetId(resource.UniqueId())
	d.Set("request_id", resp.OK.ResponseContext.RequestId)

	return d.Set("route_tables", routeTables)
}
