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

func dataSourceOutscaleOAPIPolicyUserLink() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIPolicyUserLinkRead,

		Schema: map[string]*schema.Schema{
			"user_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_arn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPIPolicyUserLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	user := d.Get("user_name").(string)

	request := &eim.ListAttachedUserPoliciesInput{
		UserName: aws.String(user),
	}

	if v, ok := d.GetOk("path"); ok {
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

	if len(resp.ListAttachedUserPoliciesResult.AttachedPolicies) < 1 {
		return fmt.Errorf("No results")
	}

	p := resp.ListAttachedUserPoliciesResult.AttachedPolicies

	d.Set("policy_arn", aws.StringValue(p[0].PolicyArn))
	d.Set("policy_name", aws.StringValue(p[0].PolicyName))
	d.SetId(resource.UniqueId())
	d.Set("request_id", resp.ResponseMetadata.RequestID)

	return d.Set("attached_policies", p)
}
