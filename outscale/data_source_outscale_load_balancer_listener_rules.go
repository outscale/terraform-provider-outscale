package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func attrLBListenerRules() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"listener_rules": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
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
				},
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceOutscaleOAPILoadBalancerLDRules() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPILoadBalancerLDRulesRead,
		Schema: getDataSourceSchemas(attrLBListenerRules()),
	}
}

func dataSourceOutscaleOAPILoadBalancerLDRulesRead(d *schema.ResourceData, meta interface{}) error {
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

	result := *resp.ListenerRules
	result_len := len(result)
	if result_len == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	lrs_ret := make([]map[string]interface{}, result_len)
	for k, lr := range result {
		l := make(map[string]interface{})
		if lr.Action != nil {
			l["action"] = lr.Action
		}
		if lr.HostNamePattern != nil {
			l["host_name_pattern"] = lr.HostNamePattern
		}
		if lr.ListenerRuleName != nil {
			l["listener_rule_name"] = lr.ListenerRuleName
		}
		if lr.PathPattern != nil {
			l["path_pattern"] = lr.PathPattern
		}

		if lr.ListenerRuleId != nil {
			l["listener_rule_id"] = lr.ListenerRuleId
		}
		if lr.ListenerId != nil {
			l["listener_id"] = lr.ListenerId
		}

		if lr.Priority != nil {
			l["priority"] = lr.Priority
		} else {
			fmt.Errorf("Malformated listener rule")
		}

		if lr.VmIds != nil {
			l["vm_ids"] = utils.StringSlicePtrToInterfaceSlice(lr.VmIds)
		} else {
			fmt.Errorf("Malformated listener rule")
		}
		lrs_ret[k] = l
	}

	d.Set("listener_rules", lrs_ret)
	d.SetId(resource.UniqueId())

	return nil
}
