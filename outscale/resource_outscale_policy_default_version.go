package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func resourceOutscalePolicyDefaultVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscalePolicyDefaultVersionCreate,
		Read:   resourceOutscalePolicyDefaultVersionRead,
		Delete: resourceOutscalePolicyDefaultVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_arn": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscalePolicyDefaultVersionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM

	var err error

	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.API.SetDefaultPolicyVersion(&eim.SetDefaultPolicyVersionInput{
			PolicyArn: aws.String(d.Get("policy_arn").(string)),
			VersionId: aws.String(d.Get("version_id").(string)),
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error setting default policy %s", err)
	}

	d.SetId(d.Get("policy_arn").(string))

	return resourceOutscalePolicyDefaultVersionRead(d, meta)
}

func resourceOutscalePolicyDefaultVersionRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOutscalePolicyDefaultVersionDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	return nil
}
