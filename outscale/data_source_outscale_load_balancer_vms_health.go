package outscale

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

	describeElbOpts := &oscgo.ReadVmsHealthOpts{
		ReadVmsHealthRequest: optional.NewInterface(req),
	}

	var resp oscgo.ReadVmsHealthResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.LoadBalancerApi.ReadVmsHealth(
			context.Background(),
			describeElbOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if isLoadBalancerNotFound(err) {
			d.SetId("")
			return fmt.Errorf("Unknow error")
		}

		return fmt.Errorf("Error retrieving ELB: %s", err)
	}

	d.Set("request_id", resp.ResponseContext.RequestId)
	return nil
}
