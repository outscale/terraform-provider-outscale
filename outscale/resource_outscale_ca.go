package outscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/spf13/cast"
)

func ResourceOutscaleCa() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleCaCreate,
		Read:   ResourceOutscaleCaRead,
		Update: ResourceOutscaleCaUpdate,
		Delete: ResourceOutscaleCaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ca_pem": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ca_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ca_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceOutscaleCaCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	if _, ok := d.GetOk("ca_pem"); ok == false {
		return fmt.Errorf("[DEBUG] Error 'ca_pem' field is require for certificate authority creation")
	}

	req := oscgo.CreateCaRequest{
		CaPem: d.Get("ca_pem").(string),
	}
	if _, ok := d.GetOk("description"); ok {
		req.SetDescription(d.Get("description").(string))
	}

	var resp oscgo.CreateCaResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.CaApi.CreateCa(context.Background()).CreateCaRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(cast.ToString(resp.Ca.GetCaId()))

	return ResourceOutscaleCaRead(d, meta)
}

func ResourceOutscaleCaRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadCasRequest{
		Filters: &oscgo.FiltersCa{CaIds: &[]string{d.Id()}},
	}

	var resp oscgo.ReadCasResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.CaApi.ReadCas(context.Background()).ReadCasRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error reading certificate authority id (%s)", utils.GetErrorResponse(err))
	}
	if !resp.HasCas() {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}
	if utils.IsResponseEmpty(len(resp.GetCas()), "Ca", d.Id()) {
		d.SetId("")
		return nil
	}

	ca := resp.GetCas()[0]
	if err := d.Set("ca_fingerprint", ca.GetCaFingerprint()); err != nil {
		return err
	}
	if err := d.Set("ca_id", ca.GetCaId()); err != nil {
		return err
	}
	if err := d.Set("description", ca.GetDescription()); err != nil {
		return err
	}
	return nil
}

func ResourceOutscaleCaUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.UpdateCaRequest{
		CaId: d.Get("ca_id").(string),
	}

	if d.HasChange("description") {
		req.SetDescription(d.Get("description").(string))
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.CaApi.UpdateCa(context.Background()).UpdateCaRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return ResourceOutscaleCaRead(d, meta)
}

func ResourceOutscaleCaDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.DeleteCaRequest{
		CaId: d.Get("ca_id").(string),
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.CaApi.DeleteCa(context.Background()).DeleteCaRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	return err
}
