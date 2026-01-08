package oapi

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func ResourceOutscaleUser() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleUserCreate,
		Read:   ResourceOutscaleUserRead,
		Update: ResourceOutscaleUserUpdate,
		Delete: ResourceOutscaleUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"user_name": {
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
			"user_email": {
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
		},
	}
}

func ResourceOutscaleUserCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	req := oscgo.NewCreateUserRequest(d.Get("user_name").(string))
	req.SetPath(d.Get("path").(string))
	if email := d.Get("user_email").(string); email != "" {
		req.SetUserEmail(email)
	}

	var resp oscgo.CreateUserResponse
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.UserApi.CreateUser(context.Background()).CreateUserRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(*resp.GetUser().UserId)
	if policiesToAdd, ok := d.GetOk("policy"); ok {
		reqAddPolicy := oscgo.LinkPolicyRequest{}

		for _, v := range policiesToAdd.(*schema.Set).List() {
			policy := v.(map[string]interface{})
			reqAddPolicy.SetUserName(d.Get("user_name").(string))
			reqAddPolicy.SetPolicyOrn(policy["policy_orn"].(string))

			err := retry.Retry(2*time.Minute, func() *retry.RetryError {
				_, httpResp, err := conn.PolicyApi.LinkPolicy(context.Background()).LinkPolicyRequest(reqAddPolicy).Execute()
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

	return ResourceOutscaleUserRead(d, meta)
}

func ResourceOutscaleUserRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	req := oscgo.ReadUsersRequest{
		Filters: &oscgo.FiltersUsers{UserIds: &[]string{d.Id()}},
	}

	var resp oscgo.ReadUsersResponse
	err := retry.Retry(1*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.UserApi.ReadUsers(context.Background()).ReadUsersRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}
	users := resp.GetUsers()
	if len(users) == 0 {
		d.SetId("")
		return nil
	}
	linkReq := oscgo.NewReadLinkedPoliciesRequest(users[0].GetUserName())
	var linkResp oscgo.ReadLinkedPoliciesResponse
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.PolicyApi.ReadLinkedPolicies(context.Background()).ReadLinkedPoliciesRequest(*linkReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		linkResp = rp
		return nil
	})
	if err != nil {
		return err
	}

	if err := d.Set("user_name", users[0].GetUserName()); err != nil {
		return err
	}
	if err := d.Set("user_id", users[0].GetUserId()); err != nil {
		return err
	}
	if err := d.Set("path", users[0].GetPath()); err != nil {
		return err
	}
	if err := d.Set("user_email", users[0].GetUserEmail()); err != nil {
		return err
	}
	if err := d.Set("creation_date", users[0].GetCreationDate()); err != nil {
		return err
	}
	if err := d.Set("last_modification_date", users[0].GetLastModificationDate()); err != nil {
		return err
	}

	uPolicies := linkResp.GetPolicies()
	policies := make([]map[string]interface{}, len(uPolicies))
	if len(uPolicies) > 0 {
		for i, v := range uPolicies {
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

func ResourceOutscaleUserUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	req := oscgo.UpdateUserRequest{UserName: d.Get("user_name").(string)}
	isUpdateUser := false
	if d.HasChange("user_name") {
		oldN, newN := d.GetChange("user_name")
		if oldName := oldN.(string); oldName != "" {
			req.SetUserName(oldName)
			isUpdateUser = true
		}
		if newName := newN.(string); newName != "" && oldN.(string) != newN.(string) {
			req.SetNewUserName(newName)
			isUpdateUser = true
		}
	}
	if d.HasChange("path") {
		req.SetNewPath(d.Get("path").(string))
		if req.GetUserName() == "" {
			req.SetUserName(d.Get("user_name").(string))
		}
		isUpdateUser = true
	}
	if d.HasChange("user_email") {
		_, newM := d.GetChange("user_email")
		req.SetNewUserEmail(newM.(string))
		isUpdateUser = true
	}
	if d.HasChange("policy") {
		oldPolicies, newPolicies := d.GetChange("policy")
		inter := oldPolicies.(*schema.Set).Intersection(newPolicies.(*schema.Set))
		toCreate := newPolicies.(*schema.Set).Difference(inter)
		toRemove := oldPolicies.(*schema.Set).Difference(inter)

		if len(toRemove.List()) > 0 {
			unlinkReq := oscgo.UnlinkPolicyRequest{}
			oldN, _ := d.GetChange("user_name")
			unlinkReq.SetUserName(oldN.(string))

			for _, v := range toRemove.List() {
				policy := v.(map[string]interface{})
				unlinkReq.SetPolicyOrn(policy["policy_orn"].(string))
				err := retry.Retry(2*time.Minute, func() *retry.RetryError {
					_, httpResp, err := conn.PolicyApi.UnlinkPolicy(context.Background()).UnlinkPolicyRequest(unlinkReq).Execute()
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
			linkReq := oscgo.LinkPolicyRequest{}
			oldN, _ := d.GetChange("user_name")
			linkReq.SetUserName(oldN.(string))

			for _, v := range toCreate.List() {
				policy := v.(map[string]interface{})
				linkReq.SetPolicyOrn(policy["policy_orn"].(string))

				err := retry.Retry(2*time.Minute, func() *retry.RetryError {
					_, httpResp, err := conn.PolicyApi.LinkPolicy(context.Background()).LinkPolicyRequest(linkReq).Execute()
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
	if isUpdateUser {
		err := retry.Retry(2*time.Minute, func() *retry.RetryError {
			_, httpResp, err := conn.UserApi.UpdateUser(context.Background()).UpdateUserRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return ResourceOutscaleUserRead(d, meta)
}

func ResourceOutscaleUserDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	if _, ok := d.GetOk("policy"); ok {
		policies := d.Get("policy")
		for _, v := range policies.(*schema.Set).List() {
			unlinkReq := oscgo.UnlinkPolicyRequest{}
			unlinkReq.SetUserName(d.Get("user_name").(string))
			policy := v.(map[string]interface{})
			unlinkReq.SetPolicyOrn(policy["policy_orn"].(string))
			err := retry.Retry(2*time.Minute, func() *retry.RetryError {
				_, httpResp, err := conn.PolicyApi.UnlinkPolicy(context.Background()).UnlinkPolicyRequest(unlinkReq).Execute()
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

	req := oscgo.DeleteUserRequest{
		UserName: d.Get("user_name").(string),
	}
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, httpResp, err := conn.UserApi.DeleteUser(context.Background()).DeleteUserRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting outscale user %s: %s", d.Id(), err)
	}

	return nil
}
