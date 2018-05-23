package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleOAPILBUAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILBUAttachmentCreate,
		Read:   resourceOutscaleOAPILBUAttachmentRead,
		Delete: resourceOutscaleOAPILBUAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"load_balancer_name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"backend_vm_id": &schema.Schema{
				Type:     schema.TypeList,
				ForceNew: true,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceOutscaleOAPILBUAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	e, eok := d.GetOk("load_balancer_name")
	i, iok := d.GetOk("backend_vm_id")

	if !eok && !iok {
		return fmt.Errorf("please provide the required attributes load_balancer_name and backend_vm_id")
	}

	lb := make([]*lbu.Instance, len(i.([]interface{})))

	for k, v := range i.([]interface{}) {
		lb[k] = &lbu.Instance{InstanceId: aws.String(v.(string))}
	}

	registerInstancesOpts := lbu.RegisterInstancesWithLoadBalancerInput{
		LoadBalancerName: aws.String(e.(string)),
		Instances:        lb,
	}

	log.Printf("[INFO] registering i %s with ELB %s", i, e)

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.API.RegisterInstancesWithLoadBalancer(&registerInstancesOpts)

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
		return fmt.Errorf("Failure registering backend_vm_id with ELB: %s", err)
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-", e)))

	return nil
}

func resourceOutscaleOAPILBUAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU
	e := d.Get("load_balancer_name").(string)

	// only add the backend_vm_id that was previously defined for this resource
	expected := d.Get("backend_vm_id").([]interface{})

	// Retrieve the ELB properties to get a list of attachments
	describeElbOpts := &lbu.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{aws.String(e)},
	}

	var resp *lbu.DescribeLoadBalancersOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.DescribeLoadBalancers(describeElbOpts)

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
		if isLoadBalancerNotFound(err) {
			log.Printf("[ERROR] ELB %s not found", e)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving ELB: %s", err)
	}
	if len(resp.LoadBalancerDescriptions) != 1 {
		log.Printf("[ERROR] Unable to find ELB: %v", resp.LoadBalancerDescriptions)
		d.SetId("")
		return nil
	}

	found := false
	for k, i := range resp.LoadBalancerDescriptions[0].Instances {
		instance := expected[k].(map[string]interface{})["instance_id"].(string)

		if instance == *i.InstanceId {
			d.Set("backend_vm_id", expected)
			found = true
		}
	}

	if !found {
		log.Printf("[WARN] i %s not found in lbu attachments", expected)
		d.SetId("")
	}

	return nil
}

func resourceOutscaleOAPILBUAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU
	e := d.Get("load_balancer_name").(string)
	i := d.Get("backend_vm_id").([]interface{})

	lb := make([]*lbu.Instance, len(i))

	for k, v := range i {
		lb[k] = &lbu.Instance{InstanceId: aws.String(v.(string))}
	}

	deRegisterInstancesOpts := lbu.DeregisterInstancesFromLoadBalancerInput{
		LoadBalancerName: aws.String(e),
		Instances:        lb,
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.API.DeregisterInstancesFromLoadBalancer(&deRegisterInstancesOpts)

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
		return fmt.Errorf("Failure deregistering backend_vm_id from ELB: %s", err)
	}

	return nil
}
