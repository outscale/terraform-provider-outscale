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

func dataSourceOutscaleVpnConnection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleVpnConnectionRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vpn_connection_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"options": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"static_routes_only": {
							Type:     schema.TypeBool,
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleVpnConnectionRead(d *schema.ResourceData, meta interface{}) error {
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
		params.VpnConnectionIds = []*string{aws.String(vpn.(string))}
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
	if len(resp.VpnConnections) > 1 {
		return fmt.Errorf("multiple VPN connections matched; use additional constraints to reduce matches to a single VPN connection")
	}

	vpnConnection := resp.VpnConnections[0]

	options := map[string]interface{}{
		"static_routes_only": vpnConnection.Options.StaticRoutesOnly,
	}

	d.Set("options", options)
	d.Set("customer_gateway_configuration", vpnConnection.CustomerGatewayConfiguration)

	routes := make([]map[string]interface{}, len(vpnConnection.Routes))

	for k, v := range vpnConnection.Routes {
		route := make(map[string]interface{})

		route["destination_cidr_block"] = *v.DestinationCidrBlock
		route["source"] = *v.Source
		route["state"] = *v.State

		routes[k] = route
	}

	d.Set("routes", routes)
	d.Set("tag_set", tagsToMap(vpnConnection.Tags))

	d.Set("state", vpnConnection.State)

	vgws := make([]map[string]interface{}, len(vpnConnection.VgwTelemetry))

	for k, v := range vpnConnection.VgwTelemetry {
		vgw := make(map[string]interface{})

		vgw["accepted_route_count"] = *v.AcceptedRouteCount
		vgw["outside_ip_address"] = *v.OutsideIpAddress
		vgw["status"] = *v.Status
		vgw["status_message"] = *v.StatusMessage

		vgws[k] = vgw
	}

	d.Set("vgw_telemetry", vgws)
	d.Set("vpn_connection_id", vpnConnection.VpnConnectionId)
	d.Set("vpn_gateway_id", vpnConnection.VpnGatewayId)
	d.Set("type", vpnConnection.Type)
	d.Set("request_id", resp.RequestId)
	d.SetId(resource.UniqueId())

	return nil
}
