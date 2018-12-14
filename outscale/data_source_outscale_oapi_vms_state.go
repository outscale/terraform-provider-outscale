package outscale

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func dataSourceOutscaleOAPIVMSState() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPIVMSStateRead,
		Schema: getOAPIVMSStateDataSourceSchema(),
	}
}

func getOAPIVMSStateDataSourceSchema() map[string]*schema.Schema {
	wholeSchema := map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
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

func dataSourceOutscaleOAPIVMSStateRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	filters, filtersOk := d.GetOk("filter")
	instanceIds, instanceIdsOk := d.GetOk("vm_ids")

	if !instanceIdsOk && !filtersOk {
		return errors.New("vm_id or filter must be set")
	}

	params := oapi.ReadVmsStateRequest{}
	if filtersOk {
		params.Filters = buildOutscaleOAPIDataSourceVmStateFilters(filters.(*schema.Set))
	}
	if instanceIdsOk {
		params.Filters.VmIds = oapiExpandStringList(instanceIds.([]interface{}))
	}

	var resp *oapi.POST_ReadVmsStateResponses
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_ReadVmsState(params)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		return err
	}

	filteredStates := resp.OK.VmStates[:]

	if len(filteredStates) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return statusDescriptionOAPIVMSStateAttributes(d, filteredStates)
}

func statusDescriptionOAPIVMSStateAttributes(d *schema.ResourceData, status []oapi.VmStates) error {

	d.SetId(resource.UniqueId())

	states := make([]map[string]interface{}, len(status))

	for k, v := range status {
		state := make(map[string]interface{})

		setterFunc := func(key string, value interface{}) error {
			state[key] = value
			return nil
		}

		statusDescriptionOAPIVMStateAttributes(setterFunc, &v)

		states[k] = state
	}

	return d.Set("vm_states", states)
}
