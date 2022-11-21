package outscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceApiAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApiAccessPolicyRead,
		Schema: map[string]*schema.Schema{
			"max_access_key_expiration_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"require_trusted_env": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceApiAccessPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).OSCAPI

	req := oscgo.ReadApiAccessPolicyRequest{}

	var resp oscgo.ReadApiAccessPolicyResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.ApiAccessPolicyApi.ReadApiAccessPolicy(context.Background()).ReadApiAccessPolicyRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp.StatusCode, err)
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
