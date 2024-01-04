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

func resourceOutscaleLoadBalancerListenerRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleLoadBalancerListenerRuleCreate,
		Read:   resourceOutscaleLoadBalancerListenerRuleRead,
		Update: resourceOutscaleLoadBalancerListenerRuleUpdate,
		Delete: resourceOutscaleLoadBalancerListenerRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceOutscaleLoadBalancerListenerRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := &oscgo.CreateListenerRuleRequest{}

	if vids, ok := d.GetOk("vm_ids"); ok {
		req.SetVmIds(utils.SetToStringSlice(vids.(*schema.Set)))
	} else {
		return fmt.Errorf("expect vm_ids")
	}

	if li, lok := d.GetOk("listener"); lok {
		ls := li.([]interface{})
		l := ls[0].(map[string]interface{})
		ll := oscgo.LoadBalancerLight{}
		if l["load_balancer_name"] == nil || l["load_balancer_port"] == nil {
			return fmt.Errorf("listener missing argument ")
		}
		lbpii := l["load_balancer_port"].(int)
		lbpi := int32(lbpii)
		ll.SetLoadBalancerName(l["load_balancer_name"].(string))
		ll.SetLoadBalancerPort(lbpi)
		req.SetListener(ll)
	} else {
		return fmt.Errorf("expect listener")
	}

	if lri, lok := d.GetOk("listener_rule"); lok {
		lrs := lri.([]interface{})
		lr := lrs[0].(map[string]interface{})

		lrfc := oscgo.ListenerRuleForCreation{}
		if lr["priority"] == nil {
			return fmt.Errorf("listener priority argument missing")
		}
		if lr["action"] != nil {
			lrfc.SetAction(lr["action"].(string))
		}
		if lr["path_pattern"] != nil {
			lrfc.SetPathPattern(lr["path_pattern"].(string))
		}
		if lr["host_name_pattern"] != nil {
			lrfc.SetHostNamePattern(lr["host_name_pattern"].(string))
		}
		if lr["listener_rule_name"] != nil {
			lrfc.SetListenerRuleName(lr["listener_rule_name"].(string))
		}
		p := lr["priority"].(int)
		lrfc.SetPriority(int32(p))
		req.SetListenerRule(lrfc)
	} else {
		return fmt.Errorf("expect listener rule")
	}

	var err error
	var resp oscgo.CreateListenerRuleResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.ListenerApi.CreateListenerRule(
			context.Background()).CreateListenerRuleRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(*resp.ListenerRule.ListenerRuleName)

	return resourceOutscaleLoadBalancerListenerRuleRead(d, meta)
}

func resourceOutscaleLoadBalancerListenerRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	lrName := d.Id()

	filter := &oscgo.FiltersListenerRule{
		ListenerRuleNames: &[]string{lrName},
	}

	req := oscgo.ReadListenerRulesRequest{
		Filters: filter,
	}

	var resp oscgo.ReadListenerRulesResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.ListenerApi.ReadListenerRules(
			context.Background()).ReadListenerRulesRequest(req).
			Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}
	if utils.IsResponseEmpty(len(resp.GetListenerRules()), "LoadBalancerListenerRule", d.Id()) {
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
		return err
	}
	if lr.VmIds != nil {
		d.Set("vm_ids", utils.StringSlicePtrToInterfaceSlice(lr.VmIds))
	}
	return nil
}

func resourceOutscaleLoadBalancerListenerRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	if d.HasChange("listener_rule") {
		var err error
		nw := d.Get("listener_rule").([]interface{})
		if len(nw) != 1 {
			return fmt.Errorf("Error Multiple listener_rule matched or empty: %s", err)
		}
		check := nw[0].(map[string]interface{})
		req := oscgo.UpdateListenerRuleRequest{
			ListenerRuleName: d.Id(),
		}
		if check["host_name_pattern"] != nil {
			req.SetHostPattern(check["host_name_pattern"].(string))
		} else {
			req.SetHostPattern("")
		}
		if check["listener_rule_name"] != nil {
			req.SetListenerRuleName(check["listener_rule_name"].(string))
		} else {
			req.SetListenerRuleName("")
		}
		if check["path_pattern"] != nil {
			req.SetPathPattern(check["path_pattern"].(string))
		} else {
			req.SetPathPattern("")
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, httpResp, err := conn.ListenerApi.UpdateListenerRule(
				context.Background()).UpdateListenerRuleRequest(req).
				Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})

		if err != nil {
			return err
		}
	}
	return resourceOutscaleLoadBalancerListenerRuleRead(d, meta)
}

func resourceOutscaleLoadBalancerListenerRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	log.Printf("[INFO] Deleting Listener Rule: %s", d.Id())

	// Destroy the listener rule
	req := oscgo.DeleteListenerRuleRequest{
		ListenerRuleName: d.Id(),
	}

	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.ListenerApi.DeleteListenerRule(
			context.Background()).DeleteListenerRuleRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting listener rule: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"ready"},
		Target:  []string{},
		Refresh: func() (interface{}, string, error) {
			filter := &oscgo.FiltersListenerRule{
				ListenerRuleNames: &[]string{d.Id()},
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

			if err != nil || len(*resp.ListenerRules) < 1 {
				return nil, "", nil
			}

			return &(*resp.ListenerRules)[0], "ready", nil
		},
		Timeout:    5 * time.Minute,
		MinTimeout: 10 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for listener rule (%s) to become nil: %s", d.Id(), err)
	}

	return nil
}
