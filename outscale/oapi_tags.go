package outscale

import (
	"context"
	"fmt"

	"regexp"
	"strings"
	"time"

	"github.com/antihax/optional"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"github.com/outscale/osc-go/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/osc/common"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func tagsSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeMap},
	}
}

func tagsOAPISchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Computed: true,
	}
}

func setOAPITags(conn *oapi.Client, d *schema.ResourceData) error {

	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.([]interface{})
		n := nraw.([]interface{})
		create, remove := diffOAPITags(tagsOAPIFromSliceMap(o), tagsOAPIFromSliceMap(n))

		// Set tag
		if len(remove) > 0 {
			err := resource.Retry(60*time.Second, func() *resource.RetryError {
				_, err := conn.POST_DeleteTags(oapi.DeleteTagsRequest{
					ResourceIds: []string{d.Id()},
					Tags:        remove,
				})
				if err != nil {
					if strings.Contains(fmt.Sprint(err), ".NotFound") {
						return resource.RetryableError(err) // retry
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		if len(create) > 0 {
			err := resource.Retry(60*time.Second, func() *resource.RetryError {
				_, err := conn.POST_CreateTags(oapi.CreateTagsRequest{
					ResourceIds: []string{d.Id()},
					Tags:        create,
				})
				if err != nil {
					if strings.Contains(fmt.Sprint(err), ".NotFound") {
						return resource.RetryableError(err) // retry
					}
					return resource.NonRetryableError(err)
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

func setOSCAPITags(conn *oscgo.APIClient, d *schema.ResourceData) error {

	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.([]interface{})
		n := nraw.([]interface{})
		create, remove := diffOSCAPITags(tagsFromSliceMap(o), tagsFromSliceMap(n))

		// Set tag
		if len(remove) > 0 {
			err := resource.Retry(60*time.Second, func() *resource.RetryError {
				_, _, err := conn.TagApi.DeleteTags(context.Background(), &oscgo.DeleteTagsOpts{DeleteTagsRequest: optional.NewInterface(oscgo.DeleteTagsRequest{
					ResourceIds: []string{d.Id()},
					Tags:        remove,
				})})
				if err != nil {
					if strings.Contains(fmt.Sprint(err), ".NotFound") {
						return resource.RetryableError(err) // retry
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		if len(create) > 0 {
			err := resource.Retry(60*time.Second, func() *resource.RetryError {
				_, _, err := conn.TagApi.CreateTags(context.Background(), &oscgo.CreateTagsOpts{CreateTagsRequest: optional.NewInterface(oscgo.CreateTagsRequest{
					ResourceIds: []string{d.Id()},
					Tags:        create,
				})})
				if err != nil {
					if strings.Contains(fmt.Sprint(err), ".NotFound") {
						return resource.RetryableError(err) // retry
					}
					return resource.NonRetryableError(err)
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

func tagsListOAPISchemaForceNew() *schema.Schema {
	tagsSchema := tagsListOAPISchema()
	tagsSchema.ForceNew = true
	return tagsSchema
}

func tagsOAPIListSchemaComputed() *schema.Schema {
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

func tagsOAPISchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: true,
		ForceNew: true,
	}
}

func tagsListOAPISchema() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeList,
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

// tagsOAPI	ToMap turns the list of tag into a map.
func tagsOAPIToMap(ts []oapi.ResourceTag) []map[string]string {
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

func tagsOAPIFromMap(m map[string]interface{}) []oapi.ResourceTag {
	result := make([]oapi.ResourceTag, 0, len(m))
	for k, v := range m {
		t := oapi.ResourceTag{
			Key:   k,
			Value: v.(string),
		}
		result = append(result, t)
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

// diffOAPITags takes our tag locally and the ones remotely and returns
// the set of tag that must be created, and the set of tag that must
// be destroyed.
func diffOAPITags(oldTags, newTags []oapi.ResourceTag) ([]oapi.ResourceTag, []oapi.ResourceTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []oapi.ResourceTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			remove = append(remove, t)
		}
	}

	return tagsOAPIFromMap(create), remove
}

func diffOSCAPITags(oldTags, newTags []oscgo.ResourceTag) ([]oscgo.ResourceTag, []oscgo.ResourceTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []oscgo.ResourceTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			remove = append(remove, t)
		}
	}

	return tagsOSCAPIFromMap(create), remove
}

func tagsOAPIFromSliceMap(m []interface{}) []oapi.ResourceTag {
	result := make([]oapi.ResourceTag, 0, len(m))
	for _, v := range m {
		tag := v.(map[string]interface{})
		t := oapi.ResourceTag{
			Key:   tag["key"].(string),
			Value: tag["value"].(string),
		}
		result = append(result, t)
	}

	return result
}

func tagsFromSliceMap(m []interface{}) []oscgo.ResourceTag {
	result := make([]oscgo.ResourceTag, 0, len(m))
	for _, v := range m {
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
		if r, _ := regexp.MatchString(v, t.GetKey()); r == true {
			return true
		}
	}
	return false
}

func assignOapiTags(tag []interface{}, resourceID string, conn *oapi.Client) error {
	request := oapi.CreateTagsRequest{}
	request.Tags = tagsOAPIFromSliceMap(tag)
	request.ResourceIds = []string{resourceID}
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
	return nil
}

func unAssignOapiTags(tag []interface{}, resourceID string, conn *oapi.Client) error {
	request := oapi.DeleteTagsRequest{}
	request.Tags = tagsOAPIFromSliceMap(tag)
	request.ResourceIds = []string{resourceID}
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		_, err := conn.POST_DeleteTags(request)
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
	return nil
}
func assignTags(tag []interface{}, resourceID string, conn *oscgo.APIClient) error {
	request := oscgo.CreateTagsRequest{}
	request.Tags = tagsFromSliceMap(tag)
	request.ResourceIds = []string{resourceID}
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		_, _, err := conn.TagApi.CreateTags(context.Background(), &oscgo.CreateTagsOpts{
			CreateTagsRequest: optional.NewInterface(request),
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "NotFound") {
				return resource.RetryableError(err)
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

//TODO: remove the following function after oapi integration

func tagsToMapC(ts []*common.Tag) []map[string]string {
	result := make([]map[string]string, len(ts))
	if len(ts) > 0 {
		for k, t := range ts {
			tag := make(map[string]string)
			tag["key"] = *t.Key
			tag["value"] = *t.Value
			result[k] = tag
		}
	} else {
		result = make([]map[string]string, 0)
	}

	return result
}

func tagsDescToMap(ts []*fcu.TagDescription) map[string]string {
	result := make(map[string]string)
	for _, t := range ts {
		if !tagDescIgnored(t) {
			result[*t.Key] = *t.Value
		}
	}

	return result
}

func tagsDescToList(ts []*fcu.TagDescription) []map[string]string {
	result := make([]map[string]string, len(ts))
	for k, t := range ts {
		if !tagDescIgnored(t) {
			r := map[string]string{}
			r["load_balancer_name"] = *t.Key
			r["value"] = *t.Value
			r["resource_id"] = *t.ResourceId
			r["resource_type"] = *t.ResourceType

			result[k] = r
		}
	}

	return result
}

func tagSetDescToList(ts []*fcu.TagDescription) []map[string]string {
	result := make([]map[string]string, len(ts))
	for k, t := range ts {
		if !tagDescIgnored(t) {
			r := map[string]string{}
			r["key"] = *t.Key
			r["value"] = *t.Value
			r["resource_id"] = *t.ResourceId
			r["resource_type"] = *t.ResourceType

			result[k] = r
		}
	}

	return result
}

// tagsToMap turns the list of tag into a map.
func tagsToMap(ts []*fcu.Tag) []map[string]string {
	result := make([]map[string]string, len(ts))
	if len(ts) > 0 {
		for k, t := range ts {
			tag := make(map[string]string)
			tag["key"] = *t.Key
			tag["value"] = *t.Value
			result[k] = tag
		}
	} else {
		result = make([]map[string]string, 0)
	}

	return result
}

func dataSourceTagsSchema() *schema.Schema {
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

func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		ForceNew: true,
	}
}

func tagDescIgnored(t *fcu.TagDescription) bool {
	filter := []string{"^outscale:"}
	for _, v := range filter {
		if r, _ := regexp.MatchString(v, *t.Key); r == true {
			return true
		}
	}
	return false
}

func setTags(conn *fcu.Client, d *schema.ResourceData) error {

	if d.HasChange("tag") {
		oraw, nraw := d.GetChange("tag")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := diffTags(tagsFromMap(o), tagsFromMap(n))

		// Set tag
		if len(remove) > 0 {
			err := resource.Retry(60*time.Second, func() *resource.RetryError {
				_, err := conn.VM.DeleteTags(&fcu.DeleteTagsInput{
					Resources: []*string{aws.String(d.Id())},
					Tags:      remove,
				})
				if err != nil {
					if strings.Contains(fmt.Sprint(err), ".NotFound") {
						return resource.RetryableError(err) // retry
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		if len(create) > 0 {
			err := resource.Retry(60*time.Second, func() *resource.RetryError {
				_, err := conn.VM.CreateTags(&fcu.CreateTagsInput{
					Resources: []*string{aws.String(d.Id())},
					Tags:      create,
				})
				if err != nil {
					if strings.Contains(fmt.Sprint(err), ".NotFound") {
						return resource.RetryableError(err) // retry
					}
					return resource.NonRetryableError(err)
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

// diffTags takes our tag locally and the ones remotely and returns
// the set of tag that must be created, and the set of tag that must
// be destroyed.
func diffTags(oldTags, newTags []*fcu.Tag) ([]*fcu.Tag, []*fcu.Tag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[*t.Key] = *t.Value
	}

	// Build the list of what to remove
	var remove []*fcu.Tag
	for _, t := range oldTags {
		old, ok := create[*t.Key]
		if !ok || old != *t.Value {
			remove = append(remove, t)
		}
	}

	return tagsFromMap(create), remove
}

func tagsFromMap(m map[string]interface{}) []*fcu.Tag {
	result := make([]*fcu.Tag, 0, len(m))
	for k, v := range m {
		t := &fcu.Tag{
			Key:   aws.String(k),
			Value: aws.String(v.(string)),
		}
		result = append(result, t)
	}

	return result
}

func diffTagsCommon(oldTags, newTags []*common.Tag) ([]*common.Tag, []*common.Tag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[*t.Key] = *t.Value
	}

	// Build the list of what to remove
	var remove []*common.Tag
	for _, t := range oldTags {
		old, ok := create[*t.Key]
		if !ok || old != *t.Value {
			remove = append(remove, t)
		}
	}

	return tagsFromMapCommon(create), remove
}

// tagsFromMap returns the tag for the given map of data.

func tagsFromMapCommon(m map[string]interface{}) []*common.Tag {
	result := make([]*common.Tag, 0, len(m))
	for k, v := range m {
		t := &common.Tag{
			Key:   aws.String(k),
			Value: aws.String(v.(string)),
		}
		result = append(result, t)
	}

	return result
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

	if tags != nil {
		for _, t := range tags {
			tag := map[string]interface{}{}

			tag["key"] = t.Key
			tag["value"] = t.Value

			res = append(res, tag)
		}
	}

	return res
}
