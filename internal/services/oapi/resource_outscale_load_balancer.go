package oapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func lb_sg_schema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"security_group_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"security_group_account_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func ResourceOutscaleLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleLoadBalancerCreate,
		ReadContext:   ResourceOutscaleLoadBalancerRead,
		UpdateContext: ResourceOutscaleLoadBalancerUpdate,
		DeleteContext: ResourceOutscaleLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout),
			Update: schema.DefaultTimeout(UpdateDefaultTimeout),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"subregion_names": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"load_balancer_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subnets": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": TagsSchemaSDK(),

			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_log": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"publication_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"health_check": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"backend_vm_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"backend_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"listeners": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: lb_listener_schema(false),
				},
			},
			"source_security_group": lb_sg_schema(),
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secured_cookies": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"net_id": {
				Type:     schema.TypeString,
				Computed: true,
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Flattens an array of Listeners into a []map[string]interface{}
func flattenOAPIListeners(list *[]osc.Listener) []map[string]interface{} {
	if list == nil {
		return make([]map[string]interface{}, 0)
	}

	result := make([]map[string]interface{}, 0, len(*list))

	for _, i := range *list {
		listener := map[string]interface{}{
			"backend_port":           int(i.BackendPort),
			"backend_protocol":       i.BackendProtocol,
			"load_balancer_port":     int(i.LoadBalancerPort),
			"load_balancer_protocol": i.LoadBalancerProtocol,
		}
		if i.ServerCertificateId != nil {
			listener["server_certificate_id"] = *i.ServerCertificateId
		}
		listener["policy_names"] = utils.StringSlicePtrToInterfaceSlice(&i.PolicyNames)
		result = append(result, listener)
	}
	return result
}

func expandListeners(configured []interface{}) ([]*osc.Listener, error) {
	listeners := make([]*osc.Listener, 0, len(configured))

	for _, lRaw := range configured {
		data := lRaw.(map[string]interface{})

		ip := data["backend_port"].(int)
		lp := data["load_balancer_port"].(int)
		bproto := data["backend_protocol"].(string)
		lproto := data["load_balancer_protocol"].(string)
		l := &osc.Listener{
			BackendPort:          ip,
			BackendProtocol:      bproto,
			LoadBalancerPort:     lp,
			LoadBalancerProtocol: lproto,
		}

		if v, ok := data["server_certificate_id"]; ok && v != "" {
			protocolNeedCerticate := []string{"https", "ssl"}
			if !slices.Contains(protocolNeedCerticate, strings.ToLower(l.BackendProtocol)) &&
				!slices.Contains(protocolNeedCerticate, strings.ToLower(l.LoadBalancerProtocol)) {
				return nil, errors.New("LBU Listener: server_certificate_id may be set only when protocol is 'https' or 'ssl'")
			}
			l.ServerCertificateId = new(v.(string))
		}
		listeners = append(listeners, l)
	}
	return listeners, nil
}

func expandListenerForCreation(configured []interface{}) ([]osc.ListenerForCreation, error) {
	listeners := make([]osc.ListenerForCreation, 0, len(configured))

	for _, lRaw := range configured {
		data := lRaw.(map[string]interface{})

		ip := data["backend_port"].(int)
		lp := data["load_balancer_port"].(int)
		bproto := data["backend_protocol"].(string)
		lproto := data["load_balancer_protocol"].(string)
		l := osc.ListenerForCreation{
			BackendPort:          ip,
			BackendProtocol:      &bproto,
			LoadBalancerPort:     lp,
			LoadBalancerProtocol: lproto,
		}

		if v, ok := data["server_certificate_id"]; ok && v != "" {
			protocolNeedCerticate := []string{"https", "ssl"}
			if !slices.Contains(protocolNeedCerticate, strings.ToLower(ptr.From(l.BackendProtocol))) &&
				!slices.Contains(protocolNeedCerticate, strings.ToLower(l.LoadBalancerProtocol)) {
				return nil, errors.New("LBU Listener: server_certificate_id may be set only when protocol is 'https' or 'ssl'")
			}
			l.ServerCertificateId = new(v.(string))
		}
		listeners = append(listeners, l)
	}

	return listeners, nil
}

func mk_elem(computed bool, required bool,
	t schema.ValueType,
) *schema.Schema {
	if computed {
		return &schema.Schema{
			Type:     t,
			Computed: true,
		}
	} else if required {
		return &schema.Schema{
			Type:     t,
			Required: true,
		}
	} else {
		return &schema.Schema{
			Type:     t,
			Optional: true,
		}
	}
}

