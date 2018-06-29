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

func resourceOutscalePolicyGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscalePolicyGroupPut,
		Read:   resourceOutscalePolicyGroupRead,
		Delete: resourceOutscalePolicyGroupDelete,

		Schema: map[string]*schema.Schema{
			"policy_document": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_name": &schema.Schema{
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

func resourceOutscalePolicyGroupPut(d *schema.ResourceData, meta interface{}) error {
	eimconn := meta.(*OutscaleClient).EIM

	request := &eim.PutGroupPolicyInput{
		GroupName:      aws.String(d.Get("group_name").(string)),
		PolicyDocument: aws.String(d.Get("policy_document").(string)),
		PolicyName:     aws.String(d.Get("policy_name").(string)),
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = eimconn.API.PutGroupPolicy(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error putting IAM group policy %s: %s", *request.PolicyName, err)
	}

	d.SetId(fmt.Sprintf("%s:%s", *request.GroupName, *request.PolicyName))
	return nil
}

func resourceOutscalePolicyGroupRead(d *schema.ResourceData, meta interface{}) error {
	eimconn := meta.(*OutscaleClient).EIM

	group, name := resourceOutscalePolicyGroupParseID(d.Id())

	request := &eim.GetGroupPolicyInput{
		PolicyName: aws.String(name),
		GroupName:  aws.String(group),
	}

	var getResp *eim.GetGroupPolicyOutput
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		getResp, err = eimconn.API.GetGroupPolicy(request)

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
		return fmt.Errorf("Error reading IAM policy %s from group %s: %s", name, group, err)
	}

	if getResp.GetGroupPolicyResult.PolicyDocument == nil {
		return fmt.Errorf("GetGroupPolicy returned a nil policy document")
	}

	d.Set("request_id", getResp.ResponseMetadata.RequestID)

	return nil
}

func resourceOutscalePolicyGroupDelete(d *schema.ResourceData, meta interface{}) error {
	eimconn := meta.(*OutscaleClient).EIM

	group, name := resourceOutscalePolicyGroupParseID(d.Id())

	request := &eim.DeleteGroupPolicyInput{
		PolicyName: aws.String(name),
		GroupName:  aws.String(group),
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = eimconn.API.DeleteGroupPolicy(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting IAM group policy %s: %s", d.Id(), err)
	}
	return nil
}

func resourceOutscalePolicyGroupParseID(id string) (groupName, policyName string) {
	parts := strings.SplitN(id, ":", 2)
	groupName = parts[0]
	policyName = parts[1]
	return
}
