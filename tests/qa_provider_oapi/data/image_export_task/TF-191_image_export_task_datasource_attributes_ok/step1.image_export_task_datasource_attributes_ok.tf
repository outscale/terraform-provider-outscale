resource "outscale_vm" "vm-TF191" {
    image_id              = var.image_id
    vm_type               = var.vm_type
}

resource "outscale_image" "image-TF191" {
    image_name      = "terraform_export_task"
    vm_id           = outscale_vm.vm-TF191.vm_id
    no_reboot       = "true"
}

resource "outscale_image_export_task" "image_export_task-TF191" {
    image_id                     = outscale_image.image-TF191.image_id
    osu_export {
         disk_image_format       = "qcow2"
         osu_bucket              = var.osu_bucket_name
         osu_prefix              = "export-TF191"
         osu_api_key {
               api_key_id        = var.access_key_id
               secret_key        = var.secret_key_id
          }    
     }
} 

data "outscale_image_export_task" "outscale_image_export_task" {
   filter {
        name   = "task_ids"
        values = [outscale_image_export_task.image_export_task-TF191.task_id]
    }
}
