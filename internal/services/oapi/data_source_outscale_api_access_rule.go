package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleApiAccessRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleApiAccessRuleRead,
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

func DataSourceOutscaleApiAccessRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return diag.FromErr(ErrFilterRequired)
	}

	filterParams, err := buildOutscaleApiAccessRuleFilters(filters.(*schema.Set))
	if err != nil {
		return diag.FromErr(err)
	}
	req := osc.ReadApiAccessRulesRequest{
		Filters: filterParams,
	}

	resp, err := client.ReadApiAccessRules(ctx, req, options.WithRetryTimeout(120*time.Second))
	if err != nil {
		return diag.Errorf("error reading api access rule id (%s)", err)
	}
	apiAccessRules := ptr.From(resp.ApiAccessRules)[:]
	if len(apiAccessRules) < 1 {
		d.SetId("")
		return diag.FromErr(ErrNoResults)
	}
	if len(apiAccessRules) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	accRule := apiAccessRules[0]
	if err := d.Set("api_access_rule_id", ptr.From(accRule.ApiAccessRuleId)); err != nil {
		return diag.FromErr(err)
	}
	if accRule.CaIds != nil {
		if err := d.Set("ca_ids", ptr.From(accRule.CaIds)); err != nil {
			return diag.FromErr(err)
		}
	}

	if accRule.Cns != nil {
		if err := d.Set("cns", ptr.From(accRule.Cns)); err != nil {
			return diag.FromErr(err)
		}
	}
	if accRule.IpRanges != nil {
		if err := d.Set("ip_ranges", ptr.From(accRule.IpRanges)); err != nil {
			return diag.FromErr(err)
		}
	}
	if accRule.Description != nil {
		if err := d.Set("description", ptr.From(accRule.Description)); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(ptr.From(accRule.ApiAccessRuleId))
	return nil
}

func buildOutscaleApiAccessRuleFilters(set *schema.Set) (*osc.FiltersApiAccessRule, error) {
	var filters osc.FiltersApiAccessRule
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "api_access_rule_ids":
			filters.ApiAccessRuleIds = &filterValues
		case "ca_ids":
			filters.CaIds = &filterValues
		case "cns":
			filters.Cns = &filterValues
		case "descriptions":
			filters.Descriptions = &filterValues
		case "ip_ranges":
			filters.IpRanges = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
