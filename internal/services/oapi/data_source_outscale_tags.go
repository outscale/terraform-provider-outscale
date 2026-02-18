package oapi

import (
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleTags() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleTagsRead,
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

func DataSourceOutscaleTagsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.OutscaleClient).OSC

	// Build up search parameters
	params := osc.ReadTagsRequest{}
	filters, filtersOk := d.GetOk("filter")

	var err error
	if filtersOk {
		params.Filters, err = oapiBuildOutscaleDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp osc.ReadTagsResponse
	err = retry.Retry(60*time.Second, func() *retry.RetryError {
		rp, httpResp, err := client.TagApi.ReadTags(ctx).ReadTagsRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	if err := d.Set("tags", flattenOAPITagsDescSDK(resp.Tags)); err != nil {
		return err
	}
	d.SetId(id.UniqueId())

	return err
}
