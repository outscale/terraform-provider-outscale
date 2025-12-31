package oapi

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceEntitiesLinkedToPolicy() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceEntitiesLinkedToPoliciesRead,
		Schema: map[string]*schema.Schema{
			"policy_orn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entities_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"policy_entities": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"users": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"orn": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"groups": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"orn": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"accounts": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"orn": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func DataSourceEntitiesLinkedToPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	orn := d.Get("policy_orn").(string)
	req := oscgo.ReadEntitiesLinkedToPolicyRequest{PolicyOrn: orn}
	if entities := utils.SetToStringSlice(d.Get("entities_type").(*schema.Set)); len(entities) > 0 {
		req.SetEntitiesType(entities)
	}

	var resp oscgo.ReadEntitiesLinkedToPolicyResponse
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.PolicyApi.ReadEntitiesLinkedToPolicy(context.Background()).ReadEntitiesLinkedToPolicyRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}
	entities, ok := resp.GetPolicyEntitiesOk()
	if !ok {
		return fmt.Errorf("unable to find Entities linked to policy")
	}
	d.SetId(id.UniqueId())

	users := make([]map[string]interface{}, len(entities.GetUsers()))
	groups := make([]map[string]interface{}, len(entities.GetGroups()))
	accounts := make([]map[string]interface{}, len(entities.GetAccounts()))
	if respUsers, ok := entities.GetUsersOk(); ok {
		for i, v := range *respUsers {
			user := make(map[string]interface{})
			user["id"] = v.GetId()
			user["name"] = v.GetName()
			user["orn"] = v.GetOrn()
			users[i] = user
		}
	}
	if respGroups, ok := entities.GetGroupsOk(); ok {
		for i, v := range *respGroups {
			group := make(map[string]interface{})
			group["name"] = v.GetName()
			group["id"] = v.GetId()
			group["orn"] = v.GetOrn()
			groups[i] = group
		}
	}
	if respAccounts, ok := entities.GetAccountsOk(); ok {
		for i, v := range *respAccounts {
			account := make(map[string]interface{})
			account["name"] = v.GetName()
			account["id"] = v.GetId()
			account["orn"] = v.GetOrn()
			accounts[i] = account
		}
	}

	return d.Set("policy_entities", []map[string]interface{}{{
		"users":    users,
		"groups":   groups,
		"accounts": accounts,
	}})
}
