package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
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
	conn := meta.(*OutscaleClient).OAPI

	req := &oapi.UpdateNetRequest{}

	req.NetId = d.Get("net_id").(string)

	if c, ok := d.GetOk("dhcp_options_set_id"); ok {
		req.DhcpOptionsSetId = c.(string)
	}

	var err error
	var resp *oapi.POST_UpdateNetResponses
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.POST_UpdateNet(*req)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var errString string

	if err != nil || resp.OK == nil {
		if err != nil {
			errString = err.Error()
		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("Status Code: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("Status Code: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("Status: 500, %s", utils.ToJSONString(resp.Code500))
		}
		return fmt.Errorf("[DEBUG] Error creating net attribute. Details: %s", errString)
	}

	d.Set("request_id", resp.OK.ResponseContext.RequestId)

	d.SetId(resource.UniqueId())

	return resourceOutscaleOAPILinAttrRead(d, meta)
}

func resourceOutscaleOAPILinAttrUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	req := &oapi.UpdateNetRequest{}

	if d.HasChange("net_id") && !d.IsNewResource() {
		req.NetId = d.Get("net_id").(string)
	}
	if d.HasChange("dhcp_options_set_id") && !d.IsNewResource() {
		req.DhcpOptionsSetId = d.Get("dhcp_options_set_id").(string)
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.POST_UpdateNet(*req)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[DEBUG] Error creating lin (%s)", err)
		return err
	}

	return resourceOutscaleOAPILinAttrRead(d, meta)
}

func resourceOutscaleOAPILinAttrRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	filters := oapi.FiltersNet{
		NetIds: []string{d.Get("net_id").(string)},
	}

	req := oapi.ReadNetsRequest{
		Filters: filters,
	}

	var rs *oapi.POST_ReadNetsResponses
	var resp *oapi.ReadNetsResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rs, err = conn.POST_ReadNets(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		log.Printf("[DEBUG] Error reading lin (%s)", err)
	}

	resp = rs.OK

	if resp == nil {
		d.SetId("")
		return fmt.Errorf("oAPI Lin not found")
	}

	if len(resp.Nets) == 0 {
		d.SetId("")
		return fmt.Errorf("oAPI Lin not found")
	}

	d.Set("net_id", resp.Nets[0].NetId)
	d.Set("dhcp_options_set_id", resp.Nets[0].DhcpOptionsSetId)

	return nil
}

func resourceOutscaleOAPILinAttrDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	return nil
}
