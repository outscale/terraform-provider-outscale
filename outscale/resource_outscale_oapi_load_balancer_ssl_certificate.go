package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
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
	conn := meta.(*OutscaleClient).LBU

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

	opts := lbu.SetLoadBalancerListenerSSLCertificateInput{
		LoadBalancerName: aws.String(ename.(string)),
		LoadBalancerPort: aws.Int64(int64(port.(int))),
		SSLCertificateId: aws.String(ssl.(string)),
	}
	var err error
	var resp = &lbu.SetLoadBalancerListenerSSLCertificateOutput{}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.SetLoadBalancerListenerSSLCertificate(&opts)

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

	if resp.ResponseMetadata != nil {
		d.Set("request_id", resp.ResponseMetadata.RequestId)
	}

	d.SetId(ename.(string))

	return resourceOutscaleOAPILoadBalancerSSLCertificateRead(d, meta)
}

func resourceOutscaleOAPILoadBalancerSSLCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("server_certificate_id") {
		return nil
	}

	conn := meta.(*OutscaleClient).LBU

	opts := lbu.SetLoadBalancerListenerSSLCertificateInput{
		LoadBalancerName: aws.String(d.Get("load_balancer_name").(string)),
		LoadBalancerPort: aws.Int64(d.Get("server_certificate_id").(int64)),
		SSLCertificateId: aws.String(d.Get("server_certificate_id").(string)),
	}
	var err error
	var resp = &lbu.SetLoadBalancerListenerSSLCertificateOutput{}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.SetLoadBalancerListenerSSLCertificate(&opts)

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

	if resp.ResponseMetadata != nil {
		d.Set("request_id", resp.ResponseMetadata.RequestId)
	}

	d.SetId(aws.StringValue(opts.LoadBalancerName))

	return resourceOutscaleOAPILoadBalancerSSLCertificateRead(d, meta)
}

func resourceOutscaleOAPILoadBalancerSSLCertificateRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOutscaleOAPILoadBalancerSSLCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}
