resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF146"
}

resource "outscale_volume" "my_volume" {
    subregion_name = "${var.region}a"
    size           = 20
}

resource "outscale_snapshot" "my_snapshot" {
    volume_id = outscale_volume.my_volume.volume_id
}

## Test Public  VM with Block Device Mapping with multiple volumes  ##

resource "outscale_vm" "outscale_vm2" {
    image_id            = var.image_id
    vm_type             = var.vm_type
    keypair_name        = outscale_keypair.my_keypair.keypair_name
    block_device_mappings {
      device_name = "/dev/sda1"   # resizing bootdisk volume
      bsu {
      volume_size = "100"
      volume_type = "standard"
      delete_on_vm_deletion = false
      }
    }
     block_device_mappings {
     device_name = "/dev/sdb"
     bsu  {
         volume_size=15
         volume_type = "standard"
       }
    }
    block_device_mappings {
     device_name = "/dev/sdc"
     bsu  {
         volume_size=22
         volume_type = "io1"
         iops      = 150
         snapshot_id = outscale_snapshot.my_snapshot.snapshot_id
         delete_on_vm_deletion = true
      }
    }
}
