package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/outscale/osc-go/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleOAPILinPeeringConnectionAccepter() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILinPeeringAccepterCreate,
		Read:   resourceOutscaleOAPILinPeeringRead,
		Delete: resourceOutscaleOAPILinPeeringAccepterDelete,
		Update: resourceOutscaleOAPILinPeeringAccepterUpdate,
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
			"tags":         tagsListOAPISchema(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPILinPeeringAccepterCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	id := d.Get("net_peering_id").(string)
	d.SetId(id)

	req := &oapi.AcceptNetPeeringRequest{
		NetPeeringId: id,
	}

	var err error
	var resp *oapi.POST_AcceptNetPeeringResponses
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_AcceptNetPeering(*req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var errString string

	if err != nil || resp.OK == nil {
		if err != nil {
			errString = err.Error()
		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("Status Code: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("Status Code: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("Status: 500, %s", utils.ToJSONString(resp.Code500))
		}
		return fmt.Errorf("Error creating Net Peering accepter. Details: %s", errString)
	}

	if err := setOAPITags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

	return resourceOutscaleOAPILinPeeringRead(d, meta)
}

func resourceOutscaleOAPILinPeeringAccepterUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	d.Partial(true)

	if err := setOAPITags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

	d.Partial(false)
	return resourceOutscaleOAPILinPeeringRead(d, meta)
}

func resourceOutscaleOAPILinPeeringAccepterDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Will not delete VPC peering connection. Terraform will remove this resource from the state file, however resources may remain.")
	d.SetId("")
	return nil
}
