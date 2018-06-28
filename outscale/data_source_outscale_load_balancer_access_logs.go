package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleLoadBalancerAccessLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleLoadBalancerAccessLogsRead,

		Schema: map[string]*schema.Schema{
			"emit_interval": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"s3_bucket_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"s3_bucket_prefix": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_balancer_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleLoadBalancerAccessLogsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU
	elbName, ok1 := d.GetOk("load_balancer_name")

	if !ok1 {
		return fmt.Errorf("please provide the load_balancer_name required attribute")
	}

	describeElbOpts := &lbu.DescribeLoadBalancerAttributesInput{
		LoadBalancerName: aws.String(elbName.(string)),
	}

	var describeResp *lbu.DescribeLoadBalancerAttributesOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		describeResp, err = conn.API.DescribeLoadBalancerAttributes(describeElbOpts)

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

		return fmt.Errorf("Error retrieving LBU Attr: %s", err)
	}

	if describeResp.LoadBalancerAttributes == nil {
		return fmt.Errorf("NO Attributes FOUND")
	}

	utils.PrintToJSON(describeResp, "RESPONSE =>")

	a := describeResp.LoadBalancerAttributes.AccessLog

	d.Set("emit_interval", aws.Int64Value(a.EmitInterval))
	d.Set("enabled", aws.BoolValue(a.Enabled))
	d.Set("s3_bucket_name", aws.StringValue(a.S3BucketName))
	d.Set("s3_bucket_prefix", aws.StringValue(a.S3BucketPrefix))

	d.SetId(elbName.(string))
	// d.Set("request_id", describeResp.ResponseMetadata.RequestId)

	return nil
}
