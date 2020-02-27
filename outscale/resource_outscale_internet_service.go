package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/optional"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	oscgo "github.com/marinsalinas/osc-sdk-go"
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

	resp, _, err := conn.InternetServiceApi.CreateInternetService(context.Background(), &oscgo.CreateInternetServiceOpts{
		CreateInternetServiceRequest: optional.NewInterface(oscgo.CreateInternetServiceRequest{}),
	})
	if err != nil {
		return fmt.Errorf("[DEBUG] Error creating Internet Service: %s", utils.GetErrorResponse(err))
	}

	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.([]interface{}), resp.InternetService.GetInternetServiceId(), conn)
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

	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		r, _, err := conn.InternetServiceApi.ReadInternetServices(context.Background(), &oscgo.ReadInternetServicesOpts{ReadInternetServicesRequest: optional.NewInterface(req)})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		resp = r
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error reading Internet Service id (%s)", utils.GetErrorResponse(err))

	}
	if !resp.HasInternetServices() {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	if err := d.Set("request_id", resp.ResponseContext.GetRequestId()); err != nil {
		return err
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

	d.Partial(true)

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

	d.Partial(false)
	return resourceOutscaleOAPIInternetServiceRead(d, meta)
}

func resourceOutscaleOAPIInternetServiceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.DeleteInternetServiceRequest{
		InternetServiceId: d.Id(),
	}

	_, _, err := conn.InternetServiceApi.DeleteInternetService(context.Background(), &oscgo.DeleteInternetServiceOpts{
		DeleteInternetServiceRequest: optional.NewInterface(req),
	})
	if err != nil {
		return fmt.Errorf("[DEBUG] Error deleting Internet Service id (%s)", err)
	}

	return nil
}
