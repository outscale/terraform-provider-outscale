package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOutscaleOAPILBUAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILBUAttachmentCreate,
		Read:   resourceOutscaleOAPILBUAttachmentRead,
		Update: resourceOutscaleOAPILBUAttachmentUpdate,
		Delete: resourceOutscaleOAPILBUAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"load_balancer_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"backend_vm_ids": {
				Type:     schema.TypeSet,
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

	a := make([]string, i.(*schema.Set).Len())
	for k, v := range i.(*schema.Set).List() {
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
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
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
	e := d.Get("load_balancer_name").(string)
	lb, _, err := readResourceLb(conn, e)
	expected := d.Get("backend_vm_ids").(*schema.Set)

	if err != nil {
		return err
	}

	all_backends := schema.Set{F: expected.F}
	for _, v := range *lb.BackendVmIds {
		all_backends.Add(v)
	}

	managed := all_backends.Intersection(expected)
	d.Set("backend_vm_ids", managed)

	if managed.Len() == 0 {
		log.Printf("[WARN] not expected attachments found in LBU %e", e)
		log.Printf("[WARN] lbu current attachments are %#v", all_backends)
		log.Printf("[WARN] we would manage only these attachments %#v", expected)
		log.Printf("[WARN] no managed attachments are present.")
		d.SetId("")
	}

	return nil
}

func resourceOutscaleOAPILBUDiffBackendVmIds(oldBackends *schema.Set, newBackends *schema.Set) (*schema.Set, *schema.Set) {

	// Strange, but if you insist...
	if newBackends == nil {
		if oldBackends != nil {
			return nil, oldBackends
		}
		return nil, nil
	}
	// Start by supposing that we create everything and remove nothing.
	create := schema.CopySet(newBackends)
	remove := schema.NewSet(create.F, []interface{}{})

	for _, backend := range oldBackends.List() {
		// When old set contains backends not in the new set,
		// they are to be removed.
		if !create.Contains(backend) {
			remove.Add(backend)
		} else {
			create.Remove(backend)
		}
	}

	return create, remove
}

func resourceOutscaleOAPILBUAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	lbu_name := d.Get("load_balancer_name").(string)
	var err error

	if !d.HasChange("backend_vm_ids") {
		return nil
	}

	oldBackends, newBackends := d.GetChange("backend_vm_ids")
	create, remove := resourceOutscaleOAPILBUDiffBackendVmIds(oldBackends.(*schema.Set), newBackends.(*schema.Set))

	if create != nil && create.Len() > 0 {
		// Convert the Set to a string list
		createStrings := make([]string, 0, create.Len())
		for _, val := range create.List() {
			createStrings = append(createStrings, val.(string))
		}
		// Make the Register request
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, _, err = conn.LoadBalancerApi.
				RegisterVmsInLoadBalancer(context.Background()).
				RegisterVmsInLoadBalancerRequest(
					oscgo.RegisterVmsInLoadBalancerRequest{
						LoadBalancerName: lbu_name,
						BackendVmIds:     createStrings,
					}).
				Execute()
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
			return fmt.Errorf("Failure registering new backend_vm_ids with LBU: %s", err)
		}
	}
	if remove != nil && remove.Len() > 0 {
		// Convert the Set to a string list
		removeStrings := make([]string, 0, remove.Len())
		for _, val := range remove.List() {
			removeStrings = append(removeStrings, val.(string))
		}

		// Make the Deregister request
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, _, err := conn.LoadBalancerApi.
				DeregisterVmsInLoadBalancer(context.Background()).
				DeregisterVmsInLoadBalancerRequest(
					oscgo.DeregisterVmsInLoadBalancerRequest{
						LoadBalancerName: lbu_name,
						BackendVmIds:     removeStrings,
					}).
				Execute()
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
			return fmt.Errorf("Failure deregistering old backend_vm_ids from LBU: %s", err)
		}
	}
	return resourceOutscaleOAPILBUAttachmentRead(d, meta)
}

func resourceOutscaleOAPILBUAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	e := d.Get("load_balancer_name").(string)
	i := d.Get("backend_vm_ids").(*schema.Set)

	lb := make([]string, i.Len())

	for k, v := range i.List() {
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
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failure deregistering backend_vm_ids from LBU: %s", err)
	}

	return nil
}
