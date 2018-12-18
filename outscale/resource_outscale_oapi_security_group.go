package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleOAPISecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPISecurityGroupCreate,
		Read:   resourceOutscaleOAPISecurityGroupRead,
		Delete: resourceOutscaleOAPISecurityGroupDelete,
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
				"groups": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeMap},
				},
			},
		},
	}
}

func resourceOutscaleOAPISecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	securityGroupOpts := &oapi.CreateSecurityGroupRequest{}

	if v, ok := d.GetOk("net_id"); ok {
		securityGroupOpts.NetId = v.(string)
	}

	if v := d.Get("description"); v != nil {
		securityGroupOpts.Description = v.(string)
	} else {
		return fmt.Errorf("please provide a group description, its a required argument")
	}

	var groupName string
	if v, ok := d.GetOk("security_group_name"); ok {
		groupName = v.(string)
	} else {
		groupName = resource.UniqueId()
	}
	securityGroupOpts.SecurityGroupName = groupName

	fmt.Printf(
		"[DEBUG] Security Group create configuration: %#v", securityGroupOpts)

	var createResp *oapi.CreateSecurityGroupResponse
	var resp *oapi.POST_CreateSecurityGroupResponses
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_CreateSecurityGroup(*securityGroupOpts)

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
			errString = err.Error()
		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
		}

		return fmt.Errorf("Error creating Security Group: %s", errString)
	}

	createResp = resp.OK

	d.SetId(createResp.SecurityGroup.SecurityGroupId)

	fmt.Printf("\n\n[INFO] Security Group ID: %s", d.Id())

	// Wait for the security group to truly exist
	fmt.Printf("\n\n[DEBUG] Waiting for Security Group (%s) to exist", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{""},
		Target:  []string{"exists"},
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
		if err := setOAPITags(conn, d); err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	return resourceOutscaleOAPISecurityGroupRead(d, meta)
}

func resourceOutscaleOAPISecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	sgRaw, _, err := SGOAPIStateRefreshFunc(conn, d.Id())()
	if err != nil {
		return err
	}
	if sgRaw == nil {
		d.SetId("")
		return nil
	}

	group := sgRaw.(oapi.SecurityGroup)

	req := &oapi.ReadSecurityGroupsRequest{}
	req.Filters = oapi.FiltersSecurityGroup{SecurityGroupIds: []string{group.SecurityGroupId}}

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

	if len(result.SecurityGroups) > 1 {
		return fmt.Errorf("multiple results returned, please use a more specific criteria in your query")
	}

	sg := result.SecurityGroups[0]

	d.SetId(sg.SecurityGroupId)
	d.Set("security_group_id", sg.SecurityGroupId)
	d.Set("description", sg.Description)
	if sg.SecurityGroupName != "" {
		d.Set("security_group_name", sg.SecurityGroupName)
	}
	d.Set("net_id", sg.NetId)
	d.Set("account_id", sg.AccountId)
	d.Set("tags", tagsOAPIToMap(sg.Tags))
	d.Set("request_id", result.ResponseContext.RequestId)

	if err := d.Set("inbound_rules", flattenOAPISecurityGroupRule(sg.InboundRules)); err != nil {
		return err
	}

	return d.Set("outbound_rules", flattenOAPISecurityGroupRule(sg.OutboundRules))
}

func resourceOutscaleOAPISecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	fmt.Printf("\n\n[DEBUG] Security Group destroy: %v", d.Id())

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.POST_DeleteSecurityGroup(oapi.DeleteSecurityGroupRequest{
			SecurityGroupId: d.Id(),
		})

		var errString string

		if err != nil || resp.OK == nil {
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") || strings.Contains(err.Error(), "DependencyViolation") {
					return resource.RetryableError(err)
				} else if strings.Contains(err.Error(), "InvalidGroup.NotFound") {
					return nil
				}
				return resource.NonRetryableError(err)

			} else if resp.Code401 != nil {
				errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
			} else if resp.Code400 != nil {
				errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
			} else if resp.Code500 != nil {
				errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
			}

			return resource.NonRetryableError(fmt.Errorf("Error on SGStateRefresh: %s", errString))
		}

		return nil
	})
}

func flattenOAPISecurityGroups(list []*fcu.UserIdGroupPair, ownerID *string) []*fcu.GroupIdentifier {
	result := make([]*fcu.GroupIdentifier, 0, len(list))
	for _, g := range list {
		var userID *string
		if g.UserId != nil && *g.UserId != "" && (ownerID == nil || *ownerID != *g.UserId) {
			userID = g.UserId
		}
		// userid nil here for same vpc groups

		vpc := g.GroupName == nil || *g.GroupName == ""
		var id *string
		if vpc {
			id = g.GroupId
		} else {
			id = g.GroupName
		}

		// id is groupid for vpcs
		// id is groupname for non vpc (classic)

		if userID != nil {
			id = aws.String(*userID + "/" + *id)
		}

		if vpc {
			result = append(result, &fcu.GroupIdentifier{
				GroupId: id,
			})
		} else {
			result = append(result, &fcu.GroupIdentifier{
				GroupId:   g.GroupId,
				GroupName: id,
			})
		}
	}
	return result
}

// SGOAPIStateRefreshFunc ...
func SGOAPIStateRefreshFunc(conn *oapi.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := &oapi.ReadSecurityGroupsRequest{
			Filters: oapi.FiltersSecurityGroup{
				SecurityGroupIds: []string{id},
			},
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

			return nil, "", fmt.Errorf("Error on SGStateRefresh: %s", errString)
		}

		if resp == nil {
			return nil, "", nil
		}

		group := resp.OK.SecurityGroups[0]
		return group, "exists", nil
	}
}
