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

func dataSourceOutscaleOAPILoadBalancerAttr() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPILoadBalancerAttrRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"load_balancer_attributes": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_log": &schema.Schema{
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"publication_interval": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_enabled": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"osu_bucket_name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"osu_bucket_prefix": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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

func dataSourceOutscaleOAPILoadBalancerAttrRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU
	ename, ok := d.GetOk("load_balancer_name")

	if !ok {
		return fmt.Errorf("please provide the name of the load balancer")
	}

	elbName := ename.(string)

	describeElbOpts := &lbu.DescribeLoadBalancerAttributesInput{
		LoadBalancerName: aws.String(elbName),
	}

	var describeResp *lbu.DescribeLoadBalancerAttributesResult
	var resp *lbu.DescribeLoadBalancerAttributesOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.DescribeLoadBalancerAttributes(describeElbOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if resp.DescribeLoadBalancerAttributesResult != nil {
			describeResp = resp.DescribeLoadBalancerAttributesResult
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

	if describeResp.LoadBalancerAttributes == nil {
		return fmt.Errorf("NO ELB FOUND")
	}

	a := describeResp.LoadBalancerAttributes.AccessLog

	ld := make([]map[string]interface{}, 1)
	acc := make(map[string]interface{})

	acc["publication_interval"] = strconv.Itoa(int(aws.Int64Value(a.EmitInterval)))
	acc["is_enabled"] = strconv.FormatBool(aws.BoolValue(a.Enabled))
	acc["osu_bucket_name"] = aws.StringValue(a.S3BucketName)
	acc["osu_bucket_prefix"] = aws.StringValue(a.S3BucketPrefix)

	ld[0] = map[string]interface{}{"access_log": acc}

	d.Set("request_id", resp.ResponseMetadata.RequestID)
	d.SetId(resource.UniqueId())

	return d.Set("load_balancer_attributes", ld)
}
