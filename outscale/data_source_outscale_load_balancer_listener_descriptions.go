package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceOutscaleOAPILoadBalancerLDs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPILoadBalancerLDsRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_names": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"listener": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"listener": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"backend_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"backend_protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"load_balancer_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"load_balancer_protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"server_certificate_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"policy_name": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPILoadBalancerLDsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	e, ok := d.GetOk("load_balancer_names")

	if !ok {
		return fmt.Errorf("please provide the name of the load balancers")
	}

	elbName := e.(string)

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

		return fmt.Errorf("Error retrieving LBU: %s", err)
	}

	lds := make([]map[string]interface{}, len(*resp.LoadBalancers))
	for k, v1 := range *resp.LoadBalancers {
		ld := make(map[string]interface{})
		ls := make([]map[string]interface{}, len(*v1.Listeners))

		for k1, v2 := range *v1.Listeners {
			l := make(map[string]interface{})
			l["backend_port"] = v2.BackendPort
			l["backend_protocol"] = v2.BackendProtocol
			l["load_balancer_port"] = v2.LoadBalancerPort
			l["load_balancer_protocol"] = v2.LoadBalancerProtocol
			l["server_certificate_id"] = v2.ServerCertificateId
			l["policy_name"] = flattenStringList(v2.PolicyNames)
			ls[k1] = l
		}

		ld["listener"] = ls

		lds[k] = ld
	}

	d.Set("request_id", resp.ResponseContext.RequestId)
	d.SetId(resource.UniqueId())

	return d.Set("listener", lds)
}
