package outscale

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cast"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	oscgo "github.com/outscale/osc-sdk-go/v2"
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

	var resp oscgo.CreateClientGatewayResponse
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		var err error
		rp, httpResp, err := conn.ClientGatewayApi.CreateClientGateway(context.Background()).CreateClientGatewayRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.(*schema.Set), *resp.GetClientGateway().ClientGatewayId, conn)
		if err != nil {
			return err
		}
	}

	d.SetId(*resp.GetClientGateway().ClientGatewayId)

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
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	r, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Outscale Client Gateway (%s) to become ready: %s", clientGatewayID, err)
	}

	resp := r.(oscgo.ReadClientGatewaysResponse)
	if !resp.HasClientGateways() || utils.IsResponseEmpty(len(resp.GetClientGateways()), "ClientGateway", d.Id()) {
		d.SetId("")
		return nil
	}

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

	return nil
}

func resourceOutscaleClientGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}

	return resourceOutscaleClientGatewayRead(d, meta)
}

func resourceOutscaleClientGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	gatewayID := d.Id()
	req := oscgo.DeleteClientGatewayRequest{
		ClientGatewayId: gatewayID,
	}

	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.ClientGatewayApi.DeleteClientGateway(context.Background()).DeleteClientGatewayRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{"deleted", "failed"},
		Refresh:    clientGatewayRefreshFunc(conn, &gatewayID),
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
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
		var resp oscgo.ReadClientGatewaysResponse
		var statusCode int
		err := resource.Retry(120*time.Second, func() *resource.RetryError {
			rp, httpResp, err := conn.ClientGatewayApi.ReadClientGateways(context.Background()).ReadClientGatewaysRequest(filter).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			statusCode = httpResp.StatusCode
			return nil
		})
		if err != nil || len(resp.GetClientGateways()) == 0 {
			switch {
			case statusCode == http.StatusServiceUnavailable || statusCode == http.StatusConflict:
				return nil, "pending", nil
			case statusCode == http.StatusNotFound || len(resp.GetClientGateways()) == 0:
				return nil, "deleted", nil
			default:
				return nil, "failed", fmt.Errorf("Error on clientGatewayRefresh: %s", err)
			}
		}

		gateway := resp.GetClientGateways()[0]

		return resp, gateway.GetState(), nil
	}
}
