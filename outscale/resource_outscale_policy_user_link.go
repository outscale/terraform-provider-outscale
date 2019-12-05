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

func resourceOutscaleOAPIPolicyUserLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIPolicyUserLinkCreate,
		Read:   resourceOutscaleOAPIPolicyUserLinkRead,
		Delete: resourceOutscaleOAPIPolicyUserLinkDelete,

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
			"policy_id": &schema.Schema{
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

func resourceOutscaleOAPIPolicyUserLinkCreate(d *schema.ResourceData, meta interface{}) error {
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
			if strings.Contains(fmt.Sprint(err), "Throttling") {
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

	return resourceOutscaleOAPIPolicyUserLinkRead(d, meta)
}

func resourceOutscaleOAPIPolicyUserLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	user := d.Get("user_name").(string)

	var err error
	var resp *eim.GetUserPolicyOutput
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.API.GetUserPolicy(&eim.GetUserPolicyInput{
			UserName: aws.String(user),
		})

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
		return err
	}

	d.Set("policy_id", aws.StringValue(resp.PolicyName))

	return nil
}

func resourceOutscaleOAPIPolicyUserLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	user := d.Get("user_name").(string)
	arn := d.Get("policy_arn").(string)

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.API.DetachUserPolicy(&eim.DetachUserPolicyInput{
			UserName:  aws.String(user),
			PolicyArn: aws.String(arn),
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
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
