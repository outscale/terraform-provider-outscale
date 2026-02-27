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

func DataSourceUserGroupsPerUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceUserGroupsPerUserRead,
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
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
					},
				},
			},
		},
	}
}

func DataSourceUserGroupsPerUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadUserGroupsPerUserRequest{
		UserName: d.Get("user_name").(string),
	}
	if userPath := d.Get("user_path").(string); userPath != "" {
		req.UserPath = &userPath
	}
	resp, err := client.ReadUserGroupsPerUser(ctx, req, options.WithRetryTimeout(2*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.UserGroups == nil {
		return diag.Errorf("unable to find user groups")
	}
	d.SetId(id.UniqueId())
	userGps := ptr.From(resp.UserGroups)
	userGroups := make([]map[string]interface{}, len(userGps))

	for i, v := range userGps {
		userGroup := make(map[string]interface{})
		userGroup["user_group_name"] = v.Name
		userGroup["user_group_id"] = v.UserGroupId
		userGroup["path"] = v.Path
		userGroup["orn"] = v.Orn
		userGroup["creation_date"] = from.ISO8601(v.CreationDate)
		userGroup["last_modification_date"] = from.ISO8601(v.LastModificationDate)
		userGroups[i] = userGroup
	}
	return diag.FromErr(d.Set("user_groups", userGroups))
}
