package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func attrLBListenerRule() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
			Optional: true,
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

	lrNamei, nameOk := d.GetOk("listener_rule_name")
	filters, filtersOk := d.GetOk("filter")
	filter := &oscgo.FiltersListenerRule{}

	if !nameOk && !filtersOk {
		return fmt.Errorf("listener_rule_name must be assigned")
	}

	if filtersOk {
		set := filters.(*schema.Set)

		if set.Len() < 1 {
			return fmt.Errorf("filter can't be empty")
		}
		for _, v := range set.List() {
			m := v.(map[string]interface{})
			filterValues := make([]string, 0)
			for _, e := range m["values"].([]interface{}) {
				filterValues = append(filterValues, e.(string))
			}

			switch name := m["name"].(string); name {
			case "listener_rule_name":
				filter.ListenerRuleNames = &filterValues
			default:
				filter.ListenerRuleNames = &filterValues
				log.Printf("[Debug] Unknown Filter Name: %s. default to 'load_balancer_name'", name)
			}
		}
	} else {
		filter = &oscgo.FiltersListenerRule{
			ListenerRuleNames: &[]string{lrNamei.(string)},
		}
	}

	req := oscgo.ReadListenerRulesRequest{
		Filters: filter,
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

	if len(*resp.ListenerRules) < 1 {
		return fmt.Errorf("can't find listener rule")
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
