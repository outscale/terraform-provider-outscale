package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func napdSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
		"net_access_point_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"net_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"service_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tags": TagsSchemaComputedSDK(),
		"route_table_ids": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func DataSourceOutscaleNetAccessPoint() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleNetAccessPointRead,

		Schema: getDataSourceSchemas(napdSchema()),
	}
}

func DataSourceOutscaleNetAccessPointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return diag.FromErr(ErrFilterRequired)
	}

	req := osc.ReadNetAccessPointsRequest{}

	var err error
	req.Filters, err = buildOutscaleDataSourcesNAPFilters(filters.(*schema.Set))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ReadNetAccessPoints(ctx, req, options.WithRetryTimeout(30*time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.NetAccessPoints == nil || len(*resp.NetAccessPoints) == 0 {
		return diag.FromErr(ErrNoResults)
	}
	if len(*resp.NetAccessPoints) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	nap := (*resp.NetAccessPoints)[0]

	d.Set("net_access_point_id", ptr.From(nap.NetAccessPointId))
	d.Set("route_table_ids", utils.StringSlicePtrToInterfaceSlice(nap.RouteTableIds))
	d.Set("net_id", ptr.From(nap.NetId))
	d.Set("service_name", ptr.From(nap.ServiceName))
	d.Set("state", ptr.From(nap.State))
	d.Set("tags", FlattenOAPITagsSDK(ptr.From(nap.Tags)))

	id := *nap.NetAccessPointId
	d.SetId(id)

	return nil
}
