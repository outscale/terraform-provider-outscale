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

func DataSourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceUserGroupRead,
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

func DataSourceUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.NewReadUserGroupRequest(d.Get("user_group_name").(string))
	if path := d.Get("path").(string); path != "" {
		req.SetPath(path)
	}
	var resp oscgo.ReadUserGroupResponse
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.UserGroupApi.ReadUserGroup(context.Background()).ReadUserGroupRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	if _, ok := resp.GetUserGroupOk(); !ok {
		return fmt.Errorf("Unable to find user group")
	}
	d.SetId(id.UniqueId())
	group := resp.GetUserGroup()
	users := resp.GetUsers()

	if err := d.Set("user_group_name", group.GetName()); err != nil {
		return err
	}
	if err := d.Set("user_group_id", group.GetUserGroupId()); err != nil {
		return err
	}
	if err := d.Set("orn", group.GetOrn()); err != nil {
		return err
	}
	if err := d.Set("path", group.GetPath()); err != nil {
		return err
	}
	if err := d.Set("creation_date", group.GetCreationDate()); err != nil {
		return err
	}
	if err := d.Set("last_modification_date", group.GetLastModificationDate()); err != nil {
		return err
	}
	if len(users) > 0 {
		userState := make([]map[string]interface{}, len(users))

		for i, v := range users {
			user := make(map[string]interface{})
			user["user_name"] = v.GetUserName()
			user["user_id"] = v.GetUserId()
			user["path"] = v.GetPath()
			user["user_email"] = v.GetUserEmail()
			user["creation_date"] = v.GetCreationDate()
			user["last_modification_date"] = v.GetLastModificationDate()
			userState[i] = user
		}
		if err := d.Set("user", userState); err != nil {
			return err
		}
	}
	return nil
}
