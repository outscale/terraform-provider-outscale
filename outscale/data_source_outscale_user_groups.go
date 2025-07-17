package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func DataSourceUserGroups() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceUserGroupsRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
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

func DataSourceUserGroupsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.ReadUserGroupsRequest{}
	filters, filtersOk := d.GetOk("filter")
	if filtersOk {
		filterReq := buildUserGroupsFilters(filters.(*schema.Set))
		req.SetFilters(*filterReq)
	}
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

func buildUserGroupsFilters(set *schema.Set) *oscgo.FiltersUserGroup {
	var filters oscgo.FiltersUserGroup
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "path_prefix":
			filters.SetPathPrefix(filterValues[0])
		case "user_group_ids":
			filters.SetUserGroupIds(filterValues)
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return &filters
}
