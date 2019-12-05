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

func dataSourceOutscaleOAPIPolicyVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIPolicyVersionsRead,

		Schema: map[string]*schema.Schema{
			"policy_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default_version": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPIPolicyVersionsRead(d *schema.ResourceData, meta interface{}) error {
	eimconn := meta.(*OutscaleClient).EIM

	policyArn, policyArnOk := d.GetOk("policy_arn")

	if !policyArnOk {
		return fmt.Errorf("policy_arn must be provided")
	}

	request := &eim.ListPolicyVersionsInput{
		PolicyArn: aws.String(policyArn.(string)),
	}

	var getPolicyVersionResponse *eim.ListPolicyVersionsOutput
	var err error

	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		getPolicyVersionResponse, err = eimconn.API.ListPolicyVersions(request)

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

	versions := make([]map[string]interface{}, len(getPolicyVersionResponse.ListPolicyVersionsResult.Versions))

	for k, v := range getPolicyVersionResponse.ListPolicyVersionsResult.Versions {
		version := make(map[string]interface{})

		if v.Document != nil {
			version["document"] = aws.StringValue(v.Document)
		}

		if v.VersionId != nil {
			version["version_id"] = aws.StringValue(v.VersionId)
		}

		if v.IsDefaultVersion != nil {
			version["is_default_version"] = aws.BoolValue(v.IsDefaultVersion)
		}

		versions[k] = version

	}

	d.Set("versions", versions)
	d.SetId(resource.UniqueId())
	return nil
}
