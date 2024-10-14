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

func DataSourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceUserGroupRead,
		Schema: map[string]*schema.Schema{
			"user_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_group_id": {
				Type:     schema.TypeString,
				Required: true,
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

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return fmt.Errorf("One of filters, user_group_id or path_refix must be assigned")
	}

	filterReq := buildUserGroupsFilters(filters.(*schema.Set))
	req := oscgo.ReadUserGroupsRequest{}
	req.SetFilters(*filterReq)
	var resp oscgo.ReadUserGroupsResponse
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.UserGroupApi.ReadUserGroups(context.Background()).ReadUserGroupsRequest(req).Execute()
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
		return fmt.Errorf("Unable to find user groups")
	}
	d.SetId(resource.UniqueId())
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
