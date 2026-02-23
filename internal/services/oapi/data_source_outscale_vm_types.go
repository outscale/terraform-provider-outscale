package oapi

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"
)

func DataSourceOutscaleVMTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleVMTypesRead,

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

func DataSourceOutscaleVMTypesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filter, filterOk := d.GetOk("filter")

	var req osc.ReadVmTypesRequest
	var err error
	if filterOk {
		req.Filters, err = buildOutscaleDataSourceVMTypesFilters(filter.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadVmTypes(ctx, req, options.WithRetryTimeout(30*time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	filteredTypes := ptr.From(resp.VmTypes)[:]

	if len(filteredTypes) < 1 {
		return diag.FromErr(ErrNoResults)
	}

	return diag.FromErr(statusDescriptionOAPIVMTypesAttributes(d, filteredTypes))
}

func setOAPIVMTypeAttributes(set AttributeSetter, vType *osc.VmType) error {
	if err := set("bsu_optimized", vType.BsuOptimized); err != nil {
		return err
	}
	if err := set("max_private_ips", vType.MaxPrivateIps); err != nil {
		return err
	}
	if err := set("memory_size", vType.MemorySize); err != nil {
		return err
	}
	if err := set("vcore_count", vType.VcoreCount); err != nil {
		return err
	}
	if err := set("vm_type_name", vType.VmTypeName); err != nil {
		return err
	}
	if err := set("volume_count", vType.VolumeCount); err != nil {
		return err
	}
	if err := set("volume_size", vType.VolumeSize); err != nil {
		return err
	}

	return nil
}

func statusDescriptionOAPIVMTypesAttributes(d *schema.ResourceData, fTypes []osc.VmType) error {
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

func buildOutscaleDataSourceVMTypesFilters(set *schema.Set) (*osc.FiltersVmType, error) {
	var filters osc.FiltersVmType
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "bsu_optimized":
			filters.BsuOptimized = new(cast.ToBool(filterValues[0]))
		case "ephemerals_types":
			filters.EphemeralsTypes = &filterValues
		case "eths":
			filters.Eths = new(utils.StringSliceToIntSlice(filterValues))
		case "gpus":
			filters.Gpus = new(utils.StringSliceToIntSlice(filterValues))
		case "memory_sizes":
			filters.MemorySizes = new(utils.StringSliceToFloat32Slice(filterValues))
		case "vcore_counts":
			filters.VcoreCounts = new(utils.StringSliceToIntSlice(filterValues))
		case "vm_type_names":
			filters.VmTypeNames = &filterValues
		case "volume_counts":
			filters.VolumeCounts = new(utils.StringSliceToIntSlice(filterValues))
		case "volume_sizes":
			filters.VolumeSizes = new(utils.StringSliceToIntSlice(filterValues))
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
