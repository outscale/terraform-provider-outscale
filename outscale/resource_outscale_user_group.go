package outscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func ResourceOutscaleUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleUserGroupCreate,
		Read:   ResourceOutscaleUserGroupRead,
		Update: ResourceOutscaleUserGroupUpdate,
		Delete: ResourceOutscaleUserGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
			"users": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
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

func ResourceOutscaleUserGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.NewCreateUserGroupRequest(d.Get("user_group_name").(string))
	groupPath := d.Get("path").(string)
	if groupPath != "" {
		req.Path = &groupPath
	}
	var resp oscgo.CreateUserGroupResponse
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.UserGroupApi.CreateUserGroup(context.Background()).CreateUserGroupRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(resource.UniqueId())
	// Remove d.Set when read user_group return user_group_id
	if err := d.Set("user_group_id", resp.GetUserGroup().UserGroupId); err != nil {
		return err
	}
	if usersToAdd, ok := d.GetOk("users"); ok {
		reqUserAdd := oscgo.AddUserToUserGroupRequest{}
		reqUserAdd.SetUserGroupName(d.Get("user_group_name").(string))
		if path := d.Get("path").(string); path != "" {
			reqUserAdd.UserGroupPath = &path
		}

		for _, v := range usersToAdd.(*schema.Set).List() {
			user := v.(map[string]interface{})
			if userName := user["user_name"].(string); userName != "" {
				reqUserAdd.SetUserName(userName)
			}
			if path := user["path"].(string); path != "" {
				reqUserAdd.SetUserPath(path)
			}
			if groupPath != "" {
				reqUserAdd.SetUserGroupPath(groupPath)
			}
			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, httpResp, err := conn.UserGroupApi.AddUserToUserGroup(context.Background()).AddUserToUserGroupRequest(reqUserAdd).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})
			if err != nil {
				return err
			}

		}

	}
	return ResourceOutscaleUserGroupRead(d, meta)
}

func ResourceOutscaleUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.NewReadUserGroupRequest(d.Get("user_group_name").(string))
	if path := d.Get("path").(string); path != "" {
		req.Path = &path
	}

	var resp oscgo.ReadUserGroupResponse
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
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
		d.SetId("")
		return nil
	}
	userGroup := resp.GetUserGroup()
	if err := d.Set("user_group_name", userGroup.GetName()); err != nil {
		return err
	}
	/* Remove comment when read user_group return user_group_id
	if err := d.Set("user_group_id", userGroup.GetUserGroupId()); err != nil {
		return err
	}
	*/
	if err := d.Set("path", userGroup.GetPath()); err != nil {
		return err
	}
	if err := d.Set("orn", userGroup.GetOrn()); err != nil {
		return err
	}
	if err := d.Set("creation_date", (userGroup.GetCreationDate())); err != nil {
		return err
	}
	if err := d.Set("last_modification_date", (userGroup.GetLastModificationDate())); err != nil {
		return err
	}
	usrs := resp.GetUsers()
	users := make([]map[string]interface{}, len(usrs))
	if len(usrs) > 0 {
		usrs := resp.GetUsers()
		for i, v := range usrs {
			user := make(map[string]interface{})
			user["user_id"] = v.GetUserId()
			user["user_name"] = v.GetUserName()
			user["path"] = v.GetPath()
			user["creation_date"] = v.GetCreationDate()
			user["last_modification_date"] = v.GetLastModificationDate()
			users[i] = user
		}
	}
	if err := d.Set("users", users); err != nil {
		return err
	}
	return nil
}

func ResourceOutscaleUserGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.UpdateUserGroupRequest{}
	isUpdateGroup := false

	oldN, newN := d.GetChange("user_group_name")
	if oldName := oldN.(string); oldName != "" {
		req.SetUserGroupName(oldName)
	}
	if newName := newN.(string); newName != "" && oldN.(string) != newN.(string) {
		req.SetNewUserGroupName(newName)
		isUpdateGroup = true
	}

	oldP, newP := d.GetChange("path")
	if oldPath := oldP.(string); oldPath != "" {
		req.SetPath(oldPath)
	}
	if newPath := newP.(string); newPath != "" && oldP.(string) != newP.(string) {
		req.SetNewPath(newPath)
		isUpdateGroup = true
	}
	if d.HasChange("users") {
		oldUsers, newUsers := d.GetChange("users")
		inter := oldUsers.(*schema.Set).Intersection(newUsers.(*schema.Set))
		toCreate := newUsers.(*schema.Set).Difference(inter)
		toRemove := oldUsers.(*schema.Set).Difference(inter)

		if len(toRemove.List()) > 0 {
			rmUserReq := oscgo.RemoveUserFromUserGroupRequest{}
			oldN, _ := d.GetChange("user_group_name")
			rmUserReq.SetUserGroupName(oldN.(string))
			oldP, _ := d.GetChange("path")
			rmUserReq.SetUserGroupPath(oldP.(string))

			for _, v := range toRemove.List() {
				user := v.(map[string]interface{})
				if userName := user["user_name"].(string); userName != "" {
					rmUserReq.SetUserName(userName)
				}
				if path := user["path"].(string); path != "" {
					rmUserReq.SetUserGroupPath(path)
				}
				err := resource.Retry(5*time.Minute, func() *resource.RetryError {
					_, httpResp, err := conn.UserGroupApi.RemoveUserFromUserGroup(context.Background()).RemoveUserFromUserGroupRequest(rmUserReq).Execute()
					if err != nil {
						return utils.CheckThrottling(httpResp, err)
					}
					return nil
				})
				if err != nil {
					return err
				}
			}
		}
		if len(toCreate.List()) > 0 {
			addUserReq := oscgo.AddUserToUserGroupRequest{}
			oldN, _ := d.GetChange("user_group_name")
			addUserReq.SetUserGroupName(oldN.(string))
			oldP, _ := d.GetChange("path")
			addUserReq.SetUserGroupPath(oldP.(string))

			for _, v := range toCreate.List() {
				user := v.(map[string]interface{})
				if userName := user["user_name"].(string); userName != "" {
					addUserReq.SetUserName(userName)
				}
				if path := user["path"].(string); path != "" {
					addUserReq.SetUserGroupPath(path)
				}
				err := resource.Retry(2*time.Minute, func() *resource.RetryError {
					_, httpResp, err := conn.UserGroupApi.AddUserToUserGroup(context.Background()).AddUserToUserGroupRequest(addUserReq).Execute()
					if err != nil {
						return utils.CheckThrottling(httpResp, err)
					}
					return nil
				})
				if err != nil {
					return err
				}

			}
		}
	}
	if isUpdateGroup {
		err := resource.Retry(2*time.Minute, func() *resource.RetryError {
			_, httpResp, err := conn.UserGroupApi.UpdateUserGroup(context.Background()).UpdateUserGroupRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return ResourceOutscaleUserGroupRead(d, meta)
}

func ResourceOutscaleUserGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	forceDeletion := true
	req := oscgo.DeleteUserGroupRequest{
		UserGroupName: d.Get("user_group_name").(string),
		Force:         &forceDeletion,
	}
	if path := d.Get("path").(string); path != "" {
		req.Path = &path
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.UserGroupApi.DeleteUserGroup(context.Background()).DeleteUserGroupRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error deleting Outscale User Group %s: %s", d.Id(), err)
	}

	return nil
}
