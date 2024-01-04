package outscale

import (
	"context"
	"fmt"

	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func setOSCAPITags(conn *oscgo.APIClient, d *schema.ResourceData) error {

	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(*schema.Set)
		n := nraw.(*schema.Set)
		create, remove := diffOSCAPITags(tagsFromSliceMap(o), tagsFromSliceMap(n))

		// Set tag
		if len(remove) > 0 {
			err := resource.Retry(60*time.Second, func() *resource.RetryError {
				_, httpResp, err := conn.TagApi.DeleteTags(context.Background()).DeleteTagsRequest(oscgo.DeleteTagsRequest{
					ResourceIds: []string{d.Id()},
					Tags:        remove,
				}).Execute()
				if err != nil {
					if strings.Contains(fmt.Sprint(err), ".NotFound") {
						return resource.RetryableError(err) // retry
					}
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		if len(create) > 0 {
			err := resource.Retry(60*time.Second, func() *resource.RetryError {
				_, httpResp, err := conn.TagApi.CreateTags(context.Background()).CreateTagsRequest(oscgo.CreateTagsRequest{
					ResourceIds: []string{d.Id()},
					Tags:        create,
				}).Execute()
				if err != nil {
					if strings.Contains(fmt.Sprint(err), ".NotFound") {
						return resource.RetryableError(err) // retry
					}
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func tagsOAPIListSchemaComputed() *schema.Schema {
	return &schema.Schema{
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
			},
		},
	}
}

func tagsListOAPISchema2(computed bool) *schema.Schema {
	stype := schema.TypeSet

	if computed {
		stype = schema.TypeList
	}

	return &schema.Schema{
		Type: stype,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": mk_elem(computed, false, !computed,
					schema.TypeString),
				"value": mk_elem(computed, false, !computed,
					schema.TypeString),
			},
		},
		Optional: true,
		Computed: computed,
	}
}

func tagsListOAPISchema() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"value": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
			},
		},
		Optional: true,
	}
}

// tagsOSCsAPI	ToMap turns the list of tag into a map.
func tagsOSCAPIToMap(ts []oscgo.ResourceTag) []map[string]string {
	result := make([]map[string]string, len(ts))
	if len(ts) > 0 {
		for k, t := range ts {
			tag := make(map[string]string)
			tag["key"] = t.Key
			tag["value"] = t.Value
			result[k] = tag
		}
	} else {
		result = make([]map[string]string, 0)
	}

	return result
}

func tagsOSCAPIFromMap(m map[string]interface{}) []oscgo.ResourceTag {
	result := make([]oscgo.ResourceTag, 0, len(m))
	for k, v := range m {
		t := oscgo.ResourceTag{
			Key:   k,
			Value: v.(string),
		}
		result = append(result, t)
	}

	return result
}

// diffOSCAPITags takes our tag locally and the ones remotely and returns
// the set of tag that must be created, and the set of tag that must
// be destroyed.
func diffOSCAPITags(oldTags, newTags []oscgo.ResourceTag) ([]oscgo.ResourceTag, []oscgo.ResourceTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	stateTags := make(map[string]interface{})
	for _, t := range oldTags {
		stateTags[t.Key] = t.Value
	}

	tagsToCreate := make(map[string]interface{})
	for _, t := range newTags {
		old, ok := stateTags[t.Key]
		if !ok || old != t.Value {
			tagsToCreate[t.Key] = t.Value
		}
	}

	// Build the list of what to remove
	var remove []oscgo.ResourceTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			remove = append(remove, t)
		}
	}

	return tagsOSCAPIFromMap(tagsToCreate), remove
}

func tagsFromMapLBU(m map[string]interface{}) *[]oscgo.ResourceTag {
	result := make([]oscgo.ResourceTag, 0, len(m))
	for k, v := range m {
		t := oscgo.ResourceTag{
			Key:   k,
			Value: v.(string),
		}
		result = append(result, t)
	}

	return &result
}

func tagsFromSliceMap(m *schema.Set) []oscgo.ResourceTag {
	result := make([]oscgo.ResourceTag, 0, m.Len())
	for _, v := range m.List() {
		tag := v.(map[string]interface{})
		t := oscgo.ResourceTag{
			Key:   tag["key"].(string),
			Value: tag["value"].(string),
		}
		result = append(result, t)
	}

	return result
}

func oapiTagsDescToList(ts []oscgo.Tag) []map[string]interface{} {
	res := make([]map[string]interface{}, len(ts))

	for i, t := range ts {
		if !oapiTagDescIgnored(&t) {
			res[i] = map[string]interface{}{
				"key":           t.Key,
				"value":         t.Value,
				"resource_id":   t.ResourceId,
				"resource_type": t.ResourceType,
			}
		}
	}
	return res
}

func oapiTagDescIgnored(t *oscgo.Tag) bool {
	filter := []string{"^outscale:"}
	for _, v := range filter {
		if r, _ := regexp.MatchString(v, t.GetKey()); r {
			return true
		}
	}
	return false
}

func assignTags(tag *schema.Set, resourceID string, conn *oscgo.APIClient) error {
	request := oscgo.CreateTagsRequest{}
	request.Tags = tagsFromSliceMap(tag)
	request.ResourceIds = []string{resourceID}
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.TagApi.CreateTags(context.Background()).CreateTagsRequest(request).Execute()

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "NotFound") {
				return resource.RetryableError(err)
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

func dataSourceTagsSchema() *schema.Schema {
	return &schema.Schema{
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
			},
		},
	}
}

func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		ForceNew: true,
	}
}

func getOapiTagSet(tags *[]oscgo.ResourceTag) []map[string]interface{} {
	res := []map[string]interface{}{}

	if tags != nil {
		for _, t := range *tags {
			tag := map[string]interface{}{}

			tag["key"] = t.Key
			tag["value"] = t.Value

			res = append(res, tag)
		}
	}

	return res
}

func getOscAPITagSet(tags []oscgo.ResourceTag) []map[string]interface{} {
	res := []map[string]interface{}{}

	for _, t := range tags {
		tag := map[string]interface{}{}

		tag["key"] = t.Key
		tag["value"] = t.Value

		res = append(res, tag)
	}

	return res
}
