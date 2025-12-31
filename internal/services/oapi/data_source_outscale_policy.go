package oapi

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourcePolicy() *schema.Resource {
	return &schema.Resource{
		Read: DataSourcePolicyRead,
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

func DataSourcePolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	req := oscgo.NewReadPolicyRequest(d.Get("policy_orn").(string))

	var resp oscgo.ReadPolicyResponse
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.PolicyApi.ReadPolicy(context.Background()).ReadPolicyRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	if _, ok := resp.GetPolicyOk(); !ok {
		d.SetId("")
		return nil
	}
	policy := resp.GetPolicy()
	d.SetId(id.UniqueId())
	if err := d.Set("policy_name", policy.GetPolicyName()); err != nil {
		return err
	}
	if err := d.Set("policy_id", policy.GetPolicyId()); err != nil {
		return err
	}
	if err := d.Set("path", policy.GetPath()); err != nil {
		return err
	}
	if err := d.Set("orn", policy.GetOrn()); err != nil {
		return err
	}
	if err := d.Set("resources_count", policy.GetResourcesCount()); err != nil {
		return err
	}
	if err := d.Set("is_linkable", policy.GetIsLinkable()); err != nil {
		return err
	}
	if err := d.Set("policy_default_version_id", policy.GetPolicyDefaultVersionId()); err != nil {
		return err
	}
	if err := d.Set("description", policy.GetDescription()); err != nil {
		return err
	}
	if err := d.Set("creation_date", (policy.GetCreationDate())); err != nil {
		return err
	}
	if err := d.Set("last_modification_date", (policy.GetLastModificationDate())); err != nil {
		return err
	}
	return nil
}
