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

func resourceOutscaleLoadBalancerListeners() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleLoadBalancerListenersCreate,
		Read:   resourceOutscaleLoadBalancerListenersRead,
		Delete: resourceOutscaleLoadBalancerListenersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"listeners": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"instance_protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"load_balancer_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"ssl_certificate_id": &schema.Schema{
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
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleLoadBalancerListenersCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	elbOpts := &lbu.CreateLoadBalancerListenersInput{}

	listeners, err := expandListeners(d.Get("listeners").([]interface{}))
	if err != nil {
		return err
	}

	elbOpts.Listeners = listeners

	if v, ok := d.GetOk("load_balancer_name"); ok {
		elbOpts.LoadBalancerName = aws.String(v.(string))
	}

	resp := &lbu.CreateLoadBalancerListenersOutput{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.CreateLoadBalancerListeners(elbOpts)

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

	d.Set("load_balancer_name", d.Id())
	d.Set("request_id", *resp.ResponseMetadata.RequestID)

	return resourceOutscaleLoadBalancerListenersRead(d, meta)
}

func resourceOutscaleLoadBalancerListenersRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOutscaleLoadBalancerListenersUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	d.Partial(true)

	if d.HasChange("listeners") {
		o, n := d.GetChange("listeners")
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
				return fmt.Errorf("Failure removing outdated ELB listeners: %s", err)
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
				return fmt.Errorf("Failure adding new or updated ELB listeners: %s", err)
			}
		}

		d.SetPartial("listeners")
	}

	d.Partial(false)

	return resourceOutscaleLoadBalancerListenersRead(d, meta)
}

func resourceOutscaleLoadBalancerListenersDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	remove, _ := expandListeners(d.Get("listeners").([]interface{}))

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
		return fmt.Errorf("Failure removing outdated ELB listeners: %s", err)
	}

	d.SetId("")

	return nil
}
