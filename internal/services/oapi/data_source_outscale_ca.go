package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscaleCa() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleCaRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"ca_pem": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceOutscaleCaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return diag.FromErr(ErrFilterRequired)
	}

	params := osc.ReadCasRequest{}
	if filtersOk {
		filterParams, err := buildOutscaleDataSourceCaFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
		params.Filters = filterParams
	}

	resp, err := client.ReadCas(ctx, params, options.WithRetryTimeout(120*time.Second))
	if err != nil {
		return diag.Errorf("error reading certificate authority id (%s)", err)
	}

	if resp.Cas == nil || len(*resp.Cas) == 0 {
		d.SetId("")
		return diag.FromErr(ErrNoResults)
	}

	if len(*resp.Cas) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	ca := (*resp.Cas)[0]
	if err := d.Set("ca_fingerprint", ptr.From(ca.CaFingerprint)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ca_id", ptr.From(ca.CaId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", ptr.From(ca.Description)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ptr.From(ca.CaId))
	return nil
}

func buildOutscaleDataSourceCaFilters(set *schema.Set) (*osc.FiltersCa, error) {
	var filters osc.FiltersCa
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "ca_fingerprints":
			filters.CaFingerprints = &filterValues
		case "ca_ids":
			filters.CaIds = &filterValues
		case "descriptions":
			filters.Descriptions = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
