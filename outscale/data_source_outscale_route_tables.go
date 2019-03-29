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

func dataSourceOutscaleRouteTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleRouteTablesRead,

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
			"route_table_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_set": tagsSchemaComputed(),
						"route_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destination_cidr_block": {
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

									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_owner_id": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"vpc_peering_connection_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"origin": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"state": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"network_interface_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"association_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"route_table_association_id": {
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
						"propagating_vgw_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gateway_id": {
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

func dataSourceOutscaleRouteTablesRead(d *schema.ResourceData, meta interface{}) error {
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

	rts := make([]map[string]interface{}, len(resp.RouteTables))

	for k, v := range resp.RouteTables {
		rt := make(map[string]interface{})

		propagatingVGWs := make([]string, 0, len(v.PropagatingVgws))
		for _, vgw := range v.PropagatingVgws {
			propagatingVGWs = append(propagatingVGWs, *vgw.GatewayId)
		}
		rt["propagating_vgw_set"] = propagatingVGWs

		rt["route_table_id"] = *v.RouteTableId

		rt["vpc_id"] = *v.VpcId

		rt["tag_set"] = tagsToMap(v.Tags)

		rt["route_set"] = setRouteSet(v.Routes)

		rt["association_set"] = setAssociactionSet(v.Associations)

		rts[k] = rt
	}

	d.Set("request_id", resp.RequestId)

	d.SetId(resource.UniqueId())

	return d.Set("route_table_set", rts)
}
