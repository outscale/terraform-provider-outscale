package outscale

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func resourceOutscaleLoadBalancerAttributes() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleLoadBalancerAttributesCreate,
		Read:   resourceOutscaleLoadBalancerAttributesRead,
		Update: resourceOutscaleLoadBalancerAttributesUpdate,
		Delete: resourceOutscaleLoadBalancerAttributesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_attributes": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_log": &schema.Schema{
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"emit_interval": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"enabled": &schema.Schema{
										Type:     schema.TypeBool,
										Required: true,
									},

									"s3_bucket_name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"s3_bucket_prefix": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
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

func resourceOutscaleLoadBalancerAttributesCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	v, ok := d.GetOk("load_balancer_attributes")
	v1, ok1 := d.GetOk("load_balancer_name")

	if !ok && ok1 {
		return fmt.Errorf("please provide the required attributes")
	}

	elbOpts := &lbu.ModifyLoadBalancerAttributesInput{
		LoadBalancerName: aws.String(v1.(string)),
	}

	lb := v.([]interface{})[0].(map[string]interface{})

	a := lb["access_log"].(map[string]interface{})

	if v, ok := a["enabled"]; ok && v == "" {
		return fmt.Errorf("please provide the enable attribute")
	}

	b, err := strconv.ParseBool(a["enabled"].(string))
	if err != nil {
		return err
	}

	access := &lbu.AccessLog{
		Enabled: aws.Bool(b),
	}

	if v, ok := a["emit_interval"]; ok && v != "" {
		i, err := strconv.Atoi(v.(string))
		if err != nil {
			return err
		}
		access.EmitInterval = aws.Int64(int64(i))
	}
	if v, ok := a["s3_bucket_name"]; ok && v != "" {
		access.S3BucketName = aws.String(v.(string))
	}
	if v, ok := a["s3_bucket_prefix"]; ok && v != "" {
		access.S3BucketPrefix = aws.String(v.(string))
	}

	elbOpts.LoadBalancerAttributes = &lbu.LoadBalancerAttributes{
		AccessLog: access,
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.API.ModifyLoadBalancerAttributes(elbOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error creating LBU Attr Listener with SSL Cert, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(*elbOpts.LoadBalancerName)
	log.Printf("[INFO] LBU Attr ID: %s", d.Id())

	return resourceOutscaleLoadBalancerAttributesRead(d, meta)
}

func resourceOutscaleLoadBalancerAttributesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU
	elbName := d.Id()

	// Retrieve the LBU Attr properties for updating the state
	describeElbOpts := &lbu.DescribeLoadBalancerAttributesInput{
		LoadBalancerName: aws.String(elbName),
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

	a := describeResp.LoadBalancerAttributes.AccessLog

	access := make(map[string]string)
	ac := make(map[string]interface{})
	access["emit_interval"] = strconv.Itoa(int(aws.Int64Value(a.EmitInterval)))
	access["enabled"] = strconv.FormatBool(aws.BoolValue(a.Enabled))
	access["s3_bucket_name"] = aws.StringValue(a.S3BucketName)
	access["s3_bucket_prefix"] = aws.StringValue(a.S3BucketPrefix)
	ac["access_log"] = access

	l := make([]map[string]interface{}, 1)
	l[0] = ac

	// d.Set("request_id", resp.ResponseMetadata.RequestID)

	return d.Set("load_balancer_attributes", l)
}

func resourceOutscaleLoadBalancerAttributesUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	d.Partial(true)

	if d.HasChange("load_balancer_attributes") {
		v := d.Get("load_balancer_attributes")
		v1 := d.Get("load_balancer_name")

		elbOpts := &lbu.ModifyLoadBalancerAttributesInput{
			LoadBalancerName: aws.String(v1.(string)),
		}

		lb := v.([]interface{})[0].(map[string]interface{})

		a := lb["access_log"].(map[string]interface{})
		b, err := strconv.ParseBool(a["enabled"].(string))
		if err != nil {
			return err
		}

		access := &lbu.AccessLog{
			Enabled: aws.Bool(b),
		}

		if v, ok := a["emit_interval"]; ok {
			i, err := strconv.Atoi(v.(string))
			if err != nil {
				return err
			}
			access.EmitInterval = aws.Int64(int64(i))
		}
		if v, ok := a["s3_bucket_name"]; ok {
			access.S3BucketName = aws.String(v.(string))
		}
		if v, ok := a["s3_bucket_prefix"]; ok {
			access.S3BucketPrefix = aws.String(v.(string))
		}

		elbOpts.LoadBalancerAttributes = &lbu.LoadBalancerAttributes{
			AccessLog: access,
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = conn.API.ModifyLoadBalancerAttributes(elbOpts)

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "Throttling") {
					return resource.RetryableError(
						fmt.Errorf("[WARN] Error creating LBU Attr Listener with SSL Cert, retrying: %s", err))
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return err
		}

		d.SetPartial("load_balancer_attributes")
	}

	d.Partial(false)

	return resourceOutscaleLoadBalancerAttributesRead(d, meta)
}

func resourceOutscaleLoadBalancerAttributesDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	return nil
}
