package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/outscale/osc-go/oapi"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleOAPITags() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPITagsCreate,
		Read:   resourceOutscaleOAPITagsRead,
		Delete: resourceOutscaleOAPITagsDelete,
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

func resourceOutscaleOAPITagsCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	request := oapi.CreateTagsRequest{}

	tag, tagsOk := d.GetOk("tag")
	resourceIds, resourceIdsOk := d.GetOk("resource_ids")
	if !tagsOk && !resourceIdsOk {
		return fmt.Errorf("One tag and resource id, must be assigned")
	}

	request.Tags = tagsOAPIFromMap(tag.(map[string]interface{}))

	var rids []string
	sgs := resourceIds.(*schema.Set).List()
	for _, v := range sgs {
		str := v.(string)
		rids = append(rids, str)
	}

	request.ResourceIds = rids

	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		_, err := conn.POST_CreateTags(request)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), ".NotFound") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(resource.UniqueId())

	return resourceOutscaleOAPITagsRead(d, meta)
}

func resourceOutscaleOAPITagsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	// Build up search parameters
	params := oapi.ReadTagsRequest{
		Filters: oapi.FiltersTag{},
	}

	tag, tagsOk := d.GetOk("tag")
	if tagsOk {
		tgs := tagsOAPIFromMap(tag.(map[string]interface{}))
		keys := make([]string, 0, len(tgs))
		values := make([]string, 0, len(tgs))
		for _, t := range tgs {
			keys = append(keys, t.Key)
			values = append(values, t.Value)
		}

		params.Filters.Keys = keys
		params.Filters.Values = values

	}

	resourceIds, resourceIdsOk := d.GetOk("resource_ids")
	if resourceIdsOk {
		var rids []string
		sgs := resourceIds.(*schema.Set).List()
		for _, v := range sgs {
			str := v.(string)
			rids = append(rids, str)
		}

		params.Filters.ResourceIds = rids
	}

	var resp *oapi.POST_ReadTagsResponses
	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		resp, err = conn.POST_ReadTags(params)
		return resource.RetryableError(err)
	})

	if err != nil {
		return err
	}

	tg := oapiTagsDescToList(resp.OK.Tags)
	err = d.Set("tags", tg)

	if err := d.Set("request_id", resp.OK.ResponseContext.RequestId); err != nil {
		return err
	}

	return err
}

func resourceOutscaleOAPITagsDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	request := oapi.DeleteTagsRequest{}

	tag, tagsOk := d.GetOk("tag")

	resourceIds, resourceIdsOk := d.GetOk("resource_ids")

	if tagsOk == false && resourceIdsOk == false {
		return fmt.Errorf("One tag and resource id, must be assigned")
	}

	if tagsOk {
		request.Tags = tagsOAPIFromMap(tag.(map[string]interface{}))
	}
	if resourceIdsOk {
		var rids []string
		sgs := resourceIds.(*schema.Set).List()
		for _, v := range sgs {
			str := v.(string)
			rids = append(rids, str)
		}

		request.ResourceIds = rids
	}

	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		_, err := conn.POST_DeleteTags(request)
		if err != nil {
			ec2err, ok := err.(awserr.Error)
			if ok && strings.Contains(ec2err.Code(), ".NotFound") {
				return resource.RetryableError(err) // retry
			}
			return resource.NonRetryableError(err)
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
			Type:     schema.TypeMap,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
					},
					"value": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
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
	}
}
