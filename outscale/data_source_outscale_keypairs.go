package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/outscale/osc-go/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func datasourceOutscaleOAPiKeyPairsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI
	req := &oapi.ReadKeypairsRequest{
		Filters: oapi.FiltersKeypair{},
	}

	//filters, filtersOk := d.GetOk("filter")
	KeyName, KeyNameisOk := d.GetOk("keypair_names")

	if KeyNameisOk {
		var names []string
		for _, v := range KeyName.([]interface{}) {
			names = append(names, v.(string))
		}
		req.Filters.KeypairNames = names
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

	d.SetId(resource.UniqueId())

	if response.ResponseContext.RequestId != "" {
		d.Set("request_id", response.ResponseContext.RequestId)
	}

	keypairs := make([]map[string]interface{}, len(response.Keypairs))
	for k, v := range response.Keypairs {
		keypair := make(map[string]interface{})
		if v.KeypairName != "" {
			keypair["keypair_name"] = v.KeypairName
		}
		if v.KeypairFingerprint != "" {
			keypair["keypair_fingerprint"] = v.KeypairFingerprint
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
