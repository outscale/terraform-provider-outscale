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

func resourceOutscalePolicyVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscalePolicyVersionUpdate,
		Read:   resourceOutscalePolicyVersionRead,
		Delete: resourceOutscalePolicyVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"set_as_default": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
			},
			"policy_document": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_arn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"document": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_id": {
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

func resourceOutscalePolicyVersionRead(d *schema.ResourceData, meta interface{}) error {
	eimconn := meta.(*OutscaleClient).EIM

	getPolicyVersionRequest := &eim.GetPolicyVersionInput{
		PolicyArn: aws.String(d.Get("policy_arn").(string)),
		VersionId: aws.String(d.Id()),
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

	return nil
}

func resourceOutscalePolicyVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	eimconn := meta.(*OutscaleClient).EIM

	request := &eim.CreatePolicyVersionInput{
		PolicyArn:      aws.String(d.Get("policy_arn").(string)),
		PolicyDocument: aws.String(d.Get("policy_document").(string)),
	}

	if v, ok := d.GetOk("set_as_default"); ok {
		request.SetAsDefault = aws.Bool(v.(bool))
	}

	var err error
	var resp *eim.CreatePolicyVersionOutput

	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = eimconn.API.CreatePolicyVersion(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error in policy version %s: %s", *resp.CreatePolicyVersionResult.PolicyVersion.VersionId, err)
	}

	d.SetId(*resp.CreatePolicyVersionResult.PolicyVersion.VersionId)

	return nil
}

func resourceOutscalePolicyVersionDelete(d *schema.ResourceData, meta interface{}) error {
	eimconn := meta.(*OutscaleClient).EIM

	return eimPolicyDeleteVersion(d.Get("policy_arn").(string), d.Id(), eimconn)
}
