package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
)

func DataSourceOutscalePublicCatalog() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscalePublicCatalogRead,
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

func DataSourceOutscalePublicCatalogRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadPublicCatalogRequest{}

	resp, err := client.ReadPublicCatalog(ctx, req, options.WithRetryTimeout(20*time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	cs := resp.Catalog
	entries := ptr.From(cs.Entries)
	e_ret := make([]map[string]any, len(entries))

	for k, v := range entries {
		m := make(map[string]any)
		m["category"] = ptr.From(v.Category)
		if v.Flags != nil {
			m["flags"] = v.Flags
		}
		m["operation"] = ptr.From(v.Operation)
		m["service"] = ptr.From(v.Service)
		m["subregion_name"] = ptr.From(v.SubregionName)
		m["title"] = ptr.From(v.Title)
		m["type"] = ptr.From(v.Type)
		m["unit_price"] = ptr.From(v.UnitPrice)
		e_ret[k] = m
	}

	c_set := make(map[string]any, 1)
	c_set["entries"] = e_ret

	c_ret := make([]any, 1)
	c_ret[0] = c_set

	if err := d.Set("catalog", c_ret); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.UniqueId())

	return nil
}
