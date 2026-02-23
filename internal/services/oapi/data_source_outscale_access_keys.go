package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceOutscaleAccessKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleAccessKeysRead,
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

func DataSourceOutscaleAccessKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	accessKeyID, accessKeyOk := d.GetOk("access_key_ids")
	state, stateOk := d.GetOk("states")
	filterReq := &osc.FiltersAccessKeys{}

	var err error
	if filtersOk {
		filterReq, err = buildOutscaleDataSourceAccessKeyFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if accessKeyOk {
		filterReq.AccessKeyIds = new(utils.InterfaceSliceToStringSlice(accessKeyID.([]interface{})))
	}
	if stateOk {
		filterReq.States = new(utils.SliceToSuperStringSlice[osc.AccessKeyState](state.([]interface{})))
	}
	req := osc.ReadAccessKeysRequest{
		Filters: filterReq,
	}

	if userName := d.Get("user_name").(string); userName != "" {
		req.UserName = &userName
	}
	resp, err := client.ReadAccessKeys(ctx, req, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.AccessKeys == nil || len(*resp.AccessKeys) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	if err := d.Set("access_keys", flattenAccessKeys(*resp.AccessKeys)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.UniqueId())
	return nil
}

func flattenAccessKeys(accessKeys []osc.AccessKey) []map[string]interface{} {
	accessKeysMap := make([]map[string]interface{}, len(accessKeys))

	for i, ak := range accessKeys {
		accessKeysMap[i] = map[string]interface{}{
			"access_key_id":          ak.AccessKeyId,
			"creation_date":          from.ISO8601(ak.CreationDate),
			"expiration_date":        from.ISO8601(ak.ExpirationDate),
			"last_modification_date": from.ISO8601(ak.LastModificationDate),
			"state":                  ptr.From(ak.State),
		}
	}
	return accessKeysMap
}
