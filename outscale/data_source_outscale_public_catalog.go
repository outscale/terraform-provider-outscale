package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/icu"
)

func dataSourceOutscaleOAPIPublicCatalog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIPublicCatalogRead,
		Schema: map[string]*schema.Schema{
			"catalog": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"catalog_attributes": {
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
						"catalog_entries": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"catalog_attributes": {
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
									"entry_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"short_description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"entry_value": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
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

func dataSourceOutscaleOAPIPublicCatalogRead(d *schema.ResourceData, meta interface{}) error {
	icuconn := meta.(*OutscaleClient).ICU

	request := &icu.ReadCatalogInput{}

	var getResp *icu.ReadCatalogOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		getResp, err = icuconn.API.ReadPublicCatalog(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading access key: %s", err)
	}

	// utils.PrintToJSON(getResp, "ReadCatalog")

	catalog := make(map[string]interface{})
	catalog["catalog_attributes"] = flattenOAPIAttritbutes(getResp.Catalog.Attributes)
	catalog["catalog_entries"] = flattenOAPIEntries(getResp.Catalog.Entries)
	catList := make([]map[string]interface{}, 1)
	catList[0] = catalog

	if err := d.Set("catalog", catList); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())

	return d.Set("request_id", getResp.ResponseMetadata.RequestID)
}
