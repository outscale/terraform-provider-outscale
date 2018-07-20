package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleInboundRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleInboundRuleCreate,
		Read:   resourceOutscaleInboundRuleRead,
		Delete: resourceOutscaleInboundRuleDelete,

		Schema: map[string]*schema.Schema{
			"cidr_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"from_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_security_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_security_group_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"to_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"ip_permissions": getIPPermissionsSchema(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleInboundRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	sgID := d.Get("group_id").(string)

	awsMutexKV.Lock(sgID)
	defer awsMutexKV.Unlock(sgID)

	sg, _, err := findResourceSecurityGroup(conn, sgID)
	if err != nil {
		return err
	}

	perms, err := expandIPPermIngress(d, sg)
	if err != nil {
		return err
	}

	ippems := d.Get("ip_permissions").([]interface{})

	if err := validateAwsSecurityGroupRule(ippems); err != nil {
		return err
	}

	ruleType := "ingress"
	isVPC := sg.VpcId != nil && *sg.VpcId != ""

	var autherr error
	req := &fcu.AuthorizeSecurityGroupIngressInput{
		GroupId:       sg.GroupId,
		IpPermissions: perms,
	}

	autherr = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		_, err = conn.VM.AuthorizeSecurityGroupIngress(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if autherr != nil {
		if strings.Contains(fmt.Sprint(autherr), "InvalidPermission.Duplicate") {
			return fmt.Errorf(`[WARN] A duplicate Security Group rule was found on (%s). This may be
a side effect of a now-fixed Terraform issue causing two security groups with
identical attributes but different source_security_group_ids to overwrite each
other in the state. See https://github.com/hashicorp/terraform/pull/2376 for more
information and instructions for recovery. Error message: %s`, sgID, autherr)
		}

		return fmt.Errorf(
			"Error authorizing security group rule type %s: %s",
			ruleType, autherr)
	}

	id := ipPermissionIDHash(sgID, ruleType, perms)

	retErr := resource.Retry(5*time.Minute, func() *resource.RetryError {
		sg, _, err := findResourceSecurityGroup(conn, sgID)

		if err != nil {
			return resource.NonRetryableError(err)
		}

		var rules []*fcu.IpPermission
		rules = sg.IpPermissions

		rule := findRuleMatch(perms, rules, isVPC)

		if rule == nil {
			return resource.RetryableError(fmt.Errorf("No match found"))
		}

		return nil
	})

	if retErr != nil {
		return fmt.Errorf("Error finding matching %s Security Group Rule (%s) for Group %s",
			ruleType, id, sgID)
	}

	d.SetId(id)
	return nil
}

func resourceOutscaleInboundRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	sgID := d.Get("group_id").(string)
	sg, reqID, err := findResourceSecurityGroup(conn, sgID)
	if _, notFound := err.(securityGroupNotFound); notFound {
		// The security group containing this rule no longer exists.
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error finding security group (%s) for rule (%s): %s", sgID, d.Id(), err)
	}

	isVPC := sg.VpcId != nil && *sg.VpcId != ""

	var rule *fcu.IpPermission
	var rules []*fcu.IpPermission
	ruleType := "ingress"
	rules = sg.IpPermissions

	p, err := expandIPPermIngress(d, sg)
	if err != nil {
		return err
	}

	if len(rules) == 0 {
		log.Printf("[WARN] No %s rules were found for Security Group (%s) looking for Security Group Rule (%s)",
			ruleType, *sg.GroupName, d.Id())
		d.SetId("")
		return nil
	}

	rule = findRuleMatch(p, rules, isVPC)

	if rule == nil {
		log.Printf("[DEBUG] Unable to find matching %s Security Group Rule (%s) for Group %s",
			ruleType, d.Id(), sgID)
		d.SetId("")
		return nil
	}

	if ips, err := setFromIPPerm(d, sg, p); err != nil {
		return d.Set("ip_permissions", ips)
	}
	return d.Set("request_id", aws.StringValue(reqID))
}

func resourceOutscaleInboundRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	sgID := d.Get("group_id").(string)

	awsMutexKV.Lock(sgID)
	defer awsMutexKV.Unlock(sgID)

	sg, _, err := findResourceSecurityGroup(conn, sgID)
	if err != nil {
		return err
	}

	perms, err := expandIPPermIngress(d, sg)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Revoking security group %#v %s rule: %#v",
		sgID, "ingress", perms)
	req := &fcu.RevokeSecurityGroupIngressInput{
		GroupId:       sg.GroupId,
		IpPermissions: perms,
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.RevokeSecurityGroupIngress(req)

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

func expandIPPermIngress(d *schema.ResourceData, sg *fcu.SecurityGroup) ([]*fcu.IpPermission, error) {
	ippems := d.Get("ip_permissions").([]interface{})
	perms := make([]*fcu.IpPermission, len(ippems))

	return expandIPPerm(d, sg, perms, ippems)
}
