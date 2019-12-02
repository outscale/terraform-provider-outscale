package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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

//Create SubNet
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
	err = resource.Retry(40*time.Minute, func() *resource.RetryError {
		r, _, err := conn.SubnetApi.CreateSubnet(context.Background(), &oscgo.CreateSubnetOpts{CreateSubnetRequest: optional.NewInterface(req)})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
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
		err := assignTags(tags.([]interface{}), result.GetSubnetId(), conn)
		if err != nil {
			return err
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "ending/wait"},
		Target:     []string{"available"},
		Refresh:    SubnetStateOApiRefreshFunc(conn, result.GetSubnetId()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for subnet (%s) to become created: %s", d.Id(), err)
	}

	d.SetId(result.GetSubnetId())

	return resourceOutscaleOAPISubNetRead(d, meta)
}

//Read SubNet

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
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		r, _, err := conn.SubnetApi.ReadSubnets(context.Background(), &oscgo.ReadSubnetsOpts{ReadSubnetsRequest: optional.NewInterface(req)})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		resp = r
		return nil
	})

	if err != nil {
		errString := err.Error()

		return fmt.Errorf("[DEBUG] Error reading Subnet (%s)", errString)
	}

	d.Set("request_id", resp.ResponseContext.GetRequestId())

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
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, _, err = conn.SubnetApi.DeleteSubnet(context.Background(), &oscgo.DeleteSubnetOpts{DeleteSubnetRequest: optional.NewInterface(req)})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		log.Printf("[DEBUG] Error deleting Subnet(%s)", err)
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "ending/wait"},
		Target:     []string{"deleted"},
		Refresh:    SubnetStateOApiRefreshFunc(conn, id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

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

	return d.Set("tags", tagsOSCAPIToMap(subnet.GetTags()))
}

func SubnetStateOApiRefreshFunc(conn *oscgo.APIClient, subnetID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, _, err := conn.SubnetApi.ReadSubnets(context.Background(), &oscgo.ReadSubnetsOpts{
			ReadSubnetsRequest: optional.NewInterface(oscgo.ReadSubnetsRequest{
				Filters: &oscgo.FiltersSubnet{
					SubnetIds: &[]string{subnetID},
				},
			}),
		})

		if err != nil {
			log.Printf("[ERROR] error on SubnetStateRefresh: %s", err)
			return nil, "", err
		}

		if !resp.HasSubnets() || len(resp.GetSubnets()) == 0 {
			return nil, "deleted", nil
		}

		subnet := resp.GetSubnets()[0]

		return subnet, subnet.GetState(), nil
	}
}

func getOAPISubNetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		//This is attribute part for schema SubNet
		"net_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"ip_range": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"subregion_name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		//This is arguments part for schema SubNet
		"available_ips_count": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},

		"state": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"subnet_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"request_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"tags": tagsListOAPISchema(),
	}
}
