package outscale

import (
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleOAPILinPeeringConnectionAccepter() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILinPeeringAccepterCreate,
		Read:   resourceOutscaleOAPILinPeeringRead,
		Delete: resourceOutscaleOAPILinPeeringAccepterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"lin_peering_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"accepter_lin": vpcPeeringConnectionOptionsSchema(),
			"source_lin":   vpcPeeringConnectionOptionsSchema(),
			"tag":          tagsSchema(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPILinPeeringAccepterCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	id := d.Get("lin_peering_id").(string)
	d.SetId(id)

	req := &fcu.AcceptVpcPeeringConnectionInput{
		VpcPeeringConnectionId: aws.String(id),
	}

	var resp *fcu.AcceptVpcPeeringConnectionOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.AcceptVpcPeeringConnection(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errwrap.Wrapf("Error creating VPC Peering Connection accepter: {{err}}", err)
	}

	if err := setTags(conn, d); err != nil {
		return err
	} else {
		d.SetPartial("tag_set")
	}

	return resourceOutscaleOAPILinPeeringRead(d, meta)
}

func resourceOutscaleOAPILinPeeringAccepterDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Will not delete VPC peering connection. Terraform will remove this resource from the state file, however resources may remain.")
	d.SetId("")
	return nil
}
