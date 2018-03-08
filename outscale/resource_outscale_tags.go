package outscale

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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

	tags, tagsOk := d.GetOk("tags")

	resourceIds, resourceIdsOk := d.GetOk("resource_ids")

	if tagsOk == false && resourceIdsOk == false {
		return fmt.Errorf("One tag and resource id, must be assigned")
	}

	if tagsOk {
		request.Tags = tagsFromMap(tags.(map[string]interface{}))
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

	tags, tagsOk := d.GetOk("tags")
	if tagsOk {
		tgs := tagsFromMap(tags.(map[string]interface{}))
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

	tg := tagsDescToList(resp.Tags)
	err = d.Set("tag_set", tg)

	return err
}

func resourceOutscaleTagsDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	request := &fcu.DeleteTagsInput{}

	tags, tagsOk := d.GetOk("tags")

	resourceIds, resourceIdsOk := d.GetOk("resource_ids")

	if tagsOk == false && resourceIdsOk == false {
		return fmt.Errorf("One tag and resource id, must be assigned")
	}

	if tagsOk {
		request.Tags = tagsFromMap(tags.(map[string]interface{}))
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

func getTagsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"resource_ids": {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"tags": {
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
				log.Printf("[DEBUG] Removing tags: %#v from %s", remove, d.Id())
				_, err := conn.VM.DeleteTags(&fcu.DeleteTagsInput{
					Resources: []*string{aws.String(d.Id())},
					Tags:      remove,
				})
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
		}
		if len(create) > 0 {
			err := resource.Retry(60*time.Second, func() *resource.RetryError {
				log.Printf("[DEBUG] Creating tags: %s for %s", create, d.Id())
				_, err := conn.VM.CreateTags(&fcu.CreateTagsInput{
					Resources: []*string{aws.String(d.Id())},
					Tags:      create,
				})
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

// tagsFromMap returns the tag for the given map of data.
func tagsFromMap(m map[string]interface{}) []*fcu.Tag {
	result := make([]*fcu.Tag, 0, len(m))
	for k, v := range m {
		t := &fcu.Tag{
			Key:   aws.String(k),
			Value: aws.String(v.(string)),
		}
		if !tagIgnored(t) {
			result = append(result, t)
		}
	}

	return result
}

// tagsToMap turns the list of tags into a map.
func tagsToMap(ts []*fcu.Tag) []map[string]string {
	result := make([]map[string]string, len(ts))
	for k, t := range ts {
		r := map[string]string{}
		r["key"] = *t.Value
		r["value"] = *t.Value

		result[k] = r
	}

	fmt.Printf("[DEBUG] TAG_SET %s", result)

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
	result := []map[string]string{}
	for _, t := range ts {
		if !tagDescIgnored(t) {
			r := map[string]string{}
			r["key"] = *t.Value
			r["value"] = *t.Value
			r["resource_id"] = *t.ResourceId
			r["resource_type"] = *t.ResourceType

			result = append(result, r)
		}
	}

	return result
}

// tagIgnored compares a s against a list of strings and checks if it should
// be ignored or not
func tagIgnored(t *fcu.Tag) bool {
	filter := []string{"^outscale:"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching %v with %v\n", v, *t.Key)
		if r, _ := regexp.MatchString(v, *t.Key); r == true {
			log.Printf("[DEBUG] Found Outscale specific s %s (val: %s), ignoring.\n", *t.Key, *t.Value)
			return true
		}
	}
	return false
}

func tagDescIgnored(t *fcu.TagDescription) bool {
	filter := []string{"^outscale:"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching %v with %v\n", v, *t.Key)
		if r, _ := regexp.MatchString(v, *t.Key); r == true {
			log.Printf("[DEBUG] Found AWS specific s %s (val: %s), ignoring.\n", *t.Key, *t.Value)
			return true
		}
	}
	return false
}

func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}
}
