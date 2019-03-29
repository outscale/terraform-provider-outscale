package outscale

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleQuotaRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"quota_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"quota_set": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_quota_value": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"used_quota_value": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"reference": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleQuotaRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	quota, quotaOk := d.GetOk("quota_name")

	if !filtersOk && !quotaOk {
		return fmt.Errorf("One of quota_name or filters must be assigned")
	}

	params := &fcu.DescribeQuotasInput{}

	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}

	if quotaOk {
		params.QuotaName = aws.StringSlice([]string{quota.(string)})
	}

	var resp *fcu.DescribeQuotasOutput
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.VM.DescribeQuotas(params)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})

	if err != nil {
		return err
	}
	if resp == nil || len(resp.ReferenceQuotaSet) == 0 {
		return fmt.Errorf("no matching quotas list found; the quotas list ID or name may be invalid or not exist in the current region")
	}

	if len(resp.ReferenceQuotaSet) > 1 {
		return fmt.Errorf("multiple Quotas matched; use additional constraints to reduce matches to a single Quotas")
	}

	pl := resp.ReferenceQuotaSet[0]

	d.SetId(resource.UniqueId())

	quotas := make([]map[string]interface{}, len(pl.QuotaSet))
	for k, v := range pl.QuotaSet {
		quota := make(map[string]interface{})
		quota["description"] = aws.StringValue(v.Description)
		quota["display_name"] = aws.StringValue(v.DisplayName)
		quota["group_name"] = aws.StringValue(v.GroupName)
		i, err := strconv.Atoi(*v.MaxQuotaValue)
		if err != nil {
			return err
		}
		quota["max_quota_value"] = i
		quota["name"] = aws.StringValue(v.Name)
		quota["owner_id"] = aws.StringValue(v.OwnerId)
		i2, err := strconv.Atoi(*v.MaxQuotaValue)
		if err != nil {
			return err
		}
		quota["used_quota_value"] = i2
		quotas[k] = quota
	}

	if err := d.Set("quota_set", quotas); err != nil {
		return err
	}
	d.Set("reference", pl.Reference)
	d.Set("request_id", resp.RequestId)

	return nil
}
