package oapi

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceUserGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceUserGroupsRead,
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

func DataSourceUserGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadUserGroupsRequest{}

	var err error
	filters, filtersOk := d.GetOk("filter")
	if filtersOk {
		req.Filters, err = buildUserGroupsFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadUserGroups(ctx, req, options.WithRetryTimeout(2*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.UserGroups == nil {
		return diag.Errorf("unable to find user groups")
	}
	d.SetId(id.UniqueId())
	userGps := ptr.From(resp.UserGroups)
	userGroups := make([]map[string]interface{}, len(userGps))

	for i, v := range userGps {
		userGroup := make(map[string]interface{})
		userGroup["user_group_name"] = v.Name
		userGroup["user_group_id"] = v.UserGroupId
		userGroup["path"] = v.Path
		userGroup["orn"] = v.Orn
		userGroup["creation_date"] = from.ISO8601(v.CreationDate)
		userGroup["last_modification_date"] = from.ISO8601(v.LastModificationDate)
		userGroups[i] = userGroup
	}
	return diag.FromErr(d.Set("user_groups", userGroups))
}

func buildUserGroupsFilters(set *schema.Set) (*osc.FiltersUserGroup, error) {
	var filters osc.FiltersUserGroup
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "path_prefix":
			filters.PathPrefix = &filterValues[0]
		case "user_group_ids":
			filters.UserGroupIds = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
