package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceOutscaleAccessKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleAccessKeysRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"access_key_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"states": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
				},
			},
			"access_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiration_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modification_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
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

func dataSourceOutscaleAccessKeysRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	accessKeyID, accessKeyOk := d.GetOk("access_key_ids")
	state, stateOk := d.GetOk("states")

	if !filtersOk && !accessKeyOk && !stateOk {
		return fmt.Errorf("One of filters, access_key_ids or states must be assigned")
	}

	filterReq := &oscgo.FiltersAccessKeys{}
	if filtersOk {
		filterReq = buildOutscaleDataSourceAccessKeyFilters(filters.(*schema.Set))
	}
	if accessKeyOk {
		filterReq.SetAccessKeyIds(utils.InterfaceSliceToStringSlice(accessKeyID.([]interface{})))
	}
	if stateOk {
		filterReq.SetStates(utils.InterfaceSliceToStringSlice(state.([]interface{})))
	}

	var resp oscgo.ReadAccessKeysResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.AccessKeyApi.ReadAccessKeys(context.Background()).ReadAccessKeysRequest(oscgo.ReadAccessKeysRequest{
			Filters: filterReq,
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	if len(resp.GetAccessKeys()) == 0 {
		return fmt.Errorf("Unable to find Access Keys")
	}

	if err := d.Set("access_keys", flattenAccessKeys(resp.GetAccessKeys())); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenAccessKeys(accessKeys []oscgo.AccessKey) []map[string]interface{} {
	accessKeysMap := make([]map[string]interface{}, len(accessKeys))

	for i, ak := range accessKeys {
		accessKeysMap[i] = map[string]interface{}{
			"access_key_id":          ak.GetAccessKeyId(),
			"creation_date":          ak.GetCreationDate(),
			"expiration_date":        ak.GetExpirationDate(),
			"last_modification_date": ak.GetLastModificationDate(),
			"state":                  ak.GetState(),
		}
	}
	return accessKeysMap
}
