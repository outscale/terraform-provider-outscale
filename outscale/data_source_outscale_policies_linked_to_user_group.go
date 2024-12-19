package outscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func DataSourcePoliciesLinkedToUserGroup() *schema.Resource {
	return &schema.Resource{
		Read: DataSourcePoliciesLinkedToUserGroupRead,
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

func DataSourcePoliciesLinkedToUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.ReadManagedPoliciesLinkedToUserGroupRequest{}
	req.SetUserGroupName(d.Get("user_group_name").(string))

	if filters, filtersOk := d.GetOk("filter"); filtersOk {
		filterReq := buildUserGroupsFilters(filters.(*schema.Set))
		req.SetFilters(*filterReq)
	}

	var resp oscgo.ReadManagedPoliciesLinkedToUserGroupResponse
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.PolicyApi.ReadManagedPoliciesLinkedToUserGroup(context.Background()).ReadManagedPoliciesLinkedToUserGroupRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}
	if _, ok := resp.GetPoliciesOk(); !ok {
		return fmt.Errorf("unable to find policies linked to user group")
	}
	policiesResp := resp.GetPolicies()
	d.SetId(resource.UniqueId())
	policies := make([]map[string]interface{}, len(policiesResp))
	for i, v := range policiesResp {
		policy := make(map[string]interface{})
		policy["policy_name"] = v.GetPolicyName()
		policy["policy_id"] = v.GetPolicyId()
		policy["orn"] = v.GetOrn()
		policy["creation_date"] = v.GetCreationDate()
		policy["last_modification_date"] = v.GetLastModificationDate()
		policies[i] = policy
	}

	return d.Set("policies", policies)
}
