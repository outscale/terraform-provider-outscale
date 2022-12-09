package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleAccessKeys() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleAccessKeysRead,
		Schema: GetDataSourcesSchema("access_keys", AccessKeySchema()),
	}
}

func dataSourceOutscaleAccessKeysRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")

	filterReq := &oscgo.FiltersAccessKeys{}
	if filtersOk {
		filterReq = buildOutscaleDataSourceAccessKeyFilters(filters.(*schema.Set))
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
