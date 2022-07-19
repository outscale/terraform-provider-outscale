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

func dataSourceOutscaleOAPIApiAccessRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIApiAccessRulesRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"api_access_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_access_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ca_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"cns": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_ranges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPIApiAccessRulesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	req := oscgo.ReadApiAccessRulesRequest{}
	if filtersOk {
		req.Filters = buildOutscaleApiAccessRuleFilters(filters.(*schema.Set))
	}

	var resp oscgo.ReadApiAccessRulesResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.ApiAccessRuleApi.ReadApiAccessRules(context.Background()).ReadApiAccessRulesRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error reading api access rule id (%s)", utils.GetErrorResponse(err))
	}
	apiAccessRules := resp.GetApiAccessRules()[:]
	if len(apiAccessRules) < 1 {
		d.SetId("")
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}
	blockRules := make([]map[string]interface{}, len(apiAccessRules))
	for key, val := range apiAccessRules {
		rule := make(map[string]interface{})

		rule["api_access_rule_id"] = val.GetApiAccessRuleId()
		if val.HasCaIds() {
			rule["ca_ids"] = val.GetCaIds()
		}
		if val.HasCns() {
			rule["cns"] = val.GetCns()
		}
		if val.HasIpRanges() {
			rule["ip_ranges"] = val.GetIpRanges()
		}
		if val.HasDescription() {
			rule["description"] = val.GetDescription()
		}
		blockRules[key] = rule
	}
	d.SetId(resource.UniqueId())
	return d.Set("api_access_rules", blockRules)
}
