package oapi

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/samber/lo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
)

func DataSourceOutscaleSecurityGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleSecurityGroupsRead,

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
						"tags": TagsSchemaComputedSDK(),
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

func DataSourceOutscaleSecurityGroupsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadSecurityGroupsRequest{}

	filters, filtersOk := d.GetOk("filter")
	gn, gnOk := d.GetOk("security_group_names")
	gid, gidOk := d.GetOk("security_group_ids")
	var filter osc.FiltersSecurityGroup
	if gnOk {
		g := lo.Map(gn.([]any), func(v any, _ int) string {
			return v.(string)
		})
		filter.SecurityGroupNames = &g
		req.Filters = &filter
	}

	if gidOk {
		g := lo.Map(gid.([]any), func(v any, _ int) string {
			return v.(string)
		})
		filter.SecurityGroupIds = &g
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
			errString = err.Error()
		}

		return diag.Errorf("error on sgstaterefresh: %s", errString)
	}

	if resp.SecurityGroups == nil || len(*resp.SecurityGroups) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	sg := make([]map[string]any, len(*resp.SecurityGroups))

	for k, v := range *resp.SecurityGroups {
		s := make(map[string]any)

		s["security_group_id"] = v.SecurityGroupId
		s["security_group_name"] = v.SecurityGroupName
		s["description"] = v.Description
		if ptr.From(v.NetId) != "" {
			s["net_id"] = v.NetId
		}
		s["account_id"] = v.AccountId
		s["tags"] = FlattenOAPITagsSDK(v.Tags)
		s["inbound_rules"] = flattenOAPISecurityGroupRule(v.InboundRules)
		s["outbound_rules"] = flattenOAPISecurityGroupRule(v.OutboundRules)
		sg[k] = s
	}

	log.Printf("[DEBUG] security_groups %+v", sg)

	d.SetId(id.UniqueId())

	err = d.Set("security_groups", sg)

	return diag.FromErr(err)
}

func flattenOAPISecurityGroupRule(p []osc.SecurityGroupRule) []map[string]any {
	ips := make([]map[string]any, len(p))

	for k, v := range p {
		ip := make(map[string]any)
		if v.FromPortRange != 0 {
			ip["from_port_range"] = v.FromPortRange
		}
		if v.IpProtocol != "" {
			ip["ip_protocol"] = v.IpProtocol
		}
		if v.ToPortRange != 0 {
			ip["to_port_range"] = v.ToPortRange
		}

		if len(v.IpRanges) > 0 {
			ip["ip_ranges"] = v.IpRanges
		}

		/*if v.PrefixListIds != nil && len(v.PrefixListIds) > 0 {
			ip["prefix_list_ids"] = v.PrefixListIds
		}*/

		if len(v.SecurityGroupsMembers) > 0 {
			grp := make([]map[string]any, len(v.SecurityGroupsMembers))
			for i, v := range v.SecurityGroupsMembers {
				g := make(map[string]any)

				if ptr.From(v.AccountId) != "" {
					g["account_id"] = v.AccountId
				}
				if ptr.From(v.SecurityGroupName) != "" {
					g["security_group_name"] = v.SecurityGroupName
				}
				if v.SecurityGroupId != "" {
					g["security_group_id"] = v.SecurityGroupId
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
