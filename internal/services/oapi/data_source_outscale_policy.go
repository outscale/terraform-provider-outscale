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
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
)

func DataSourcePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourcePolicyRead,
		Schema: map[string]*schema.Schema{
			"policy_orn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"document": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resources_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policy_default_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_linkable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"orn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modification_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourcePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadPolicyRequest{
		PolicyOrn: d.Get("policy_orn").(string),
	}

	resp, err := client.ReadPolicy(ctx, req, options.WithRetryTimeout(2*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Policy == nil {
		d.SetId("")
		return nil
	}
	policy := resp.Policy
	d.SetId(id.UniqueId())
	if err := d.Set("policy_name", ptr.From(policy.PolicyName)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("policy_id", ptr.From(policy.PolicyId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("path", ptr.From(policy.Path)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("orn", ptr.From(policy.Orn)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("resources_count", ptr.From(policy.ResourcesCount)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_linkable", ptr.From(policy.IsLinkable)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("policy_default_version_id", ptr.From(policy.PolicyDefaultVersionId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", ptr.From(policy.Description)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("creation_date", from.ISO8601(policy.CreationDate)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_modification_date", from.ISO8601(policy.LastModificationDate)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
