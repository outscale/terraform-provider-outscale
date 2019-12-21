package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"net_peering_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"state": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
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
			"accepter_net": vpcOAPIPeeringConnectionOptionsSchema(),
			"source_net":   vpcOAPIPeeringConnectionOptionsSchema(),
			"tags":         tagsOAPIListSchemaComputed(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPILinPeeringAccepterCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	id := d.Get("net_peering_id").(string)
	d.SetId(id)

	req := oscgo.AcceptNetPeeringRequest{
		NetPeeringId: id,
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.NetPeeringApi.AcceptNetPeering(context.Background(), &oscgo.AcceptNetPeeringOpts{AcceptNetPeeringRequest: optional.NewInterface(req)})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var errString string

	if err != nil {
		errString = err.Error()
		return fmt.Errorf("Error creating Net Peering accepter. Details: %s", errString)
	}

	return resourceOutscaleOAPILinPeeringRead(d, meta)
}

func resourceOutscaleOAPILinPeeringAccepterDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Will not delete VPC peering connection. Terraform will remove this resource from the state file, however resources may remain.")
	d.SetId("")
	return nil
}
