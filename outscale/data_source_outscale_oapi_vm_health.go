package outscale

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOutscaleOAPIVMHealth() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIVMHealthRead,

		Schema: map[string]*schema.Schema{
			"backend_vm_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vm_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backend_vm_health": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
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

func dataSourceOutscaleOAPIVMHealthRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	params := &lbu.DescribeInstanceHealthInput{
		LoadBalancerName: aws.String(d.Get("load_balancer_name").(string)),
	}

	if v, ok := d.GetOk("backend_vm_id"); ok {
		in := make([]*lbu.Instance, len(v.([]interface{})))

		for k, v1 := range v.([]interface{}) {
			i := v1.(map[string]interface{})
			in[k] = &lbu.Instance{
				InstanceId: aws.String(i["vm_id"].(string)),
			}
		}
		params.Instances = in
	}

	var resp *lbu.DescribeInstanceHealthResult
	var rs *lbu.DescribeInstanceHealthOutput
	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		rs, err = conn.API.DescribeInstanceHealth(params)
		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if rs != nil {
			resp = rs.DescribeInstanceHealthResult
		}
		return nil
	})

	if err != nil {
		return err
	}

	health := make([]map[string]interface{}, len(resp.InstanceStates))

	for k, v := range resp.InstanceStates {
		h := make(map[string]interface{})
		h["description"] = aws.StringValue(v.Description)
		h["vm_id"] = aws.StringValue(v.InstanceId)
		h["comment"] = aws.StringValue(v.ReasonCode)
		h["state"] = aws.StringValue(v.State)
		health[k] = h
	}

	d.Set("request_id", rs.ResponseMetadata.RequestID)
	d.SetId(resource.UniqueId())

	return d.Set("backend_vm_health", health)
}
