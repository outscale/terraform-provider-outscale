package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

	resp, _, err := conn.InternetServiceApi.LinkInternetService(context.Background()).LinkInternetServiceRequest(req).Execute()
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
		MinTimeout: 30 * time.Second,
		// Delay:      3 * time.Minute,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for NAT Service (%s) to become available: %s", internetServiceID, err)
	}

	d.SetId(internetServiceID)
	if err := d.Set("request_id", resp.ResponseContext.GetRequestId()); err != nil {
		return err
	}

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
		MinTimeout: 30 * time.Second,
		// Delay:      3 * time.Minute,
	}

	value, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for NAT Service (%s) to become available: %s", d.Id(), err)
	}

	resp := value.(oscgo.ReadInternetServicesResponse)
	if !resp.HasInternetServices() || len(resp.GetInternetServices()) == 0 {
		return fmt.Errorf("Error retrieving Internet Service Link: not found")
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

		return d.Set("request_id", resp.ResponseContext.RequestId)
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
		Timeout:    10 * time.Minute,
		MinTimeout: 30 * time.Second,
		// Delay:      3 * time.Minute,
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

	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		_, res, err := conn.InternetServiceApi.UnlinkInternetService(context.Background()).UnlinkInternetServiceRequest(req).Execute()
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Failed Dependency") {
				return resource.RetryableError(err) // retry
			}
			return resource.NonRetryableError(err)
		}
		if res != nil && res.StatusCode == 424 {
			return resource.RetryableError(err) // retry
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
		resp, _, err := client.InternetServiceApi.ReadInternetServices(context.Background()).ReadInternetServicesRequest(req).Execute()
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
