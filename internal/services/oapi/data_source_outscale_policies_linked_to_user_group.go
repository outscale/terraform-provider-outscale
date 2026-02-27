package oapi

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
)

func DataSourcePoliciesLinkedToUserGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourcePoliciesLinkedToUserGroupRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"user_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_id": {
							Type:     schema.TypeString,
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
				},
			},
		},
	}
}

func DataSourcePoliciesLinkedToUserGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadManagedPoliciesLinkedToUserGroupRequest{}
	req.UserGroupName = d.Get("user_group_name").(string)

	var err error
	if filters, filtersOk := d.GetOk("filter"); filtersOk {
		req.Filters, err = buildUserGroupsFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadManagedPoliciesLinkedToUserGroup(ctx, req, options.WithRetryTimeout(2*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.Policies == nil {
		return diag.Errorf("unable to find policies linked to user group")
	}
	policiesResp := resp.Policies
	d.SetId(id.UniqueId())
	policies := make([]map[string]interface{}, len(*policiesResp))
	for i, v := range *policiesResp {
		policy := make(map[string]interface{})
		policy["policy_name"] = v.PolicyName
		policy["policy_id"] = v.PolicyId
		policy["orn"] = v.Orn
		policy["creation_date"] = from.ISO8601(v.CreationDate)
		policy["last_modification_date"] = from.ISO8601(v.LastModificationDate)
		policies[i] = policy
	}

	return diag.FromErr(d.Set("policies", policies))
}
