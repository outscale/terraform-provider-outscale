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

func dataSourceOutscaleFirewallRulesSets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleFirewallRulesSetsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"group_name": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"group_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"security_group_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
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
						"ip_permissions":        getDSIPPerms(),
						"ip_permissions_egress": getDSIPPerms(),
						"owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_set": tagsSchemaComputed(),
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

func getDSIPPerms() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		// Set:      resourceOutscaleSecurityGroupRuleHash,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"from_port": {
					Type:     schema.TypeInt,
					Computed: true,
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
						Type: schema.TypeString,
					},
				},
				"groups": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeMap},
				},
				"self": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func dataSourceOutscaleFirewallRulesSetsRead(d *schema.ResourceData, meta interface{}) error {
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
		var g []*string
		for _, v := range gid.([]interface{}) {
			g = append(g, aws.String(v.(string)))
		}
		req.GroupIds = g
	}

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

	sg := make([]map[string]interface{}, len(resp.SecurityGroups))

	for k, v := range resp.SecurityGroups {
		s := make(map[string]interface{})

		s["group_id"] = *v.GroupId
		s["group_name"] = *v.GroupName
		s["group_description"] = *v.Description
		if v.VpcId != nil {
			s["vpc_id"] = *v.VpcId
		}
		s["owner_id"] = *v.OwnerId
		s["tag_set"] = tagsToMap(v.Tags)
		s["ip_permissions"] = flattenIPPermissions(v.IpPermissions)
		s["ip_permissions_egress"] = flattenIPPermissions(v.IpPermissionsEgress)

		sg[k] = s
	}

	fmt.Printf("[DEBUG] security_group_info %s", sg)

	d.Set("request_id", resp.RequestId)

	d.SetId(resource.UniqueId())

	err = d.Set("security_group_info", sg)

	return err
}

func flattenIPPermissions(p []*fcu.IpPermission) []map[string]interface{} {
	ips := make([]map[string]interface{}, len(p))

	// s := &schema.Set{
	// 	F: resourceOutscaleSecurityGroupRuleHash,
	// }

	for k, v := range p {
		ip := make(map[string]interface{})
		if v.FromPort != nil {
			ip["from_port"] = *v.FromPort
		}
		if v.IpProtocol != nil {
			ip["ip_protocol"] = *v.IpProtocol
		}
		if v.ToPort != nil {
			ip["to_port"] = *v.ToPort
		}

		if v.IpRanges != nil && len(v.IpRanges) > 0 {
			ipr := make([]string, len(v.IpRanges))
			for i, v := range v.IpRanges {
				if v.CidrIp != nil {
					ipr[i] = *v.CidrIp
				}
			}
			ip["ip_ranges"] = ipr
		}

		if v.PrefixListIds != nil && len(v.PrefixListIds) > 0 {
			prx := make([]string, len(v.PrefixListIds))
			for i, v := range v.PrefixListIds {
				if v.PrefixListId != nil {
					prx[i] = *v.PrefixListId
				}
			}
			ip["prefix_list_ids"] = prx
		}

		if v.UserIdGroupPairs != nil && len(v.UserIdGroupPairs) > 0 {
			grp := make([]map[string]interface{}, len(v.UserIdGroupPairs))
			for i, v := range v.UserIdGroupPairs {
				g := make(map[string]interface{})

				if v.UserId != nil {
					g["user_id"] = *v.UserId
				}
				if v.GroupName != nil {
					g["group_name"] = *v.GroupName
				}
				if v.GroupId != nil {
					g["group_id"] = *v.GroupId
				}

				grp[i] = g
			}
			ip["groups"] = grp
		}

		ips[k] = ip
		// s.Add(ip)
	}

	return ips
}
