package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOutscaleDHCPOption() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleDHCPOptionCreate,
		Read:   resourceOutscaleDHCPOptionRead,
		Update: resourceOutscaleDHCPOptionUpdate,
		Delete: resourceOutscaleDHCPOptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"domain_name_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ntp_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dhcp_options_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsListOAPISchema(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleDHCPOptionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	createOpts := oscgo.CreateDhcpOptionsRequest{}

	domainName, okDomainName := d.GetOk("domain_name")
	domainNameServers, okDomainNameServers := d.GetOk("domain_name_servers")
	ntpServers, okNTPServers := d.GetOk("ntp_servers")

	if !okDomainName && !okDomainNameServers && !okNTPServers {
		return fmt.Errorf("Insufficient parameters provided out of: DomainName, domainNameServers, ntpServers. Expected at least: 1")
	}
	if okDomainName {
		createOpts.SetDomainName(domainName.(string))
	}
	if okDomainNameServers {
		createOpts.SetDomainNameServers(expandStringValueList(domainNameServers.([]interface{})))
	}
	if okNTPServers {
		createOpts.SetNtpServers(expandStringValueList(ntpServers.([]interface{})))
	}

	dhcp, _, err := createDhcpOption(conn, createOpts)
	if err != nil {
		return err
	}

	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.(*schema.Set), dhcp.GetDhcpOptionsSetId(), conn)
		if err != nil {
			return err
		}
	}

	d.SetId(dhcp.GetDhcpOptionsSetId())

	return resourceOutscaleDHCPOptionRead(d, meta)
}

func resourceOutscaleDHCPOptionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	dhcpID := d.Id()

	_, resp, err := readDhcpOption(conn, dhcpID)
	if err != nil {
		return err
	}

	dhcps := resp.GetDhcpOptionsSets()
	if len(dhcps) == 0 {
		d.SetId("")
		return nil
	}
	dhcp := dhcps[0]

	if err := d.Set("domain_name", dhcp.GetDomainName()); err != nil {
		return err
	}
	if err := d.Set("domain_name_servers", dhcp.GetDomainNameServers()); err != nil {
		return err
	}
	if err := d.Set("ntp_servers", dhcp.GetNtpServers()); err != nil {
		return err
	}
	if err := d.Set("default", dhcp.GetDefault()); err != nil {
		return err
	}
	if err := d.Set("dhcp_options_set_id", dhcp.GetDhcpOptionsSetId()); err != nil {
		return err
	}
	if err := d.Set("tags", tagsOSCAPIToMap(dhcp.GetTags())); err != nil {
		return err
	}
	if err := d.Set("request_id", resp.ResponseContext.GetRequestId()); err != nil {
		return err
	}
	return nil
}

func resourceOutscaleDHCPOptionUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	d.Partial(true)

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

	d.Partial(false)

	return resourceOutscaleDHCPOptionRead(d, meta)
}

func resourceOutscaleDHCPOptionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	dhcpID := d.Id()

	nets, err := getAttachedDHCPs(conn, dhcpID)
	if err != nil {
		return err
	}

	if err := detachDHCPs(conn, nets); err != nil {
		return err
	}

	// Deletes the dhcp option
	if err := deleteDhcpOptions(conn, dhcpID); err != nil {
		return err
	}

	return nil
}

func createDhcpOption(conn *oscgo.APIClient, dhcp oscgo.CreateDhcpOptionsRequest) (*oscgo.DhcpOptionsSet, *oscgo.CreateDhcpOptionsResponse, error) {
	var resp oscgo.CreateDhcpOptionsResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.DhcpOptionApi.CreateDhcpOptions(context.Background()).CreateDhcpOptionsRequest(dhcp).Execute()
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return resp.DhcpOptionsSet, &resp, err
}

func readDhcpOption(conn *oscgo.APIClient, dhcpID string) (*oscgo.DhcpOptionsSet, *oscgo.ReadDhcpOptionsResponse, error) {
	filterRequest := oscgo.ReadDhcpOptionsRequest{
		Filters: &oscgo.FiltersDhcpOptions{DhcpOptionsSetIds: &[]string{dhcpID}},
	}

	var resp oscgo.ReadDhcpOptionsResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.DhcpOptionApi.ReadDhcpOptions(context.Background()).ReadDhcpOptionsRequest(filterRequest).Execute()
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, &resp, err
	}

	dhcps := resp.GetDhcpOptionsSets()
	if len(dhcps) == 0 {
		return nil, &resp, fmt.Errorf("the Outscale DHCP Option is not found %s", dhcpID)
	}

	return &dhcps[0], &resp, err
}

func deleteDhcpOptions(conn *oscgo.APIClient, dhcpID string) error {
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err := conn.DhcpOptionApi.DeleteDhcpOptions(context.Background()).DeleteDhcpOptionsRequest(oscgo.DeleteDhcpOptionsRequest{
			DhcpOptionsSetId: dhcpID,
		}).Execute()
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return err
}

func getAttachedDHCPs(conn *oscgo.APIClient, dhcpID string) ([]oscgo.Net, error) {
	// Validate if the DHCP  Option is attached to a Net
	var resp oscgo.ReadNetsResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, _, err = conn.NetApi.ReadNets(context.Background()).ReadNetsRequest(oscgo.ReadNetsRequest{
			Filters: &oscgo.FiltersNet{
				DhcpOptionsSetIds: &[]string{dhcpID},
			},
		}).Execute()

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("[DEBUG] Error reading network (%s)", err)
	}

	return resp.GetNets(), nil
}

func detachDHCPs(conn *oscgo.APIClient, nets []oscgo.Net) error {
	// Detaching the dhcp of the nets
	for _, net := range nets {
		_, _, err := conn.NetApi.UpdateNet(context.Background()).UpdateNetRequest(oscgo.UpdateNetRequest{
			DhcpOptionsSetId: "default",
			NetId:            net.GetNetId(),
		}).Execute()
		if err != nil {
			return fmt.Errorf("Error updating net(%s) in DHCP Option resource: %s", net.GetNetId(), err)
		}
	}
	return nil
}
