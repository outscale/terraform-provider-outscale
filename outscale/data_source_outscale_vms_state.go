package outscale

import (
	"context"
	"errors"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	instanceIds, instanceIdsOk := d.GetOk("vm_ids")

	if !instanceIdsOk && !filtersOk {
		return errors.New("vm_id or filter must be set")
	}

	params := oscgo.ReadVmsStateRequest{}
	if filtersOk {
		params.SetFilters(buildOutscaleOAPIDataSourceVMStateFilters(filters.(*schema.Set)))
	}
	if instanceIdsOk {
		filter := oscgo.FiltersVmsState{}
		filter.SetVmIds(oapiExpandStringList(instanceIds.([]interface{})))
		params.SetFilters(filter)
	}

	var resp oscgo.ReadVmsStateResponse
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.VmApi.ReadVmsState(context.Background(), &oscgo.ReadVmsStateOpts{ReadVmsStateRequest: optional.NewInterface(params)})
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

	filteredStates := resp.GetVmStates()[:]

	if len(filteredStates) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}
	d.Set("request_id", resp.ResponseContext.GetRequestId())
	return statusDescriptionOAPIVMSStateAttributes(d, filteredStates)
}

func statusDescriptionOAPIVMSStateAttributes(d *schema.ResourceData, status []oscgo.VmStates) error {
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
