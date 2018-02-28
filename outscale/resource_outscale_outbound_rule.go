package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		// Create: resourceAwsSecurityGroupRuleCreate,
		// Read:   resourceAwsSecurityGroupRuleRead,
		// Delete: resourceAwsSecurityGroupRuleDelete,

		Schema: map[string]*schema.Schema{
			"cidr_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"from_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
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
				Computed: true,
				ForceNew: true,
			},
			"source_security_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"source_security_group_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"to_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ip_permissions": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"groups": {
							Type:     schema.TypeMap,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
									"group_name": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
									"user_id": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
								},
							},
						},
						"to_port": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"ip_protocol": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"ip_ranges": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_ip": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"prefix_list_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"prefix_list_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
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

var awsMutexKV = mutexkv.NewMutexKV()

func resourceAwsSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
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

	if err := validateAwsSecurityGroupRule(d); err != nil {
		return err
	}

	ruleType := d.Get("type").(string)
	isVPC := sg.VpcId != nil && *sg.VpcId != ""

	var autherr error
	switch ruleType {
	case "ingress":
		log.Printf("[DEBUG] Authorizing security group %s %s rule: %s",
			sg_id, "Ingress", perm)

		req := &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId:       sg.GroupId,
			IpPermissions: []*ec2.IpPermission{perm},
		}

		if !isVPC {
			req.GroupId = nil
			req.GroupName = sg.GroupName
		}

		_, autherr = conn.AuthorizeSecurityGroupIngress(req)

	case "egress":
		log.Printf("[DEBUG] Authorizing security group %s %s rule: %#v",
			sg_id, "Egress", perm)

		req := &ec2.AuthorizeSecurityGroupEgressInput{
			GroupId:       sg.GroupId,
			IpPermissions: []*ec2.IpPermission{perm},
		}

		_, autherr = conn.AuthorizeSecurityGroupEgress(req)

	default:
		return fmt.Errorf("Security Group Rule must be type 'ingress' or type 'egress'")
	}

	if autherr != nil {
		if awsErr, ok := autherr.(awserr.Error); ok {
			if awsErr.Code() == "InvalidPermission.Duplicate" {
				return fmt.Errorf(`[WARN] A duplicate Security Group rule was found on (%s). This may be
a side effect of a now-fixed Terraform issue causing two security groups with
identical attributes but different source_security_group_ids to overwrite each
other in the state. See https://github.com/hashicorp/terraform/pull/2376 for more
information and instructions for recovery. Error message: %s`, sg_id, awsErr.Message())
			}
		}

		return fmt.Errorf(
			"Error authorizing security group rule type %s: %s",
			ruleType, autherr)
	}

	id := ipPermissionIDHash(sg_id, ruleType, perm)
	log.Printf("[DEBUG] Computed group rule ID %s", id)

	retErr := resource.Retry(5*time.Minute, func() *resource.RetryError {
		sg, err := findResourceSecurityGroup(conn, sg_id)

		if err != nil {
			log.Printf("[DEBUG] Error finding Security Group (%s) for Rule (%s): %s", sg_id, id, err)
			return resource.NonRetryableError(err)
		}

		var rules []*ec2.IpPermission
		switch ruleType {
		case "ingress":
			rules = sg.IpPermissions
		default:
			rules = sg.IpPermissionsEgress
		}

		rule := findRuleMatch(perm, rules, isVPC)

		if rule == nil {
			log.Printf("[DEBUG] Unable to find matching %s Security Group Rule (%s) for Group %s",
				ruleType, id, sg_id)
			return resource.RetryableError(fmt.Errorf("No match found"))
		}

		log.Printf("[DEBUG] Found rule for Security Group Rule (%s): %s", id, rule)
		return nil
	})

	if retErr != nil {
		return fmt.Errorf("Error finding matching %s Security Group Rule (%s) for Group %s",
			ruleType, id, sg_id)
	}

	d.SetId(id)
	return nil
}

func findResourceSecurityGroup(conn *fcu.Client, id string) (*fcu.SecurityGroup, error) {
	req := &fcu.DescribeSecurityGroupsInput{
		GroupIds: []*string{aws.String(id)},
	}
	resp, err := conn.VM.DescribeSecurityGroups(req)
	if strings.Contains(fmt.Sprint(err), "InvalidGroup.NotFound") {
		return nil, securityGroupNotFound{id, nil}
	}
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, securityGroupNotFound{id, nil}
	}
	if len(resp.SecurityGroups) != 1 || resp.SecurityGroups[0] == nil {
		return nil, securityGroupNotFound{id, resp.SecurityGroups}
	}

	return resp.SecurityGroups[0], nil
}
