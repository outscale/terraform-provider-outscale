package oapi

import (
	"context"
	"fmt"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleDHCPOption() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleDHCPOptionCreate,
		ReadContext:   ResourceOutscaleDHCPOptionRead,
		UpdateContext: ResourceOutscaleDHCPOptionUpdate,
		DeleteContext: ResourceOutscaleDHCPOptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Update: schema.DefaultTimeout(UpdateDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
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

func ResourceOutscaleDHCPOptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutCreate)

	createReq := osc.CreateDhcpOptionsRequest{}

	domainName, okDomainName := d.GetOk("domain_name")
	domainNameServers, okDomainNameServers := d.GetOk("domain_name_servers")
	logServers, okLogServers := d.GetOk("log_servers")
	ntpServers, okNTPServers := d.GetOk("ntp_servers")

	if !okDomainName && !okDomainNameServers && !okLogServers && !okNTPServers {
		return diag.Errorf("insufficient parameters provided out of: domainname, domainnameservers, logservers, ntpservers - expected at least: 1")
	}
	if okDomainName {
		createReq.DomainName = new(domainName.(string))
	}
	if okDomainNameServers {
		createReq.DomainNameServers = new(utils.InterfaceSliceToStringSlice(domainNameServers.([]interface{})))
	}
	if okLogServers {
		createReq.LogServers = new(utils.InterfaceSliceToStringSlice(logServers.([]interface{})))
	}
	if okNTPServers {
		createReq.NtpServers = new(utils.InterfaceSliceToStringSlice(ntpServers.([]interface{})))
	}

	dhcp, _, err := createDhcpOption(ctx, client, timeout, createReq)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ptr.From(dhcp.DhcpOptionsSetId))

	err = createOAPITagsSDK(ctx, client, timeout, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceOutscaleDHCPOptionRead(ctx, d, meta)
}

func ResourceOutscaleDHCPOptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutRead)

	dhcpID := d.Id()

	_, resp, err := readDhcpOption(ctx, client, timeout, dhcpID)
	if err != nil {
		return diag.FromErr(err)
	}

	dhcps := ptr.From(resp.DhcpOptionsSets)
	if utils.IsResponseEmpty(len(dhcps), "DhcpOption", d.Id()) {
		d.SetId("")
		return nil
	}
	dhcp := dhcps[0]

	if err := d.Set("domain_name", ptr.From(dhcp.DomainName)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("domain_name_servers", ptr.From(dhcp.DomainNameServers)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("log_servers", ptr.From(dhcp.LogServers)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ntp_servers", ptr.From(dhcp.NtpServers)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default", ptr.From(dhcp.Default)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dhcp_options_set_id", ptr.From(dhcp.DhcpOptionsSetId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(ptr.From(dhcp.Tags))); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceOutscaleDHCPOptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutUpdate)

	if err := updateOAPITagsSDK(ctx, client, timeout, d); err != nil {
		return diag.FromErr(err)
	}

	return ResourceOutscaleDHCPOptionRead(ctx, d, meta)
}

func ResourceOutscaleDHCPOptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutDelete)

	dhcpID := d.Id()

	nets, err := getAttachedDHCPs(ctx, client, timeout, dhcpID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := detachDHCPs(ctx, client, timeout, nets); err != nil {
		return diag.FromErr(err)
	}

	// Deletes the dhcp option
	if err := deleteDhcpOptions(ctx, client, timeout, dhcpID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func createDhcpOption(ctx context.Context, client *osc.Client, timeout time.Duration, dhcp osc.CreateDhcpOptionsRequest) (*osc.DhcpOptionsSet, *osc.CreateDhcpOptionsResponse, error) {
	resp, err := client.CreateDhcpOptions(ctx, dhcp, options.WithRetryTimeout(timeout))
	if err != nil {
		return nil, nil, err
	}

	return resp.DhcpOptionsSet, resp, err
}

func readDhcpOption(ctx context.Context, client *osc.Client, timeout time.Duration, dhcpID string) (*osc.DhcpOptionsSet, *osc.ReadDhcpOptionsResponse, error) {
	filterRequest := osc.ReadDhcpOptionsRequest{
		Filters: &osc.FiltersDhcpOptions{DhcpOptionsSetIds: &[]string{dhcpID}},
	}

	resp, err := client.ReadDhcpOptions(ctx, filterRequest, options.WithRetryTimeout(timeout))
	if err != nil {
		return nil, resp, err
	}

	dhcps := ptr.From(resp.DhcpOptionsSets)
	if len(dhcps) == 0 {
		return nil, resp, fmt.Errorf("the outscale dhcp option is not found %s", dhcpID)
	}

	return &dhcps[0], resp, err
}

func deleteDhcpOptions(ctx context.Context, client *osc.Client, timeout time.Duration, dhcpID string) error {
	_, err := client.DeleteDhcpOptions(ctx, osc.DeleteDhcpOptionsRequest{
		DhcpOptionsSetId: dhcpID,
	}, options.WithRetryTimeout(timeout))
	return err
}

func getAttachedDHCPs(ctx context.Context, client *osc.Client, timeout time.Duration, dhcpID string) ([]osc.Net, error) {
	// Validate if the DHCP  Option is attached to a Net
	resp, err := client.ReadNets(ctx, osc.ReadNetsRequest{
		Filters: &osc.FiltersNet{
			DhcpOptionsSetIds: &[]string{dhcpID},
		},
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return nil, fmt.Errorf("error reading network (%s)", err)
	}

	return ptr.From(resp.Nets), nil
}

func detachDHCPs(ctx context.Context, client *osc.Client, timeout time.Duration, nets []osc.Net) error {
	// Detaching the dhcp of the nets
	for _, net := range nets {
		_, err := client.UpdateNet(ctx, osc.UpdateNetRequest{
			DhcpOptionsSetId: "default",
			NetId:            net.NetId,
		}, options.WithRetryTimeout(timeout))
		if err != nil {
			return fmt.Errorf("error updating net(%s) in dhcp option resource: %s", net.NetId, err)
		}
	}
	return nil
}
