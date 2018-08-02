package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/aws/aws-sdk-go/aws"
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
	conn := meta.(*OutscaleClient).FCU

	request := &fcu.CreateTagsInput{}

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

	return resourceOutscaleOAPITagsRead(d, meta)
}

func resourceOutscaleOAPITagsRead(d *schema.ResourceData, meta interface{}) error {
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

	tg := tagsDescToList(resp.Tags)
	err = d.Set("tags", tg)

	return err
}

func resourceOutscaleOAPITagsDelete(d *schema.ResourceData, meta interface{}) error {
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

//TODO: OAPI
// func setOAPITags(conn *fcu.Client, d *schema.ResourceData) error {

// 	if d.HasChange("tag") {
// 		oraw, nraw := d.GetChange("tag")
// 		o := oraw.(map[string]interface{})
// 		n := nraw.(map[string]interface{})
// 		create, remove := diffTags(tagsFromMap(o), tagsFromMap(n))

// 		// Set tag
// 		if len(remove) > 0 {
// 			err := resource.Retry(60*time.Second, func() *resource.RetryError {
// 				log.Printf("[DEBUG] Removing tag: %#v from %s", remove, d.Id())
// 				_, err := conn.VM.DeleteTags(&fcu.DeleteTagsInput{
// 					Resources: []*string{aws.String(d.Id())},
// 					Tags:      remove,
// 				})
// 				if err != nil {
// 					ec2err, ok := err.(awserr.Error)
// 					if ok && strings.Contains(ec2err.Code(), ".NotFound") {
// 						return resource.RetryableError(err) // retry
// 					}
// 					return resource.NonRetryableError(err)
// 				}
// 				return nil
// 			})
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		if len(create) > 0 {
// 			err := resource.Retry(60*time.Second, func() *resource.RetryError {
// 				log.Printf("[DEBUG] Creating tag: %v for %s", create, d.Id())
// 				_, err := conn.VM.CreateTags(&fcu.CreateTagsInput{
// 					Resources: []*string{aws.String(d.Id())},
// 					Tags:      create,
// 				})
// 				if err != nil {
// 					ec2err, ok := err.(awserr.Error)
// 					if ok && strings.Contains(ec2err.Code(), ".NotFound") {
// 						return resource.RetryableError(err) // retry
// 					}
// 					return resource.NonRetryableError(err)
// 				}
// 				return nil
// 			})
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }
