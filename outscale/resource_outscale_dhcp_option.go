package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

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
			"dhcp_options_name": {
				Type:     schema.TypeString,
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

	var resp oscgo.CreateDhcpOptionsResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.DhcpOptionApi.CreateDhcpOptions(context.Background(), &oscgo.CreateDhcpOptionsOpts{
			CreateDhcpOptionsRequest: optional.NewInterface(createOpts),
		})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.([]interface{}), *resp.GetDhcpOptionsSet().DhcpOptionsSetId, conn)
		if err != nil {
			return err
		}
	}

	d.SetId(*resp.GetDhcpOptionsSet().DhcpOptionsSetId)

	return resourceOutscaleDHCPOptionRead(d, meta)
}

func resourceOutscaleDHCPOptionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	dhcpID := d.Id()

	filterRequest := oscgo.ReadDhcpOptionsRequest{
		Filters: &oscgo.FiltersDhcpOptions{DhcpOptionsSetIds: &[]string{dhcpID}},
	}

	var resp oscgo.ReadDhcpOptionsResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.DhcpOptionApi.ReadDhcpOptions(context.Background(), &oscgo.ReadDhcpOptionsOpts{
			ReadDhcpOptionsRequest: optional.NewInterface(filterRequest),
		})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	dhcps, ok := resp.GetDhcpOptionsSetsOk()
	if !ok {
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
	if err := d.Set("dhcp_options_name", dhcp.GetDhcpOptionsName()); err != nil {
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

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err := conn.DhcpOptionApi.DeleteDhcpOptions(context.Background(), &oscgo.DeleteDhcpOptionsOpts{
			DeleteDhcpOptionsRequest: optional.NewInterface(oscgo.DeleteDhcpOptionsRequest{
				DhcpOptionsSetId: dhcpID,
			}),
		})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
