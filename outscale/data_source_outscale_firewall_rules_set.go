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
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeMap},
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
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeMap},
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
							Type: schema.TypeList,
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

	d.SetId(resource.UniqueId())

	err = d.Set("security_group_info", sg)

	return err
}
