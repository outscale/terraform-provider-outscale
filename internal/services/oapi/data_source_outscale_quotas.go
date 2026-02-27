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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleQuotasRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"used_value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"quota_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_collection": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"short_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
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

func DataSourceOutscaleQuotasRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadQuotasRequest{}

	filters, filtersOk := d.GetOk("filter")

	var err error
	if filtersOk {
		req.Filters, err = buildOutscaleQuotaDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadQuotas(ctx, req, options.WithRetryTimeout(120*time.Second))

	var errString string
	if err != nil {
		errString = err.Error()
		return diag.Errorf("error reading quotas type (%s)", errString)
	}

	if resp.QuotaTypes == nil || len(*resp.QuotaTypes) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	quotas := make([]map[string]interface{}, 0)

	for _, quotaType := range *resp.QuotaTypes {
		if quotaType.Quotas == nil || len(*quotaType.Quotas) == 0 {
			return diag.Errorf("no matching quotas found")
		}

		for _, quota := range *quotaType.Quotas {
			quotaMap := make(map[string]interface{})
			if ptr.From(quota.Name) != "" {
				quotaMap["name"] = quota.Name
			}
			if ptr.From(quota.Description) != "" {
				quotaMap["description"] = quota.Description
			}

			quotaMap["max_value"] = quota.MaxValue

			quotaMap["used_value"] = quota.UsedValue

			if ptr.From(quotaType.QuotaType) != "" {
				quotaMap["quota_type"] = quotaType.QuotaType
			}
			if ptr.From(quota.QuotaCollection) != "" {
				quotaMap["quota_collection"] = quota.QuotaCollection
			}
			if ptr.From(quota.ShortDescription) != "" {
				quotaMap["short_description"] = quota.ShortDescription
			}
			if ptr.From(quota.AccountId) != "" {
				quotaMap["account_id"] = quota.AccountId
			}
			quotas = append(quotas, quotaMap)
		}
	}

	if err := d.Set("quotas", quotas); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id.UniqueId())

	return nil
}

func buildOutscaleQuotaDataSourceFilters(set *schema.Set) (*osc.FiltersQuota, error) {
	var filters osc.FiltersQuota
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "quota_types":
			filters.QuotaTypes = &filterValues
		case "quota_names":
			filters.QuotaNames = &filterValues
		case "collections":
			filters.Collections = &filterValues
		case "short_descriptions":
			filters.ShortDescriptions = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
