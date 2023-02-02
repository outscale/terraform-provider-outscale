package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func attrLBListenerRule() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
		"action": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"host_name_pattern": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"listener_rule_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"path_pattern": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"priority": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"listener_rule_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"listener_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"vm_ids": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceOutscaleOAPILoadBalancerLDRule() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPILoadBalancerLDRuleRead,
		Schema: getDataSourceSchemas(attrLBListenerRule()),
	}
}

func dataSourceOutscaleOAPILoadBalancerLDRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.ReadListenerRulesRequest{}
	if filters, filtersOk := d.GetOk("filter"); filtersOk {
		req.SetFilters(buildOutscaleOAPILoadBalancerListenerRuleDataSourceFilters(filters.(*schema.Set)))
	}

	var resp oscgo.ReadListenerRulesResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.ListenerApi.ReadListenerRules(
			context.Background()).
			ReadListenerRulesRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	if err := utils.IsResponseEmptyOrMutiple(len(resp.GetListenerRules()), "Listener Rule"); err != nil {
		return err
	}

	lr := (*resp.ListenerRules)[0]
	if lr.Action != nil {
		d.Set("action", lr.Action)
	}
	if lr.HostNamePattern != nil {
		d.Set("host_name_pattern", lr.HostNamePattern)
	}
	if lr.ListenerRuleName != nil {
		d.Set("listener_rule_name", lr.ListenerRuleName)
	}
	if lr.PathPattern != nil {
		d.Set("path_pattern", lr.PathPattern)
	}
	if lr.ListenerRuleId != nil {
		d.Set("listener_rule_id", lr.ListenerRuleId)
	}
	if lr.ListenerId != nil {
		d.Set("listener_id", lr.ListenerId)
	}
	if lr.Priority != nil {
		d.Set("priority", lr.Priority)
	} else {
		fmt.Errorf("Malformated listener rule")
	}
	if lr.VmIds != nil {
		d.Set("vm_ids", utils.StringSlicePtrToInterfaceSlice(lr.VmIds))
	} else {
		fmt.Errorf("Malformated listener rule")
	}

	d.SetId(resource.UniqueId())

	return nil
}

func buildOutscaleOAPILoadBalancerListenerRuleDataSourceFilters(set *schema.Set) oscgo.FiltersListenerRule {
	var filters oscgo.FiltersListenerRule
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}
		switch name := m["name"].(string); name {
		case "listener_rule_name":
			filters.ListenerRuleNames = &filterValues
		default:
			log.Printf("[Debug] Unknown Filter Name: %s. default to 'load_balancer_name'", name)
		}
	}
	return filters
}
