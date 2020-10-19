package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	oscgo "github.com/outscale/osc-sdk-go/osc"
)

func resourceOutscaleClientGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleClientGatewayCreate,
		Read:   resourceOutscaleClientGatewayRead,
		Update: resourceOutscaleClientGatewayUpdate,
		Delete: resourceOutscaleClientGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bgp_asn": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"connection_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsListOAPISchema(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleClientGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.CreateClientGatewayRequest{
		BgpAsn:         cast.ToInt32(d.Get("bgp_asn")),
		ConnectionType: d.Get("connection_type").(string),
		PublicIp:       d.Get("public_ip").(string),
	}

	client, _, err := conn.ClientGatewayApi.CreateClientGateway(context.Background()).CreateClientGatewayRequest(req).Execute()
	if err != nil {
		return err
	}

	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.(*schema.Set), *client.GetClientGateway().ClientGatewayId, conn)
		if err != nil {
			return err
		}
	}

	d.SetId(*client.GetClientGateway().ClientGatewayId)

	return resourceOutscaleClientGatewayRead(d, meta)
}

func resourceOutscaleClientGatewayRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	clientGatewayID := d.Id()

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available", "failed"},
		Refresh:    clientGatewayRefreshFunc(conn, &clientGatewayID),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	r, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Outscale Client Gateway (%s) to become ready: %s", clientGatewayID, err)
	}

	resp := r.(oscgo.ReadClientGatewaysResponse)
	clientGateway := resp.GetClientGateways()[0]

	if err := d.Set("bgp_asn", clientGateway.GetBgpAsn()); err != nil {
		return err
	}
	if err := d.Set("connection_type", clientGateway.GetConnectionType()); err != nil {
		return err
	}
	if err := d.Set("public_ip", clientGateway.GetPublicIp()); err != nil {
		return err
	}
	if err := d.Set("client_gateway_id", clientGateway.GetClientGatewayId()); err != nil {
		return err
	}
	if err := d.Set("state", clientGateway.GetState()); err != nil {
		return err
	}
	if err := d.Set("tags", tagsOSCAPIToMap(clientGateway.GetTags())); err != nil {
		return err
	}
	if err := d.Set("request_id", resp.ResponseContext.GetRequestId()); err != nil {
		return err
	}

	return nil
}

func resourceOutscaleClientGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	d.Partial(true)

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

	d.Partial(false)
	return resourceOutscaleClientGatewayRead(d, meta)
}

func resourceOutscaleClientGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	gatewayID := d.Id()
	req := oscgo.DeleteClientGatewayRequest{
		ClientGatewayId: gatewayID,
	}

	_, _, err := conn.ClientGatewayApi.DeleteClientGateway(context.Background()).DeleteClientGatewayRequest(req).Execute()
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{"deleted", "failed"},
		Refresh:    clientGatewayRefreshFunc(conn, &gatewayID),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Outscale Client Gateway (%s) to become deleted: %s", gatewayID, err)
	}

	return nil
}

func clientGatewayRefreshFunc(conn *oscgo.APIClient, gatewayID *string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		filter := oscgo.ReadClientGatewaysRequest{
			Filters: &oscgo.FiltersClientGateway{
				ClientGatewayIds: &[]string{*gatewayID},
			},
		}

		resp, _, err := conn.ClientGatewayApi.ReadClientGateways(context.Background()).ReadClientGatewaysRequest(filter).Execute()
		if err != nil {
			switch {
			case strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:"):
				return nil, "pending", nil
			case strings.Contains(fmt.Sprint(err), "404"):
				return nil, "deleted", nil
			default:
				return nil, "failed", fmt.Errorf("Error on clientGatewayRefresh: %s", err)
			}
		}

		gateway := resp.GetClientGateways()[0]

		return resp, gateway.GetState(), nil
	}
}
