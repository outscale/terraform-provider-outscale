package outscale

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleOAPIVpcEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIVpcEndpointCreate,
		Read:   resourceOutscaleOAPIVpcEndpointRead,
		Update: resourceOutscaleOAPIVpcEndpointUpdate,
		Delete: resourceOutscaleOAPIVpcEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"lin_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"prefix_list_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"route_table_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"lin_api_access_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prefix_list_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_ranges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"request_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceOutscaleOAPIVpcEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.CreateVpcEndpointInput{
		VpcId:       aws.String(d.Get("lin_id").(string)),
		ServiceName: aws.String(d.Get("prefix_list_name").(string)),
	}

	setVpcEndpointCreateList(d, "route_table_id", &req.RouteTableIds)

	log.Printf("[DEBUG] Creating VPC Endpoint: %#v", req)

	var err error
	var resp *fcu.CreateVpcEndpointOutput

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.CreateVpcEndpoint(req)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating VPC Endpoint: %s", err.Error())
	}

	vpce := resp.VpcEndpoint
	d.SetId(aws.StringValue(vpce.VpcEndpointId))

	if err := vpcEndpointWaitUntilAvailable(d, conn); err != nil {
		return err
	}

	return resourceOutscaleVpcEndpointRead(d, meta)
}

func resourceOutscaleOAPIVpcEndpointRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var resp *fcu.DescribeVpcEndpointsOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeVpcEndpoints(&fcu.DescribeVpcEndpointsInput{
			VpcEndpointIds: aws.StringSlice([]string{d.Id()}),
		})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	vpce := resp.VpcEndpoints[0]
	state := *vpce.State

	if err != nil && state != "failed" {
		return fmt.Errorf("Error reading VPC Endpoint: %s", err.Error())
	}

	terminalStates := map[string]bool{
		"deleted":  true,
		"deleting": true,
		"failed":   true,
		"expired":  true,
		"rejected": true,
	}
	if _, ok := terminalStates[state]; ok {
		log.Printf("[WARN] VPC Endpoint (%s) in state (%s), removing from state", d.Id(), state)
		d.SetId("")
		return nil
	}

	d.Set("request_id", *resp.RequestId)
	return vpcEndpointAttributes(d, vpce, conn)
}

func resourceOutscaleOAPIVpcEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.ModifyVpcEndpointInput{
		VpcEndpointId: aws.String(d.Id()),
	}

	if d.HasChange("policy") {
		policy, err := structure.NormalizeJsonString(d.Get("policy"))
		if err != nil {
			return errwrap.Wrapf("policy contains an invalid JSON: {{err}}", err)
		}

		if policy == "" {
			req.ResetPolicy = aws.Bool(true)
		} else {
			req.PolicyDocument = aws.String(policy)
		}
	}

	setVpcEndpointUpdateLists(d, "route_table_id", &req.AddRouteTableIds, &req.RemoveRouteTableIds)

	log.Printf("[DEBUG] Updating VPC Endpoint: %#v", req)

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.ModifyVpcEndpoint(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error updating VPC Endpoint: %s", err.Error())
	}

	if err := vpcEndpointWaitUntilAvailable(d, conn); err != nil {
		return err
	}

	return resourceOutscaleVpcEndpointRead(d, meta)
}

func resourceOutscaleOAPIVpcEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	log.Printf("[DEBUG] Deleting VPC Endpoint: %s", d.Id())

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.DeleteVpcEndpoints(&fcu.DeleteVpcEndpointsInput{
			VpcEndpointIds: aws.StringSlice([]string{d.Id()}),
		})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidVpcEndpointId.NotFound") {
			log.Printf("[DEBUG] VPC Endpoint %s is already gone", d.Id())
		} else {
			return fmt.Errorf("Error deleting VPC Endpoint: %s", err.Error())
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "pending", "deleting"},
		Target:     []string{"deleted"},
		Refresh:    vpcEndpointStateRefresh(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for VPC Endpoint %s to delete: %s", d.Id(), err.Error())
	}

	return nil
}

