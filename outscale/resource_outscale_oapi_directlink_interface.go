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

func resourceOutscaleOAPIDirectLinkInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIDirectLinkInterfaceCreate,
		Read:   resourceOutscaleOAPIDirectLinkInterfaceRead,
		Delete: resourceOutscaleOAPIDirectLinkInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceOutscaleOAPIDirectLinkInterfaceImport,
		},

		Schema: map[string]*schema.Schema{
			"direct_link_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"direct_link_interface": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"outscale_private_ip": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"bgp_asn": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"bgp_key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"client_private_ip": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"vpn_gateway_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"direct_link_interface_name": {
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
			"outscale_private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bgp_asn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bgp_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_private_ip": {
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
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtual_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"direct_link_interface_name": {
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

func resourceOutscaleOAPIDirectLinkInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).DL

	req := &dl.CreatePrivateVirtualInterfaceInput{
		ConnectionID: aws.String(d.Get("direct_link_id").(string)),
	}

	npv := d.Get("direct_link_interface").(map[string]interface{})

	if v, ok := npv["vpn_gateway_id"]; ok && v.(string) != "" {
		req.NewPrivateVirtualInterface.VirtualGatewayID = aws.String(v.(string))
	} else {
		return fmt.Errorf("please provide the vpn_gateway_id attribute of direct_link_interface it is required")
	}
	if v, ok := npv["vlan"]; ok && v.(string) != "" {
		i, _ := strconv.Atoi(v.(string))
		req.NewPrivateVirtualInterface.Vlan = aws.Int64(int64(i))
	} else {
		return fmt.Errorf("please provide the vlan attribute of direct_link_interface it is required")
	}
	if v, ok := npv["direct_link_interface_name"]; ok && v.(string) != "" {
		req.NewPrivateVirtualInterface.VirtualInterfaceName = aws.String(v.(string))
	} else {
		return fmt.Errorf("please provide the direct_link_interface_name attribute of direct_link_interface it is required")
	}
	if v, ok := npv["bgp_key"]; ok && v.(string) != "" {
		req.NewPrivateVirtualInterface.AuthKey = aws.String(v.(string))
	}
	if v, ok := npv["bgp_asn"]; ok && v.(string) != "" {
		i, _ := strconv.Atoi(v.(string))
		req.NewPrivateVirtualInterface.Asn = aws.Int64(int64(i))
	} else {
		return fmt.Errorf("please provide the bgp_asn attribute of direct_link_interface it is required")
	}
	if v, ok := npv["client_private_ip"]; ok && v.(string) != "" {
		req.NewPrivateVirtualInterface.CustomerAddress = aws.String(v.(string))
	}
	if v, ok := npv["outscale_private_ip"]; ok && v.(string) != "" {
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

	return resourceOutscaleOAPIDirectLinkInterfaceRead(d, meta)
}

func resourceOutscaleOAPIDirectLinkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
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

	d.Set("name", vif.VirtualInterfaceName)
	d.Set("vlan", vif.Vlan)
	d.Set("bgp_asn", vif.Asn)
	d.Set("bgp_key", vif.AuthKey)
	d.Set("address_family", vif.AddressFamily)
	d.Set("client_private_ip", vif.CustomerAddress)
	d.Set("location", vif.Location)
	d.Set("owner_account", vif.OwnerAccount)
	d.Set("outscale_private_ip", vif.AmazonAddress)
	d.Set("vpn_gateway_id", vif.VirtualGatewayID)
	d.Set("virtual_interface_id", vif.VirtualInterfaceID)
	d.Set("direct_link_interface_name", vif.VirtualInterfaceName)
	d.Set("virtual_interface_state", vif.VirtualInterfaceState)
	d.Set("virtual_interface_type", vif.VirtualInterfaceType)
	d.Set("vlan", vif.Vlan)
	d.Set("request_id", resp.RequestID)

	return nil
}

func resourceOutscaleOAPIDirectLinkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceOutscaleOAPIDirectLinkInterfaceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
