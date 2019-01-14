package outscale

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleOAPIInternetServiceLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIInternetServiceLinkCreate,
		Read:   resourceOutscaleOAPIInternetServiceLinkRead,
		Delete: resourceOutscaleOAPIInternetServiceLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: getOAPIInternetServiceLinkSchema(),
	}
}

func resourceOutscaleOAPIInternetServiceLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	netId := d.Get("net_id").(string)
	igID := d.Get("internet_service_id").(string)

	req := &oapi.LinkInternetServiceRequest{
		NetId:             netId,
		InternetServiceId: igID,
	}

	var err error
	var resp *oapi.POST_LinkInternetServiceResponses
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.POST_LinkInternetService(*req)

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

		return fmt.Errorf("[DEBUG] Error linking internet service id (%s)", errString)
	}

	d.SetId(igID)

	return resourceOutscaleOAPIInternetServiceLinkRead(d, meta)
}

func resourceOutscaleOAPIInternetServiceLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	id := d.Id()

	d.SetId(id)

	log.Printf("Reading Internet Service id (%s)", id)

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
	err, result := getInternetService(resp.OK.InternetServices, id)

	if resp == nil {
		d.SetId("")
		return errors.New("Got a nil response for internet service Link")
	}

	if err := d.Set("tags", tagsOAPIToMap(result.Tags)); err != nil {
		return err
	}

	d.Set("state", result.State)
	d.Set("internet_service_id", result.InternetServiceId)
	d.Set("request_id", resp.OK.ResponseContext.RequestId)

	return nil
}

func resourceOutscaleOAPIInternetServiceLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	netId := d.Get("net_id").(string)
	igID := d.Get("internet_service_id").(string)

	req := &oapi.UnlinkInternetServiceRequest{
		NetId:             netId,
		InternetServiceId: igID,
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.POST_UnlinkInternetService(*req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("Error dettaching internet service id (%s)", err)
	}

	d.SetId("")

	return nil
}

func getOAPIInternetServiceLinkSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
	}
}
