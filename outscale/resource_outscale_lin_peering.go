package outscale

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleLinPeeringConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleLinPeeringCreate,
		Read:   resourceOutscaleLinPeeringRead,
		Delete: resourceOutscaleLinPeeringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"peer_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"peer_vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_peering_connection_id": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceOutscaleLinPeeringCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	// Create the vpc peering connection
	createOpts := &fcu.CreateVpcPeeringConnectionInput{
		PeerVpcId: aws.String(d.Get("peer_vpc_id").(string)),
		VpcId:     aws.String(d.Get("vpc_id").(string)),
	}

	if v, ok := d.GetOk("peer_owner_id"); ok {
		createOpts.PeerOwnerId = aws.String(v.(string))
	}

	log.Printf("[DEBUG] VPC Peering Create options: %#v", createOpts)

	var resp *fcu.CreateVpcPeeringConnectionOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.CreateVpcPeeringConnection(createOpts)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errwrap.Wrapf("Error creating VPC Peering Connection: {{err}}", err)
	}

	// Get the ID and store it
	rt := resp.VpcPeeringConnection
	d.SetId(*rt.VpcPeeringConnectionId)

	if err := setTags(conn, d); err != nil {
		return err
	} else {
		d.SetPartial("tag_set")
	}

	log.Printf("[INFO] VPC Peering Connection ID: %s", d.Id())

	// Wait for the vpc peering connection to become available
	log.Printf("[DEBUG] Waiting for VPC Peering Connection (%s) to become available.", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"initiating-request", "provisioning", "pending"},
		Target:  []string{"pending-acceptance", "active"},
		Refresh: resourceOutscaleLinPeeringConnectionStateRefreshFunc(conn, d.Id()),
		Timeout: 1 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return errwrap.Wrapf(fmt.Sprintf(
			"Error waiting for VPC Peering Connection (%s) to become available: {{err}}",
			d.Id()), err)
	}

	return resourceOutscaleLinPeeringRead(d, meta)
}

func resourceOutscaleLinPeeringRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var resp *fcu.DescribeVpcPeeringConnectionsOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeVpcPeeringConnections(&fcu.DescribeVpcPeeringConnectionsInput{
			VpcPeeringConnectionIds: []*string{aws.String(d.Id())},
		})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidVpcPeeringConnectionID.NotFound") {
			resp = nil
		} else {
			log.Printf("Error reading VPC Peering Connection details: %s", err)
			return err
		}
	}

	pc := resp.VpcPeeringConnections[0]

	// Allow a failed VPC Peering Connection to fallthrough,
	// to allow rest of the logic below to do its work.
	if err != nil && *pc.Status.Code != "failed" {
		return err
	}

	if resp == nil {
		d.SetId("")
		return nil
	}

	// The failed status is a status that we can assume just means the
	// connection is gone. Destruction isn't allowed, and it eventually
	// just "falls off" the console. See GH-2322
	if pc.Status != nil {
		status := map[string]bool{
			"deleted":  true,
			"deleting": true,
			"expired":  true,
			"failed":   true,
			"rejected": true,
		}
		if _, ok := status[*pc.Status.Code]; ok {
			log.Printf("[DEBUG] VPC Peering Connection (%s) in state (%s), removing.",
				d.Id(), *pc.Status.Code)
			d.SetId("")
			return nil
		}
	}
	log.Printf("[DEBUG] VPC Peering Connection response: %#v", pc)

	log.Printf("[DEBUG] VPC PeerConn Requester %s, Accepter %s", *pc.RequesterVpcInfo.OwnerId, *pc.AccepterVpcInfo.OwnerId)

	accepter := make(map[string]interface{})
	requester := make(map[string]interface{})
	stat := make(map[string]interface{})

	if pc.AccepterVpcInfo != nil {
		accepter["cidr_block"] = aws.StringValue(pc.AccepterVpcInfo.CidrBlock)
		accepter["owner_id"] = aws.StringValue(pc.AccepterVpcInfo.OwnerId)
		accepter["vpc_id"] = aws.StringValue(pc.AccepterVpcInfo.VpcId)
	}
	if pc.RequesterVpcInfo != nil {
		requester["cidr_block"] = aws.StringValue(pc.AccepterVpcInfo.CidrBlock)
		requester["owner_id"] = aws.StringValue(pc.AccepterVpcInfo.OwnerId)
		requester["vpc_id"] = aws.StringValue(pc.AccepterVpcInfo.VpcId)
	}
	if pc.Status != nil {
		stat["code"] = aws.StringValue(pc.Status.Code)
		stat["message"] = aws.StringValue(pc.Status.Message)
	}

	if err := d.Set("accepter_vpc_info", accepter); err != nil {
		return err
	}
	if err := d.Set("requester_vpc_info", requester); err != nil {
		return err
	}
	if err := d.Set("status", stat); err != nil {
		return err
	}
	if err := d.Set("vpc_peering_connection_id", pc.VpcPeeringConnectionId); err != nil {
		return err
	}
	if err := d.Set("tag_set", tagsToMap(pc.Tags)); err != nil {
		return errwrap.Wrapf("Error setting VPC Peering Connection tags: {{err}}", err)
	}

	d.Set("request_id", resp.RequestId)

	return nil
}

func resourceVPCPeeringConnectionAccept(conn *fcu.Client, id string) (string, error) {
	log.Printf("[INFO] Accept VPC Peering Connection with ID: %s", id)

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
		return "", err
	}
	pc := resp.VpcPeeringConnection

	return *pc.Status.Code, nil
}

func resourceOutscaleLinPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.VM.DeleteVpcPeeringConnection(
			&fcu.DeleteVpcPeeringConnectionInput{
				VpcPeeringConnectionId: aws.String(d.Id()),
			})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	return err
}

// resourceOutscaleLinPeeringConnectionStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a VPCPeeringConnection.
func resourceOutscaleLinPeeringConnectionStateRefreshFunc(conn *fcu.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		var resp *fcu.DescribeVpcPeeringConnectionsOutput
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeVpcPeeringConnections(&fcu.DescribeVpcPeeringConnectionsInput{
				VpcPeeringConnectionIds: []*string{aws.String(id)},
			})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidVpcPeeringConnectionID.NotFound") {
				resp = nil
			} else {
				log.Printf("Error reading VPC Peering Connection details: %s", err)
				return nil, "error", err
			}
		}

		if resp == nil {
			// Sometimes AWS just has consistency issues and doesn't see
			// our instance yet. Return an empty state.
			return nil, "", nil
		}

		pc := resp.VpcPeeringConnections[0]

		// A VPC Peering Connection can exist in a failed state due to
		// incorrect VPC ID, account ID, or overlapping IP address range,
		// thus we short circuit before the time out would occur.
		if pc != nil && *pc.Status.Code == "failed" {
			return nil, "failed", errors.New(*pc.Status.Message)
		}

		return pc, *pc.Status.Code, nil
	}
}

func vpcPeeringConnectionOptionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cidr_block": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"owner_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"vpc_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}
