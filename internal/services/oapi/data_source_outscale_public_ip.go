package oapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscalePublicIP() *schema.Resource {
	return &schema.Resource{
		Read:   DataSourceOutscalePublicIPRead,
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
			Computed: true,
		},
		"public_ip": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
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
			Computed: true,
		},
		"tags": TagsSchemaComputedSDK(),
	}
}

func DataSourceOutscalePublicIPRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	req := oscgo.ReadPublicIpsRequest{
		Filters: &oscgo.FiltersPublicIp{},
	}

	if p, ok := d.GetOk("public_ip_id"); ok {
		req.Filters.SetPublicIpIds([]string{p.(string)})
	}

	if id, ok := d.GetOk("public_ip"); ok {
		req.Filters.SetPublicIps([]string{id.(string)})
	}

	var err error
	filters, filtersOk := d.GetOk("filter")
	if filtersOk {
		req.Filters, err = buildOutscaleDataSourcePublicIpsFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var response oscgo.ReadPublicIpsResponse
	var statusCode int
	err = retry.Retry(60*time.Second, func() *retry.RetryError {
		var err error
		rp, httpResp, err := conn.PublicIpApi.ReadPublicIps(context.Background()).ReadPublicIpsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		response = rp
		statusCode = httpResp.StatusCode
		return nil
	})

	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving EIP: %s", err)
	}

	// Verify Outscale returned our EIP
	if err := utils.IsResponseEmptyOrMutiple(len(response.GetPublicIps()), "PublicIp"); err != nil {
		return err
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

	if err := d.Set("public_ip_id", address.GetPublicIpId()); err != nil {
		return err
	}

	if err := d.Set("tags", FlattenOAPITagsSDK(address.GetTags())); err != nil {
		return fmt.Errorf("Error setting PublicIp tags: %s", err)
	}

	d.Set("public_ip", address.PublicIp)

	d.SetId(address.GetPublicIp())

	return nil
}

func buildOutscaleDataSourcePublicIpsFilters(set *schema.Set) (*oscgo.FiltersPublicIp, error) {
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
		case "link_public_ip_ids":
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
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		case "tags":
			filters.SetTags(filterValues)
		default:
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}
