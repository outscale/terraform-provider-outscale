package outscale

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
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
		"tags":                  tagsListOAPISchema2(true),
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
			Optional: true,
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

	for k, v := range attrsSchema {
		wholeSchema[k] = v
	}

	return wholeSchema

}

func dataSourceOutscaleOAPILoadBalancer() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPILoadBalancerRead,
		Schema: getDataSourceSchemas(attrLBchema()),
	}
}

func buildOutscaleDataSourceLBFilters(set *schema.Set) oscgo.FiltersLoadBalancer {
	filters := oscgo.FiltersLoadBalancer{}

	for _, v := range set.List() {
		m := v.(map[string]interface{})
		filterValues := make([]string, 0)
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "load_balancer_name":
			filters.LoadBalancerNames = &filterValues
		default:
			filters.LoadBalancerNames = &filterValues
			log.Printf("[Debug] Unknown Filter Name: %s. default to 'load_balancer_name'", name)
		}
	}
	return filters
}

func readLbs(conn *oscgo.APIClient, d *schema.ResourceData) (*oscgo.ReadLoadBalancersResponse, *string, error) {
	return readLbs_(conn, d, schema.TypeString)
}

func readLbs_(conn *oscgo.APIClient, d *schema.ResourceData, t schema.ValueType) (*oscgo.ReadLoadBalancersResponse, *string, error) {
	ename, nameOk := d.GetOk("load_balancer_name")
	filters, filtersOk := d.GetOk("filter")
	req := oscgo.ReadLoadBalancersRequest{
		Filters: &oscgo.FiltersLoadBalancer{},
	}

	if !nameOk && !filtersOk {
		return nil, nil, fmt.Errorf("One of filters, or load_balancer_name must be assigned")
	}

	if filtersOk {
		req.SetFilters(buildOutscaleDataSourceLBFilters(filters.(*schema.Set)))
	} else if t == schema.TypeString {
		req.Filters.SetLoadBalancerNames([]string{ename.(string)})
	} else { /* assuming typelist */
		req.Filters = &oscgo.FiltersLoadBalancer{
			LoadBalancerNames: utils.InterfaceSliceToStringSlicePtr(ename.([]interface{})),
		}
	}
	elbName := (*req.Filters.LoadBalancerNames)[0]
	var resp oscgo.ReadLoadBalancersResponse
	var err error
	var statusCode int
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.LoadBalancerApi.
			ReadLoadBalancers(context.Background()).
			ReadLoadBalancersRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		statusCode = httpResp.StatusCode
		return nil
	})

	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("")
			return nil, nil, fmt.Errorf("Loadbalancer Not Found")
		}

		return nil, nil, fmt.Errorf("Error retrieving ELB: %s", err)
	}
	return &resp, &elbName, nil
}

func readLbs0(conn *oscgo.APIClient, d *schema.ResourceData) (*oscgo.LoadBalancer, *oscgo.ReadLoadBalancersResponse, error) {
	resp, _, err := readLbs(conn, d)
	if err != nil {
		return nil, nil, err
	}

	if err := utils.IsResponseEmptyOrMutiple(len(resp.GetLoadBalancers()), "LoadBalancer"); err != nil {
		return nil, nil, err
	}

	lbs := *resp.LoadBalancers
	return &lbs[0], resp, nil
}

func dataSourceOutscaleOAPILoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	lb, _, err := readLbs0(conn, d)

	if err != nil {
		return err
	}

	d.Set("subregion_names", utils.StringSlicePtrToInterfaceSlice(lb.SubregionNames))
	d.Set("dns_name", lb.DnsName)
	d.Set("health_check", flattenOAPIHealthCheck(lb.HealthCheck))
	d.Set("access_log", flattenOAPIAccessLog(lb.AccessLog))

	d.Set("backend_vm_ids", utils.StringSlicePtrToInterfaceSlice(lb.BackendVmIds))
	if err := d.Set("listeners", flattenOAPIListeners(lb.Listeners)); err != nil {
		return err
	}
	d.Set("load_balancer_name", lb.LoadBalancerName)

	if lb.ApplicationStickyCookiePolicies != nil {
		app := make([]map[string]interface{}, len(*lb.ApplicationStickyCookiePolicies))
		for k, v := range *lb.ApplicationStickyCookiePolicies {
			a := make(map[string]interface{})
			a["cookie_name"] = v.CookieName
			a["policy_name"] = v.PolicyName
			app[k] = a
		}
		d.Set("application_sticky_cookie_policies", app)
	} else {
		app := make([]map[string]interface{}, 0)
		d.Set("application_sticky_cookie_policies", app)
	}
	if lb.LoadBalancerStickyCookiePolicies != nil {
		lbc := make([]map[string]interface{}, len(*lb.LoadBalancerStickyCookiePolicies))
		for k, v := range *lb.LoadBalancerStickyCookiePolicies {
			a := make(map[string]interface{})
			a["policy_name"] = v.PolicyName
			lbc[k] = a
		}
		d.Set("load_balancer_sticky_cookie_policies", lbc)
	} else {
		lbc := make([]map[string]interface{}, 0)
		d.Set("load_balancer_sticky_cookie_policies", lbc)
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
	d.Set("load_balancer_type", lb.LoadBalancerType)
	if lb.SecurityGroups != nil {
		d.Set("security_groups", utils.StringSlicePtrToInterfaceSlice(lb.SecurityGroups))
	} else {
		d.Set("security_groups", make([]map[string]interface{}, 0))
	}

	if lb.SourceSecurityGroup != nil {
		d.Set("source_security_group", flattenSource_sg(lb.SourceSecurityGroup))
	} else {
		d.Set("source_security_group", make([]map[string]interface{}, 0))
	}

	d.Set("public_ip", lb.PublicIp)
	d.Set("secured_cookies", lb.SecuredCookies)
	d.Set("net_id", lb.NetId)
	d.Set("subnets", utils.StringSlicePtrToInterfaceSlice(lb.Subnets))
	d.SetId(*lb.LoadBalancerName)

	return nil
}
