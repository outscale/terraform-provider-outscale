package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/outscale/osc-go/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleOAPINetworkInterfacePrivateIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPINetworkInterfacePrivateIPCreate,
		Read:   resourceOutscaleOAPINetworkInterfacePrivateIPRead,
		Delete: resourceOutscaleOAPINetworkInterfacePrivateIPDelete,

		Schema: map[string]*schema.Schema{
			"allow_relink": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"secondary_private_ip_count": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"nic_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_ips": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"primary_private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPINetworkInterfacePrivateIPCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	input := &oapi.LinkPrivateIpsRequest{
		NicId: d.Get("nic_id").(string),
	}

	if v, ok := d.GetOk("allow_relink"); ok {
		input.AllowRelink = v.(bool)
	}

	if v, ok := d.GetOk("secondary_private_ip_count"); ok {
		input.SecondaryPrivateIpCount = int64(v.(int) - 1)
	}

	if v, ok := d.GetOk("private_ips"); ok {
		input.PrivateIps = expandStringValueList(v.([]interface{}))
	}

	var err error
	var resp *oapi.POST_LinkPrivateIpsResponses
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_LinkPrivateIps(*input)
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
		return fmt.Errorf("Failure to assign Private IPs: %s", errString)

	}

	d.SetId(input.NicId)

	return resourceOutscaleOAPINetworkInterfacePrivateIPRead(d, meta)
}

func resourceOutscaleOAPINetworkInterfacePrivateIPRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	interfaceID := d.Get("nic_id").(string)

	req := &oapi.ReadNicsRequest{
		Filters: oapi.FiltersNic{NicIds: []string{interfaceID}},
	}

	var describeResp *oapi.POST_ReadNicsResponses
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		describeResp, err = conn.POST_ReadNics(*req)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var errString string

	if err != nil || describeResp.OK == nil {
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidNetworkInterfaceID.NotFound") {
				// The ENI is gone now, so just remove the attachment from the state
				d.SetId("")
				return nil
			}
			errString = err.Error()
		} else if describeResp.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(describeResp.Code401))
		} else if describeResp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(describeResp.Code400))
		} else if describeResp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(describeResp.Code500))
		}
		return fmt.Errorf("Could not find network interface: %s", errString)

	}

	result := describeResp.OK

	if len(result.Nics) != 1 {
		return fmt.Errorf("Unable to find ENI (%s): %#v", interfaceID, result.Nics)
	}

	eni := result.Nics[0]

	if eni.NicId == "" {
		// Interface is no longer attached, remove from state
		d.SetId("")
		return nil
	}

	var ips []string

	// We need to avoid to store inside private_ips when private IP is the primary IP
	//because the primary can't remove.
	var primaryPrivateID string
	for _, v := range eni.PrivateIps {
		if v.IsPrimary {
			primaryPrivateID = v.PrivateIp
		} else {
			ips = append(ips, v.PrivateIp)
		}
	}

	_, ok := d.GetOk("allow_relink")

	d.Set("allow_relink", ok)
	d.Set("private_ips", ips)
	d.Set("secondary_private_ip_count", len(eni.PrivateIps))
	d.Set("nic_id", eni.NicId)
	d.Set("primary_private_ip", primaryPrivateID)
	d.Set("request_id", result.ResponseContext.RequestId)

	return nil
}

func resourceOutscaleOAPINetworkInterfacePrivateIPDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	input := &oapi.UnlinkPrivateIpsRequest{
		NicId: d.Id(),
	}

	if v, ok := d.GetOk("private_ips"); ok {
		input.PrivateIps = expandStringValueList(v.([]interface{}))
	}

	var err error
	var resp *oapi.POST_UnlinkPrivateIpsResponses
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_UnlinkPrivateIps(*input)
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
		return fmt.Errorf("Failure to unassign Private IPs: %s", errString)

	}

	d.SetId("")

	return nil
}
