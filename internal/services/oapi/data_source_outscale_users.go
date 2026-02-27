package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceUsersRead,
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

func DataSourceUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")

	var err error
	req := osc.ReadUsersRequest{}
	if filtersOk {
		req.Filters, err = buildUsersFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadUsers(ctx, req, options.WithRetryTimeout(2*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}
	users := ptr.From(resp.Users)
	d.SetId(id.UniqueId())
	if len(users) == 0 {
		return diag.FromErr(ErrNoResults)
	}
	d.SetId(id.UniqueId())
	usersToSet := make([]map[string]interface{}, len(users))
	for i, v := range users {
		user := make(map[string]interface{})

		user["user_id"] = v.UserId
		user["user_name"] = v.UserName
		user["user_email"] = v.UserEmail
		user["path"] = v.Path
		user["creation_date"] = from.ISO8601(v.CreationDate)
		user["last_modification_date"] = from.ISO8601(v.LastModificationDate)
		usersToSet[i] = user
	}
	return diag.FromErr(d.Set("users", usersToSet))
}
