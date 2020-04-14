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

func resourceOutscaleOAPILoadBalancerListeners() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILoadBalancerListenersCreate,
		Read:   resourceOutscaleOAPILoadBalancerListenersRead,
		Delete: resourceOutscaleOAPILoadBalancerListenersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"listener": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backend_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"backend_protocol": {
							Type:     schema.TypeString,
							Required: true,
						},

						"load_balancer_port": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"load_balancer_protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"server_certificate_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"load_balancer_name": {
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

func expandListeners(configured []interface{}) ([]*oscgo.Listener, error) {
	listeners := make([]*oscgo.Listener, 0, len(configured))

	for _, lRaw := range configured {
		data := lRaw.(map[string]interface{})

		ip := int64(data["backend_port"].(int))
		lp := int64(data["load_balancer_port"].(int))
		bproto := data["backend_protocol"].(string)
		lproto := data["load_balancer_protocol"].(string)
		l := &oscgo.Listener{
			BackendPort:          &ip,
			BackendProtocol:      &bproto,
			LoadBalancerPort:     &lp,
			LoadBalancerProtocol: &lproto,
		}

		if v, ok := data["ssl_certificate_id"]; ok && v != "" {
			vs := v.(string)
			l.ServerCertificateId = &vs
		}

		var valid bool
		if l.ServerCertificateId != nil && *l.ServerCertificateId != "" {
			// validate the protocol is correct
			for _, p := range []string{"https", "ssl"} {
				if (strings.ToLower(*l.BackendProtocol) == p) ||
					(strings.ToLower(*l.LoadBalancerProtocol) == p) {
					valid = true
				}
			}
		} else {
			valid = true
		}

		if valid {
			listeners = append(listeners, l)
		} else {
			return nil, fmt.Errorf("[ERR] ELB Listener: ssl_certificate_id may be set only when protocol is 'https' or 'ssl'")
		}
	}

	return listeners, nil
}

func expandListenerForCreation(configured []interface{}) ([]oscgo.ListenerForCreation, error) {
	listeners := make([]oscgo.ListenerForCreation, 0, len(configured))

	for _, lRaw := range configured {
		data := lRaw.(map[string]interface{})

		ip := int64(data["backend_port"].(int))
		lp := int64(data["load_balancer_port"].(int))
		bproto := data["backend_protocol"].(string)
		lproto := data["load_balancer_protocol"].(string)
		l := oscgo.ListenerForCreation{
			BackendPort:          ip,
			BackendProtocol:      &bproto,
			LoadBalancerPort:     lp,
			LoadBalancerProtocol: lproto,
		}

		if v, ok := data["ssl_certificate_id"]; ok && v != "" {
			vs := v.(string)
			l.ServerCertificateId = &vs
		}

		var valid bool
		if l.ServerCertificateId != nil && *l.ServerCertificateId != "" {
			// validate the protocol is correct
			for _, p := range []string{"https", "ssl"} {
				if (strings.ToLower(*l.BackendProtocol) == p) ||
					(strings.ToLower(l.LoadBalancerProtocol) == p) {
					valid = true
				}
			}
		} else {
			valid = true
		}

		if valid {
			listeners = append(listeners, l)
		} else {
			return nil, fmt.Errorf("[ERR] ELB Listener: ssl_certificate_id may be set only when protocol is 'https' or 'ssl'")
		}
	}

	return listeners, nil
}

func resourceOutscaleOAPILoadBalancerListenersCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.CreateLoadBalancerListenersRequest{}

	listener, err := expandListenerForCreation(d.Get("listener").([]interface{}))
	if err != nil {
		return err
	}

	req.Listeners = listener

	if v, ok := d.GetOk("load_balancer_name"); ok {
		req.LoadBalancerName = v.(string)
	}

	elbOpts := oscgo.CreateLoadBalancerListenersOpts{
		optional.NewInterface(req),
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.ListenerApi.CreateLoadBalancerListeners(
			context.Background(), &elbOpts)

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

	d.SetId(req.LoadBalancerName)
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
			"backend_port":           i.BackendPort,
			"backend_protocol":       i.BackendProtocol,
			"load_balancer_port":     i.LoadBalancerPort,
			"load_balancer_protocol": i.LoadBalancerProtocol,
			"server_certificate_id":  i.ServerCertificateId,
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
	conn := meta.(*OutscaleClient).OSCAPI

	d.Partial(true)

	if d.HasChange("listener") {
		o, n := d.GetChange("listener")
		os := o.([]interface{})
		ns := n.([]interface{})

		remove, _ := expandListeners(ns)
		add, _ := expandListenerForCreation(os)

		if len(remove) > 0 {
			ports := make([]int64, 0, len(remove))
			for _, listener := range remove {
				ports = append(ports, *listener.LoadBalancerPort)
			}

			req := oscgo.DeleteLoadBalancerListenersRequest{
				LoadBalancerName:  d.Id(),
				LoadBalancerPorts: ports,
			}

			deleteListenersOpts := &oscgo.DeleteLoadBalancerListenersOpts{
				optional.NewInterface(req),
			}

			var err error
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, _, err = conn.ListenerApi.DeleteLoadBalancerListeners(
					context.Background(), deleteListenersOpts)

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
			req := oscgo.CreateLoadBalancerListenersRequest{
				LoadBalancerName: d.Id(),
				Listeners:        add,
			}

			createListenersOpts := &oscgo.CreateLoadBalancerListenersOpts{
				optional.NewInterface(req),
			}

			var err error
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, _, err = conn.ListenerApi.CreateLoadBalancerListeners(
					context.Background(), createListenersOpts)
				if err != nil {
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
	conn := meta.(*OutscaleClient).OSCAPI

	remove, _ := expandListeners(d.Get("listener").([]interface{}))

	ports := make([]int64, 0, len(remove))
	for _, listener := range remove {
		ports = append(ports, *listener.LoadBalancerPort)
	}

	req := oscgo.DeleteLoadBalancerListenersRequest{
		LoadBalancerName:  d.Id(),
		LoadBalancerPorts: ports,
	}

	deleteListenersOpts := &oscgo.DeleteLoadBalancerListenersOpts{
		optional.NewInterface(req),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.ListenerApi.DeleteLoadBalancerListeners(
			context.Background(),
			deleteListenersOpts)

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
