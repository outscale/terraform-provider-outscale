package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOutscalePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscalePolicyRead,

		Schema: map[string]*schema.Schema{
			"policy_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_document": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_attachable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscalePolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM

	policyArn, policyArnOk := d.GetOk("policy_arn")

	if policyArnOk == false {
		return fmt.Errorf("policy_arn must be provided")
	}

	input := &eim.GetPolicyInput{PolicyArn: aws.String(policyArn.(string))}

	var err error
	var output *eim.GetPolicyResult
	var rs *eim.GetPolicyOutput

	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rs, err = conn.API.GetPolicy(input)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if rs.GetPolicyResult != nil {
			output = rs.GetPolicyResult
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading policy with arn %s: %s", policyArn, err)
	}

	d.Set("arn", aws.StringValue(output.Policy.Arn))
	d.Set("attachment_count", aws.Int64Value(output.Policy.AttachmentCount))
	d.Set("default_version_id", aws.StringValue(output.Policy.DefaultVersionId))
	d.Set("description", aws.StringValue(output.Policy.Description))
	d.Set("is_attachable", aws.BoolValue(output.Policy.IsAttachable))
	d.Set("path", aws.StringValue(output.Policy.Path))
	d.Set("policy_id", aws.StringValue(output.Policy.PolicyId))
	d.Set("policy_name", aws.StringValue(output.Policy.PolicyName))
	d.Set("request_id", aws.StringValue(rs.ResponseMetadata.RequestID))
	d.SetId(*output.Policy.PolicyId)

	return nil
}
