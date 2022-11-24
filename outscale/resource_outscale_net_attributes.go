package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLinAttributes() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinAttrCreate,
		Read:   resourceLinAttrRead,
		Update: resourceLinAttrUpdate,
		Delete: resourceLinAttrDelete,
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
			"tags": tagsListSchemaComputed(),
		},
	}
}

func resourceLinAttrCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).OSCAPI

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
			return utils.CheckThrottling(httpResp.StatusCode, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error creating net attribute. Details: %s", utils.GetErrorResponse(err))
	}

	d.SetId(resp.Net.GetNetId())

	return resourceLinAttrRead(d, meta)
}

func resourceLinAttrUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).OSCAPI

	req := oscgo.UpdateNetRequest{
		NetId:            d.Get("net_id").(string),
		DhcpOptionsSetId: d.Get("dhcp_options_set_id").(string),
	}

	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.NetApi.UpdateNet(context.Background()).UpdateNetRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp.StatusCode, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("[DEBUG] Error creating lin (%s)", utils.GetErrorResponse(err))
	}

	d.SetId(d.Get("net_id").(string))

	return resourceLinAttrRead(d, meta)
}

func resourceLinAttrRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).OSCAPI

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
			return utils.CheckThrottling(httpResp.StatusCode, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		log.Printf("[DEBUG] Error reading lin (%s)", utils.GetErrorResponse(err))
	}

	if len(resp.GetNets()) == 0 {
		d.SetId("")
		return fmt.Errorf("network is not found")
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

	return d.Set("tags", tagsToMap(resp.GetNets()[0].GetTags()))
}

func resourceLinAttrDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}
