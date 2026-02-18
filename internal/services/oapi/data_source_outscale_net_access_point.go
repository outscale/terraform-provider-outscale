package oapi

import (
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/osc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
		Read: DataSourceOutscaleNetAccessPointRead,

		Schema: getDataSourceSchemas(napdSchema()),
	}
}

func DataSourceOutscaleNetAccessPointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return ErrFilterRequired
	}

	var resp osc.ReadNetAccessPointsResponse
	var err error
	req := osc.ReadNetAccessPointsRequest{}

	req.Filters, err = buildOutscaleDataSourcesNAPFilters(filters.(*schema.Set))
	if err != nil {
		return err
	}

	err = retry.Retry(30*time.Second, func() *retry.RetryError {
		rp, httpResp, err := client.NetAccessPointApi.ReadNetAccessPoints(
			ctx).
			ReadNetAccessPointsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	if len(resp.GetNetAccessPoints()) == 0 {
		return ErrNoResults
	}
	if len(resp.GetNetAccessPoints()) > 1 {
		return ErrMultipleResults
	}

	nap := resp.GetNetAccessPoints()[0]

	d.Set("net_access_point_id", nap.NetAccessPointId)
	d.Set("route_table_ids", utils.StringSlicePtrToInterfaceSlice(nap.RouteTableIds))
	d.Set("net_id", nap.NetId)
	d.Set("service_name", nap.ServiceName)
	d.Set("state", nap.State)
	d.Set("tags", FlattenOAPITagsSDK(nap.Tags))

	id := *nap.NetAccessPointId
	d.SetId(id)

	return nil
}
