package outscale

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleVpnConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleVpnConnectionsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vpn_connection_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpn_connection_set": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpn_connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"customer_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"options": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"static_routes_only": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"customer_gateway_configuration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"routes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destination_cidr_block": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source": {
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

						"tag_set": tagsSchemaComputed(),
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vgw_telemetry": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"accepted_route_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"outside_ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status_message": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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

func dataSourceOutscaleVpnConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	vpn, vpnOk := d.GetOk("vpn_connection_id")

	if !filtersOk && !vpnOk {
		return fmt.Errorf("One of vpn_connection_id or filters must be assigned")
	}

	params := &fcu.DescribeVpnConnectionsInput{}

	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if vpnOk {
		var ids []*string

		for _, id := range vpn.([]interface{}) {
			ids = append(ids, aws.String(id.(string)))
		}
		params.VpnConnectionIds = ids
	}

	var resp *fcu.DescribeVpnConnectionsOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeVpnConnections(params)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		return err
	}
	if resp == nil || len(resp.VpnConnections) == 0 {
		return fmt.Errorf("no matching VPN connection found: %#v", params)
	}

	vcs := make([]map[string]interface{}, len(resp.VpnConnections))

	for k, v := range resp.VpnConnections {
		vc := make(map[string]interface{})
		options := make(map[string]interface{})
		if v.Options != nil {
			options["static_routes_only"] = strconv.FormatBool(aws.BoolValue(v.Options.StaticRoutesOnly))
		} else {
			options["static_routes_only"] = strconv.FormatBool(false)
		}
		vc["options"] = options
		vc["customer_gateway_configuration"] = *v.CustomerGatewayConfiguration
		vc["customer_gateway_id"] = *v.CustomerGatewayId

		routes := make([]map[string]interface{}, len(v.Routes))

		for k1, v1 := range v.Routes {
			route := make(map[string]interface{})

			route["destination_cidr_block"] = *v1.DestinationCidrBlock
			route["source"] = *v1.Source
			route["state"] = *v1.State

			routes[k1] = route
		}
		vc["routes"] = routes
		vc["tag_set"] = tagsToMap(v.Tags)
		vc["state"] = *v.State

		vgws := make([]map[string]interface{}, len(v.VgwTelemetry))

		for k1, v1 := range v.VgwTelemetry {
			vgw := make(map[string]interface{})

			vgw["accepted_route_count"] = *v1.AcceptedRouteCount
			vgw["outside_ip_address"] = *v1.OutsideIpAddress
			vgw["status"] = *v1.Status
			vgw["status_message"] = *v1.StatusMessage

			vgws[k1] = vgw
		}
		vc["vgw_telemetry"] = vgws
		vc["vpn_connection_id"] = *v.VpnConnectionId
		vc["vpn_gateway_id"] = *v.VpnGatewayId
		vc["type"] = *v.Type

		vcs[k] = vc
	}

	if err := d.Set("vpn_connection_set", vcs); err != nil {
		return err
	}
	d.Set("request_id", resp.RequestId)
	d.SetId(resource.UniqueId())

	return nil
}
