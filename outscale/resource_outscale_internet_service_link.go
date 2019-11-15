package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/optional"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	oscgo "github.com/marinsalinas/osc-sdk-go"
)

func resourceOutscaleOAPIInternetServiceLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIInternetServiceLinkCreate,
		Read:   resourceOutscaleOAPIInternetServiceLinkRead,
		Delete: resourceOutscaleOAPIInternetServiceLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// Arguments
			"net_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"internet_service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// Attributes
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": dataSourceTagsSchema(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPIInternetServiceLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	log.Println("[DEBUG] Creating Internet Service")
	resp, _, err := conn.InternetServiceApi.CreateInternetService(context.Background(), &oscgo.CreateInternetServiceOpts{
		CreateInternetServiceRequest: optional.NewInterface(oscgo.CreateInternetServiceRequest{}),
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error creating Internet Service: %s", err)
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

	return resourceOutscaleOAPIInternetServiceLinkRead(d, meta)
}

func resourceOutscaleOAPIInternetServiceLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	log.Printf("[DEBUG] Reading Internet Service id (%s)", d.Id())

	req := oscgo.ReadInternetServicesRequest{
		Filters: &oscgo.FiltersInternetService{InternetServiceIds: &[]string{d.Id()}},
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
	err, result := getInternetServiceOSC(resp.GetInternetServices(), d.Id())

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

func resourceOutscaleOAPIInternetServiceLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := &oscgo.UnlinkInternetServiceRequest{
		InternetServiceId: d.Get("internet_service_id").(string),
		NetId:             d.Get("net_id").(string),
	}

	var err error

	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, _, err := conn.InternetServiceApi.DeleteInternetService(context.Background(), &oscgo.DeleteInternetServiceOpts{
			DeleteInternetServiceRequest: optional.NewInterface(req),
		})

		if err != nil {
			if strings.Contains(err.Error(), "DependencyProblem") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var errString string

	if err != nil {
		errString = err.Error()

		return fmt.Errorf("[DEBUG] Error deleting Internet Service id (%s)", errString)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "ending/wait"},
		Target:     []string{"not available"},
		Refresh:    InternetServiceStateOApiRefreshFunc(conn, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	return nil
}
