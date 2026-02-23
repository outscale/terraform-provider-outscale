package oapi

import (
	"context"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleOAPiKeyPairsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadKeypairsRequest{
		Filters: &osc.FiltersKeypair{},
	}

	// filters, filtersOk := d.GetOk("filter")
	KeyName, KeyNameisOk := d.GetOk("keypair_names")

	if KeyNameisOk {
		var names []string
		for _, v := range KeyName.([]interface{}) {
			names = append(names, v.(string))
		}
		filter := osc.FiltersKeypair{}
		filter.KeypairNames = &names
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

	if resp.Keypairs == nil || len(*resp.Keypairs) < 1 {
		return diag.FromErr(ErrNoResults)
	}

	d.SetId(id.UniqueId())

	keypairs := make([]map[string]interface{}, len(*resp.Keypairs))
	for k, v := range *resp.Keypairs {
		keypair := make(map[string]interface{})
		if v.KeypairName != nil {
			keypair["keypair_name"] = v.KeypairName
		}
		if v.KeypairFingerprint != nil {
			keypair["keypair_fingerprint"] = v.KeypairFingerprint
		}
		if v.KeypairId != nil {
			keypair["keypair_id"] = v.KeypairId
		}
		if v.KeypairType != nil {
			keypair["keypair_type"] = v.KeypairType
		}
		if v.Tags != nil {
			keypair["tags"] = FlattenOAPITagsSDK(ptr.From(v.Tags))
		}
		keypairs[k] = keypair
	}

	return diag.FromErr(d.Set("keypairs", keypairs))
}

func DataSourceOutscaleKeyPairs() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleOAPiKeyPairsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			// Attributes
			"keypair_names": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"keypairs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keypair_fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"keypair_name": {
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
