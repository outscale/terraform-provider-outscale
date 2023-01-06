package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOutscaleOAPIFlexibleGpus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIFlexibleGpusRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"flexible_gpus": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_on_vm_deletion": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"model_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"generation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subregion_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flexible_gpu_id": {
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

func dataSourceOutscaleOAPIFlexibleGpusRead(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).OSCAPI
	filters, filtersOk := d.GetOk("filter")
	_, IDOk := d.GetOk("flexible_gpu_id")

	if !filtersOk && !IDOk {
		return fmt.Errorf("One of filters, or flexible_gpu_id must be assigned")
	}

	req := oscgo.ReadFlexibleGpusRequest{}
	req.SetFilters(buildOutscaleOAPIDataSourceFlexibleGpuFilters(filters.(*schema.Set)))

	var resp oscgo.ReadFlexibleGpusResponse
	var err error

	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.FlexibleGpuApi.ReadFlexibleGpus(
			context.Background()).ReadFlexibleGpusRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		errString := err.Error()
		return fmt.Errorf("[DEBUG] Error reading flexible gpu (%s)", errString)
	}

	flexgps := resp.GetFlexibleGpus()[:]

	if len(flexgps) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	d.SetId(resource.UniqueId())

	return setOAPIFlexibleGpuAttributes(d, flexgps)
}

func setOAPIFlexibleGpuAttributes(d *schema.ResourceData, fg []oscgo.FlexibleGpu) error {

	fgpus := make([]map[string]interface{}, len(fg))
	for k, v := range fg {
		fgpu := make(map[string]interface{})

		fgpu["delete_on_vm_deletion"] = v.GetDeleteOnVmDeletion()
		if v.GetFlexibleGpuId() != "" {
			fgpu["flexible_gpu_id"] = v.GetFlexibleGpuId()
		}
		if v.GetGeneration() != "" {
			fgpu["generation"] = v.GetGeneration()
		}
		if v.GetModelName() != "" {
			fgpu["model_name"] = v.GetModelName()
		}
		if v.GetState() != "" {
			fgpu["state"] = v.GetState()
		}
		if v.GetSubregionName() != "" {
			fgpu["subregion_name"] = v.GetSubregionName()
		}
		if v.GetVmId() != "" {
			fgpu["vm_id"] = v.GetVmId()
		}
		fgpus[k] = fgpu
	}

	return d.Set("flexible_gpus", fgpus)
}
