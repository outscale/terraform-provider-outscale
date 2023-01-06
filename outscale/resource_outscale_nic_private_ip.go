package outscale

import (
	"context"
	"fmt"
	"net/http"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOutscaleOAPINetworkInterfacePrivateIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPINetworkInterfacePrivateIPCreate,
		Read:   resourceOutscaleOAPINetworkInterfacePrivateIPRead,
		Delete: resourceOutscaleOAPINetworkInterfacePrivateIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
	conn := meta.(*OutscaleClient).OSCAPI

	input := oscgo.LinkPrivateIpsRequest{
		NicId: d.Get("nic_id").(string),
	}

	if v, ok := d.GetOk("allow_relink"); ok {
		input.SetAllowRelink(v.(bool))
	}

	if v, ok := d.GetOk("secondary_private_ip_count"); ok {
		input.SetSecondaryPrivateIpCount(int32(v.(int)))
	}

	if v, ok := d.GetOk("private_ips"); ok {
		input.SetPrivateIps(utils.InterfaceSliceToStringSlice(v.([]interface{})))
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.NicApi.LinkPrivateIps(context.Background()).LinkPrivateIpsRequest(input).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		errString := err.Error()
		return fmt.Errorf("Failure to assign Private IPs: %s", errString)
	}

	d.SetId(input.NicId)

	return resourceOutscaleOAPINetworkInterfacePrivateIPRead(d, meta)
}

func resourceOutscaleOAPINetworkInterfacePrivateIPRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadNicsRequest{
		Filters: &oscgo.FiltersNic{NicIds: &[]string{d.Id()}},
	}

	var resp oscgo.ReadNicsResponse
	var err error
	var statusCode int
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.NicApi.ReadNics(context.Background()).ReadNicsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		statusCode = httpResp.StatusCode
		resp = rp
		return nil
	})
	if err != nil {
		if statusCode == http.StatusNotFound {
			// The ENI is gone now, so just remove the attachment from the state
			d.SetId("")
			return nil
		}
		errString := err.Error()
		return fmt.Errorf("Could not find network interface: %s", errString)

	}
	if utils.IsResponseEmpty(len(resp.GetNics()), "NicPrivateIp", d.Id()) {
		d.SetId("")
		return nil
	}
	eni := resp.GetNics()[0]

	if eni.GetNicId() == "" {
		// Interface is no longer attached, remove from state
		d.SetId("")
		return nil
	}

	var ips []string

	// We need to avoid to store inside private_ips when private IP is the primary IP
	//because the primary can't remove.
	var primaryPrivateID string
	secondary_private_ip_count := 0
	for _, v := range eni.GetPrivateIps() {
		if v.GetIsPrimary() {
			primaryPrivateID = v.GetPrivateIp()
		} else {
			ips = append(ips, v.GetPrivateIp())
			secondary_private_ip_count += 1
		}
	}

	_, ok := d.GetOk("allow_relink")

	if err := d.Set("allow_relink", ok); err != nil {
		return err
	}
	if err := d.Set("private_ips", ips); err != nil {
		return err
	}
	if err := d.Set("secondary_private_ip_count", secondary_private_ip_count); err != nil {
		return err
	}
	if err := d.Set("nic_id", eni.GetNicId()); err != nil {
		return err
	}
	if err := d.Set("primary_private_ip", primaryPrivateID); err != nil {
		return err
	}

	return nil
}

func resourceOutscaleOAPINetworkInterfacePrivateIPDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	input := oscgo.UnlinkPrivateIpsRequest{
		NicId: d.Id(),
	}

	if v, ok := d.GetOk("private_ips"); ok {
		input.SetPrivateIps(utils.InterfaceSliceToStringSlice(v.([]interface{})))
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.NicApi.UnlinkPrivateIps(context.Background()).UnlinkPrivateIpsRequest(input).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		errString := err.Error()
		return fmt.Errorf("Failure to unassign Private IPs: %s", errString)
	}
	d.SetId("")
	return nil
}
