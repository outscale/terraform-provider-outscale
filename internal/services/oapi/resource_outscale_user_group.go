package oapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func ResourceOutscaleUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleUserGroupCreate,
		Read:   ResourceOutscaleUserGroupRead,
		Update: ResourceOutscaleUserGroupUpdate,
		Delete: ResourceOutscaleUserGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"user_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					pathVal := val.(string)
					if err := utils.CheckPath(pathVal); err != nil {
						errs = append(errs, fmt.Errorf("%v, got:%v", err, pathVal))
					}
					return
				},
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
			"policy": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_orn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"default_version_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_id": {
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
			"user": {
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
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								pathVal := val.(string)
								if err := utils.CheckPath(pathVal); err != nil {
									errs = append(errs, fmt.Errorf("%v, got:%v", err, pathVal))
								}
								return
							},
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
	conn := meta.(*client.OutscaleClient).OSCAPI

	req := oscgo.NewCreateUserGroupRequest(d.Get("user_group_name").(string))
	groupPath := d.Get("path").(string)
	req.SetPath(groupPath)

	var resp oscgo.CreateUserGroupResponse
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
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

	d.SetId(*resp.GetUserGroup().UserGroupId)
	if usersToAdd, ok := d.GetOk("user"); ok {
		reqUserAdd := oscgo.AddUserToUserGroupRequest{}
		reqUserAdd.SetUserGroupName(d.Get("user_group_name").(string))
		reqUserAdd.SetUserGroupPath(groupPath)

		for _, v := range usersToAdd.(*schema.Set).List() {
			user := v.(map[string]interface{})
			if userName := user["user_name"].(string); userName != "" {
				reqUserAdd.SetUserName(userName)
			}

			if path := user["path"].(string); path != "" {
				reqUserAdd.SetUserPath(path)
			}
			reqUserAdd.SetUserGroupPath(groupPath)
			err := retry.Retry(1*time.Minute, func() *retry.RetryError {
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
	if policiesToAdd, ok := d.GetOk("policy"); ok {
		reqAddPolicy := oscgo.LinkManagedPolicyToUserGroupRequest{}

		for _, v := range policiesToAdd.(*schema.Set).List() {
			policy := v.(map[string]interface{})

			reqAddPolicy.SetUserGroupName(d.Get("user_group_name").(string))
			reqAddPolicy.SetPolicyOrn(policy["policy_orn"].(string))
			err := retry.Retry(1*time.Minute, func() *retry.RetryError {
				_, httpResp, err := conn.PolicyApi.LinkManagedPolicyToUserGroup(context.Background()).LinkManagedPolicyToUserGroupRequest(reqAddPolicy).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})
			if err != nil {
				return err
			}
			if versionId := policy["default_version_id"].(string); versionId != "" {
				if err := setDefaultPolicyVersion(conn, policy["policy_orn"].(string), versionId); err != nil {
					return err
				}
			}
		}
	}
	return ResourceOutscaleUserGroupRead(d, meta)
}

func ResourceOutscaleUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	filter := oscgo.FiltersUserGroup{
		UserGroupIds: &[]string{d.Id()},
	}
	req := oscgo.ReadUserGroupsRequest{}
	req.SetFilters(filter)
	var statusCode int
	var resp oscgo.ReadUserGroupsResponse
	err := retry.Retry(1*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.UserGroupApi.ReadUserGroups(context.Background()).ReadUserGroupsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		statusCode = httpResp.StatusCode
		return nil
	})
	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return err
	}

	if groups, ok := resp.GetUserGroupsOk(); !ok || len(*groups) == 0 {
		d.SetId("")
		return nil
	}
	userGroup := resp.GetUserGroups()[0]
	reqUser := oscgo.NewReadUserGroupRequest(userGroup.GetName())

	var groupUsers []oscgo.User
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.UserGroupApi.ReadUserGroup(context.Background()).ReadUserGroupRequest(*reqUser).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		groupUsers = rp.GetUsers()
		statusCode = httpResp.StatusCode
		return nil
	})
	if err != nil {
		return err
	}

	linkReq := oscgo.NewReadManagedPoliciesLinkedToUserGroupRequest(userGroup.GetName())
	var linkResp oscgo.ReadManagedPoliciesLinkedToUserGroupResponse
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.PolicyApi.ReadManagedPoliciesLinkedToUserGroup(context.Background()).ReadManagedPoliciesLinkedToUserGroupRequest(*linkReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		linkResp = rp
		return nil
	})
	if err != nil {
		return err
	}

	if err := d.Set("user_group_name", userGroup.GetName()); err != nil {
		return err
	}
	if err := d.Set("path", userGroup.GetPath()); err != nil {
		return err
	}
	if err := d.Set("user_group_id", userGroup.GetUserGroupId()); err != nil {
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

	users := make([]map[string]interface{}, len(groupUsers))
	if len(groupUsers) > 0 {
		for i, v := range groupUsers {
			user := make(map[string]interface{})
			user["user_id"] = v.GetUserId()
			user["user_name"] = v.GetUserName()
			user["path"] = v.GetPath()
			user["creation_date"] = v.GetCreationDate()
			user["last_modification_date"] = v.GetLastModificationDate()
			users[i] = user
		}
	}
	if err := d.Set("user", users); err != nil {
		return err
	}

	gPolicies := linkResp.GetPolicies()
	policies := make([]map[string]interface{}, len(gPolicies))
	if len(gPolicies) > 0 {
		for i, v := range gPolicies {
			policy := make(map[string]interface{})
			policy["policy_id"] = v.GetPolicyId()
			policy["policy_name"] = v.GetPolicyName()
			policy["policy_orn"] = v.GetOrn()
			policy["creation_date"] = v.GetCreationDate()
			policy["last_modification_date"] = v.GetLastModificationDate()
			versionId, err := getPolicyVersion(conn, v.GetOrn())
			if err != nil {
				return err
			}
			policy["default_version_id"] = versionId
			policies[i] = policy
		}
	}
	if err := d.Set("policy", policies); err != nil {
		return err
	}
	return nil
}

func ResourceOutscaleUserGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	req := oscgo.UpdateUserGroupRequest{}
	isUpdateGroup := false
	if d.HasChange("user_group_name") {
		oldName, newName := d.GetChange("user_group_name")
		req.SetUserGroupName(oldName.(string))
		req.SetNewUserGroupName(newName.(string))
		isUpdateGroup = true
	}
	if d.HasChange("path") {
		oldPath, newPath := d.GetChange("path")
		oldName, _ := d.GetChange("user_group_name")
		req.SetPath(oldPath.(string))
		req.SetNewPath(newPath.(string))
		req.SetUserGroupName(oldName.(string))
		isUpdateGroup = true
	}
	if d.HasChange("user") {
		oldUsers, newUsers := d.GetChange("user")
		inter := oldUsers.(*schema.Set).Intersection(newUsers.(*schema.Set))
		toCreate := newUsers.(*schema.Set).Difference(inter)
		toRemove := oldUsers.(*schema.Set).Difference(inter)

		if len(toRemove.List()) > 0 {
			rmUserReq := oscgo.RemoveUserFromUserGroupRequest{}
			oldN, _ := d.GetChange("user_group_name")
			rmUserReq.SetUserGroupName(oldN.(string))
			oldP, _ := d.GetChange("path")
			rmUserReq.SetUserGroupPath(oldP.(string))
			_, checkUpdate, _, err := getUsersLinkedToGroup(conn, toCreate.List(), oldN.(string))
			if err != nil {
				return err
			}
			for _, v := range toRemove.List() {
				user := v.(map[string]interface{})
				if len(checkUpdate) != 0 {
					if checkUpdate[user["user_id"].(string)] {
						continue
					}
				}
				if userName := user["user_name"].(string); userName != "" {
					rmUserReq.SetUserName(userName)
				}
				if path := user["path"].(string); path != "" {
					rmUserReq.SetUserPath(path)
				}
				err = retry.Retry(2*time.Minute, func() *retry.RetryError {
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
			_, _, checkUpdate, err := getUsersLinkedToGroup(conn, toCreate.List(), oldN.(string))
			if err != nil {
				return err
			}
			for _, v := range toCreate.List() {
				user := v.(map[string]interface{})
				if len(checkUpdate) != 0 {
					if checkUpdate[user["user_name"].(string)] {
						continue
					}
				}
				if userName := user["user_name"].(string); userName != "" {
					addUserReq.SetUserName(userName)
				}
				if path := user["path"].(string); path != "" {
					addUserReq.SetUserPath(path)
				}
				err := retry.Retry(2*time.Minute, func() *retry.RetryError {
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

	if d.HasChange("policy") {
		oldPolicies, newPolicies := d.GetChange("policy")
		inter := oldPolicies.(*schema.Set).Intersection(newPolicies.(*schema.Set))
		toCreate := newPolicies.(*schema.Set).Difference(inter)
		toRemove := oldPolicies.(*schema.Set).Difference(inter)

		if len(toRemove.List()) > 0 {
			unlinkReq := oscgo.UnlinkManagedPolicyFromUserGroupRequest{}
			oldN, _ := d.GetChange("user_group_name")
			unlinkReq.SetUserGroupName(oldN.(string))
			for _, v := range toRemove.List() {
				policy := v.(map[string]interface{})
				unlinkReq.SetPolicyOrn(policy["policy_orn"].(string))
				err := retry.Retry(2*time.Minute, func() *retry.RetryError {
					_, httpResp, err := conn.PolicyApi.UnlinkManagedPolicyFromUserGroup(context.Background()).UnlinkManagedPolicyFromUserGroupRequest(unlinkReq).Execute()
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
			linkReq := oscgo.LinkManagedPolicyToUserGroupRequest{}
			oldN, _ := d.GetChange("user_group_name")
			linkReq.SetUserGroupName(oldN.(string))

			for _, v := range toCreate.List() {
				policy := v.(map[string]interface{})
				linkReq.SetPolicyOrn(policy["policy_orn"].(string))
				err := retry.Retry(2*time.Minute, func() *retry.RetryError {
					_, httpResp, err := conn.PolicyApi.LinkManagedPolicyToUserGroup(context.Background()).LinkManagedPolicyToUserGroupRequest(linkReq).Execute()
					if err != nil {
						return utils.CheckThrottling(httpResp, err)
					}
					return nil
				})
				if err != nil {
					return err
				}
				if versionId := policy["default_version_id"].(string); versionId != "" {
					if err := setDefaultPolicyVersion(conn, policy["policy_orn"].(string), versionId); err != nil {
						return err
					}
				}
			}
		}
	}
	if isUpdateGroup {
		err := retry.Retry(2*time.Minute, func() *retry.RetryError {
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
	conn := meta.(*client.OutscaleClient).OSCAPI
	forceDeletion := true
	req := oscgo.DeleteUserGroupRequest{
		UserGroupName: d.Get("user_group_name").(string),
		Force:         &forceDeletion,
	}

	err := retry.Retry(3*time.Minute, func() *retry.RetryError {
		_, httpResp, err := conn.UserGroupApi.DeleteUserGroup(context.Background()).DeleteUserGroupRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting outscale user group %s: %s", d.Id(), err)
	}

	return nil
}

func getPolicyVersion(conn *oscgo.APIClient, policyOrn string) (string, error) {
	version_id := ""
	err := retry.Retry(3*time.Minute, func() *retry.RetryError {
		resp, httpResp, err := conn.PolicyApi.ReadPolicy(context.Background()).ReadPolicyRequest(
			oscgo.ReadPolicyRequest{PolicyOrn: policyOrn}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		policy := resp.GetPolicy()
		version_id = policy.GetPolicyDefaultVersionId()
		return nil
	})
	if err != nil {
		return version_id, err
	}
	return version_id, err
}

func setDefaultPolicyVersion(conn *oscgo.APIClient, policyOrn, version string) error {
	err := retry.Retry(1*time.Minute, func() *retry.RetryError {
		_, httpResp, err := conn.PolicyApi.SetDefaultPolicyVersion(context.Background()).SetDefaultPolicyVersionRequest(
			oscgo.SetDefaultPolicyVersionRequest{
				PolicyOrn: policyOrn,
				VersionId: version,
			}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	return err
}

func getUsersLinkedToGroup(conn *oscgo.APIClient, toCreate []interface{}, groupName string) ([]oscgo.User, map[string]bool, map[string]bool, error) {
	reqUser := oscgo.NewReadUserGroupRequest(groupName)

	var users []oscgo.User
	err := retry.Retry(1*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.UserGroupApi.ReadUserGroup(context.Background()).ReadUserGroupRequest(*reqUser).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		users = rp.GetUsers()
		return nil
	})
	checkIds := make(map[string]bool)
	checkName := make(map[string]bool)
	for _, user := range users {
		if len(toCreate) > 0 {
			for _, v := range toCreate {
				userTo := v.(map[string]interface{})
				if userTo["user_name"] == user.GetUserName() {
					checkIds[user.GetUserId()] = true
					checkName[user.GetUserName()] = true
				}
			}
		}
	}
	return users, checkIds, checkName, err
}
