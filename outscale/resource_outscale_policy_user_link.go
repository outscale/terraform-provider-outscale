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

func resourceOutscalePolicyUserLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscalePolicyUserLinkCreate,
		Read:   resourceOutscalePolicyUserLinkRead,
		Delete: resourceOutscalePolicyUserLinkDelete,

		Schema: map[string]*schema.Schema{
			"user_name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"policy_arn": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"policy_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			// "request_id": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
		},
	}
}

func resourceOutscalePolicyUserLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	user := d.Get("user_name").(string)
	arn := d.Get("policy_arn").(string)

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.API.AttachUserPolicy(&eim.AttachUserPolicyInput{
			UserName:  aws.String(user),
			PolicyArn: aws.String(arn),
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") || strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-", user)))

	return resourceOutscalePolicyUserLinkRead(d, meta)
}

func resourceOutscalePolicyUserLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	user := d.Get("user_name").(string)
	arn := d.Get("policy_arn").(string)

	var err error
	var attachedPolicies *eim.ListAttachedUserPoliciesOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		attachedPolicies, err = conn.API.ListAttachedUserPolicies(&eim.ListAttachedUserPoliciesInput{
			UserName: aws.String(user),
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") || strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	var po string
	var pn string
	for _, p := range attachedPolicies.AttachedPolicies {
		if *p.PolicyArn == arn {
			po = *p.PolicyArn
			pn = *p.PolicyName
		}
	}

	if po == "" {
		d.SetId("")
		return fmt.Errorf("No such User found for Policy Attachment: (%s)", err)
	}

	d.Set("policy_name", pn)

	return nil
}

func resourceOutscalePolicyUserLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	user := d.Get("user_name").(string)
	arn := d.Get("policy_arn").(string)

	req := &eim.DetachUserPolicyInput{
		UserName:  aws.String(user),
		PolicyArn: aws.String(arn),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.API.DetachUserPolicy(req)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") || strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
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
