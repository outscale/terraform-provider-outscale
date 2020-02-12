package outscale

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleOAPIVpnConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIVpnConnectionCreate,
		Read:   resourceOutscaleOAPIVpnConnectionRead,
		Delete: resourceOutscaleOAPIVpnConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// Argumentos
			"client_endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpn_connection_option": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"static_vpn_static_route_only": {
							Type:     schema.TypeString,
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

			"tag":  tagsSchemaComputed(),
			"tags": tagsSchema(),

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
						"state": {
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

func resourceOutscaleOAPIVpnConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	cgid, ok := d.GetOk("client_endpoint_id")
	vpngid, ok2 := d.GetOk("vpn_gateway_id")
	typev, ok3 := d.GetOk("type")
	vpn, ok4 := d.GetOk("vpn_connection_option")

	if !ok && !ok2 && ok3 {
		return fmt.Errorf("please provide the required attributes client_endpoint_id, vpn_gateway_id and type")
	}

	createOpts := &fcu.CreateVpnConnectionInput{
		CustomerGatewayId: aws.String(cgid.(string)),
		Type:              aws.String(typev.(string)),
		VpnGatewayId:      aws.String(vpngid.(string)),
	}

	if ok4 {
		opt := vpn.(map[string]interface{})
		option := opt["static_vpn_static_route_only"]

		b, err := strconv.ParseBool(option.(string))
		if err != nil {
			return err
		}

		createOpts.Options = &fcu.VpnConnectionOptionsSpecification{
			StaticRoutesOnly: aws.Bool(b),
		}
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

	return resourceOutscaleOAPIVpnConnectionRead(d, meta)
}

func vpnConnectionRefreshFunc(conn *fcu.Client, connectionID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		var resp *fcu.DescribeVpnConnectionsOutput
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeVpnConnections(&fcu.DescribeVpnConnectionsInput{
				VpnConnectionIds: []*string{aws.String(connectionID)},
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

func resourceOutscaleOAPIVpnConnectionRead(d *schema.ResourceData, meta interface{}) error {
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
		}
		log.Printf("[ERROR] Error finding VPN connection: %s", err)
		return err
	}

	if len(resp.VpnConnections) != 1 {
		return fmt.Errorf("[ERROR] Error finding VPN connection: %s", d.Id())
	}

	vpnConnection := resp.VpnConnections[0]
	if vpnConnection == nil || *vpnConnection.State == "deleted" {
		d.SetId("")
		return nil
	}
	vpn := make(map[string]interface{})
	if vpnConnection.Options != nil {
		vpn["static_vpn_static_route_only"] = strconv.FormatBool(aws.BoolValue(vpnConnection.Options.StaticRoutesOnly))
	} else {
		vpn["static_vpn_static_route_only"] = strconv.FormatBool(false)
	}
	if err := d.Set("vpn_connection_option", vpn); err != nil {
		return err
	}
	d.Set("client_endpoint_configuration", vpnConnection.CustomerGatewayConfiguration)

	vpns := make([]map[string]interface{}, len(vpnConnection.Routes))

	for k, v := range vpnConnection.Routes {
		route := make(map[string]interface{})

		route["destination_ip_range"] = *v.DestinationCidrBlock
		route["type"] = *v.Source
		route["state"] = *v.State

		vpns[k] = route
	}

	d.Set("vpn_static_route", vpns)
	d.Set("tag", tagsToMap(vpnConnection.Tags))

	d.Set("state", vpnConnection.State)

	vgws := make([]map[string]interface{}, len(vpnConnection.VgwTelemetry))

	for k, v := range vpnConnection.VgwTelemetry {
		vgw := make(map[string]interface{})

		vgw["accepted_routes_count"] = *v.AcceptedRouteCount
		vgw["outscale_side_ip"] = *v.OutsideIpAddress
		vgw["state"] = *v.Status
		vgw["comment"] = *v.StatusMessage

		vgws[k] = vgw
	}

	d.Set("vpn_tunnel_description", vgws)
	d.Set("vpn_connection_id", vpnConnection.VpnConnectionId)
	d.Set("vpn_gateway_id", vpnConnection.VpnGatewayId)
	d.Set("request_id", resp.RequestId)

	return nil
}

func resourceOutscaleOAPIVpnConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.DeleteVpnConnection(&fcu.DeleteVpnConnectionInput{
			VpnConnectionId: aws.String(d.Id()),
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
		}
		fmt.Printf("[ERROR] Error deleting VPN connection: %s", err)
		return err
	}

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
			"Error waiting for VPN connection (%s) to delete: %s", d.Id(), stateErr)
	}

	return nil
}
