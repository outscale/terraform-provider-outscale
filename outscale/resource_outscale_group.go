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

func resourceOutscaleOAPIGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIGroupCreate,
		Read:   resourceOutscaleOAPIGroupRead,
		Update: resourceOutscaleOAPIGroupUpdate,
		Delete: resourceOutscaleOAPIGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOutscaleOAPIGroupName,
			},
			"group_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/",
			},
		},
	}
}
func resourceOutscaleOAPIGroupCreate(d *schema.ResourceData, meta interface{}) error {
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
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating IAM Group %s: %s", name, err)
	}
	d.SetId(*createResp.CreateGroupResult.Group.GroupName)
	return resourceOutscaleOAPIGroupRead(d, meta)
}

func resourceOutscaleOAPIGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	request := &eim.GetGroupInput{
		GroupName: aws.String(d.Id()),
	}

	var getResp *eim.GetGroupOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		getResp, err = conn.API.GetGroup(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
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

	d.Set("group_id", aws.StringValue(getResp.GetGroupResult.Group.GroupId))
	d.Set("group_name", aws.StringValue(getResp.GetGroupResult.Group.GroupName))
	d.Set("path", aws.StringValue(getResp.GetGroupResult.Group.Path))

	return nil
}

func resourceOutscaleOAPIGroupUpdate(d *schema.ResourceData, meta interface{}) error {

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
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("Error updating IAM Group %s: %s", d.Id(), err)
		}
		return resourceOutscaleOAPIGroupRead(d, meta)
	}
	return nil
}

func resourceOutscaleOAPIGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	request := &eim.DeleteGroupInput{
		GroupName: aws.String(d.Id()),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.API.DeleteGroup(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
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

func validateOutscaleOAPIGroupName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[0-9A-Za-z=,.@\-_+]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only alphanumeric characters, hyphens, underscores, commas, periods, @ symbols, plus and equals signs allowed in %q: %q",
			k, value))
	}
	return
}
