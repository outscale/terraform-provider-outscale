package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOutscaleOAPIPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIPolicyRead,

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: false,
			},
			"is_linked": {
				Type:     schema.TypeBool,
				Required: false,
			},
			"path": {
				Type:     schema.TypeString,
				Required: false,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: false,
			},
			"resources_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policy_default_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_linkable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPIPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM

	groupName, groupNameOk := d.GetOk("group_name")
	isLinked, isLinkedOk := d.GetOk("is_linked")
	path, pathOk := d.GetOk("path")
	userName, userNameOk := d.GetOk("user_name")

	if groupNameOk == false && isLinkedOk == false && pathOk == false && userNameOk == false {
		return fmt.Errorf("At least one the following arguments must be provided: group_name, is_linked, path or user_name")
	}

	input := &eim.GetPolicyInput{}

	if groupNameOk == true {
		input.GroupName = aws.String(groupName.(string))
	}

	if isLinkedOk == true {
		input.IsLinked = aws.Bool(isLinked.(bool))
	}

	if pathOk == true {
		input.Path = aws.String(path.(string))
	}

	if userNameOk == true {
		input.UserName = aws.String(userName.(string))
	}

	var err error
	var output *eim.GetPolicyResult
	var rs *eim.GetPolicyOutput
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rs, err = conn.API.GetPolicy(input)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		if rs.GetPolicyResult != nil {
			output = rs.GetPolicyResult
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading policy with arn %s: %s", groupName, err)
	}

	d.Set("resources_count", aws.Int64Value(output.Policy.AttachmentCount))
	d.Set("policy_default_version_id", aws.StringValue(output.Policy.DefaultVersionId))
	d.Set("description", aws.StringValue(output.Policy.Description))
	d.Set("is_linkable", aws.BoolValue(output.Policy.IsAttachable))
	d.Set("path", aws.StringValue(output.Policy.Path))
	d.Set("policy_id", aws.StringValue(output.Policy.PolicyId))
	d.Set("policy_name", aws.StringValue(output.Policy.PolicyName))
	d.Set("request_id", aws.StringValue(rs.ResponseMetadata.RequestID))
	d.SetId(*output.Policy.PolicyId)

	return nil
}
