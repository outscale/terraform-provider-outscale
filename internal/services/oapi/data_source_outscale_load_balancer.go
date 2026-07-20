package oapi

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/samber/lo"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func attrLBchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"subregion_names": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"load_balancer_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"load_balancer_type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"security_groups": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"subnets": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"source_security_group": lb_sg_schema(),
		"tags":                  TagsSchemaComputedSDK(),
		"dns_name": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
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
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
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
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"backend_ips": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"listeners": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			Elem: &schema.Resource{
				Schema: lb_listener_schema(true),
			},
		},
		"public_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"secured_cookies": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"net_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
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
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func getDataSourceSchemas(attrsSchema map[string]*schema.Schema) map[string]*schema.Schema {
	wholeSchema := map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
	}

	maps.Copy(wholeSchema, attrsSchema)

	return wholeSchema
}

func DataSourceOutscaleLoadBalancer() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleLoadBalancerRead,
		Schema:      getDataSourceSchemas(attrLBchema()),
	}
}

func buildOutscaleDataSourceLBFilters(set *schema.Set) (*osc.FiltersLoadBalancer, error) {
	filters := osc.FiltersLoadBalancer{}

	for _, v := range set.List() {
		m := v.(map[string]any)
		filterValues := lo.Map(m["values"].([]any), func(e any, _ int) string {
			return e.(string)
		})

		switch name := m["name"].(string); name {
		case "load_balancer_name":
			filters.LoadBalancerNames = &filterValues
		case "load_balancer_names":
			filters.LoadBalancerNames = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}

func readLbs(ctx context.Context, client *osc.Client, d *schema.ResourceData) (*osc.ReadLoadBalancersResponse, *string, error) {
	return readLbs_(ctx, client, d, schema.TypeString)
}

func readLbs_(ctx context.Context, client *osc.Client, d *schema.ResourceData, t schema.ValueType) (*osc.ReadLoadBalancersResponse, *string, error) {
	ename, nameOk := d.GetOk("load_balancer_name")
	filters, filtersOk := d.GetOk("filter")
	req := osc.ReadLoadBalancersRequest{
		Filters: &osc.FiltersLoadBalancer{},
	}

	if !nameOk && !filtersOk {
		return nil, nil, errors.New("one of filters, or load_balancer_name must be assigned")
	}

	var err error
	switch {
	case filtersOk:
		req.Filters, err = buildOutscaleDataSourceLBFilters(filters.(*schema.Set))
		if err != nil {
			return nil, nil, err
		}
	case t == schema.TypeString:
		req.Filters.LoadBalancerNames = &[]string{ename.(string)}
	default: /* assuming typelist */
		req.Filters = &osc.FiltersLoadBalancer{
			LoadBalancerNames: utils.InterfaceSliceToStringSlicePtr(ename.([]any)),
		}
	}
	elbName := (*req.Filters.LoadBalancerNames)[0]

	resp, err := client.ReadLoadBalancers(ctx, req, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving elb: %w", err)
	}
	return resp, &elbName, nil
}

func readLbs0(ctx context.Context, client *osc.Client, d *schema.ResourceData) (*osc.LoadBalancer, *osc.ReadLoadBalancersResponse, error) {
	resp, _, err := readLbs(ctx, client, d)
	if err != nil {
		return nil, nil, err
	}

	if resp.LoadBalancers == nil || len(*resp.LoadBalancers) == 0 {
		return nil, nil, ErrNoResults
	}
	if len(*resp.LoadBalancers) > 1 {
		return nil, nil, ErrMultipleResults
	}

	lbs := *resp.LoadBalancers
	return &lbs[0], resp, nil
}

func DataSourceOutscaleLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	lb, _, err := readLbs0(ctx, client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("subregion_names", utils.StringSlicePtrToInterfaceSlice(&lb.SubregionNames)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dns_name", lb.DnsName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("health_check", flattenOAPIHealthCheck(&lb.HealthCheck)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("access_log", flattenOAPIAccessLog(&lb.AccessLog)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("backend_vm_ids", utils.StringSlicePtrToInterfaceSlice(&lb.BackendVmIds)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("backend_ips", utils.StringSlicePtrToInterfaceSlice(&lb.BackendIps)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("listeners", flattenOAPIListeners(&lb.Listeners)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("load_balancer_name", lb.LoadBalancerName); err != nil {
		return diag.FromErr(err)
	}

	if lb.ApplicationStickyCookiePolicies != nil {
		app := make([]map[string]any, len(lb.ApplicationStickyCookiePolicies))
		for k, v := range lb.ApplicationStickyCookiePolicies {
			a := make(map[string]any)
			a["cookie_name"] = v.CookieName
			a["policy_name"] = v.PolicyName
			app[k] = a
		}
		if err := d.Set("application_sticky_cookie_policies", app); err != nil {
			return diag.FromErr(err)
		}
	} else {
		app := make([]map[string]any, 0)
		if err := d.Set("application_sticky_cookie_policies", app); err != nil {
			return diag.FromErr(err)
		}
	}
	if lb.LoadBalancerStickyCookiePolicies != nil {
		lbc := make([]map[string]any, len(lb.LoadBalancerStickyCookiePolicies))
		for k, v := range lb.LoadBalancerStickyCookiePolicies {
			a := make(map[string]any)
			a["policy_name"] = v.PolicyName
			lbc[k] = a
		}
		if err := d.Set("load_balancer_sticky_cookie_policies", lbc); err != nil {
			return diag.FromErr(err)
		}
	} else {
		lbc := make([]map[string]any, 0)
		if err := d.Set("load_balancer_sticky_cookie_policies", lbc); err != nil {
			return diag.FromErr(err)
		}
	}

	if lb.Tags != nil {
		ta := make([]map[string]any, len(lb.Tags))
		for k1, v1 := range lb.Tags {
			t := make(map[string]any)
			t["key"] = v1.Key
			t["value"] = v1.Value
			ta[k1] = t
		}
		if err := d.Set("tags", ta); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("tags", make([]map[string]any, 0)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("load_balancer_type", lb.LoadBalancerType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("security_groups", utils.StringSlicePtrToInterfaceSlice(&lb.SecurityGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("source_security_group", flattenSource_sg(&lb.SourceSecurityGroup)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("public_ip", ptr.From(lb.PublicIp)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("secured_cookies", lb.SecuredCookies); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("net_id", ptr.From(lb.NetId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("subnets", utils.StringSlicePtrToInterfaceSlice(&lb.Subnets)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", lb.State); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(lb.LoadBalancerName)

	return nil
}

// Legacy helper, will be removed once the datasource is migrated to the Plugin Framework
func lb_listener_schema(withPolicyNames bool) map[string]*schema.Schema {
	listenerSchema := map[string]*schema.Schema{
		"backend_port": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"backend_protocol": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"load_balancer_port": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"load_balancer_protocol": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"server_certificate_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
	if withPolicyNames {
		listenerSchema["policy_names"] = &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		}
	}

	return listenerSchema
}

func lb_sg_schema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"security_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		}},
	}
}

func flattenOAPIAccessLog(accessLog *osc.AccessLog) []map[string]any {
	if accessLog == nil {
		return nil
	}

	return []map[string]any{{
		"is_enabled":           accessLog.IsEnabled,
		"osu_bucket_name":      accessLog.OsuBucketName,
		"osu_bucket_prefix":    accessLog.OsuBucketPrefix,
		"publication_interval": accessLog.PublicationInterval,
	}}
}

func flattenOAPIHealthCheck(healthCheck *osc.HealthCheck) []map[string]any {
	if healthCheck == nil {
		return nil
	}

	return []map[string]any{{
		"healthy_threshold":   healthCheck.HealthyThreshold,
		"unhealthy_threshold": healthCheck.UnhealthyThreshold,
		"path":                healthCheck.Path,
		"check_interval":      healthCheck.CheckInterval,
		"port":                healthCheck.Port,
		"protocol":            healthCheck.Protocol,
		"timeout":             healthCheck.Timeout,
	}}
}

func flattenOAPIListeners(listeners *[]osc.Listener) []map[string]any {
	if listeners == nil {
		return nil
	}

	flattened := make([]map[string]any, len(*listeners))
	for i, listener := range *listeners {
		flattened[i] = map[string]any{
			"backend_port":           listener.BackendPort,
			"backend_protocol":       listener.BackendProtocol,
			"load_balancer_port":     listener.LoadBalancerPort,
			"load_balancer_protocol": listener.LoadBalancerProtocol,
			"server_certificate_id":  listener.ServerCertificateId,
			"policy_names":           listener.PolicyNames,
		}
	}

	return flattened
}

func flattenSource_sg(sourceSecurityGroup *osc.SourceSecurityGroup) []map[string]any {
	if sourceSecurityGroup == nil {
		return nil
	}

	return []map[string]any{{
		"security_group_name":       sourceSecurityGroup.SecurityGroupName,
		"security_group_account_id": sourceSecurityGroup.SecurityGroupAccountId,
	}}
}
