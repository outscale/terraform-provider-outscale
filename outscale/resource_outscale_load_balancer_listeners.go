package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func resourceOutscaleOAPILoadBalancerListeners() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILoadBalancerListenersCreate,
		Read:   resourceOutscaleOAPILoadBalancerListenersRead,
		Delete: resourceOutscaleOAPILoadBalancerListenersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"listener": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backend_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"backend_protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"load_balancer_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"load_balancer_protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"server_certificate_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"load_balancer_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// "request_id": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
		},
	}
}

func resourceOutscaleOAPILoadBalancerListenersCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	elbOpts := &lbu.CreateLoadBalancerListenersInput{}

	listener, err := expandListeners(d.Get("listener").([]interface{}))
	if err != nil {
		return err
	}

	elbOpts.Listeners = listener

	if v, ok := d.GetOk("load_balancer_name"); ok {
		elbOpts.LoadBalancerName = aws.String(v.(string))
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.API.CreateLoadBalancerListeners(elbOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "DuplicateListener") {
				log.Printf("[DEBUG] Duplicate listener found for ELB (%s), retrying", d.Id())
				return resource.RetryableError(err)
			}
			if strings.Contains(fmt.Sprint(err), "CertificateNotFound") && strings.Contains(fmt.Sprint(err), "Server Certificate not found for the key: arn") {
				log.Printf("[DEBUG] SSL Cert not found for given ARN, retrying")
				return resource.RetryableError(err)
			}
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error creating ELB Listener with SSL Cert, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(*elbOpts.LoadBalancerName)
	log.Printf("[INFO] ELB ID: %s", d.Id())

	return resourceOutscaleOAPILoadBalancerListenersRead(d, meta)
}

func resourceOutscaleOAPILoadBalancerListenersRead(d *schema.ResourceData, meta interface{}) error {
	listener, err := expandListeners(d.Get("listener").([]interface{}))
	if err != nil {
		return err
	}

	result := make([]map[string]interface{}, 0, len(listener))
	for _, i := range listener {
		listener := map[string]interface{}{
			"backend_port":           aws.Int64Value(i.InstancePort),
			"backend_protocol":       aws.StringValue(i.InstanceProtocol),
			"load_balancer_port":     aws.Int64Value(i.LoadBalancerPort),
			"load_balancer_protocol": aws.StringValue(i.Protocol),
			"server_certificate_id":  aws.StringValue(i.SSLCertificateId),
		}
		result = append(result, listener)
	}
	if err := d.Set("listener", result); err != nil {
		return err
	}
	d.Set("load_balancer_name", d.Id())

	// d.Set("request_id", describeResp.ResponseMetadata.RequestID)

	return nil
}

func resourceOutscaleOAPILoadBalancerListenersUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	d.Partial(true)

	if d.HasChange("listener") {
		o, n := d.GetChange("listener")
		os := o.([]interface{})
		ns := n.([]interface{})

		remove, _ := expandListeners(ns)
		add, _ := expandListeners(os)

		if len(remove) > 0 {
			ports := make([]*int64, 0, len(remove))
			for _, listener := range remove {
				ports = append(ports, listener.LoadBalancerPort)
			}

			deleteListenersOpts := &lbu.DeleteLoadBalancerListenersInput{
				LoadBalancerName:  aws.String(d.Id()),
				LoadBalancerPorts: ports,
			}

			var err error
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err = conn.API.DeleteLoadBalancerListeners(deleteListenersOpts)

				if err != nil {
					if strings.Contains(err.Error(), "Throttling:") {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("Failure removing outdated ELB listener: %s", err)
			}
		}

		if len(add) > 0 {
			createListenersOpts := &lbu.CreateLoadBalancerListenersInput{
				LoadBalancerName: aws.String(d.Id()),
				Listeners:        add,
			}

			var err error
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err = conn.API.CreateLoadBalancerListeners(createListenersOpts)
				if err != nil {
					if err, ok := err.(awserr.Error); ok {
						if strings.Contains(fmt.Sprint(err), "DuplicateListener") {
							log.Printf("[DEBUG] Duplicate listener found for ELB (%s), retrying", d.Id())
							return resource.RetryableError(err)
						}
						if strings.Contains(fmt.Sprint(err), "CertificateNotFound") && strings.Contains(fmt.Sprint(err), "Server Certificate not found for the key: arn") {
							log.Printf("[DEBUG] SSL Cert not found for given ARN, retrying")
							return resource.RetryableError(err)
						}
						if strings.Contains(fmt.Sprint(err), "Throttling") && strings.Contains(fmt.Sprint(err), "Server Certificate not found for the key: arn") {
							log.Printf("[DEBUG] SSL Cert not found for given ARN, retrying")
							return resource.RetryableError(err)
						}
					}

					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("Failure adding new or updated ELB listener: %s", err)
			}
		}

		d.SetPartial("listener")
	}

	d.Partial(false)

	return resourceOutscaleOAPILoadBalancerListenersRead(d, meta)
}

func resourceOutscaleOAPILoadBalancerListenersDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	remove, _ := expandListeners(d.Get("listener").([]interface{}))

	ports := make([]*int64, 0, len(remove))
	for _, listener := range remove {
		ports = append(ports, listener.LoadBalancerPort)
	}

	deleteListenersOpts := &lbu.DeleteLoadBalancerListenersInput{
		LoadBalancerName:  aws.String(d.Id()),
		LoadBalancerPorts: ports,
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.API.DeleteLoadBalancerListeners(deleteListenersOpts)

		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failure removing outdated ELB listener: %s", err)
	}

	d.SetId("")

	return nil
}
