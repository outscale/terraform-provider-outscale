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
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
)

func DataSourceUserGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceUserGroupRead,
		Schema: map[string]*schema.Schema{
			"user_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_group_id": {
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
			"user": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
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

func DataSourceUserGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadUserGroupRequest{
		UserGroupName: d.Get("user_group_name").(string),
	}
	if path := d.Get("path").(string); path != "" {
		req.Path = &path
	}
	resp, err := client.ReadUserGroup(ctx, req, options.WithRetryTimeout(2*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.UserGroup == nil {
		return diag.FromErr(ErrNoResults)
	}
	d.SetId(id.UniqueId())
	group := ptr.From(resp.UserGroup)
	users := ptr.From(resp.Users)

	if err := d.Set("user_group_name", ptr.From(group.Name)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("user_group_id", ptr.From(group.UserGroupId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("orn", ptr.From(group.Orn)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("path", ptr.From(group.Path)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("creation_date", from.ISO8601(group.CreationDate)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_modification_date", from.ISO8601(group.LastModificationDate)); err != nil {
		return diag.FromErr(err)
	}
	if len(users) > 0 {
		userState := make([]map[string]interface{}, len(users))

		for i, v := range users {
			user := make(map[string]interface{})
			user["user_name"] = v.UserName
			user["user_id"] = v.UserId
			user["path"] = v.Path
			user["user_email"] = v.UserEmail
			user["creation_date"] = from.ISO8601(v.CreationDate)
			user["last_modification_date"] = from.ISO8601(v.LastModificationDate)
			userState[i] = user
		}
		if err := d.Set("user", userState); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}
