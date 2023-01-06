package outscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleOAPICas() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPICasRead,
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

func dataSourceOutscaleOAPICasRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	params := oscgo.ReadCasRequest{}

	if filtersOk {
		params.Filters = buildOutscaleOAPIDataSourceCaFilters(filters.(*schema.Set))
	}
	var resp oscgo.ReadCasResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.CaApi.ReadCas(context.Background()).ReadCasRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error reading certificate authority id (%s)", utils.GetErrorResponse(err))
	}
	respCas := resp.GetCas()[:]
	if len(respCas) < 1 {
		d.SetId("")
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	blockCas := make([]map[string]interface{}, len(respCas))
	for k, v := range respCas {
		ca := make(map[string]interface{})
		if v.GetCaFingerprint() != "" {
			ca["ca_fingerprint"] = v.GetCaFingerprint()
		}
		if v.GetCaId() != "" {
			ca["ca_id"] = v.GetCaId()
		}
		if v.GetDescription() != "" {
			ca["description"] = v.GetDescription()
		}
		blockCas[k] = ca
	}
	d.SetId(resource.UniqueId())

	return d.Set("cas", blockCas)
}
