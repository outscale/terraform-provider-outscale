package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antihax/optional"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	oscgo "github.com/marinsalinas/osc-sdk-go"
)

func resourceOutscaleOAPILoadBalancerSSLCertificate() *schema.Resource {
	return &schema.Resource{
		Read:   resourceOutscaleOAPILoadBalancerSSLCertificateRead,
		Create: resourceOutscaleOAPILoadBalancerSSLCertificateCreate,
		Update: resourceOutscaleOAPILoadBalancerSSLCertificateUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Delete: resourceOutscaleOAPILoadBalancerSSLCertificateDelete,

		Schema: map[string]*schema.Schema{
			"load_balancer_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"load_balancer_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"server_certificate_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPILoadBalancerSSLCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	ename, ok := d.GetOk("load_balancer_name")
	port, pok := d.GetOk("load_balancer_port")
	ssl, sok := d.GetOk("server_certificate_id")

	if !ok {
		return fmt.Errorf("please provide the name of the load balancer")
	}

	if !pok {
		return fmt.Errorf("please provide the load_balancer_port argument")
	}

	if !sok {
		return fmt.Errorf("please provide server_certificate_id argument")
	}
	port_i := port.(int)
	port_i64 := int64(port_i)
	ssl_s := ssl.(string)
	req := oscgo.UpdateLoadBalancerRequest{
		LoadBalancerName:    ename.(string),
		LoadBalancerPort:    &port_i64,
		ServerCertificateId: &ssl_s,
	}

	opts := oscgo.UpdateLoadBalancerOpts{
		optional.NewInterface(req),
	}

	var err error
	var resp = oscgo.UpdateLoadBalancerResponse{}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.LoadBalancerApi.UpdateLoadBalancer(
			context.Background(), &opts)

		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failure setting Load Balancer Listeners SSL Certificate for LBU: %s", err)
	}

	if resp.ResponseContext != nil {
		d.Set("request_id", resp.ResponseContext.RequestId)
	}

	d.SetId(ename.(string))

	return resourceOutscaleOAPILoadBalancerSSLCertificateRead(d, meta)
}

func resourceOutscaleOAPILoadBalancerSSLCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("server_certificate_id") {
		return nil
	}

	conn := meta.(*OutscaleClient).OSCAPI

	port := d.Get("load_balancer_port").(int64)
	ssl := d.Get("server_certificate_id").(string)

	req := oscgo.UpdateLoadBalancerRequest{
		LoadBalancerName:    d.Get("load_balancer_name").(string),
		LoadBalancerPort:    &port,
		ServerCertificateId: &ssl,
	}

	opts := oscgo.UpdateLoadBalancerOpts{
		optional.NewInterface(req),
	}

	var err error
	var resp = oscgo.UpdateLoadBalancerResponse{}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.LoadBalancerApi.UpdateLoadBalancer(
			context.Background(), &opts)

		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failure setting Load Balancer Listeners SSL Certificate for LBU: %s", err)
	}

	if resp.ResponseContext != nil {
		d.Set("request_id", resp.ResponseContext.RequestId)
	}

	d.SetId(req.LoadBalancerName)

	return resourceOutscaleOAPILoadBalancerSSLCertificateRead(d, meta)
}

func resourceOutscaleOAPILoadBalancerSSLCertificateRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOutscaleOAPILoadBalancerSSLCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}
