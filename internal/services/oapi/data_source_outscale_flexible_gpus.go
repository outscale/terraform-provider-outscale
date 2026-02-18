package oapi

import (
	"fmt"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleFlexibleGpus() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleFlexibleGpusRead,

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

func DataSourceOutscaleFlexibleGpusRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	_, IDOk := d.GetOk("flexible_gpu_id")

	if !filtersOk && !IDOk {
		return fmt.Errorf("one of filters, or flexible_gpu_id must be assigned")
	}

	var err error
	req := osc.ReadFlexibleGpusRequest{}
	if filtersOk {
		req.Filters, err = buildOutscaleDataSourceFlexibleGpuFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp osc.ReadFlexibleGpusResponse
	err = retry.Retry(30*time.Second, func() *retry.RetryError {
		rp, httpResp, err := client.FlexibleGpuApi.ReadFlexibleGpus(
			ctx).ReadFlexibleGpusRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		errString := err.Error()
		return fmt.Errorf("error reading flexible gpu (%s)", errString)
	}

	flexgps := resp.GetFlexibleGpus()[:]

	if len(flexgps) < 1 {
		return ErrNoResults
	}

	d.SetId(id.UniqueId())

	return setOAPIFlexibleGpuAttributes(d, flexgps)
}

func setOAPIFlexibleGpuAttributes(d *schema.ResourceData, fg []osc.FlexibleGpu) error {
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
