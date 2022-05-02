package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleOAPIQuotas() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIQuotasRead,

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

func dataSourceOutscaleOAPIQuotasRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadQuotasRequest{}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters = buildOutscaleOAPIQuotaDataSourceFilters(filters.(*schema.Set))
	}

	var resp oscgo.ReadQuotasResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, _, err = conn.QuotaApi.ReadQuotas(context.Background()).ReadQuotasRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(err)
		}
		return nil
	})

	var errString string
	if err != nil {
		errString = err.Error()
		return fmt.Errorf("[DEBUG] Error reading Quotas type (%s)", errString)
	}

	if len(resp.GetQuotaTypes()) == 0 {
		return fmt.Errorf("no matching Quotas type found")
	}

	quotas := make([]map[string]interface{}, 0)

	for _, quotaType := range resp.GetQuotaTypes() {
		if len(quotaType.GetQuotas()) == 0 {
			return fmt.Errorf("no matching quotas found")
		}

		for _, quota := range quotaType.GetQuotas() {
			quotaMap := make(map[string]interface{})
			if quota.GetName() != "" {
				quotaMap["name"] = quota.GetName()
			}
			if quota.GetDescription() != "" {
				quotaMap["description"] = quota.GetDescription()
			}

			quotaMap["max_value"] = quota.GetMaxValue()

			quotaMap["used_value"] = quota.GetUsedValue()

			if quotaType.GetQuotaType() != "" {
				quotaMap["quota_type"] = quotaType.GetQuotaType()
			}
			if quota.GetQuotaCollection() != "" {
				quotaMap["quota_collection"] = quota.GetQuotaCollection()
			}
			if quota.GetShortDescription() != "" {
				quotaMap["short_description"] = quota.GetShortDescription()
			}
			if quota.GetAccountId() != "" {
				quotaMap["account_id"] = quota.GetAccountId()
			}
			quotas = append(quotas, quotaMap)
		}
	}

	if err := d.Set("quotas", quotas); err != nil {
		return err
	}
	d.SetId(resource.UniqueId())

	return nil
}
