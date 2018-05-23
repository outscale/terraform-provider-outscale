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

func resourceOutscaleLBUAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleLBUAttachmentCreate,
		Read:   resourceOutscaleLBUAttachmentRead,
		Delete: resourceOutscaleLBUAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"load_balancer_name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"instances": &schema.Schema{
				Type:     schema.TypeList,
				ForceNew: true,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceOutscaleLBUAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	e, eok := d.GetOk("load_balancer_name")
	i, iok := d.GetOk("instances")

	if !eok && !iok {
		return fmt.Errorf("please provide the required attributes load_balancer_name and instances")
	}

	lb := make([]*lbu.Instance, len(i.([]interface{})))

	for k, v := range i.([]interface{}) {
		ins := v.(map[string]interface{})["instance_id"]
		lb[k] = &lbu.Instance{InstanceId: aws.String(ins.(string))}
	}

	registerInstancesOpts := lbu.RegisterInstancesWithLoadBalancerInput{
		LoadBalancerName: aws.String(e.(string)),
		Instances:        lb,
	}

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
		return fmt.Errorf("Failure registering instances with LBU: %s", err)
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-", e)))

	return nil
}

func resourceOutscaleLBUAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	e := d.Get("load_balancer_name").(string)
	expected := d.Get("instances").([]interface{})

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
			log.Printf("[ERROR] LBU %s not found", e)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving LBU: %s", err)
	}
	if len(resp.LoadBalancerDescriptions) != 1 {
		log.Printf("[ERROR] Unable to find LBU: %v", resp.LoadBalancerDescriptions)
		d.SetId("")
		return nil
	}

	found := false
	for k, i := range resp.LoadBalancerDescriptions[0].Instances {
		instance := expected[k].(map[string]interface{})["instance_id"].(string)

		if instance == *i.InstanceId {
			d.Set("instances", expected)
			found = true
		}
	}

	if !found {
		log.Printf("[WARN] i %s not found in lbu attachments", expected)
		d.SetId("")
	}

	return nil
}

func resourceOutscaleLBUAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU
	e := d.Get("load_balancer_name").(string)
	i := d.Get("instances").([]interface{})

	lb := make([]*lbu.Instance, len(i))

	for k, v := range i {
		ins := v.(map[string]interface{})["instance_id"]
		lb[k] = &lbu.Instance{InstanceId: aws.String(ins.(string))}
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
		return fmt.Errorf("Failure deregistering instances from LBU: %s", err)
	}

	return nil
}
