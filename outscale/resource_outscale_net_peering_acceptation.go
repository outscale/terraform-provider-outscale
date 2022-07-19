package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:     schema.TypeList,
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

func resourceOutscaleOAPILinPeeringAccepterCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	id := d.Get("net_peering_id").(string)
	d.SetId(id)

	req := oscgo.AcceptNetPeeringRequest{
		NetPeeringId: id,
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.NetPeeringApi.AcceptNetPeering(context.Background()).AcceptNetPeeringRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
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
