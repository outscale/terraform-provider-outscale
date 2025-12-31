package oapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func ResourceOutscaleVirtualGateway() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleVirtualGatewayCreate,
		Read:   ResourceOutscaleVirtualGatewayRead,
		Update: ResourceOutscaleVirtualGatewayUpdate,
		Delete: ResourceOutscaleVirtualGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"connection_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"net_to_virtual_gateway_links": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"virtual_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": TagsSchemaSDK(),
		},
	}
}

func ResourceOutscaleVirtualGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	connectType, connecTypeOk := d.GetOk("connection_type")
	createOpts := oscgo.CreateVirtualGatewayRequest{}
	if connecTypeOk {
		createOpts.SetConnectionType(connectType.(string))
	}

	var resp oscgo.CreateVirtualGatewayResponse
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		var err error
		rp, httpResp, err := conn.VirtualGatewayApi.CreateVirtualGateway(context.Background()).CreateVirtualGatewayRequest(createOpts).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error creating VPN gateway: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"pending", "ending/wait"},
		Target:     []string{"available"},
		Refresh:    virtualGatewayStateRefreshFunc(conn, resp.VirtualGateway.GetVirtualGatewayId(), "terminated"),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to become created: %s", d.Id(), err)
	}

	virtualGateway := resp.GetVirtualGateway()
	d.SetId(virtualGateway.GetVirtualGatewayId())

	if d.IsNewResource() {
		if err := updateOAPITagsSDK(conn, d); err != nil {
			return err
		}
	}
	return ResourceOutscaleVirtualGatewayRead(d, meta)
}

func ResourceOutscaleVirtualGatewayRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	var resp oscgo.ReadVirtualGatewaysResponse
	var err error
	var statusCode int

	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.VirtualGatewayApi.ReadVirtualGateways(context.Background()).ReadVirtualGatewaysRequest(oscgo.ReadVirtualGatewaysRequest{
			Filters: &oscgo.FiltersVirtualGateway{VirtualGatewayIds: &[]string{d.Id()}},
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		statusCode = httpResp.StatusCode
		return nil
	})
	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		fmt.Printf("\n\n[ERROR] Error finding VpnGateway: %s", err)
		return err
	}

	if utils.IsResponseEmpty(len(resp.GetVirtualGateways()), "VirtualGateway", d.Id()) {
		d.SetId("")
		return nil
	}
	virtualGateway := resp.GetVirtualGateways()[0]
	if virtualGateway.GetState() == "deleted" {
		d.SetId("")
		return nil
	}
	if virtualGateway.HasNetToVirtualGatewayLinks() {
		vs := make([]map[string]interface{}, len(virtualGateway.GetNetToVirtualGatewayLinks()))
		for k, v := range virtualGateway.GetNetToVirtualGatewayLinks() {
			vp := make(map[string]interface{})
			vp["state"] = v.GetState()
			vp["net_id"] = v.GetNetId()
			vs[k] = vp
		}
		d.Set("net_to_virtual_gateway_links", vs)
	}

	d.Set("connection_type", virtualGateway.GetConnectionType())
	d.Set("virtual_gateway_id", virtualGateway.GetVirtualGatewayId())

	d.Set("state", virtualGateway.State)
	d.Set("tags", FlattenOAPITagsSDK(virtualGateway.GetTags()))

	return nil
}

func ResourceOutscaleVirtualGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	if err := updateOAPITagsSDK(conn, d); err != nil {
		return err
	}
	return nil
}

func ResourceOutscaleVirtualGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	return retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, httpResp, err := conn.VirtualGatewayApi.DeleteVirtualGateway(context.Background()).DeleteVirtualGatewayRequest(
			oscgo.DeleteVirtualGatewayRequest{VirtualGatewayId: d.Id()}).Execute()
		if err != nil {
			if httpResp.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
			return utils.CheckThrottling(httpResp, err)
		}
		d.SetId("")
		return nil
	})
}

// vpnGatewayAttachStateRefreshFunc returns a retry.StateRefreshFunc that is used to watch
// the state of a VPN gateway's attachment
func vpnGatewayAttachStateRefreshFunc(conn *oscgo.APIClient, id string, expected string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp oscgo.ReadVirtualGatewaysResponse
		var err error
		var statusCode int

		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			rp, httpResp, err := conn.VirtualGatewayApi.ReadVirtualGateways(context.Background()).ReadVirtualGatewaysRequest(oscgo.ReadVirtualGatewaysRequest{
				Filters: &oscgo.FiltersVirtualGateway{VirtualGatewayIds: &[]string{id}},
			}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			statusCode = httpResp.StatusCode
			return nil
		})
		if err != nil {
			if statusCode == http.StatusNotFound {
				resp.SetVirtualGateways(nil)
			} else {
				fmt.Printf("[ERROR] Error on VpnGatewayStateRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp.GetVirtualGateways() == nil {
			return nil, "", nil
		}

		virtualGateway := resp.GetVirtualGateways()[0]
		if len(virtualGateway.GetNetToVirtualGatewayLinks()) == 0 {
			return virtualGateway, "detached", nil
		}

		vpnAttachment := oapiVpnGatewayGetLink(virtualGateway)
		return virtualGateway, vpnAttachment.GetState(), nil
	}
}

func oapiVpnGatewayGetLink(vgw oscgo.VirtualGateway) *oscgo.NetToVirtualGatewayLink {
	for _, v := range vgw.GetNetToVirtualGatewayLinks() {
		if v.GetState() == "attached" {
			return &v
		}
	}
	return &oscgo.NetToVirtualGatewayLink{State: ptr.To("detached")}
}

func virtualGatewayStateRefreshFunc(conn *oscgo.APIClient, instanceID, failState string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp oscgo.ReadVirtualGatewaysResponse
		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			var err error
			rp, httpResp, err := conn.VirtualGatewayApi.ReadVirtualGateways(context.Background()).ReadVirtualGatewaysRequest(oscgo.ReadVirtualGatewaysRequest{
				Filters: &oscgo.FiltersVirtualGateway{
					VirtualGatewayIds: &[]string{instanceID},
				},
			}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] error on InstanceStateRefresh: %s", err)
			return nil, "", err
		}

		if !resp.HasVirtualGateways() {
			return nil, "", nil
		}

		virtualGateway := resp.GetVirtualGateways()[0]
		state := virtualGateway.GetState()

		if state == failState {
			return virtualGateway, state, fmt.Errorf("Failed to reach target state. Reason: %v", *virtualGateway.State)
		}

		return virtualGateway, state, nil
	}
}
