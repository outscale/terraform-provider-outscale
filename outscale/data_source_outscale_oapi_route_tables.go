package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
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
			"route_table": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lin_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag": tagsSchemaComputed(),
						"route": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destination_ip_range": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"destinaton_prefix_list_id": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"vpn_gateway_id": {
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

									"lin_peering_id": {
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
						"link": {
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
						"route_propagating_vpn_gateway": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpn_gateway_id": {
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
	conn := meta.(*OutscaleClient).FCU
	req := &fcu.DescribeRouteTablesInput{}
	rtbID, rtbOk := d.GetOk("route_table_id")
	filter, filterOk := d.GetOk("filter")

	if !filterOk && !rtbOk {
		return fmt.Errorf("One of route_table_id or filters must be assigned")
	}

	if rtbOk {
		var ids []*string
		for _, v := range rtbID.([]interface{}) {
			ids = append(ids, aws.String(v.(string)))
		}

		req.RouteTableIds = ids
	}

	if filterOk {
		req.Filters = buildOutscaleDataSourceFilters(filter.(*schema.Set))
	}

	var resp *fcu.DescribeRouteTablesOutput
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.VM.DescribeRouteTables(req)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})

	if err != nil {
		return err
	}
	if resp == nil || len(resp.RouteTables) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	routeTables := make([]map[string]interface{}, len(resp.RouteTables))

	for k, v := range resp.RouteTables {
		routeTable := make(map[string]interface{})

		propagatingVGWs := make([]string, 0, len(v.PropagatingVgws))
		for _, vgw := range v.PropagatingVgws {
			propagatingVGWs = append(propagatingVGWs, *vgw.GatewayId)
		}
		routeTable["route_propagating_vpn_gateway"] = propagatingVGWs

		routeTable["route_table_id"] = *v.RouteTableId

		routeTable["lin_id"] = *v.VpcId

		routeTable["tag"] = tagsToMap(v.Tags)

		// routeTable["routes"] = setOAPIRoutes(v.Routes)

		// routeTable["link_route_tables"] = setOAPIAssociactionSet(v.Associations) // FIXME

		routeTables[k] = routeTable
	}

	d.SetId(resource.UniqueId())

	return d.Set("route_table", routeTables)
}
