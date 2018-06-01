package outscale

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func dataSourceOutscaleLoadBalancerLDs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleLoadBalancerLDsRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_names": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"listener_descriptions": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"listener": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_port": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_protocol": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"load_balancer_port": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"ssl_certificate_id": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"policy_names": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleLoadBalancerLDsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU
	e, ok := d.GetOk("load_balancer_names")

	if !ok {
		return fmt.Errorf("please provide the name of the load balancers")
	}

	// Retrieve the ELB properties for updating the state
	describeElbOpts := &lbu.DescribeLoadBalancersInput{
		LoadBalancerNames: expandStringList(e.([]interface{})),
	}

	var describeResp *lbu.DescribeLoadBalancersOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		describeResp, err = conn.API.DescribeLoadBalancers(describeElbOpts)

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

	if describeResp.LoadBalancerDescriptions == nil {
		return fmt.Errorf("NO LBU FOUND")
	}

	if len(describeResp.LoadBalancerDescriptions) < 1 {
		return fmt.Errorf("Unable to find LBUS: %#v", describeResp.LoadBalancerDescriptions)
	}

	lds := make([]map[string]interface{}, len(describeResp.LoadBalancerDescriptions))

	for k, v1 := range describeResp.LoadBalancerDescriptions {
		ld := make(map[string]interface{})
		ls := make([]map[string]interface{}, len(v1.ListenerDescriptions))

		for k1, v2 := range v1.ListenerDescriptions {
			l := make(map[string]interface{})
			l["instance_port"] = strconv.Itoa(int(aws.Int64Value(v2.Listener.InstancePort)))
			l["instance_protocol"] = aws.StringValue(v2.Listener.InstanceProtocol)
			l["load_balancer_port"] = strconv.Itoa(int(aws.Int64Value(v2.Listener.LoadBalancerPort)))
			l["protocol"] = aws.StringValue(v2.Listener.Protocol)
			l["ssl_certificate_id"] = aws.StringValue(v2.Listener.SSLCertificateId)
			ls[k1] = l
		}

		ld["listener"] = ls
		ld["policy_names"] = flattenStringList(v1.ListenerDescriptions[0].PolicyNames)

		lds[k] = ld
	}

	// d.Set("request_id", resp.ResponseMetadata.RequestID)
	d.SetId(resource.UniqueId())

	return d.Set("listener_descriptions", lds)
}
