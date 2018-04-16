package outscale

import (
	"encoding/xml"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

type XmlVpnConnectionConfig struct {
	Tunnels []XmlIpsecTunnel `xml:"ipsec_tunnel"`
}

type XmlIpsecTunnel struct {
	OutsideAddress string `xml:"vpn_gateway>tunnel_outside_address>ip_address"`
	PreSharedKey   string `xml:"ike>pre_shared_key"`
}

type TunnelInfo struct {
	Tunnel1Address      string
	Tunnel1PreSharedKey string
	Tunnel2Address      string
	Tunnel2PreSharedKey string
}

func (slice XmlVpnConnectionConfig) Len() int {
	return len(slice.Tunnels)
}

func (slice XmlVpnConnectionConfig) Less(i, j int) bool {
	return slice.Tunnels[i].OutsideAddress < slice.Tunnels[j].OutsideAddress
}

func (slice XmlVpnConnectionConfig) Swap(i, j int) {
	slice.Tunnels[i], slice.Tunnels[j] = slice.Tunnels[j], slice.Tunnels[i]
}

func resourceOutscaleVpnConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleVpnConnectionCreate,
		Read:   resourceOutscaleVpnConnectionRead,
		Delete: resourceOutscaleVpnConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// Argumentos
			"customer_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"options": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"static_routes_only": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// Atributos

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
			"tag":     tagsSchema(),

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

			"vpn_connection_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleVpnConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	cgid, ok := d.GetOk("customer_gateway_id")
	vpngid, ok2 := d.GetOk("vpn_gateway_id")

	if !ok && !ok2 {
		return fmt.Errorf("please provide the required attributes customer_gateway_id and vpn_gateway_id")
	}

	createOpts := &fcu.CreateVpnConnectionInput{
		CustomerGatewayId: aws.String(cgid.(string)),
		Type:              aws.String("ipsec.1"),
		VpnGatewayId:      aws.String(vpngid.(string)),
	}

	var resp *fcu.CreateVpnConnectionOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.CreateVpnConnection(createOpts)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("Error creating vpn connection: %s", err)
	}

	vpnConnection := resp.VpnConnection
	d.SetId(*vpnConnection.VpnConnectionId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available"},
		Refresh:    vpnConnectionRefreshFunc(conn, *vpnConnection.VpnConnectionId),
		Timeout:    30 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for VPN connection (%s) to become ready: %s",
			*vpnConnection.VpnConnectionId, err)
	}

	if err := setTags(conn, d); err != nil {
		return err
	}

	return resourceOutscaleVpnConnectionRead(d, meta)
}

func vpnConnectionRefreshFunc(conn *fcu.Client, connectionId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := conn.VM.DescribeVpnConnections(&fcu.DescribeVpnConnectionsInput{
			VpnConnectionIds: []*string{aws.String(connectionId)},
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidVpnConnectionID.NotFound") {
				resp = nil
			} else {
				log.Printf("Error on VPNConnectionRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil || len(resp.VpnConnections) == 0 {
			return nil, "", nil
		}

		connection := resp.VpnConnections[0]
		return connection, *connection.State, nil
	}
}

func resourceOutscaleVpnConnectionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var resp *fcu.DescribeVpnConnectionsOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeVpnConnections(&fcu.DescribeVpnConnectionsInput{
			VpnConnectionIds: []*string{aws.String(d.Id())},
		})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidVpnConnectionID.NotFound") {
			d.SetId("")
			return nil
		} else {
			log.Printf("[ERROR] Error finding VPN connection: %s", err)
			return err
		}
	}

	if len(resp.VpnConnections) != 1 {
		return fmt.Errorf("[ERROR] Error finding VPN connection: %s", d.Id())
	}

	vpnConnection := resp.VpnConnections[0]
	if vpnConnection == nil || *vpnConnection.State == "deleted" {
		d.SetId("")
		return nil
	}

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
	d.Set("request_id", resp.RequestId)

	return nil
}

func resourceOutscaleVpnConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	// Update tags if required.
	if err := setTags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

	return resourceOutscaleVpnConnectionRead(d, meta)
}

func resourceOutscaleVpnConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	_, err := conn.DeleteVpnConnection(&fcu.DeleteVpnConnectionInput{
		VpnConnectionId: aws.String(d.Id()),
	})
	if err != nil {
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidVpnConnectionID.NotFound" {
			d.SetId("")
			return nil
		} else {
			log.Printf("[ERROR] Error deleting VPN connection: %s", err)
			return err
		}
	}

	// These things can take quite a while to tear themselves down and any
	// attempt to modify resources they reference (e.g. CustomerGateways or
	// VPN Gateways) before deletion will result in an error. Furthermore,
	// they don't just disappear. The go into "deleted" state. We need to
	// wait to ensure any other modifications the user might make to their
	// VPC stack can safely run.
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{"deleted"},
		Refresh:    vpnConnectionRefreshFunc(conn, d.Id()),
		Timeout:    30 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for VPN connection (%s) to delete: %s", d.Id(), err)
	}

	return nil
}

// routesToMapList turns the list of routes into a list of maps.
func routesToMapList(routes []*fcu.VpnStaticRoute) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(routes))
	for _, r := range routes {
		staticRoute := make(map[string]interface{})
		staticRoute["destination_cidr_block"] = *r.DestinationCidrBlock
		staticRoute["state"] = *r.State

		if r.Source != nil {
			staticRoute["source"] = *r.Source
		}

		result = append(result, staticRoute)
	}

	return result
}

// telemetryToMapList turns the VGW telemetry into a list of maps.
func telemetryToMapList(telemetry []*fcu.VgwTelemetry) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(telemetry))
	for _, t := range telemetry {
		vgw := make(map[string]interface{})
		vgw["accepted_route_count"] = *t.AcceptedRouteCount
		vgw["outside_ip_address"] = *t.OutsideIpAddress
		vgw["status"] = *t.Status
		vgw["status_message"] = *t.StatusMessage

		// LastStatusChange is a time.Time(). Convert it into a string
		// so it can be handled by schema's type system.
		vgw["last_status_change"] = t.LastStatusChange.String()
		result = append(result, vgw)
	}

	return result
}

func xmlConfigToTunnelInfo(xmlConfig string) (*TunnelInfo, error) {
	var vpnConfig XmlVpnConnectionConfig
	if err := xml.Unmarshal([]byte(xmlConfig), &vpnConfig); err != nil {
		return nil, errwrap.Wrapf("Error Unmarshalling XML: {{err}}", err)
	}

	// don't expect consistent ordering from the XML
	sort.Sort(vpnConfig)

	tunnelInfo := TunnelInfo{
		Tunnel1Address:      vpnConfig.Tunnels[0].OutsideAddress,
		Tunnel1PreSharedKey: vpnConfig.Tunnels[0].PreSharedKey,

		Tunnel2Address:      vpnConfig.Tunnels[1].OutsideAddress,
		Tunnel2PreSharedKey: vpnConfig.Tunnels[1].PreSharedKey,
	}

	return &tunnelInfo, nil
}
