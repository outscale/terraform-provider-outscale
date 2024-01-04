package outscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/spf13/cast"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleOAPIApiAccessRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIApiAccessRuleCreate,
		Read:   resourceOutscaleOAPIApiAccessRuleRead,
		Update: resourceOutscaleOAPIApiAccessRuleUpdate,
		Delete: resourceOutscaleOAPIApiAccessRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"api_access_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ca_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cns": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_ranges": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPIApiAccessRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	var checkParam = false
	req := oscgo.CreateApiAccessRuleRequest{}

	if val, ok := d.GetOk("ca_ids"); ok {
		checkParam = true
		req.CaIds = utils.SetToStringSlicePtr(val.(*schema.Set))
	}
	if val, ok := d.GetOk("ip_ranges"); ok {
		checkParam = true
		req.IpRanges = utils.SetToStringSlicePtr(val.(*schema.Set))
	}
	if !checkParam {
		return fmt.Errorf("[DEBUG] Error 'ca_ids' or 'ip_ranges' field is require for API Access Rules creation")
	}

	if val, ok := d.GetOk("cns"); ok {
		req.Cns = utils.SetToStringSlicePtr(val.(*schema.Set))
	}
	if v, ok := d.GetOk("description"); ok {
		req.SetDescription(v.(string))
	}

	var resp oscgo.CreateApiAccessRuleResponse
	var err error
	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.ApiAccessRuleApi.CreateApiAccessRule(context.Background()).CreateApiAccessRuleRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(cast.ToString(resp.ApiAccessRule.GetApiAccessRuleId()))

	return resourceOutscaleOAPIApiAccessRuleRead(d, meta)
}

func resourceOutscaleOAPIApiAccessRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadApiAccessRulesRequest{
		Filters: &oscgo.FiltersApiAccessRule{ApiAccessRuleIds: &[]string{d.Id()}},
	}

	var resp oscgo.ReadApiAccessRulesResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.ApiAccessRuleApi.ReadApiAccessRules(context.Background()).ReadApiAccessRulesRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error reading api access rule id (%s)", utils.GetErrorResponse(err))
	}
	if !resp.HasApiAccessRules() {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}
	if utils.IsResponseEmpty(len(resp.GetApiAccessRules()), "ApiAccessRule", d.Id()) {
		d.SetId("")
		return nil
	}

	accRule := resp.GetApiAccessRules()[0]

	if err := d.Set("api_access_rule_id", accRule.GetApiAccessRuleId()); err != nil {
		return err
	}
	if accRule.HasCaIds() {
		if err := d.Set("ca_ids", accRule.GetCaIds()); err != nil {
			return err
		}
	}

	if accRule.HasCns() {
		if err := d.Set("cns", accRule.GetCns()); err != nil {
			return err
		}
	}
	if accRule.HasIpRanges() {
		if err := d.Set("ip_ranges", accRule.GetIpRanges()); err != nil {
			return err
		}
	}
	if accRule.HasDescription() {
		if err := d.Set("description", accRule.GetDescription()); err != nil {
			return err
		}
	}
	return nil
}

func resourceOutscaleOAPIApiAccessRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	accRid, isIdOk := d.GetOk("api_access_rule_id")
	if !isIdOk {
		return fmt.Errorf("[DEBUG] Error 'api_access_rule_id' field is required to update API Access Rules")
	}

	var checkParam = false
	req := oscgo.UpdateApiAccessRuleRequest{
		ApiAccessRuleId: accRid.(string),
	}

	if val, ok := d.GetOk("ca_ids"); ok {
		checkParam = true
		req.CaIds = utils.SetToStringSlicePtr(val.(*schema.Set))
	}
	if val, ok := d.GetOk("ip_ranges"); ok {
		checkParam = true
		req.IpRanges = utils.SetToStringSlicePtr(val.(*schema.Set))
	}

	if !checkParam {
		return fmt.Errorf("[DEBUG] Error 'ca_ids' or 'ip_ranges' field is require to update API Access Rules")
	}

	if val, ok := d.GetOk("cns"); ok {
		req.Cns = utils.SetToStringSlicePtr(val.(*schema.Set))
	}
	if v, ok := d.GetOk("description"); ok {
		req.SetDescription(v.(string))
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.ApiAccessRuleApi.UpdateApiAccessRule(context.Background()).UpdateApiAccessRuleRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return resourceOutscaleOAPIApiAccessRuleRead(d, meta)
}

func resourceOutscaleOAPIApiAccessRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.DeleteApiAccessRuleRequest{
		ApiAccessRuleId: d.Id(),
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.ApiAccessRuleApi.DeleteApiAccessRule(context.Background()).DeleteApiAccessRuleRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	return err
}

func getParameters(d *schema.ResourceData, param string) []string {
	_, val := d.GetChange(param)
	m := val.([]interface{})
	a := make([]string, len(m))
	for k, v := range m {
		a[k] = v.(string)
	}
	return a
}
