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

func DataSourceUserGroupsPerUser() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceUserGroupsPerUserRead,
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

func DataSourceUserGroupsPerUserRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	req := oscgo.NewReadUserGroupsPerUserRequest(d.Get("user_name").(string))
	if userPath := d.Get("user_path").(string); userPath != "" {
		req.SetUserPath(userPath)
	}
	var resp oscgo.ReadUserGroupsPerUserResponse
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.UserGroupApi.ReadUserGroupsPerUser(context.Background()).ReadUserGroupsPerUserRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	if _, ok := resp.GetUserGroupsOk(); !ok {
		return fmt.Errorf("unable to find user groups")
	}
	d.SetId(id.UniqueId())
	userGps := resp.GetUserGroups()
	userGroups := make([]map[string]interface{}, len(userGps))

	for i, v := range userGps {
		userGroup := make(map[string]interface{})
		userGroup["user_group_name"] = v.GetName()
		userGroup["user_group_id"] = v.GetUserGroupId()
		userGroup["path"] = v.GetPath()
		userGroup["orn"] = v.GetOrn()
		userGroup["creation_date"] = v.GetCreationDate()
		userGroup["last_modification_date"] = v.GetLastModificationDate()
		userGroups[i] = userGroup
	}
	return d.Set("user_groups", userGroups)
}
