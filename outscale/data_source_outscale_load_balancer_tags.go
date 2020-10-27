package outscale

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceOutscaleOAPILBUTags() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPILBUTagsRead,

		Schema: getDataSourceSchemas(getDSOAPILBUTagsSchema()),
	}
}

func dataSourceOutscaleOAPILBUTagsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	resp, _, err := readLbs(conn, d)
	if err != nil {
		return err
	}
	lbs := resp.LoadBalancers
	l := 0

	for _, v := range *lbs {
		l += len(*v.Tags)
	}

	ta := make([]map[string]interface{}, l)
	for _, v := range *lbs {
		for k1, v1 := range *v.Tags {
			t := make(map[string]interface{})
			t["key"] = v1.Key
			t["value"] = v1.Value
			t["load_balancer_name"] = v.LoadBalancerName
			ta[k1] = t
		}
	}

	d.Set("tags", ta)
	d.SetId(resource.UniqueId())
	d.Set("request_id", resp.ResponseContext.RequestId)
	return nil
}

func getDSOAPILBUTagsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"load_balancer_name": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"tags": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"load_balancer_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"value": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
