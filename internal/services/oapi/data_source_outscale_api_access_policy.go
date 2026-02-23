package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
)

func DataSourceOutscaleApiAccessPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleApiAccessPolicyRead,
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

func DataSourceOutscaleApiAccessPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadApiAccessPolicyRequest{}

	resp, err := client.ReadApiAccessPolicy(ctx, req, options.WithRetryTimeout(120*time.Second))
	if err != nil {
		return diag.Errorf("error reading api access policy id (%s)", err)
	}

	if resp.ApiAccessPolicy == nil {
		d.SetId("")
		return diag.FromErr(ErrNoResults)
	}

	policy := resp.ApiAccessPolicy
	if err := d.Set("max_access_key_expiration_seconds", ptr.From(policy.MaxAccessKeyExpirationSeconds)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("require_trusted_env", ptr.From(policy.RequireTrustedEnv)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id.UniqueId())
	return nil
}
