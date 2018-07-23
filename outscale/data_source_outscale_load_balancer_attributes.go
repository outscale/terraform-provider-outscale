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

func dataSourceOutscaleLoadBalancerAttr() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleLoadBalancerAttrRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"access_log_emit_interval": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_log_enabled": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_log_s3_bucket_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_log_s3_bucket_prefix": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleLoadBalancerAttrRead(d *schema.ResourceData, meta interface{}) error {
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
	d.SetId(elbName)

	a := describeResp.LoadBalancerAttributes.AccessLog
	d.Set("access_log_emit_interval", strconv.Itoa(int(aws.Int64Value(a.EmitInterval))))
	d.Set("access_log_enabled", strconv.FormatBool(aws.BoolValue(a.Enabled)))
	d.Set("access_log_s3_bucket_name", aws.StringValue(a.S3BucketName))
	d.Set("access_log_s3_bucket_prefix", aws.StringValue(a.S3BucketPrefix))

	return d.Set("request_id", resp.ResponseMetadata.RequestID)
}
