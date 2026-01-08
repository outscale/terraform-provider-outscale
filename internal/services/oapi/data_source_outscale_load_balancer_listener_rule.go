package oapi

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
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

func DataSourceOutscaleLoadBalancerLDRule() *schema.Resource {
	return &schema.Resource{
		Read:   DataSourceOutscaleLoadBalancerLDRuleRead,
		Schema: getDataSourceSchemas(attrLBListenerRule()),
	}
}

func DataSourceOutscaleLoadBalancerLDRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

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
			case "listener_rule_names":
				filter.ListenerRuleNames = &filterValues
			default:
				return utils.UnknownDataSourceFilterError(context.Background(), name)
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
	var err = retry.Retry(5*time.Minute, func() *retry.RetryError {
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
		return ErrNoResults
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
	}

	if lr.VmIds != nil {
		d.Set("vm_ids", utils.StringSlicePtrToInterfaceSlice(lr.VmIds))
	}

	d.SetId(id.UniqueId())

	return nil
}
