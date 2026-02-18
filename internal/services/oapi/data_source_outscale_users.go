package oapi

import (
	"context"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceUsers() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceUsersRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
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
						"user_email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
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

func DataSourceUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.OutscaleClient).OSC
	
	filters, filtersOk := d.GetOk("filter")

	var err error
	req := osc.NewReadUsersRequest()
	if filtersOk {
		req.Filters, err = buildUsersFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp osc.ReadUsersResponse
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := client.UserApi.ReadUsers(ctx).ReadUsersRequest(*req).Execute()
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
	d.SetId(id.UniqueId())
	if len(users) == 0 {
		return ErrNoResults
	}
	d.SetId(id.UniqueId())
	usersToSet := make([]map[string]interface{}, len(users))
	for i, v := range users {
		user := make(map[string]interface{})

		user["user_id"] = v.GetUserId()
		user["user_name"] = v.GetUserName()
		user["user_email"] = v.GetUserEmail()
		user["path"] = v.GetPath()
		user["creation_date"] = v.GetCreationDate()
		user["last_modification_date"] = v.GetLastModificationDate()
		usersToSet[i] = user
	}
	return d.Set("users", usersToSet)
}
