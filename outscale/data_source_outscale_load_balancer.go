package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/outscale/osc-sdk-go/osc"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func attrLBSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"subregion_name": {
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
		"security_groups_member": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"subnets_member": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"tag": tagsSchema(),

		"dns_name": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"health_check": {
			Type:     schema.TypeMap,
			Computed: true,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"healthy_threshold": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"unhealthy_threshold": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"checked_vm": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"check_interval": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"timeout": {
						Type:     schema.TypeString,
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
			Type:     schema.TypeSet,
			Computed: true,
			Optional: true,
			Elem: &schema.Resource{
				Schema: lb_listener_schema(),
			},
		},
		"firewall_rules_set_name": {
			Type:     schema.TypeMap,
			Computed: true,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"firewall_rules_set_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"account_alias": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"net_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"policies": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"application_sticky_cookie_policy": {
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
					"load_balancer_sticky_cookie_policy": {
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
		Schema: getDataSourceSchemas(attrLBSchema()),
	}
}

func buildOutscaleDataSourceLBFilters(set *schema.Set) *oscgo.FiltersLoadBalancer {
	filters := new(oscgo.FiltersLoadBalancer)

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
	ename, nameOk := d.GetOk("load_balancer_name")
	filters, filtersOk := d.GetOk("filter")
	filter := new(oscgo.FiltersLoadBalancer)

	if !nameOk && !filtersOk {
		return nil, nil, fmt.Errorf("One of filters, or listener_rule_name must be assigned")
	}

	if filtersOk {
		filter = buildOutscaleDataSourceLBFilters(filters.(*schema.Set))
	} else {
		elbName := ename.(string)
		filter = &oscgo.FiltersLoadBalancer{
			LoadBalancerNames: &[]string{elbName},
		}
	}
	elbName := (*filter.LoadBalancerNames)[0]

	req := oscgo.ReadLoadBalancersRequest{
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
			return nil, nil, fmt.Errorf("Unknow error")
		}

		return nil, nil, fmt.Errorf("Error retrieving ELB: %s", err)
	}
	return &resp, &elbName, nil
}

func readLbs0(conn *oscgo.APIClient, d *schema.ResourceData) (*oscgo.LoadBalancer, *oscgo.ReadLoadBalancersResponse, error) {
	resp, elbName, err := readLbs(conn, d)
	if err != nil {
		return nil, nil, err
	}
	lbs := *resp.LoadBalancers
	if len(lbs) != 1 {
		return nil, nil, fmt.Errorf("Unable to find LBU: %s", *elbName)
	}
	return &lbs[0], resp, nil
}

func dataSourceOutscaleOAPILoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	lb, resp, err := readLbs0(conn, d)
	if err != nil {
		return err
	}

	d.Set("subregion_name", flattenStringList(lb.SubregionNames))
	d.Set("dns_name", lb.DnsName)
	d.Set("health_check", flattenOAPIHealthCheck(nil, lb.HealthCheck))

	d.Set("backend_vm_ids", flattenStringList(lb.BackendVmIds))
	if lb.Listeners != nil {
		if err := d.Set("listeners", flattenOAPIListeners(lb.Listeners)); err != nil {
			return err
		}
	} else {
		if err := d.Set("listeners", make([]map[string]interface{}, 0)); err != nil {
			return err
		}
	}
	d.Set("load_balancer_name", lb.LoadBalancerName)

	policies := make(map[string]interface{})
	if lb.ApplicationStickyCookiePolicies != nil {
		app := make([]map[string]interface{}, len(*lb.ApplicationStickyCookiePolicies))
		for k, v := range *lb.ApplicationStickyCookiePolicies {
			a := make(map[string]interface{})
			a["cookie_name"] = v.CookieName
			a["policy_name"] = v.PolicyName
			app[k] = a
		}
		policies["application_sticky_cookie_policy"] = app
	}
	if lb.LoadBalancerStickyCookiePolicies != nil {
		lbc := make([]map[string]interface{}, len(*lb.LoadBalancerStickyCookiePolicies))
		for k, v := range *lb.LoadBalancerStickyCookiePolicies {
			a := make(map[string]interface{})
			a["policy_name"] = v.PolicyName
			lbc[k] = a
		}
		policies["load_balancer_sticky_cookie_policy"] = lbc
	} else {
		lbc := make([]map[string]interface{}, 0)
		policies["load_balancer_sticky_cookie_policy"] = lbc
	}
	d.Set("policies", policies)

	d.Set("load_balancer_type", lb.LoadBalancerType)
	if lb.SecurityGroups != nil {
		d.Set("security_groups_member", flattenStringList(lb.SecurityGroups))
	} else {
		d.Set("security_groups_member", make([]map[string]interface{}, 0))
	}
	ssg := make(map[string]string)
	if lb.SourceSecurityGroup != nil {
		ssg["security_group_account_id"] = *lb.SourceSecurityGroup.SecurityGroupAccountId
		ssg["security_group_name"] = *lb.SourceSecurityGroup.SecurityGroupName
	}
	d.Set("firewall_rules_set_name", ssg)
	d.Set("subnets_member", flattenStringList(lb.Subnets))
	d.Set("net_id", lb.NetId)
	d.Set("request_id", resp.ResponseContext.RequestId)
	d.SetId(*lb.LoadBalancerName)

	return nil
}
