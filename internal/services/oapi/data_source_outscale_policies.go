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
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"
)

func DataSourcePolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourcePoliciesRead,
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

func DataSourcePoliciesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")

	var err error
	req := osc.ReadPoliciesRequest{}
	if filtersOk {
		req.Filters, err = buildPoliciesFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	resp, err := client.ReadPolicies(ctx, req, options.WithRetryTimeout(2*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}
	policyResp := resp.Policies
	if policyResp == nil || len(*policyResp) == 0 {
		return diag.FromErr(ErrNoResults)
	}
	d.SetId(id.UniqueId())

	policies := make([]map[string]interface{}, len(*policyResp))

	for i, v := range *policyResp {
		policy := make(map[string]interface{})
		policy["policy_name"] = v.PolicyName
		policy["policy_id"] = v.PolicyId
		policy["path"] = v.Path
		policy["orn"] = v.Orn
		policy["resources_count"] = v.ResourcesCount
		policy["is_linkable"] = v.IsLinkable
		policy["policy_default_version_id"] = v.PolicyDefaultVersionId
		policy["description"] = v.Description
		policy["creation_date"] = from.ISO8601(v.CreationDate)
		policy["last_modification_date"] = from.ISO8601(v.LastModificationDate)
		policies[i] = policy
	}
	return diag.FromErr(d.Set("policies", policies))
}

func buildPoliciesFilters(set *schema.Set) (*osc.ReadPoliciesFilters, error) {
	var filters osc.ReadPoliciesFilters
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "only_linked":
			filters.OnlyLinked = new(cast.ToBool(filterValues[0]))
		case "path_prefix":
			filters.PathPrefix = &filterValues[0]
		case "scope":
			filters.Scope = new(osc.ReadPoliciesFiltersScope(filterValues[0]))
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
