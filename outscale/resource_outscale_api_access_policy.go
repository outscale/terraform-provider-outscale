package outscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleOAPIApiAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIApiAccessPolicyCreate,
		Read:   resourceOutscaleOAPIApiAccessPolicyRead,
		Update: resourceOutscaleOAPIApiAccessPolicyUpdate,
		Delete: resourceOutscaleOAPIApiAccessPolicyDelete,
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

func resourceOutscaleOAPIApiAccessPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

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
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.ApiAccessPolicyApi.UpdateApiAccessPolicy(context.Background()).UpdateApiAccessPolicyRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return resourceOutscaleOAPIApiAccessPolicyRead(d, meta)
}

func resourceOutscaleOAPIApiAccessPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadApiAccessPolicyRequest{}

	var resp oscgo.ReadApiAccessPolicyResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
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
	d.SetId(resource.UniqueId())
	return nil
}

func resourceOutscaleOAPIApiAccessPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

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
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.ApiAccessPolicyApi.UpdateApiAccessPolicy(context.Background()).UpdateApiAccessPolicyRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return resourceOutscaleOAPIApiAccessPolicyRead(d, meta)
}

func resourceOutscaleOAPIApiAccessPolicyDelete(d *schema.ResourceData, _ interface{}) error {
	d.SetId("")
	return nil
}
