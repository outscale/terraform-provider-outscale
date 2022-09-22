package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleOAPISubNet() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPISubNetCreate,
		Read:   resourceOutscaleOAPISubNetRead,
		Update: resourceOutscaleOAPISubNetUpdate,
		Delete: resourceOutscaleOAPISubNetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: getOAPISubNetSchema(),
	}
}

// Create SubNet
func resourceOutscaleOAPISubNetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.CreateSubnetRequest{
		IpRange: d.Get("ip_range").(string),
		NetId:   d.Get("net_id").(string),
	}
	if a, aok := d.GetOk("subregion_name"); aok {
		req.SetSubregionName(a.(string))
	}
	var resp oscgo.CreateSubnetResponse
	var err error
	err = resource.Retry(40*time.Second, func() *resource.RetryError {
		r, _, err := conn.SubnetApi.CreateSubnet(context.Background()).CreateSubnetRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(err)
		}
		resp = r
		return nil
	})
	if err != nil {
		errString := err.Error()
		return fmt.Errorf("[DEBUG] Error creating Subnet (%s)", errString)
	}
	result := resp.GetSubnet()
	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.(*schema.Set), result.GetSubnetId(), conn)
		if err != nil {
			return err
		}
	}
	if result.GetState() != "available" {
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending"},
			Target:     []string{"available"},
			Refresh:    SubnetStateOApiRefreshFunc(conn, result.GetSubnetId()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      6 * time.Second,
			MinTimeout: 1 * time.Second,
		}
		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for subnet (%s) to become created: %s", d.Id(), err)
		}
	}
	d.SetId(result.GetSubnetId())
	return resourceOutscaleOAPISubNetRead(d, meta)
}

// Read SubNet
func resourceOutscaleOAPISubNetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	id := d.Id()
	log.Printf("[DEBUG] Reading Subnet(%s)", id)
	req := oscgo.ReadSubnetsRequest{
		Filters: &oscgo.FiltersSubnet{
			SubnetIds: &[]string{id},
		},
	}
	var resp oscgo.ReadSubnetsResponse
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		r, _, err := conn.SubnetApi.ReadSubnets(context.Background()).ReadSubnetsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(err)
		}
		resp = r
		return nil
	})
	if err != nil {
		errString := err.Error()
		return fmt.Errorf("[DEBUG] Error reading Subnet (%s)", errString)
	}
	if len(resp.GetSubnets()) > 0 {
		return readOutscaleOAPISubNet(d, &resp.GetSubnets()[0])
	}
	return fmt.Errorf("No subnet (%s) found", d.Id())
}
func resourceOutscaleOAPISubNetUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	d.Partial(true)
	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}
	d.SetPartial("tags")
	d.Partial(false)
	return resourceOutscaleOAPISubNetRead(d, meta)
}
func resourceOutscaleOAPISubNetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	id := d.Id()
	log.Printf("[DEBUG] Deleting Subnet (%s)", id)
	req := oscgo.DeleteSubnetRequest{
		SubnetId: id,
	}
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.SubnetApi.DeleteSubnet(context.Background()).DeleteSubnetRequest(req).Execute()
		if err != nil {
			if strings.Contains(err.Error(), utils.ResourceConflict) {
				log.Printf("[DEBUG] Subnet waiting delete: (%s)", err)
				return resource.RetryableError(err)
			}
			return utils.CheckThrottling(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[DEBUG] Error deleting Subnet(%s)", err)
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "available"},
		Target:     []string{"deleted"},
		Refresh:    SubnetStateOApiRefreshFunc(conn, id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      2 * time.Second,
		MinTimeout: 1 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
func readOutscaleOAPISubNet(d *schema.ResourceData, subnet *oscgo.Subnet) error {
	if err := d.Set("subregion_name", subnet.GetSubregionName()); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet1 (%s)", err)
		return err
	}
	if err := d.Set("available_ips_count", subnet.GetAvailableIpsCount()); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet2 (%s)", err)
		return err
	}
	if err := d.Set("ip_range", subnet.GetIpRange()); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet (%s)", err)
		return err
	}
	if err := d.Set("state", subnet.GetState()); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet4 (%s)", err)
		return err
	}
	if err := d.Set("subnet_id", subnet.GetSubnetId()); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet5 (%s)", err)
		return err
	}
	if err := d.Set("net_id", subnet.GetNetId()); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet6 (%s)", err)
		return err
	}
	if err := d.Set("map_public_ip_on_launch", subnet.GetMapPublicIpOnLaunch()); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet6 (%s)", err)
		return err
	}
	return d.Set("tags", tagsOSCAPIToMap(subnet.GetTags()))
}

func SubnetStateOApiRefreshFunc(conn *oscgo.APIClient, subnetID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, _, err := conn.SubnetApi.ReadSubnets(context.Background()).ReadSubnetsRequest(oscgo.ReadSubnetsRequest{
			Filters: &oscgo.FiltersSubnet{
				SubnetIds: &[]string{subnetID},
			},
		}).Execute()
		if err != nil {
			log.Printf("[ERROR] error on SubnetStateRefresh: %s", err)
			return nil, "error", err
		}
		if len(resp.GetSubnets()) == 0 {
			return oscgo.Subnet{}, "deleted", nil
		}
		return resp.GetSubnets()[0], resp.GetSubnets()[0].GetState(), nil
	}
}
func getOAPISubNetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		//This is attribute part for schema SubNet
		"net_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"ip_range": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"subregion_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		//This is arguments part for schema SubNet
		"available_ips_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"subnet_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"map_public_ip_on_launch": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"tags": tagsListOAPISchema(),
	}
}
