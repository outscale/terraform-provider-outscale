package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLinPeeringConnectionAccepter() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinPeeringAccepterCreate,
		Read:   resourceLinPeeringRead,
		Delete: resourceLinPeeringAccepterDelete,
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
			"accepter_net": vpcPeeringConnectionOptionsSchema(),
			"source_net":   vpcPeeringConnectionOptionsSchema(),
			"tags":         tagsListSchemaComputed(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"accepter_net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLinPeeringAccepterCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).OSCAPI

	id := d.Get("net_peering_id").(string)
	d.SetId(id)

	req := oscgo.AcceptNetPeeringRequest{
		NetPeeringId: id,
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.NetPeeringApi.AcceptNetPeering(context.Background()).AcceptNetPeeringRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp.StatusCode, err)
		}
		return nil
	})

	var errString string
	if err != nil {
		errString = err.Error()
		return fmt.Errorf("Error creating Net Peering accepter. Details: %s", errString)
	}

	return resourceLinPeeringRead(d, meta)
}

func resourceLinPeeringAccepterDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Will not delete VPC peering connection. Terraform will remove this resource from the state file, however resources may remain.")
	d.SetId("")
	return nil
}
