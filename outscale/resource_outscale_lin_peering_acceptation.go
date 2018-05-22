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

func resourceOutscaleLinPeeringConnectionAccepter() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleLinPeeringAccepterCreate,
		Read:   resourceOutscaleLinPeeringRead,
		Delete: resourceOutscaleLinPeeringAccepterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_peering_connection_id": {
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
			"accepter_vpc_info":  vpcPeeringConnectionOptionsSchema(),
			"requester_vpc_info": vpcPeeringConnectionOptionsSchema(),
			"tag_set":            tagsSchemaComputed(),
			"tag":                tagsSchema(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleLinPeeringAccepterCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	id := d.Get("vpc_peering_connection_id").(string)
	d.SetId(id)

	req := &fcu.AcceptVpcPeeringConnectionInput{
		VpcPeeringConnectionId: aws.String(id),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.AcceptVpcPeeringConnection(req)

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
	}
	d.SetPartial("tag_set")

	return resourceOutscaleLinPeeringRead(d, meta)
}

func resourceOutscaleLinPeeringAccepterDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Will not delete VPC peering connection. Terraform will remove this resource from the state file, however resources may remain.")
	d.SetId("")
	return nil
}
