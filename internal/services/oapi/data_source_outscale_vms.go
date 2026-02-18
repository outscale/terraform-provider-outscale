package oapi

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscaleVMS() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleVMSRead,

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

func DataSourceOutscaleVMSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	vmID, vmIDOk := d.GetOk("vm_id")
	var err error
	if !filtersOk && !vmIDOk {
		return fmt.Errorf("one of filters, and vm id must be assigned")
	}

	// Build up search parameters
	params := osc.ReadVmsRequest{}
	if filtersOk {
		params.Filters, err = buildOutscaleDataSourceVMFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}
	if vmIDOk {
		params.Filters.VmIds = &[]string{vmID.(string)}
	}

	var resp osc.ReadVmsResponse
	err = retry.Retry(30*time.Second, func() *retry.RetryError {
		rp, httpResp, err := client.VmApi.ReadVms(ctx).ReadVmsRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading the vms %s", err)
	}

	// If no instances were returned, return
	if !resp.HasVms() {
		return ErrNoResults
	}

	var filteredVms []osc.Vm

	// loop through reservations, and remove terminated instances, populate vm slice
	for _, res := range resp.GetVms() {
		if res.GetState() != "terminated" {
			filteredVms = append(filteredVms, res)
		}
	}

	if len(filteredVms) == 0 {
		return ErrNoResults
	}

	d.SetId(id.UniqueId())
	return d.Set("vms", dataSourceOAPIVMS(filteredVms, client))
}

func dataSourceOAPIVMS(i []osc.Vm, client *osc.Client) []map[string]interface{} {
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
		mapsTags, _ := oapihelpers.GetBsuTagsMaps(v, client)
		vm["block_device_mappings_created"] = getOscAPIVMBlockDeviceMapping(mapsTags, v.GetBlockDeviceMappings())

		vm["tags"] = FlattenOAPITagsSDK(v.Tags)
		vms[index] = vm
	}
	return vms
}
