package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOutscaleOAPIVirtualGatewayLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIVirtualGatewayLinkCreate,
		Read:   resourceOutscaleOAPIVirtualGatewayLinkRead,
		Delete: resourceOutscaleOAPIVirtualGatewayLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"net_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"virtual_gateway_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dry_run": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"net_to_virtual_gateway_links": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPIVirtualGatewayLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	vgwID := d.Get("virtual_gateway_id").(string)

	var resp oscgo.ReadVirtualGatewaysResponse
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.VirtualGatewayApi.ReadVirtualGateways(context.Background(), &oscgo.ReadVirtualGatewaysOpts{ReadVirtualGatewaysRequest: optional.NewInterface(oscgo.ReadVirtualGatewaysRequest{
			Filters: &oscgo.FiltersVirtualGateway{VirtualGatewayIds: &[]string{vgwID}},
		})})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "InvalidVPNGatewayID.NotFound" {
			log.Printf("[WARN] VPN Gateway %q not found.", vgwID)
			d.SetId("")
			return nil
		}
		return err
	}

	vgw := resp.GetVirtualGateways()[0]
	if vgw.GetState() == "deleted" {
		log.Printf("[INFO] VPN Gateway %q appears to have been deleted.", vgwID)
		d.SetId("")
		return nil
	}

	vga := oapiVpnGatewayGetLink(vgw)
	if len(vgw.GetNetToVirtualGatewayLinks()) == 0 || vga.GetState() == "detached" {
		//d.Set("net_id", "")
		return nil
	}
	vs := make([]map[string]interface{}, len(vgw.GetNetToVirtualGatewayLinks()))

	for k, v := range vgw.GetNetToVirtualGatewayLinks() {
		vp := make(map[string]interface{})
		vp["state"] = v.GetState()
		vp["net_id"] = v.GetNetId()

		vs[k] = vp
	}
	d.Set("net_to_virtual_gateway_links", vs)
	d.Set("request_id", resp.ResponseContext.RequestId)

	return nil
}

func resourceOutscaleOAPIVirtualGatewayLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	netID := d.Get("net_id").(string)
	vgwID := d.Get("virtual_gateway_id").(string)

	createOpts := oscgo.LinkVirtualGatewayRequest{
		NetId:            netID,
		VirtualGatewayId: vgwID,
	}
	log.Printf("[DEBUG] VPN Gateway attachment options: %#v", createOpts)

	var err error

	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		_, _, err = conn.VirtualGatewayApi.LinkVirtualGateway(context.Background(), &oscgo.LinkVirtualGatewayOpts{LinkVirtualGatewayRequest: optional.NewInterface(createOpts)})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidVirtualGatewayID.NotFound") {
				return resource.RetryableError(
					fmt.Errorf("Gateway not found, retry for eventual consistancy"))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error attaching Virtual Gateway %q to VPC %q: %s",
			vgwID, netID, err)
	}

	d.SetId(vpnGatewayLinkID(netID, vgwID))

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"detached", "attaching"},
		Target:     []string{"attached"},
		Refresh:    vpnGatewayLinkStateRefresh(conn, netID, vgwID),
		Timeout:    15 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Virtual Gateway %q to attach to VPC %q: %s",
			vgwID, netID, err)
	}
	log.Printf("[DEBUG] VPN Gateway %q attached to VPC %q.", vgwID, netID)

	return resourceOutscaleOAPIVirtualGatewayLinkRead(d, meta)
}

func resourceOutscaleOAPIVirtualGatewayLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	// Get the old VPC ID to detach from
	netID, _ := d.GetChange("net_id")

	if netID.(string) == "" {
		fmt.Printf(
			"[DEBUG] Not detaching Virtual Gateway '%s' as no VPC ID is set",
			d.Get("virtual_gateway_id").(string))
		return nil
	}

	fmt.Printf(
		"[INFO] Detaching Virtual Gateway '%s' from VPC '%s'",
		d.Get("virtual_gateway_id").(string),
		netID.(string))

	wait := true

	var err error
	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		_, _, err = conn.VirtualGatewayApi.UnlinkVirtualGateway(context.Background(), &oscgo.UnlinkVirtualGatewayOpts{UnlinkVirtualGatewayRequest: optional.NewInterface(oscgo.UnlinkVirtualGatewayRequest{
			VirtualGatewayId: d.Get("virtual_gateway_id").(string),
			NetId:            netID.(string),
		})})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidVpnGatewayID.NotFound") {
				return resource.RetryableError(
					fmt.Errorf("Gateway not found, retry for eventual consistancy"))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidVpnGatewayID.NotFound") {
			err = nil
			wait = false
		} else if strings.Contains(fmt.Sprint(err), "InvalidVpnGatewayAttachment.NotFound") {
			err = nil
			wait = false
		}

		if err != nil {
			return err
		}
	}

	if !wait {
		return nil
	}

	// Wait for it to be fully detached before continuing
	log.Printf("[DEBUG] Waiting for VPN gateway (%s) to detach", d.Get("virtual_gateway_id").(string))
	stateConf := &resource.StateChangeConf{
		Pending: []string{"attached", "detaching", "available"},
		Target:  []string{"detached"},
		Refresh: vpnGatewayAttachStateRefreshFunc(conn, d.Get("virtual_gateway_id").(string), "detached"),
		Timeout: 5 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for vpn gateway (%s) to detach: %s",
			d.Get("virtual_gateway_id").(string), err)
	}

	return nil
}

func vpnGatewayLinkStateRefresh(conn *oscgo.APIClient, vpcID, vgwID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var err error
		var resp oscgo.ReadVirtualGatewaysResponse
		err = resource.Retry(30*time.Second, func() *resource.RetryError {
			resp, _, err = conn.VirtualGatewayApi.ReadVirtualGateways(context.Background(), &oscgo.ReadVirtualGatewaysOpts{ReadVirtualGatewaysRequest: optional.NewInterface(
				oscgo.ReadVirtualGatewaysRequest{Filters: &oscgo.FiltersVirtualGateway{
					VirtualGatewayIds: &[]string{vgwID},
				}})})
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "InvalidVpnGatewayID.NotFound") {
					return resource.RetryableError(
						fmt.Errorf("Gateway not found, retry for eventual consistancy"))
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			awsErr, ok := err.(awserr.Error)
			if ok {
				switch awsErr.Code() {
				case "InvalidVPNGatewayID.NotFound":
					fallthrough
				case "InvalidVpnGatewayAttachment.NotFound":
					return nil, "", nil
				}
			}

			return nil, "", err
		}

		vgw := resp.GetVirtualGateways()[0]
		if len(vgw.GetNetToVirtualGatewayLinks()) == 0 {
			return vgw, "detached", nil
		}

		vga := oapiVpnGatewayGetLink(vgw)

		log.Printf("[DEBUG] VPN Gateway %q attachment status: %s", vgwID, *vga.State)
		return vgw, *vga.State, nil
	}
}

func vpnGatewayLinkID(vpcID, vgwID string) string {
	return fmt.Sprintf("vpn-attachment-%x", hashcode.String(fmt.Sprintf("%s-%s", vpcID, vgwID)))
}
