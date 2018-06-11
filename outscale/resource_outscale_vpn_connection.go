package outscale

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleVpnConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleVpnConnectionCreate,
		Read:   resourceOutscaleVpnConnectionRead,
		Delete: resourceOutscaleVpnConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
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
	typev, ok3 := d.GetOk("type")
	options, ok4 := d.GetOk("options")

	if !ok && !ok2 && ok3 {
		return fmt.Errorf("please provide the required attributes customer_gateway_id, vpn_gateway_id and type")
	}

	createOpts := &fcu.CreateVpnConnectionInput{
		CustomerGatewayId: aws.String(cgid.(string)),
		Type:              aws.String(typev.(string)),
		VpnGatewayId:      aws.String(vpngid.(string)),
	}

	if ok4 {
		opt := options.(map[string]interface{})
		option := opt["static_routes_only"]

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
		return nil
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
		return nil
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
	options := make(map[string]interface{})
	opt := "false"
	if vpnConnection.Options != nil && *vpnConnection.Options.StaticRoutesOnly {
		opt = "true"
	}
	options["static_routes_only"] = opt
	if err := d.Set("options", options); err != nil {
		return err
	}

	d.Set("customer_gateway_configuration", aws.StringValue(vpnConnection.CustomerGatewayConfiguration))

	routes := make([]map[string]interface{}, len(vpnConnection.Routes))

	for k, v := range vpnConnection.Routes {
		route := make(map[string]interface{})

		route["destination_cidr_block"] = aws.StringValue(v.DestinationCidrBlock)
		route["source"] = aws.StringValue(v.Source)
		route["state"] = aws.StringValue(v.State)

		routes[k] = route
	}

	d.Set("routes", routes)
	d.Set("tag_set", tagsToMap(vpnConnection.Tags))
	d.Set("state", aws.StringValue(vpnConnection.State))

	vgws := make([]map[string]interface{}, len(vpnConnection.VgwTelemetry))

	for k, v := range vpnConnection.VgwTelemetry {
		vgw := make(map[string]interface{})
		vgw["accepted_route_count"] = aws.Int64Value(v.AcceptedRouteCount)
		vgw["outside_ip_address"] = aws.StringValue(v.OutsideIpAddress)
		vgw["status"] = aws.StringValue(v.Status)
		vgw["status_message"] = aws.StringValue(v.StatusMessage)

		vgws[k] = vgw
	}

	d.Set("vpn_connection_id", vpnConnection.VpnConnectionId)
	d.Set("request_id", resp.RequestId)

	return d.Set("vgw_telemetry", vgws)
}

func resourceOutscaleVpnConnectionDelete(d *schema.ResourceData, meta interface{}) error {
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
		Pending:    []string{"pending", "deleting"},
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
