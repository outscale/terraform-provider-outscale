package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceOutscaleOAPIServerCertificate() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOutscaleOAPIServerCertificateRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"expiration_date": {
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"upload_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceOutscaleOAPIServerCertificateRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.ReadServerCertificatesRequest{}
	if filters, filtersOk := d.GetOk("filter"); filtersOk {
		req.Filters = buildOutscaleOSCAPIDataSourceServerCertificateFilters(filters.(*schema.Set))
	}

	var resp oscgo.ReadServerCertificatesResponse
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		var err error
		rp, httpResp, err := conn.ServerCertificateApi.ReadServerCertificates(context.Background()).ReadServerCertificatesRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error reading Server Certificate id (%s)", utils.GetErrorResponse(err))
	}
	if err = utils.IsResponseEmptyOrMutiple(len(resp.GetServerCertificates()), "Server Certificate"); err != nil {
		return err
	}

	result := resp.GetServerCertificates()[0]

	d.Set("expiration_date", result.GetExpirationDate())
	d.Set("name", result.GetName())
	d.Set("orn", result.GetOrn())
	d.Set("path", result.GetPath())
	d.Set("upload_date", result.GetUploadDate())

	d.SetId(result.GetId())

	return nil
}

func buildOutscaleOSCAPIDataSourceServerCertificateFilters(set *schema.Set) *oscgo.FiltersServerCertificate {
	var filters oscgo.FiltersServerCertificate
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "paths":
			filters.SetPaths(filterValues)
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return &filters
}
