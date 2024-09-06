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

func ResourceOutscaleUser() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleUserCreate,
		Read:   ResourceOutscaleUserRead,
		Update: ResourceOutscaleUserUpdate,
		Delete: ResourceOutscaleUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func ResourceOutscaleUserCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.NewCreateUserRequest(d.Get("user_name").(string))
	if path := d.Get("path").(string); path != "" {
		req.Path = &path
	}

	var resp oscgo.CreateUserResponse
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
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
	d.SetId(resource.UniqueId())
	if err := d.Set("user_id", *resp.GetUser().UserId); err != nil {
		return err
	}

	return ResourceOutscaleUserRead(d, meta)
}

func ResourceOutscaleUserRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.NewReadUsersRequest()

	var resp oscgo.ReadUsersResponse
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.UserApi.ReadUsers(context.Background()).ReadUsersRequest(*req).Execute()
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
	for _, user := range users {
		if user.GetUserId() == d.Get("user_id") {

			if err := d.Set("user_name", user.GetUserName()); err != nil {
				return err
			}
			if err := d.Set("user_id", user.GetUserId()); err != nil {
				return err
			}
			if err := d.Set("path", user.GetPath()); err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func ResourceOutscaleUserUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.UpdateUserRequest{}

	oldN, newN := d.GetChange("user_name")
	if oldName := oldN.(string); oldName != "" {
		req.SetUserName(oldName)
	}
	if newName := newN.(string); newName != "" && oldN.(string) != newN.(string) {
		req.SetNewUserName(newName)
	}

	if d.HasChange("path") {
		path := d.Get("path").(string)
		req.NewPath = &path
	}

	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.UserApi.UpdateUser(context.Background()).UpdateUserRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return ResourceOutscaleUserRead(d, meta)
}

func ResourceOutscaleUserDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.DeleteUserRequest{
		UserName: d.Get("user_name").(string),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.UserApi.DeleteUser(context.Background()).DeleteUserRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error deleting Outscale Access Key %s: %s", d.Id(), err)
	}

	return nil
}
