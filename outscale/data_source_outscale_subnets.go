package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOutscaleOAPISubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPISubnetsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"subnet_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"available_ips_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"map_public_ip_on_launch": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPISubnetsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadSubnetsRequest{}

	if id := d.Get("subnet_ids"); id != "" {
		var ids []string
		for _, v := range id.([]interface{}) {
			ids = append(ids, v.(string))
		}
		req.SetFilters(oscgo.FiltersSubnet{SubnetIds: &ids})
	}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters = buildOutscaleOAPISubnetDataSourceFilters(filters.(*schema.Set))
	}

	var resp oscgo.ReadSubnetsResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.SubnetApi.ReadSubnets(context.Background()).ReadSubnetsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	var errString string

	if err != nil {
		errString = err.Error()

		return fmt.Errorf("[DEBUG] Error reading Subnet (%s)", errString)
	}

	if len(resp.GetSubnets()) == 0 {
		return fmt.Errorf("no matching subnet found")
	}

	subnets := make([]map[string]interface{}, len(resp.GetSubnets()))

	for k, v := range resp.GetSubnets() {
		subnet := make(map[string]interface{})

		if v.GetSubregionName() != "" {
			subnet["subregion_name"] = v.GetSubregionName()
		}
		//if v.AvailableIpsCount != 0 {
		subnet["available_ips_count"] = v.GetAvailableIpsCount()
		//}
		if v.GetIpRange() != "" {
			subnet["ip_range"] = v.GetIpRange()
		}
		if v.GetState() != "" {
			subnet["state"] = v.GetState()
		}
		if v.GetSubnetId() != "" {
			subnet["subnet_id"] = v.GetSubnetId()
		}
		if v.GetTags() != nil {
			subnet["tags"] = tagsOSCAPIToMap(v.GetTags())
		}
		if v.GetNetId() != "" {
			subnet["net_id"] = v.GetNetId()
		}
		subnet["map_public_ip_on_launch"] = v.GetMapPublicIpOnLaunch()

		subnets[k] = subnet
	}

	if err := d.Set("subnets", subnets); err != nil {
		return err
	}
	d.SetId(resource.UniqueId())

	return nil
}
