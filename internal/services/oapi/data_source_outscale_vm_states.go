package oapi

import (
	"context"
	"errors"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
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
	conn := meta.(*client.OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	instanceIds, instanceIdsOk := d.GetOk("vm_ids")

	if !instanceIdsOk && !filtersOk {
		return errors.New("vm_id or filter must be set")
	}

	var err error
	params := oscgo.ReadVmsStateRequest{}
	if filtersOk {
		params.Filters, err = buildOutscaleDataSourceVMStateFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	if instanceIdsOk {
		filter := oscgo.FiltersVmsState{}
		filter.SetVmIds(utils.InterfaceSliceToStringSlice(instanceIds.([]interface{})))
		params.SetFilters(filter)
	}
	params.SetAllVms(d.Get("all_vms").(bool))
	var resp oscgo.ReadVmsStateResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.VmApi.ReadVmsState(context.Background()).ReadVmsStateRequest(params).Execute()
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

func statusDescriptionOAPIVMStatesAttributes(d *schema.ResourceData, status []oscgo.VmStates) error {
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
