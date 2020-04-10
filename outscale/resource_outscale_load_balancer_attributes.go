package outscale

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"publication_interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"osu_bucket_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"osu_bucket_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"load_balancer_attributes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_log": {
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"publication_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"is_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"osu_bucket_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"osu_bucket_prefix": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
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

func resourceOutscaleOAPILoadBalancerAttributesCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	v, ok := d.GetOk("is_enabled")
	v1, ok1 := d.GetOk("load_balancer_name")

	if !ok && !ok1 {
		return fmt.Errorf("please provide the is_enabled and load_balancer_name required attributes")
	}

	req := &oscgo.UpdateLoadBalancerRequest{
		LoadBalancerName: v1.(string),
	}

	is_enable := v.(bool)
	access := &oscgo.AccessLog{
		IsEnabled: &is_enable,
	}

	if v, ok := d.GetOk("publication_interval"); ok {
		pi := int64(v.(int))
		access.PublicationInterval = &pi
	}
	if v, ok := d.GetOk("osu_bucket_name"); ok {
		obn := v.(string)
		access.OsuBucketName = &obn
	}
	if v, ok := d.GetOk("osu_bucket_prefix"); ok {
		obp := v.(string)
		access.OsuBucketPrefix = &obp
	}

	req.AccessLog = access

	elbOpts := oscgo.UpdateLoadBalancerOpts{
		optional.NewInterface(req),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.LoadBalancerApi.UpdateLoadBalancer(
			context.Background(), &elbOpts)

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

	d.SetId(req.LoadBalancerName)
	log.Printf("[INFO] LBU Attr ID: %s", d.Id())

	return resourceOutscaleOAPILoadBalancerAttributesRead(d, meta)
}

func resourceOutscaleOAPILoadBalancerAttributesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	elbName := d.Id()

	// Retrieve the LBU Attr properties for updating the state
	filter := &oscgo.FiltersLoadBalancer{
		LoadBalancerNames: &[]string{elbName},
	}

	req := &oscgo.ReadLoadBalancersRequest{
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

		return fmt.Errorf("Error retrieving LBU Attr: %s", err)
	}

	if resp.LoadBalancers == nil {
		return fmt.Errorf("NO ELB FOUND")
	}

	if len(*resp.LoadBalancers) != 1 {
		return fmt.Errorf("Unable to find ELB: %#v", resp.LoadBalancers)
	}

	lb_resp := (*resp.LoadBalancers)[0]
	if lb_resp.AccessLog == nil {
		return fmt.Errorf("NO Attributes FOUND")
	}

	a := lb_resp.AccessLog

	access := make(map[string]string)
	ac := make(map[string]interface{})
	access["publication_interval"] = strconv.Itoa(int(*a.PublicationInterval))
	access["is_enabled"] = strconv.FormatBool(*a.IsEnabled)
	access["osu_bucket_name"] = *a.OsuBucketName
	access["osu_bucket_prefix"] = *a.OsuBucketPrefix
	ac["access_log"] = access

	l := make([]map[string]interface{}, 1)
	l[0] = ac

	d.Set("request_id", resp.ResponseContext.RequestId)

	return d.Set("load_balancer_attributes", l)
}

func resourceOutscaleOAPILoadBalancerAttributesUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := &oscgo.UpdateLoadBalancerRequest{}
	access := &oscgo.AccessLog{}
	if d.HasChange("load_balancer_name") {
		_, n := d.GetChange("load_balancer_name")

		req.LoadBalancerName = n.(string)
	}
	if d.HasChange("is_enabled") {
		_, n := d.GetChange("is_enabled")

		b, err := strconv.ParseBool(n.(string))
		if err != nil {
			return err
		}

		access.IsEnabled = &b
	}
	if d.HasChange("publication_interval") {
		_, n := d.GetChange("publication_interval")

		i, err := strconv.Atoi(n.(string))
		if err != nil {
			return err
		}
		i64 := int64(i)
		access.PublicationInterval = &i64
	}
	if d.HasChange("osu_bucket_name") {
		_, n := d.GetChange("osu_bucket_name")

		s := n.(string)
		access.OsuBucketName = &s
	}
	if d.HasChange("osu_bucket_prefix") {
		_, n := d.GetChange("osu_bucket_prefix")

		s := n.(string)
		access.OsuBucketPrefix = &s
	}

	req.AccessLog = access

	elbOpts := &oscgo.UpdateLoadBalancerOpts{
		optional.NewInterface(req),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.LoadBalancerApi.UpdateLoadBalancer(
			context.Background(), elbOpts)

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
