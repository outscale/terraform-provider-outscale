package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
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
		ReadContext: DataSourceOutscaleLoadBalancerLDRuleRead,
		Schema:      getDataSourceSchemas(attrLBListenerRule()),
	}
}

func DataSourceOutscaleLoadBalancerLDRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	lrNamei, nameOk := d.GetOk("listener_rule_name")
	filters, filtersOk := d.GetOk("filter")
	filter := &osc.FiltersListenerRule{}

	if !nameOk && !filtersOk {
		return diag.Errorf("listener_rule_name must be assigned")
	}

	if filtersOk {
		set := filters.(*schema.Set)

		if set.Len() < 1 {
			return diag.Errorf("filter can't be empty")
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
				return diag.FromErr(utils.UnknownDataSourceFilterError(name))
			}
		}
	} else {
		filter = &osc.FiltersListenerRule{
			ListenerRuleNames: &[]string{lrNamei.(string)},
		}
	}

	req := osc.ReadListenerRulesRequest{
		Filters: filter,
	}

	resp, err := client.ReadListenerRules(ctx, req, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.ListenerRules == nil || len(*resp.ListenerRules) < 1 {
		return diag.FromErr(ErrNoResults)
	}
	lr := (*resp.ListenerRules)[0]
	if lr.Action != nil {
		d.Set("action", ptr.From(lr.Action))
	}
	if lr.HostNamePattern != nil {
		d.Set("host_name_pattern", ptr.From(lr.HostNamePattern))
	}
	if lr.ListenerRuleName != nil {
		d.Set("listener_rule_name", ptr.From(lr.ListenerRuleName))
	}
	if lr.PathPattern != nil {
		d.Set("path_pattern", ptr.From(lr.PathPattern))
	}

	if lr.ListenerRuleId != nil {
		d.Set("listener_rule_id", ptr.From(lr.ListenerRuleId))
	}
	if lr.ListenerId != nil {
		d.Set("listener_id", ptr.From(lr.ListenerId))
	}

	if lr.Priority != nil {
		d.Set("priority", ptr.From(lr.Priority))
	}

	if lr.VmIds != nil {
		d.Set("vm_ids", utils.StringSlicePtrToInterfaceSlice(lr.VmIds))
	}

	d.SetId(id.UniqueId())

	return nil
}
