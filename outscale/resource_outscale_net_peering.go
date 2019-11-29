package outscale

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/outscale/osc-go/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleOAPILinPeeringConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILinPeeringCreate,
		Read:   resourceOutscaleOAPILinPeeringRead,
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
			"tags":         tagsSchemaComputed(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPILinPeeringCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	// Create the vpc peering connection
	createOpts := &oapi.CreateNetPeeringRequest{
		AccepterNetId: d.Get("accepter_net_id").(string),
		SourceNetId:   d.Get("source_net_id").(string),
	}

	log.Printf("[DEBUG] VPC Peering Create options: %#v", createOpts)

	var resp *oapi.POST_CreateNetPeeringResponses
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_CreateNetPeering(*createOpts)

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
		return fmt.Errorf("Error creating Net Peering. Details: %s", errString)
	}

	// Get the ID and store it
	rt := resp.OK
	d.SetId(rt.NetPeering.NetPeeringId)

	if err := setOAPITags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

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
	conn := meta.(*OutscaleClient).OAPI
	var resp *oapi.POST_ReadNetPeeringsResponses
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_ReadNetPeerings(oapi.ReadNetPeeringsRequest{
			Filters: oapi.FiltersNetPeering{NetPeeringIds: []string{d.Id()}},
		})

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
			if strings.Contains(fmt.Sprint(err), "InvalidVpcPeeringConnectionID.NotFound") {
				d.SetId("")
				return nil
			}

			// Allow a failed Net Peering to fallthrough,
			// to allow rest of the logic below to do its work.
			//TODO: improve logic
			//FIXME: check if it is Name or Message
			if resp.OK != nil {
				if err != nil && resp.OK.NetPeerings[0].State.Name != "failed" {
					return err
				}
			}
			errString = err.Error()
		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("Status Code: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("Status Code: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("Status: 500, %s", utils.ToJSONString(resp.Code500))
		}
		return fmt.Errorf("Error reading Net Peering details: %s", errString)
	}

	result := resp.OK

	pc := result.NetPeerings[0]

	// The failed status is a status that we can assume just means the
	// connection is gone. Destruction isn't allowed, and it eventually
	// just "falls off" the console. See GH-2322
	if !reflect.DeepEqual(pc.State, oapi.NetPeeringState{}) {
		status := map[string]bool{
			"deleted":  true,
			"deleting": true,
			"expired":  true,
			"failed":   true,
			"rejected": true,
		}
		if _, ok := status[pc.State.Name]; ok {
			log.Printf("[DEBUG] Net Peering (%s) in state (%s), removing.",
				d.Id(), pc.State.Name)
			d.SetId("")
			return nil
		}
	}
	log.Printf("[DEBUG] Net Peering response: %#v", pc)

	log.Printf("[DEBUG] VPC PeerConn Source %s, Accepter %s", pc.SourceNet.AccountId, pc.AccepterNet.AccountId)

	accepter := make(map[string]interface{})
	requester := make(map[string]interface{})
	stat := make(map[string]interface{})

	if !reflect.DeepEqual(pc.AccepterNet, oapi.AccepterNet{}) {
		accepter["ip_range"] = pc.AccepterNet.IpRange
		accepter["account_id"] = pc.AccepterNet.AccountId
		accepter["net_id"] = pc.AccepterNet.NetId
	}
	if !reflect.DeepEqual(pc.SourceNet, oapi.SourceNet{}) {
		requester["ip_range"] = pc.SourceNet.IpRange
		requester["account_id"] = pc.SourceNet.AccountId
		requester["net_id"] = pc.SourceNet.NetId
	}
	if pc.State.Name != "" {
		stat["name"] = pc.State.Name
		stat["message"] = pc.State.Message
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
	if err := d.Set("net_peering_id", pc.NetPeeringId); err != nil {
		return err
	}
	if err := d.Set("tags", tagsOAPIToMap(pc.Tags)); err != nil {
		return errwrap.Wrapf("Error setting Net Peering tags: {{err}}", err)
	}

	d.Set("request_id", result.ResponseContext.RequestId)

	return nil
}

func resourceOutscaleOAPILinPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	var err error
	var resp *oapi.POST_DeleteNetPeeringResponses
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_DeleteNetPeering(oapi.DeleteNetPeeringRequest{
			NetPeeringId: d.Id(),
		})

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
		return fmt.Errorf("Error deleteting Net Peering. Details: %s", errString)
	}

	return nil
}

// resourceOutscaleOAPILinPeeringConnectionStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a VPCPeeringConnection.
func resourceOutscaleOAPILinPeeringConnectionStateRefreshFunc(conn *oapi.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		var resp *oapi.POST_ReadNetPeeringsResponses
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.POST_ReadNetPeerings(oapi.ReadNetPeeringsRequest{
				Filters: oapi.FiltersNetPeering{NetPeeringIds: []string{id}},
			})

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
				if strings.Contains(fmt.Sprint(err), "InvalidVpcPeeringConnectionID.NotFound") {
					// Sometimes AWS just has consistency issues and doesn't see
					// our instance yet. Return an empty state.
					return nil, "", nil
				}
				errString = err.Error()
			} else if resp.Code401 != nil {
				errString = fmt.Sprintf("Status Code: 401, %s", utils.ToJSONString(resp.Code401))
			} else if resp.Code400 != nil {
				errString = fmt.Sprintf("Status Code: 400, %s", utils.ToJSONString(resp.Code400))
			} else if resp.Code500 != nil {
				errString = fmt.Sprintf("Status: 500, %s", utils.ToJSONString(resp.Code500))
			}
			return nil, "error", fmt.Errorf("Error reading Net Peering details: %s", errString)
		}

		result := resp.OK

		pc := result.NetPeerings[0]

		// A Net Peering can exist in a failed state due to
		// incorrect VPC ID, account ID, or overlapping IP address range,
		// thus we short circuit before the time out would occur.
		if pc.State.Name == "failed" {
			return nil, "failed", errors.New(pc.State.Message)
		}

		return pc, pc.State.Name, nil
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
