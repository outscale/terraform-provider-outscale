package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

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
			Optional: true,
			Computed: true,
		},
		"dhcp_options_set_id": {
			Type:     schema.TypeString,
			Optional: true,
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
	d.SetId(resp.Net.GetNetId())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available"},
		Refresh:    netStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for net (%s) to create: %s", d.Id(), err)
	}

	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.(*schema.Set), d.Id(), conn)
		if err != nil {
			return err
		}
	}

	if _, ok := d.GetOk("dhcp_options_set_id"); ok {
		return updateNet(d, conn)
	}

	return resourceOutscaleOAPINetRead(d, meta)
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
	return netSetter(d, &resp.GetNets()[0])
}

func updateNet(d *schema.ResourceData, conn *oscgo.APIClient) error {
	req := oscgo.UpdateNetRequest{
		NetId:            d.Id(),
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
		return fmt.Errorf("[DEBUG] Error updating net (%s)", utils.GetErrorResponse(err))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available"},
		Refresh:    netStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	resp, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for net (%s) to create: %s", d.Id(), err)
	}
	rp := resp.(oscgo.ReadNetsResponse)
	return netSetter(d, &rp.GetNets()[0])
}

func resourceOutscaleOAPINetUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}

	if d.HasChange("dhcp_options_set_id") {
		return updateNet(d, conn)
	}
	return resourceOutscaleOAPINetRead(d, meta)
}

func resourceOutscaleOAPINetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	id := d.Id()

	req := oscgo.DeleteNetRequest{
		NetId: id,
	}
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.NetApi.DeleteNet(context.Background()).DeleteNetRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting Net Service(%s): %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"deleted", "failed"},
		Refresh:    netStateRefreshFunc(conn, id),
		Timeout:    10 * time.Minute,
		MinTimeout: 30 * time.Second,
		Delay:      20 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error deleting Net Service(%s): %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}

func netStateRefreshFunc(conn *oscgo.APIClient, netID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp oscgo.ReadNetsResponse
		filters := oscgo.FiltersNet{
			NetIds: &[]string{netID},
		}

		req := oscgo.ReadNetsRequest{
			Filters: &filters,
		}

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
			return nil, "failed", fmt.Errorf("[DEBUG] Error reading network (%s)", err)
		}
		if !resp.HasNets() || len(resp.GetNets()) < 1 {
			return resp, "deleted", nil
		}
		net := &resp.GetNets()[0]
		state := net.GetState()
		return resp, state, nil
	}
}

func netSetter(d *schema.ResourceData, net *oscgo.Net) error {
	if err := d.Set("ip_range", net.GetIpRange()); err != nil {
		return err
	}
	if err := d.Set("tenancy", net.GetTenancy()); err != nil {
		return err
	}
	if err := d.Set("dhcp_options_set_id", net.GetDhcpOptionsSetId()); err != nil {
		return err
	}
	if err := d.Set("net_id", net.GetNetId()); err != nil {
		return err
	}
	if err := d.Set("state", net.GetState()); err != nil {
		return err
	}
	return d.Set("tags", tagsOSCAPIToMap(net.GetTags()))
}
