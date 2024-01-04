package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOutscaleOAPISecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPISecurityGroupCreate,
		Read:   resourceOutscaleOAPISecurityGroupRead,
		Delete: resourceOutscaleOAPISecurityGroupDelete,
		Update: resourceOutscaleOAPISecurityGroupUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Managed by Terraform",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) > 255 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 255 characters", k))
					}
					return
				},
			},
			"security_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"net_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			// comouted
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"inbound_rules":  getOAPIIPPerms(),
			"outbound_rules": getOAPIIPPerms(),
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remove_default_outbound_rule": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"tags": tagsListOAPISchema(),
			"tag":  tagsSchema(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getOAPIIPPerms() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"from_port_range": {
					Type:     schema.TypeInt,
					Computed: true,
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
					},
				},
				"security_groups_members": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeMap},
				},
			},
		},
	}
}

func resourceOutscaleOAPISecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	securityGroupOpts := oscgo.CreateSecurityGroupRequest{}

	net, have_net := d.GetOk("net_id")
	if have_net {
		securityGroupOpts.SetNetId(net.(string))
	}

	if v := d.Get("description"); v != nil {
		securityGroupOpts.SetDescription(v.(string))
	} else {
		return fmt.Errorf("please provide a group description, its a required argument")
	}

	var groupName string
	if v, ok := d.GetOk("security_group_name"); ok {
		groupName = v.(string)
	} else {
		groupName = resource.UniqueId()
	}
	securityGroupOpts.SetSecurityGroupName(groupName)

	empty := d.Get("remove_default_outbound_rule").(bool)
	if empty && !have_net {
		return fmt.Errorf("remove_default_outbound_rule is useless when net_id is not set (so with public SG)")
	}

	log.Printf("[DEBUG] Security Group create configuration: %#v", securityGroupOpts)

	var resp oscgo.CreateSecurityGroupResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.SecurityGroupApi.CreateSecurityGroup(context.Background()).CreateSecurityGroupRequest(securityGroupOpts).Execute()

		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	var errString string

	if err != nil {
		errString = err.Error()

		return fmt.Errorf("Error creating Security Group: %s", errString)
	}

	id := resp.SecurityGroup.GetSecurityGroupId()
	if empty {
		ipRange := "0.0.0.0/0"
		ipProtocol := "-1"
		emptierOpts := oscgo.DeleteSecurityGroupRuleRequest{
			Flow:            "Outbound",
			SecurityGroupId: id,
			IpRange:         &ipRange,
			IpProtocol:      &ipProtocol,
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, httpResp, err := conn.SecurityGroupRuleApi.DeleteSecurityGroupRule(context.Background()).DeleteSecurityGroupRuleRequest(emptierOpts).Execute()

			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("Error Emptying Security Group: %s", err.Error())
		}

	}

	d.SetId(id)

	log.Printf("[INFO] Security Group ID: %s", d.Id())

	// Wait for the security group to truly exist
	log.Printf("[DEBUG] Waiting for Security Group (%s) to exist", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{""},
		Target:  []string{"exists", "failed"},
		Refresh: SGOAPIStateRefreshFunc(conn, d.Id()),
		Timeout: 3 * time.Minute,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for Security Group (%s) to become available: %s",
			d.Id(), err)
	}

	if d.IsNewResource() {
		if err := setOSCAPITags(conn, d); err != nil {
			return err
		}
	}
	return resourceOutscaleOAPISecurityGroupRead(d, meta)
}

func resourceOutscaleOAPISecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	sg, _, err := readSecurityGroups(conn, d.Id())
	if err != nil {
		return err
	}
	if sg == nil {
		utils.LogManuallyDeleted("SecurityGroup", d.Id())
		d.SetId("")
		return nil
	}

	if err := d.Set("security_group_id", sg.GetSecurityGroupId()); err != nil {
		return err
	}
	if err := d.Set("description", sg.GetDescription()); err != nil {
		return err
	}
	if sg.GetSecurityGroupName() != "" {
		if err := d.Set("security_group_name", sg.GetSecurityGroupName()); err != nil {
			return err
		}
	}
	if err := d.Set("net_id", sg.GetNetId()); err != nil {
		return err
	}
	if err := d.Set("account_id", sg.GetAccountId()); err != nil {
		return err
	}
	if err := d.Set("tags", tagsOSCAPIToMap(sg.GetTags())); err != nil {
		return err
	}

	if err := d.Set("inbound_rules", flattenOAPISecurityGroupRule(sg.GetInboundRules())); err != nil {
		return err
	}

	d.SetId(sg.GetSecurityGroupId())

	return d.Set("outbound_rules", flattenOAPISecurityGroupRule(sg.GetOutboundRules()))
}

func resourceOutscaleOAPISecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	log.Printf("[DEBUG] Security Group destroy: %v", d.Id())
	securityGroupID := d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.SecurityGroupApi.DeleteSecurityGroup(context.Background()).DeleteSecurityGroupRequest(oscgo.DeleteSecurityGroupRequest{
			SecurityGroupId: &securityGroupID,
		}).Execute()
		if err != nil {
			if strings.Contains(err.Error(), "DependencyProblem") {
				return resource.RetryableError(err)
			}
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
}

func SGOAPIStateRefreshFunc(conn *oscgo.APIClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		securityGroup, _, err := readSecurityGroups(conn, id)
		if err != nil {
			return nil, "failed", err
		}
		if securityGroup == nil {
			return nil, "failed", fmt.Errorf("SecurityGroup Not Found")
		}
		return securityGroup, "exists", nil
	}
}

func resourceOutscaleOAPISecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}
	return resourceOutscaleOAPISecurityGroupRead(d, meta)
}

func readSecurityGroups(client *oscgo.APIClient, securityGroupID string) (*oscgo.SecurityGroup, *oscgo.ReadSecurityGroupsResponse, error) {
	filters := oscgo.ReadSecurityGroupsRequest{
		Filters: &oscgo.FiltersSecurityGroup{
			SecurityGroupIds: &[]string{securityGroupID},
		},
	}

	var err error
	var resp oscgo.ReadSecurityGroupsResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := client.SecurityGroupApi.ReadSecurityGroups(context.Background()).ReadSecurityGroupsRequest(filters).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error reading the Outscale Security Group(%s): %s", securityGroupID, err)
	}

	if len(*resp.SecurityGroups) == 0 {
		return nil, nil, nil
	}
	return &resp.GetSecurityGroups()[0], &resp, nil
}
