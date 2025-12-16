package outscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/spf13/cast"
)

func DataSourcePolicies() *schema.Resource {
	return &schema.Resource{
		Read: DataSourcePoliciesRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"policy_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
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
				},
			},
		},
	}
}

func DataSourcePoliciesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")

	var err error
	req := oscgo.NewReadPoliciesRequest()
	if filtersOk {
		req.Filters, err = buildPoliciesFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}
	var resp oscgo.ReadPoliciesResponse
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.PolicyApi.ReadPolicies(context.Background()).ReadPoliciesRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}
	policyResp := resp.GetPolicies()
	if len(policyResp) == 0 {
		return fmt.Errorf("Unable to find Policies with fileters: %v", filters.(*schema.Set))
	}
	d.SetId(id.UniqueId())

	policies := make([]map[string]interface{}, len(policyResp))

	for i, v := range policyResp {
		policy := make(map[string]interface{})
		policy["policy_name"] = v.GetPolicyName()
		policy["policy_id"] = v.GetPolicyId()
		policy["path"] = v.GetPath()
		policy["orn"] = v.GetOrn()
		policy["resources_count"] = v.GetResourcesCount()
		policy["is_linkable"] = v.GetIsLinkable()
		policy["policy_default_version_id"] = v.GetPolicyDefaultVersionId()
		policy["description"] = v.GetDescription()
		policy["creation_date"] = v.GetCreationDate()
		policy["last_modification_date"] = v.GetLastModificationDate()
		policies[i] = policy
	}
	return d.Set("policies", policies)
}

func buildPoliciesFilters(set *schema.Set) (*oscgo.ReadPoliciesFilters, error) {
	var filters oscgo.ReadPoliciesFilters
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "only_linked":
			filters.SetOnlyLinked(cast.ToBool(filterValues[0]))
		case "path_prefix":
			filters.SetPathPrefix(filterValues[0])
		case "scope":
			filters.SetScope(filterValues[0])
		default:
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}
