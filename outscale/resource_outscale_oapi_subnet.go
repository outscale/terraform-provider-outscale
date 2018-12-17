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

func resourceOutscaleOAPISubNet() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPISubNetCreate,
		Read:   resourceOutscaleOAPISubNetRead,
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
	conn := meta.(*OutscaleClient).OAPI
	req := &oapi.CreateSubnetRequest{
		IpRange: d.Get("ip_range").(string),
		NetId:   d.Get("net_id").(string),
	}
	if a, aok := d.GetOk("subregion_name"); aok {
		req.SubregionName = a.(string)
	}

	var resp *oapi.POST_CreateSubnetResponses
	var err error
	err = resource.Retry(40*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_CreateSubnet(*req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("[INFO] Request limit exceeded")
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
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
		}

		return fmt.Errorf("[DEBUG] Error creating Subnet (%s)", errString)
	}

	result := resp.OK

	d.SetId(result.Subnet.SubnetId)

	return resourceOutscaleOAPISubNetRead(d, meta)
}

//Read SubNet

func resourceOutscaleOAPISubNetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	id := d.Id()

	log.Printf("[DEBUG] Reading Subnet(%s)", id)

	req := &oapi.ReadSubnetsRequest{
		Filters: oapi.FiltersSubnet{
			SubnetIds: []string{id},
		},
	}

	var resp *oapi.POST_ReadSubnetsResponses
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.POST_ReadSubnets(*req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
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
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
		}

		return fmt.Errorf("[DEBUG] Error reading Subnet (%s)", errString)
	}

	response := resp.OK

	log.Printf("[DEBUG] Setting Subnet (%s)", err)

	d.Set("request_id", response.ResponseContext.RequestId)
	return readOutscaleOAPISubNet(d, &response.Subnets[0])
}

func resourceOutscaleOAPISubNetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	id := d.Id()
	log.Printf("[DEBUG] Deleting Subnet (%s)", id)

	req := &oapi.DeleteSubnetRequest{
		SubnetId: id,
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.POST_DeleteSubnet(*req)

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

	return nil
}

func readOutscaleOAPISubNet(d *schema.ResourceData, subnet *oapi.Subnet) error {

	if err := d.Set("subregion_name", subnet.SubregionName); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet1 (%s)", err)

		return err
	}
	if err := d.Set("available_ips_count", subnet.AvailableIpsCount); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet2 (%s)", err)

		return err
	}
	if err := d.Set("ip_range", subnet.IpRange); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet3 (%s)", err)

		return err
	}
	if err := d.Set("state", subnet.State); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet4 (%s)", err)

		return err
	}
	if err := d.Set("subnet_id", subnet.SubnetId); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet5 (%s)", err)

		return err
	}

	if err := d.Set("net_id", subnet.NetId); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleSubNet6 (%s)", err)

		return err
	}

	return d.Set("tags", tagsOAPIToMap(subnet.Tags))
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
		"tags": dataSourceTagsSchema(),
	}
}
