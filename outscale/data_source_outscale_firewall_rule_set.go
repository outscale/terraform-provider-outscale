package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleFirewallRuleSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleFirewallRuleSetRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"group_name": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"group_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"groups": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"to_port": {
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
								Type:         schema.TypeString,
								ValidateFunc: validateCIDRNetworkAddress,
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
			"ip_permissions_egress": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"groups": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"to_port": {
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
								Type:         schema.TypeString,
								ValidateFunc: validateCIDRNetworkAddress,
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
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag_set": {
				Type: schema.TypeMap,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleFirewallRuleSetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeSecurityGroupsInput{}

	filters, filtersOk := d.GetOk("filter")
	gn, gnOk := d.GetOk("group_name")
	gid, gidOk := d.GetOk("group_id")

	if filtersOk {
		req.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if gnOk {
		var g []*string
		for _, v := range gn.([]interface{}) {
			g = append(g, aws.String(v.(string)))
		}
		req.GroupNames = g
	}
	if gidOk {
		req.GroupIds = []*string{aws.String(gid.(string))}
	}

	fmt.Printf("[DEBUG] REQ %s", req)

	var resp *fcu.DescribeSecurityGroupsOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeSecurityGroups(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "InvalidSecurityGroupID.NotFound") || strings.Contains(err.Error(), "InvalidGroup.NotFound") {
			resp = nil
			err = nil
		}

		if err != nil {
			return fmt.Errorf("\nError on SGStateRefresh: %s", err)
		}
	}

	if resp == nil || len(resp.SecurityGroups) == 0 {
		return fmt.Errorf("Unable to find Security Group")
	}

	if len(resp.SecurityGroups) > 1 {
		return fmt.Errorf("multiple results returned, please use a more specific criteria in your query")
	}

	sg := resp.SecurityGroups[0]

	d.SetId(*sg.GroupId)
	d.Set("group_id", sg.GroupId)
	d.Set("group_description", sg.Description)
	d.Set("group_name", sg.GroupName)
	d.Set("vpc_id", sg.VpcId)
	d.Set("owner_id", sg.OwnerId)
	d.Set("tag_set", tagsToMap(sg.Tags))
	d.Set("ip_permissions", flattenIPPermissions(sg.IpPermissions))
	d.Set("ip_permissions_egress", flattenIPPermissions(sg.IpPermissionsEgress))

	return nil
}

func flattenIPPermissions(p []*fcu.IpPermission) []map[string]interface{} {
	ips := make([]map[string]interface{}, len(p))

	for k, v := range p {
		ip := make(map[string]interface{})
		ip["from_port"] = v.FromPort
		ip["ip_protocol"] = v.IpProtocol
		ip["to_port"] = v.ToPort

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
				"user_id":    v.UserId,
				"group_name": v.GroupName,
				"group_id":   v.GroupId,
			}
		}
		ip["groups"] = grp

		ips[k] = ip
	}

	return ips
}
