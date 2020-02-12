package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOutscaleOAPIVpnConnectionRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIVpnConnectionRouteCreate,
		Read:   resourceOutscaleOAPIVpnConnectionRouteRead,
		Delete: resourceOutscaleOAPIVpnConnectionRouteDelete,

		Schema: map[string]*schema.Schema{
			"destination_ip_range": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpn_connection_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPIVpnConnectionRouteCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	createOpts := &fcu.CreateVpnConnectionRouteInput{
		DestinationCidrBlock: aws.String(d.Get("destination_ip_range").(string)),
		VpnConnectionId:      aws.String(d.Get("vpn_connection_id").(string)),
	}

	// Create the route.
	log.Printf("[DEBUG] Creating VPN connection route")

	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.CreateVpnConnectionRoute(createOpts)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("Error creating VPN connection route: %s", err)
	}

	// Store the ID by the only two data we have available to us.
	d.SetId(fmt.Sprintf("%s:%s", *createOpts.DestinationCidrBlock, *createOpts.VpnConnectionId))

	return resourceOutscaleOAPIVpnConnectionRouteRead(d, meta)
}

func resourceOutscaleOAPIVpnConnectionRouteRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	cidrBlock, vpnConnectionID := resourceOutscaleOAPIVpnConnectionRouteParseID(d.Id())

	routeFilters := []*fcu.Filter{
		&fcu.Filter{
			Name:   aws.String("route.destination-cidr-block"),
			Values: []*string{aws.String(cidrBlock)},
		},
		&fcu.Filter{
			Name:   aws.String("vpn-connection-ID"),
			Values: []*string{aws.String(vpnConnectionID)},
		},
	}

	var err error
	var resp *fcu.DescribeVpnConnectionsOutput

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeVpnConnections(&fcu.DescribeVpnConnectionsInput{
			Filters: routeFilters,
		})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidVpnConnectionID.NotFound") {
			d.SetId("")
			return nil
		}
		log.Printf("[ERROR] Error finding VPN connection route: %s", err)
		return err
	}
	if resp == nil || len(resp.VpnConnections) == 0 {
		return fmt.Errorf("No VPN connections returned")
	}

	vpnConnection := resp.VpnConnections[0]

	var found bool
	for _, r := range vpnConnection.Routes {
		if *r.DestinationCidrBlock == cidrBlock {
			d.Set("destination_ip_range", *r.DestinationCidrBlock)
			d.Set("vpn_connection_id", *vpnConnection.VpnConnectionId)
			found = true
		}
	}
	if !found {
		// Something other than tersraform eliminated the route.
		d.SetId("")
	}

	d.Set("request_id", resp.RequestId)

	return nil
}

func resourceOutscaleOAPIVpnConnectionRouteDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.DeleteVpnConnectionRoute(&fcu.DeleteVpnConnectionRouteInput{
			DestinationCidrBlock: aws.String(d.Get("destination_ip_range").(string)),
			VpnConnectionId:      aws.String(d.Get("vpn_connection_id").(string)),
		})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidVpnConnectionID.NotFound") {
			d.SetId("")
			return nil
		}
		log.Printf("[ERROR] Error deleting VPN connection route: %s", err)
		return err
	}

	return nil
}

func resourceOutscaleOAPIVpnConnectionRouteParseID(ID string) (string, string) {
	parts := strings.SplitN(ID, ":", 2)
	return parts[0], parts[1]
}
