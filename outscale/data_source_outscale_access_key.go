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

func dataSourceOutscaleAccessKey() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleAccessKeyRead,
		Schema: getDataSourceSchemas(AccessKeySchema()),
	}
}

func dataSourceOutscaleAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")

	filterReq := &oscgo.FiltersAccessKeys{}
	if filtersOk {
		filterReq = buildOutscaleDataSourceAccessKeyFilters(filters.(*schema.Set))
	}

	var resp oscgo.ReadAccessKeysResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.AccessKeyApi.ReadAccessKeys(context.Background()).ReadAccessKeysRequest(oscgo.ReadAccessKeysRequest{Filters: filterReq}).Execute()
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
		return fmt.Errorf("Unable to find Access Key")
	}

	if len(resp.GetAccessKeys()) > 1 {
		return fmt.Errorf("multiple results returned, please use a more specific criteria in your query")
	}

	accessKey := resp.GetAccessKeys()[0]

	if err := d.Set("access_key_id", accessKey.GetAccessKeyId()); err != nil {
		return err
	}
	if err := d.Set("creation_date", accessKey.GetCreationDate()); err != nil {
		return err
	}
	if err := d.Set("expiration_date", accessKey.GetExpirationDate()); err != nil {
		return err
	}
	if err := d.Set("last_modification_date", accessKey.GetLastModificationDate()); err != nil {
		return err
	}
	if err := d.Set("state", accessKey.GetState()); err != nil {
		return err
	}

	d.SetId(accessKey.GetAccessKeyId())

	return nil
}

func buildOutscaleDataSourceAccessKeyFilters(set *schema.Set) *oscgo.FiltersAccessKeys {
	var filters oscgo.FiltersAccessKeys
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "access_key_ids":
			filters.SetAccessKeyIds(filterValues)
		case "states":
			filters.SetStates(filterValues)
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return &filters
}
