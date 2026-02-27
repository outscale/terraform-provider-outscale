package oapi

import (
	"context"
	"log"
	"time"

	"github.com/outscale/goutils/sdk/ptr"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscalePublicIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscalePublicIPRead,
		Schema:      getOAPIPublicIPDataSourceSchema(),
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

func DataSourceOutscalePublicIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadPublicIpsRequest{
		Filters: &osc.FiltersPublicIp{},
	}

	if p, ok := d.GetOk("public_ip_id"); ok {
		req.Filters.PublicIpIds = &[]string{p.(string)}
	}

	if id, ok := d.GetOk("public_ip"); ok {
		req.Filters.PublicIps = &[]string{id.(string)}
	}

	var err error
	filters, filtersOk := d.GetOk("filter")
	if filtersOk {
		req.Filters, err = buildOutscaleDataSourcePublicIpsFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	response, err := client.ReadPublicIps(ctx, req, options.WithRetryTimeout(60*time.Second))
	if err != nil {
		return diag.Errorf("error retrieving eip: %s", err)
	}

	if response.PublicIps == nil || len(*response.PublicIps) == 0 {
		return diag.FromErr(ErrNoResults)
	}
	if len(*response.PublicIps) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	address := (*response.PublicIps)[0]

	log.Printf("[DEBUG] EIP read configuration: %+v", address)

	if err := d.Set("link_public_ip_id", ptr.From(address.LinkPublicIpId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vm_id", ptr.From(address.VmId)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nic_id", ptr.From(address.NicId)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nic_account_id", ptr.From(address.NicAccountId)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("private_ip", ptr.From(address.PrivateIp)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("public_ip_id", address.PublicIpId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tags", FlattenOAPITagsSDK(address.Tags)); err != nil {
		return diag.Errorf("error setting publicip tags: %s", err)
	}

	d.Set("public_ip", address.PublicIp)

	d.SetId(address.PublicIp)

	return nil
}

func buildOutscaleDataSourcePublicIpsFilters(set *schema.Set) (*osc.FiltersPublicIp, error) {
	var filters osc.FiltersPublicIp
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "public_ip_ids":
			filters.PublicIpIds = &filterValues
		case "link_public_ip_ids":
			filters.LinkPublicIpIds = &filterValues
		case "placements":
			filters.Placements = &filterValues
		case "vm_ids":
			filters.VmIds = &filterValues
		case "nic_ids":
			filters.NicIds = &filterValues
		case "nic_account_ids":
			filters.NicAccountIds = &filterValues
		case "private_ips":
			filters.PrivateIps = &filterValues
		case "public_ips":
			filters.PublicIps = &filterValues
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		case "tags":
			filters.Tags = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
