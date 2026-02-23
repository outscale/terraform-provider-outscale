package oapi

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscaleSecurityGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleSecurityGroupRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"security_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"inbound_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port_range": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"security_groups_members": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_group_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"to_port_range": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ip_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_ranges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								// ValidateFunc: validateCIDRNetworkAddress,
							},
						},
						"prefix_list_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"outbound_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port_range": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"security_groups_members": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_group_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"to_port_range": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ip_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_ranges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								// ValidateFunc: validateCIDRNetworkAddress,
							},
						},
						"prefix_list_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": TagsSchemaComputedSDK(),
		},
	}
}

func DataSourceOutscaleSecurityGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadSecurityGroupsRequest{}

	filters, filtersOk := d.GetOk("filter")
	gn, gnOk := d.GetOk("security_group_name")
	gid, gidOk := d.GetOk("security_group_id")

	var filter osc.FiltersSecurityGroup
	if gnOk {
		filter.SecurityGroupNames = &[]string{gn.(string)}
		req.Filters = &filter
	}

	if gidOk {
		filter.SecurityGroupIds = &[]string{gid.(string)}
		req.Filters = &filter
	}

	var err error
	if filtersOk {
		req.Filters, err = buildOutscaleDataSourceSecurityGroupFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadSecurityGroups(ctx, req, options.WithRetryTimeout(5*time.Minute))

	var errString string
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidSecurityGroupID.NotFound") ||
			strings.Contains(fmt.Sprint(err), "InvalidGroup.NotFound") {
			resp.SecurityGroups = nil
			err = nil
		} else {
			// fmt.Printf("\n\nError on SGStateRefresh: %s", err)
			errString = err.Error()
		}

		return diag.Errorf("error on sgstaterefresh: %s", errString)
	}

	if resp.SecurityGroups == nil || len(*resp.SecurityGroups) == 0 {
		return diag.FromErr(ErrNoResults)
	}
	if len(*resp.SecurityGroups) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	sg := (*resp.SecurityGroups)[0]

	d.SetId(sg.SecurityGroupId)
	if err := d.Set("security_group_id", sg.SecurityGroupId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", sg.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("security_group_name", sg.SecurityGroupName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("net_id", ptr.From(sg.NetId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("account_id", sg.AccountId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(sg.Tags)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("inbound_rules", flattenOAPISecurityGroupRule(sg.InboundRules)); err != nil {
		return diag.FromErr(err)
	}
	return diag.FromErr(d.Set("outbound_rules", flattenOAPISecurityGroupRule(sg.OutboundRules)))
}

func buildOutscaleDataSourceSecurityGroupFilters(set *schema.Set) (*osc.FiltersSecurityGroup, error) {
	var filters osc.FiltersSecurityGroup
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "descriptions":
			filters.Descriptions = &filterValues
		case "inbound_rule_account_ids":
			filters.InboundRuleAccountIds = &filterValues
		case "inbound_rule_from_port_ranges":
			filters.InboundRuleFromPortRanges = new(utils.StringSliceToIntSlice(filterValues))
		case "inbound_rule_ip_ranges":
			filters.InboundRuleIpRanges = &filterValues
		case "inbound_rule_protocols":
			filters.InboundRuleProtocols = &filterValues
		case "inbound_rule_security_group_ids":
			filters.InboundRuleSecurityGroupIds = &filterValues
		case "inbound_rule_security_group_names":
			filters.InboundRuleSecurityGroupNames = &filterValues
		case "inbound_rule_to_port_ranges":
			filters.InboundRuleToPortRanges = new(utils.StringSliceToIntSlice(filterValues))
		case "net_ids":
			filters.NetIds = &filterValues
		case "outbound_rule_account_ids":
			filters.OutboundRuleAccountIds = &filterValues
		case "outbound_rule_from_port_ranges":
			filters.OutboundRuleFromPortRanges = new(utils.StringSliceToIntSlice(filterValues))
		case "outbound_rule_ip_ranges":
			filters.OutboundRuleIpRanges = &filterValues
		case "outbound_rule_protocols":
			filters.OutboundRuleProtocols = &filterValues
		case "outbound_rule_security_group_ids":
			filters.OutboundRuleSecurityGroupIds = &filterValues
		case "outbound_rule_security_group_names":
			filters.OutboundRuleSecurityGroupNames = &filterValues
		case "outbound_rule_to_port_ranges":
			filters.OutboundRuleToPortRanges = new(utils.StringSliceToIntSlice(filterValues))
		case "security_group_ids":
			filters.SecurityGroupIds = &filterValues
		case "security_group_names":
			filters.SecurityGroupNames = &filterValues
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		case "tags":
			filters.Tags = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
