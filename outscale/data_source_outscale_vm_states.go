package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOutscaleOAPIVMStates() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPIVMStatesRead,
		Schema: getOAPIVMStatesDataSourceSchema(),
	}
}

func getOAPIVMStatesDataSourceSchema() map[string]*schema.Schema {
	wholeSchema := map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
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

func dataSourceOutscaleOAPIVMStatesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.ReadVmsStateRequest{}
	if filters, filtersOk := d.GetOk("filter"); filtersOk {
		req.SetFilters(buildOutscaleOAPIDataSourceVMStateFilters(filters.(*schema.Set)))
	}
	var resp oscgo.ReadVmsStateResponse
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.VmApi.ReadVmsState(context.Background()).ReadVmsStateRequest(req).Execute()
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
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return statusDescriptionOAPIVMStatesAttributes(d, filteredStates)
}

func statusDescriptionOAPIVMStatesAttributes(d *schema.ResourceData, status []oscgo.VmStates) error {
	d.SetId(resource.UniqueId())

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
