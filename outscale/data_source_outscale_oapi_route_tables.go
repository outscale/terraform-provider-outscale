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
	rtbId, rtbOk := d.GetOk("route_table_id")
	filter, filterOk := d.GetOk("filter")

	if !filterOk && !rtbOk {
		return fmt.Errorf("One of route_table_id or filters must be assigned")
	}

	if rtbOk {
		var ids []*string
		for _, v := range rtbId.([]interface{}) {
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
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	route_tables := make([]map[string]interface{}, len(resp.RouteTables))

	for k, v := range resp.RouteTables {
		route_table := make(map[string]interface{})

		propagatingVGWs := make([]string, 0, len(v.PropagatingVgws))
		for _, vgw := range v.PropagatingVgws {
			propagatingVGWs = append(propagatingVGWs, *vgw.GatewayId)
		}
		route_table["route_propagating_vpn_gateway"] = propagatingVGWs

		route_table["route_table_id"] = *v.RouteTableId

		route_table["lin_id"] = *v.VpcId

		route_table["tag"] = tagsToMap(v.Tags)

		route_table["route"] = setOAPIRouteSet(v.Routes)

		route_table["link"] = setOAPIAssociactionSet(v.Associations)

		route_tables[k] = route_table
	}

	d.SetId(resource.UniqueId())

	return d.Set("route_table", route_tables)
}
