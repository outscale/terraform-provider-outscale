package oapi

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
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
		return diag.FromErr(fmt.Errorf("error 'body' field is require for server certificate creation"))
	}

	if _, ok := d.GetOk("private_key"); !ok {
		return diag.FromErr(fmt.Errorf("error 'private_key' field is require for server certificate creation"))
	}

	if _, ok := d.GetOk("chain"); ok {
		req.Chain = ptr.To(d.Get("chain").(string))
	}
	if _, ok := d.GetOk("dry_run"); ok {
		req.DryRun = ptr.To(d.Get("dry_run").(bool))
	}
	if _, ok := d.GetOk("path"); ok {
		req.Path = ptr.To(d.Get("path").(string))
	}
	resp, err := client.CreateServerCertificate(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading server certificate: %s", utils.GetErrorResponse(err)))
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
		return diag.FromErr(fmt.Errorf("error reading server certificate id (%s)", utils.GetErrorResponse(err)))
	}
	if resp.ServerCertificates == nil {
		return diag.FromErr(ErrNoResults)
	}

	if len(*resp.ServerCertificates) == 0 {
		utils.LogManuallyDeleted("ServerCertificate", d.Id())
		d.SetId("")
		return nil
	}

	var server osc.ServerCertificate

	for _, serv := range *resp.ServerCertificates {
		if serv.Id == ptr.To(d.Id()) {
			server = serv
		}
	}

	d.Set("expiration_date", server.ExpirationDate)
	d.Set("name", server.Name)
	d.Set("orn", server.Orn)
	d.Set("path", server.Path)
	d.Set("upload_date", server.UploadDate)

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
		req.NewName = ptr.To(d.Get("name").(string))
	}
	if d.HasChange("path") {
		req.NewPath = ptr.To(d.Get("path").(string))
	}

	_, err := client.UpdateServerCertificate(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error update server certificate: %s", utils.GetErrorResponse(err)))
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
