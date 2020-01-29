## Test Public  VM with Block Device Mapping with multiple volumes  ##

resource "outscale_vm" "outscale_vm" {
    image_id            = var.image_id
    vm_type             = var.vm_type
    keypair_name        = var.keypair_name
    block_device_mappings {
     device_name = "/dev/sdb"
     bsu = {
         volume_size=15
         volume_type = "gp2"
         snapshot_id = var.snapshot_id
       }
    }
    block_device_mappings {
     device_name = "/dev/sdc"
     bsu = {
         volume_size=22
         volume_type = "io1"
         iops      = 150
         snapshot_id = var.snapshot_id
         delete_on_vm_deletion = true
      }
    }
    tags {
    key = "name"
    value = "VM with multiple BDM"
    }
}