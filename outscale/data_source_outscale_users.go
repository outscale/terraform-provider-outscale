package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceUsers() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceUsersRead,
		Schema: map[string]*schema.Schema{
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceUsersRead(d *schema.ResourceData, meta interface{}) error {
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
	d.SetId(resource.UniqueId())
	if len(users) == 0 {
		return fmt.Errorf("Unable to find users")
	}
	d.SetId(resource.UniqueId())
	usersToSet := make([]map[string]interface{}, len(users))
	for i, v := range users {
		user := make(map[string]interface{})

		user["user_id"] = v.GetUserId()
		user["user_name"] = v.GetUserName()
		user["path"] = v.GetPath()
		usersToSet[i] = user
	}
	return d.Set("users", usersToSet)
}
