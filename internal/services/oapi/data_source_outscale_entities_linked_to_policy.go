package oapi

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceEntitiesLinkedToPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceEntitiesLinkedToPoliciesRead,
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

func DataSourceEntitiesLinkedToPoliciesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	orn := d.Get("policy_orn").(string)
	req := osc.ReadEntitiesLinkedToPolicyRequest{PolicyOrn: orn}
	if entities := utils.SetToSuperStringSlice[osc.ReadEntitiesLinkedToPolicyRequestEntitiesType](d.Get("entities_type").(*schema.Set)); len(entities) > 0 {
		req.EntitiesType = &entities
	}

	resp, err := client.ReadEntitiesLinkedToPolicy(ctx, req, options.WithRetryTimeout(2*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.PolicyEntities == nil {
		return diag.Errorf("unable to find entities linked to policy")
	}
	d.SetId(id.UniqueId())

	entities := *resp.PolicyEntities

	users := make([]map[string]interface{}, len(ptr.From(entities.Users)))
	groups := make([]map[string]interface{}, len(ptr.From(entities.Groups)))
	accounts := make([]map[string]interface{}, len(ptr.From(entities.Accounts)))
	if entities.Users != nil {
		for i, v := range *entities.Users {
			user := make(map[string]interface{})
			user["id"] = v.Id
			user["name"] = v.Name
			user["orn"] = v.Orn
			users[i] = user
		}
	}
	if entities.Groups != nil {
		for i, v := range *entities.Groups {
			group := make(map[string]interface{})
			group["name"] = v.Name
			group["id"] = v.Id
			group["orn"] = v.Orn
			groups[i] = group
		}
	}
	if entities.Accounts != nil {
		for i, v := range *entities.Accounts {
			account := make(map[string]interface{})
			account["name"] = v.Name
			account["id"] = v.Id
			account["orn"] = v.Orn
			accounts[i] = account
		}
	}

	return diag.FromErr(d.Set("policy_entities", []map[string]interface{}{{
		"users":    users,
		"groups":   groups,
		"accounts": accounts,
	}}))
}
