package oapi

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscalePublicIPLink() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscalePublicIPLinkCreate,
		ReadContext:   ResourceOutscalePublicIPLinkRead,
		DeleteContext: ResourceOutscalePublicIPLinkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
		},

		Schema: getOAPIPublicIPLinkSchema(),
	}
}

func ResourceOutscalePublicIPLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutCreate)

	request := osc.LinkPublicIpRequest{}

	if v, ok := d.GetOk("public_ip_id"); ok {
		fmt.Println(v.(string))
		request.PublicIpId = new(v.(string))
	}
	if v, ok := d.GetOk("allow_relink"); ok {
		request.AllowRelink = new(v.(bool))
	}
	if v, ok := d.GetOk("vm_id"); ok {
		request.VmId = new(v.(string))
	}
	if v, ok := d.GetOk("nic_id"); ok {
		request.NicId = new(v.(string))
	}
	if v, ok := d.GetOk("private_ip"); ok {
		request.PrivateIp = new(v.(string))
	}
	if v, ok := d.GetOk("public_ip"); ok {
		request.PublicIp = new(v.(string))
	}

	log.Printf("[DEBUG] EIP association configuration: %#v", request)

	resp, err := client.LinkPublicIp(ctx, request, options.WithRetryTimeout(timeout))
	if err != nil {
		log.Printf("[WARN] ERROR ResourceOutscalePublicIPLinkCreate (%s)", err)
		return diag.FromErr(err)
	}
	// Using validation with request.
	if resp.LinkPublicIpId != nil && len(*resp.LinkPublicIpId) > 0 {
		d.SetId(ptr.From(resp.LinkPublicIpId))
	} else {
		d.SetId(ptr.From(request.PublicIp))
	}

	return ResourceOutscalePublicIPLinkRead(ctx, d, meta)
}

func ResourceOutscalePublicIPLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutRead)

	id := d.Id()
	var request osc.ReadPublicIpsRequest

	if strings.Contains(id, "eipassoc") {
		request = osc.ReadPublicIpsRequest{
			Filters: &osc.FiltersPublicIp{
				LinkPublicIpIds: &[]string{id},
			},
		}
	} else {
		request = osc.ReadPublicIpsRequest{
			Filters: &osc.FiltersPublicIp{
				PublicIps: &[]string{id},
			},
		}
	}

	response, err := client.ReadPublicIps(ctx, request, options.WithRetryTimeout(timeout))
	if err != nil {
		log.Printf("[WARN] ERROR ResourceOutscalePublicIPLinkRead (%s)", err)
		return diag.Errorf("error reading outscale vm public ip %s: %#v", d.Get("public_ip_id").(string), err)
	}
	if response.PublicIps == nil || utils.IsResponseEmpty(len(*response.PublicIps), "PublicIpLink", d.Id()) {
		d.SetId("")
		return nil
	}

	if err := d.Set("tags", FlattenOAPITagsSDK((*response.PublicIps)[0].Tags)); err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(readOutscalePublicIPLink(ctx, d, (*response.PublicIps)[0]))
}

func ResourceOutscalePublicIPLinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutDelete)

	linkID := d.Get("link_public_ip_id")

	opts := osc.UnlinkPublicIpRequest{}
	opts.LinkPublicIpId = new(linkID.(string))

	_, err := client.UnlinkPublicIp(ctx, opts, options.WithRetryTimeout(timeout))
	if err != nil {
		log.Printf("[WARN] ERROR ResourceOutscalePublicIPLinkDelete (%s)", err)
		return diag.Errorf("error deleting elastic ip association: %s", err)
	}

	return nil
}

func readOutscalePublicIPLink(ctx context.Context, d *schema.ResourceData, address osc.PublicIp) error {
	if err := d.Set("vm_id", ptr.From(address.VmId)); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink2 (%s)", err)

		return err
	}
	if err := d.Set("nic_id", ptr.From(address.NicId)); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink3 (%s)", err)

		return err
	}
	if err := d.Set("private_ip", ptr.From(address.PrivateIp)); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink4 (%s)", err)

		return err
	}
	if err := d.Set("public_ip", address.PublicIp); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink (%s)", err)

		return err
	}

	if err := d.Set("link_public_ip_id", ptr.From(address.LinkPublicIpId)); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink (%s)", err)

		return err
	}

	if err := d.Set("nic_account_id", ptr.From(address.NicAccountId)); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink (%s)", err)

		return err
	}

	if err := d.Set("public_ip_id", address.PublicIpId); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink (%s)", err)

		return err
	}

	if err := d.Set("tags", FlattenOAPITagsSDK(address.Tags)); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink TAGS PROBLEME (%s)", err)
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
		"tags": TagsSchemaComputedSDK(),
	}
}
