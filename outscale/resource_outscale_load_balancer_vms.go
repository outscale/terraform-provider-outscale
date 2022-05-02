package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOutscaleOAPILBUAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILBUAttachmentCreate,
		Read:   resourceOutscaleOAPILBUAttachmentRead,
		Delete: resourceOutscaleOAPILBUAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"load_balancer_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"backend_vm_ids": {
				Type:     schema.TypeList,
				ForceNew: true,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPILBUAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	e, eok := d.GetOk("load_balancer_name")
	i, iok := d.GetOk("backend_vm_ids")

	if !eok && !iok {
		return fmt.Errorf("please provide the required attributes load_balancer_name and backend_vm_id")
	}

	m := i.([]interface{})
	a := make([]string, len(m))
	for k, v := range m {
		a[k] = v.(string)
	}

	req := oscgo.RegisterVmsInLoadBalancerRequest{
		LoadBalancerName: e.(string),
		BackendVmIds:     a,
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.LoadBalancerApi.
			RegisterVmsInLoadBalancer(context.Background()).
			RegisterVmsInLoadBalancerRequest(req).
			Execute()

		if err != nil {
			return utils.CheckThrottling(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failure registering backend_vm_ids with LBU: %s", err)
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-", e)))

	return resourceOutscaleOAPILBUAttachmentRead(d, meta)
}

func resourceOutscaleOAPILBUAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	found := false
	e := d.Get("load_balancer_name").(string)
	lb, _, err := readResourceLb(conn, e)
	expected := d.Get("backend_vm_ids").([]interface{})

	if err != nil {
		return err
	}
	for _, v := range *lb.BackendVmIds {
		for k1 := range expected {
			sid := expected[k1].(string)
			if sid == v {
				d.Set("backend_vm_ids", expected)
				found = true
			}
		}
	}

	if !found {
		log.Printf("[WARN] i %s not found in lbu attachments", expected)
		d.SetId("")
	}

	return nil
}

func resourceOutscaleOAPILBUAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	e := d.Get("load_balancer_name").(string)
	i := d.Get("backend_vm_ids").([]interface{})

	lb := make([]string, len(i))

	for k, v := range i {
		lb[k] = v.(string)
	}

	req := oscgo.DeregisterVmsInLoadBalancerRequest{
		LoadBalancerName: e,
		BackendVmIds:     lb,
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err := conn.LoadBalancerApi.
			DeregisterVmsInLoadBalancer(context.Background()).
			DeregisterVmsInLoadBalancerRequest(req).
			Execute()

		if err != nil {
			return utils.CheckThrottling(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failure deregistering backend_vm_ids from LBU: %s", err)
	}

	return nil
}
