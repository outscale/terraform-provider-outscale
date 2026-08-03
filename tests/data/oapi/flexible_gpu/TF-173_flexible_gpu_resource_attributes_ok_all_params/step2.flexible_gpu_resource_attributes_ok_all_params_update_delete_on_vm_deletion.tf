resource "outscale_flexible_gpu" "fGPU-1" {
  model_name            = "nvidia-p6"
  generation            = "v5"
  delete_on_vm_deletion = false
}
