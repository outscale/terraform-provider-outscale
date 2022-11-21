package outscale

import (
	"context"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceTags() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTagsRead,
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

func dataSourceTagsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).OSCAPI

	// Build up search parameters
	params := oscgo.ReadTagsRequest{}
	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		params.SetFilters(buildDataSourceFilters(filters.(*schema.Set)))
	}

	var resp oscgo.ReadTagsResponse
	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.TagApi.ReadTags(context.Background()).ReadTagsRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp.StatusCode, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	if err := d.Set("tags", tagsDescToList(resp.GetTags())); err != nil {
		return err
	}
	d.SetId(resource.UniqueId())

	return err
}
