package outscale

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func dataSourceOutscaleOAPITags() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPITagsRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},
		},
	}
}

func dataSourceOutscaleOAPITagsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	// Build up search parameters
	params := oapi.ReadTagsRequest{}
	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		params.Filters = oapiBuildOutscaleDataSourceFilters(filters.(*schema.Set))
	}

	var resp *oapi.POST_ReadTagsResponses
	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		resp, err = conn.POST_ReadTags(params)
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

	d.Set("tags", oapiTagsDescToList(resp.OK.Tags))
	d.SetId(resource.UniqueId())

	return err
}
