package oapi

import (
	"errors"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleVMStates() *schema.Resource {
	return &schema.Resource{
		Read:   DataSourceOutscaleVMStatesRead,
		Schema: getOAPIVMStatesDataSourceSchema(),
	}
}

func getOAPIVMStatesDataSourceSchema() map[string]*schema.Schema {
	wholeSchema := map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
		"all_vms": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"vm_ids": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"vm_states": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: getVMStateAttrsSchema(),
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	return wholeSchema
}

func DataSourceOutscaleVMStatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	instanceIds, instanceIdsOk := d.GetOk("vm_ids")

	if !instanceIdsOk && !filtersOk {
		return errors.New("vm_id or filter must be set")
	}

	var err error
	params := osc.ReadVmsStateRequest{}
	if filtersOk {
		params.Filters, err = buildOutscaleDataSourceVMStateFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	if instanceIdsOk {
		filter := osc.FiltersVmsState{}
		filter.SetVmIds(utils.InterfaceSliceToStringSlice(instanceIds.([]interface{})))
		params.SetFilters(filter)
	}
	params.SetAllVms(d.Get("all_vms").(bool))
	var resp osc.ReadVmsStateResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := client.VmApi.ReadVmsState(ctx).ReadVmsStateRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	filteredStates := resp.GetVmStates()[:]

	if len(filteredStates) < 1 {
		return ErrNoResults
	}

	return statusDescriptionOAPIVMStatesAttributes(d, filteredStates)
}

func statusDescriptionOAPIVMStatesAttributes(d *schema.ResourceData, status []osc.VmStates) error {
	d.SetId(id.UniqueId())

	states := make([]map[string]interface{}, len(status))

	for k, v := range status {
		state := make(map[string]interface{})

		setterFunc := func(key string, value interface{}) error {
			state[key] = value
			return nil
		}

		if err := statusDescriptionOAPIVMStateAttributes(setterFunc, &v); err != nil {
			return err
		}

		states[k] = state
	}

	return d.Set("vm_states", states)
}
