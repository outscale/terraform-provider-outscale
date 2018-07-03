package outscale

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOutscaleVMHealth() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleVMHealthRead,

		Schema: map[string]*schema.Schema{
			"instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
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
			"instance_states": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reason_code": {
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

func dataSourceOutscaleVMHealthRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	params := &lbu.DescribeInstanceHealthInput{
		LoadBalancerName: aws.String(d.Get("load_balancer_name").(string)),
	}

	if v, ok := d.GetOk("instances"); ok {
		in := make([]*lbu.Instance, len(v.([]interface{})))

		for k, v1 := range v.([]interface{}) {
			i := v1.(map[string]interface{})
			in[k] = &lbu.Instance{
				InstanceId: aws.String(i["instance_id"].(string)),
			}
		}
		params.Instances = in
	}

	var rs *lbu.DescribeInstanceHealthOutput
	var resp *lbu.DescribeInstanceHealthResult
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
		h["instance_id"] = aws.StringValue(v.InstanceId)
		h["reason_code"] = aws.StringValue(v.ReasonCode)
		h["state"] = aws.StringValue(v.State)
		health[k] = h
	}

	d.Set("request_id", rs.ResponseMetadata.RequestID)
	d.SetId(resource.UniqueId())

	return d.Set("instance_states", health)
}
