package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleNetAccessPoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleNetAccessPointCreate,
		Read:   resourceOutscaleNetAccessPointRead,
		Delete: resourceOutscaleNetAccessPointDelete,
		Update: resourceOutscaleNetAccessPointUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"net_access_point_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"route_table_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"service_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsListOAPISchema(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleNetAccessPointUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	if d.HasChange("route_table_ids") {
		o, n := d.GetChange("route_table_ids")

		log.Printf("[DEBUG] it change !: %v %v", o, n)
		oo := utils.SetToStringSlicePtr(o.(*schema.Set))
		nn := utils.SetToStringSlicePtr(n.(*schema.Set))
		destroy := make([]string, 0)
		add := make([]string, 0)

		for _, v := range *oo {
			to_destroy := true
			for _, v2 := range *nn {
				if v2 == v {
					to_destroy = false
					break
				}
			}
			if to_destroy {
				destroy = append(destroy, v)
			}
		}

		for _, v := range *nn {
			to_add := true
			for _, v2 := range *oo {
				if v2 == v {
					to_add = false
					break
				}
			}
			if to_add {
				add = append(add, v)
			}
		}

		req := &oscgo.UpdateNetAccessPointRequest{
			AddRouteTableIds:    &add,
			RemoveRouteTableIds: &destroy,
			NetAccessPointId:    d.Id(),
		}

		var err error
		err = resource.Retry(60*time.Second, func() *resource.RetryError {
			_, httpResp, err := conn.NetAccessPointApi.UpdateNetAccessPoint(context.Background()).UpdateNetAccessPointRequest(*req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		if err := setOSCAPITags(conn, d); err != nil {
			return err
		}
	}
	return resourceOutscaleNetAccessPointRead(d, meta)
}

func napStateRefreshFunc(conn *oscgo.APIClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp oscgo.ReadNetAccessPointsResponse
		var err error

		err = resource.Retry(60*time.Second, func() *resource.RetryError {
			rp, httpResp, err := conn.NetAccessPointApi.
				ReadNetAccessPoints(context.Background()).
				ReadNetAccessPointsRequest(oscgo.ReadNetAccessPointsRequest{
					Filters: &oscgo.FiltersNetAccessPoint{
						NetAccessPointIds: &[]string{id},
					},
				}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil {
			log.Printf("[ERROR] error on NetAccessPointStateRefresh: %s", err)
			return nil, "", err
		}

		if !resp.HasNetAccessPoints() {
			return nil, "", nil
		}

		nap := resp.GetNetAccessPoints()[0]
		state := nap.GetState()

		return nap, state, nil
	}
}

func resourceOutscaleNetAccessPointCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := &oscgo.CreateNetAccessPointRequest{}

	if v, ok := d.GetOk("route_table_ids"); ok {
		req.RouteTableIds = utils.SetToStringSlicePtr(v.(*schema.Set))
	}

	nid := d.Get("net_id")
	req.SetNetId(nid.(string))

	sn := d.Get("service_name")
	req.SetServiceName(sn.(string))

	var resp oscgo.CreateNetAccessPointResponse
	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.NetAccessPointApi.CreateNetAccessPoint(
			context.Background()).
			CreateNetAccessPointRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	//SetTags
	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.(*schema.Set), resp.NetAccessPoint.GetNetAccessPointId(), conn)
		if err != nil {
			return err
		}
	}

	id := *resp.NetAccessPoint.NetAccessPointId
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available"},
		Refresh:    napStateRefreshFunc(conn, id),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for instance (%s) to become ready: %s", id, err)
	}

	d.Set("net_access_point_id", id)
	d.SetId(id)

	return resourceOutscaleNetAccessPointRead(d, meta)
}

func resourceOutscaleNetAccessPointRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	napid := d.Id()

	filter := &oscgo.FiltersNetAccessPoint{
		NetAccessPointIds: &[]string{napid},
	}

	req := &oscgo.ReadNetAccessPointsRequest{
		Filters: filter,
	}

	var resp oscgo.ReadNetAccessPointsResponse
	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.NetAccessPointApi.ReadNetAccessPoints(
			context.Background()).
			ReadNetAccessPointsRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}
	if utils.IsResponseEmpty(len(resp.GetNetAccessPoints()), "NetAccessPoint", d.Id()) {
		d.SetId("")
		return nil
	}
	nap := (*resp.NetAccessPoints)[0]

	d.Set("route_table_ids", utils.StringSlicePtrToInterfaceSlice(nap.RouteTableIds))
	d.Set("net_id", nap.NetId)
	d.Set("service_name", nap.ServiceName)
	d.Set("state", nap.State)
	d.Set("tags", tagsOSCAPIToMap(nap.GetTags()))
	d.Set("net_access_point_id", nap.GetNetAccessPointId())

	return nil
}

func resourceOutscaleNetAccessPointDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := &oscgo.DeleteNetAccessPointRequest{
		NetAccessPointId: d.Id(),
	}

	var err error

	err = resource.Retry(70*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.NetAccessPointApi.DeleteNetAccessPoint(
			context.Background()).
			DeleteNetAccessPointRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "available", "deleting"},
		Target:     []string{"deleted"},
		Refresh:    napStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for instance (%s) to become ready: %s", d.Id(), err)
	}
	d.SetId("")
	return nil

}
