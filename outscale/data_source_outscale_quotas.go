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

func dataSourceOutscaleQuotas() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleQuotasRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"quota_name": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"reference_quota_set": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"reference": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleQuotasRead(d *schema.ResourceData, meta interface{}) error {
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
		var ids []*string
		for _, v := range quota.([]interface{}) {
			ids = append(ids, aws.String(v.(string)))
		}
		params.QuotaName = ids
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

	d.SetId(resource.UniqueId())

	qs := make([]map[string]interface{}, len(resp.ReferenceQuotaSet))

	for k, v := range resp.ReferenceQuotaSet {
		q := make(map[string]interface{})
		q["reference"] = *v.Reference

		quotas := make([]map[string]interface{}, len(v.QuotaSet))
		for k, v := range v.QuotaSet {
			quota := make(map[string]interface{})
			quota["description"] = *v.Description
			quota["display_name"] = *v.DisplayName
			quota["group_name"] = *v.GroupName
			i, err := strconv.Atoi(*v.MaxQuotaValue)
			if err != nil {
				return err
			}
			quota["max_quota_value"] = i
			quota["name"] = *v.Name
			quota["owner_id"] = *v.OwnerId
			i2, err := strconv.Atoi(*v.UsedQuotaValue)
			if err != nil {
				return err
			}
			quota["used_quota_value"] = i2
			quotas[k] = quota
		}

		q["quota_set"] = quotas

		qs[k] = q
	}

	if err := d.Set("reference_quota_set", qs); err != nil {
		return err
	}

	d.Set("request_id", resp.RequestId)

	return nil
}