func lb_listener_schema(computed bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"backend_port": mk_elem(computed, !computed,
			schema.TypeInt),
		"backend_protocol": mk_elem(computed, !computed,
			schema.TypeString),
		"load_balancer_port": mk_elem(computed, !computed,
			schema.TypeInt),
		"load_balancer_protocol": mk_elem(computed, !computed,
			schema.TypeString),
		"server_certificate_id": mk_elem(computed, false,
			schema.TypeString),
		"policy_names": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func ResourceOutscaleLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceOutscaleLoadBalancerCreate_(ctx, d, meta, false)
}

func ResourceOutscaleLoadBalancerCreate_(ctx context.Context, d *schema.ResourceData, meta interface{}, isUpdate bool) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutCreate)

	req := &osc.CreateLoadBalancerRequest{}

	listeners, err := expandListenerForCreation(d.Get("listeners").(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}

	req.Listeners = listeners

	if v, ok := d.GetOk("load_balancer_name"); ok {
		req.LoadBalancerName = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		r := expandOAPITagsSDK(v.(*schema.Set))
		req.Tags = &r
	}

	if v, ok := d.GetOk("load_balancer_type"); ok {
		s := v.(string)
		req.LoadBalancerType = &s
	}

	if v, ok := d.GetOk("public_ip"); ok {
		s := v.(string)
		req.PublicIp = &s
	}

	if v, ok := d.GetOk("security_groups"); ok {
		req.SecurityGroups = utils.SetToStringSlicePtr(v.(*schema.Set))
	}

	v_sb, sb_ok := d.GetOk("subnets")
	if sb_ok {
		req.Subnets = utils.InterfaceSliceToStringList(v_sb.([]interface{}))
	}

	v_srn, srn_ok := d.GetOk("subregion_names")
	if sb_ok && srn_ok {
		return diag.Errorf("can't use both 'subregion_names' and 'subnets'")
	}

	if srn_ok && !sb_ok {
		req.SubregionNames = utils.InterfaceSliceToStringList(v_srn.([]interface{}))
	}

	log.Printf("[DEBUG] Load Balancer request configuration: %#v", *req)
	_, err = client.CreateLoadBalancer(ctx, *req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(err)
	}

	// Assign the lbu's unique identifier for use later
	d.SetId(req.LoadBalancerName)
	log.Printf("[INFO] Load Balancer ID: %s", d.Id())

	if scVal, scOk := d.GetOk("secured_cookies"); scOk {
		req := osc.UpdateLoadBalancerRequest{
			LoadBalancerName: d.Id(),
		}
		req.SecuredCookies = new(scVal.(bool))
		_, err := client.UpdateLoadBalancer(ctx, req, options.WithRetryTimeout(timeout))
		if err != nil {
			return diag.Errorf("failure updating secruedcookies: %v", err)
		}
	}

	return ResourceOutscaleLoadBalancerRead(ctx, d, meta)
}

func readResourceLb(ctx context.Context, client *osc.Client, elbName string, timeout time.Duration) (*osc.LoadBalancer, *osc.ReadLoadBalancersResponse, error) {
	filter := &osc.FiltersLoadBalancer{
		LoadBalancerNames: &[]string{elbName},
	}

	req := osc.ReadLoadBalancersRequest{
		Filters: filter,
	}

	resp, err := client.ReadLoadBalancers(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving load balancer: %w", err)
	}
	if resp == nil || len(*resp.LoadBalancers) == 0 {
		return nil, nil, nil
	}

	lb := (*resp.LoadBalancers)[0]
	return &lb, resp, nil
}

func ResourceOutscaleLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutRead)
	elbName := d.Id()

	lb, _, err := readResourceLb(ctx, client, elbName, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	if lb == nil {
		utils.LogManuallyDeleted("LoadBalancer", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("subregion_names", utils.StringSlicePtrToInterfaceSlice(&lb.SubregionNames))
	d.Set("dns_name", lb.DnsName)
	d.Set("health_check", flattenOAPIHealthCheck(&lb.HealthCheck))
	d.Set("access_log", flattenOAPIAccessLog(&lb.AccessLog))

	d.Set("backend_vm_ids", utils.StringSlicePtrToInterfaceSlice(&lb.BackendVmIds))
	d.Set("backend_ips", utils.StringSlicePtrToInterfaceSlice(&lb.BackendIps))
	if err := d.Set("listeners", flattenOAPIListeners(&lb.Listeners)); err != nil {
		return diag.FromErr(err)
	}
	d.Set("load_balancer_name", lb.LoadBalancerName)

	if lb.Tags != nil {
		ta := make([]map[string]interface{}, len(lb.Tags))
		for k1, v1 := range lb.Tags {
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
			len(lb.ApplicationStickyCookiePolicies))
		for k, v := range lb.ApplicationStickyCookiePolicies {
			a := make(map[string]interface{})
			a["cookie_name"] = v.CookieName
			a["policy_name"] = v.PolicyName
			app[k] = a
		}
		d.Set("application_sticky_cookie_policies", app)
	}
	if lb.LoadBalancerStickyCookiePolicies != nil {
		lbc := make([]map[string]interface{},
			len(lb.LoadBalancerStickyCookiePolicies))
		for k, v := range lb.LoadBalancerStickyCookiePolicies {
			a := make(map[string]interface{})
			a["policy_name"] = v.PolicyName
			lbc[k] = a
		}
		d.Set("load_balancer_sticky_cookie_policies", lbc)
	}

	d.Set("load_balancer_type", lb.LoadBalancerType)
	if lb.SecurityGroups != nil {
		d.Set("security_groups", utils.StringSlicePtrToInterfaceSlice(&lb.SecurityGroups))
	} else {
		d.Set("security_groups", make([]map[string]interface{}, 0))
	}

	d.Set("source_security_group", flattenSource_sg(&lb.SourceSecurityGroup))
	d.Set("subnets", utils.StringSlicePtrToInterfaceSlice(&lb.Subnets))
	d.Set("public_ip", ptr.From(lb.PublicIp))
	d.Set("secured_cookies", lb.SecuredCookies)

	d.Set("net_id", ptr.From(lb.NetId))

	return nil
}

func ResourceOutscaleLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutUpdate)
	var err error

	if d.HasChange("security_groups") {
		req := osc.UpdateLoadBalancerRequest{
			LoadBalancerName: d.Id(),
		}
		nSg, _ := d.GetOk("security_groups")
		req.SecurityGroups = utils.SetToStringSlicePtr(nSg.(*schema.Set))

		err := oapihelpers.RetryOnCodes(ctx, []string{"6031"}, func() (resp any, err error) {
			return client.UpdateLoadBalancer(ctx, req, options.WithRetryTimeout(timeout))
		}, timeout)
		if err != nil {
			return diag.Errorf("failure updating securitygroups: %v", err)
		}
	}

	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(*schema.Set)
		n := nraw.(*schema.Set)
		create := expandOAPITagsSDK(n)
		var remove []osc.ResourceLoadBalancerTag
		for _, t := range o.List() {
			tag := t.(map[string]interface{})
			s := tag["key"].(string)
			remove = append(remove,
				osc.ResourceLoadBalancerTag{
					Key: s,
				})
		}
		if len(remove) < 1 {
			goto skip_delete
		}

		_, err = client.DeleteLoadBalancerTags(ctx,
			osc.DeleteLoadBalancerTagsRequest{
				LoadBalancerNames: []string{d.Id()},
				Tags:              remove,
			}, options.WithRetryTimeout(timeout))
		if err != nil {
			return diag.FromErr(err)
		}

	skip_delete:
		if len(create) < 1 {
			goto skip_create
		}

		_, err = client.CreateLoadBalancerTags(ctx,
			osc.CreateLoadBalancerTagsRequest{
				LoadBalancerNames: []string{d.Id()},
				Tags:              create,
			}, options.WithRetryTimeout(timeout))
		if err != nil {
			return diag.FromErr(err)
		}

	skip_create:
	}

	if d.HasChange("listeners") {
		oldListeners, newListeners := d.GetChange("listeners")
		inter := oldListeners.(*schema.Set).Intersection(newListeners.(*schema.Set))
		lCreate := newListeners.(*schema.Set).Difference(inter)
		lRemoved := oldListeners.(*schema.Set).Difference(inter)
		var toRemove []*osc.Listener
		var toCreate []osc.ListenerForCreation
		var err error

		if lRemoved.Len() > 0 {
			toRemove, err = expandListeners(lRemoved.List())
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if lCreate.Len() > 0 {
			toCreate, err = expandListenerForCreation(lCreate.List())
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if len(toRemove) > 0 {
			ports := make([]int, 0, len(toRemove))
			for _, listener := range toRemove {
				ports = append(ports, listener.LoadBalancerPort)
			}

			req := osc.DeleteLoadBalancerListenersRequest{
				LoadBalancerName:  d.Id(),
				LoadBalancerPorts: ports,
			}

			log.Printf("[DEBUG] Load Balancer Delete Listeners")
			_, err := client.DeleteLoadBalancerListeners(ctx, req, options.WithRetryTimeout(timeout))
			if err != nil {
				return diag.Errorf("failure removing outdated load balancer listeners: %v", err)
			}
		}

		if len(toCreate) > 0 {
			req := osc.CreateLoadBalancerListenersRequest{
				LoadBalancerName: d.Id(),
				Listeners:        toCreate,
			}

			// Occasionally AWS will error with a 'duplicate listener', without any
			// other listeners on the Load Balancer. Retry here to eliminate that.
			_, err := client.CreateLoadBalancerListeners(ctx, req, options.WithRetryTimeout(timeout))
			if err != nil {
				return diag.Errorf("failure adding new or updated load balancer listeners: %v", err)
			}
		}
	}

	if d.HasChange("health_check") {
		hc := d.Get("health_check").([]interface{})
		if len(hc) > 0 {
			check := hc[0].(map[string]interface{})
			req := osc.UpdateLoadBalancerRequest{
				LoadBalancerName: d.Id(),
				HealthCheck: &osc.HealthCheck{
					HealthyThreshold:   check["healthy_threshold"].(int),
					UnhealthyThreshold: check["unhealthy_threshold"].(int),
					CheckInterval:      check["check_interval"].(int),
					Protocol:           check["protocol"].(string),
					Port:               check["port"].(int),
					Timeout:            check["timeout"].(int),
				},
			}
			if check["path"] != nil {
				p := check["path"].(string)
				req.HealthCheck.Path = &p
			}

			err := oapihelpers.RetryOnCodes(ctx, []string{"6031"}, func() (resp any, err error) {
				return client.UpdateLoadBalancer(ctx, req, options.WithRetryTimeout(timeout))
			}, timeout)
			if err != nil {
				return diag.Errorf("failure configuring health check for load balancer: %v", err)
			}
		}
	}

	if d.HasChange("access_log") {
		acg := d.Get("access_log").([]interface{})
		if len(acg) > 0 {

			aclg := acg[0].(map[string]interface{})
			isEnabled := aclg["is_enabled"].(bool)
			osuBucketName := aclg["osu_bucket_name"].(string)
			osuBucketPrefix := aclg["osu_bucket_prefix"].(string)
			publicationInterval := aclg["publication_interval"].(int)
			req := osc.UpdateLoadBalancerRequest{
				LoadBalancerName: d.Id(),
				AccessLog: &osc.AccessLog{
					IsEnabled:           isEnabled,
					OsuBucketName:       &osuBucketName,
					OsuBucketPrefix:     &osuBucketPrefix,
					PublicationInterval: &publicationInterval,
				},
			}

			err := oapihelpers.RetryOnCodes(ctx, []string{"6031"}, func() (resp any, err error) {
				return client.UpdateLoadBalancer(ctx, req, options.WithRetryTimeout(timeout))
			}, timeout)
			if err != nil {
				return diag.Errorf("failure configuring access log for load balancer: %v", err)
			}
		}
	}

	if d.HasChange("secured_cookies") {
		req := osc.UpdateLoadBalancerRequest{
			LoadBalancerName: d.Id(),
		}
		req.SecuredCookies = new(d.Get("secured_cookies").(bool))

		_, err := client.UpdateLoadBalancer(ctx, req, options.WithRetryTimeout(timeout))
		if err != nil {
			return diag.Errorf("failure updating secruedcookies: %v", err)
		}
	}

	return ResourceOutscaleLoadBalancerRead(ctx, d, meta)
}

func ResourceOutscaleLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutDelete)

	log.Printf("[INFO] Deleting Load Balancer: %s", d.Id())

	// Destroy the load balancer
	req := osc.DeleteLoadBalancerRequest{
		LoadBalancerName: d.Id(),
	}

	_, err := client.DeleteLoadBalancer(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error deleting load balancer: %v", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"ready"},
		Target:  []string{},
		Timeout: timeout,
		Refresh: func() (interface{}, string, error) {
			lb, _, _ := readResourceLb(ctx, client, d.Id(), timeout)
			if lb == nil {
				return nil, "", nil
			}
			return lb, "ready", nil
		},
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for load balancer (%s) to become null: %v", d.Id(), err)
	}

	// Remove this when bug will be fix
	if _, ok := d.GetOk("public_ip"); ok {
		time.Sleep(5 * time.Second)
	}

	return nil
}

func flattenOAPIHealthCheck(check *osc.HealthCheck) []map[string]interface{} {
	return []map[string]interface{}{{
		"healthy_threshold":   check.HealthyThreshold,
		"unhealthy_threshold": check.UnhealthyThreshold,
		"path":                check.Path,
		"check_interval":      check.CheckInterval,
		"port":                check.Port,
		"protocol":            check.Protocol,
		"timeout":             check.Timeout,
	}}
}

func flattenOAPIAccessLog(aclog *osc.AccessLog) []map[string]interface{} {
	return []map[string]interface{}{{
		"is_enabled":           aclog.IsEnabled,
		"osu_bucket_name":      aclog.OsuBucketName,
		"osu_bucket_prefix":    aclog.OsuBucketPrefix,
		"publication_interval": aclog.PublicationInterval,
	}}
}

func flattenSource_sg(ssg *osc.SourceSecurityGroup) []map[string]interface{} {
	return []map[string]interface{}{{
		"security_group_name":       ssg.SecurityGroupName,
		"security_group_account_id": ssg.SecurityGroupAccountId,
	}}
}
