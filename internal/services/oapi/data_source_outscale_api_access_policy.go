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

func DataSourceOutscaleApiAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleApiAccessPolicyRead,
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

func DataSourceOutscaleApiAccessPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	req := oscgo.ReadApiAccessPolicyRequest{}

	var resp oscgo.ReadApiAccessPolicyResponse
	err := retry.Retry(120*time.Second, func() *retry.RetryError {
		rp, httpResp, err := conn.ApiAccessPolicyApi.ReadApiAccessPolicy(context.Background()).ReadApiAccessPolicyRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("error reading api access policy id (%s)", utils.GetErrorResponse(err))
	}

	if !resp.HasApiAccessPolicy() {
		d.SetId("")
		return ErrNoResults
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
