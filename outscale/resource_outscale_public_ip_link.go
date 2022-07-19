package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOutscaleOAPIPublicIPLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIPublicIPLinkCreate,
		Read:   resourceOutscaleOAPIPublicIPLinkRead,
		Delete: resourceOutscaleOAPIPublicIPLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: getOAPIPublicIPLinkSchema(),
	}
}

func resourceOutscaleOAPIPublicIPLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	request := oscgo.LinkPublicIpRequest{}

	if v, ok := d.GetOk("public_ip_id"); ok {
		fmt.Println(v.(string))
		request.SetPublicIpId(v.(string))
	}
	if v, ok := d.GetOk("allow_relink"); ok {
		request.SetAllowRelink(v.(bool))
	}
	if v, ok := d.GetOk("vm_id"); ok {
		request.SetVmId(v.(string))
	}
	if v, ok := d.GetOk("nic_id"); ok {
		request.SetNicId(v.(string))
	}
	if v, ok := d.GetOk("private_ip"); ok {
		request.SetPrivateIp(v.(string))
	}
	if v, ok := d.GetOk("public_ip"); ok {
		request.SetPublicIp(v.(string))
	}

	log.Printf("[DEBUG] EIP association configuration: %#v", request)

	var resp oscgo.LinkPublicIpResponse
	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.PublicIpApi.LinkPublicIp(context.Background()).LinkPublicIpRequest(request).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		log.Printf("[WARN] ERROR resourceOutscaleOAPIPublicIPLinkCreate (%s)", err)
		return err
	}
	//Using validation with request.
	if resp.GetLinkPublicIpId() != "" && len(resp.GetLinkPublicIpId()) > 0 {
		d.SetId(resp.GetLinkPublicIpId())
	} else {
		d.SetId(request.GetPublicIp())
	}

	return resourceOutscaleOAPIPublicIPLinkRead(d, meta)
}

func resourceOutscaleOAPIPublicIPLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	id := d.Id()
	var request oscgo.ReadPublicIpsRequest

	if strings.Contains(id, "eipassoc") {
		request = oscgo.ReadPublicIpsRequest{
			Filters: &oscgo.FiltersPublicIp{
				LinkPublicIpIds: &[]string{id},
			},
		}
	} else {
		request = oscgo.ReadPublicIpsRequest{
			Filters: &oscgo.FiltersPublicIp{
				PublicIps: &[]string{id},
			},
		}
	}

	var response oscgo.ReadPublicIpsResponse
	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		resp, httpResp, err := conn.PublicIpApi.ReadPublicIps(context.Background()).ReadPublicIpsRequest(request).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		response = resp
		return nil
	})

	if err != nil {
		log.Printf("[WARN] ERROR resourceOutscaleOAPIPublicIPLinkRead (%s)", err)
		return fmt.Errorf("Error reading Outscale VM Public IP %s: %#v", d.Get("public_ip_id").(string), err)
	}
	if utils.IsResponseEmpty(len(response.GetPublicIps()), "PublicIpLink", d.Id()) {
		d.SetId("")
		return nil
	}

	if err := d.Set("tags", getOapiTagSet(response.GetPublicIps()[0].Tags)); err != nil {
		return err
	}

	return readOutscaleOAPIPublicIPLink(d, &response.GetPublicIps()[0])
}

func resourceOutscaleOAPIPublicIPLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	linkID := d.Get("link_public_ip_id")

	opts := oscgo.UnlinkPublicIpRequest{}
	opts.SetLinkPublicIpId(linkID.(string))

	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.PublicIpApi.UnlinkPublicIp(context.Background()).UnlinkPublicIpRequest(opts).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	if err != nil {
		log.Printf("[WARN] ERROR resourceOutscaleOAPIPublicIPLinkDelete (%s)", err)
		return fmt.Errorf("Error deleting Elastic IP association: %s", err)
	}

	return nil
}

func readOutscaleOAPIPublicIPLink(d *schema.ResourceData, address *oscgo.PublicIp) error {
	// if err := d.Set("public_ip_id", address.ReservationId); err != nil {
	// 	log.Printf("[WARN] ERROR readOutscalePublicIPLink1 (%s)", err)

	// 	return err
	// }
	if err := d.Set("vm_id", address.GetVmId()); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink2 (%s)", err)

		return err
	}
	if err := d.Set("nic_id", address.GetNicId()); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink3 (%s)", err)

		return err
	}
	if err := d.Set("private_ip", address.GetPrivateIp()); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink4 (%s)", err)

		return err
	}
	if err := d.Set("public_ip", address.GetPublicIp()); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleOAPIPublicIPLink (%s)", err)

		return err
	}

	if err := d.Set("link_public_ip_id", address.GetLinkPublicIpId()); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleOAPIPublicIPLink (%s)", err)

		return err
	}

	if err := d.Set("nic_account_id", address.GetNicAccountId()); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleOAPIPublicIPLink (%s)", err)

		return err
	}

	if err := d.Set("public_ip_id", address.GetPublicIpId()); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleOAPIPublicIPLink (%s)", err)

		return err
	}

	if err := d.Set("tags", tagsOSCAPIToMap(address.GetTags())); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleOAPIPublicIPLink TAGS PROBLEME (%s)", err)
	}

	return nil
}

func getOAPIPublicIPLinkSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"public_ip_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		"allow_relink": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"vm_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		"nic_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		"private_ip": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		"public_ip": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		"link_public_ip_id": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"nic_account_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tags": tagsOAPIListSchemaComputed(),
	}
}
