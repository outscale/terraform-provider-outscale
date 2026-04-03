package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/samber/lo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceUserRead,
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

func DataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return diag.Errorf("filters: user_ids must be assigned")
	}

	var err error
	req := osc.ReadUsersRequest{}

	req.Filters, err = buildUsersFilters(filters.(*schema.Set))
	if err != nil {
		return diag.FromErr(err)
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
	if len(users) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}
	user := users[0]

	if err := d.Set("user_name", ptr.From(user.UserName)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("user_email", ptr.From(user.UserEmail)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("user_id", ptr.From(user.UserId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("path", ptr.From(user.Path)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("creation_date", from.ISO8601(user.CreationDate)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_modification_date", from.ISO8601(user.LastModificationDate)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func buildUsersFilters(set *schema.Set) (*osc.FiltersUsers, error) {
	var filters osc.FiltersUsers
	for _, v := range set.List() {
		m := v.(map[string]any)
		filterValues := lo.Map(m["values"].([]any), func(e any, _ int) string {
			return e.(string)
		})

		switch name := m["name"].(string); name {
		case "user_ids":
			filters.UserIds = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
