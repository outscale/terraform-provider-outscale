package outscale

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func resourceOutscaleGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleGroupCreate,
		Read:   resourceOutscaleGroupRead,
		Update: resourceOutscaleGroupUpdate,
		Delete: resourceOutscaleGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group": &schema.Schema{
				Type:     schema.TypeMap,
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
			"users": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arn": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": &schema.Schema{
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
			"group_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOutscaleGroupName,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/",
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceOutscaleGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	name := d.Get("group_name").(string)

	request := &eim.CreateGroupInput{
		GroupName: aws.String(name),
	}

	if path, ok := d.GetOk("path"); ok {
		request.Path = aws.String(path.(string))
	}

	var createResp *eim.CreateGroupOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		createResp, err = conn.API.CreateGroup(request)
		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating IAM Group %s: %s", name, err)
	}
	d.SetId(*createResp.Group.GroupName)
	return resourceOutscaleGroupRead(d, meta)
}

func resourceOutscaleGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	request := &eim.GetGroupInput{
		GroupName: aws.String(d.Id()),
	}

	var getResp *eim.GetGroupOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		getResp, err = conn.API.GetGroup(request)
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
		return fmt.Errorf("Error reading IAM Group %s: %s", d.Id(), err)
	}

	grp := make(map[string]interface{})
	grp["arn"] = aws.StringValue(getResp.Group.Arn)
	grp["group_id"] = aws.StringValue(getResp.Group.GroupId)
	grp["group_name"] = aws.StringValue(getResp.Group.GroupName)
	grp["path"] = aws.StringValue(getResp.Group.Path)

	usr := make([]map[string]interface{}, len(getResp.Users))
	for k, v := range getResp.Users {
		us := make(map[string]interface{})
		us["arn"] = aws.StringValue(v.Arn)
		us["user_id"] = aws.StringValue(v.UserId)
		us["user_name"] = aws.StringValue(v.UserName)
		us["path"] = aws.StringValue(v.Path)
		usr[k] = us
	}

	if err := d.Set("group", grp); err != nil {
		return err
	}

	return d.Set("users", usr)
}

func resourceOutscaleGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("group_name") || d.HasChange("path") {
		conn := meta.(*OutscaleClient).EIM
		on, nn := d.GetChange("group_name")
		_, np := d.GetChange("path")

		request := &eim.UpdateGroupInput{
			GroupName:    aws.String(on.(string)),
			NewGroupName: aws.String(nn.(string)),
			NewPath:      aws.String(np.(string)),
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = conn.API.UpdateGroup(request)
			if err != nil {
				if strings.Contains(err.Error(), "Throttling:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("Error updating IAM Group %s: %s", d.Id(), err)
		}
		return resourceOutscaleGroupRead(d, meta)
	}
	return nil
}

func resourceOutscaleGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	request := &eim.DeleteGroupInput{
		GroupName: aws.String(d.Id()),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.API.DeleteGroup(request)
		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting IAM Group %s: %s", d.Id(), err)
	}
	return nil
}

func validateOutscaleGroupName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[0-9A-Za-z=,.@\-_+]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only alphanumeric characters, hyphens, underscores, commas, periods, @ symbols, plus and equals signs allowed in %q: %q",
			k, value))
	}
	return
}
