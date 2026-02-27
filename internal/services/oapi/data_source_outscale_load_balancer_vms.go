package oapi

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscaleLoadBalancerVms() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleLoadBalancerVmsRead,
		Schema: getDataSourceSchemas(map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"load_balancer_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backend_vm_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func DataSourceOutscaleLoadBalancerVmsRead(ctx context.Context, d *schema.ResourceData,
	meta interface{},
) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	lb, _, err := readLbs0(ctx, client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("backend_vm_ids", utils.StringSlicePtrToInterfaceSlice(&lb.BackendVmIds))
	return nil
}
