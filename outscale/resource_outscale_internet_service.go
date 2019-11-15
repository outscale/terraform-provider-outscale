package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	"github.com/outscale/osc-go/oapi"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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

		Schema: getOAPIInternetServiceSchema(),
	}
}

func resourceOutscaleOAPIInternetServiceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	log.Println("[DEBUG] Creating Internet Service")
	resp, _, err := conn.InternetServiceApi.CreateInternetService(context.Background(), &oscgo.CreateInternetServiceOpts{CreateInternetServiceRequest: optional.NewInterface(oscgo.CreateInternetServiceRequest{})})

	var errString string

	if err != nil {
		errString = err.Error()

		return fmt.Errorf("[DEBUG] Error creating Internet Service: %s", errString)
	}

	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.([]interface{}), resp.InternetService.GetInternetServiceId(), conn)
		if err != nil {
			return err
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "ending/wait"},
		Target:     []string{"available"},
		Refresh:    InternetServiceStateOApiRefreshFunc(conn, resp.InternetService.GetInternetServiceId()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to become created: %s", d.Id(), err)
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
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
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
		return fmt.Errorf("[DEBUG] Error reading Internet Service id (%s)", err.Error())

	}
	if !resp.HasInternetServices() {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}
	// Workaround to get the desired internet_service instance. TODO: Remove getInternetService
	// once filters work again. And use resp.OK.InternetServices[0]
	err, result := getInternetServiceOSC(resp.GetInternetServices(), id)

	d.Set("request_id", resp.ResponseContext.GetRequestId())
	d.Set("internet_service_id", result.GetInternetServiceId())

	if err := d.Set("net_id", result.GetNetId()); err != nil {
		return err
	}

	if err := d.Set("state", result.GetState()); err != nil {
		return err
	}

	return d.Set("tags", tagsOSCAPIToMap(result.GetTags()))
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

func getInternetService(internetServices []oapi.InternetService, id string) (oapi.InternetService, error) {
	for _, element := range internetServices {
		if element.InternetServiceId == id {
			return element, nil
		}
	}
	return oapi.InternetService{}, fmt.Errorf("InternetService %+s not found", id)
}

func getInternetServiceOSC(internetServices []oscgo.InternetService, id string) (error, oscgo.InternetService) {
	for _, element := range internetServices {
		if *element.InternetServiceId == id {
			return nil, element
		}
	}
	return fmt.Errorf("InternetService %+s not found", id), oscgo.InternetService{}
}

func resourceOutscaleOAPIInternetServiceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	id := d.Id()
	log.Printf("[DEBUG] Deleting Internet Service id (%s)", id)

	req := oscgo.DeleteInternetServiceRequest{
		InternetServiceId: id,
	}

	var err error

	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, _, err := conn.InternetServiceApi.DeleteInternetService(context.Background(), &oscgo.DeleteInternetServiceOpts{DeleteInternetServiceRequest: optional.NewInterface(req)})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})

	var errString string

	if err != nil {
		errString = err.Error()

		return fmt.Errorf("[DEBUG] Error deleting Internet Service id (%s)", errString)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "ending/wait"},
		Target:     []string{"not available"},
		Refresh:    InternetServiceStateOApiRefreshFunc(conn, id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	return nil
}

func InternetServiceStateOApiRefreshFunc(conn *oscgo.APIClient, internetServiceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, _, err := conn.InternetServiceApi.ReadInternetServices(context.Background(), &oscgo.ReadInternetServicesOpts{
			ReadInternetServicesRequest: optional.NewInterface(oscgo.ReadInternetServicesRequest{
				Filters: &oscgo.FiltersInternetService{
					InternetServiceIds: &[]string{internetServiceID},
				},
			}),
		})

		if err != nil {
			log.Printf("[ERROR] error on InstanceStateRefresh: %s", err)
			return nil, "", err
		}

		if !resp.HasInternetServices() || len(resp.GetInternetServices()) == 0 {
			return nil, "not available", nil
		}

		internetService := resp.GetInternetServices()[0]

		return internetService, "available", nil
	}
}

func getOAPIInternetServiceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
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
	}
}
