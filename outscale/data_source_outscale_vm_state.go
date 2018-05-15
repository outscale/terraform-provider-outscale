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

func dataSourceOutscaleVMState() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleVMStateRead,
		Schema: getVMStateDataSourceSchema(),
	}
}

func dataSourceOutscaleVMStateRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	instanceIds, instanceIdsOk := d.GetOk("instance_id")

	if !instanceIdsOk && !filtersOk {
		return errors.New("instance_id or filter must be set")
	}

	params := &fcu.DescribeInstanceStatusInput{}
	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if instanceIdsOk {
		params.InstanceIds = []*string{aws.String(instanceIds.(string))}
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
		return nil
	})

	if err != nil {
		return err
	}

	filteredStates := resp.InstanceStatuses[:]

	var state *fcu.InstanceStatus
	if len(filteredStates) < 1 {
		return fmt.Errorf("our query returned no results, please change your search criteria and try again")
	}

	if len(filteredStates) > 1 {
		return fmt.Errorf("our query returned more than one result, please try a more " +
			"specific search criteria")
	}

	state = filteredStates[0]

	log.Printf("[DEBUG] outscale_vm_state - Single State found: %s", *state.InstanceId)

	d.Set("request_id", *resp.RequestId)

	return statusDescriptionAttributes(d, state)
}

func statusDescriptionAttributes(d *schema.ResourceData, status *fcu.InstanceStatus) error {

	d.SetId(*status.InstanceId)

	d.Set("availability_zone", aws.StringValue(status.AvailabilityZone))

	events := eventsSet(status.Events)
	if err := d.Set("events_set", events); err != nil {
		return err
	}

	state := flattenedState(status.InstanceState)
	if err := d.Set("instance_state", state); err != nil {
		return err
	}

	st := statusSet(status.InstanceStatus)
	if err := d.Set("instance_status", st); err != nil {
		return err
	}

	sst := statusSet(status.SystemStatus)

	return d.Set("system_status", sst)
}

func statusSet(status *fcu.InstanceStatusSummary) []map[string]interface{} {
	st := make([]map[string]interface{}, 1)

	s := make(map[string]interface{})
	s["status"] = aws.StringValue(status.Status)
	s["details"] = detailsSet(status.Details)

	st[0] = s

	return st
}

func detailsSet(details []*fcu.InstanceStatusDetails) []map[string]interface{} {
	s := make([]map[string]interface{}, len(details))

	for k, v := range details {

		status := map[string]interface{}{
			"name":   aws.StringValue(v.Name),
			"status": aws.StringValue(v.Status),
		}
		s[k] = status
	}

	return s
}

func flattenedState(state *fcu.InstanceState) map[string]interface{} {
	return map[string]interface{}{
		"code": fmt.Sprintf("%d", aws.Int64Value(state.Code)),
		"name": aws.StringValue(state.Name),
	}
}

func eventsSet(events []*fcu.InstanceStatusEvent) []map[string]interface{} {
	s := make([]map[string]interface{}, len(events))

	for k, v := range events {
		status := map[string]interface{}{
			"code":        aws.StringValue(v.Code),
			"description": aws.StringValue(v.Description),
			"not_before":  v.NotBefore.Format(time.RFC3339),
			"not_after":   v.NotAfter.Format(time.RFC3339),
		}
		s[k] = status
	}
	return s
}

func getVMStateDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Arguments
		"filter": dataSourceFiltersSchema(),
		"instance_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"include_all_instances": {
			Type:     schema.TypeBool,
			Optional: true,
		},

		// Attributes
		"availability_zone": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"events_set": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"code": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"description": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"not_before": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"not_after": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"instance_state": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"code": {
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
		"instance_status": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"details": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"status": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"status": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"system_status": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"details": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"details": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"status": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"status": {
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
