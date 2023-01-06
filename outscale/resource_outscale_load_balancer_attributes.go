package outscale

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOutscaleOAPILoadBalancerAttributes() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILoadBalancerAttributesCreate,
		Update: resourceOutscaleOAPILoadBalancerAttributesUpdate,
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
							Default:  60,
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
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								htVal := val.(int)
								if htVal < 5 || htVal > 600 {
									errs = append(errs, fmt.Errorf("%q must be between 5 and 600 inclusive, got: %d", key, htVal))
								}
								return
							},
						},
						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								uhtVal := val.(int)
								if uhtVal < 2 || uhtVal > 10 {
									errs = append(errs, fmt.Errorf("%q must be between 2 and 10 inclusive, got: %d", key, uhtVal))
								}
								return
							},
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
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								portVal := val.(int)
								if portVal < utils.MinPort || portVal > utils.MaxPort {
									errs = append(errs, fmt.Errorf("%q must be between %d and %d inclusive, got: %d", key, utils.MinPort, utils.MaxPort, portVal))
								}
								return
							},
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"check_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								ciVal := val.(int)
								if ciVal < 5 || ciVal > 600 {
									errs = append(errs, fmt.Errorf("%q must be between 5 and 600 inclusive, got: %d", key, ciVal))
								}
								return
							},
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								tVal := val.(int)
								if tVal < 5 || tVal > 60 {
									errs = append(errs, fmt.Errorf("%q must be between 5 and 60 inclusive, got: %d", key, tVal))
								}
								return
							},
						},
					},
				},
			},
			"listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: lb_listener_schema(true),
				},
			},

			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"subregion_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"load_balancer_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"tags": tagsListOAPISchema2(true),
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"security_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"server_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"source_security_group": lb_sg_schema(),

			"backend_vm_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

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

func resourceOutscaleOAPILoadBalancerAttributesUpdate(d *schema.ResourceData,
	meta interface{}) error {
	return resourceOutscaleOAPILoadBalancerAttributesCreate_(d, meta, true)
}

func resourceOutscaleOAPILoadBalancerAttributesCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceOutscaleOAPILoadBalancerAttributesCreate_(d, meta, false)
}

func loadBalancerAttributesDoRequest(d *schema.ResourceData, meta interface{}, req oscgo.UpdateLoadBalancerRequest) error {
	conn := meta.(*OutscaleClient).OSCAPI
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.LoadBalancerApi.UpdateLoadBalancer(
			context.Background()).UpdateLoadBalancerRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
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

func resourceOutscaleOAPILoadBalancerAttributesCreate_(d *schema.ResourceData, meta interface{}, isUpdate bool) error {
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

	if pol_names, plnok := d.GetOk("policy_names"); plnok {
		m := pol_names.([]interface{})
		a := make([]string, len(m))
		for k, v := range m {
			a[k] = v.(string)
		}
		req.PolicyNames = &a
	} else if isUpdate {
		a := make([]string, 0)
		req.PolicyNames = &a
	}
	if isUpdate {
		return loadBalancerAttributesDoRequest(d, meta, req)
	}

	if ssl, sok := d.GetOk("server_certificate_id"); sok {
		ssl_s := ssl.(string)
		req.ServerCertificateId = &ssl_s
	}

	if al, alok := d.GetOk("access_log"); alok {
		dals := al.([]interface{})
		dal := dals[0].(map[string]interface{})
		check, _ := dal["is_enabled"]
		access := &oscgo.AccessLog{}

		if check != nil {
			is_enable := check.(bool)
			access.IsEnabled = &is_enable
		}

		if v := dal["publication_interval"]; v != nil {
			pi := int32(v.(int))
			access.PublicationInterval = &pi
		}

		obn := dal["osu_bucket_name"]
		if obn != nil && obn.(string) != "" {
			obn_s := obn.(string)
			access.OsuBucketName = &obn_s
		}

		obp := dal["osu_bucket_prefix"]
		if obp != nil && obp.(string) != "" {
			obp_s := obp.(string)
			access.OsuBucketPrefix = &obp_s
		}
		req.AccessLog = access
	}

	hcs, hok := d.GetOk("health_check")
	if hok {
		hc := hcs.([]interface{})
		check := hc[0].(map[string]interface{})
		var healthCheck oscgo.HealthCheck
		healthCheck.SetHealthyThreshold(int32(check["healthy_threshold"].(int)))
		healthCheck.SetUnhealthyThreshold(int32(check["unhealthy_threshold"].(int)))
		healthCheck.SetCheckInterval(int32(check["check_interval"].(int)))
		protocol := check["protocol"].(string)
		if protocol == "" {
			return fmt.Errorf("please provide protocol in health_check argument")
		}
		healthCheck.SetProtocol(protocol)
		if path := check["path"].(string); path != "" {
			healthCheck.SetPath(path)
		}
		healthCheck.SetPort(int32(check["port"].(int)))
		healthCheck.SetTimeout(int32(check["timeout"].(int)))
		req.SetHealthCheck(healthCheck)
	}

	return loadBalancerAttributesDoRequest(d, meta, req)
}

func resourceOutscaleOAPILoadBalancerAttributesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	elbName := d.Id()

	lb, _, err := readResourceLb(conn, elbName)
	if err != nil {
		return err
	}
	if lb == nil {
		utils.LogManuallyDeleted("LoadBalancerAttributes", elbName)
		d.SetId("")
		return nil
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

	if lb.SourceSecurityGroup != nil {
		d.Set("source_security_group", flattenSource_sg(lb.SourceSecurityGroup))
	} else {
		d.Set("security_groups", make([]map[string]interface{}, 0))
	}

	if lb.SecurityGroups != nil {
		d.Set("security_groups", utils.StringSlicePtrToInterfaceSlice(lb.SecurityGroups))
	} else {
		d.Set("security_groups", make([]map[string]interface{}, 0))
	}

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

	d.Set("backend_vm_ids", utils.StringSlicePtrToInterfaceSlice(lb.BackendVmIds))

	d.Set("subnets", utils.StringSlicePtrToInterfaceSlice(lb.Subnets))

	d.Set("subregion_names", utils.StringSlicePtrToInterfaceSlice(lb.SubregionNames))

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

	d.Set("health_check", flattenOAPIHealthCheck(lb.HealthCheck))
	d.Set("listeners", flattenOAPIListeners(lb.Listeners))
	d.Set("dns_name", lb.DnsName)

	return nil
}

func resourceOutscaleOAPILoadBalancerAttributesDelete(d *schema.ResourceData, meta interface{}) error {
	var err error

	conn := meta.(*OutscaleClient).OSCAPI
	ename, ok := d.GetOk("load_balancer_name")

	if !ok {
		return fmt.Errorf("required load_balancer_name attributes")
	}

	_, ok = d.GetOk("policy_names")
	if !ok {
		return nil
	}
	a := make([]string, 0)

	p, pok := d.GetOk("load_balancer_port")
	if !pok {
		return fmt.Errorf("required load_balancer_port attributes")
	}
	p32 := int32(p.(int))

	req := oscgo.UpdateLoadBalancerRequest{
		LoadBalancerName: ename.(string),
		PolicyNames:      &a,
		LoadBalancerPort: &p32,
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.LoadBalancerApi.UpdateLoadBalancer(
			context.Background()).UpdateLoadBalancerRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
