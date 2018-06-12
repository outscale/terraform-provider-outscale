package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func resourceOutscaleGroupUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleGroupUserCreate,
		Read:   resourceOutscaleGroupUserRead,
		Delete: resourceOutscaleGroupUserDelete,
		Schema: map[string]*schema.Schema{
			"group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"groups": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arn": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceOutscaleGroupUserCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	group := d.Get("group_name").(string)
	user := d.Get("user_name").(string)

	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.API.AddUserToGroup(&eim.AddUserToGroupInput{
			UserName:  aws.String(user),
			GroupName: aws.String(group),
		})
		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(d.Get("group_name").(string))
	return resourceOutscaleGroupUserRead(d, meta)
}

func resourceOutscaleGroupUserRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	group := d.Get("group_name").(string)
	var ul []map[string]interface{}

	var err error
	var resp *eim.GetGroupOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.GetGroup(&eim.GetGroupInput{
			GroupName: aws.String(group),
		})
		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
			d.SetId("")
			return nil
		}
		return err
	}

	ul = append(ul, map[string]interface{}{
		"arn":        aws.StringValue(resp.Group.Arn),
		"group_id":   aws.StringValue(resp.Group.GroupId),
		"group_name": aws.StringValue(resp.Group.GroupName),
		"path":       aws.StringValue(resp.Group.Path),
	})

	return d.Set("groups", ul)
}

func resourceOutscaleGroupUserDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM

	user := d.Get("user_name").(string)
	group := d.Get("group_name").(string)

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.API.RemoveUserFromGroup(&eim.RemoveUserFromGroupInput{
			UserName:  aws.String(user),
			GroupName: aws.String(group),
		})
		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
			return nil
		}
		return err
	}

	return nil
}
