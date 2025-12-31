package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_DataSourceFlexibleGpus_basic(t *testing.T) {
	datasourcesName := "data.outscale_flexible_gpus.data_fGPU-1"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleFlexibleGpusConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourcesName, "flexible_gpus.0.model_name"),
					resource.TestCheckResourceAttrSet(datasourcesName, "flexible_gpus.0.generation"),
					resource.TestCheckResourceAttrSet(datasourcesName, "flexible_gpus.0.subregion_name"),
					resource.TestCheckResourceAttr(datasourcesName, "flexible_gpus.0.delete_on_vm_deletion", "true"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleFlexibleGpusConfig(region string) string {
	return fmt.Sprintf(`
                resource "outscale_flexible_gpu" "fGPUS-1" {
                        model_name             =  "nvidia-p6"
                        generation             =  "v5"
                        subregion_name         =  "%sa"
                        delete_on_vm_deletion  =   true
                }

		data "outscale_flexible_gpus" "data_fGPU-1" {
			filter {
				name = "flexible_gpu_ids"
				values = [outscale_flexible_gpu.fGPUS-1.flexible_gpu_id]
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
