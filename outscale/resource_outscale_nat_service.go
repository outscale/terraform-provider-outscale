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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// Attributes
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
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNatServiceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	createOpts := &fcu.CreateNatGatewayInput{
		AllocationId: aws.String(d.Get("allocation_id").(string)),
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

	ng := natResp.NatGateway
	d.SetId(*ng.NatGatewayId)

	log.Printf("\n\n[DEBUG] Waiting for NAT Gateway (%s) to become available", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"available"},
		Refresh: ngStateRefreshFunc(conn, d.Id()),
		Timeout: 10 * time.Minute,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for NAT Gateway (%s) to become available: %s", d.Id(), err)
	}

	return resourceNatServiceRead(d, meta)
}

func resourceNatServiceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	// Refresh the NAT Gateway state
	ngRaw, state, err := ngStateRefreshFunc(conn, d.Id())()
	if err != nil {
		return err
	}

	status := map[string]bool{
		"deleted":  true,
		"deleting": true,
		"failed":   true,
	}

	if _, ok := status[strings.ToLower(state)]; ngRaw == nil || ok {
		log.Printf("\n\n[INFO] Removing %s from Terraform state as it is not found or in the deleted state.", d.Id())
		d.SetId("")
		return nil
	}

	opts := &fcu.DescribeNatGatewaysInput{
		NatGatewayIds: []*string{aws.String(d.Id())},
	}
	var resp *fcu.DescribeNatGatewaysOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
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

	nat := resp.NatGateways[0]

	d.Set("nat_gateway_id", aws.StringValue(nat.NatGatewayId))
	d.Set("state", aws.StringValue(nat.State))
	d.Set("subnet_id", aws.StringValue(nat.SubnetId))
	d.Set("vpc_id", aws.StringValue(nat.VpcId))

	addresses := make([]map[string]interface{}, len(nat.NatGatewayAddresses))
	if nat.NatGatewayAddresses != nil {
		for k, v := range nat.NatGatewayAddresses {
			address := make(map[string]interface{})
			address["allocation_id"] = aws.StringValue(v.AllocationId)
			address["public_ip"] = aws.StringValue(v.PublicIp)

			addresses[k] = address
		}
	}

	d.Set("request_id", resp.RequestId)

	return d.Set("nat_gateway_address", addresses)
}

func resourceNatServiceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	deleteOpts := &fcu.DeleteNatGatewayInput{
		NatGatewayId: aws.String(d.Id()),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
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
		if strings.Contains(fmt.Sprint(err), "NatGatewayNotFound:") {
			return nil
		}

		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{"deleted"},
		Refresh:    ngStateRefreshFunc(conn, d.Id()),
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

func ngStateRefreshFunc(conn *fcu.Client, ID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := &fcu.DescribeNatGatewaysInput{
			NatGatewayIds: []*string{aws.String(ID)},
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
			}
			log.Printf("\n\nError on NGStateRefresh: %s", err)
			return nil, "", err
		}

		ng := resp.NatGateways[0]
		return ng, *ng.State, nil
	}
}
