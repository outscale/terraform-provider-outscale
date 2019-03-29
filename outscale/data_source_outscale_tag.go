package outscale

import (
	"errors"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOutscaleTag() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleTagRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleTagRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	// Build up search parameters
	params := &fcu.DescribeTagsInput{}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}

	var resp *fcu.DescribeTagsOutput
	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		resp, err = conn.VM.DescribeTags(params)
		return resource.RetryableError(err)
	})

	if err != nil {
		return err
	}

	if len(resp.Tags) > 1 {
		return errors.New("Your query returned more than one result. Please try a more " +
			"specific search criteria")
	}

	tag := resp.Tags[0]

	d.Set("key", *tag.Key)
	d.Set("value", *tag.Value)
	d.Set("resource_id", *tag.ResourceId)
	d.Set("resource_type", *tag.ResourceType)
	d.Set("request_id", resp.RequestId)

	d.SetId(resource.UniqueId())

	return err
}
