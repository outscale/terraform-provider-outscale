resource "outscale_flexible_gpu" "fGPU-2" {
   model_name             =  "nvidia-k2"
   subregion_name         =  "${var.region}a"
}

