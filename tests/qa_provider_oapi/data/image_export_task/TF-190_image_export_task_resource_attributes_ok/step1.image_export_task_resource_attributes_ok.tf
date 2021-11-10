resource "outscale_vm" "vm-TF190" {
    image_id              = var.image_id
    vm_type               = var.vm_type
}

resource "outscale_image" "image-TF190" {
    image_name      = "terraform_TF190"
    vm_id           = outscale_vm.vm-TF190.vm_id
    no_reboot       = "true"
}

resource "outscale_image_export_task" "image_export_task-TF190" {
    image_id                     = outscale_image.image-TF190.image_id
    osu_export {
         disk_image_format       = "qcow2"
         osu_bucket              = var.osu_bucket_name
         osu_prefix              = "export-TF190"
         osu_api_key {
               api_key_id        = var.access_key_id
               secret_key        = var.secret_key_id
          }    
     }
} 
