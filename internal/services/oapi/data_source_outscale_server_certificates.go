package oapi

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
)

func DataSourceOutscaleServerCertificates() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleServerCertificatesRead,
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

func DataSourceOutscaleServerCertificatesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")

	// Build up search parameters
	params := osc.ReadServerCertificatesRequest{}

	if filtersOk {
		filters, err := buildOutscaleOSCAPIDataSourceServerCertificateFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
		params.Filters = filters
	}

	resp, err := client.ReadServerCertificates(ctx, params, options.WithRetryTimeout(120*time.Second))
	var errString string
	if err != nil {
		errString = err.Error()
		return diag.Errorf("error reading server certificates (%s)", errString)
	}

	log.Printf("[DEBUG] Setting Server Certificates id (%s)", err)
	d.Set("server_certificates", flattenServerCertificates(ptr.From(resp.ServerCertificates)))
	d.SetId(id.UniqueId())
	return nil
}

func flattenServerCertificate(apiObject osc.ServerCertificate) map[string]interface{} {
	tfMap := map[string]interface{}{}
	tfMap["expiration_date"] = from.ISO8601(apiObject.ExpirationDate)
	tfMap["id"] = apiObject.Id
	tfMap["name"] = apiObject.Name
	tfMap["orn"] = apiObject.Orn
	tfMap["path"] = apiObject.Path
	tfMap["upload_date"] = from.ISO8601(apiObject.UploadDate)

	return tfMap
}

func flattenServerCertificates(apiObjects []osc.ServerCertificate) []map[string]interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []map[string]interface{}
	for _, apiObject := range apiObjects {
		tfList = append(tfList, flattenServerCertificate(apiObject))
	}
	return tfList
}
