package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOutscaleOAPILinAttributes() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILinAttrCreate,
		Read:   resourceOutscaleOAPILinAttrRead,
		Update: resourceOutscaleOAPILinAttrUpdate,
		Delete: resourceOutscaleOAPILinAttrDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"dhcp_options_set_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"net_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_range": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenancy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsOAPIListSchemaComputed(),
		},
	}
}

func resourceOutscaleOAPILinAttrCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.UpdateNetRequest{
		NetId:            d.Get("net_id").(string),
		DhcpOptionsSetId: "default",
	}

	if v, ok := d.GetOk("dhcp_options_set_id"); ok {
		req.SetDhcpOptionsSetId(v.(string))
	}

	var err error
	var resp oscgo.UpdateNetResponse
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		var err error
		rp, httpResp, err := conn.NetApi.UpdateNet(context.Background()).UpdateNetRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error creating net attribute. Details: %s", utils.GetErrorResponse(err))
	}

	d.SetId(resp.Net.GetNetId())

	return resourceOutscaleOAPILinAttrRead(d, meta)
}

func resourceOutscaleOAPILinAttrUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.UpdateNetRequest{
		NetId:            d.Get("net_id").(string),
		DhcpOptionsSetId: d.Get("dhcp_options_set_id").(string),
	}

	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.NetApi.UpdateNet(context.Background()).UpdateNetRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("[DEBUG] Error creating lin (%s)", utils.GetErrorResponse(err))
	}

	d.SetId(d.Get("net_id").(string))

	return resourceOutscaleOAPILinAttrRead(d, meta)
}

func resourceOutscaleOAPILinAttrRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadNetsRequest{
		Filters: &oscgo.FiltersNet{
			NetIds: &[]string{d.Id()},
		},
	}

	var resp oscgo.ReadNetsResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		var err error
		rp, httpResp, err := conn.NetApi.ReadNets(context.Background()).ReadNetsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		log.Printf("[DEBUG] Error reading lin (%s)", utils.GetErrorResponse(err))
	}
	if utils.IsResponseEmpty(len(resp.GetNets()), "NetAttributes", d.Id()) {
		d.SetId("")
		return nil
	}
	if err := d.Set("net_id", resp.GetNets()[0].GetNetId()); err != nil {
		return err
	}
	if err := d.Set("dhcp_options_set_id", resp.GetNets()[0].GetDhcpOptionsSetId()); err != nil {
		return err
	}

	d.Set("ip_range", resp.GetNets()[0].GetIpRange())
	d.Set("tenancy", resp.GetNets()[0].Tenancy)
	d.Set("dhcp_options_set_id", resp.GetNets()[0].GetDhcpOptionsSetId())
	d.Set("net_id", resp.GetNets()[0].GetNetId())
	d.Set("state", resp.GetNets()[0].GetState())

	return d.Set("tags", tagsOSCAPIToMap(resp.GetNets()[0].GetTags()))
}

func resourceOutscaleOAPILinAttrDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}
