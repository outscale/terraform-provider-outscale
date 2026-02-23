package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleRouteTables() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleRouteTablesRead,

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
						"tags": TagsSchemaComputedSDK(),
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

func DataSourceOutscaleRouteTablesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	rtbID, rtbOk := d.GetOk("route_table_id")
	filter, filterOk := d.GetOk("filter")
	if !filterOk && !rtbOk {
		return diag.Errorf("one of route_table_id or filters must be assigned")
	}

	params := osc.ReadRouteTablesRequest{
		Filters: &osc.FiltersRouteTable{},
	}

	if rtbOk {
		i := rtbID.([]interface{})
		in := make([]string, len(i))
		for k, v := range i {
			in[k] = v.(string)
		}
		filter := osc.FiltersRouteTable{}
		filter.RouteTableIds = &in
		params.Filters = &filter
	}

	var err error
	if filterOk {
		params.Filters, err = buildOutscaleDataSourceRouteTableFilters(filter.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadRouteTables(ctx, params, options.WithRetryTimeout(60*time.Second))

	var errString string
	if err != nil {
		errString = err.Error()
		return diag.Errorf("error reading internet services (%s)", errString)
	}

	rt := ptr.From(resp.RouteTables)
	if len(rt) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	routeTables := make([]map[string]interface{}, len(rt))

	for k, v := range rt {
		routeTable := make(map[string]interface{})
		routeTable["route_propagating_virtual_gateways"] = setOSCAPIPropagatingVirtualGateways(v.RoutePropagatingVirtualGateways)
		routeTable["route_table_id"] = v.RouteTableId
		routeTable["net_id"] = v.NetId
		routeTable["tags"] = FlattenOAPITagsSDK(v.Tags)
		routeTable["routes"] = setOSCAPIRoutes(v.Routes)
		routeTable["link_route_tables"] = setOSCAPILinkRouteTables(v.LinkRouteTables)
		routeTables[k] = routeTable
	}

	d.SetId(id.UniqueId())

	return diag.FromErr(d.Set("route_tables", routeTables))
}
