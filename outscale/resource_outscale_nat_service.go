package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleNatService() *schema.Resource {
	return &schema.Resource{
		Create: resourceNatServiceCreate,
		Read:   resourceNatServiceRead,
		Delete: resourceNatServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// Arguments
			"allocation_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_token": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// Attributes
			"nat_gateway": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nat_gateway_address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allocation_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"nat_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnetId": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceNatServiceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	// Create the NAT Gateway
	createOpts := &fcu.CreateNatGatewayInput{
		AllocationId: aws.String(d.Get("allocation_id").(string)),
		SubnetId:     aws.String(d.Get("subnet_id").(string)),
	}

	log.Printf("[DEBUG] Create NAT Gateway: %s", *createOpts)

	var natResp *fcu.CreateNatGatewayOutput

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		natResp, err = conn.VM.CreateNatGateway(createOpts)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {

		return fmt.Errorf("Error creating NAT Gateway: %s", err)
	}

	// Get the ID and store it
	ng := natResp.NatGateway
	d.SetId(*ng.NatGatewayId)
	log.Printf("[INFO] NAT Gateway ID: %s", d.Id())

	// Wait for the NAT Gateway to become available
	log.Printf("[DEBUG] Waiting for NAT Gateway (%s) to become available", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"available"},
		Refresh: NGStateRefreshFunc(conn, d.Id()),
		Timeout: 10 * time.Minute,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for NAT Gateway (%s) to become available: %s", d.Id(), err)
	}

	// Update our attributes and return
	return resourceNatServiceRead(d, meta)
}

func resourceNatServiceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	// Refresh the NAT Gateway state
	ngRaw, state, err := NGStateRefreshFunc(conn, d.Id())()
	if err != nil {
		return err
	}

	status := map[string]bool{
		"deleted":  true,
		"deleting": true,
		"failed":   true,
	}

	if _, ok := status[strings.ToLower(state)]; ngRaw == nil || ok {
		log.Printf("[INFO] Removing %s from Terraform state as it is not found or in the deleted state.", d.Id())
		d.SetId("")
		return nil
	}

	// Set NAT Gateway attributes
	ng := ngRaw.(*fcu.NatGateway)
	ngGateway := make(map[string]interface{})
	if ng.NatGatewayAddresses != nil && len(ng.NatGatewayAddresses) > 0 {
		addresses := make([]map[string]interface{}, len(ng.NatGatewayAddresses))

		for k, v := range ng.NatGatewayAddresses {
			address := make(map[string]interface{})
			if v.AllocationId != nil {
				address["allocation_id"] = *v.AllocationId
			}
			if v.PublicIp != nil {
				address["public_ip"] = *v.PublicIp
			}
			addresses[k] = address
		}
		ngGateway["nat_gateway_address"] = addresses
	}
	if ng.NatGatewayId != nil {
		ngGateway["nat_gateway_id"] = *ng.NatGatewayId
	}
	if ng.State != nil {
		ngGateway["state"] = *ng.State
	}
	if ng.SubnetId != nil {
		ngGateway["subnet_id"] = *ng.SubnetId
	}
	if ng.VpcId != nil {
		ngGateway["vpc_id"] = *ng.VpcId
	}

	return d.Set("nat_gateway", ngGateway)
}

func resourceNatServiceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	deleteOpts := &fcu.DeleteNatGatewayInput{
		NatGatewayId: aws.String(d.Id()),
	}
	log.Printf("[INFO] Deleting NAT Gateway: %s", d.Id())

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		_, err = conn.VM.DeleteNatGateway(deleteOpts)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "NatGatewayNotFound:") {
			return nil
		}

		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{"deleted"},
		Refresh:    NGStateRefreshFunc(conn, d.Id()),
		Timeout:    30 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf("Error waiting for NAT Gateway (%s) to delete: %s", d.Id(), err)
	}

	return nil
}

// NGStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a NAT Gateway.
func NGStateRefreshFunc(conn *fcu.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := &fcu.DescribeNatGatewaysInput{
			NatGatewayIds: []*string{aws.String(id)},
		}
		var resp *fcu.DescribeNatGatewaysOutput
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error

			resp, err = conn.VM.DescribeNatGateways(opts)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "NatGatewayNotFound") {
				resp = nil
			} else {
				log.Printf("Error on NGStateRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil {
			// Sometimes AWS just has consistency issues and doesn't see
			// our instance yet. Return an empty state.
			return nil, "", nil
		}

		ng := resp.NatGateways[0]
		return ng, *ng.State, nil
	}
}
