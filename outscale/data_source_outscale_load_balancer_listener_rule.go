package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			Required: true,
		},
		"path_pattern": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"priority": {
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

	lrName := lrNamei.(string)
	if !nameOk {
		return fmt.Errorf("listener_rule_name must be assigned")
	}

	filter := &oscgo.FiltersListenerRule{
		ListenerRuleNames: &[]string{lrName},
	}

	req := oscgo.ReadListenerRulesRequest{
		Filters: filter,
	}

	describeElbOpts := &oscgo.ReadListenerRulesOpts{
		ReadListenerRulesRequest: optional.NewInterface(req),
	}

	var resp oscgo.ReadListenerRulesResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.ListenerApi.ReadListenerRules(
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

	if lr.Priority != nil {
		d.Set("priority", lr.Priority)
	} else {
		fmt.Errorf("Malformated listener rule")
	}

	if lr.VmIds != nil {
		d.Set("vm_ids", flattenStringList(lr.VmIds))
	} else {
		fmt.Errorf("Malformated listener rule")
	}

	d.Set("request_id", resp.ResponseContext.RequestId)
	d.SetId(resource.UniqueId())

	return nil
}
