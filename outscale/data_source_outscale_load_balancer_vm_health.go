package outscale

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOutscaleLoadBalancerVmsHeals() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleLoadBalancerVmsHealRead,
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

func dataSourceOutscaleLoadBalancerVmsHealRead(d *schema.ResourceData,
	meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	ename, ok := d.GetOk("load_balancer_name")
	if ok == false {
		return errors.New("load_balancer_name is require")
	}

	req := oscgo.ReadVmsHealthRequest{
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

	var resp oscgo.ReadVmsHealthResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.LoadBalancerApi.ReadVmsHealth(
			context.Background()).ReadVmsHealthRequest(req).
			Execute()

		if err != nil {
			log.Printf("[DEBUG] err: (%s)", err)
			if strings.Contains(fmt.Sprint(err), "InvalidResource") ||
				strings.Contains(fmt.Sprint(err), "Bad Request") {
				return resource.RetryableError(err)
			}
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balacer Vms Heal: %s", err)
	}

	if resp.BackendVmHealth == nil {
		return fmt.Errorf("lb.BackendVmHealth not found")
	}
	lbvh := make([]map[string]interface{}, len(*resp.BackendVmHealth))
	for k, v := range *resp.BackendVmHealth {
		a := make(map[string]interface{})
		a["description"] = v.Description
		a["state"] = v.State
		a["state_reason"] = v.StateReason
		a["vm_id"] = v.VmId
		lbvh[k] = a
	}
	d.Set("backend_vm_health", lbvh)
	//  ename.(string) "-heal-" resource.UniqueId()
	id := ename.(string) + "-heal-"
	d.SetId(resource.PrefixedUniqueId(id))
	return nil
}
