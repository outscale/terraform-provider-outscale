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

func dataSourceOutscalePolicyUserLink() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscalePolicyUserLinkRead,

		Schema: map[string]*schema.Schema{
			"user_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"path_prefix": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"attached_policies": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_arn": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			// "request_id": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
		},
	}
}

func dataSourceOutscalePolicyUserLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	user := d.Get("user_name").(string)

	request := &eim.ListAttachedUserPoliciesInput{
		UserName: aws.String(user),
	}

	if v, ok := d.GetOk("path_prefix"); ok {
		request.PathPrefix = aws.String(v.(string))
	}

	var err error
	var resp *eim.ListAttachedUserPoliciesOutput
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.API.ListAttachedUserPolicies(request)

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

	if len(resp.AttachedPolicies) < 1 {
		return fmt.Errorf("No results")
	}

	p := make([]map[string]interface{}, len(resp.AttachedPolicies))
	for k, v := range resp.AttachedPolicies {
		a := make(map[string]interface{})
		a["policy_arn"] = aws.StringValue(v.PolicyArn)
		a["policy_name"] = aws.StringValue(v.PolicyName)
		p[k] = a
	}

	d.SetId(resource.UniqueId())

	return d.Set("attached_policies", p)
}