func vpcEndpointStateRefreshOAPI(conn *fcu.Client, vpceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Reading VPC Endpoint: %s", vpceID)

		var resp *fcu.DescribeVpcEndpointsOutput
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeVpcEndpoints(&fcu.DescribeVpcEndpointsInput{
				VpcEndpointIds: aws.StringSlice([]string{vpceID}),
			})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidVpcEndpointId.NotFound") {
				return false, "deleted", nil
			}

			return nil, "", err
		}

		vpce := resp.VpcEndpoints[0]
		state := aws.StringValue(vpce.State)
		// No use in retrying if the endpoint is in a failed state.
		if state == "failed" {
			return nil, state, errors.New("VPC Endpoint is in a failed state")
		}
		return vpce, state, nil
	}
}

func vpcEndpointWaitUntilAvailableOAPI(d *schema.ResourceData, conn *fcu.Client) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available", "pendingAcceptance"},
		Refresh:    vpcEndpointStateRefreshOAPI(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for VPC Endpoint %s to become available: %s", d.Id(), err.Error())
	}

	return nil
}

func vpcEndpointAttributesOAPI(d *schema.ResourceData, vpce *fcu.VpcEndpoint, conn *fcu.Client) error {
	d.Set("state", vpce.State)
	d.Set("lin_id", vpce.VpcId)

	serviceName := aws.StringValue(vpce.ServiceName)
	d.Set("prefix_list_name", serviceName)
	d.Set("lin_api_access_id", aws.StringValue(vpce.VpcEndpointId))

	policy, err := structure.NormalizeJsonString(aws.StringValue(vpce.PolicyDocument))
	if err != nil {
		return errwrap.Wrapf("policy contains an invalid JSON: {{err}}", err)
	}
	d.Set("policy", policy)

	d.Set("route_table_id", flattenStringList(vpce.RouteTableIds))

	req := &fcu.DescribePrefixListsInput{}
	req.Filters = buildFCUAttributeFilterList(
		map[string]string{
			"prefix-list-name": serviceName,
		},
	)

	var resp *fcu.DescribePrefixListsOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribePrefixLists(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return err
	}
	if resp != nil && len(resp.PrefixLists) > 0 {
		if len(resp.PrefixLists) > 1 {
			return fmt.Errorf("multiple prefix lists associated with the service name '%s'. Unexpected", serviceName)
		}

		pl := resp.PrefixLists[0]
		d.Set("prefix_list_id", pl.PrefixListId)
		d.Set("ip_ranges", flattenStringList(pl.Cidrs))
	} else {
		d.Set("ip_ranges", make([]string, 0))
	}

	return nil
}

func setVpcEndpointCreateListOAPI(d *schema.ResourceData, key string, c *[]*string) {
	if v, ok := d.GetOk(key); ok {
		list := v.(*schema.Set).List()
		if len(list) > 0 {
			*c = expandStringList(list)
		}
	}
}

func setVpcEndpointUpdateListsOAPI(d *schema.ResourceData, key string, a, r *[]*string) {
	if d.HasChange(key) {
		o, n := d.GetChange(key)
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		add := expandStringList(ns.Difference(os).List())
		if len(add) > 0 {
			*a = add
		}

		remove := expandStringList(os.Difference(ns).List())
		if len(remove) > 0 {
			*r = remove
		}
	}
}

func buildFCUAttributeFilterListOAPI(attrs map[string]string) []*fcu.Filter {
	var filters []*fcu.Filter

	// sort the filters by name to make the output deterministic
	var names []string
	for filterName := range attrs {
		names = append(names, filterName)
	}

	sort.Strings(names)

	for _, filterName := range names {
		value := attrs[filterName]
		if value == "" {
			continue
		}

		filters = append(filters, &fcu.Filter{
			Name:   aws.String(filterName),
			Values: []*string{aws.String(value)},
		})
	}

	return filters
}
