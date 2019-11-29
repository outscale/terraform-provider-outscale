package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func resourceOutscaleOAPIPolicyGroupLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIPolicyGroupLinkCreate,
		Read:   resourceOutscaleOAPIPolicyGroupLinkRead,
		Delete: resourceOutscaleOAPIPolicyGroupLinkDelete,

		Schema: map[string]*schema.Schema{
			"group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_id": &schema.Schema{
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

func resourceOutscaleOAPIPolicyGroupLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM

	group := d.Get("group_name").(string)
	arn := d.Get("policy_id").(string)

	err := attachPolicyToGroup(conn, group, arn)
	if err != nil {
		return fmt.Errorf("[WARN] Error attaching policy %s to EIM group %s: %v", arn, group, err)
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-", group)))
	return resourceOutscaleOAPIPolicyGroupLinkRead(d, meta)
}

func resourceOutscaleOAPIPolicyGroupLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	group := d.Get("group_name").(string)
	arn := d.Get("policy_id").(string)

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.API.GetGroup(&eim.GetGroupInput{
			GroupName: aws.String(group),
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
			log.Printf("[WARN] No such entity found for Policy Attachment (%s)", group)
			d.SetId("")
			return nil
		}
		return err
	}

	var attachedPolicies *eim.ListAttachedGroupPoliciesOutput
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		attachedPolicies, err = conn.API.ListAttachedGroupPolicies(&eim.ListAttachedGroupPoliciesInput{
			GroupName: aws.String(group),
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

	var policy string
	for _, p := range attachedPolicies.ListAttachedGroupPoliciesResult.AttachedPolicies {
		if *p.PolicyArn == arn {
			policy = *p.PolicyArn
		}
	}

	if policy == "" {
		log.Printf("[WARN] No such policy found for Group Policy Attachment (%s)", group)
		d.SetId("")
	}

	d.Set("request_id", attachedPolicies.ResponseMetadata.RequestID)

	return nil
}

func resourceOutscaleOAPIPolicyGroupLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	group := d.Get("group_name").(string)
	arn := d.Get("policy_id").(string)

	if err := detachPolicyFromGroup(conn, group, arn); err != nil {
		return fmt.Errorf("[WARN] Error removing policy %s from IAM Group %s: %v", arn, group, err)
	}
	return nil
}

func attachPolicyToGroup(conn *eim.Client, group string, arn string) error {

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.API.AttachGroupPolicy(&eim.AttachGroupPolicyInput{
			GroupName: aws.String(group),
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

func detachPolicyFromGroup(conn *eim.Client, group string, arn string) error {

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.API.DetachGroupPolicy(&eim.DetachGroupPolicyInput{
			GroupName: aws.String(group),
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
