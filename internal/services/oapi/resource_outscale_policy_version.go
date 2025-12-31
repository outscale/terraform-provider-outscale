package oapi

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func ResourceOutscalePolicyVersion() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscalePolicyVersionCreate,
		Read:   ResourceOutscalePolicyVersionRead,
		Delete: ResourceOutscalePolicyVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"policy_orn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"document": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"set_as_default": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_version": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"body": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceOutscalePolicyVersionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	pvDocument := d.Get("document").(string)
	req := oscgo.NewCreatePolicyVersionRequest(pvDocument, d.Get("policy_orn").(string))

	if asDefault, ok := d.GetOk("set_as_default"); ok {
		req.SetSetAsDefault(asDefault.(bool))
	}

	var resp oscgo.CreatePolicyVersionResponse
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.PolicyApi.CreatePolicyVersion(context.Background()).CreatePolicyVersionRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(id.UniqueId())
	pVersion := resp.GetPolicyVersion()
	if err := d.Set("version_id", pVersion.GetVersionId()); err != nil {
		return err
	}

	return ResourceOutscalePolicyVersionRead(d, meta)
}

func ResourceOutscalePolicyVersionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	req := oscgo.NewReadPolicyVersionRequest(d.Get("policy_orn").(string), d.Get("version_id").(string))

	var resp oscgo.ReadPolicyVersionResponse
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.PolicyApi.ReadPolicyVersion(context.Background()).ReadPolicyVersionRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	if _, ok := resp.GetPolicyVersionOk(); !ok {
		d.SetId("")
		return nil
	}
	pVersion := resp.GetPolicyVersion()
	if err := d.Set("default_version", pVersion.GetDefaultVersion()); err != nil {
		return err
	}
	if err := d.Set("creation_date", (pVersion.GetCreationDate())); err != nil {
		return err
	}
	if err := d.Set("body", pVersion.GetBody()); err != nil {
		return err
	}
	return nil
}

func ResourceOutscalePolicyVersionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	if d.Get("default_version").(bool) {
		resps, err := readPolicyVersions(d, meta)
		if err != nil {
			return err
		}
		if _, ok := resps.GetPolicyVersionsOk(); !ok {
			d.SetId("")
			return nil
		}
		pVersions := resps.GetPolicyVersions()
		if len(pVersions) <= 1 {
			return fmt.Errorf("cannot delete the default policy version\n It will be deleted with the policy")
		}
		for _, version := range pVersions {
			if version.GetVersionId() == "v1" {
				req := oscgo.NewSetDefaultPolicyVersionRequest(d.Get("policy_orn").(string), version.GetVersionId())
				err = retry.Retry(2*time.Minute, func() *retry.RetryError {
					_, httpResp, err := conn.PolicyApi.SetDefaultPolicyVersion(context.Background()).SetDefaultPolicyVersionRequest(*req).Execute()
					if err != nil {
						return utils.CheckThrottling(httpResp, err)
					}
					return nil
				})
				if err != nil {
					return err
				}
				break
			}
		}
	}

	req := oscgo.NewDeletePolicyVersionRequest(d.Get("policy_orn").(string), d.Get("version_id").(string))

	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		_, httpResp, err := conn.PolicyApi.DeletePolicyVersion(context.Background()).DeletePolicyVersionRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting Outscale Policy version %s: %s", d.Id(), err)
	}

	return nil
}

func readPolicyVersions(d *schema.ResourceData, meta interface{}) (oscgo.ReadPolicyVersionsResponse, error) {
	conn := meta.(*client.OutscaleClient).OSCAPI
	req := oscgo.NewReadPolicyVersionsRequest(d.Get("policy_orn").(string))

	var resp oscgo.ReadPolicyVersionsResponse
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.PolicyApi.ReadPolicyVersions(context.Background()).ReadPolicyVersionsRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	return resp, err
}
