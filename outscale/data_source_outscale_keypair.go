package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).OSCAPI
	req := oscgo.ReadKeypairsRequest{
		Filters: &oscgo.FiltersKeypair{KeypairNames: &[]string{d.Id()}},
	}

	KeyName, KeyNameisOk := d.GetOk("keypair_name")
	if KeyNameisOk {
		filter := oscgo.FiltersKeypair{}
		filter.SetKeypairNames([]string{KeyName.(string)})
		req.SetFilters(filter)
	}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.SetFilters(buildKeyPairsDataSourceFilters(filters.(*schema.Set)))
	}

	var resp oscgo.ReadKeypairsResponse
	var statusCode int
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		var err error
		rp, httpResp, err := conn.KeypairApi.ReadKeypairs(context.Background()).ReadKeypairsRequest(req).Execute()

		if err != nil {
			return utils.CheckThrottling(httpResp.StatusCode, err)
		}
		resp = rp
		statusCode = httpResp.StatusCode
		return nil
	})

	var errString string

	if err != nil {
		if statusCode == utils.ResourceNotFound {
			d.SetId("")
			return nil
		}
		errString = err.Error()

		return fmt.Errorf("Error retrieving KeyPair: %s", errString)
	}

	if len(resp.GetKeypairs()) < 1 {
		return fmt.Errorf("Unable to find key pair, please provide a better query criteria ")
	}
	if len(resp.GetKeypairs()) > 1 {

		return fmt.Errorf("Found to many key pairs, please provide a better query criteria ")
	}

	keypair := resp.GetKeypairs()[0]
	if err := d.Set("keypair_name", keypair.GetKeypairName()); err != nil {
		return err
	}
	if err := d.Set("keypair_fingerprint", keypair.GetKeypairFingerprint()); err != nil {
		return err
	}

	d.SetId(keypair.GetKeypairName())
	return nil
}

func dataSourceKeyPair() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeyPairRead,

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

func buildKeyPairsDataSourceFilters(set *schema.Set) oscgo.FiltersKeypair {
	var filters oscgo.FiltersKeypair
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "keypair_fingerprints":
			filters.SetKeypairFingerprints(filterValues)
		case "keypair_names":
			filters.SetKeypairNames(filterValues)
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
