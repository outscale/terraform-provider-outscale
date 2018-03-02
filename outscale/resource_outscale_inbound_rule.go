package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleInboundRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleInboundRuleCreate,
		Read:   resourceOutscaleInboundRuleRead,
		Delete: resourceOutscaleInboundRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
			"ip_permissions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"groups": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"group_name": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"user_id": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"to_port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"ip_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"ip_ranges": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_ip": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"prefix_list_ids": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"prefix_list_id": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceOutscaleInboundRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	sg_id := d.Get("group_id").(string)

	awsMutexKV.Lock(sg_id)
	defer awsMutexKV.Unlock(sg_id)

	sg, err := findResourceSecurityGroup(conn, sg_id)

	fmt.Println("\n[DEDUG] ERROR resourceOutscaleInboundRuleCreate 1 =>", err)

	if err != nil {
		return err
	}

	perm, err := expandIPPerm(d, sg)

	fmt.Println("\n[DEDUG] ERROR resourceOutscaleInboundRuleCreate 2 =>", err)

	if err != nil {
		return err
	}

	if err := validateOutscaleSecurityGroupRule(d); err != nil {
		fmt.Println("\n[DEDUG] ERROR resourceOutscaleInboundRuleCreate 3 =>", err)

		return err
	}

	isVPC := sg.VpcId != nil && *sg.VpcId != ""
	ruleType := "ingress"
	var autherr error
	fmt.Printf("\n\n[DEBUG] Authorizing security group %s %s rule: %#v", sg_id, "Ingress", perm)

	req := &fcu.AuthorizeSecurityGroupIngressInput{
		GroupId:       sg.GroupId,
		IpPermissions: []*fcu.IpPermission{perm},
	}

	_, autherr = conn.VM.AuthorizeSecurityGroupIngress(req)

	fmt.Println("\n[DEDUG] ERROR resourceOutscaleInboundRuleCreate 4 =>", err)

	if autherr != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidPermission.Duplicate") {
			return fmt.Errorf(`[WARN] A duplicate Security Group rule was found on (%s). This may be
a side effect of a now-fixed Terraform issue causing two security groups with
identical attributes but different source_security_group_ids to overwrite each
other in the state. See https://github.com/hashicorp/terraform/pull/2376 for more
information and instructions for recovery. Error message: %s`, sg_id, fmt.Sprint(err))

		}

		return fmt.Errorf(
			"Error authorizing security group rule type %s: %s",
			ruleType, autherr)
	}

	id := ipPermissionIDHash(sg_id, ruleType, perm)
	fmt.Printf("\n\n[DEBUG] Computed group rule ID %s", id)

	retErr := resource.Retry(5*time.Minute, func() *resource.RetryError {
		sg, err := findResourceSecurityGroup(conn, sg_id)

		if err != nil {
			fmt.Println("\n[DEDUG] ERROR resourceOutscaleInboundRuleCreate 6 =>", err)

			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("\n\n[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}

			fmt.Printf("\n\n[DEBUG] Error finding Security Group (%s) for Rule (%s): %s", sg_id, id, err)
			return resource.NonRetryableError(err)
		}

		var rules []*fcu.IpPermission
		rules = sg.IpPermissions

		rule := findRuleMatch(perm, rules, isVPC)

		if rule == nil {
			fmt.Printf("\n\n[DEBUG] Unable to find matching %s Security Group Rule (%s) for Group %s",
				ruleType, id, sg_id)
			return resource.RetryableError(fmt.Errorf("No match found"))
		}

		fmt.Printf("\n\n[DEBUG] Found rule for Security Group Rule (%s): %s", id, rule)
		return nil
	})

	fmt.Println("\n[DEDUG] ERROR resourceOutscaleInboundRuleCreate 7 =>", retErr)

	if retErr != nil {
		return fmt.Errorf("Error finding matching %s Security Group Rule (%s) for Group %s",
			ruleType, id, sg_id)
	}

	d.SetId(id)
	return nil
}

func resourceOutscaleInboundRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	sg_id := d.Get("group_id").(string)
	sg, err := findResourceSecurityGroup(conn, sg_id)
	if _, notFound := err.(securityGroupNotFound); notFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error finding security group (%s) for rule (%s): %s", sg_id, d.Id(), err)
	}

	isVPC := sg.VpcId != nil && *sg.VpcId != ""

	var rule *fcu.IpPermission
	var rules []*fcu.IpPermission
	ruleType := "ingress"
	rules = sg.IpPermissions

	p, err := expandIPPerm(d, sg)
	if err != nil {
		return err
	}

	if len(rules) == 0 {
		fmt.Printf("\n\n[WARN] No %s rules were found for Security Group (%s) looking for Security Group Rule (%s)",
			ruleType, *sg.GroupName, d.Id())
		d.SetId("")
		return nil
	}

	rule = findRuleMatch(p, rules, isVPC)

	if rule == nil {
		fmt.Printf("\n\n[DEBUG] Unable to find matching %s Security Group Rule (%s) for Group %s",
			ruleType, d.Id(), sg_id)
		d.SetId("")
		return nil
	}

	fmt.Printf("\n\n[DEBUG] Found rule for Security Group Rule (%s): %s", d.Id(), rule)

	if err := setFromIPPerm(d, sg, p); err != nil {
		return errwrap.Wrapf("Error setting IP Permission for Security Group Rule: {{err}}", err)
	}

	return nil
}

func resourceOutscaleInboundRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	sg_id := d.Get("group_id").(string)

	awsMutexKV.Lock(sg_id)
	defer awsMutexKV.Unlock(sg_id)

	sg, err := findResourceSecurityGroup(conn, sg_id)
	if err != nil {
		return err
	}

	perm, err := expandIPPerm(d, sg)
	if err != nil {
		return err
	}
	fmt.Printf("\n\n[DEBUG] Revoking security group %#v %s rule: %#v",
		sg_id, "ingress", perm)
	req := &fcu.RevokeSecurityGroupIngressInput{
		GroupId:       sg.GroupId,
		IpPermissions: []*fcu.IpPermission{perm},
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.RevokeSecurityGroupIngress(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("\n\n[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf(
			"Error revoking security group %s rules: %s",
			sg_id, err)
	}

	d.SetId("")

	return nil
}
