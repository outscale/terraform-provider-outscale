package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleOAPIInternetService() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIInternetServiceCreate,
		Read:   resourceOutscaleOAPIInternetServiceRead,
		Update: resourceOutscaleOAPIInternetServiceUpdate,
		Delete: resourceOutscaleOAPIInternetServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsListOAPISchema(),
		},
	}
}

func resourceOutscaleOAPIInternetServiceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	var resp oscgo.CreateInternetServiceResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.InternetServiceApi.CreateInternetService(context.Background()).CreateInternetServiceRequest(oscgo.CreateInternetServiceRequest{}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error creating Internet Service: %s", utils.GetErrorResponse(err))
	}

	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.(*schema.Set), resp.InternetService.GetInternetServiceId(), conn)
		if err != nil {
			return err
		}
	}

	d.SetId(resp.InternetService.GetInternetServiceId())

	return resourceOutscaleOAPIInternetServiceRead(d, meta)
}

func resourceOutscaleOAPIInternetServiceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	id := d.Id()

	log.Printf("[DEBUG] Reading Internet Service id (%s)", id)

	req := oscgo.ReadInternetServicesRequest{
		Filters: &oscgo.FiltersInternetService{InternetServiceIds: &[]string{id}},
	}

	var resp oscgo.ReadInternetServicesResponse

	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		r, httpResp, err := conn.InternetServiceApi.ReadInternetServices(context.Background()).ReadInternetServicesRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = r
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error reading Internet Service id (%s)", utils.GetErrorResponse(err))

	}
	if utils.IsResponseEmpty(len(resp.GetInternetServices()), "InternetService", d.Id()) {
		d.SetId("")
		return nil
	}
	if err := d.Set("internet_service_id", resp.GetInternetServices()[0].GetInternetServiceId()); err != nil {
		return err
	}

	if err := d.Set("net_id", resp.GetInternetServices()[0].GetNetId()); err != nil {
		return err
	}

	if err := d.Set("state", resp.GetInternetServices()[0].GetState()); err != nil {
		return err
	}

	return d.Set("tags", tagsOSCAPIToMap(resp.GetInternetServices()[0].GetTags()))
}

func resourceOutscaleOAPIInternetServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}
	return resourceOutscaleOAPIInternetServiceRead(d, meta)
}

func resourceOutscaleOAPIInternetServiceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	internetServiceID := d.Id()
	req := oscgo.DeleteInternetServiceRequest{
		InternetServiceId: internetServiceID,
	}

	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.InternetServiceApi.DeleteInternetService(context.Background()).DeleteInternetServiceRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("[DEBUG] Error deleting Internet Service id (%s)", err)
	}
	return nil
}
