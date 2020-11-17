package outscale

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/osc"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOutscaleOAPILoadBalancerAttributes() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILoadBalancerAttributesCreate,
		Read:   resourceOutscaleOAPILoadBalancerAttributesRead,
		Delete: resourceOutscaleOAPILoadBalancerAttributesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"access_log": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"osu_bucket_name": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"osu_bucket_prefix": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"publication_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
					},
				},
			},
			"health_check": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"healthy_threshold": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"unhealthy_threshold": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"path": {
							Type:     schema.TypeString,
							ForceNew: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if new == "" && old == "/" {
									return true
								}
								return old == new
							},
							Optional: true,
						},
						"port": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"check_interval": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"timeout": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"listeners": {
				Type:     schema.TypeSet,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: lb_listener_schema(),
				},
			},
			"load_balancer_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"server_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_security_group": lb_sg_schema(),
			"application_sticky_cookie_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cookie_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"load_balancer_sticky_cookie_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"load_balancer_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"load_balancer_type": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"policy_names": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func isLoadBalancerNotFound(err error) bool {
	return strings.Contains(fmt.Sprint(err), "LoadBalancerNotFound")
}

func lb_atoi_at(hc map[string]interface{}, el string) (int, bool) {
	hc_el := hc[el]

	if hc_el == nil {
		return 0, false
	}

	r, err := strconv.Atoi(hc_el.(string))
	return r, err == nil
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
		port_i := int32(port.(int))
		req.LoadBalancerPort = &port_i
	}

	if ssl, sok := d.GetOk("server_certificate_id"); sok {
		ssl_s := ssl.(string)
		req.ServerCertificateId = &ssl_s
	}

	if pol_names, plnok := d.GetOk("policy_names"); plnok {
		m := pol_names.([]interface{})
		a := make([]string, len(m))
		for k, v := range m {
			a[k] = v.(string)
		}
		req.PolicyNames = &a
	}

	if al, alok := d.GetOk("access_log"); alok {
		dals := al.([]interface{})
		dal := dals[0].(map[string]interface{})
		check, _ := dal["is_enabled"]
		is_enable := false
		if check == "true" {
			is_enable = true
		}
		access := &oscgo.AccessLog{
			IsEnabled: &is_enable,
		}

		if v, ok := lb_atoi_at(dal, "publication_interval"); ok {
			pi := int32(v)
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

	hcs, hok := d.GetOk("health_check")
	if hok {
		hc := hcs.([]interface{})
		check := hc[0].(map[string]interface{})

		ht, ut, sucess := 0, 0, false
		if ht, sucess = lb_atoi_at(check, "healthy_threshold"); sucess == false {
			return fmt.Errorf("please provide an number in health_check.healthy_threshold argument")

		}

		if ut, sucess = lb_atoi_at(check, "unhealthy_threshold"); sucess == false {
			return fmt.Errorf("please provide an number in health_check.unhealthy_threshold argument")
		}

		i, ierr := lb_atoi_at(check, "check_interval")
		t, terr := lb_atoi_at(check, "timeout")
		p, perr := lb_atoi_at(check, "port")

		if ierr != true {
			return fmt.Errorf("please provide an number in health_check.check_interval argument")
		}

		if terr != true {
			return fmt.Errorf("please provide an number in health_check.timeout argument")
		}

		if perr != true {
			return fmt.Errorf("please provide an number in health_check.port argument")
		}

		var hc_req oscgo.HealthCheck
		hc_req.HealthyThreshold = int32(ht)
		hc_req.UnhealthyThreshold = int32(ut)
		hc_req.CheckInterval = int32(i)
		hc_req.Protocol = check["protocol"].(string)
		if check["path"] != nil {
			p := check["path"].(string)
			if p != "" {
				hc_req.Path = &p
			}
		}
		hc_req.Port = int32(p)
		hc_req.Timeout = int32(t)
		req.HealthCheck = &hc_req

	}

	var err error
	var resp oscgo.UpdateLoadBalancerResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.LoadBalancerApi.UpdateLoadBalancer(
			context.Background()).UpdateLoadBalancerRequest(req).Execute()

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "400 Bad Request") {
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(
				fmt.Errorf("[WARN] Error creating LBU Attr: %s", err))
		}
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(req.LoadBalancerName)
	log.Printf("[INFO] LBU Attr ID: %s", d.Id())

	d.Set("request_id", resp.ResponseContext.RequestId)
	return resourceOutscaleOAPILoadBalancerAttributesRead(d, meta)
}

func resourceOutscaleOAPILoadBalancerAttributesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	elbName := d.Id()

	lb, _, err := readResourceLb(conn, elbName)
	if err != nil {
		return err
	}

	a := lb.AccessLog

	if a != nil {
		ac := make([]interface{}, 1)
		access := make(map[string]interface{})
		access["publication_interval"] = int(*a.PublicationInterval)
		access["is_enabled"] = *a.IsEnabled
		access["osu_bucket_name"] = a.OsuBucketName
		access["osu_bucket_prefix"] = a.OsuBucketPrefix
		ac[0] = access
		err := d.Set("access_log", ac)
		if err != nil {
			return err
		}
	} else {
		d.Set("access_log", make([]interface{}, 0))
	}

	sgr := make(map[string]string)
	if lb.SourceSecurityGroup != nil {
		sgr["security_group_name"] = *lb.SourceSecurityGroup.SecurityGroupName
		sgr["security_group_account_id"] = *lb.SourceSecurityGroup.SecurityGroupAccountId
	}
	d.Set("source_security_group", sgr)

	if lb.Tags != nil {
		ta := make([]map[string]interface{}, len(*lb.Tags))
		for k1, v1 := range *lb.Tags {
			t := make(map[string]interface{})
			t["key"] = v1.Key
			t["value"] = v1.Value
			ta[k1] = t
		}

		d.Set("tags", ta)
	} else {
		d.Set("tags", make([]map[string]interface{}, 0))
	}

	if lb.ApplicationStickyCookiePolicies != nil {
		app := make([]map[string]interface{},
			len(*lb.ApplicationStickyCookiePolicies))
		for k, v := range *lb.ApplicationStickyCookiePolicies {
			a := make(map[string]interface{})
			a["cookie_name"] = v.CookieName
			a["policy_name"] = v.PolicyName
			app[k] = a
		}
		d.Set("application_sticky_cookie_policies", app)
	} else {
		d.Set("application_sticky_cookie_policies",
			make([]map[string]interface{}, 0))
	}

	if lb.LoadBalancerStickyCookiePolicies == nil {
		d.Set("load_balancer_sticky_cookie_policies",
			make([]map[string]interface{}, 0))
	} else {
		lbc := make([]map[string]interface{},
			len(*lb.LoadBalancerStickyCookiePolicies))
		for k, v := range *lb.LoadBalancerStickyCookiePolicies {
			a := make(map[string]interface{})
			a["policy_name"] = v.PolicyName
			lbc[k] = a
		}
		d.Set("load_balancer_sticky_cookie_policies", lbc)
	}

	hls := make([]interface{}, 1)
	hls[0] = flattenOAPIHealthCheck(lb.HealthCheck)
	d.Set("health_check", hls)
	d.Set("listeners", flattenOAPIListeners(lb.Listeners))
	return nil
}

func resourceOutscaleOAPILoadBalancerAttributesDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	return nil
}
