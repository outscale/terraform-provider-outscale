package oapi

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
)

func DataSourceOutscaleCas() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleCasRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"cas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ca_fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ca_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
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

func DataSourceOutscaleCasRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	params := osc.ReadCasRequest{}

	var err error
	if filtersOk {
		params.Filters, err = buildOutscaleDataSourceCaFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadCas(ctx, params, options.WithRetryTimeout(120*time.Second))
	if err != nil {
		return diag.Errorf("error reading certificate authority id (%s)", err)
	}
	respCas := ptr.From(resp.Cas)[:]
	if len(respCas) < 1 {
		d.SetId("")
		return diag.FromErr(ErrNoResults)
	}

	blockCas := make([]map[string]interface{}, len(respCas))
	for k, v := range respCas {
		ca := make(map[string]interface{})
		if ptr.From(v.CaFingerprint) != "" {
			ca["ca_fingerprint"] = v.CaFingerprint
		}
		if ptr.From(v.CaId) != "" {
			ca["ca_id"] = v.CaId
		}
		if ptr.From(v.Description) != "" {
			ca["description"] = v.Description
		}
		blockCas[k] = ca
	}
	d.SetId(id.UniqueId())

	return diag.FromErr(d.Set("cas", blockCas))
}
