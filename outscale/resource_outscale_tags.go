package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"strings"
	"time"

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
	conn := meta.(*OutscaleClient).OSCAPI

	request := oscgo.CreateTagsRequest{}

	tag, tagsOk := d.GetOk("tag")
	resourceIds, resourceIdsOk := d.GetOk("resource_ids")
	if !tagsOk && !resourceIdsOk {
		return fmt.Errorf("One tag and resource id, must be assigned")
	}

	if tagsOk {
		request.SetTags(tagsFromSliceMap(tag.([]interface{})))
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

	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		_, _, err := conn.TagApi.CreateTags(context.Background(), &oscgo.CreateTagsOpts{CreateTagsRequest: optional.NewInterface(request)})
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
	conn := meta.(*OutscaleClient).OSCAPI

	// Build up search parameters
	params := oscgo.ReadTagsRequest{
		Filters: &oscgo.FiltersTag{},
	}

	tag, tagsOk := d.GetOk("tag")
	filter := oscgo.FiltersTag{}
	if tagsOk {
		tgs := tagsFromSliceMap(tag.([]interface{}))
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

	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		resp, _, err = conn.TagApi.ReadTags(context.Background(), &oscgo.ReadTagsOpts{ReadTagsRequest: optional.NewInterface(params)})
		return resource.RetryableError(err)
	})

	if err != nil {
		return err
	}

	d.Set("request_id", resp.GetResponseContext().RequestId)
	tg := oapiTagsDescToList(resp.GetTags())
	err = d.Set("tags", tg)

	if err := d.Set("request_id", resp.ResponseContext.GetRequestId()); err != nil {
		return err
	}

	return err
}

func resourceOutscaleOAPITagsDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	request := oscgo.DeleteTagsRequest{}

	tag, tagsOk := d.GetOk("tag")

	resourceIds, resourceIdsOk := d.GetOk("resource_ids")

	if tagsOk == false && resourceIdsOk == false {
		return fmt.Errorf("One tag and resource id, must be assigned")
	}

	if tagsOk {
		request.Tags = tagsFromSliceMap(tag.([]interface{}))
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

	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		_, _, err := conn.TagApi.DeleteTags(context.Background(), &oscgo.DeleteTagsOpts{DeleteTagsRequest: optional.NewInterface(request)})
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
			Type:     schema.TypeList,
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
