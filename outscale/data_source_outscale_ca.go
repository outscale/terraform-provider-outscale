package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleOAPICa() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPICaRead,
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

func dataSourceOutscaleOAPICaRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return fmt.Errorf("filters must be assigned")
	}

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
	if !resp.HasCas() {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}
	if len(resp.GetCas()) == 0 {
		d.SetId("")
		return fmt.Errorf("Certificate authority not found")
	}

	if len(resp.GetCas()) > 1 {
		return fmt.Errorf("your query returned more than one result, please try a more specific search criteria")
	}

	ca := resp.GetCas()[0]
	if err := d.Set("ca_fingerprint", ca.GetCaFingerprint()); err != nil {
		return err
	}
	if err := d.Set("ca_id", ca.GetCaId()); err != nil {
		return err
	}
	if err := d.Set("description", ca.GetDescription()); err != nil {
		return err
	}
	d.SetId(ca.GetCaId())
	return nil
}

func buildOutscaleOAPIDataSourceCaFilters(set *schema.Set) *oscgo.FiltersCa {
	var filters oscgo.FiltersCa
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "ca_fingerprints":
			filters.SetCaFingerprints(filterValues)
		case "ca_ids":
			filters.SetCaIds(filterValues)
		case "descriptions":
			filters.SetDescriptions(filterValues)
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return &filters
}
