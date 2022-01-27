
resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF174"
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
