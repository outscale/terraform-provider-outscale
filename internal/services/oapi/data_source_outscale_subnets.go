package oapi

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleSubnets() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleSubnetsRead,

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

						"tags": TagsSchemaComputedSDK(),

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

func DataSourceOutscaleSubnetsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	req := oscgo.ReadSubnetsRequest{}

	if id := d.Get("subnet_ids"); id != "" {
		var ids []string
		for _, v := range id.([]interface{}) {
			ids = append(ids, v.(string))
		}
		req.SetFilters(oscgo.FiltersSubnet{SubnetIds: &ids})
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
			subnet["tags"] = FlattenOAPITagsSDK(v.GetTags())
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
	d.SetId(id.UniqueId())

	return nil
}
