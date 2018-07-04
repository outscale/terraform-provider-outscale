package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/icu"
)

func dataSourceOutscalePublicCatalog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscalePublicCatalogRead,
		Schema: map[string]*schema.Schema{
			"catalog": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attributes": {
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
						"entries": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attributes": {
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
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"title": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
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

func dataSourceOutscalePublicCatalogRead(d *schema.ResourceData, meta interface{}) error {
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
	catalog["attributes"] = flattenAttritbutes(getResp.Catalog.Attributes)
	catalog["entries"] = flattenEntries(getResp.Catalog.Entries)
	catList := make([]map[string]interface{}, 1)
	catList[0] = catalog

	if err := d.Set("catalog", catList); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())

	return d.Set("request_id", getResp.ResponseMetadata.RequestID)
}

// func flattenAttritbutes(attrs []*icu.CatalogAttribute) []map[string]interface{} {
// 	mapList := make([]map[string]interface{}, len(attrs))

// 	for k, v := range attrs {
// 		attrItem := make(map[string]interface{})
// 		attrItem["key"] = aws.StringValue(v.Key)
// 		attrItem["value"] = aws.StringValue(v.Value)
// 		mapList[k] = attrItem
// 	}
// 	return mapList
// }

// func flattenEntries(entries []*icu.CatalogEntry) []map[string]interface{} {
// 	mapList := make([]map[string]interface{}, len(entries))

// 	for k, v := range entries {
// 		attrItem := make(map[string]interface{})
// 		attrItem["attributes"] = flattenAttritbutes(v.Attributes)
// 		attrItem["value"] = int(aws.Int64Value(v.Value))
// 		attrItem["key"] = aws.StringValue(v.Key)
// 		attrItem["title"] = aws.StringValue(v.Title)
// 		mapList[k] = attrItem
// 	}
// 	return mapList
// }
