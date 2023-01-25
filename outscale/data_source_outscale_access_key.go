package outscale

import (
	"context"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleAccessKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleAccessKeyRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filterReq := &oscgo.FiltersAccessKeys{}
	if filters, filtersOk := d.GetOk("filter"); filtersOk {
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

	if err = utils.IsResponseEmptyOrMutiple(len(resp.GetAccessKeys()), "Access Key"); err != nil {
		return err
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
