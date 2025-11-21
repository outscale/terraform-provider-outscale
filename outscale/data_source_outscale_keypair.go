package outscale

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
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

	var err error
	if filtersOk {
		req.Filters, err = buildOutscaleKeyPairsDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp oscgo.ReadKeypairsResponse
	var statusCode int
	err = retry.RetryContext(context.Background(), utils.ReadDefaultTimeout, func() *retry.RetryError {
		rp, httpResp, err := conn.KeypairApi.ReadKeypairs(context.Background()).ReadKeypairsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		statusCode = httpResp.StatusCode
		return nil
	})

	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Keypair: %w", err)
	}

	if len(resp.GetKeypairs()) < 1 {
		return errors.New("Unable to find keypair, please provide a better query criteria")
	}
	if len(resp.GetKeypairs()) > 1 {
		return errors.New("Found to many keypairs, please provide a better query criteria")
	}

	keypair := resp.GetKeypairs()[0]
	if err := d.Set("keypair_name", keypair.GetKeypairName()); err != nil {
		return err
	}
	if err := d.Set("keypair_fingerprint", keypair.GetKeypairFingerprint()); err != nil {
		return err
	}
	if err := d.Set("keypair_type", keypair.GetKeypairType()); err != nil {
		return err
	}
	if err := d.Set("keypair_id", keypair.GetKeypairId()); err != nil {
		return err
	}
	if err := d.Set("tags", tagsOSCAPIToMap(keypair.GetTags())); err != nil {
		return err
	}
	d.SetId(keypair.GetKeypairId())
	return nil
}

func DataSourceOutscaleKeyPair() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleKeyPairRead,

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
			"tags": dataSourceTagsSchema(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildOutscaleKeyPairsDataSourceFilters(set *schema.Set) (*oscgo.FiltersKeypair, error) {
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
		case "keypair_ids":
			filters.SetKeypairIds(filterValues)
		case "keypair_types":
			filters.SetKeypairTypes(filterValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		case "tags":
			filters.SetTags(filterValues)
		default:
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}
