package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleQuota() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleQuotaRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceOutscaleQuotaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	if err != nil {
		errString := err.Error()
		return diag.Errorf("error reading quotatype (%s)", errString)
	}

	if resp.QuotaTypes == nil || len(*resp.QuotaTypes) == 0 {
		return diag.FromErr(ErrNoResults)
	}
	if len(*resp.QuotaTypes) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	quotaType := (*resp.QuotaTypes)[0]

	d.SetId(id.UniqueId())
	if err := d.Set("quota_type", ptr.From(quotaType.QuotaType)); err != nil {
		return diag.FromErr(err)
	}

	if quotaType.Quotas == nil || len(*quotaType.Quotas) == 0 {
		return diag.FromErr(ErrNoResults)
	}
	if len(*quotaType.Quotas) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	quota := (*quotaType.Quotas)[0]

	if err := d.Set("name", ptr.From(quota.Name)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", ptr.From(quota.Description)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("max_value", ptr.From(quota.MaxValue)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("used_value", ptr.From(quota.UsedValue)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("quota_collection", ptr.From(quota.ShortDescription)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("short_description", ptr.From(quota.ShortDescription)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("account_id", ptr.From(quota.AccountId)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
