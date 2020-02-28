package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceOutscaleOAPIPublicIP() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPIPublicIPRead,
		Schema: getOAPIPublicIPDataSourceSchema(),
	}
}

func getOAPIPublicIPDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"filter": dataSourceFiltersSchema(),
		"public_ip_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"public_ip": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"link_public_ip_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vm_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"nic_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"nic_account_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"private_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func dataSourceOutscaleOAPIPublicIPRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadPublicIpsRequest{
		Filters: &oscgo.FiltersPublicIp{},
	}

	if p, ok := d.GetOk("public_ip_id"); ok {
		req.Filters.SetPublicIpIds([]string{p.(string)})
	}

	if id, ok := d.GetOk("public_ip"); ok {
		req.Filters.SetPublicIps([]string{id.(string)})
	}

	filters, filtersOk := d.GetOk("filter")
	if filtersOk {
		req.Filters = buildOutscaleOAPIDataSourcePublicIpsFilters(filters.(*schema.Set))
	}

	var response oscgo.ReadPublicIpsResponse
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		var err error
		response, _, err = conn.PublicIpApi.ReadPublicIps(context.Background(), &oscgo.ReadPublicIpsOpts{ReadPublicIpsRequest: optional.NewInterface(req)})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving EIP: %s", err)
	}

	// Verify Outscale returned our EIP
	if len(response.GetPublicIps()) == 0 {
		return fmt.Errorf("Unable to find Public IP: %#v", req)
	}

	if len(response.GetPublicIps()) > 1 {
		return fmt.Errorf("multiple Public IPs matched; you can either use additional constraints to reduce matches to a single Public IP or use public_ips data source instead.")
	}

	address := response.GetPublicIps()[0]

	log.Printf("[DEBUG] EIP read configuration: %+v", address)

	if err := d.Set("link_public_ip_id", address.GetLinkPublicIpId()); err != nil {
		return err
	}
	if err := d.Set("vm_id", address.GetVmId()); err != nil {
		return err
	}

	if err := d.Set("nic_id", address.GetNicId()); err != nil {
		return err
	}

	if err := d.Set("nic_account_id", address.GetNicAccountId()); err != nil {
		return err
	}

	if err := d.Set("private_ip", address.GetPrivateIp()); err != nil {
		return err
	}

	if err := d.Set("request_id", response.ResponseContext.GetRequestId()); err != nil {
		return err
	}

	if err := d.Set("public_ip_id", address.GetPublicIpId()); err != nil {
		return err
	}

	if err := d.Set("public_ip", address.GetPublicIp()); err != nil {
		return err
	}

	d.SetId(address.GetPublicIp())

	return d.Set("request_id", response.ResponseContext.GetRequestId())
}

func buildOutscaleOAPIDataSourcePublicIpsFilters(set *schema.Set) *oscgo.FiltersPublicIp {
	var filters oscgo.FiltersPublicIp
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "public_ip_ids":
			filters.SetPublicIpIds(filterValues)
		case "link_public_ip_id":
			filters.SetLinkPublicIpIds(filterValues)
		case "placements":
			filters.SetPlacements(filterValues)
		case "vm_ids":
			filters.SetVmIds(filterValues)
		case "nic_ids":
			filters.SetNicIds(filterValues)
		case "nic_account_ids":
			filters.SetNicAccountIds(filterValues)
		case "private_ips":
			filters.SetPrivateIps(filterValues)
		case "public_ips":
			filters.SetPublicIps(filterValues)
		case "tags":
			filters.SetTags(filterValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return &filters
}
