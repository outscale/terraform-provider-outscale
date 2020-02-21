package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	conn := meta.(*OutscaleClient).OSCAPI

	input := oscgo.LinkPrivateIpsRequest{
		NicId: d.Get("nic_id").(string),
	}

	if v, ok := d.GetOk("allow_relink"); ok {
		input.SetAllowRelink(v.(bool))
	}

	if v, ok := d.GetOk("secondary_private_ip_count"); ok {
		input.SetSecondaryPrivateIpCount(int64(v.(int) - 1))
	}

	if v, ok := d.GetOk("private_ips"); ok {
		input.SetPrivateIps(expandStringValueList(v.([]interface{})))
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.NicApi.LinkPrivateIps(context.Background(), &oscgo.LinkPrivateIpsOpts{LinkPrivateIpsRequest: optional.NewInterface(input)})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
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

	interfaceID := d.Get("nic_id").(string)

	req := oscgo.ReadNicsRequest{
		Filters: &oscgo.FiltersNic{NicIds: &[]string{interfaceID}},
	}

	var resp oscgo.ReadNicsResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.NicApi.ReadNics(context.Background(), &oscgo.ReadNicsOpts{ReadNicsRequest: optional.NewInterface(req)})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidNetworkInterfaceID.NotFound") {
			// The ENI is gone now, so just remove the attachment from the state
			d.SetId("")
			return nil
		}
		errString := err.Error()
		return fmt.Errorf("Could not find network interface: %s", errString)

	}
	if len(resp.GetNics()) != 1 {
		return fmt.Errorf("Unable to find ENI (%s): %#v", interfaceID, resp.GetNics())
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
	for _, v := range eni.GetPrivateIps() {
		if v.GetIsPrimary() {
			primaryPrivateID = v.GetPrivateIp()
		} else {
			ips = append(ips, v.GetPrivateIp())
		}
	}

	_, ok := d.GetOk("allow_relink")

	if err := d.Set("allow_relink", ok); err != nil {
		return err
	}
	if err := d.Set("private_ips", ips); err != nil {
		return err
	}
	if err := d.Set("secondary_private_ip_count", len(eni.GetPrivateIps())); err != nil {
		return err
	}
	if err := d.Set("nic_id", eni.GetNicId()); err != nil {
		return err
	}
	if err := d.Set("primary_private_ip", primaryPrivateID); err != nil {
		return err
	}
	if err := d.Set("request_id", resp.ResponseContext.GetRequestId()); err != nil {
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
		input.SetPrivateIps(expandStringValueList(v.([]interface{})))
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.NicApi.UnlinkPrivateIps(context.Background(), &oscgo.UnlinkPrivateIpsOpts{UnlinkPrivateIpsRequest: optional.NewInterface(input)})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
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
