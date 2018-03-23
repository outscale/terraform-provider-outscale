package outscale

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleVMSState() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleVMSStateRead,
		Schema: getVMSStateDataSourceSchema(),
	}
}

func dataSourceOutscaleVMSStateRead(d *schema.ResourceData, meta interface{}) error {
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

	if len(filteredStates) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	states := filteredStates

	log.Printf("[DEBUG] outscale_vms_state - states found: %s", spew.Sdump(filteredStates))

	return statusesDescriptionAttributes(d, states)
}

func statusesDescriptionAttributes(d *schema.ResourceData, status []*fcu.InstanceStatus) error {

	d.SetId(resource.UniqueIdPrefix)

	statuses := make([]map[string]interface{}, len(status))

	for i, s := range status {
		statuses[i] = map[string]interface{}{
			"instance_id":       *s.InstanceId,
			"availability_zone": *s.AvailabilityZone,
			"events_set":        eventsSet(s.Events),
			"instance_state":    flattenedState(s.InstanceState),
			"instance_status":   statusSet(s.InstanceStatus),
			"system_status":     statusSet(s.SystemStatus),
		}
	}

	err := d.Set("instance_status_set", statuses)

	return err
}

func statusSetHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["instance_id"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["availability_zone"].(string)))
	return hashcode.String(buf.String())
}

func getVMSStateDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Arguments
		"filter": dataSourceFiltersSchema(),
		"instance_id": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"include_all_instances": {
			Type:     schema.TypeBool,
			Optional: true,
		},

		// Attributes
		"instance_status_set": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"availability_zone": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"events_set": {
						Type:     schema.TypeSet,
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

					// Need to check this
					"instance_id": {
						Type:     schema.TypeString,
						Optional: true,
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
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"details": {
									Type:     schema.TypeSet,
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
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"details": {
									Type:     schema.TypeSet,
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
				},
			},
		},

		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
