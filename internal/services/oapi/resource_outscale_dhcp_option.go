package oapi

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleDHCPOption() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleDHCPOptionCreate,
		Read:   ResourceOutscaleDHCPOptionRead,
		Update: ResourceOutscaleDHCPOptionUpdate,
		Delete: ResourceOutscaleDHCPOptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"log_servers": {
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
			"tags": TagsSchemaSDK(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceOutscaleDHCPOptionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	createOpts := oscgo.CreateDhcpOptionsRequest{}

	domainName, okDomainName := d.GetOk("domain_name")
	domainNameServers, okDomainNameServers := d.GetOk("domain_name_servers")
	logServers, okLogServers := d.GetOk("log_servers")
	ntpServers, okNTPServers := d.GetOk("ntp_servers")

	if !okDomainName && !okDomainNameServers && !okLogServers && !okNTPServers {
		return fmt.Errorf("insufficient parameters provided out of: domainname, domainnameservers, logservers, ntpservers - expected at least: 1")
	}
	if okDomainName {
		createOpts.SetDomainName(domainName.(string))
	}
	if okDomainNameServers {
		createOpts.SetDomainNameServers(utils.InterfaceSliceToStringSlice(domainNameServers.([]interface{})))
	}
	if okLogServers {
		createOpts.SetLogServers(utils.InterfaceSliceToStringSlice(logServers.([]interface{})))
	}
	if okNTPServers {
		createOpts.SetNtpServers(utils.InterfaceSliceToStringSlice(ntpServers.([]interface{})))
	}

	dhcp, _, err := createDhcpOption(conn, createOpts)
	if err != nil {
		return err
	}
	d.SetId(dhcp.GetDhcpOptionsSetId())

	err = createOAPITagsSDK(conn, d)
	if err != nil {
		return err
	}

	return ResourceOutscaleDHCPOptionRead(d, meta)
}

func ResourceOutscaleDHCPOptionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	dhcpID := d.Id()

	_, resp, err := readDhcpOption(conn, dhcpID)
	if err != nil {
		return err
	}

	dhcps := resp.GetDhcpOptionsSets()
	if utils.IsResponseEmpty(len(dhcps), "DhcpOption", d.Id()) {
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
	if err := d.Set("log_servers", dhcp.GetLogServers()); err != nil {
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
	if err := d.Set("tags", FlattenOAPITagsSDK(dhcp.GetTags())); err != nil {
		return err
	}

	return nil
}

func ResourceOutscaleDHCPOptionUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	if err := updateOAPITagsSDK(conn, d); err != nil {
		return err
	}
	return ResourceOutscaleDHCPOptionRead(d, meta)
}

func ResourceOutscaleDHCPOptionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

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
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.DhcpOptionApi.CreateDhcpOptions(context.Background()).CreateDhcpOptionsRequest(dhcp).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
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
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.DhcpOptionApi.ReadDhcpOptions(context.Background()).ReadDhcpOptionsRequest(filterRequest).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return nil, &resp, err
	}

	dhcps := resp.GetDhcpOptionsSets()
	if len(dhcps) == 0 {
		return nil, &resp, fmt.Errorf("the outscale dhcp option is not found %s", dhcpID)
	}

	return &dhcps[0], &resp, err
}

func deleteDhcpOptions(conn *oscgo.APIClient, dhcpID string) error {
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, httpResp, err := conn.DhcpOptionApi.DeleteDhcpOptions(context.Background()).DeleteDhcpOptionsRequest(oscgo.DeleteDhcpOptionsRequest{
			DhcpOptionsSetId: dhcpID,
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	return err
}

func getAttachedDHCPs(conn *oscgo.APIClient, dhcpID string) ([]oscgo.Net, error) {
	// Validate if the DHCP  Option is attached to a Net
	var resp oscgo.ReadNetsResponse
	err := retry.Retry(120*time.Second, func() *retry.RetryError {
		rp, httpResp, err := conn.NetApi.ReadNets(context.Background()).ReadNetsRequest(oscgo.ReadNetsRequest{
			Filters: &oscgo.FiltersNet{
				DhcpOptionsSetIds: &[]string{dhcpID},
			},
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error reading network (%s)", err)
	}

	return resp.GetNets(), nil
}

func detachDHCPs(conn *oscgo.APIClient, nets []oscgo.Net) error {
	// Detaching the dhcp of the nets
	for _, net := range nets {
		err := retry.Retry(120*time.Second, func() *retry.RetryError {
			_, httpResp, err := conn.NetApi.UpdateNet(context.Background()).UpdateNetRequest(oscgo.UpdateNetRequest{
				DhcpOptionsSetId: "default",
				NetId:            net.GetNetId(),
			}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error updating net(%s) in dhcp option resource: %s", net.GetNetId(), err)
		}
	}
	return nil
}
