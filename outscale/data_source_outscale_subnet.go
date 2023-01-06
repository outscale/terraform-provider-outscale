package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cast"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleOAPISubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPISubnetRead,

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
			"tags": dataSourceTagsSchema(),
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

func dataSourceOutscaleOAPISubnetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadSubnetsRequest{}

	if id := d.Get("subnet_id"); id != "" {
		req.Filters = &oscgo.FiltersSubnet{SubnetIds: &[]string{id.(string)}}
	}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters = buildOutscaleOAPISubnetDataSourceFilters(filters.(*schema.Set))
	}

	var resp oscgo.ReadSubnetsResponse

	err := resource.Retry(120*time.Second, func() *resource.RetryError {
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
	if err := d.Set("tags", tagsOSCAPIToMap(subnet.GetTags())); err != nil {
		return err
	}
	if err := d.Set("available_ips_count", subnet.GetAvailableIpsCount()); err != nil {
		return err
	}

	return nil
}

func buildOutscaleOAPISubnetDataSourceFilters(set *schema.Set) *oscgo.FiltersSubnet {
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
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return &filters
}
