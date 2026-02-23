package oapi

import (
	"context"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleLBUTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleLBUTagsRead,

		Schema: getDataSourceSchemas(getDSOAPILBUTagsSchema()),
	}
}

func DataSourceOutscaleLBUTagsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	ename, nameOk := d.GetOk("load_balancer_names")
	if !nameOk {
		return diag.Errorf("load_balancer_names is required")
	}

	names := ename.([]interface{})

	req := osc.ReadLoadBalancerTagsRequest{
		LoadBalancerNames: utils.InterfaceSliceToStringSlice(names),
	}

	resp, err := client.ReadLoadBalancerTags(ctx, req, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.Tags == nil {
		return diag.Errorf("tags of lbus (%v) not found", req.LoadBalancerNames)
	}

	tags := *resp.Tags
	l := len(*resp.Tags)

	ta := make([]map[string]interface{}, l)
	for k1, v1 := range tags {
		t := make(map[string]interface{})
		t["key"] = v1.Key
		t["value"] = v1.Value
		t["load_balancer_name"] = v1.LoadBalancerName
		ta[k1] = t
	}

	d.Set("tags", ta)
	d.SetId(id.UniqueId())
	return nil
}

func getDSOAPILBUTagsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"load_balancer_names": {
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
