package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
						"tags": dataSourceTagsSchema(),
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPISecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadSecurityGroupsRequest{}

	filters, filtersOk := d.GetOk("filter")
	gn, gnOk := d.GetOk("security_group_names")
	gid, gidOk := d.GetOk("security_group_ids")
	var filter oscgo.FiltersSecurityGroup
	if gnOk {
		var g []string
		for _, v := range gn.([]interface{}) {
			g = append(g, v.(string))
		}
		filter.SetSecurityGroupNames(g)
		req.SetFilters(filter)
	}

	if gidOk {
		var g []string
		for _, v := range gid.([]interface{}) {
			g = append(g, v.(string))
		}
		filter.SetSecurityGroupIds(g)
		req.SetFilters(filter)
	}

	if filtersOk {
		req.SetFilters(buildOutscaleOAPIDataSourceSecurityGroupFilters(filters.(*schema.Set)))
	}

	var err error
	var resp oscgo.ReadSecurityGroupsResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.SecurityGroupApi.ReadSecurityGroups(context.Background()).ReadSecurityGroupsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	var errString string
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidSecurityGroupID.NotFound") ||
			strings.Contains(fmt.Sprint(err), "InvalidGroup.NotFound") {
			resp.SetSecurityGroups(nil)
			err = nil
		} else {
			errString = err.Error()
		}

		return fmt.Errorf("Error on SGStateRefresh: %s", errString)
	}

	if resp.GetSecurityGroups() == nil || len(resp.GetSecurityGroups()) == 0 {
		return fmt.Errorf("Unable to find Security Groups by the following %s", utils.ToJSONString(req.Filters))
	}

	sg := make([]map[string]interface{}, len(resp.GetSecurityGroups()))

	for k, v := range resp.GetSecurityGroups() {
		s := make(map[string]interface{})

		s["security_group_id"] = v.GetSecurityGroupId()
		s["security_group_name"] = v.GetSecurityGroupName()
		s["description"] = v.GetDescription()
		if v.GetNetId() != "" {
			s["net_id"] = v.GetNetId()
		}
		s["account_id"] = v.GetAccountId()
		s["tags"] = tagsOSCAPIToMap(v.GetTags())
		s["inbound_rules"] = flattenOAPISecurityGroupRule(v.GetInboundRules())
		s["outbound_rules"] = flattenOAPISecurityGroupRule(v.GetOutboundRules())
		sg[k] = s
	}

	log.Printf("[DEBUG] security_groups %+v", sg)

	d.SetId(resource.UniqueId())

	err = d.Set("security_groups", sg)

	return err
}

func flattenOAPISecurityGroupRule(p []oscgo.SecurityGroupRule) []map[string]interface{} {
	ips := make([]map[string]interface{}, len(p))

	for k, v := range p {
		ip := make(map[string]interface{})
		if v.GetFromPortRange() != 0 {
			ip["from_port_range"] = v.GetFromPortRange()
		}
		if v.GetIpProtocol() != "" {
			ip["ip_protocol"] = v.GetIpProtocol()
		}
		if v.GetToPortRange() != 0 {
			ip["to_port_range"] = v.GetToPortRange()
		}

		if v.GetIpRanges() != nil && len(v.GetIpRanges()) > 0 {
			ip["ip_ranges"] = v.GetIpRanges()
		}

		/*if v.PrefixListIds != nil && len(v.PrefixListIds) > 0 {
			ip["prefix_list_ids"] = v.PrefixListIds
		}*/

		if v.GetSecurityGroupsMembers() != nil && len(v.GetSecurityGroupsMembers()) > 0 {
			grp := make([]map[string]interface{}, len(v.GetSecurityGroupsMembers()))
			for i, v := range v.GetSecurityGroupsMembers() {
				g := make(map[string]interface{})

				if v.GetAccountId() != "" {
					g["account_id"] = v.GetAccountId()
				}
				if v.GetSecurityGroupName() != "" {
					g["security_group_name"] = v.GetSecurityGroupName()
				}
				if v.GetSecurityGroupId() != "" {
					g["security_group_id"] = v.GetSecurityGroupId()
				}

				grp[i] = g
			}
			ip["security_groups_members"] = grp
		}

		ips[k] = ip
		// s.Add(ip)
	}

	return ips
}
