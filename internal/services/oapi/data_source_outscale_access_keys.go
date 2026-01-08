package oapi

import (
	"context"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceOutscaleAccessKeys() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleAccessKeysRead,
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
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
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

func DataSourceOutscaleAccessKeysRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	accessKeyID, accessKeyOk := d.GetOk("access_key_ids")
	state, stateOk := d.GetOk("states")
	filterReq := &oscgo.FiltersAccessKeys{}

	var err error
	if filtersOk {
		filterReq, err = buildOutscaleDataSourceAccessKeyFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}
	if accessKeyOk {
		filterReq.SetAccessKeyIds(utils.InterfaceSliceToStringSlice(accessKeyID.([]interface{})))
	}
	if stateOk {
		filterReq.SetStates(utils.InterfaceSliceToStringSlice(state.([]interface{})))
	}
	req := oscgo.ReadAccessKeysRequest{
		Filters: filterReq,
	}

	if userName := d.Get("user_name").(string); userName != "" {
		req.SetUserName(userName)
	}
	var resp oscgo.ReadAccessKeysResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
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
		return ErrNoResults
	}

	if err := d.Set("access_keys", flattenAccessKeys(resp.GetAccessKeys())); err != nil {
		return err
	}

	d.SetId(id.UniqueId())
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
