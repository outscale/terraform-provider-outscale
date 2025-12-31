package oapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleTags() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleTagsCreate,
		Read:   ResourceOutscaleTagsRead,
		Delete: ResourceOutscaleTagsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: getOAPITagsSchema(),
	}
}

func ResourceOutscaleTagsCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	request := oscgo.CreateTagsRequest{}

	tag, tagsOk := d.GetOk("tag")
	resourceIds, resourceIdsOk := d.GetOk("resource_ids")
	if !tagsOk && !resourceIdsOk {
		return fmt.Errorf("One tag and resource id, must be assigned")
	}

	if tagsOk {
		request.SetTags(expandOAPITagsSDK(tag.(*schema.Set)))
	}
	if resourceIdsOk {
		var rids []string
		sgs := resourceIds.(*schema.Set).List()
		for _, v := range sgs {
			str := v.(string)
			rids = append(rids, str)
		}

		request.SetResourceIds(rids)
	}

	err := retry.Retry(60*time.Second, func() *retry.RetryError {
		_, httpResp, err := conn.TagApi.CreateTags(context.Background()).CreateTagsRequest(request).Execute()
		if err != nil {
			if httpResp.StatusCode == http.StatusNotFound {
				return retry.RetryableError(err)
			}
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(id.UniqueId())

	return ResourceOutscaleTagsRead(d, meta)
}

func ResourceOutscaleTagsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	// Build up search parameters
	params := oscgo.ReadTagsRequest{
		Filters: &oscgo.FiltersTag{},
	}

	tag, tagsOk := d.GetOk("tag")
	filter := oscgo.FiltersTag{}
	if tagsOk {
		tgs := expandOAPITagsSDK(tag.(*schema.Set))
		keys := make([]string, 0, len(tgs))
		values := make([]string, 0, len(tgs))
		for _, t := range tgs {
			keys = append(keys, t.Key)
			values = append(values, t.Value)
		}
		filter.SetKeys(keys)
		filter.SetValues(values)
		params.SetFilters(filter)

	}

	resourceIds, resourceIdsOk := d.GetOk("resource_ids")
	if resourceIdsOk {
		var rids []string
		sgs := resourceIds.(*schema.Set).List()
		for _, v := range sgs {
			str := v.(string)
			rids = append(rids, str)
		}

		filter.SetResourceIds(rids)
		params.SetFilters(filter)
	}

	var resp oscgo.ReadTagsResponse
	var err error

	err = retry.Retry(60*time.Second, func() *retry.RetryError {
		rp, httpResp, err := conn.TagApi.ReadTags(context.Background()).ReadTagsRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	tg := flattenOAPITagsDescSDK(resp.GetTags())
	if err := d.Set("tags", tg); err != nil {
		return err
	}

	return err
}

func ResourceOutscaleTagsDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	request := oscgo.DeleteTagsRequest{}

	tag, tagsOk := d.GetOk("tag")

	resourceIds, resourceIdsOk := d.GetOk("resource_ids")

	if !tagsOk && !resourceIdsOk {
		return fmt.Errorf("One tag and resource id, must be assigned")
	}

	if tagsOk {
		request.Tags = expandOAPITagsSDK(tag.(*schema.Set))
	}
	if resourceIdsOk {
		var rids []string
		sgs := resourceIds.(*schema.Set).List()
		for _, v := range sgs {
			str := v.(string)
			rids = append(rids, str)
		}

		request.SetResourceIds(rids)
	}

	err := retry.Retry(60*time.Second, func() *retry.RetryError {
		_, httpResp, err := conn.TagApi.DeleteTags(context.Background()).DeleteTagsRequest(request).Execute()
		if err != nil {
			if httpResp.StatusCode == http.StatusNotFound {
				return retry.RetryableError(err) // retry
			}
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func getOAPITagsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"resource_ids": {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"tag": {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
						ForceNew: true,
					},
					"value": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
					},
				},
			},
		},
		"tags": {
			Type:     schema.TypeSet,
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
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
