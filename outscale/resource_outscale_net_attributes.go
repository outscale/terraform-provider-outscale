package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		},
	}
}

func resourceOutscaleOAPILinAttrCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.UpdateNetRequest{}

	req.SetNetId(d.Get("net_id").(string))

	if c, ok := d.GetOk("dhcp_options_set_id"); ok {
		req.SetDhcpOptionsSetId(c.(string))
	}

	var err error
	var resp oscgo.UpdateNetResponse
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, _, err = conn.NetApi.UpdateNet(context.Background(), &oscgo.UpdateNetOpts{UpdateNetRequest: optional.NewInterface(req)})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error creating net attribute. Details: %s", utils.GetErrorResponse(err))
	}

	d.Set("request_id", resp.ResponseContext.GetRequestId())

	d.SetId(resource.UniqueId())

	return resourceOutscaleOAPILinAttrRead(d, meta)
}

func resourceOutscaleOAPILinAttrUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.UpdateNetRequest{}

	if d.HasChange("net_id") && !d.IsNewResource() {
		req.SetNetId(d.Get("net_id").(string))
	}
	if d.HasChange("dhcp_options_set_id") && !d.IsNewResource() {
		req.SetDhcpOptionsSetId(d.Get("dhcp_options_set_id").(string))
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, _, err = conn.NetApi.UpdateNet(context.Background(), &oscgo.UpdateNetOpts{UpdateNetRequest: optional.NewInterface(req)})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("[DEBUG] Error creating lin (%s)", utils.GetErrorResponse(err))

	}

	return resourceOutscaleOAPILinAttrRead(d, meta)
}

func resourceOutscaleOAPILinAttrRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters := oscgo.FiltersNet{
		NetIds: &[]string{d.Get("net_id").(string)},
	}

	req := oscgo.ReadNetsRequest{
		Filters: &filters,
	}

	var resp oscgo.ReadNetsResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, _, err = conn.NetApi.ReadNets(context.Background(), &oscgo.ReadNetsOpts{ReadNetsRequest: optional.NewInterface(req)})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		log.Printf("[DEBUG] Error reading lin (%s)", utils.GetErrorResponse(err))
	}

	if len(resp.GetNets()) == 0 {
		d.SetId("")
		return fmt.Errorf("oAPI Lin not found")
	}

	d.Set("net_id", resp.GetNets()[0].GetNetId())
	d.Set("dhcp_options_set_id", resp.GetNets()[0].GetDhcpOptionsSetId())

	return nil
}

func resourceOutscaleOAPILinAttrDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}
