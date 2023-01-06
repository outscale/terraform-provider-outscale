package outscale

import (
	"context"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleOAPIPublicCatalog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIPublicCatalogRead,
		Schema: map[string]*schema.Schema{
			"catalog": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entries": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"flags": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"operation": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"service": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subregion_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"title": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"unit_price": {
										Type:     schema.TypeFloat,
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
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadPublicCatalogRequest{}

	var resp oscgo.ReadPublicCatalogResponse
	var err error
	err = resource.Retry(20*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.PublicCatalogApi.ReadPublicCatalog(context.Background()).ReadPublicCatalogRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	cs := resp.GetCatalog()
	entries := cs.GetEntries()[:]
	e_ret := make([]map[string]interface{}, len(entries))

	for k, v := range entries {
		m := make(map[string]interface{})
		m["category"] = v.GetCategory()
		if v.HasFlags() {
			m["flags"] = v.GetFlags()
		}
		m["operation"] = v.GetOperation()
		m["service"] = v.GetService()
		m["subregion_name"] = v.GetSubregionName()
		m["title"] = v.GetTitle()
		m["type"] = v.GetType()
		m["unit_price"] = v.GetUnitPrice()
		e_ret[k] = m
	}

	c_set := make(map[string]interface{}, 1)
	c_set["entries"] = e_ret

	c_ret := make([]interface{}, 1)
	c_ret[0] = c_set

	if err := d.Set("catalog", c_ret); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())

	return nil
}
