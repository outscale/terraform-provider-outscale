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

func lb_listener_schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"policy_names": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func resourceOutscaleOAPILoadBalancerListeners() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILoadBalancerListenersCreate,
		Read:   resourceOutscaleOAPILoadBalancerListenersRead,
		Delete: resourceOutscaleOAPILoadBalancerListenersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"listeners": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: lb_listener_schema(),
				},
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Flattens an array of Listeners into a []map[string]interface{}
func flattenOAPIListeners(list *[]oscgo.Listener) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(*list))

	for _, i := range *list {
		listener := map[string]interface{}{
			"backend_port":           int(*i.BackendPort),
			"backend_protocol":       *i.BackendProtocol,
			"load_balancer_port":     int(*i.LoadBalancerPort),
			"load_balancer_protocol": *i.LoadBalancerProtocol,
		}
		if i.ServerCertificateId != nil {
			listener["server_certificate_id"] =
				*i.ServerCertificateId
		}
		listener["policy_names"] = flattenStringList(i.PolicyNames)
		result = append(result, listener)
	}
	return result
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

		if v, ok := data["server_certificate_id"]; ok && v != "" {
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
			return nil, fmt.Errorf("[ERR] ELB Listener: server_certificate_id may be set only when protocol is 'https' or 'ssl'")
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

		if v, ok := data["server_certificate_id"]; ok && v != "" {
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
			return nil, fmt.Errorf("[ERR] ELB Listener: server_certificate_id may be set only when protocol is 'https' or 'ssl'")
		}
	}

	return listeners, nil
}

func resourceOutscaleOAPILoadBalancerListenersCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.CreateLoadBalancerListenersRequest{}

	listener, err := expandListenerForCreation(d.Get("listeners").(*schema.Set).List())
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
	conn := meta.(*OutscaleClient).OSCAPI
	ename, ok := d.GetOk("load_balancer_name")

	if !ok {
		return fmt.Errorf("please provide the name of the load balancer")
	}

	elbName := ename.(string)

	// Retrieve the ELB properties for updating the state
	filter := &oscgo.FiltersLoadBalancer{
		LoadBalancerNames: &[]string{elbName},
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
			if strings.Contains(fmt.Sprint(err), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if isLoadBalancerNotFound(err) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving ELB: %s", err)
	}

	if resp.LoadBalancers == nil {
		return fmt.Errorf("NO ELB FOUND")
	}

	if len(*resp.LoadBalancers) != 1 {
		return fmt.Errorf("Unable to find ELB: %#v", resp.LoadBalancers)
	}

	lb := (*resp.LoadBalancers)[0]

	log.Printf("[DEBUG] read lb.Listeners %v", lb.Listeners)
	if lb.Listeners != nil {
		if err := d.Set("listeners", flattenOAPIListeners(lb.Listeners)); err != nil {
			log.Printf("[DEBUG] out err %v", err)
			return err
		}
	} else {
		return fmt.Errorf("Something very wrong happen on Load Balancer: %s", elbName)
	}

	d.Set("load_balancer_name", *lb.LoadBalancerName)
	d.Set("request_id", *resp.ResponseContext.RequestId)

	return nil
}

func resourceOutscaleOAPILoadBalancerListenersDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	remove, _ := expandListeners(d.Get("listeners").(*schema.Set).List())

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
