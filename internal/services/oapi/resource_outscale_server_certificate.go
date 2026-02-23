package oapi

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func ResourceOutscaleServerCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleServerCertificateCreate,
		ReadContext:   ResourceOutscaleServerCertificateRead,
		UpdateContext: ResourceOutscaleServerCertificateUpdate,
		DeleteContext: ResourceOutscaleServerCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Update: schema.DefaultTimeout(UpdateDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"body": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"chain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dry_run": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"orn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"upload_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceOutscaleServerCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutCreate)

	req := osc.CreateServerCertificateRequest{
		Body:       d.Get("body").(string),
		Name:       d.Get("name").(string),
		PrivateKey: d.Get("private_key").(string),
	}

	if _, ok := d.GetOk("body"); !ok {
		return diag.Errorf("error 'body' field is require for server certificate creation")
	}

	if _, ok := d.GetOk("private_key"); !ok {
		return diag.Errorf("error 'private_key' field is require for server certificate creation")
	}

	if _, ok := d.GetOk("chain"); ok {
		req.Chain = new(d.Get("chain").(string))
	}
	if _, ok := d.GetOk("dry_run"); ok {
		req.DryRun = new(d.Get("dry_run").(bool))
	}
	if _, ok := d.GetOk("path"); ok {
		req.Path = new(d.Get("path").(string))
	}
	resp, err := client.CreateServerCertificate(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error reading server certificate: %s", err)
	}

	d.SetId(cast.ToString(resp.ServerCertificate.Id))

	return ResourceOutscaleServerCertificateRead(ctx, d, meta)
}

func ResourceOutscaleServerCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutRead)

	id := d.Id()

	log.Printf("[DEBUG] Reading Server Certificate id (%s)", id)

	resp, err := client.ReadServerCertificates(ctx, osc.ReadServerCertificatesRequest{}, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error reading server certificate id (%s)", err)
	}
	if resp.ServerCertificates == nil {
		return diag.FromErr(ErrNoResults)
	}

	if len(*resp.ServerCertificates) == 0 {
		utils.LogManuallyDeleted("ServerCertificate", d.Id())
		d.SetId("")
		return nil
	}

	server, ok := lo.Find(*resp.ServerCertificates, func(s osc.ServerCertificate) bool {
		return ptr.From(s.Id) == d.Id()
	})
	if !ok {
		utils.LogManuallyDeleted("ServerCertificate", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("expiration_date", from.ISO8601(server.ExpirationDate))
	d.Set("name", ptr.From(server.Name))
	d.Set("orn", ptr.From(server.Orn))
	d.Set("path", ptr.From(server.Path))
	d.Set("upload_date", from.ISO8601(server.UploadDate))

	return nil
}

func ResourceOutscaleServerCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutUpdate)

	oldName, _ := d.GetChange("name")
	req := osc.UpdateServerCertificateRequest{
		Name: oldName.(string),
	}

	if d.HasChange("name") {
		req.NewName = new(d.Get("name").(string))
	}
	if d.HasChange("path") {
		req.NewPath = new(d.Get("path").(string))
	}

	_, err := client.UpdateServerCertificate(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error update server certificate: %s", err)
	}

	return ResourceOutscaleServerCertificateRead(ctx, d, meta)
}

func ResourceOutscaleServerCertificateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutDelete)

	req := osc.DeleteServerCertificateRequest{
		Name: d.Get("name").(string),
	}

	_, err := client.DeleteServerCertificate(ctx, req, options.WithRetryTimeout(timeout))

	return diag.FromErr(err)
}
