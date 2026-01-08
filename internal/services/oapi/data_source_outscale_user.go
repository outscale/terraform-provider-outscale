package oapi

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
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
	}
}

func DataSourceUserRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return fmt.Errorf("filters: user_ids must be assigned")
	}

	var err error
	req := oscgo.NewReadUsersRequest()

	req.Filters, err = buildUsersFilters(filters.(*schema.Set))
	if err != nil {
		return err
	}

	var resp oscgo.ReadUsersResponse
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
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
	d.SetId(id.UniqueId())
	if len(users) == 0 {
		return ErrNoResults
	}
	if len(users) > 1 {
		return ErrMultipleResults
	}

	if err := d.Set("user_name", users[0].GetUserName()); err != nil {
		return err
	}
	if err := d.Set("user_email", users[0].GetUserEmail()); err != nil {
		return err
	}
	if err := d.Set("user_id", users[0].GetUserId()); err != nil {
		return err
	}
	if err := d.Set("path", users[0].GetPath()); err != nil {
		return err
	}
	if err := d.Set("creation_date", users[0].GetCreationDate()); err != nil {
		return err
	}
	if err := d.Set("last_modification_date", users[0].GetLastModificationDate()); err != nil {
		return err
	}
	return nil
}

func buildUsersFilters(set *schema.Set) (*oscgo.FiltersUsers, error) {
	var filters oscgo.FiltersUsers
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "user_ids":
			filters.SetUserIds(filterValues)
		default:
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}
