package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleOAPIPrefixList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIPrefixListRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"prefix_list_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"prefix_list_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_range": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceOutscaleOAPIPrefixListRead(d *schema.ResourceData, meta interface{}) error {
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
		params.PrefixListIds = aws.StringSlice([]string{prefix.(string)})
	}

	log.Printf("[DEBUG] DescribePrefixLists %s\n", params)

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

	if len(resp.PrefixLists) > 1 {
		return fmt.Errorf("multiple Prefix matched; use additional constraints to reduce matches to a single Prefix")
	}

	pl := resp.PrefixLists[0]

	d.SetId(*pl.PrefixListId)
	d.Set("prefix_list_id", pl.PrefixListId)
	d.Set("prefix_list_name", pl.PrefixListName)

	cidrs := make([]string, len(pl.Cidrs))
	for i, v := range pl.Cidrs {
		cidrs[i] = *v
	}
	d.Set("ip_range", cidrs)
	d.Set("request_id", resp.RequestId)

	return nil
}
