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

func resourceOutscaleOAPILoadBalancerAttributes() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILoadBalancerAttributesCreate,
		Read:   resourceOutscaleOAPILoadBalancerAttributesRead,
		Update: resourceOutscaleOAPILoadBalancerAttributesUpdate,
		Delete: resourceOutscaleOAPILoadBalancerAttributesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"publication_interval": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"is_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"osu_bucket_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"osu_bucket_prefix": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"load_balancer_attributes": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_log": &schema.Schema{
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"publication_interval": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"is_enabled": &schema.Schema{
										Type:     schema.TypeBool,
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

func resourceOutscaleOAPILoadBalancerAttributesCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	v, ok := d.GetOk("is_enabled")
	v1, ok1 := d.GetOk("load_balancer_name")

	if !ok && !ok1 {
		return fmt.Errorf("please provide the is_enabled and load_balancer_name required attributes")
	}

	elbOpts := &lbu.ModifyLoadBalancerAttributesInput{
		LoadBalancerName: aws.String(v1.(string)),
	}
	access := &lbu.AccessLog{
		Enabled: aws.Bool(v.(bool)),
	}

	if v, ok := d.GetOk("publication_interval"); ok {
		access.EmitInterval = aws.Int64(int64(v.(int)))
	}
	if v, ok := d.GetOk("osu_bucket_name"); ok {
		access.S3BucketName = aws.String(v.(string))
	}
	if v, ok := d.GetOk("osu_bucket_prefix"); ok {
		access.S3BucketPrefix = aws.String(v.(string))
	}

	elbOpts.LoadBalancerAttributes = &lbu.LoadBalancerAttributes{
		AccessLog: access,
	}

	var err error
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

	return resourceOutscaleOAPILoadBalancerAttributesRead(d, meta)
}

func resourceOutscaleOAPILoadBalancerAttributesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU
	elbName := d.Id()

	// Retrieve the LBU Attr properties for updating the state
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

		return fmt.Errorf("Error retrieving LBU Attr: %s", err)
	}

	if describeResp.LoadBalancerAttributes == nil {
		return fmt.Errorf("NO Attributes FOUND")
	}

	a := describeResp.LoadBalancerAttributes.AccessLog

	access := make(map[string]string)
	ac := make(map[string]interface{})
	access["publication_interval"] = strconv.Itoa(int(aws.Int64Value(a.EmitInterval)))
	access["is_enabled"] = strconv.FormatBool(aws.BoolValue(a.Enabled))
	access["osu_bucket_name"] = aws.StringValue(a.S3BucketName)
	access["osu_bucket_prefix"] = aws.StringValue(a.S3BucketPrefix)
	ac["access_log"] = access

	l := make([]map[string]interface{}, 1)
	l[0] = ac

	d.Set("request_id", resp.ResponseMetadata.RequestID)

	return d.Set("load_balancer_attributes", l)
}

func resourceOutscaleOAPILoadBalancerAttributesUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	elbOpts := &lbu.ModifyLoadBalancerAttributesInput{}
	access := &lbu.AccessLog{}
	if d.HasChange("load_balancer_name") {
		_, n := d.GetChange("load_balancer_name")

		elbOpts.LoadBalancerName = aws.String(n.(string))
	}
	if d.HasChange("is_enabled") {
		_, n := d.GetChange("is_enabled")

		b, err := strconv.ParseBool(n.(string))
		if err != nil {
			return err
		}

		access.Enabled = aws.Bool(b)
	}
	if d.HasChange("publication_interval") {
		_, n := d.GetChange("publication_interval")

		i, err := strconv.Atoi(n.(string))
		if err != nil {
			return err
		}
		access.EmitInterval = aws.Int64(int64(i))
	}
	if d.HasChange("osu_bucket_name") {
		_, n := d.GetChange("osu_bucket_name")

		access.S3BucketName = aws.String(n.(string))
	}
	if d.HasChange("osu_bucket_prefix") {
		_, n := d.GetChange("osu_bucket_prefix")
		access.S3BucketPrefix = aws.String(n.(string))
	}

	elbOpts.LoadBalancerAttributes = &lbu.LoadBalancerAttributes{
		AccessLog: access,
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.API.ModifyLoadBalancerAttributes(elbOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceOutscaleOAPILoadBalancerAttributesRead(d, meta)
}

func resourceOutscaleOAPILoadBalancerAttributesDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	return nil
}
