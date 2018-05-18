package outscale

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleOAPIVMSState() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPIVMSStateRead,
		Schema: getOAPIVMSStateDataSourceSchema(),
	}
}

func dataSourceOutscaleOAPIVMSStateRead(d *schema.ResourceData, meta interface{}) error {
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
		params.InstanceIds = expandStringList(instanceIds.([]interface{}))
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

	filteredStates := resp.InstanceStatuses

	if len(filteredStates) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return statusDescriptionOAPIVMSStateAttributes(d, filteredStates)
}

func statusDescriptionOAPIVMSStateAttributes(d *schema.ResourceData, status []*fcu.InstanceStatus) error {

	d.SetId(resource.UniqueId())

	states := make([]map[string]interface{}, len(status))

	for k, v := range status {
		state := make(map[string]interface{})

		state["sub_region_name"] = *v.AvailabilityZone

		events := eventsSetOAPIVMSState(v.Events)
		state["maintenance_event"] = events

		st := flattenedStateOAPIVMSState(v.InstanceState)
		state["state"] = st

		st1 := statusSetOAPIVMSState(v.InstanceStatus)
		state["comment"] = st1

		states[k] = state
	}

	return d.Set("vm_state_set", states)
}

func statusSetOAPIVMSState(status *fcu.InstanceStatusSummary) map[string]interface{} {
	s := make(map[string]interface{})
	s["state"] = aws.StringValue(status.Status)
	s["item"] = detailsSetOAPIVMSState(status.Details)

	return s
}

func detailsSetOAPIVMSState(details []*fcu.InstanceStatusDetails) []map[string]interface{} {
	s := make([]map[string]interface{}, len(details))

	for k, v := range details {
		status := map[string]interface{}{
			"name":  aws.StringValue(v.Name),
			"state": aws.StringValue(v.Status),
		}
		s[k] = status
	}

	return s
}

func flattenedStateOAPIVMSState(state *fcu.InstanceState) map[string]interface{} {
	return map[string]interface{}{
		"code": fmt.Sprintf("%d", aws.Int64Value(state.Code)),
		"name": aws.StringValue(state.Name),
	}
}

func eventsSetOAPIVMSState(events []*fcu.InstanceStatusEvent) []map[string]interface{} {
	s := make([]map[string]interface{}, len(events))
	for k, v := range events {
		status := map[string]interface{}{
			"state_code":  aws.StringValue(v.Code),
			"description": aws.StringValue(v.Description),
			"not_before":  v.NotBefore.Format(time.RFC3339),
			"not_after":   v.NotAfter.Format(time.RFC3339),
		}
		s[k] = status
	}
	return s
}

func getOAPIVMSStateDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Arguments
		"filter": dataSourceFiltersSchema(),
		"vm_id": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"vm_state_set": { //events_set
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"sub_region_name": { //availability_zone
						Type:     schema.TypeString,
						Computed: true,
					},
					"maintenance_event": { //events_set
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"state_code": { //code
									Type:     schema.TypeString,
									Computed: true,
								},
								"description": { //
									Type:     schema.TypeString,
									Computed: true,
								},
								"not_after": { // not_before
									Type:     schema.TypeString,
									Computed: true,
								},
								"not_before": { // not_after
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},

					"vm_id": { //instance_id
						Type:     schema.TypeString,
						Computed: true,
					},
					"state": { //instance_state
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"state_code": { // code
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
					"comment": { // instance_status
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"item": { //details
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"name": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"state": { //status
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"state": { //state
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
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
