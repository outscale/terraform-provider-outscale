
resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF177"
}

resource "outscale_vm" "MaVM" {
   image_id                       = var.image_id
   vm_type                        = var.vm_type
   keypair_name                   = outscale_keypair.my_keypair.keypair_name
   placement_subregion_name       = "${var.region}a"
   vm_initiated_shutdown_behavior = "restart"
}

resource "outscale_flexible_gpu" "fGPU-1" {
   model_name                   =  "nvidia-k2"
   generation                   =  "v4"
   subregion_name               =  "${var.region}a"
   delete_on_vm_deletion        =   true
}


resource "outscale_flexible_gpu_link" "link_fGPU" {
   flexible_gpu_id              =  outscale_flexible_gpu.fGPU-1.flexible_gpu_id
    vm_id                        = outscale_vm.MaVM.vm_id
}

data "outscale_flexible_gpu" "data-fGPU-1" {
filter {
        name     = "flexible_gpu_ids"
        values   = [outscale_flexible_gpu.fGPU-1.flexible_gpu_id]
    }
depends_on =[outscale_flexible_gpu_link.link_fGPU]
}

data "outscale_flexible_gpu" "data-fGPU-2" {

filter {
        name     = "delete_on_vm_deletion"
        values   = [true]
    }
  filter {
        name     = "generations"
        values   = [ "v4"]
    }
  filter {
        name     = "states"
        values   = ["attached"]
    }
  filter {
        name     = "model_names"
        values   = ["nvidia-k2"]
    }
  filter {
        name     = "subregion_names"
        values   = ["${var.region}a"] 
    }

depends_on =[outscale_flexible_gpu_link.link_fGPU]
}

data "outscale_flexible_gpu" "data-fGPU-3" {
filter {
        name     = "vm_ids"
        values   = [outscale_vm.MaVM.vm_id]
    }
depends_on =[outscale_flexible_gpu_link.link_fGPU]
}
