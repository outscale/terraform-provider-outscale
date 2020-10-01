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
			"load_balancer_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"server_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_log": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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

	ename, ok := d.GetOk("load_balancer_name")

	if !ok {
		return fmt.Errorf("please provide the is_enabled and load_balancer_name required attributes")
	}

	req := oscgo.UpdateLoadBalancerRequest{
		LoadBalancerName: ename.(string),
	}

	if port, pok := d.GetOk("load_balancer_port"); pok {
		port_i := int64(port.(int))
		req.LoadBalancerPort = &port_i
	}
	if ssl, sok := d.GetOk("server_certificate_id"); sok {
		ssl_s := ssl.(string)
		req.ServerCertificateId = &ssl_s
	}

	if al, alok := d.GetOk("access_log"); alok {
		dal := al.(map[string]interface{})

		is_enable := dal["is_enable"].(bool)
		access := &oscgo.AccessLog{
			IsEnabled: &is_enable,
		}

		if v, ok := lb_atoi_at(dal, "publication_interval"); ok {
			pi := int64(v)
			access.PublicationInterval = &pi
		}
		obn := dal["osu_bucket_name"]
		if obn != nil {
			obn_s := obn.(string)
			access.OsuBucketName = &obn_s
		}
		obp := dal["osu_bucket_prefix"]
		if obp != nil {
			obp_s := obp.(string)
			access.OsuBucketPrefix = &obp_s
		}
		req.AccessLog = access
	}

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

	lb_resp, resp, err := readResourceLb(conn, elbName)
	if err != nil {
		return err
	}

	if lb_resp.AccessLog == nil {
		return fmt.Errorf("NO Attributes FOUND")
	}

	a := lb_resp.AccessLog

	if a != nil {
		access := make(map[string]string)
		access["publication_interval"] = strconv.Itoa(int(*a.PublicationInterval))
		access["is_enabled"] = strconv.FormatBool(*a.IsEnabled)
		access["osu_bucket_name"] = *a.OsuBucketName
		access["osu_bucket_prefix"] = *a.OsuBucketPrefix
		d.Set("access_log", access)
	}

	d.Set("request_id", resp.ResponseContext.RequestId)
	return nil
}

func resourceOutscaleOAPILoadBalancerAttributesUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.UpdateLoadBalancerRequest{}
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
