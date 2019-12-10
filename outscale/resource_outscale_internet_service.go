package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/outscale/osc-go/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleOAPIInternetService() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIInternetServiceCreate,
		Read:   resourceOutscaleOAPIInternetServiceRead,
		Delete: resourceOutscaleOAPIInternetServiceDelete,
		Update: resourceOutscaleOAPIInternetServiceUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: getOAPIInternetServiceSchema(),
	}
}

func resourceOutscaleOAPIInternetServiceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	log.Println("[DEBUG] Creating Internet Service")
	r, err := conn.POST_CreateInternetService(oapi.CreateInternetServiceRequest{})

	var errString string

	if err != nil || r.OK == nil {
		if err != nil {
			errString = err.Error()

		} else if r.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(r.Code401))
		} else if r.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(r.Code400))
		} else if r.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(r.Code500))
		}

		return fmt.Errorf("[DEBUG] Error creating Internet Service: %s", errString)
	}

	result := r.OK

	if tags, ok := d.GetOk("tags"); ok {
		err := assignOapiTags(tags.([]interface{}), result.InternetService.InternetServiceId, conn)
		if err != nil {
			return err
		}
	}

	d.SetId(result.InternetService.InternetServiceId)

	return resourceOutscaleOAPIInternetServiceRead(d, meta)
}

func resourceOutscaleOAPIInternetServiceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	id := d.Id()

	log.Printf("[DEBUG] Reading Internet Service id (%s)", id)

	req := &oapi.ReadInternetServicesRequest{
		Filters: oapi.FiltersInternetService{InternetServiceIds: []string{id}},
	}

	var resp *oapi.POST_ReadInternetServicesResponses
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.POST_ReadInternetServices(*req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var errString string

	if err != nil || resp.OK == nil {
		if err != nil {
			errString = err.Error()

		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
		}

		return fmt.Errorf("[DEBUG] Error reading Internet Service id (%s)", errString)
	}

	// Workaround to get the desired internet_service instance. TODO: Remove getInternetService
	// once filters work again. And use resp.OK.InternetServices[0]
	result, err := getInternetService(resp.OK.InternetServices, id)

	d.Set("request_id", resp.OK.ResponseContext.RequestId)
	d.Set("internet_service_id", result.InternetServiceId)

	if err := d.Set("net_id", result.NetId); err != nil {
		return err
	}

	if err := d.Set("state", result.State); err != nil {
		return err
	}

	return d.Set("tags", tagsOAPIToMap(result.Tags))
}

func resourceOutscaleOAPIInternetServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	d.Partial(true)

	if err := setOAPITags(conn, d); err != nil {
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

func resourceOutscaleOAPIInternetServiceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	id := d.Id()
	log.Printf("[DEBUG] Deleting Internet Service id (%s)", id)

	req := &oapi.DeleteInternetServiceRequest{
		InternetServiceId: id,
	}

	var err error
	var resp *oapi.POST_DeleteInternetServiceResponses

	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.POST_DeleteInternetService(*req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})

	var errString string

	if err != nil || resp.OK == nil {
		if err != nil {
			errString = err.Error()

		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
		}

		return fmt.Errorf("[DEBUG] Error deleting Internet Service id (%s)", errString)
	}

	return nil
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
