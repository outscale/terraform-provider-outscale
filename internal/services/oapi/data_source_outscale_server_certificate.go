package oapi

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleServerCertificate() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleServerCertificateRead,
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

func DataSourceOutscaleServerCertificateRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return ErrFilterRequired
	}

	// Build up search parameters
	params := oscgo.ReadServerCertificatesRequest{}

	if filtersOk {
		filterParams, err := buildOutscaleOSCAPIDataSourceServerCertificateFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
		params.Filters = filterParams
	}

	var resp oscgo.ReadServerCertificatesResponse
	err := retry.Retry(120*time.Second, func() *retry.RetryError {
		var err error
		rp, httpResp, err := conn.ServerCertificateApi.ReadServerCertificates(context.Background()).ReadServerCertificatesRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading server certificate id (%s)", utils.GetErrorResponse(err))
	}

	if !resp.HasServerCertificates() || len(resp.GetServerCertificates()) == 0 {
		return ErrNoResults
	}

	if len(resp.GetServerCertificates()) > 1 {
		return ErrMultipleResults
	}

	result := resp.GetServerCertificates()[0]

	log.Printf("[DEBUG] Setting Server Certificate id (%s)", err)

	d.Set("expiration_date", result.GetExpirationDate())
	d.Set("name", result.GetName())
	d.Set("orn", result.GetOrn())
	d.Set("path", result.GetPath())
	d.Set("upload_date", result.GetUploadDate())

	d.SetId(result.GetId())

	return nil
}

func buildOutscaleOSCAPIDataSourceServerCertificateFilters(set *schema.Set) (*oscgo.FiltersServerCertificate, error) {
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
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}
