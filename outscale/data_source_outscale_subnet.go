package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		r, _, err := conn.SubnetApi.ReadSubnets(context.Background(), &oscgo.ReadSubnetsOpts{ReadSubnetsRequest: optional.NewInterface(req)})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		resp = r
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
	if err := d.Set("tags", tagsOSCAPIToMap(subnet.GetTags())); err != nil {
		return err
	}
	if err := d.Set("available_ips_count", subnet.GetAvailableIpsCount()); err != nil {
		return err
	}
	if err := d.Set("request_id", resp.ResponseContext.GetRequestId()); err != nil {
		return err
	}

	return nil
}

func buildOutscaleOAPISubnetDataSourceFilters(set *schema.Set) *oscgo.FiltersSubnet {
	var filters oscgo.FiltersSubnet
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		// case "available_ips_counts":
		// 	filters.AvailableIpsCounts = filterValues
		case "ip_ranges":
			filters.IpRanges = &filterValues
		case "net_ids":
			filters.NetIds = &filterValues
		case "states":
			filters.States = &filterValues
		case "subnet_ids":
			filters.SubnetIds = &filterValues
		case "subregion_names":
			filters.SubregionNames = &filterValues
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return &filters
}
