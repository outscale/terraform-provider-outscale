package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceOutscaleOAPiKeyPairsRead(d *schema.ResourceData, meta interface{}) error {
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
		req.SetFilters(buildOutscaleOAPIKeyPairsDataSourceFilters(filters.(*schema.Set)))
	}

	var resp oscgo.ReadKeypairsResponse
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		var err error
		resp, _, err = conn.KeypairApi.ReadKeypairs(context.Background(), &oscgo.ReadKeypairsOpts{ReadKeypairsRequest: optional.NewInterface(req)})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var errString string

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidOAPIKeyPair.NotFound") {
			d.SetId("")
			return nil
		}
		errString = err.Error()

		return fmt.Errorf("Error retrieving OAPIKeyPair: %s", errString)
	}

	if len(resp.GetKeypairs()) < 1 {
		return fmt.Errorf("Unable to find key pair, please provide a better query criteria ")
	}

	d.SetId(resource.UniqueId())

	if resp.ResponseContext.GetRequestId() != "" {
		d.Set("request_id", resp.ResponseContext.GetRequestId())
	}

	keypairs := make([]map[string]interface{}, len(resp.GetKeypairs()))
	for k, v := range resp.GetKeypairs() {
		keypair := make(map[string]interface{})
		if v.GetKeypairName() != "" {
			keypair["keypair_name"] = v.GetKeypairName()
		}
		if v.GetKeypairFingerprint() != "" {
			keypair["keypair_fingerprint"] = v.GetKeypairFingerprint()
		}
		keypairs[k] = keypair
	}
	d.Set("keypairs", keypairs)
	return nil
}

func datasourceOutscaleOAPIKeyPairs() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOutscaleOAPiKeyPairsRead,

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
