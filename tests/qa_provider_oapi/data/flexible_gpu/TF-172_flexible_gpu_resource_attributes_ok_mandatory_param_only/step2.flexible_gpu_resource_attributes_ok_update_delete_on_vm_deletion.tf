resource "outscale_flexible_gpu" "fGPU-2" {
   model_name             =  "nvidia-p6"
   subregion_name         =  "${var.region}a"
   delete_on_vm_deletion  =   true
}

