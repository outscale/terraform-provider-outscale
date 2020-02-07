package outscale

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleOAPIQuotas() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIQuotasRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"quota_name": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"quota_type": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quota": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"short_description": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"firewall_rules_set_name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"max_value": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"account_id": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"used_value": &schema.Schema{
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

func dataSourceOutscaleOAPIQuotasRead(d *schema.ResourceData, meta interface{}) error {
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
			quota["short_description"] = *v.DisplayName
			quota["firewall_rules_set_name"] = *v.GroupName
			i, err := strconv.Atoi(*v.MaxQuotaValue)
			if err != nil {
				return err
			}
			quota["max_value"] = i
			quota["name"] = *v.Name
			quota["account_id"] = *v.OwnerId
			i2, err := strconv.Atoi(*v.UsedQuotaValue)
			if err != nil {
				return err
			}
			quota["used_value"] = i2
			quotas[k] = quota
		}

		q["quota"] = quotas

		qs[k] = q
	}

	if err := d.Set("quota_type", qs); err != nil {
		return err
	}

	d.Set("request_id", resp.RequestId)

	return nil
}
