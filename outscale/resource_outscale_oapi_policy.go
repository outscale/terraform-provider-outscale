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

func resourceOutscaleOAPIPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIPolicyCreate,
		Read:   resourceOutscaleOAPIPolicyRead,
		Delete: resourceOutscaleOAPIPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"document": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) > 128 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 128 characters", k))
					}
					if !regexp.MustCompile("^[\\w+=,.@-]*$").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must match [\\w+=,.@-]", k))
					}
					return
				},
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resources_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policy_default_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_linkable": {
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

func resourceOutscaleOAPIPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	var name string

	request := &eim.CreatePolicyInput{
		Path: aws.String(d.Get("path").(string)),
	}

	if v, ok := d.GetOk("policy_name"); ok {
		request.PolicyName = aws.String(v.(string))
	} else {
		request.PolicyName = aws.String(resource.UniqueId())
	}
	if v, ok := d.GetOk("description"); ok {
		request.Description = aws.String(v.(string))
	}
	if v, ok := d.GetOk("document"); ok {
		request.PolicyDocument = aws.String(v.(string))
	}

	var err error
	var response *eim.CreatePolicyResult
	var rs *eim.CreatePolicyOutput
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rs, err = conn.API.CreatePolicy(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if rs.CreatePolicyResult != nil {
			response = rs.CreatePolicyResult
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating IAM policy %s: %s", name, err)
	}

	d.SetId(*response.Policy.Arn)

	return resourceOutscaleOAPIPolicyRead(d, meta)
}

func resourceOutscaleOAPIPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM

	getPolicyRequest := &eim.GetPolicyInput{
		PolicyArn: aws.String(d.Id()),
	}

	var err error
	var getPolicyResponse *eim.GetPolicyResult
	var rs *eim.GetPolicyOutput

	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rs, err = conn.API.GetPolicy(getPolicyRequest)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if rs.GetPolicyResult != nil {
			getPolicyResponse = rs.GetPolicyResult
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading IAM policy %s: %s", d.Id(), err)
	}

	d.Set("arn", aws.StringValue(getPolicyResponse.Policy.Arn))
	d.Set("resources_count", aws.Int64Value(getPolicyResponse.Policy.AttachmentCount))
	d.Set("policy_default_version_id", aws.StringValue(getPolicyResponse.Policy.DefaultVersionId))
	d.Set("description", aws.StringValue(getPolicyResponse.Policy.Description))
	d.Set("is_linkable", aws.BoolValue(getPolicyResponse.Policy.IsAttachable))
	d.Set("path", aws.StringValue(getPolicyResponse.Policy.Path))
	d.Set("policy_id", aws.StringValue(getPolicyResponse.Policy.PolicyId))
	d.Set("policy_name", aws.StringValue(getPolicyResponse.Policy.PolicyName))
	d.Set("request_id", aws.StringValue(rs.ResponseMetadata.RequestID))

	return nil
}

func resourceOutscaleOAPIPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	if err := eimPolicyDeleteNondefaultVersions(d.Id(), conn); err != nil {
		return err
	}
	request := &eim.DeletePolicyInput{
		PolicyArn: aws.String(d.Id()),
	}

	var err error

	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.API.DeletePolicy(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
			return nil
		}
		return fmt.Errorf("Error deleting IAM policy %s: %#v", d.Id(), err)
	}

	return nil
}

func eimPolicyDeleteNondefaultVersions(arn string, EIM *eim.Client) error {
	versions, err := eimPolicyListVersions(arn, EIM)
	if err != nil {
		return err
	}
	for _, version := range versions {
		if *version.IsDefaultVersion {
			continue
		}
		if err := eimPolicyDeleteVersion(arn, *version.VersionId, EIM); err != nil {
			return err
		}
	}
	return nil
}

func eimPolicyDeleteVersion(arn, versionID string, EIM *eim.Client) error {
	request := &eim.DeletePolicyVersionInput{
		PolicyArn: aws.String(arn),
		VersionId: aws.String(versionID),
	}

	var err error

	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = EIM.API.DeletePolicyVersion(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting version %s from IAM policy %s: %s", versionID, arn, err)
	}
	return nil
}

func eimPolicyListVersions(arn string, conn *eim.Client) ([]*eim.PolicyVersion, error) {
	request := &eim.ListPolicyVersionsInput{
		PolicyArn: aws.String(arn),
	}

	var err error
	var response *eim.ListPolicyVersionsOutput
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		response, err = conn.API.ListPolicyVersions(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error listing versions for IAM policy %s: %s", arn, err)
	}
	return response.ListPolicyVersionsResult.Versions, nil
}
