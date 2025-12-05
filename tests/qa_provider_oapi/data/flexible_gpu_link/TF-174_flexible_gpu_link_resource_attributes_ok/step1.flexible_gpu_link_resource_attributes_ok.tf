resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF174"
}

resource "outscale_security_group" "my_sgfg_link" {
   description = "test security group"
   security_group_name = "SG-TF174"
}

resource "outscale_vm" "MaVM" {
   image_id                       = var.image_id
   vm_type                        = var.fgpu_vm_type
   keypair_name                   = outscale_keypair.my_keypair.keypair_name
   security_group_ids             = [outscale_security_group.my_sgfg_link.security_group_id]
   placement_subregion_name       = "${var.region}a"
   vm_initiated_shutdown_behavior = "restart"
}

resource "outscale_flexible_gpu" "fGPU-1" {
  model_name            = "nvidia-p6"
  generation            = var.fgpu_gen
  subregion_name        = "${var.region}a"
  delete_on_vm_deletion = true
}

resource "outscale_flexible_gpu_link" "link_fGPU" {
  flexible_gpu_ids = [outscale_flexible_gpu.fGPU-1.flexible_gpu_id]
  vm_id            = outscale_vm.MaVM.vm_id
}
