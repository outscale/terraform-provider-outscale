package oapi

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
)

func DataSourceOutscaleVMS() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleVMSRead,

		Schema: DataSourceOutscaleVMSSchema(),
	}
}

func dataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},

				"values": {
					Type:     schema.TypeList,
					Required: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func DataSourceOutscaleVMSSchema() map[string]*schema.Schema {
	wholeSchema := map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
		"vms": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: getOApiVMAttributesSchema(),
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	return wholeSchema
}

func DataSourceOutscaleVMSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutRead)

	filters, filtersOk := d.GetOk("filter")
	vmID, vmIDOk := d.GetOk("vm_id")
	var err error
	if !filtersOk && !vmIDOk {
		return diag.Errorf("one of filters, and vm id must be assigned")
	}

	// Build up search parameters
	params := osc.ReadVmsRequest{}
	if filtersOk {
		params.Filters, err = buildOutscaleDataSourceVMFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if vmIDOk {
		params.Filters.VmIds = &[]string{vmID.(string)}
	}

	resp, err := client.ReadVms(ctx, params, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error reading the vms %s", err)
	}

	// If no instances were returned, return
	if resp.Vms == nil {
		return diag.FromErr(ErrNoResults)
	}

	var filteredVms []osc.Vm

	// loop through reservations, and remove terminated instances, populate vm slice
	for _, res := range *resp.Vms {
		if res.State != "terminated" {
			filteredVms = append(filteredVms, res)
		}
	}

	if len(filteredVms) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	d.SetId(id.UniqueId())

	return diag.FromErr(d.Set("vms", dataSourceOAPIVMS(ctx, client, timeout, filteredVms)))
}

func dataSourceOAPIVMS(ctx context.Context, client *osc.Client, timeout time.Duration, i []osc.Vm) []map[string]interface{} {
	vms := make([]map[string]interface{}, len(i))
	for index, v := range i {
		vm := make(map[string]interface{})

		setterFunc := func(key string, value interface{}) error {
			vm[key] = value
			return nil
		}

		if err := oapiVMDescriptionAttributes(setterFunc, &v); err != nil {
			log.Fatalf("[DEBUG] oapiVMDescriptionAttributes ERROR %+v", err)
		}
		mapsTags, _ := oapihelpers.GetBsuTagsMaps(ctx, client, timeout, v)
		vm["block_device_mappings_created"] = getOscAPIVMBlockDeviceMapping(mapsTags, v.BlockDeviceMappings)

		vm["tags"] = FlattenOAPITagsSDK(v.Tags)
		vms[index] = vm
	}
	return vms
}
