package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_FlexibleGpu_DataSource(t *testing.T) {
	t.Parallel()
	region := fmt.Sprintf("%sa", os.Getenv("OUTSCALE_REGION"))
	dataSourceName := "data.outscale_flexible_gpu.fgpu"
	dataSourcesName := "data.outscale_flexible_gpus.fgpus"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_FlexibleGpu_DataSource_Config(region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "flexible_gpu_id"),
					resource.TestCheckResourceAttrSet(dataSourcesName, "flexible_gpus.#"),
				),
			},
		},
	})
}

func testAcc_FlexibleGpu_DataSource_Config(region string) string {
	return fmt.Sprintf(`
	resource "outscale_flexible_gpu" "fGPU-1" {
		model_name             =  "nvidia-p6"
		generation             =  "v5"
		subregion_name         =  "%[1]s"
		delete_on_vm_deletion  =   true
	}

	data "outscale_flexible_gpu" "fgpu" {
		filter {
			name = "flexible_gpu_ids"
			values = [outscale_flexible_gpu.fGPU-1.flexible_gpu_id]
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
			values = ["%[1]s"]
		}
	}

	data "outscale_flexible_gpus" "fgpus" {
		filter {
			name = "flexible_gpu_ids"
			values = [outscale_flexible_gpu.fGPU-1.flexible_gpu_id]
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
			values = ["%[1]s"]
		}
	}
	`, region)
}
