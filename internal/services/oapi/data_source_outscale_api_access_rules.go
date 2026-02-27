package oapi

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
)

func DataSourceOutscaleApiAccessRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleApiAccessRulesRead,
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

func DataSourceOutscaleApiAccessRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	req := osc.ReadApiAccessRulesRequest{}
	if filtersOk {
		filterParams, err := buildOutscaleApiAccessRuleFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
		req.Filters = filterParams
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
	blockRules := make([]map[string]interface{}, len(apiAccessRules))
	for key, val := range apiAccessRules {
		rule := make(map[string]interface{})

		rule["api_access_rule_id"] = val.ApiAccessRuleId
		if val.CaIds != nil {
			rule["ca_ids"] = val.CaIds
		}
		if val.Cns != nil {
			rule["cns"] = val.Cns
		}
		if val.IpRanges != nil {
			rule["ip_ranges"] = val.IpRanges
		}
		if val.Description != nil {
			rule["description"] = val.Description
		}
		blockRules[key] = rule
	}
	d.SetId(id.UniqueId())

	return diag.FromErr(d.Set("api_access_rules", blockRules))
}
