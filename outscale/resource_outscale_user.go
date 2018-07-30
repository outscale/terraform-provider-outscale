package outscale

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleUserCreate,
		Read:   resourceOutscaleUserRead,
		Update: resourceOutscaleUserUpdate,
		Delete: resourceOutscaleUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOutscaleUserName,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOutscaleUserCreate(d *schema.ResourceData, meta interface{}) error {
	iamconn := meta.(*OutscaleClient).EIM
	name := d.Get("user_name").(string)

	request := &eim.CreateUserInput{
		Path:     aws.String("/"),
		UserName: aws.String(name),
	}

	if v, ok := d.GetOk("path"); ok {
		request.Path = aws.String(v.(string))
	}

	var err error
	var createResp *eim.CreateUserResult
	var resp *eim.CreateUserOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = iamconn.API.CreateUser(request)
		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		if resp.CreateUserResult != nil {
			createResp = resp.CreateUserResult
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating IAM User %s: %s", name, err)
	}

	d.SetId(*createResp.User.UserName)

	return resourceOutscaleUserRead(d, meta)
}

func resourceOutscaleUserRead(d *schema.ResourceData, meta interface{}) error {
	iamconn := meta.(*OutscaleClient).EIM

	request := &eim.GetUserInput{
		UserName: aws.String(d.Id()),
	}

	var err error
	var getResp *eim.GetUserResult
	var resp *eim.GetUserOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = iamconn.API.GetUser(request)
		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		if resp.GetUserResult != nil {
			getResp = resp.GetUserResult
		}

		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading IAM User %s: %s", d.Id(), err)
	}

	d.Set("user_name", aws.StringValue(getResp.User.UserName))
	d.Set("arn", aws.StringValue(getResp.User.Arn))
	d.Set("path", aws.StringValue(getResp.User.Path))
	d.Set("request_id", aws.StringValue(resp.ResponseMetadata.RequestID))

	return d.Set("user_id", aws.StringValue(getResp.User.UserId))
}

func resourceOutscaleUserUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("user_name") || d.HasChange("path") {
		iamconn := meta.(*OutscaleClient).EIM
		on, nn := d.GetChange("user_name")
		_, np := d.GetChange("path")

		request := &eim.UpdateUserInput{
			UserName:    aws.String(on.(string)),
			NewUserName: aws.String(nn.(string)),
			NewPath:     aws.String(np.(string)),
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = iamconn.API.UpdateUser(request)
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
			return fmt.Errorf("Error updating IAM User %s: %s", d.Id(), err)
		}
		return resourceOutscaleUserRead(d, meta)
	}
	return nil
}

func resourceOutscaleUserDelete(d *schema.ResourceData, meta interface{}) error {
	iamconn := meta.(*OutscaleClient).EIM

	request := &eim.DeleteUserInput{
		UserName: aws.String(d.Id()),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = iamconn.API.DeleteUser(request)
		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting IAM User %s: %s", d.Id(), err)
	}
	return nil
}

func validateOutscaleUserName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[0-9A-Za-z=,.@\-_+]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only alphanumeric characters, hyphens, underscores, commas, periods, @ symbols, plus and equals signs allowed in %q: %q",
			k, value))
	}
	return
}
