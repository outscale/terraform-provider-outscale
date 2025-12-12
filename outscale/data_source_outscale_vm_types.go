package outscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/spf13/cast"
)

func DataSourceOutscaleVMTypes() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleVMTypesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vm_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bsu_optimized": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"max_private_ips": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vcore_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vm_type_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_size": {
							Type:     schema.TypeInt,
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

func DataSourceOutscaleVMTypesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filter, filterOk := d.GetOk("filter")

	var req oscgo.ReadVmTypesRequest
	var err error
	if filterOk {
		req.Filters, err = buildOutscaleDataSourceVMTypesFilters(filter.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp oscgo.ReadVmTypesResponse
	err = retry.Retry(30*time.Second, func() *retry.RetryError {
		var err error
		rp, httpResp, err := conn.VmApi.ReadVmTypes(context.Background()).ReadVmTypesRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	filteredTypes := resp.GetVmTypes()[:]

	if len(filteredTypes) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	return statusDescriptionOAPIVMTypesAttributes(d, filteredTypes)

}

func setOAPIVMTypeAttributes(set AttributeSetter, vType *oscgo.VmType) error {

	if err := set("bsu_optimized", vType.GetBsuOptimized()); err != nil {
		return err
	}
	if err := set("max_private_ips", vType.GetMaxPrivateIps()); err != nil {
		return err
	}
	if err := set("memory_size", vType.GetMemorySize()); err != nil {
		return err
	}
	if err := set("vcore_count", vType.GetVcoreCount()); err != nil {
		return err
	}
	if err := set("vm_type_name", vType.GetVmTypeName()); err != nil {
		return err
	}
	if err := set("volume_count", vType.GetVolumeCount()); err != nil {
		return err
	}
	if err := set("volume_size", vType.GetVolumeSize()); err != nil {
		return err
	}

	return nil
}

func statusDescriptionOAPIVMTypesAttributes(d *schema.ResourceData, fTypes []oscgo.VmType) error {
	d.SetId(id.UniqueId())

	vTypes := make([]map[string]interface{}, len(fTypes))

	for k, v := range fTypes {
		vtype := make(map[string]interface{})

		setterFunc := func(key string, value interface{}) error {
			vtype[key] = value
			return nil
		}

		if err := setOAPIVMTypeAttributes(setterFunc, &v); err != nil {
			return err
		}

		vTypes[k] = vtype
	}

	return d.Set("vm_types", vTypes)

}

func buildOutscaleDataSourceVMTypesFilters(set *schema.Set) (*oscgo.FiltersVmType, error) {
	var filters oscgo.FiltersVmType
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "bsu_optimized":
			filters.SetBsuOptimized(cast.ToBool(filterValues[0]))
		case "ephemerals_types":
			filters.SetEphemeralsTypes(filterValues)
		case "eths":
			filters.SetEths(utils.StringSliceToInt32Slice(filterValues))
		case "gpus":
			filters.SetGpus(utils.StringSliceToInt32Slice(filterValues))
		case "memory_sizes":
			filters.SetMemorySizes(utils.StringSliceToFloat32Slice(filterValues))
		case "vcore_counts":
			filters.SetVcoreCounts(utils.StringSliceToInt32Slice(filterValues))
		case "vm_type_names":
			filters.SetVmTypeNames(filterValues)
		case "volume_counts":
			filters.SetVolumeCounts(utils.StringSliceToInt32Slice(filterValues))
		case "volume_sizes":
			filters.SetVolumeSizes(utils.StringSliceToInt32Slice(filterValues))
		default:
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}
