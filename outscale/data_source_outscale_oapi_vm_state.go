package outscale

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleOAPIVMState() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPIVMStateRead,
		Schema: getOAPIVMStateDataSourceSchema(),
	}
}

func dataSourceOutscaleOAPIVMStateRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	instanceIds, instanceIdsOk := d.GetOk("vm_id")

	if !instanceIdsOk && !filtersOk {
		return errors.New("vm_id or filter must be set")
	}

	params := &fcu.DescribeInstanceStatusInput{}
	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if instanceIdsOk {
		var ids []*string

		for _, id := range instanceIds.(*schema.Set).List() {
			ids = append(ids, aws.String(id.(string)))
		}

		params.InstanceIds = ids
	}

	params.IncludeAllInstances = aws.Bool(false)

	var resp *fcu.DescribeInstanceStatusOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeInstanceStatus(params)
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

	filteredStates := resp.InstanceStatuses[:]

	var state *fcu.InstanceStatus
	if len(filteredStates) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	if len(filteredStates) > 1 {
		return fmt.Errorf("Your query returned more than one result. Please try a more " +
			"specific search criteria.")
	}

	state = filteredStates[0]

	log.Printf("[DEBUG] outscale_oapi_vm_state - Single State found: %s", *state.InstanceId)

	return statusDescriptionOAPIVMStateAttributes(d, state)
}

func statusDescriptionOAPIVMStateAttributes(d *schema.ResourceData, status *fcu.InstanceStatus) error {

	d.SetId(*status.InstanceId)

	d.Set("sub_region_name", status.AvailabilityZone)

	events := eventsSet(status.Events)
	err := d.Set("maintenance_event", events)
	if err != nil {
		return err
	}

	state := flattenedState(status.InstanceState)
	err = d.Set("state", state)
	if err != nil {
		return err
	}

	st := statusSet(status.InstanceStatus)
	err = d.Set("comment", st)
	if err != nil {
		return err
	}

	return nil
}

func statusSetOAPIVMState(status *fcu.InstanceStatusSummary) map[string]interface{} {

	st := map[string]interface{}{
		"state": *status.Status,
		"item":  detailsSetOAPIVMState(status.Details),
	}

	return st
}

func detailsSetOAPIVMState(details []*fcu.InstanceStatusDetails) []map[string]interface{} {
	s := make([]map[string]interface{}, len(details))

	for k, v := range details {

		status := map[string]interface{}{
			"name":  *v.Name,
			"state": *v.Status,
		}
		s[k] = status
	}

	return s
}

func flattenedStateOAPIVMState(state *fcu.InstanceState) map[string]interface{} {
	return map[string]interface{}{
		"code": fmt.Sprintf("%d", *state.Code),
		"name": *state.Name,
	}
}

func eventsSetOAPIVMState(events []*fcu.InstanceStatusEvent) []map[string]interface{} {

	s := make([]map[string]interface{}, len(events))

	for k, v := range events {

		status := map[string]interface{}{
			"state_code":  *v.Code,
			"description": *v.Description,
			"not_before":  v.NotBefore.Format(time.RFC3339),
			"not_after":   v.NotAfter.Format(time.RFC3339),
		}
		s[k] = status
	}
	return s
}

func getOAPIVMStateDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
		"sub_region_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"maintenance_event": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"state_code": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"description": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"not_after": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"not_before": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},

		"vm_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"state": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"state_code": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"comment": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"item": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
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
	}
}
