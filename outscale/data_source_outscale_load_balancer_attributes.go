package outscale

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceOutscaleOAPILoadBalancerAttr() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPILoadBalancerAttrRead,
		Schema: getDataSourceSchemas(lbAttrAttributes()),
	}
}

func lbAttrAttributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"load_balancer_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"load_balancer_attributes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"access_log": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"publication_interval": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"is_enabled": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"osu_bucket_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"osu_bucket_prefix": {
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

func dataSourceOutscaleOAPILoadBalancerAttrRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	resp, elbName, err := readLbs(conn, d)
	if err != nil {
		return err
	}

	lbs := *resp.LoadBalancers
	if len(lbs) != 1 {
		return fmt.Errorf("Unable to find LBU: %s", *elbName)
	}

	lb := (lbs)[0]

	a := lb.AccessLog

	ld := make([]map[string]interface{}, 1)
	acc := make(map[string]interface{})

	acc["publication_interval"] = a.PublicationInterval
	acc["is_enabled"] = a.IsEnabled
	acc["osu_bucket_name"] = a.OsuBucketName
	acc["osu_bucket_prefix"] = a.OsuBucketPrefix

	ld[0] = map[string]interface{}{"access_log": acc}

	d.Set("request_id", resp.ResponseContext.RequestId)
	d.SetId(resource.UniqueId())

	return d.Set("load_balancer_attributes", ld)
}
