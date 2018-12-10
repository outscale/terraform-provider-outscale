package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func datasourceOutscaleOApiKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI
	req := &oapi.ReadKeypairsRequest{
		Filters: oapi.FiltersKeypair{KeypairNames: []string{d.Id()}},
	}

	KeyName, KeyNameisOk := d.GetOk("keypair_name")
	if KeyNameisOk {
		req.Filters.KeypairNames = []string{KeyName.(string)}
	}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters = buildOutscaleOAPIKeyPairsDataSourceFilters(filters.(*schema.Set))
	}

	var response *oapi.ReadKeypairsResponse
	var resp *oapi.POST_ReadKeypairsResponses
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.POST_ReadKeypairs(*req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var errString string

	if err != nil || resp.OK == nil {
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidOAPIKeyPair.NotFound") {
				d.SetId("")
				return nil
			}
			errString = err.Error()
		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
		}

		return fmt.Errorf("Error retrieving OAPIKeyPair: %s", errString)
	}

	response = resp.OK

	if len(response.Keypairs) < 1 {
		return fmt.Errorf("Unable to find key pair, please provide a better query criteria ")
	}
	if len(response.Keypairs) > 1 {

		return fmt.Errorf("Found to many key pairs, please provide a better query criteria ")
	}

	if response.ResponseContext.RequestId != "" {
		d.Set("request_id", response.ResponseContext.RequestId)
	}

	keypair := response.Keypairs[0]
	d.Set("keypair_name", keypair.KeypairName)
	d.Set("keypair_fingerprint", keypair.KeypairFingerprint)
	d.SetId(keypair.KeypairName)
	return nil
}

func datasourceOutscaleOAPIKeyPair() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOutscaleOApiKeyPairRead,

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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildOutscaleOAPIKeyPairsDataSourceFilters(set *schema.Set) oapi.FiltersKeypair {
	var filters oapi.FiltersKeypair
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "keypair_fingerprints":
			filters.KeypairFingerprints = filterValues
		case "keypair_names":
			filters.KeypairNames = filterValues
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
