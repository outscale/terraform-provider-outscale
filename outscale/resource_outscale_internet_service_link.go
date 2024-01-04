package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
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
		},
	}
}

func resourceOutscaleOAPIInternetServiceLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	internetServiceID := d.Get("internet_service_id").(string)
	req := oscgo.LinkInternetServiceRequest{
		InternetServiceId: internetServiceID,
		NetId:             d.Get("net_id").(string),
	}
	var resp oscgo.LinkInternetServiceResponse
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.InternetServiceApi.LinkInternetService(context.Background()).LinkInternetServiceRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error Link Internet Service: %s", err.Error())
	}

	if !resp.HasResponseContext() {
		return fmt.Errorf("Error there is not Link Internet Service (%s)", err)
	}

	filterReq := oscgo.ReadInternetServicesRequest{
		Filters: &oscgo.FiltersInternetService{InternetServiceIds: &[]string{internetServiceID}},
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available"},
		Refresh:    LISOAPIStateRefreshFunction(conn, filterReq, "failed"),
		Timeout:    10 * time.Minute,
		MinTimeout: 20 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for NAT Service (%s) to become available: %s", internetServiceID, err)
	}

	d.SetId(internetServiceID)

	return resourceOutscaleOAPIInternetServiceLinkRead(d, meta)
}

func resourceOutscaleOAPIInternetServiceLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	internetServiceID := d.Id()
	filterReq := oscgo.ReadInternetServicesRequest{
		Filters: &oscgo.FiltersInternetService{InternetServiceIds: &[]string{internetServiceID}},
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"deleted", "available"},
		Refresh:    LISOAPIStateRefreshFunction(conn, filterReq, "failed"),
		Timeout:    10 * time.Minute,
		MinTimeout: 20 * time.Second,
	}

	value, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for NAT Service (%s) to become available: %s", d.Id(), err)
	}
	resp := value.(oscgo.ReadInternetServicesResponse)
	if utils.IsResponseEmpty(len(resp.GetInternetServices()), "InternetServiceLink", d.Id()) {
		d.SetId("")
		return nil
	}
	internetService := resp.GetInternetServices()[0]

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(internetService.GetInternetServiceId())

		if err := set("internet_service_id", internetService.GetInternetServiceId()); err != nil {
			return err
		}
		if err := set("net_id", internetService.GetNetId()); err != nil {
			return err
		}
		if err := set("state", internetService.GetState()); err != nil {
			return err
		}

		if err := set("tags", getOapiTagSet(internetService.Tags)); err != nil {
			return err
		}

		return nil
	})
}

func resourceOutscaleOAPIInternetServiceLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	internetServiceID := d.Id()
	filterReq := oscgo.ReadInternetServicesRequest{
		Filters: &oscgo.FiltersInternetService{InternetServiceIds: &[]string{internetServiceID}},
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"deleted", "available"},
		Refresh:    LISOAPIStateRefreshFunction(conn, filterReq, "failed"),
		Timeout:    5 * time.Minute,
		MinTimeout: 10 * time.Second,
	}

	value, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for NAT Service (%s) to become available: %s", d.Id(), err)
	}

	resp := value.(oscgo.ReadInternetServicesResponse)
	internetService := resp.GetInternetServices()[0]

	req := oscgo.UnlinkInternetServiceRequest{
		InternetServiceId: internetService.GetInternetServiceId(),
		NetId:             internetService.GetNetId(),
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.InternetServiceApi.UnlinkInternetService(
			context.Background()).
			UnlinkInternetServiceRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error unlink Internet Service (%s):  %s", d.Id(), err)
	}
	return nil
}

// LISOAPIStateRefreshFunction returns a resource.StateRefreshFunc that is used to watch
// a Link Internet Service.
func LISOAPIStateRefreshFunction(client *oscgo.APIClient, req oscgo.ReadInternetServicesRequest, failState string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp oscgo.ReadInternetServicesResponse
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			rp, httpResp, err := client.InternetServiceApi.ReadInternetServices(context.Background()).ReadInternetServicesRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil {
			return nil, "failed", err
		}

		state := "deleted"
		if resp.HasInternetServices() && len(resp.GetInternetServices()) > 0 {
			natServices := resp.GetInternetServices()
			state = natServices[0].GetState()

			if state == failState {
				return natServices[0], state, fmt.Errorf("Failed to reach target state. Reason: %v", state)
			}

			if state == "" {
				return resp, "available", nil
			}
		}
		return resp, state, nil
	}
}
