package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleOAPISecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPISecurityGroupRead,

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
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
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
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
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
			"tags": tagsOAPIListSchemaComputed(),
		},
	}
}

func dataSourceOutscaleOAPISecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI
	req := &oapi.ReadSecurityGroupsRequest{}

	filters, filtersOk := d.GetOk("filter")
	gn, gnOk := d.GetOk("security_group_name")
	gid, gidOk := d.GetOk("security_group_id")

	if gnOk {
		req.Filters.SecurityGroupNames = []string{gn.(string)}
	}

	if gidOk {
		req.Filters.SecurityGroupIds = []string{gid.(string)}
	}

	if filtersOk {
		req.Filters = buildOutscaleOAPIDataSourceSecurityGroupFilters(filters.(*schema.Set))
		fmt.Printf("%+v\n", req)
	}

	var err error
	var resp *oapi.POST_ReadSecurityGroupsResponses
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_ReadSecurityGroups(*req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	var errString string

	if err != nil || resp.OK == nil {
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidSecurityGroupID.NotFound") ||
				strings.Contains(fmt.Sprint(err), "InvalidGroup.NotFound") {
				resp = nil
				err = nil
			} else {
				//fmt.Printf("\n\nError on SGStateRefresh: %s", err)
				errString = err.Error()
			}

		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
		}

		return fmt.Errorf("Error on SGStateRefresh: %s", errString)
	}

	result := resp.OK

	if resp == nil || len(result.SecurityGroups) == 0 {
		return fmt.Errorf("Unable to find Security Group")
	}

	if len(result.SecurityGroups) > 1 {
		return fmt.Errorf("multiple results returned, please use a more specific criteria in your query")
	}

	sg := result.SecurityGroups[0]

	d.SetId(sg.SecurityGroupId)
	d.Set("security_group_id", sg.SecurityGroupId)
	d.Set("description", sg.Description)
	d.Set("security_group_name", sg.SecurityGroupName)
	d.Set("net_id", sg.NetId)
	d.Set("account_id", sg.AccountId)
	d.Set("tags", tagsOAPIToMap(sg.Tags))
	d.Set("inbound_rules", flattenOAPISecurityGroupRule(sg.InboundRules))
	d.Set("request_id", result.ResponseContext.RequestId)
	return d.Set("outbound_rules", flattenOAPISecurityGroupRule(sg.OutboundRules))
}

func flattenOAPIIPPermissions(p []*fcu.IpPermission) []map[string]interface{} {
	ips := make([]map[string]interface{}, len(p))

	for k, v := range p {
		ip := make(map[string]interface{})
		ip["from_port_range"] = v.FromPort
		ip["ip_protocol"] = v.IpProtocol
		ip["to_port_range"] = v.ToPort

		ipr := make([]map[string]interface{}, len(v.IpRanges))
		for i, v := range v.IpRanges {
			ipr[i] = map[string]interface{}{"cidr_ip": v.CidrIp}
		}
		ip["ip_ranges"] = ipr

		prx := make([]map[string]interface{}, len(v.PrefixListIds))
		for i, v := range v.PrefixListIds {
			prx[i] = map[string]interface{}{"prefix_list_id": v.PrefixListId}
		}
		ip["prefix_list_ids"] = prx

		grp := make([]map[string]interface{}, len(v.UserIdGroupPairs))
		for i, v := range v.UserIdGroupPairs {
			grp[i] = map[string]interface{}{
				"account_id":          v.UserId,
				"security_group_name": v.GroupName,
				"security_group_id":   v.GroupId,
			}
		}
		ip["security_groups_members"] = grp

		ips[k] = ip
	}

	return ips
}

func buildOutscaleOAPIDataSourceSecurityGroupFilters(set *schema.Set) oapi.FiltersSecurityGroup {
	var filters oapi.FiltersSecurityGroup
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "account_ids":
			filters.AccountIds = filterValues
		case "descriptions":
			filters.Descriptions = filterValues
		case "inbound_rule_account_ids":
			filters.InboundRuleAccountIds = filterValues
		//case "inbound-rule-from-port-ranges-ids":
		//	filters.InboundRuleFromPortRanges = filterValues
		case "inbound_rule_ip_ranges":
			filters.InboundRuleIpRanges = filterValues
		case "inbound_rule_protocols":
			filters.InboundRuleProtocols = filterValues
		case "inbound_rule_security_group_ids":
			filters.InboundRuleSecurityGroupIds = filterValues
		case "inbound_rule_security_group_names":
			filters.InboundRuleSecurityGroupNames = filterValues
		// case "InboundRuleToPortRanges":
		// 	filters.InboundRuleToPortRanges = filterValues
		case "net_ids":
			filters.NetIds = filterValues
		case "outbound_rule_account_ids":
			filters.OutboundRuleAccountIds = filterValues
		// case "OutboundRuleFromPortRanges":
		// 	filters.OutboundRuleFromPortRanges = filterValues
		case "outbound_rule_ip_ranges":
			filters.OutboundRuleIpRanges = filterValues
		case "outbound_rule_protocols":
			filters.OutboundRuleProtocols = filterValues
		case "outbound_rule_security_group_ids":
			filters.OutboundRuleSecurityGroupIds = filterValues
		case "outbound_rule_recurity_group_names":
			filters.OutboundRuleSecurityGroupNames = filterValues
		// case "OutboundRuleToPortRanges":
		// 	filters.OutboundRuleToPortRanges = filterValues
		case "security_group_ids":
			filters.SecurityGroupIds = filterValues
		case "security_group_names":
			filters.SecurityGroupNames = filterValues
		case "tag_keys":
			filters.TagKeys = filterValues
		case "tag_values":
			filters.TagValues = filterValues
		case "tags":
			filters.Tags = filterValues
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
