package outscale

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleOAPIVpnConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIVpnConnectionsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vpn_connection_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpn_connection": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpn_connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_connection_option": {
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
						"client_endpoint_configuration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_static_route": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destination_ip_range": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
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

						"tag": tagsSchemaComputed(),
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_tunnel_description": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"accepted_routes_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"outscale_side_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"comment": {
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

func dataSourceOutscaleOAPIVpnConnectionsRead(d *schema.ResourceData, meta interface{}) error {
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
		vpn := make(map[string]interface{})
		if v.Options != nil {
			vpn["static_routes_only"] = strconv.FormatBool(aws.BoolValue(v.Options.StaticRoutesOnly))
		} else {
			vpn["static_routes_only"] = strconv.FormatBool(false)
		}
		vc["vpn_connection_option"] = vpn
		vc["client_endpoint_configuration"] = *v.CustomerGatewayConfiguration
		vc["client_endpoint_id"] = *v.CustomerGatewayId

		vr := make([]map[string]interface{}, len(v.Routes))

		for k1, v1 := range v.Routes {
			route := make(map[string]interface{})

			route["destination_ip_range"] = *v1.DestinationCidrBlock
			route["type"] = *v1.Source
			route["state"] = *v1.State

			vr[k1] = route
		}
		vc["vpn_static_route"] = vr
		vc["tag"] = tagsToMap(v.Tags)
		vc["state"] = *v.State

		vgws := make([]map[string]interface{}, len(v.VgwTelemetry))

		for k1, v1 := range v.VgwTelemetry {
			vgw := make(map[string]interface{})

			vgw["accepted_routes_count"] = *v1.AcceptedRouteCount
			vgw["outscale_side_ip"] = *v1.OutsideIpAddress
			vgw["status"] = *v1.Status
			vgw["comment"] = *v1.StatusMessage

			vgws[k1] = vgw
		}
		vc["vpn_tunnel_description"] = vgws
		vc["vpn_connection_id"] = *v.VpnConnectionId
		vc["vpn_gateway_id"] = *v.VpnGatewayId
		vc["type"] = *v.Type

		vcs[k] = vc
	}

	if err := d.Set("vpn_connection", vcs); err != nil {
		return err
	}
	d.Set("request_id", resp.RequestId)
	d.SetId(resource.UniqueId())

	return nil
}
