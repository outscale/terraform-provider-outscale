package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
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
		"tags": dataSourceTagsSchema(),
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

func dataSourceOutscaleNetAccessPoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleNetAccessPointRead,

		Schema: getDataSourceSchemas(napdSchema()),
	}
}

func dataSourceOutscaleNetAccessPointRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return fmt.Errorf("One of filters must be assigned")
	}

	var resp oscgo.ReadNetAccessPointsResponse
	var err error
	req := oscgo.ReadNetAccessPointsRequest{}

	req.SetFilters(buildOutscaleDataSourcesNAPFilters(filters.(*schema.Set)))
	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.NetAccessPointApi.ReadNetAccessPoints(
			context.Background()).
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

	if err := utils.IsResponseEmptyOrMutiple(len(resp.GetNetAccessPoints()), "NetAccessPoint"); err != nil {
		return err
	}

	nap := resp.GetNetAccessPoints()[0]

	d.Set("net_access_point_id", nap.NetAccessPointId)
	d.Set("route_table_ids", utils.StringSlicePtrToInterfaceSlice(nap.RouteTableIds))
	d.Set("net_id", nap.NetId)
	d.Set("service_name", nap.ServiceName)
	d.Set("state", nap.State)
	d.Set("tags", tagsOSCAPIToMap(nap.GetTags()))

	id := *nap.NetAccessPointId
	d.SetId(id)

	return nil
}
