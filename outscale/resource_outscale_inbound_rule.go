package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/outscale/osc-go/oapi"
)

func resourceOutscaleOAPIInboundRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIInboundRuleCreate,
		Read:   resourceOutscaleOAPIInboundRuleRead,
		Delete: resourceOutscaleOAPIInboundRuleDelete,

		Schema: map[string]*schema.Schema{
			"ip_range": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"from_port_range": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_group_name_to_link": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_group_account_id_to_link": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"to_port_range": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"rules": getIPOAPIPermissionsSchema(false),
			"reques_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPIInboundRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI
	sgID := d.Get("firewall_rules_set_id").(string)

	awsMutexKV.Lock(sgID)
	defer awsMutexKV.Unlock(sgID)

	sg, _, err := oapiFindResourceSecurityGroup(conn, sgID)
	if err != nil {
		return err
	}

	perms, err := expandOAPIIPPermIngress(d, sg)
	if err != nil {
		return err
	}

	ippems := d.Get("inbound_rule").([]interface{})

	if err := validateAwsSecurityGroupRule(ippems); err != nil {
		return err
	}

	var autherr error
	log.Printf("[DEBUG] Authorizing security group %s %s rule: %#v", sgID, "Ingress", perms)

	req := oapi.CreateSecurityGroupRuleRequest{
		SecurityGroupId: sg.SecurityGroupId,
		Rules:           perms,
	}

	autherr = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		_, err = conn.POST_CreateSecurityGroupRule(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if autherr != nil {
		if awsErr, ok := autherr.(awserr.Error); ok {
			if awsErr.Code() == "InvalidPermission.Duplicate" {
				return fmt.Errorf(`[WARN] A duplicate Security Group rule was found on (%s). This may be
a side effect of a now-fixed Terraform issue causing two security groups with
identical attributes but different source_security_group_ids to overwrite each
other in the state. See https://github.com/hashicorp/terraform/pull/2376 for more
information and instructions for recovery. Error message: %s`, sgID, awsErr.Message())
			}
		}

		return fmt.Errorf(
			"Error authorizing security group rule type %s: %s",
			"", autherr)
	}

	id := ipOAPIPermissionIDHash(sgID, "", perms)
	log.Printf("[DEBUG] Computed group rule ID %s", id)

	retErr := resource.Retry(5*time.Minute, func() *resource.RetryError {
		sg, _, err := findOAPIResourceSecurityGroup(conn, sgID)

		if err != nil {
			log.Printf("[DEBUG] Error finding Security Group (%s) for Rule (%s): %s", sgID, id, err)
			return resource.NonRetryableError(err)
		}

		var rules []oapi.SecurityGroupRule
		rules = sg.InboundRules

		rule := findOAPIRuleMatch(perms, rules)

		if rule == nil {
			log.Printf("[DEBUG] Unable to find matching %s Security Group Rule (%s) for Group %s",
				"", id, sgID)
			return resource.RetryableError(fmt.Errorf("No match found"))
		}

		return nil
	})

	if retErr != nil {
		return fmt.Errorf("Error finding matching %s Security Group Rule (%s) for Group %s",
			"", id, sgID)
	}

	d.SetId(id)
	return nil
}

func resourceOutscaleOAPIInboundRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI
	sgID := d.Get("firewall_rules_set_id").(string)
	sg, reqID, err := findOAPIResourceSecurityGroup(conn, sgID)
	if _, notFound := err.(oapiSecurityGroupNotFound); notFound {
		// The security group containing this rule no longer exists.
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error finding security group (%s) for rule (%s): %s", sgID, d.Id(), err)
	}

	var rule []oapi.SecurityGroupRule
	var rules []oapi.SecurityGroupRule
	ruleType := "ingress"
	rules = sg.InboundRules

	p, err := expandOAPIIPPermIngress(d, sg)
	if err != nil {
		return err
	}

	if len(rules) == 0 {
		log.Printf("[WARN] No %s rules were found for Security Group (%s) looking for Security Group Rule (%s)",
			ruleType, sg.SecurityGroupName, d.Id())
		d.SetId("")
		return nil
	}

	rule = findOAPIRuleMatch(p, rules)

	if rule == nil {
		log.Printf("[DEBUG] Unable to find matching %s Security Group Rule (%s) for Group %s",
			ruleType, d.Id(), sgID)
		d.SetId("")
		return nil
	}

	if ips, err := setOAPIFromIPPerm(d, sg, p); err != nil {
		return d.Set("inbound_rule", ips)
	}
	return d.Set("request_id", reqID)
}

func resourceOutscaleOAPIInboundRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI
	sgID := d.Get("firewall_rules_set_id").(string)

	awsMutexKV.Lock(sgID)
	defer awsMutexKV.Unlock(sgID)

	sg, _, err := findOAPIResourceSecurityGroup(conn, sgID)
	if err != nil {
		return err
	}

	perms, err := expandOAPIIPPermIngress(d, sg)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Revoking security group %#v %s rule: %#v",
		sgID, "ingress", perms)
	req := oapi.DeleteSecurityGroupRuleRequest{
		SecurityGroupId: sg.SecurityGroupId,
		Rules:           perms,
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.POST_DeleteSecurityGroupRule(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf(
			"Error revoking security group %s rules: %s",
			sgID, err)
	}

	d.SetId("")

	return nil
}

// #################################

func expandOAPIIPPermIngress(d *schema.ResourceData, sg *oapi.SecurityGroup) ([]oapi.SecurityGroupRule, error) {
	ippems := d.Get("inbound_rule").([]interface{})
	perms := make([]oapi.SecurityGroupRule, len(ippems))

	return expandOAPIIPPerm(d, sg, perms, ippems)
}

func oapiFindResourceSecurityGroup(conn *oapi.Client, id string) (*oapi.SecurityGroup, *string, error) {
	req := oapi.ReadSecurityGroupsRequest{
		Filters: oapi.FiltersSecurityGroup{
			InboundRuleSecurityGroupIds: []string{id},
		},
	}

	var err error
	var resp *oapi.POST_ReadSecurityGroupsResponses
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_ReadSecurityGroups(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err, ok := err.(awserr.Error); ok && err.Code() == "InvalidGroup.NotFound" {
		return nil, nil, oapiSecurityGroupNotFound{id, nil}
	}
	if err != nil {
		return nil, nil, err
	}
	if resp == nil {
		return nil, nil, oapiSecurityGroupNotFound{id, nil}
	}
	if len(resp.OK.SecurityGroups) != 1 {
		return nil, nil, nil // oapiSecurityGroupNotFound{id, resp.OK.SecurityGroups}
	}

	return &resp.OK.SecurityGroups[0], &resp.OK.ResponseContext.RequestId, nil
}
