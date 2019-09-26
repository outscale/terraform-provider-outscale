package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
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

						"tags": tagsOAPIListSchemaComputed(),

						"net_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"available_ips_count": {
							Type:     schema.TypeInt,
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
	conn := meta.(*OutscaleClient).OAPI

	req := &oapi.ReadSubnetsRequest{}

	if id := d.Get("subnet_ids"); id != "" {
		var ids []string
		for _, v := range id.([]interface{}) {
			ids = append(ids, v.(string))
		}
		req.Filters = oapi.FiltersSubnet{SubnetIds: ids}
	}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters = buildOutscaleOAPISubnetDataSourceFilters(filters.(*schema.Set))
	}

	var resp *oapi.POST_ReadSubnetsResponses
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.POST_ReadSubnets(*req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var errString string

	if err != nil || resp.OK == nil {
		if err != nil {
			errString = err.Error()
		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
		}

		return fmt.Errorf("[DEBUG] Error reading Subnet (%s)", errString)
	}

	response := resp.OK

	if response == nil || len(response.Subnets) == 0 {
		return fmt.Errorf("no matching subnet found")
	}

	subnets := make([]map[string]interface{}, len(response.Subnets))

	for k, v := range response.Subnets {
		subnet := make(map[string]interface{})

		if v.SubregionName != "" {
			subnet["subregion_name"] = v.SubregionName
		}
		//if v.AvailableIpsCount != 0 {
		subnet["available_ips_count"] = v.AvailableIpsCount
		//}
		if v.IpRange != "" {
			subnet["ip_range"] = v.IpRange
		}
		if v.State != "" {
			subnet["state"] = v.State
		}
		if v.SubnetId != "" {
			subnet["subnet_id"] = v.SubnetId
		}
		if v.Tags != nil {
			subnet["tags"] = tagsOAPIToMap(v.Tags)
		}
		if v.NetId != "" {
			subnet["net_id"] = v.NetId
		}

		subnets[k] = subnet
	}

	if err := d.Set("subnets", subnets); err != nil {
		return err
	}
	d.SetId(resource.UniqueId())
	d.Set("request_id", response.ResponseContext.RequestId)

	return nil
}
