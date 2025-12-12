package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleApiAccessRule() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleApiAccessRuleRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceOutscaleApiAccessRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return fmt.Errorf("filters must be assigned")
	}

	filterParams, err := buildOutscaleApiAccessRuleFilters(filters.(*schema.Set))
	if err != nil {
		return err
	}
	req := oscgo.ReadApiAccessRulesRequest{
		Filters: filterParams,
	}

	var resp oscgo.ReadApiAccessRulesResponse
	err = retry.Retry(120*time.Second, func() *retry.RetryError {
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
	if len(apiAccessRules) > 1 {
		return fmt.Errorf("Your query returned more results. Please change your search criteria and try again")
	}

	accRule := apiAccessRules[0]
	if err := d.Set("api_access_rule_id", accRule.GetApiAccessRuleId()); err != nil {
		return err
	}
	if accRule.HasCaIds() {
		if err := d.Set("ca_ids", accRule.GetCaIds()); err != nil {
			return err
		}
	}

	if accRule.HasCns() {
		if err := d.Set("cns", accRule.GetCns()); err != nil {
			return err
		}
	}
	if accRule.HasIpRanges() {
		if err := d.Set("ip_ranges", accRule.GetIpRanges()); err != nil {
			return err
		}
	}
	if accRule.HasDescription() {
		if err := d.Set("description", accRule.GetDescription()); err != nil {
			return err
		}
	}
	d.SetId(accRule.GetApiAccessRuleId())
	return nil
}

func buildOutscaleApiAccessRuleFilters(set *schema.Set) (*oscgo.FiltersApiAccessRule, error) {
	var filters oscgo.FiltersApiAccessRule
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "api_access_rule_ids":
			filters.SetApiAccessRuleIds(filterValues)
		case "ca_ids":
			filters.SetCaIds(filterValues)
		case "cns":
			filters.SetCns(filterValues)
		case "descriptions":
			filters.SetDescriptions(filterValues)
		case "ip_ranges":
			filters.SetIpRanges(filterValues)
		default:
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}
