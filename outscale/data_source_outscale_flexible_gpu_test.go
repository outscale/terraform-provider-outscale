package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccOthers_DataSourceFlexibleGpu_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPIFlexibleGpuConfig(utils.GetRegion()),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIFlexibleGpuConfig(region string) string {
	return fmt.Sprintf(`
                resource "outscale_flexible_gpu" "dataGPU" {
                        model_name             =  "nvidia-p6"
                        generation             =  "v5"
                        subregion_name         =  "%sa"
                        delete_on_vm_deletion  =   true
                }

		data "outscale_flexible_gpu" "data_fGPU" {
			filter {
				name = "flexible_gpu_ids"
				values = [outscale_flexible_gpu.dataGPU.flexible_gpu_id]
			}
                        filter {
                                name = "delete_on_vm_deletion"
                                values = [true]
                        }
                        filter {
                                name = "generations"
                                values = ["v5"]
                        }
                        filter {
                                name = "states"
                                values = ["allocated"]
                        }
                        filter {
                                name = "model_names"
                                values = ["nvidia-p6"]
                        }
	                filter {
                                name = "subregion_names" 
                                values = ["%[1]sa"]
                        }
		}
	`, region)
}
