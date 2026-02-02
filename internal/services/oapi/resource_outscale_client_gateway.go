package oapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	oscgo "github.com/outscale/osc-sdk-go/v2"
)

func ResourceOutscaleClientGateway() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleClientGatewayCreate,
		Read:   ResourceOutscaleClientGatewayRead,
		Update: ResourceOutscaleClientGatewayUpdate,
		Delete: ResourceOutscaleClientGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Update: schema.DefaultTimeout(UpdateDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
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
			"tags": TagsSchemaSDK(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceOutscaleClientGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutCreate)

	req := oscgo.CreateClientGatewayRequest{
		BgpAsn:         cast.ToInt32(d.Get("bgp_asn")),
		ConnectionType: d.Get("connection_type").(string),
		PublicIp:       d.Get("public_ip").(string),
	}

	var resp oscgo.CreateClientGatewayResponse
	err := retry.Retry(timeout, func() *retry.RetryError {
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
	d.SetId(*resp.GetClientGateway().ClientGatewayId)

	err = createOAPITagsSDK(conn, d)
	if err != nil {
		return err
	}

	return ResourceOutscaleClientGatewayRead(d, meta)
}

func ResourceOutscaleClientGatewayRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutRead)

	clientGatewayID := d.Id()

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available", "failed", "deleted"},
		Refresh:    clientGatewayRefreshFunc(conn, timeout, &clientGatewayID),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	r, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for outscale client gateway (%s) to become ready: %s", clientGatewayID, err)
	}

	resp := r.(oscgo.ReadClientGatewaysResponse)
	if !resp.HasClientGateways() || utils.IsResponseEmpty(len(resp.GetClientGateways()), "ClientGateway", d.Id()) ||
		resp.GetClientGateways()[0].GetState() == "deleted" {
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
	if err := d.Set("tags", FlattenOAPITagsSDK(clientGateway.GetTags())); err != nil {
		return err
	}

	return nil
}

func ResourceOutscaleClientGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	if err := updateOAPITagsSDK(conn, d); err != nil {
		return err
	}

	return ResourceOutscaleClientGatewayRead(d, meta)
}

func ResourceOutscaleClientGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutDelete)

	gatewayID := d.Id()
	req := oscgo.DeleteClientGatewayRequest{
		ClientGatewayId: gatewayID,
	}

	err := retry.Retry(timeout, func() *retry.RetryError {
		_, httpResp, err := conn.ClientGatewayApi.DeleteClientGateway(context.Background()).DeleteClientGatewayRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{"deleted", "failed"},
		Refresh:    clientGatewayRefreshFunc(conn, timeout, &gatewayID),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for outscale client gateway (%s) to become deleted: %s", gatewayID, err)
	}

	return nil
}

func clientGatewayRefreshFunc(conn *oscgo.APIClient, timeout time.Duration, gatewayID *string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		filter := oscgo.ReadClientGatewaysRequest{
			Filters: &oscgo.FiltersClientGateway{
				ClientGatewayIds: &[]string{*gatewayID},
			},
		}
		var resp oscgo.ReadClientGatewaysResponse
		var statusCode int
		err := retry.Retry(timeout, func() *retry.RetryError {
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
				return nil, "failed", fmt.Errorf("error on clientgatewayrefresh: %s", err)
			}
		}

		gateway := resp.GetClientGateways()[0]

		return resp, gateway.GetState(), nil
	}
}
