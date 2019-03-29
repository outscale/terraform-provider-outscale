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

func resourceOutscaleUserPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleUserPolicyPut,
		Read:   resourceOutscaleUserPolicyRead,
		Delete: resourceOutscaleUserPolicyDelete,

		Schema: map[string]*schema.Schema{
			"policy_document": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleUserPolicyPut(d *schema.ResourceData, meta interface{}) error {
	iamconn := meta.(*OutscaleClient).EIM

	request := &eim.PutUserPolicyInput{
		UserName:       aws.String(d.Get("user_name").(string)),
		PolicyDocument: aws.String(d.Get("policy_document").(string)),
		PolicyName:     aws.String(d.Get("policy_name").(string)),
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = iamconn.API.PutUserPolicy(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error putting IAM user policy %s: %s", *request.PolicyName, err)
	}

	d.SetId(fmt.Sprintf("%s:%s", *request.UserName, *request.PolicyName))
	return nil
}

func resourceOutscaleUserPolicyRead(d *schema.ResourceData, meta interface{}) error {
	iamconn := meta.(*OutscaleClient).EIM

	user, name := resourceOutscaleUserPolicyParseID(d.Id())

	request := &eim.GetUserPolicyInput{
		PolicyName: aws.String(name),
		UserName:   aws.String(user),
	}

	var err error
	var getResp *eim.GetUserPolicyOutput
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		getResp, err = iamconn.API.GetUserPolicy(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
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
		return fmt.Errorf("Error reading IAM policy %s from user %s: %s", name, user, err)
	}

	if getResp.PolicyDocument == nil {
		return fmt.Errorf("GetUserPolicy returned a nil policy document")
	}

	return nil
}

func resourceOutscaleUserPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	iamconn := meta.(*OutscaleClient).EIM

	user, name := resourceOutscaleUserPolicyParseID(d.Id())

	request := &eim.DeleteUserPolicyInput{
		PolicyName: aws.String(name),
		UserName:   aws.String(user),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = iamconn.API.DeleteUserPolicy(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting IAM user policy %s: %s", d.Id(), err)
	}
	return nil
}

func resourceOutscaleUserPolicyParseID(id string) (userName, policyName string) {
	parts := strings.SplitN(id, ":", 2)
	userName = parts[0]
	policyName = parts[1]
	return
}
