package oapi

import (
	"context"
	"log"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleLoadBalancerListenerRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleLoadBalancerListenerRuleCreate,
		ReadContext:   ResourceOutscaleLoadBalancerListenerRuleRead,
		UpdateContext: ResourceOutscaleLoadBalancerListenerRuleUpdate,
		DeleteContext: ResourceOutscaleLoadBalancerListenerRuleDelete,
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
			"vm_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"listener": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"load_balancer_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listener_rule": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"host_name_pattern": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"listener_rule_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"listener_rule_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"listener_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"path_pattern": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func ResourceOutscaleLoadBalancerListenerRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutCreate)

	req := &osc.CreateListenerRuleRequest{}

	if vids, ok := d.GetOk("vm_ids"); ok {
		req.VmIds = utils.SetToStringSlice(vids.(*schema.Set))
	} else {
		return diag.Errorf("expect vm_ids")
	}

	if li, lok := d.GetOk("listener"); lok {
		ls := li.([]interface{})
		l := ls[0].(map[string]interface{})
		ll := osc.LoadBalancerLight{}
		if l["load_balancer_name"] == nil || l["load_balancer_port"] == nil {
			return diag.Errorf("listener missing argument ")
		}
		lbpii := l["load_balancer_port"].(int)
		ll.LoadBalancerName = l["load_balancer_name"].(string)
		ll.LoadBalancerPort = lbpii
		req.Listener = ll
	} else {
		return diag.Errorf("expect listener")
	}

	if lri, lok := d.GetOk("listener_rule"); lok {
		lrs := lri.([]interface{})
		lr := lrs[0].(map[string]interface{})

		lrfc := osc.ListenerRuleForCreation{}
		if lr["priority"] == nil {
			return diag.Errorf("listener priority argument missing")
		}
		if lr["action"] != nil {
			lrfc.Action = new(lr["action"].(string))
		}
		if lr["path_pattern"] != nil {
			lrfc.PathPattern = new(lr["path_pattern"].(string))
		}
		if lr["host_name_pattern"] != nil {
			lrfc.HostNamePattern = new(lr["host_name_pattern"].(string))
		}
		if lr["listener_rule_name"] != nil {
			lrfc.ListenerRuleName = lr["listener_rule_name"].(string)
		}
		p := lr["priority"].(int)
		lrfc.Priority = p
		req.ListenerRule = lrfc
	} else {
		return diag.Errorf("expect listener rule")
	}

	resp, err := client.CreateListenerRule(ctx, *req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ptr.From(resp.ListenerRule.ListenerRuleName))

	return ResourceOutscaleLoadBalancerListenerRuleRead(ctx, d, meta)
}

func ResourceOutscaleLoadBalancerListenerRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutRead)
	lrName := d.Id()

	filter := &osc.FiltersListenerRule{
		ListenerRuleNames: &[]string{lrName},
	}

	req := osc.ReadListenerRulesRequest{
		Filters: filter,
	}

	resp, err := client.ReadListenerRules(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.ListenerRules == nil || utils.IsResponseEmpty(len(*resp.ListenerRules), "LoadBalancerListenerRule", d.Id()) {
		d.SetId("")
		return nil
	}
	lr := (*resp.ListenerRules)[0]
	lrsl := make([]interface{}, 1)
	lrs := make(map[string]interface{})

	if lr.Action != nil {
		lrs["action"] = lr.Action
	}
	if lr.HostNamePattern != nil {
		lrs["host_name_pattern"] = lr.HostNamePattern
	}
	if lr.ListenerRuleName != nil {
		lrs["listener_rule_name"] = lr.ListenerRuleName
	}
	if lr.ListenerRuleId != nil {
		lrs["listener_rule_id"] = lr.ListenerRuleId
	}
	if lr.ListenerId != nil {
		lrs["listener_id"] = lr.ListenerId
	}
	if lr.PathPattern != nil {
		lrs["path_pattern"] = lr.PathPattern
	}
	if lr.Priority != nil {
		lrs["priority"] = lr.Priority
	}
	lrsl[0] = lrs
	err = d.Set("listener_rule", lrsl)
	if err != nil {
		return diag.FromErr(err)
	}
	if lr.VmIds != nil {
		d.Set("vm_ids", utils.StringSlicePtrToInterfaceSlice(lr.VmIds))
	}
	return nil
}

func ResourceOutscaleLoadBalancerListenerRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutUpdate)

	if d.HasChange("listener_rule") {
		var err error
		nw := d.Get("listener_rule").([]interface{})
		if len(nw) != 1 {
			return diag.Errorf("error multiple listener_rule matched or empty: %s", err)
		}
		check := nw[0].(map[string]interface{})
		req := osc.UpdateListenerRuleRequest{
			ListenerRuleName: d.Id(),
		}
		if check["host_name_pattern"] != nil {
			req.HostPattern = new(check["host_name_pattern"].(string))
		} else {
			req.HostPattern = new("")
		}
		if check["listener_rule_name"] != nil {
			req.ListenerRuleName = check["listener_rule_name"].(string)
		} else {
			req.ListenerRuleName = ""
		}
		if check["path_pattern"] != nil {
			req.PathPattern = new(check["path_pattern"].(string))
		} else {
			req.PathPattern = new("")
		}

		_, err = client.UpdateListenerRule(ctx, req, options.WithRetryTimeout(timeout))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return ResourceOutscaleLoadBalancerListenerRuleRead(ctx, d, meta)
}

func ResourceOutscaleLoadBalancerListenerRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutDelete)

	log.Printf("[INFO] Deleting Listener Rule: %s", d.Id())

	// Destroy the listener rule
	req := osc.DeleteListenerRuleRequest{
		ListenerRuleName: d.Id(),
	}

	_, err := client.DeleteListenerRule(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error deleting listener rule: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"ready"},
		Target:  []string{},
		Timeout: timeout,
		Refresh: func() (interface{}, string, error) {
			filter := &osc.FiltersListenerRule{
				ListenerRuleNames: &[]string{d.Id()},
			}

			req := osc.ReadListenerRulesRequest{
				Filters: filter,
			}

			resp, err := client.ReadListenerRules(ctx, req, options.WithRetryTimeout(timeout))

			if err != nil || len(*resp.ListenerRules) < 1 {
				return nil, "", nil
			}

			return &(*resp.ListenerRules)[0], "ready", nil
		},
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for listener rule (%s) to become nil: %s", d.Id(), err)
	}

	return nil
}
