package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleOAPIPublicIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIPublicIPCreate,
		Read:   resourceOutscaleOAPIPublicIPRead,
		Delete: resourceOutscaleOAPIPublicIPDelete,
		Update: resourceOutscaleOAPIPublicIPUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: getOAPIPublicIPSchema(),
	}
}

func resourceOutscaleOAPIPublicIPCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	allocOpts := oscgo.CreatePublicIpRequest{}

	var resp oscgo.CreatePublicIpResponse
	log.Printf("[DEBUG] EIP create configuration: %#v", allocOpts)
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.PublicIpApi.CreatePublicIp(context.Background()).CreatePublicIpRequest(allocOpts).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("error creating EIP: %s", utils.GetErrorResponse(err))
	}

	allocResp := resp

	log.Printf("[DEBUG] EIP Allocate: %#v", allocResp)

	d.SetId(allocResp.PublicIp.GetPublicIpId())

	//SetTags
	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.(*schema.Set), *allocResp.GetPublicIp().PublicIpId, conn)
		if err != nil {
			return err
		}
	}

	log.Printf("[INFO] EIP ID: %s (placement: %v)", d.Id(), allocResp.GetPublicIp())
	return resourceOutscaleOAPIPublicIPUpdate(d, meta)
}

func resourceOutscaleOAPIPublicIPRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	id := d.Id()

	req := oscgo.ReadPublicIpsRequest{
		Filters: &oscgo.FiltersPublicIp{PublicIpIds: &[]string{id}},
	}

	var response oscgo.ReadPublicIpsResponse
	var err error
	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		resp, httpResp, err := conn.PublicIpApi.ReadPublicIps(context.Background()).ReadPublicIpsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		response = resp
		return nil
	})

	if err != nil {
		if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving EIP: %s", utils.GetErrorResponse(err))
	}
	if utils.IsResponseEmpty(len(response.GetPublicIps()), "PublicIp", d.Id()) {
		d.SetId("")
		return nil
	}

	publicIP := response.GetPublicIps()[0]
	if err := d.Set("link_public_ip_id", publicIP.GetLinkPublicIpId()); err != nil {
		return err
	}
	if err := d.Set("vm_id", publicIP.GetVmId()); err != nil {
		return err
	}
	if err := d.Set("nic_id", publicIP.GetNicId()); err != nil {
		return err
	}
	if err := d.Set("nic_account_id", publicIP.GetNicAccountId()); err != nil {
		return err
	}
	if err := d.Set("private_ip", publicIP.GetPrivateIp()); err != nil {
		return err
	}
	if err := d.Set("public_ip", publicIP.GetPublicIp()); err != nil {
		return err
	}
	if err := d.Set("public_ip_id", publicIP.GetPublicIpId()); err != nil {
		return err
	}

	if err := d.Set("tags", tagsOSCAPIToMap(publicIP.GetTags())); err != nil {
		log.Printf("[WARN] error setting tags for PublicIp(%s): %s", publicIP.GetPublicIp(), err)
	}

	d.SetId(publicIP.GetPublicIpId())

	return nil
}

func resourceOutscaleOAPIPublicIPUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	vVm, okInstance := d.GetOk("vm_id")
	vNic, okInterface := d.GetOk("nic_id")
	idIP := d.Id()
	if okInstance || okInterface {
		assocOpts := oscgo.LinkPublicIpRequest{
			PublicIpId: &idIP,
		}

		if okInterface {
			assocOpts.SetNicId(vNic.(string))
		} else {
			assocOpts.SetVmId(vVm.(string))
		}

		if v, ok := d.GetOk("allow_relink"); ok {
			assocOpts.SetAllowRelink(v.(bool))
		}

		err := resource.Retry(120*time.Second, func() *resource.RetryError {
			var err error
			_, httpResp, err := conn.PublicIpApi.LinkPublicIp(context.Background()).LinkPublicIpRequest(assocOpts).Execute()

			if err != nil {
				if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
					return resource.RetryableError(err)
				}
				return utils.CheckThrottling(httpResp, err)
			}

			return nil
		})

		if err != nil {
			if err := d.Set("vm_id", ""); err != nil {
				return err
			}
			if err := d.Set("nic_id", ""); err != nil {
				return err
			}
			return fmt.Errorf("Failure associating EIP: %s", utils.GetErrorResponse(err))
		}

	}

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}
	return resourceOutscaleOAPIPublicIPRead(d, meta)
}

func unlinkPublicIp(conn *oscgo.APIClient, publicIpId *string) error {
	var err error
	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.PublicIpApi.UnlinkPublicIp(context.Background()).UnlinkPublicIpRequest(oscgo.UnlinkPublicIpRequest{
			LinkPublicIpId: publicIpId,
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	return err
}

func resourceOutscaleOAPIPublicIPDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	if err := resourceOutscaleOAPIPublicIPRead(d, meta); err != nil {
		return err
	}
	if d.Id() == "" {
		return nil
	}

	vInstance, okInstance := d.GetOk("vm_id")
	linkPublicIPID, okAssociationID := d.GetOk("link_public_ip_id")

	if (okInstance && vInstance.(string) != "") || (okAssociationID && linkPublicIPID.(string) != "") {
		log.Printf("[DEBUG] Disassociating EIP: %s", d.Id())
		var err error
		switch resourceOutscaleOAPIPublicIPDomain(d) {
		case "vpc":
			linIpId := d.Get("link_public_ip_id").(string)
			err = unlinkPublicIp(conn, &linIpId)
		case "standard":
			pIP := d.Get("public_ip").(string)
			err = unlinkPublicIp(conn, &pIP)
		}

		if err != nil {
			if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
				return nil
			}
			return err
		}
	}

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		var err error

		idIP := d.Id()
		log.Printf("[DEBUG] EIP release (destroy) address: %v", d.Id())
		_, httpResp, err := conn.PublicIpApi.DeletePublicIp(context.Background()).DeletePublicIpRequest(oscgo.DeletePublicIpRequest{
			PublicIpId: &idIP,
		}).Execute()

		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
}

func getOAPIPublicIPSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"public_ip_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"link_public_ip_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vm_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"nic_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"nic_account_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"private_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"public_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tags": tagsListOAPISchema(),
	}
}

func resourceOutscaleOAPIPublicIPDomain(d *schema.ResourceData) string {
	if v, ok := d.GetOk("placement"); ok {
		return v.(string)
	} else if strings.Contains(d.Id(), "eipalloc") {
		return "vpc"
	}

	return "standard"
}
