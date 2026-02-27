package oapi

import (
	"context"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleKeyPairRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadKeypairsRequest{
		Filters: &osc.FiltersKeypair{KeypairNames: &[]string{d.Id()}},
	}

	KeyName, KeyNameisOk := d.GetOk("keypair_name")
	if KeyNameisOk {
		filter := osc.FiltersKeypair{}
		filter.KeypairNames = &[]string{KeyName.(string)}
		req.Filters = &filter
	}

	filters, filtersOk := d.GetOk("filter")

	var err error
	if filtersOk {
		req.Filters, err = buildOutscaleKeyPairsDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadKeypairs(ctx, req, options.WithRetryTimeout(ReadDefaultTimeout))
	if err != nil {
		return diag.Errorf("error retrieving keypair: %v", err)
	}

	if resp.Keypairs == nil || len(*resp.Keypairs) == 0 {
		return diag.FromErr(ErrNoResults)
	}
	if len(*resp.Keypairs) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	keypair := (*resp.Keypairs)[0]
	if err := d.Set("keypair_name", ptr.From(keypair.KeypairName)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("keypair_fingerprint", ptr.From(keypair.KeypairFingerprint)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("keypair_type", ptr.From(keypair.KeypairType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("keypair_id", ptr.From(keypair.KeypairId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(ptr.From(keypair.Tags))); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ptr.From(keypair.KeypairId))
	return nil
}

func DataSourceOutscaleKeyPair() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleKeyPairRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			// Attributes
			"keypair_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"keypair_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"keypair_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"keypair_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": TagsSchemaComputedSDK(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildOutscaleKeyPairsDataSourceFilters(set *schema.Set) (*osc.FiltersKeypair, error) {
	var filters osc.FiltersKeypair
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "keypair_fingerprints":
			filters.KeypairFingerprints = &filterValues
		case "keypair_names":
			filters.KeypairNames = &filterValues
		case "keypair_ids":
			filters.KeypairIds = &filterValues
		case "keypair_types":
			filters.KeypairTypes = &filterValues
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
