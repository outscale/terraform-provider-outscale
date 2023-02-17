package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOutscaleOAPINet() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPINetCreate,
		Read:   resourceOutscaleOAPINetRead,
		Update: resourceOutscaleOAPINetUpdate,
		Delete: resourceOutscaleOAPINetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: getOAPINetSchema(),
	}
}

func resourceOutscaleOAPINetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.CreateNetRequest{
		IpRange: d.Get("ip_range").(string),
	}

	if c, ok := d.GetOk("tenancy"); ok {
		tenancy := c.(string)
		if tenancy == "default" || tenancy == "dedicated" {
			req.SetTenancy(tenancy)
		} else {
			return fmt.Errorf("tenancy option not supported: %s", tenancy)
		}
	}

	var resp oscgo.CreateNetResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.NetApi.CreateNet(context.Background()).CreateNetRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("error creating Outscale Net: %s", utils.GetErrorResponse(err))
	}

	//SetTags
	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.(*schema.Set), resp.Net.GetNetId(), conn)
		if err != nil {
			return err
		}
	}

	d.SetId(resp.Net.GetNetId())

	return resource.Retry(120*time.Second, func() *resource.RetryError {
		err = resourceOutscaleOAPINetRead(d, meta)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if c, ok := d.GetOk("state"); ok {
			state := c.(string)
			if strings.Compare(state, "available") != 0 {
				return resource.RetryableError(fmt.Errorf("Expected Net to be available but was in state %s", state))
			}
		}
		return nil
	})
}

func resourceOutscaleOAPINetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	id := d.Id()

	filters := oscgo.FiltersNet{
		NetIds: &[]string{id},
	}

	req := oscgo.ReadNetsRequest{
		Filters: &filters,
	}

	var resp oscgo.ReadNetsResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.NetApi.ReadNets(context.Background()).ReadNetsRequest(req).Execute()

		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("[DEBUG] Error reading network (%s)", err)
	}
	if utils.IsResponseEmpty(len(resp.GetNets()), "Net", d.Id()) {
		d.SetId("")
		return nil
	}
	if err := d.Set("ip_range", resp.GetNets()[0].GetIpRange()); err != nil {
		return err
	}
	if err := d.Set("tenancy", resp.GetNets()[0].Tenancy); err != nil {
		return err
	}
	if err := d.Set("dhcp_options_set_id", resp.GetNets()[0].GetDhcpOptionsSetId()); err != nil {
		return err
	}
	if err := d.Set("net_id", resp.GetNets()[0].GetNetId()); err != nil {
		return err
	}
	if err := d.Set("state", resp.GetNets()[0].GetState()); err != nil {
		return err
	}

	return d.Set("tags", tagsOSCAPIToMap(resp.GetNets()[0].GetTags()))
}

func resourceOutscaleOAPINetUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	d.Partial(true)

	if err := setOSCAPITags(conn, d, "tags"); err != nil {
		return err
	}

	d.SetPartial("tags")

	d.Partial(false)
	return resourceOutscaleOAPINetRead(d, meta)
}

func resourceOutscaleOAPINetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	id := d.Id()

	req := oscgo.DeleteNetRequest{
		NetId: id,
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"deleted", "failed"},
		Refresh: func() (interface{}, string, error) {
			err := resource.Retry(120*time.Second, func() *resource.RetryError {
				_, httpResp, err := conn.NetApi.DeleteNet(context.Background()).DeleteNetRequest(req).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})
			if err != nil {
				return nil, "failed", err
			}
			return "", "deleted", nil
		},
		Timeout:    10 * time.Minute,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error deleting Net Service(%s): %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}

func getOAPINetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ip_range": {
			Type:     schema.TypeString,
			ForceNew: true,
			Required: true,
		},
		"tenancy": {
			Type:     schema.TypeString,
			ForceNew: true,
			Computed: true,
			Optional: true,
		},

		// Attributes
		"dhcp_options_set_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tags": tagsListOAPISchema(),
		"net_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
