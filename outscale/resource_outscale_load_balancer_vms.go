package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

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

			"backend_vm_id": {
				Type:     schema.TypeList,
				ForceNew: true,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vm_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceOutscaleOAPILBUAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	e, eok := d.GetOk("load_balancer_name")
	i, iok := d.GetOk("backend_vm_id")

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

	registerInstancesOpts := oscgo.RegisterVmsInLoadBalancerOpts{
		optional.NewInterface(req),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.LoadBalancerApi.
			RegisterVmsInLoadBalancer(context.Background(),
				&registerInstancesOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failure registering backend_vm_id with LBU: %s", err)
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-", e)))

	return nil
}

func resourceOutscaleOAPILBUAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	e := d.Get("load_balancer_name").(string)
	expected := d.Get("backend_vm_id").([]interface{})

	filter := &oscgo.FiltersLoadBalancer{
		LoadBalancerNames: &[]string{e},
	}

	req := oscgo.ReadLoadBalancersRequest{
		Filters: filter,
	}

	describeElbOpts := &oscgo.ReadLoadBalancersOpts{
		ReadLoadBalancersRequest: optional.NewInterface(req),
	}

	var resp oscgo.ReadLoadBalancersResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.LoadBalancerApi.ReadLoadBalancers(
			context.Background(),
			describeElbOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		/*
			if isLoadBalancerNotFound(err) {
				log.Printf("[ERROR] LBU %s not found", e)
				d.SetId("")
				return nil
			}
		*/
		return fmt.Errorf("Error retrieving LBU: %s", err)
	}

	found := false
	lbs := *resp.LoadBalancers
	if len(lbs) != 1 {
		return fmt.Errorf("Unable to find LBU: %s", expected)
	}

	lb := (lbs)[0]

	for _, v := range *lb.BackendVmIds {
		for k1 := range expected {
			sid := expected[k1].(string)
			if sid == v {
				d.Set("backend_vm_id", expected)
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
	i := d.Get("backend_vm_id").([]interface{})

	lb := make([]string, len(i))

	for k, v := range i {
		lb[k] = v.(string)
	}

	req := oscgo.DeregisterVmsInLoadBalancerRequest{
		LoadBalancerName: e,
		BackendVmIds:     lb,
	}
	deRegisterInstancesOpts := oscgo.DeregisterVmsInLoadBalancerOpts{
		optional.NewInterface(req),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err := conn.LoadBalancerApi.
			DeregisterVmsInLoadBalancer(context.Background(),
				&deRegisterInstancesOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failure deregistering backend_vm_id from LBU: %s", err)
	}

	return nil
}
