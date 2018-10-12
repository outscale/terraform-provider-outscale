package outscale

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/common"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
	"github.com/terraform-providers/terraform-provider-outscale/osc/icu"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleTags() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleTagsCreate,
		Read:   resourceOutscaleTagsRead,
		Delete: resourceOutscaleTagsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: getTagsSchema(),
	}
}

func resourceOutscaleTagsCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	request := &fcu.CreateTagsInput{}

	tag, tagsOk := d.GetOk("tag")

	resourceIds, resourceIdsOk := d.GetOk("resource_ids")

	if tagsOk == false && resourceIdsOk == false {
		return fmt.Errorf("One tag and resource id, must be assigned")
	}

	request.Tags = tagsFromMap(tag.(map[string]interface{}))

	var rids []*string
	sgs := resourceIds.(*schema.Set).List()
	for _, v := range sgs {
		str := v.(string)
		rids = append(rids, aws.String(str))
	}

	request.Resources = rids

	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		_, err := conn.VM.CreateTags(request)
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

	return resourceOutscaleTagsRead(d, meta)
}

func resourceOutscaleTagsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	// Build up search parameters
	params := &fcu.DescribeTagsInput{}
	filters := []*fcu.Filter{}

	tag, tagsOk := d.GetOk("tag")
	if tagsOk {
		tgs := tagsFromMap(tag.(map[string]interface{}))
		ts := make([]*string, 0, len(tgs))
		for _, t := range tgs {
			ts = append(ts, t.Key)
		}

		f := &fcu.Filter{
			Name:   aws.String("key"),
			Values: ts,
		}

		filters = append(filters, f)

	}

	resourceIds, resourceIdsOk := d.GetOk("resource_ids")
	if resourceIdsOk {
		var rids []*string
		sgs := resourceIds.(*schema.Set).List()
		for _, v := range sgs {
			str := v.(string)
			rids = append(rids, aws.String(str))
		}

		f := &fcu.Filter{
			Name:   aws.String("resource-id"),
			Values: rids,
		}

		filters = append(filters, f)
	}

	params.Filters = filters

	var resp *fcu.DescribeTagsOutput
	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		resp, err = conn.VM.DescribeTags(params)
		return resource.RetryableError(err)
	})

	if err != nil {
		return err
	}

	d.Set("request_id", resp.RequestId)
	tg := tagSetDescToList(resp.Tags)
	err = d.Set("tag_set", tg)

	return err
}

func resourceOutscaleTagsDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	request := &fcu.DeleteTagsInput{}

	tag, tagsOk := d.GetOk("tag")

	resourceIds, resourceIdsOk := d.GetOk("resource_ids")

	if tagsOk == false && resourceIdsOk == false {
		return fmt.Errorf("One tag and resource id, must be assigned")
	}

	if tagsOk {
		request.Tags = tagsFromMap(tag.(map[string]interface{}))
	}
	if resourceIdsOk {
		var rids []*string
		sgs := resourceIds.(*schema.Set).List()
		for _, v := range sgs {
			str := v.(string)
			rids = append(rids, aws.String(str))
		}

		request.Resources = rids
	}

	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		_, err := conn.VM.DeleteTags(request)
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

func getTagsSchema() map[string]*schema.Schema {
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
		"tag_set": {
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
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
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

func setTagsCommon(c *OutscaleClient, d *schema.ResourceData) error {
	conn := c.FCU

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

// diffOAPITags takes our tag locally and the ones remotely and returns
// the set of tag that must be created, and the set of tag that must
// be destroyed.
func diffOAPITags(oldTags, newTags []oapi.Tags_0) ([]oapi.Tags_0, []oapi.Tags_0) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []oapi.Tags_0
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			remove = append(remove, t)
		}
	}

	return tagsOAPIFromMap(create), remove
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

func tagsOAPIFromMap(m map[string]interface{}) []oapi.Tags_0 {
	result := make([]oapi.Tags_0, 0, len(m))
	for k, v := range m {
		t := oapi.Tags_0{
			Key:   k,
			Value: v.(string),
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

func tagsFromMapLBU(m map[string]interface{}) []*lbu.Tag {
	result := make([]*lbu.Tag, 0, len(m))
	for k, v := range m {
		t := &lbu.Tag{
			Key:   aws.String(k),
			Value: aws.String(v.(string)),
		}
		result = append(result, t)
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

// tagsOAPI	ToMap turns the list of tag into a map.
func tagsOAPIToMap(ts []oapi.Tags_0) []map[string]string {
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

func tagsToMapI(ts []*icu.Tag) []map[string]string {
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

	fmt.Printf("[DEBUG] TAG_SET %s", result)

	return result
}

func tagsToMapL(ts []*lbu.Tag) []map[string]string {
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

// tagIgnored compares a s against a list of strings and checks if it should
// be ignored or not
func tagIgnored(t *fcu.Tag) bool {
	filter := []string{"^outscale:"}
	for _, v := range filter {
		if r, _ := regexp.MatchString(v, *t.Key); r == true {
			return true
		}
	}
	return false
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

func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		ForceNew: true,
	}
}

func tagsSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeMap},
	}
}

func setOAPITags(conn *oapi.Client, d *schema.ResourceData) error {

	if d.HasChange("tag") {
		oraw, nraw := d.GetChange("tag")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := diffOAPITags(tagsOAPIFromMap(o), tagsOAPIFromMap(n))

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
