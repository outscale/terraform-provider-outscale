package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/samber/lo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleNatService() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleNatServiceRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"nat_service_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Attributes
			"public_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_ip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": TagsSchemaComputedSDK(),
		},
	}
}

func DataSourceOutscaleNatServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	natGatewayID, natGatewayIDOK := d.GetOk("nat_service_id")

	if !filtersOk && !natGatewayIDOK {
		return diag.Errorf("filters, or owner must be assigned, or nat_service_id must be provided")
	}

	var err error
	params := osc.ReadNatServicesRequest{}

	if filtersOk {
		params.Filters, err = buildOutscaleNatServiceDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if natGatewayIDOK && natGatewayID.(string) != "" {
		filter := osc.FiltersNatService{}
		filter.NatServiceIds = &[]string{natGatewayID.(string)}
		params.Filters = &filter
	}

	resp, err := client.ReadNatServices(ctx, params, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		errString := err.Error()

		return diag.Errorf("error reading nat service (%s)", errString)
	}

	if resp.NatServices == nil || len(*resp.NatServices) < 1 {
		return diag.FromErr(ErrNoResults)
	}

	if len(*resp.NatServices) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	return diag.FromErr(ngOAPIDescriptionAttributes(d, (*resp.NatServices)[0]))
}

// populate the numerous fields that the image description returns.
func ngOAPIDescriptionAttributes(d *schema.ResourceData, ng osc.NatService) error {
	d.SetId(ng.NatServiceId)

	if err := d.Set("nat_service_id", ng.NatServiceId); err != nil {
		return err
	}

	if ng.State != "" {
		if err := d.Set("state", ng.State); err != nil {
			return err
		}
	}
	if ng.SubnetId != "" {
		if err := d.Set("subnet_id", ng.SubnetId); err != nil {
			return err
		}
	}
	if ng.NetId != "" {
		if err := d.Set("net_id", ng.NetId); err != nil {
			return err
		}
	}

	addresses := make([]map[string]interface{}, len(ng.PublicIps))

	for k, v := range ng.PublicIps {
		address := make(map[string]interface{})
		if ptr.From(v.PublicIpId) != "" {
			address["public_ip_id"] = v.PublicIpId
		}
		if ptr.From(v.PublicIp) != "" {
			address["public_ip"] = v.PublicIp
		}
		addresses[k] = address
	}
	if err := d.Set("public_ips", addresses); err != nil {
		return err
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(ng.Tags)); err != nil {
		return err
	}

	return nil
}

func buildOutscaleNatServiceDataSourceFilters(set *schema.Set) (*osc.FiltersNatService, error) {
	var filters osc.FiltersNatService
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "nat_service_ids":
			filters.NatServiceIds = &filterValues
		case "net_ids":
			filters.NetIds = &filterValues
		case "states":
			filters.States = new(lo.Map(filterValues, func(s string, _ int) osc.NatServiceState { return osc.NatServiceState(s) }))
		case "subnet_ids":
			filters.SubnetIds = &filterValues
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
