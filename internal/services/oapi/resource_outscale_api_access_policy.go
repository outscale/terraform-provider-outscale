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

func ResourceOutscaleApiAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleApiAccessPolicyCreate,
		Read:   ResourceOutscaleApiAccessPolicyRead,
		Update: ResourceOutscaleApiAccessPolicyUpdate,
		Delete: ResourceOutscaleApiAccessPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"max_access_key_expiration_seconds": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"require_trusted_env": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceOutscaleApiAccessPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	maxAcc := d.Get("max_access_key_expiration_seconds")
	trustEnv := d.Get("require_trusted_env")

	if trustEnv.(bool) == true && maxAcc == 0 {
		return fmt.Errorf("Error 'max_access_key_expiration_seconds' value must be greater than '0' if 'require_trusted_env' value is 'true'")
	}

	req := oscgo.UpdateApiAccessPolicyRequest{
		MaxAccessKeyExpirationSeconds: int64(maxAcc.(int)),
		RequireTrustedEnv:             trustEnv.(bool),
	}

	var err error
	err = retry.Retry(120*time.Second, func() *retry.RetryError {
		_, httpResp, err := conn.ApiAccessPolicyApi.UpdateApiAccessPolicy(context.Background()).UpdateApiAccessPolicyRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return ResourceOutscaleApiAccessPolicyRead(d, meta)
}

func ResourceOutscaleApiAccessPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	req := oscgo.ReadApiAccessPolicyRequest{}

	var resp oscgo.ReadApiAccessPolicyResponse
	var err error
	err = retry.Retry(120*time.Second, func() *retry.RetryError {
		rp, httpResp, err := conn.ApiAccessPolicyApi.ReadApiAccessPolicy(context.Background()).ReadApiAccessPolicyRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error reading Api Access Policy id (%s)", utils.GetErrorResponse(err))
	}

	if !resp.HasApiAccessPolicy() {
		d.SetId("")
		return fmt.Errorf("Api Access Policy not found")
	}

	policy := resp.GetApiAccessPolicy()
	if err := d.Set("max_access_key_expiration_seconds", policy.GetMaxAccessKeyExpirationSeconds()); err != nil {
		return err
	}
	if err := d.Set("require_trusted_env", policy.GetRequireTrustedEnv()); err != nil {
		return err
	}
	d.SetId(id.UniqueId())
	return nil
}

func ResourceOutscaleApiAccessPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	_, maxAcc := d.GetChange("max_access_key_expiration_seconds")
	_, trustEnv := d.GetChange("require_trusted_env")

	if trustEnv.(bool) == true && maxAcc == 0 {
		return fmt.Errorf("Error 'max_access_key_expiration_seconds' value must be greater than '0' if 'require_trusted_env' value is 'true'")
	}

	req := oscgo.UpdateApiAccessPolicyRequest{
		MaxAccessKeyExpirationSeconds: int64(maxAcc.(int)),
		RequireTrustedEnv:             trustEnv.(bool),
	}

	var err error
	err = retry.Retry(120*time.Second, func() *retry.RetryError {
		_, httpResp, err := conn.ApiAccessPolicyApi.UpdateApiAccessPolicy(context.Background()).UpdateApiAccessPolicyRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return ResourceOutscaleApiAccessPolicyRead(d, meta)
}

func ResourceOutscaleApiAccessPolicyDelete(d *schema.ResourceData, _ interface{}) error {
	d.SetId("")
	return nil
}
