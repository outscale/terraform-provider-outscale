package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	oscgo "github.com/marinsalinas/osc-sdk-go"
)

func dataSourceOutscaleVPNConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleVPNConnectionsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vpn_connection_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpn_connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpn_connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"static_routes_only": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"client_gateway_configuration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
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
									"route_type": {
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
						"tags": dataSourceTagsSchema(),
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

func dataSourceOutscaleVPNConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	vpnConnectionIDs, vpnConnectionOk := d.GetOk("vpn_connection_ids")

	if !filtersOk && !vpnConnectionOk {
		return fmt.Errorf("One of filters, or vpn_connection_ids must be assigned")
	}

	log.Printf("vpnConnectionIDs: %#+v\n", vpnConnectionIDs)
	params := oscgo.ReadVpnConnectionsRequest{}

	if vpnConnectionOk {
		params.Filters = &oscgo.FiltersVpnConnection{
			VpnConnectionIds: expandStringValueListPointer(vpnConnectionIDs.([]interface{})),
		}
	}

	if filtersOk {
		params.Filters = buildOutscaleDataSourceVPNConnectionFilters(filters.(*schema.Set))
	}

	var resp oscgo.ReadVpnConnectionsResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.VpnConnectionApi.ReadVpnConnections(context.Background(), &oscgo.ReadVpnConnectionsOpts{
			ReadVpnConnectionsRequest: optional.NewInterface(params),
		})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(resp.GetVpnConnections()) == 0 {
		return fmt.Errorf("Unable to find VPN Connections")
	}
	if err := d.Set("vpn_connections", flattenVPNConnections(resp.GetVpnConnections())); err != nil {
		return err
	}
	if err := d.Set("request_id", resp.ResponseContext.GetRequestId()); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenVPNConnections(vpnConnections []oscgo.VpnConnection) []map[string]interface{} {
	vpnConnectionsMap := make([]map[string]interface{}, len(vpnConnections))

	for i, vpnConnection := range vpnConnections {
		vpnConnectionsMap[i] = map[string]interface{}{
			"vpn_connection_id":            vpnConnection.GetVpnConnectionId(),
			"client_gateway_id":            vpnConnection.GetClientGatewayId(),
			"virtual_gateway_id":           vpnConnection.GetVirtualGatewayId(),
			"connection_type":              vpnConnection.GetConnectionType(),
			"static_routes_only":           vpnConnection.GetStaticRoutesOnly(),
			"client_gateway_configuration": vpnConnection.GetClientGatewayConfiguration(),
			"state":                        vpnConnection.GetState(),
			"routes":                       flattenVPNConnection(vpnConnection.GetRoutes()),
			"tags":                         tagsOSCAPIToMap(vpnConnection.GetTags()),
		}
	}
	return vpnConnectionsMap
}
