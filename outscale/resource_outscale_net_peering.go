package outscale

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOutscaleOAPILinPeeringConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILinPeeringCreate,
		Read:   resourceOutscaleOAPILinPeeringRead,
		Update: resourceOutscaleOAPINetPeeringUpdate,
		Delete: resourceOutscaleOAPILinPeeringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"source_net_account_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"accepter_net_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_net_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"net_peering_id": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceOutscaleOAPILinPeeringCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	// Create the vpc peering connection
	createOpts := oscgo.CreateNetPeeringRequest{
		AccepterNetId: d.Get("accepter_net_id").(string),
		SourceNetId:   d.Get("source_net_id").(string),
	}

	log.Printf("[DEBUG] VPC Peering Create options: %#v", createOpts)

	var resp oscgo.CreateNetPeeringResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.NetPeeringApi.CreateNetPeering(context.Background()).CreateNetPeeringRequest(createOpts).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	var errString string
	if err != nil {
		errString = err.Error()
		return fmt.Errorf("Error creating Net Peering. Details: %s", errString)
	}

	// Get the ID and store it
	d.SetId(resp.NetPeering.GetNetPeeringId())

	//SetTags
	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.(*schema.Set), d.Id(), conn)
		if err != nil {
			return err
		}
	}

	log.Printf("[INFO] Net Peering ID: %s", d.Id())

	// Wait for the vpc peering connection to become available
	log.Printf("[DEBUG] Waiting for Net Peering (%s) to become available.", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"initiating-request", "provisioning", "pending"},
		Target:  []string{"pending-acceptance", "active"},
		Refresh: resourceOutscaleOAPILinPeeringConnectionStateRefreshFunc(conn, d.Id()),
		Timeout: 1 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return errwrap.Wrapf(fmt.Sprintf(
			"Error waiting for Net Peering (%s) to become available: {{err}}",
			d.Id()), err)
	}

	return resourceOutscaleOAPILinPeeringRead(d, meta)
}

func resourceOutscaleOAPILinPeeringRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	var resp oscgo.ReadNetPeeringsResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.NetPeeringApi.ReadNetPeerings(context.Background()).ReadNetPeeringsRequest(oscgo.ReadNetPeeringsRequest{
			Filters: &oscgo.FiltersNetPeering{NetPeeringIds: &[]string{d.Id()}},
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	var errString string
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidVpcPeeringConnectionID.NotFound") {
			d.SetId("")
			return nil
		}
		errString = err.Error()
		return fmt.Errorf("Error reading Net Peering details: %s", errString)
	}
	if utils.IsResponseEmpty(len(resp.GetNetPeerings()), "NetPeering", d.Id()) {
		d.SetId("")
		return nil
	}
	pc := resp.GetNetPeerings()[0]

	// The failed status is a status that we can assume just means the
	// connection is gone. Destruction isn't allowed, and it eventually
	// just "falls off" the console. See GH-2322
	if !reflect.DeepEqual(pc.State, oscgo.NetPeeringState{}) {
		status := map[string]bool{
			"deleted":  true,
			"deleting": true,
			"expired":  true,
			"failed":   true,
			"rejected": true,
		}
		if _, ok := status[pc.State.GetName()]; ok {
			log.Printf("[DEBUG] Net Peering (%s) in state (%s), removing.",
				d.Id(), pc.State.GetName())
			d.SetId("")
			return nil
		}
	}
	log.Printf("[DEBUG] Net Peering response: %#v", pc)

	log.Printf("[DEBUG] VPC PeerConn Source %s, Accepter %s", pc.SourceNet.GetAccountId(), pc.AccepterNet.GetAccountId())

	accepter := make(map[string]interface{})
	requester := make(map[string]interface{})
	stat := make(map[string]interface{})

	if !reflect.DeepEqual(pc.GetAccepterNet(), oscgo.AccepterNet{}) {
		accepter["ip_range"] = pc.AccepterNet.GetIpRange()
		accepter["account_id"] = pc.AccepterNet.GetAccountId()
		accepter["net_id"] = pc.AccepterNet.GetNetId()
	}
	if !reflect.DeepEqual(pc.GetSourceNet(), oscgo.SourceNet{}) {
		requester["ip_range"] = pc.SourceNet.GetIpRange()
		requester["account_id"] = pc.SourceNet.GetAccountId()
		requester["net_id"] = pc.SourceNet.GetNetId()
	}
	if pc.State.GetName() != "" {
		stat["name"] = pc.State.GetName()
		stat["message"] = pc.State.GetMessage()
	}

	if err := d.Set("accepter_net_id", pc.GetAccepterNet().NetId); err != nil {
		return err
	}
	if err := d.Set("source_net_id", pc.GetSourceNet().NetId); err != nil {
		return err
	}
	if err := d.Set("accepter_net", accepter); err != nil {
		return err
	}
	if err := d.Set("source_net", requester); err != nil {
		return err
	}
	if err := d.Set("state", stat); err != nil {
		return err
	}
	if err := d.Set("net_peering_id", pc.GetNetPeeringId()); err != nil {
		return err
	}
	if err := d.Set("tags", tagsOSCAPIToMap(pc.GetTags())); err != nil {
		return errwrap.Wrapf("Error setting Net Peering tags: {{err}}", err)
	}

	return nil
}

func resourceOutscaleOAPINetPeeringUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	d.Partial(true)

	if err := setOSCAPITags(conn, d, "tags"); err != nil {
		return err
	}

	d.SetPartial("tags")

	d.Partial(false)
	return resourceOutscaleOAPILinPeeringRead(d, meta)
}

func resourceOutscaleOAPILinPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.NetPeeringApi.DeleteNetPeering(context.Background()).DeleteNetPeeringRequest(oscgo.DeleteNetPeeringRequest{
			NetPeeringId: d.Id(),
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	var errString string
	if err != nil {
		errString = err.Error()
		return fmt.Errorf("Error deleteting Net Peering. Details: %s", errString)
	}

	return nil
}

// resourceOutscaleOAPILinPeeringConnectionStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a VPCPeeringConnection.
func resourceOutscaleOAPILinPeeringConnectionStateRefreshFunc(conn *oscgo.APIClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp oscgo.ReadNetPeeringsResponse
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			rp, httpResp, err := conn.NetPeeringApi.ReadNetPeerings(context.Background()).ReadNetPeeringsRequest(oscgo.ReadNetPeeringsRequest{
				Filters: &oscgo.FiltersNetPeering{NetPeeringIds: &[]string{id}},
			}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		var errString string
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidVpcPeeringConnectionID.NotFound") {
				// Sometimes AWS just has consistency issues and doesn't see
				// our instance yet. Return an empty state.
				return nil, "", nil
			}
			errString = err.Error()
			return nil, "error", fmt.Errorf("Error reading Net Peering details: %s", errString)
		}

		pc := resp.GetNetPeerings()[0]

		// A Net Peering can exist in a failed state due to
		// incorrect VPC ID, account ID, or overlapping IP address range,
		// thus we short circuit before the time out would occur.
		if pc.State.GetName() == "failed" {
			return nil, "failed", errors.New(pc.State.GetMessage())
		}

		return pc, pc.State.GetName(), nil
	}
}

func vpcOAPIPeeringConnectionOptionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ip_range": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"account_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"net_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}
