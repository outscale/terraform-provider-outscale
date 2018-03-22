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

func dataSourceOutscaleOAPIFirewallRulesSets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIFirewallRulesSetsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"firewall_rules_set_name": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"firewall_rules_set_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"firewall_rules_sets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"firewall_rules_set_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"firewall_rules_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lin_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inbound_rule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"from_port_range": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"groups": {
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
						"outbound_rule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"from_port_range": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"groups": {
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
						"tags": {
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
				},
			},
		},
	}
}

func dataSourceOutscaleOAPIFirewallRulesSetsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeSecurityGroupsInput{}

	filters, filtersOk := d.GetOk("filter")
	gn, gnOk := d.GetOk("firewall_rules_set_name")
	gid, gidOk := d.GetOk("firewall_rules_set_id")

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

		s["firewall_rules_set_id"] = *v.GroupId
		s["firewall_rules_set_name"] = *v.GroupName
		s["description"] = *v.Description
		if v.VpcId != nil {
			s["lin_id"] = *v.VpcId
		}
		s["account_id"] = *v.OwnerId
		s["tags"] = tagsToMap(v.Tags)
		s["inbound_rule"] = flattenOAPIIPPermissions(v.IpPermissions)
		s["outbound_rule"] = flattenOAPIIPPermissions(v.IpPermissionsEgress)

		sg[k] = s
	}

	fmt.Printf("[DEBUG] firewall_rules_sets %s", sg)

	d.SetId(resource.UniqueId())

	err = d.Set("firewall_rules_sets", sg)

	return err
}
