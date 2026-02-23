package oapi

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func ResourceOutscalePublicIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscalePublicIPCreate,
		ReadContext:   ResourceOutscalePublicIPRead,
		DeleteContext: ResourceOutscalePublicIPDelete,
		UpdateContext: ResourceOutscalePublicIPUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Update: schema.DefaultTimeout(UpdateDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
		},

		Schema: getOAPIPublicIPSchema(),
	}
}

func ResourceOutscalePublicIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutCreate)

	allocOpts := osc.CreatePublicIpRequest{}

	log.Printf("[DEBUG] EIP create configuration: %#v", allocOpts)
	resp, err := client.CreatePublicIp(ctx, allocOpts, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error creating eip: %s", err)
	}

	allocResp := resp

	log.Printf("[DEBUG] EIP Allocate: %#v", allocResp)

	d.SetId(allocResp.PublicIp.PublicIpId)

	err = createOAPITagsSDK(ctx, client, timeout, d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] EIP ID: %s (placement: %v)", d.Id(), allocResp.PublicIp)
	return ResourceOutscalePublicIPUpdate(ctx, d, meta)
}

func ResourceOutscalePublicIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutRead)

	id := d.Id()

	req := osc.ReadPublicIpsRequest{
		Filters: &osc.FiltersPublicIp{PublicIpIds: &[]string{id}},
	}

	response, err := client.ReadPublicIps(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error retrieving eip: %s", err)
	}
	if response.PublicIps == nil || utils.IsResponseEmpty(len(*response.PublicIps), "PublicIp", d.Id()) {
		d.SetId("")
		return nil
	}

	publicIP := (*response.PublicIps)[0]
	if err := d.Set("link_public_ip_id", ptr.From(publicIP.LinkPublicIpId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vm_id", ptr.From(publicIP.VmId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nic_id", ptr.From(publicIP.NicId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nic_account_id", ptr.From(publicIP.NicAccountId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("private_ip", ptr.From(publicIP.PrivateIp)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("public_ip", publicIP.PublicIp); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("public_ip_id", publicIP.PublicIpId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tags", FlattenOAPITagsSDK(publicIP.Tags)); err != nil {
		log.Printf("[WARN] error setting tags for PublicIp(%s): %s", publicIP.PublicIp, err)
	}

	d.SetId(publicIP.PublicIpId)

	return nil
}

func ResourceOutscalePublicIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutUpdate)

	vVm, okInstance := d.GetOk("vm_id")
	vNic, okInterface := d.GetOk("nic_id")
	idIP := d.Id()
	if okInstance || okInterface {
		assocOpts := osc.LinkPublicIpRequest{
			PublicIpId: &idIP,
		}

		if okInterface {
			assocOpts.NicId = new(vNic.(string))
		} else {
			assocOpts.VmId = new(vVm.(string))
		}

		if v, ok := d.GetOk("allow_relink"); ok {
			assocOpts.AllowRelink = new(v.(bool))
		}

		_, err := client.LinkPublicIp(ctx, assocOpts, options.WithRetryTimeout(timeout))
		if err != nil {
			if err := d.Set("vm_id", ""); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("nic_id", ""); err != nil {
				return diag.FromErr(err)
			}
			return diag.Errorf("failure associating eip: %s", err)
		}

	}

	if err := updateOAPITagsSDK(ctx, client, timeout, d); err != nil {
		return diag.FromErr(err)
	}
	return ResourceOutscalePublicIPRead(ctx, d, meta)
}

func unlinkPublicIp(ctx context.Context, client *osc.Client, publicIpId *string, timeout time.Duration) error {
	_, err := client.UnlinkPublicIp(ctx, osc.UnlinkPublicIpRequest{
		LinkPublicIpId: publicIpId,
	}, options.WithRetryTimeout(timeout))
	return err
}

func ResourceOutscalePublicIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutDelete)

	if err := ResourceOutscalePublicIPRead(ctx, d, meta); err != nil {
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
		switch ResourceOutscalePublicIPDomain(d) {
		case "vpc":
			linIpId := d.Get("link_public_ip_id").(string)
			err = unlinkPublicIp(ctx, client, &linIpId, timeout)
		case "standard":
			pIP := d.Get("public_ip").(string)
			err = unlinkPublicIp(ctx, client, &pIP, timeout)
		}

		if err != nil {
			if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
				return nil
			}
			return diag.FromErr(err)
		}
	}

	idIP := d.Id()
	log.Printf("[DEBUG] EIP release (destroy) address: %v", d.Id())
	_, err := client.DeletePublicIp(ctx, osc.DeletePublicIpRequest{
		PublicIpId: &idIP,
	})

	return diag.FromErr(err)
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
		"tags": TagsSchemaSDK(),
	}
}

func ResourceOutscalePublicIPDomain(d *schema.ResourceData) string {
	if v, ok := d.GetOk("placement"); ok {
		return v.(string)
	} else if strings.Contains(d.Id(), "eipalloc") {
		return "vpc"
	}

	return "standard"
}
