package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleNatServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleNatServicesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"nat_service_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// Attributes
			"nat_services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"nat_service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_id": {
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
		},
	}
}

func DataSourceOutscaleNatServicesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	natGatewayID, natGatewayIDOK := d.GetOk("nat_service_ids")

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
	if natGatewayIDOK {
		ids := make([]string, len(natGatewayID.([]interface{})))

		for k, v := range natGatewayID.([]interface{}) {
			ids[k] = v.(string)
		}
		filter := osc.FiltersNatService{}
		filter.NatServiceIds = &ids
		params.Filters = &filter
	}

	resp, err := client.ReadNatServices(ctx, params, options.WithRetryTimeout(5*time.Minute))

	var errString string

	if err != nil {
		errString = err.Error()

		return diag.Errorf("error reading nat service (%s)", errString)
	}

	if resp.NatServices == nil || len(*resp.NatServices) < 1 {
		return diag.FromErr(ErrNoResults)
	}

	return diag.FromErr(ngsOAPIDescriptionAttributes(d, *resp.NatServices))
}

// populate the numerous fields that the image description returns.
func ngsOAPIDescriptionAttributes(d *schema.ResourceData, ngs []osc.NatService) error {
	d.SetId(id.UniqueId())

	addngs := make([]map[string]interface{}, len(ngs))

	for k, v := range ngs {
		addng := make(map[string]interface{})

		ngas := make([]interface{}, len(v.PublicIps))

		for i, w := range v.PublicIps {
			nga := make(map[string]interface{})
			if ptr.From(w.PublicIpId) != "" {
				nga["public_ip_id"] = w.PublicIpId
			}
			if ptr.From(w.PublicIp) != "" {
				nga["public_ip"] = w.PublicIp
			}
			ngas[i] = nga
		}
		addng["public_ips"] = ngas

		if v.NatServiceId != "" {
			addng["nat_service_id"] = v.NatServiceId
		}
		if v.State != "" {
			addng["state"] = v.State
		}
		if v.SubnetId != "" {
			addng["subnet_id"] = v.SubnetId
		}
		if v.NetId != "" {
			addng["net_id"] = v.NetId
		}
		if v.Tags != nil {
			addng["tags"] = FlattenOAPITagsSDK(v.Tags)
		}

		addngs[k] = addng
	}

	return d.Set("nat_services", addngs)
}
