package oapi

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"
)

func ResourceOutscaleServerCertificate() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleServerCertificateCreate,
		Read:   ResourceOutscaleServerCertificateRead,
		Update: ResourceOutscaleServerCertificateUpdate,
		Delete: ResourceOutscaleServerCertificateDelete,
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

func ResourceOutscaleServerCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutCreate)

	req := oscgo.CreateServerCertificateRequest{
		Body:       d.Get("body").(string),
		Name:       d.Get("name").(string),
		PrivateKey: d.Get("private_key").(string),
	}

	if _, ok := d.GetOk("body"); !ok {
		return fmt.Errorf("error 'body' field is require for server certificate creation")
	}

	if _, ok := d.GetOk("private_key"); !ok {
		return fmt.Errorf("error 'private_key' field is require for server certificate creation")
	}

	if _, ok := d.GetOk("chain"); ok {
		req.SetChain(d.Get("chain").(string))
	}
	if _, ok := d.GetOk("dry_run"); ok {
		req.SetDryRun(d.Get("dry_run").(bool))
	}
	if _, ok := d.GetOk("path"); ok {
		req.SetPath(d.Get("path").(string))
	}
	var resp oscgo.CreateServerCertificateResponse
	err := retry.Retry(timeout, func() *retry.RetryError {
		rp, httpResp, err := conn.ServerCertificateApi.CreateServerCertificate(context.Background()).CreateServerCertificateRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading server certificate: %s", utils.GetErrorResponse(err))
	}

	d.SetId(cast.ToString(resp.ServerCertificate.Id))

	return ResourceOutscaleServerCertificateRead(d, meta)
}

func ResourceOutscaleServerCertificateRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutRead)

	id := d.Id()

	log.Printf("[DEBUG] Reading Server Certificate id (%s)", id)

	var resp oscgo.ReadServerCertificatesResponse
	err := retry.Retry(timeout, func() *retry.RetryError {
		rp, httpResp, err := conn.ServerCertificateApi.ReadServerCertificates(context.Background()).ReadServerCertificatesRequest(oscgo.ReadServerCertificatesRequest{}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading server certificate id (%s)", utils.GetErrorResponse(err))
	}
	if !resp.HasServerCertificates() {
		return ErrNoResults
	}

	if len(resp.GetServerCertificates()) == 0 {
		utils.LogManuallyDeleted("ServerCertificate", d.Id())
		d.SetId("")
		return nil
	}

	var server oscgo.ServerCertificate

	for _, serv := range resp.GetServerCertificates() {
		if serv.GetId() == d.Id() {
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

func ResourceOutscaleServerCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutUpdate)

	oldName, _ := d.GetChange("name")
	req := oscgo.UpdateServerCertificateRequest{
		Name: oldName.(string),
	}

	if d.HasChange("name") {
		req.SetNewName(d.Get("name").(string))
	}
	if d.HasChange("path") {
		req.SetNewPath(d.Get("path").(string))
	}

	err := retry.Retry(timeout, func() *retry.RetryError {
		_, httpResp, err := conn.ServerCertificateApi.UpdateServerCertificate(context.Background()).UpdateServerCertificateRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error update server certificate: %s", utils.GetErrorResponse(err))
	}

	return ResourceOutscaleServerCertificateRead(d, meta)
}

func ResourceOutscaleServerCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutDelete)

	req := oscgo.DeleteServerCertificateRequest{
		Name: d.Get("name").(string),
	}

	err := retry.Retry(timeout, func() *retry.RetryError {
		_, httpResp, err := conn.ServerCertificateApi.DeleteServerCertificate(context.Background()).DeleteServerCertificateRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	return err
}
