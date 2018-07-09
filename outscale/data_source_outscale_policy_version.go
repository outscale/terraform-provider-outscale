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

func dataSourceOutscalePolicyVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscalePolicyVersionRead,

		Schema: map[string]*schema.Schema{
			"policy_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"document": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default_version": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscalePolicyVersionRead(d *schema.ResourceData, meta interface{}) error {
	eimconn := meta.(*OutscaleClient).EIM

	policyArn, policyArnOk := d.GetOk("policy_arn")
	versionID, versionIDOk := d.GetOk("version_id")

	if !policyArnOk || !versionIDOk {
		return fmt.Errorf("Both policy_arn and version_id must be provided")
	}

	getPolicyVersionRequest := &eim.GetPolicyVersionInput{
		PolicyArn: aws.String(policyArn.(string)),
		VersionId: aws.String(versionID.(string)),
	}

	var getPolicyVersionResponse *eim.GetPolicyVersionOutput
	var err error

	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		getPolicyVersionResponse, err = eimconn.API.GetPolicyVersion(getPolicyVersionRequest)

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
		return fmt.Errorf("Error reading IAM policy version %s: %s", d.Get("policy_arn").(string), err)
	}

	d.Set("document", aws.StringValue(getPolicyVersionResponse.GetPolicyVersionResult.PolicyVersion.Document))
	d.Set("version_id", aws.StringValue(getPolicyVersionResponse.GetPolicyVersionResult.PolicyVersion.VersionId))
	d.Set("request_id", aws.StringValue(getPolicyVersionResponse.ResponseMetadata.RequestID))
	d.Set("is_default_version", aws.BoolValue(getPolicyVersionResponse.GetPolicyVersionResult.PolicyVersion.IsDefaultVersion))

	d.SetId(aws.StringValue(getPolicyVersionResponse.GetPolicyVersionResult.PolicyVersion.VersionId))

	return nil
}
