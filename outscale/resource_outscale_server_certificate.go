package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openlyinc/pointy"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/spf13/cast"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleOAPIServerCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIServerCertificateCreate,
		Read:   resourceOutscaleOAPIServerCertificateRead,
		Update: resourceOutscaleOAPIServerCertificateUpdate,
		Delete: resourceOutscaleOAPIServerCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceOutscaleOAPIServerCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.CreateServerCertificateRequest{
		Body:       d.Get("body").(string),
		Name:       d.Get("name").(string),
		PrivateKey: d.Get("private_key").(string),
	}

	if _, ok := d.GetOk("body"); ok == false {
		return fmt.Errorf("[DEBUG] Error 'body' field is require for server certificate creation")
	}

	if _, ok := d.GetOk("private_key"); ok == false {
		return fmt.Errorf("[DEBUG] Error 'private_key' field is require for server certificate creation")
	}

	if v, ok := d.GetOk("chain"); ok {
		req.Chain = pointy.String(v.(string))
	}
	if v, ok := d.GetOk("dry_run"); ok {
		req.DryRun = pointy.Bool(v.(bool))
	}
	if v, ok := d.GetOk("path"); ok {
		req.Path = pointy.String(v.(string))
	}
	var resp oscgo.CreateServerCertificateResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.ServerCertificateApi.CreateServerCertificate(context.Background()).CreateServerCertificateRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("[DEBUG] Error reading Server Certificate: %s", utils.GetErrorResponse(err))
	}

	d.SetId(cast.ToString(resp.ServerCertificate.Id))

	return resourceOutscaleOAPIServerCertificateRead(d, meta)
}

func resourceOutscaleOAPIServerCertificateRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	id := d.Id()

	log.Printf("[DEBUG] Reading Server Certificate id (%s)", id)

	var resp oscgo.ReadServerCertificatesResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.ServerCertificateApi.ReadServerCertificates(context.Background()).ReadServerCertificatesRequest(oscgo.ReadServerCertificatesRequest{}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error reading Server Certificate id (%s)", utils.GetErrorResponse(err))

	}
	if !resp.HasServerCertificates() {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
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

func resourceOutscaleOAPIServerCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	oldName, newName := d.GetChange("name")
	req := oscgo.UpdateServerCertificateRequest{
		Name: oldName.(string),
	}

	if d.HasChange("name") {
		req.NewName = pointy.String(newName.(string))
	}
	if d.HasChange("path") {
		req.NewPath = pointy.String(d.Get("path").(string))
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.ServerCertificateApi.UpdateServerCertificate(context.Background()).UpdateServerCertificateRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("[DEBUG] Error update Server Certificate: %s", utils.GetErrorResponse(err))
	}

	return resourceOutscaleOAPIServerCertificateRead(d, meta)
}

func resourceOutscaleOAPIServerCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.DeleteServerCertificateRequest{
		Name: d.Get("name").(string),
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.ServerCertificateApi.DeleteServerCertificate(context.Background()).DeleteServerCertificateRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	return err
}
