package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func DataSourceOutscaleServerCertificates() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleServerCertificatesRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"server_certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expiration_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"orn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"upload_date": {
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

func DataSourceOutscaleServerCertificatesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")

	// Build up search parameters
	params := oscgo.ReadServerCertificatesRequest{}

	if filtersOk {
		filters, err := buildOutscaleOSCAPIDataSourceServerCertificateFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
		params.Filters = filters
	}

	var resp oscgo.ReadServerCertificatesResponse
	var err error
	err = retry.Retry(120*time.Second, func() *retry.RetryError {
		rp, httpResp, err := conn.ServerCertificateApi.ReadServerCertificates(context.Background()).ReadServerCertificatesRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	var errString string
	if err != nil {
		errString = err.Error()
		return fmt.Errorf("[DEBUG] Error reading Server Certificates (%s)", errString)
	}

	log.Printf("[DEBUG] Setting Server Certificates id (%s)", err)
	d.Set("server_certificates", flattenServerCertificates(resp.GetServerCertificates()))
	d.SetId(id.UniqueId())
	return nil
}

func flattenServerCertificate(apiObject oscgo.ServerCertificate) map[string]interface{} {
	tfMap := map[string]interface{}{}
	tfMap["expiration_date"] = apiObject.GetExpirationDate()
	tfMap["id"] = apiObject.GetId()
	tfMap["name"] = apiObject.GetName()
	tfMap["orn"] = apiObject.GetOrn()
	tfMap["path"] = apiObject.GetPath()
	tfMap["upload_date"] = apiObject.GetUploadDate()

	return tfMap
}

func flattenServerCertificates(apiObjects []oscgo.ServerCertificate) []map[string]interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []map[string]interface{}
	for _, apiObject := range apiObjects {
		tfList = append(tfList, flattenServerCertificate(apiObject))
	}
	return tfList
}
