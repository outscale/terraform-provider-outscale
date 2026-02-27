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
	"github.com/samber/lo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceOutscaleAccessKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleAccessKeyRead,
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

func DataSourceOutscaleAccessKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	accessKeyID, accessKeyOk := d.GetOk("access_key_id")
	state, stateOk := d.GetOk("state")
	userName, userNameOk := d.GetOk("user_name")

	if !filtersOk && !accessKeyOk && !stateOk && !userNameOk {
		return diag.Errorf("one of filters: access_key_id, state or user_name must be assigned")
	}

	filterReq := &osc.FiltersAccessKeys{}

	var err error
	if filtersOk {
		filterReq, err = buildOutscaleDataSourceAccessKeyFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if accessKeyOk {
		filterReq.AccessKeyIds = &[]string{accessKeyID.(string)}
	}
	if stateOk {
		filterReq.States = &[]osc.AccessKeyState{osc.AccessKeyState(state.(string))}
	}
	req := osc.ReadAccessKeysRequest{}
	req.Filters = filterReq
	if userNameOk {
		req.UserName = new(userName.(string))
	}

	resp, err := client.ReadAccessKeys(ctx, req, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.AccessKeys == nil || len(*resp.AccessKeys) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	if len(*resp.AccessKeys) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	accessKey := (*resp.AccessKeys)[0]

	if err := d.Set("access_key_id", ptr.From(accessKey.AccessKeyId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("creation_date", from.ISO8601(accessKey.CreationDate)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("expiration_date", from.ISO8601(accessKey.ExpirationDate)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_modification_date", from.ISO8601(accessKey.LastModificationDate)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", ptr.From(accessKey.State)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ptr.From(accessKey.AccessKeyId))

	return nil
}

func buildOutscaleDataSourceAccessKeyFilters(set *schema.Set) (*osc.FiltersAccessKeys, error) {
	var filters osc.FiltersAccessKeys
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "access_key_ids":
			filters.AccessKeyIds = &filterValues
		case "states":
			filters.States = new(lo.Map(filterValues, func(s string, _ int) osc.AccessKeyState {
				return osc.AccessKeyState(s)
			}))
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
