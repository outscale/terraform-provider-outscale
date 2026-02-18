package oapi

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscaleCas() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleCasRead,
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

func DataSourceOutscaleCasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	params := osc.ReadCasRequest{}

	var err error
	if filtersOk {
		params.Filters, err = buildOutscaleDataSourceCaFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp osc.ReadCasResponse
	err = retry.Retry(120*time.Second, func() *retry.RetryError {
		rp, httpResp, err := client.CaApi.ReadCas(ctx).ReadCasRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading certificate authority id (%s)", utils.GetErrorResponse(err))
	}
	respCas := resp.GetCas()[:]
	if len(respCas) < 1 {
		d.SetId("")
		return ErrNoResults
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
	d.SetId(id.UniqueId())

	return d.Set("cas", blockCas)
}
