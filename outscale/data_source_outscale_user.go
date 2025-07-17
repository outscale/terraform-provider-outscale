package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
	conn := meta.(*OutscaleClient).OSCAPI
	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return fmt.Errorf("filters: user_ids must be assigned")
	}
	req := oscgo.NewReadUsersRequest()
	filterReq := buildUsersFilters(filters.(*schema.Set))
	req.SetFilters(*filterReq)
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
	d.SetId(id.UniqueId())
	if len(users) == 0 {
		return fmt.Errorf("Unable to find user")
	}
	if len(users) > 1 {
		return fmt.Errorf("Find To many users")
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

func buildUsersFilters(set *schema.Set) *oscgo.FiltersUsers {
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
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return &filters
}
