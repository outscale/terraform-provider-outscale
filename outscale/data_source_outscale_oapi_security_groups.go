package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleOAPISecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPISecurityGroupsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"security_group_names": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"security_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
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
				},
			},
		},
	}
}

func dataSourceOutscaleOAPISecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	req := &oapi.ReadSecurityGroupsRequest{}

	filters, filtersOk := d.GetOk("filter")
	gn, gnOk := d.GetOk("security_group_names")
	gid, gidOk := d.GetOk("security_group_ids")

	if gnOk {
		var g []string
		for _, v := range gn.([]interface{}) {
			g = append(g, v.(string))
		}
		req.Filters.SecurityGroupNames = g
	}

	if gidOk {
		var g []string
		for _, v := range gid.([]interface{}) {
			g = append(g, v.(string))
		}
		req.Filters.SecurityGroupNames = g
	}

	if filtersOk {
		req.Filters = buildOutscaleOAPIDataSourceSecurityGroupFilters(filters.(*schema.Set))
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

	if result == nil || len(result.SecurityGroups) == 0 {
		return fmt.Errorf("Unable to find Security Group")
	}

	sg := make([]map[string]interface{}, len(result.SecurityGroups))

	for k, v := range result.SecurityGroups {
		s := make(map[string]interface{})

		s["security_group_id"] = v.SecurityGroupId
		s["security_group_name"] = v.SecurityGroupName
		s["description"] = v.Description
		if v.NetId != "" {
			s["net_id"] = v.NetId
		}
		s["account_id"] = v.AccountId
		s["tags"] = tagsOAPIToMap(v.Tags)
		s["inbound_rules"] = flattenOAPISecurityGroupRule(v.InboundRules)
		s["outbound_rules"] = flattenOAPISecurityGroupRule(v.OutboundRules)
		s["tags"] = tagsOAPIToMap(v.Tags)
		sg[k] = s
	}

	fmt.Printf("[DEBUG] security_groups %s", sg)

	d.SetId(resource.UniqueId())
	d.Set("request_id", result.ResponseContext.RequestId)
	err = d.Set("security_groups", sg)

	return err
}

func flattenOAPISecurityGroupRule(p []oapi.SecurityGroupRule) []map[string]interface{} {
	ips := make([]map[string]interface{}, len(p))

	for k, v := range p {
		ip := make(map[string]interface{})
		if v.FromPortRange != 0 {
			ip["from_port_range"] = v.FromPortRange
		}
		if v.IpProtocol != "" {
			ip["ip_protocol"] = v.IpProtocol
		}
		if v.ToPortRange != 0 {
			ip["to_port_range"] = v.ToPortRange
		}

		if v.IpRanges != nil && len(v.IpRanges) > 0 {
			ip["ip_ranges"] = v.IpRanges
		}

		if v.PrefixListIds != nil && len(v.PrefixListIds) > 0 {
			ip["prefix_list_ids"] = v.PrefixListIds
		}

		if v.SecurityGroupsMembers != nil && len(v.SecurityGroupsMembers) > 0 {
			grp := make([]map[string]interface{}, len(v.SecurityGroupsMembers))
			for i, v := range v.SecurityGroupsMembers {
				g := make(map[string]interface{})

				if v.AccountId != "" {
					g["account_id"] = v.AccountId
				}
				if v.SecurityGroupName != "" {
					g["security_group_name"] = v.SecurityGroupName
				}
				if v.SecurityGroupId != "" {
					g["security_group_id"] = v.SecurityGroupId
				}

				grp[i] = g
			}
			ip["security_group_members"] = grp
		}

		ips[k] = ip
		// s.Add(ip)
	}

	return ips
}
