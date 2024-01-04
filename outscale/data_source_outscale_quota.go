package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOutscaleOAPIQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIQuotaRead,

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

func dataSourceOutscaleOAPIQuotaRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.ReadQuotasRequest{}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters = buildOutscaleOAPIQuotaDataSourceFilters(filters.(*schema.Set))
	}

	var resp oscgo.ReadQuotasResponse
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		var err error
		rp, httpResp, err := conn.QuotaApi.ReadQuotas(context.Background()).ReadQuotasRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		errString := err.Error()
		return fmt.Errorf("[DEBUG] Error reading quotaType (%s)", errString)
	}

	if len(resp.GetQuotaTypes()) == 0 {
		return fmt.Errorf("no matching quotas type found")
	}

	if len(resp.GetQuotaTypes()) > 1 {
		return fmt.Errorf("multiple quotas type matched; use additional constraints to reduce matches to a single quotaType")
	}

	quotaType := resp.GetQuotaTypes()[0]

	d.SetId(resource.UniqueId())
	if err := d.Set("quota_type", quotaType.GetQuotaType()); err != nil {
		return err
	}
	if len(quotaType.GetQuotas()) == 0 {
		return fmt.Errorf("no matching quotas found")
	}

	if len(quotaType.GetQuotas()) > 1 {
		return fmt.Errorf("multiple quotas matched; use additional constraints to reduce matches to a single quotaType")
	}

	quota := quotaType.GetQuotas()[0]

	if err := d.Set("name", quota.GetName()); err != nil {
		return err
	}
	if err := d.Set("description", quota.GetDescription()); err != nil {
		return err
	}
	if err := d.Set("max_value", quota.GetMaxValue()); err != nil {
		return err
	}
	if err := d.Set("used_value", quota.GetUsedValue()); err != nil {
		return err
	}
	if err := d.Set("quota_collection", quota.GetShortDescription()); err != nil {
		return err
	}
	if err := d.Set("short_description", quota.GetShortDescription()); err != nil {
		return err
	}
	if err := d.Set("account_id", quota.GetAccountId()); err != nil {
		return err
	}

	return nil
}
