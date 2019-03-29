package outscale

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/dl"
)

func resourceOutscaleDirectLinkInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleDirectLinkInterfaceCreate,
		Read:   resourceOutscaleDirectLinkInterfaceRead,
		Delete: resourceOutscaleDirectLinkInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceOutscaleDirectLinkInterfaceImport,
		},

		Schema: map[string]*schema.Schema{
			"connection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"new_private_virtual_interface": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"amazon_address": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"asn": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"auth_key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"customer_address": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"virtual_gateway_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"virtual_interface_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"vlan": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"amazon_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"asn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_account": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtual_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtual_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtual_interface_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtual_interface_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtual_interface_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceOutscaleDirectLinkInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).DL

	req := &dl.CreatePrivateVirtualInterfaceInput{
		ConnectionID: aws.String(d.Get("connection_id").(string)),
	}

	npv := d.Get("new_private_virtual_interface").(map[string]interface{})

	if v, ok := npv["virtual_gateway_id"]; ok && v.(string) != "" {
		req.NewPrivateVirtualInterface.VirtualGatewayID = aws.String(v.(string))
	} else {
		return fmt.Errorf("please provide the virtual_gateway_id attribute of new_private_virtual_interface it is required")
	}
	if v, ok := npv["vlan"]; ok && v.(string) != "" {
		i, _ := strconv.Atoi(v.(string))
		req.NewPrivateVirtualInterface.Vlan = aws.Int64(int64(i))
	} else {
		return fmt.Errorf("please provide the vlan attribute of new_private_virtual_interface it is required")
	}
	if v, ok := npv["virtual_interface_name"]; ok && v.(string) != "" {
		req.NewPrivateVirtualInterface.VirtualInterfaceName = aws.String(v.(string))
	} else {
		return fmt.Errorf("please provide the virtual_interface_name attribute of new_private_virtual_interface it is required")
	}
	if v, ok := npv["auth_key"]; ok && v.(string) != "" {
		req.NewPrivateVirtualInterface.AuthKey = aws.String(v.(string))
	}
	if v, ok := npv["asn"]; ok && v.(string) != "" {
		i, _ := strconv.Atoi(v.(string))
		req.NewPrivateVirtualInterface.Asn = aws.Int64(int64(i))
	} else {
		return fmt.Errorf("please provide the asn attribute of new_private_virtual_interface it is required")
	}
	if v, ok := npv["customer_address"]; ok && v.(string) != "" {
		req.NewPrivateVirtualInterface.CustomerAddress = aws.String(v.(string))
	}
	if v, ok := npv["amazon_address"]; ok && v.(string) != "" {
		req.NewPrivateVirtualInterface.AmazonAddress = aws.String(v.(string))
	}

	var resp *dl.VirtualInterface
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.CreatePrivateVirtualInterface(req)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating Direct Connect private virtual interface: %s", err.Error())
	}

	d.SetId(aws.StringValue(resp.VirtualInterfaceID))

	if err := dxPrivateVirtualInterfaceWaitUntilAvailable(d, conn); err != nil {
		return err
	}

	return resourceOutscaleDirectLinkInterfaceRead(d, meta)
}

func resourceOutscaleDirectLinkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).DL

	var err error
	var resp *dl.DescribeVirtualInterfacesOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.DescribeVirtualInterfaces(&dl.DescribeVirtualInterfacesInput{
			VirtualInterfaceID: aws.String(d.Id()),
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
		return fmt.Errorf("Error reading Direct Connect virtual interface: %s", err)
	}

	if resp == nil {
		log.Printf("[WARN] Direct Connect virtual interface (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	vif := resp.VirtualInterfaces[0]

	d.Set("asn", vif.Asn)
	d.Set("auth_key", vif.AuthKey)
	d.Set("address_family", vif.AddressFamily)
	d.Set("customer_address", vif.CustomerAddress)
	d.Set("location", vif.Location)
	d.Set("owner_account", vif.OwnerAccount)
	d.Set("amazon_address", vif.AmazonAddress)
	d.Set("virtual_gateway_id", vif.VirtualGatewayID)
	d.Set("virtual_interface_id", vif.VirtualInterfaceID)
	d.Set("virtual_interface_name", vif.VirtualInterfaceName)
	d.Set("virtual_interface_state", vif.VirtualInterfaceState)
	d.Set("virtual_interface_type", vif.VirtualInterfaceType)
	d.Set("vlan", vif.Vlan)
	d.Set("request_id", resp.RequestID)

	return nil
}

func resourceOutscaleDirectLinkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).DL

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.API.DeleteVirtualInterface(&dl.DeleteVirtualInterfaceInput{
			VirtualInterfaceID: aws.String(d.Id()),
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
		if strings.Contains(fmt.Sprint(err), "DirectConnectClientException") {
			return nil
		}
		return fmt.Errorf("Error deleting Direct Connect virtual interface: %s", err)
	}

	deleteStateConf := &resource.StateChangeConf{
		Pending: []string{
			"available",
			"confirming",
			"deleting",
			"down",
			"pending",
			"rejected",
			"verifying",
		},
		Target: []string{
			"deleted",
		},
		Refresh:    dxVirtualInterfaceStateRefresh(conn, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	_, err = deleteStateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Direct Connect virtual interface (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}

func resourceOutscaleDirectLinkInterfaceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func dxPrivateVirtualInterfaceWaitUntilAvailable(d *schema.ResourceData, conn *dl.Client) error {
	return dxVirtualInterfaceWaitUntilAvailable(
		d,
		conn,
		[]string{
			"pending",
		},
		[]string{
			"available",
			"down",
		})
}

func dxVirtualInterfaceWaitUntilAvailable(d *schema.ResourceData, conn *dl.Client, pending, target []string) error {
	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Refresh:    dxVirtualInterfaceStateRefresh(conn, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Direct Connect virtual interface (%s) to become available: %s", d.Id(), err)
	}

	return nil
}

func dxVirtualInterfaceStateRefresh(conn *dl.Client, vifID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		var err error
		var resp *dl.DescribeVirtualInterfacesOutput
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.API.DescribeVirtualInterfaces(&dl.DescribeVirtualInterfacesInput{
				VirtualInterfaceID: aws.String(vifID),
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
			return nil, "", err
		}

		n := len(resp.VirtualInterfaces)
		switch n {
		case 0:
			return "", "deleted", nil

		case 1:
			vif := resp.VirtualInterfaces[0]
			return vif, aws.StringValue(vif.VirtualInterfaceState), nil

		default:
			return nil, "", fmt.Errorf("Found %d Direct Connect virtual interfaces for %s, expected 1", n, vifID)
		}
	}
}
