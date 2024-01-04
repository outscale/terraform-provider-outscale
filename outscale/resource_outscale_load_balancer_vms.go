package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLBUAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceLBUAttachmentCreate,
		Read:   resourceLBUAttachmentRead,
		Update: resourceLBUAttachmentUpdate,
		Delete: resourceLBUAttachmentDelete,

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
func resourceLBUAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	lbuName := d.Get("load_balancer_name").(string)
	vmIds := d.Get("backend_vm_ids").(*schema.Set)

	req := oscgo.RegisterVmsInLoadBalancerRequest{
		LoadBalancerName: lbuName,
		BackendVmIds:     SetStringToListString(vmIds),
	}
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.LoadBalancerApi.
			RegisterVmsInLoadBalancer(context.Background()).
			RegisterVmsInLoadBalancerRequest(req).
			Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Failure registering backend_vm_ids with LBU: %s", err)
	}
	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-", lbuName)))
	return resourceLBUAttachmentRead(d, meta)
}

func resourceLBUAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	lbuName := d.Get("load_balancer_name").(string)
	lb, _, err := readResourceLb(conn, lbuName)
	if err != nil {
		return err
	}
	if lb == nil {
		utils.LogManuallyDeleted("LoadBalancerVms", d.Id())
		d.SetId("")
		return nil
	}

	expected := d.Get("backend_vm_ids").(*schema.Set)
	all_backends := schema.Set{F: expected.F}
	for _, v := range *lb.BackendVmIds {
		all_backends.Add(v)
	}

	managed := all_backends.Intersection(expected)
	d.Set("backend_vm_ids", managed)

	if managed.Len() == 0 {
		log.Printf("[WARN] not expected attachments found in LBU %s", lbuName)
		log.Printf("[WARN] lbu current attachments are %#v", all_backends)
		log.Printf("[WARN] we would manage only these attachments %#v", expected)
		log.Printf("[WARN] no managed attachments are present.")
		d.SetId("")
	}

	return nil
}

func resourceLBUAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	lbuName := d.Get("load_balancer_name").(string)
	var err error

	if !d.HasChange("backend_vm_ids") {
		return nil
	}

	oldBackends, newBackends := d.GetChange("backend_vm_ids")
	inter := oldBackends.(*schema.Set).Intersection(newBackends.(*schema.Set))
	created := newBackends.(*schema.Set).Difference(inter)
	removed := oldBackends.(*schema.Set).Difference(inter)

	if created.Len() > 0 {
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, httpResp, err := conn.LoadBalancerApi.
				RegisterVmsInLoadBalancer(context.Background()).
				RegisterVmsInLoadBalancerRequest(
					oscgo.RegisterVmsInLoadBalancerRequest{
						LoadBalancerName: lbuName,
						BackendVmIds:     SetStringToListString(created),
					}).
				Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Failure registering new backend_vm_ids with LBU: %s", err)
		}
	}
	if removed.Len() > 0 {
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, httpResp, err := conn.LoadBalancerApi.
				DeregisterVmsInLoadBalancer(context.Background()).
				DeregisterVmsInLoadBalancerRequest(
					oscgo.DeregisterVmsInLoadBalancerRequest{
						LoadBalancerName: lbuName,
						BackendVmIds:     SetStringToListString(removed),
					}).
				Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Failure deregistering old backend_vm_ids from LBU: %s", err)
		}
	}
	return resourceLBUAttachmentRead(d, meta)
}

func resourceLBUAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	lbuName := d.Get("load_balancer_name").(string)
	vmIds := d.Get("backend_vm_ids").(*schema.Set)

	req := oscgo.DeregisterVmsInLoadBalancerRequest{
		LoadBalancerName: lbuName,
		BackendVmIds:     SetStringToListString(vmIds),
	}
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.LoadBalancerApi.
			DeregisterVmsInLoadBalancer(context.Background()).
			DeregisterVmsInLoadBalancerRequest(req).
			Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Failure deregistering backend_vm_ids from LBU: %s", err)
	}
	return nil
}

func SetStringToListString(set *schema.Set) []string {
	result := make([]string, 0, set.Len())
	for _, val := range set.List() {
		result = append(result, val.(string))
	}
	return result
}
