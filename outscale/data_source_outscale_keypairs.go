package outscale

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleOAPiKeyPairsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.ReadKeypairsRequest{
		Filters: &oscgo.FiltersKeypair{},
	}

	//filters, filtersOk := d.GetOk("filter")
	KeyName, KeyNameisOk := d.GetOk("keypair_names")

	if KeyNameisOk {
		var names []string
		for _, v := range KeyName.([]interface{}) {
			names = append(names, v.(string))
		}
		filter := oscgo.FiltersKeypair{}
		filter.SetKeypairNames(names)
		req.SetFilters(filter)
	}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.SetFilters(buildOutscaleKeyPairsDataSourceFilters(filters.(*schema.Set)))
	}

	var resp oscgo.ReadKeypairsResponse
	var statusCode int
	err := retry.RetryContext(context.Background(), utils.ReadDefaultTimeout, func() *retry.RetryError {
		var err error
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

	d.SetId(id.UniqueId())

	keypairs := make([]map[string]interface{}, len(resp.GetKeypairs()))
	for k, v := range resp.GetKeypairs() {
		keypair := make(map[string]interface{})
		if v.HasKeypairName() {
			keypair["keypair_name"] = v.GetKeypairName()
		}
		if v.HasKeypairFingerprint() {
			keypair["keypair_fingerprint"] = v.GetKeypairFingerprint()
		}
		if v.HasKeypairId() {
			keypair["keypair_id"] = v.GetKeypairId()
		}
		if v.HasKeypairType() {
			keypair["keypair_type"] = v.GetKeypairType()
		}
		if v.HasTags() {
			keypair["tags"] = tagsOSCAPIToMap(v.GetTags())
		}
		keypairs[k] = keypair
	}

	return d.Set("keypairs", keypairs)
}

func DataSourceOutscaleKeyPairs() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleOAPiKeyPairsRead,

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
						"tags": dataSourceTagsSchema(),
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
