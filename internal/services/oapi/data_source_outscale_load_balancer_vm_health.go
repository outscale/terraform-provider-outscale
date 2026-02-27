package oapi

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	sdkid "github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
)

func DataSourceOutscaleLoadBalancerVmsHeals() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleLoadBalancerVmsHealRead,
		Schema: getDataSourceSchemas(map[string]*schema.Schema{
			"load_balancer_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"backend_vm_ids": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"backend_vm_health": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_id": {
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
		}),
	}
}

func DataSourceOutscaleLoadBalancerVmsHealRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	ename, ok := d.GetOk("load_balancer_name")
	if !ok {
		return diag.Errorf("load_balancer_name is require")
	}

	req := osc.ReadVmsHealthRequest{
		LoadBalancerName: ename.(string),
	}

	vm_ids, ok := d.GetOk("backend_vm_ids")
	if ok {
		vm_ids_i := vm_ids.([]interface{})
		vm_ids_s := make([]string, 0, len(vm_ids_i))
		for _, v := range vm_ids_i {
			vm_ids_s = append(vm_ids_s, v.(string))
		}

		req.BackendVmIds = &vm_ids_s
	}

	resp, err := client.ReadVmsHealth(ctx, req, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.Errorf("error retrieving load balacer vms health: %s", err)
	}

	if resp.BackendVmHealth == nil {
		return diag.Errorf("lb.backendvmhealth not found")
	}
	lbvh := make([]map[string]interface{}, len(*resp.BackendVmHealth))
	for k, v := range *resp.BackendVmHealth {
		a := make(map[string]interface{})
		a["description"] = v.Description
		a["state"] = ptr.From(v.State)
		a["state_reason"] = v.StateReason
		a["vm_id"] = v.VmId
		lbvh[k] = a
	}
	d.Set("backend_vm_health", lbvh)
	id := ename.(string) + "-heal-"
	d.SetId(sdkid.PrefixedUniqueId(id))
	return nil
}
