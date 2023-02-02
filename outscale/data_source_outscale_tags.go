package outscale

import (
	"context"
	"errors"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.ReadTagsRequest{}
	if filters, filtersOk := d.GetOk("filter"); filtersOk {
		req.SetFilters(oapiBuildOutscaleDataSourceFilters(filters.(*schema.Set)))
	}
	var resp oscgo.ReadTagsResponse
	var err error
	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.TagApi.ReadTags(context.Background()).ReadTagsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	if len(resp.GetTags()) < 1 {
		return errors.New("Your query returned no results. Please change your search criteria and try again")
	}

	if err := d.Set("tags", oapiTagsDescToList(resp.GetTags())); err != nil {
		return err
	}
	d.SetId(resource.UniqueId())

	return err
}
