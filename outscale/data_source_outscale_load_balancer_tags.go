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

	v, resp, err := readLbs0(conn, d)
	if err != nil {
		return err
	}

	t := make(map[string]interface{})
	t["load_balancer_name"] = v.LoadBalancerName

	ta := make([]map[string]interface{}, len(*v.Tags))
	for k1, v1 := range *v.Tags {
		t := make(map[string]interface{})
		t["key"] = v1.Key
		t["value"] = v1.Key
		ta[k1] = t
	}

	t["tag"] = ta

	d.SetId(resource.UniqueId())
	d.Set("request_id", resp.ResponseContext.RequestId)

	return d.Set("tag", t)
}

func getDSOAPILBUTagsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"load_balancer_name": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"tag": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"load_balancer_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tag": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
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
				},
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
