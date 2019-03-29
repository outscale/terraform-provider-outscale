package outscale

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleOAPICustomerGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPICustomerGatewayCreate,
		Read:   resourceOutscaleOAPICustomerGatewayRead,
		Delete: resourceOutscaleOAPICustomerGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bgp_asn": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"public_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"client_endpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchemaComputed(),
			"tag":  tagsSchema(),
		},
	}
}

func resourceOutscaleOAPICustomerGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	ipAddress := d.Get("public_ip").(string)
	vpnType := d.Get("type").(string)
	bgpAsn := d.Get("bgp_asn").(int)

	alreadyExists, err := resourceOutscaleOAPICustomerGatewayExists(vpnType, ipAddress, bgpAsn, conn)
	if err != nil {
		return err
	}

	if alreadyExists {
		return fmt.Errorf("An existing customer gateway for IpAddress: %s, VpnType: %s, BGP ASN: %d has been found", ipAddress, vpnType, bgpAsn)
	}

	createOpts := &fcu.CreateCustomerGatewayInput{
		BgpAsn:   aws.Int64(int64(bgpAsn)),
		PublicIp: aws.String(ipAddress),
		Type:     aws.String(vpnType),
	}

	// Create the Customer Gateway.
	fmt.Printf("[DEBUG] Creating customer gateway")

	var resp *fcu.CreateCustomerGatewayOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.CreateCustomerGateway(createOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error creating customer gateway: %s", err)
	}

	// Store the ID
	customerGateway := resp.CustomerGateway
	d.SetId(*customerGateway.CustomerGatewayId)
	fmt.Printf("[INFO] Customer gateway ID: %s", *customerGateway.CustomerGatewayId)

	// Wait for the CustomerGateway to be available.
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available"},
		Refresh:    customerOAPIGatewayRefreshFunc(conn, *customerGateway.CustomerGatewayId),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for customer gateway (%s) to become ready: %s",
			*customerGateway.CustomerGatewayId, err)
	}

	// Create tags.
	if err := setTags(conn, d); err != nil {
		return err
	}

	t := make([]map[string]interface{}, 0)

	d.Set("tags", t)

	return nil
}

func customerOAPIGatewayRefreshFunc(conn *fcu.Client, gatewayID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		gatewayFilter := &fcu.Filter{
			Name:   aws.String("customer-gateway-id"),
			Values: []*string{aws.String(gatewayID)},
		}

		var resp *fcu.DescribeCustomerGatewaysOutput
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeCustomerGateways(&fcu.DescribeCustomerGatewaysInput{
				Filters: []*fcu.Filter{gatewayFilter},
			})

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidCustomerGatewayID.NotFound") {
				resp = nil
			} else {
				fmt.Printf("Error on CustomerGatewayRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil || len(resp.CustomerGateways) == 0 {
			// handle consistency issues
			return nil, "", nil
		}

		gateway := resp.CustomerGateways[0]
		return gateway, *gateway.State, nil
	}
}

func resourceOutscaleOAPICustomerGatewayExists(vpnType, ipAddress string, bgpAsn int, conn *fcu.Client) (bool, error) {
	ipAddressFilter := &fcu.Filter{
		Name:   aws.String("ip-address"),
		Values: []*string{aws.String(ipAddress)},
	}

	typeFilter := &fcu.Filter{
		Name:   aws.String("type"),
		Values: []*string{aws.String(vpnType)},
	}

	bgp := strconv.Itoa(bgpAsn)
	bgpAsnFilter := &fcu.Filter{
		Name:   aws.String("bgp-asn"),
		Values: []*string{aws.String(bgp)},
	}

	var resp *fcu.DescribeCustomerGatewaysOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeCustomerGateways(&fcu.DescribeCustomerGatewaysInput{
			Filters: []*fcu.Filter{ipAddressFilter, typeFilter, bgpAsnFilter},
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return false, err
	}

	if len(resp.CustomerGateways) > 0 && *resp.CustomerGateways[0].State != "deleted" {
		return true, nil
	}

	return false, nil
}

func resourceOutscaleOAPICustomerGatewayRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	gatewayFilter := &fcu.Filter{
		Name:   aws.String("customer-gateway-id"),
		Values: []*string{aws.String(d.Id())},
	}

	var resp *fcu.DescribeCustomerGatewaysOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeCustomerGateways(&fcu.DescribeCustomerGatewaysInput{
			Filters: []*fcu.Filter{gatewayFilter},
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidCustomerGatewayID.NotFound") {
			d.SetId("")
			return nil
		}
		fmt.Printf("[ERROR] Error finding CustomerGateway: %s", err)
		return err
	}

	if len(resp.CustomerGateways) != 1 {
		return fmt.Errorf("[ERROR] Error finding CustomerGateway: %s", d.Id())
	}

	if *resp.CustomerGateways[0].State == "deleted" {
		fmt.Printf("[INFO] Customer Gateway is in `deleted` state: %s", d.Id())
		d.SetId("")
		return nil
	}

	customerGateway := resp.CustomerGateways[0]
	d.Set("public_ip", customerGateway.IpAddress)
	d.Set("type", customerGateway.Type)
	d.Set("tags", tagsToMap(customerGateway.Tags))

	if *customerGateway.BgpAsn != "" {
		val, err := strconv.ParseInt(*customerGateway.BgpAsn, 0, 0)
		if err != nil {
			return fmt.Errorf("error parsing bgp_asn: %s", err)
		}

		d.Set("bgp_asn", int(val))
	}

	return nil
}

func resourceOutscaleOAPICustomerGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var err error
	err = resource.Retry(15*time.Minute, func() *resource.RetryError {
		_, err := conn.VM.DeleteCustomerGateway(&fcu.DeleteCustomerGatewayInput{
			CustomerGatewayId: aws.String(d.Id()),
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidCustomerGatewayID.NotFound") {
			d.SetId("")
			return nil
		}
		fmt.Printf("[ERROR] Error deleting CustomerGateway: %s", err)
		return err
	}

	gatewayFilter := &fcu.Filter{
		Name:   aws.String("customer-gateway-id"),
		Values: []*string{aws.String(d.Id())},
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.VM.DescribeCustomerGateways(&fcu.DescribeCustomerGatewaysInput{
			Filters: []*fcu.Filter{gatewayFilter},
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidCustomerGatewayID.NotFound") {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if len(resp.CustomerGateways) != 1 {
			return resource.RetryableError(fmt.Errorf("[ERROR] Error finding CustomerGateway for delete: %s", d.Id()))
		}

		switch *resp.CustomerGateways[0].State {
		case "pending", "available", "deleting":
			return resource.RetryableError(fmt.Errorf("[DEBUG] Gateway (%s) in state (%s), retrying", d.Id(), *resp.CustomerGateways[0].State))
		case "deleted":
			return nil
		default:
			return resource.RetryableError(fmt.Errorf("[DEBUG] Unrecognized state (%s) for Customer Gateway delete on (%s)", *resp.CustomerGateways[0].State, d.Id()))
		}
	})

	if err != nil {
		return err
	}

	return nil
}
