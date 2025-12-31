
resource "outscale_keypair" "my_keypair" {
  keypair_name = "test-keypair-${random_string.suffix[0].result}"
}

resource "outscale_security_group" "my_sgfg" {
  description         = "test security group"
  security_group_name = "test-sg-${random_string.suffix[0].result}"
}

resource "outscale_vm" "MaVM" {
  image_id                       = var.image_id
  vm_type                        = var.fgpu_vm_type
  keypair_name                   = outscale_keypair.my_keypair.keypair_name
  security_group_ids             = [outscale_security_group.my_sgfg.security_group_id]
  placement_subregion_name       = "${var.region}a"
  vm_initiated_shutdown_behavior = "restart"
}

resource "outscale_flexible_gpu" "fGPU-1" {
  model_name            = "nvidia-p6"
  generation            = var.fgpu_gen
  subregion_name        = "${var.region}a"
  delete_on_vm_deletion = false
}


resource "outscale_flexible_gpu_link" "link_fGPU" {
  flexible_gpu_ids = [outscale_flexible_gpu.fGPU-1.flexible_gpu_id]
  vm_id            = outscale_vm.MaVM.vm_id
}

data "outscale_flexible_gpu" "data-fGPU-1" {
  filter {
    name   = "flexible_gpu_ids"
    values = [outscale_flexible_gpu.fGPU-1.flexible_gpu_id]
  }
  depends_on = [outscale_flexible_gpu_link.link_fGPU]
}

data "outscale_flexible_gpu" "data-fGPU-2" {

  filter {
    name   = "delete_on_vm_deletion"
    values = [false]
  }
  filter {
    name   = "generations"
    values = [var.fgpu_gen]
  }
  filter {
    name   = "states"
    values = ["attached"]
  }
  filter {
    name   = "model_names"
    values = ["nvidia-p6"]
  }
  filter {
    name   = "subregion_names"
    values = ["${var.region}a"]
  }

  depends_on = [outscale_flexible_gpu_link.link_fGPU]
}

data "outscale_flexible_gpu" "data-fGPU-3" {
  filter {
    name   = "vm_ids"
    values = [outscale_vm.MaVM.vm_id]
  }
  depends_on = [outscale_flexible_gpu_link.link_fGPU]
}
