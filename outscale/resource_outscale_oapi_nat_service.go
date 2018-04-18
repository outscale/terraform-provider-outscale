package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleOAPINatService() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPINatServiceCreate,
		Read:   resourceOAPINatServiceRead,
		Delete: resourceOAPINatServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// Arguments
			"reservation_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"token": {
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
			"public_ip": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reservation_id": {
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
			"nat_service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lin_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOAPINatServiceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	// Create the NAT Gateway
	createOpts := &fcu.CreateNatGatewayInput{
		AllocationId: aws.String(d.Get("reservation_id").(string)),
		SubnetId:     aws.String(d.Get("subnet_id").(string)),
	}

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
	fmt.Printf("\n\n[INFO] NAT Gateway ID: %s", d.Id())

	// Wait for the NAT Gateway to become available
	fmt.Printf("\n\n[DEBUG] Waiting for NAT Gateway (%s) to become available", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"available"},
		Refresh: NGOAPIStateRefreshFunc(conn, d.Id()),
		Timeout: 10 * time.Minute,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for NAT Gateway (%s) to become available: %s", d.Id(), err)
	}

	// Update our attributes and return
	return resourceOAPINatServiceRead(d, meta)
}

func resourceOAPINatServiceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	// Refresh the NAT Gateway state
	ngRaw, state, err := NGOAPIStateRefreshFunc(conn, d.Id())()
	if err != nil {
		return err
	}

	status := map[string]bool{
		"deleted":  true,
		"deleting": true,
		"failed":   true,
	}

	if _, ok := status[strings.ToLower(state)]; ngRaw == nil || ok {
		fmt.Printf("\n\n[INFO] Removing %s from Terraform state as it is not found or in the deleted state.", d.Id())
		d.SetId("")
		return nil
	}

	// Set NAT Gateway attributes
	ng := ngRaw.(*fcu.NatGateway)

	if ng.NatGatewayId != nil {
		d.Set("nat_service_id", *ng.NatGatewayId)
	}
	if ng.State != nil {
		d.Set("state", *ng.State)
	}
	if ng.SubnetId != nil {
		d.Set("subnet_id", *ng.SubnetId)
	}
	if ng.VpcId != nil {
		d.Set("vpc_id", *ng.VpcId)
	}

	if ng.NatGatewayAddresses != nil {
		addresses := make([]map[string]interface{}, len(ng.NatGatewayAddresses))

		for k, v := range ng.NatGatewayAddresses {
			address := make(map[string]interface{})
			if v.AllocationId != nil {
				address["reservation_id"] = *v.AllocationId
			}
			if v.PublicIp != nil {
				address["public_ip"] = *v.PublicIp
			}
			addresses[k] = address
		}
		if err := d.Set("public_ip", addresses); err != nil {
			return err
		}
	}

	return nil
}

func resourceOAPINatServiceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	deleteOpts := &fcu.DeleteNatGatewayInput{
		NatGatewayId: aws.String(d.Id()),
	}
	fmt.Printf("\n\n[INFO] Deleting NAT Gateway: %s", d.Id())

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
		Refresh:    NGOAPIStateRefreshFunc(conn, d.Id()),
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

// NGOAPIStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a NAT Gateway.
func NGOAPIStateRefreshFunc(conn *fcu.Client, id string) resource.StateRefreshFunc {
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
				return nil, "", nil
			} else {
				fmt.Printf("\n\nError on NGStateRefresh: %s", err)
				return nil, "", err
			}
		}

		ng := resp.NatGateways[0]
		return ng, *ng.State, nil
	}
}
