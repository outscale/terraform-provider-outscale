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

func dataSourceOutscaleOAPIPolicyGroupLink() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIPolicyGroupLinkRead,

		Schema: map[string]*schema.Schema{
			"group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"attached_policies": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_arn": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"policy_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPIPolicyGroupLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	group := d.Get("group_name").(string)

	request := &eim.ListAttachedGroupPoliciesInput{
		GroupName:  aws.String(group),
		PathPrefix: aws.String("/"),
	}

	if v, ok := d.GetOk("path"); ok {
		request.PathPrefix = aws.String(v.(string))
	}

	var attachedPolicies *eim.ListAttachedGroupPoliciesOutput
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		attachedPolicies, err = conn.API.ListAttachedGroupPolicies(request)

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

	a := make([]map[string]interface{}, len(attachedPolicies.ListAttachedGroupPoliciesResult.AttachedPolicies))

	for k, v := range attachedPolicies.ListAttachedGroupPoliciesResult.AttachedPolicies {
		a[k] = map[string]interface{}{
			"policy_arn":  aws.StringValue(v.PolicyArn),
			"policy_name": aws.StringValue(v.PolicyName),
		}
	}

	d.Set("request_id", attachedPolicies.ResponseMetadata.RequestID)
	d.SetId(resource.UniqueId())

	return d.Set("attached_policies", attachedPolicies)
}
