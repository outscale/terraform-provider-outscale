package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleOAPIPrefixLists() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIPrefixListsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"prefix_list_id": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"prefix_list_set": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix_list_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefix_list_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_range": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceOutscaleOAPIPrefixListsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	prefix, prefixOk := d.GetOk("prefix_list_id")

	if !filtersOk && !prefixOk {
		return fmt.Errorf("One of prefix_list_id or filters must be assigned")
	}

	params := &fcu.DescribePrefixListsInput{}

	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}

	if prefixOk {
		var ids []*string
		for _, v := range prefix.([]interface{}) {
			ids = append(ids, aws.String(v.(string)))
		}
		params.PrefixListIds = ids
	}

	var resp *fcu.DescribePrefixListsOutput
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.VM.DescribePrefixLists(params)
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
	if resp == nil || len(resp.PrefixLists) == 0 {
		return fmt.Errorf("no matching prefix list found; the prefix list ID or name may be invalid or not exist in the current region")
	}

	d.SetId(resource.UniqueId())

	pls := make([]map[string]interface{}, len(resp.PrefixLists))

	for k, v := range resp.PrefixLists {
		pl := make(map[string]interface{})
		pl["prefix_list_id"] = *v.PrefixListId
		pl["prefix_list_name"] = *v.PrefixListName
		cidrs := make([]string, len(v.Cidrs))
		for i, v1 := range v.Cidrs {
			cidrs[i] = *v1
		}
		pl["ip_range"] = cidrs
		pls[k] = pl
	}

	d.Set("prefix_list_set", pls)
	d.Set("request_id", resp.RequestId)

	return nil
}
