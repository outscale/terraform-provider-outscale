package outscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func DataSourcePoliciesLinkedToUser() *schema.Resource {
	return &schema.Resource{
		Read: DataSourcePoliciesLinkedToUserRead,
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"orn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modification_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourcePoliciesLinkedToUserRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.NewReadLinkedPoliciesRequest(d.Get("user_name").(string))
	var resp oscgo.ReadLinkedPoliciesResponse

	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.PolicyApi.ReadLinkedPolicies(context.Background()).ReadLinkedPoliciesRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}
	policiesList := resp.GetPolicies()
	if len(policiesList) == 0 {
		return fmt.Errorf("unable to find Policies linked to user: %v", d.Get("user_name").(string))
	}
	d.SetId(id.UniqueId())

	policies := make([]map[string]interface{}, len(policiesList))

	for i, v := range policiesList {
		policy := make(map[string]interface{})
		policy["policy_name"] = v.GetPolicyName()
		policy["policy_id"] = v.GetPolicyId()
		policy["orn"] = v.GetOrn()
		policy["creation_date"] = v.GetCreationDate()
		policy["last_modification_date"] = v.GetLastModificationDate()
		policies[i] = policy
	}
	return d.Set("policies", policies)
}
