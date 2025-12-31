package oapi

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cast"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscaleSubnet() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleSubnetRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"subregion_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_range": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": TagsSchemaComputedSDK(),
			"net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_ips_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"map_public_ip_on_launch": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func DataSourceOutscaleSubnetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	req := oscgo.ReadSubnetsRequest{}

	if id := d.Get("subnet_id"); id != "" {
		req.Filters = &oscgo.FiltersSubnet{SubnetIds: &[]string{id.(string)}}
	}

	filters, filtersOk := d.GetOk("filter")

	var err error
	if filtersOk {
		req.Filters, err = buildOutscaleSubnetDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp oscgo.ReadSubnetsResponse
	err = retry.Retry(120*time.Second, func() *retry.RetryError {
		var err error
		rp, httpResp, err := conn.SubnetApi.ReadSubnets(context.Background()).ReadSubnetsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		errString := err.Error()
		return fmt.Errorf("[DEBUG] Error reading Subnet (%s)", errString)
	}

	if len(resp.GetSubnets()) == 0 {
		return fmt.Errorf("no matching subnet found")
	}

	if len(resp.GetSubnets()) > 1 {
		return fmt.Errorf("multiple subnets matched; use additional constraints to reduce matches to a single subnet")
	}

	subnet := resp.GetSubnets()[0]

	d.SetId(subnet.GetSubnetId())
	if err := d.Set("subnet_id", subnet.GetSubnetId()); err != nil {
		return err
	}
	if err := d.Set("net_id", subnet.GetNetId()); err != nil {
		return err
	}
	if err := d.Set("subregion_name", subnet.GetSubregionName()); err != nil {
		return err
	}
	if err := d.Set("ip_range", subnet.GetIpRange()); err != nil {
		return err
	}
	if err := d.Set("state", subnet.GetState()); err != nil {
		return err
	}
	if err := d.Set("map_public_ip_on_launch", subnet.GetMapPublicIpOnLaunch()); err != nil {
		return err
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(subnet.GetTags())); err != nil {
		return err
	}
	if err := d.Set("available_ips_count", subnet.GetAvailableIpsCount()); err != nil {
		return err
	}

	return nil
}

func buildOutscaleSubnetDataSourceFilters(set *schema.Set) (*oscgo.FiltersSubnet, error) {
	var filters oscgo.FiltersSubnet
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		var availableIPsCounts []int32
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
			availableIPsCounts = append(availableIPsCounts, cast.ToInt32(e))
		}

		switch name := m["name"].(string); name {
		case "available_ips_counts":
			filters.SetAvailableIpsCounts(utils.StringSliceToInt32Slice(filterValues))
		case "ip_ranges":
			filters.SetIpRanges(filterValues)
		case "net_ids":
			filters.SetNetIds(filterValues)
		case "states":
			filters.SetStates(filterValues)
		case "subnet_ids":
			filters.SetSubnetIds(filterValues)
		case "subregion_names":
			filters.SetSubregionNames(filterValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		case "tags":
			filters.SetTags(filterValues)

		default:
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}
