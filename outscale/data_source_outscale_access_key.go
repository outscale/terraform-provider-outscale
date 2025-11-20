package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceOutscaleAccessKey() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleAccessKeyRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_key_id": {
				Type:     schema.TypeString,
				Optional: true,
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceOutscaleAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	accessKeyID, accessKeyOk := d.GetOk("access_key_id")
	state, stateOk := d.GetOk("state")
	userName, userNameOk := d.GetOk("user_name")

	if !filtersOk && !accessKeyOk && !stateOk && !userNameOk {
		return fmt.Errorf("one of filters: access_key_id, state or user_name must be assigned")
	}

	filterReq := &oscgo.FiltersAccessKeys{}

	var err error
	if filtersOk {
		filterReq, err = buildOutscaleDataSourceAccessKeyFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}
	if accessKeyOk {
		filterReq.SetAccessKeyIds([]string{accessKeyID.(string)})
	}
	if stateOk {
		filterReq.SetStates([]string{state.(string)})
	}
	req := oscgo.ReadAccessKeysRequest{}
	req.SetFilters(*filterReq)
	if userNameOk {
		req.SetUserName(userName.(string))
	}
	var resp oscgo.ReadAccessKeysResponse

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.AccessKeyApi.ReadAccessKeys(context.Background()).ReadAccessKeysRequest(req).Execute()
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
		return fmt.Errorf("unable to find Access Key")
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

func buildOutscaleDataSourceAccessKeyFilters(set *schema.Set) (*oscgo.FiltersAccessKeys, error) {
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
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}
