package oapi

import (
	"context"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscalePublicIPS() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscalePublicIPSRead,
		Schema:      oapiGetPublicIPSDataSourceSchema(),
	}
}

func oapiGetPublicIPSDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"filter": dataSourceFiltersSchema(),
		"public_ips": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"link_public_ip_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"public_ip_id": {
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
					"public_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tags": TagsSchemaComputedSDK(),
				},
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func DataSourceOutscalePublicIPSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadPublicIpsRequest{}

	filters, filtersOk := d.GetOk("filter")

	var err error
	if filtersOk {
		req.Filters, err = buildOutscaleDataSourcePublicIpsFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadPublicIps(ctx, req, options.WithRetryTimeout(60*time.Second))
	if err != nil {
		return diag.Errorf("error retrieving eip: %s", err)
	}

	// Verify Outscale returned our EIP
	if resp.PublicIps == nil || len(*resp.PublicIps) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	addresses := *resp.PublicIps

	address := make([]map[string]interface{}, len(addresses))

	for k, v := range addresses {
		add := make(map[string]interface{})

		add["link_public_ip_id"] = v.LinkPublicIpId
		add["public_ip_id"] = v.PublicIpId
		add["vm_id"] = v.VmId
		add["nic_id"] = v.NicId
		add["nic_account_id"] = v.NicAccountId
		add["private_ip"] = v.PrivateIp
		add["public_ip"] = v.PublicIp
		add["tags"] = FlattenOAPITagsSDK(v.Tags)
		address[k] = add
	}

	d.SetId(id.UniqueId())

	return diag.FromErr(d.Set("public_ips", address))
}
