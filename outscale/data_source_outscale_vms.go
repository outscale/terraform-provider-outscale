package outscale

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceVMS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVMSRead,

		Schema: dataSourceVMSSchema(),
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

func dataSourceVMSSchema() map[string]*schema.Schema {
	wholeSchema := map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
		"vms": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: getVMAttributesSchema(),
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	return wholeSchema
}

func dataSourceVMSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	vmID, vmIDOk := d.GetOk("vm_id")

	if !filtersOk && !vmIDOk {
		return fmt.Errorf("One of filters, and vm ID must be assigned")
	}

	// Build up search parameters
	params := oscgo.ReadVmsRequest{}
	if filtersOk {
		params.Filters = buildDataSourceVMFilters(filters.(*schema.Set))
	}
	if vmIDOk {
		params.Filters.VmIds = &[]string{vmID.(string)}
	}

	var resp oscgo.ReadVmsResponse
	err := resource.Retry(30*time.Second, func() *resource.RetryError {
		rp, httpResp, err := client.VmApi.ReadVms(context.Background()).ReadVmsRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp.StatusCode, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("error reading the VMs %s", err)
	}

	// If no instances were returned, return
	if !resp.HasVms() {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	var filteredVms []oscgo.Vm

	// loop through reservations, and remove terminated instances, populate vm slice
	for _, res := range resp.GetVms() {
		if res.GetState() != "terminated" {
			filteredVms = append(filteredVms, res)
		}
	}

	if len(filteredVms) < 1 {
		return errors.New("Your query returned no results. Please change your search criteria and try again")
	}

	d.SetId(resource.UniqueId())
	return d.Set("vms", flattenVMS(filteredVms))
}

func flattenVMS(i []oscgo.Vm) []map[string]interface{} {
	vms := make([]map[string]interface{}, len(i))
	for index, v := range i {
		vm := make(map[string]interface{})

		setterFunc := func(key string, value interface{}) error {
			vm[key] = value
			return nil
		}

		if err := setVMAttributes(setterFunc, &v); err != nil {
			log.Fatalf("[DEBUG] setVMAttributes ERROR %+v", err)
		}

		vm["tags"] = getTagSet(v.GetTags())
		vms[index] = vm
	}
	return vms
}
