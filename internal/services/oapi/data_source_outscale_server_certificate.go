package oapi

import (
	"context"
	"log"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleServerCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleServerCertificateRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"orn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"upload_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceOutscaleServerCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.FromErr(ErrFilterRequired)
	}

	// Build up search parameters
	params := osc.ReadServerCertificatesRequest{}

	if filtersOk {
		filterParams, err := buildOutscaleOSCAPIDataSourceServerCertificateFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
		params.Filters = filterParams
	}

	resp, err := client.ReadServerCertificates(ctx, params, options.WithRetryTimeout(120*time.Second))
	if err != nil {
		return diag.Errorf("error reading server certificate id (%s)", err)
	}

	if resp.ServerCertificates == nil || len(*resp.ServerCertificates) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	if len(*resp.ServerCertificates) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	result := (*resp.ServerCertificates)[0]

	log.Printf("[DEBUG] Setting Server Certificate id (%s)", err)

	d.Set("expiration_date", from.ISO8601(result.ExpirationDate))
	d.Set("name", ptr.From(result.Name))
	d.Set("orn", ptr.From(result.Orn))
	d.Set("path", ptr.From(result.Path))
	d.Set("upload_date", from.ISO8601(result.UploadDate))

	d.SetId(ptr.From(result.Id))

	return nil
}

func buildOutscaleOSCAPIDataSourceServerCertificateFilters(set *schema.Set) (*osc.FiltersServerCertificate, error) {
	var filters osc.FiltersServerCertificate
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "paths":
			filters.Paths = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
