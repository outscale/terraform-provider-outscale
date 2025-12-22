resource "outscale_keypair" "my_keypair" {
  keypair_name = "test-keypair-${random_string.suffix[0].result}"
}

resource "outscale_volume" "my_volume" {
  subregion_name = "${var.region}a"
  size           = 20
}

resource "outscale_security_group" "sg_snap" {
  description         = "test vms"
  security_group_name = "test-sg-${random_string.suffix[0].result}"
}

resource "outscale_snapshot" "my_snapshot" {
  volume_id = outscale_volume.my_volume.volume_id
}

## Test Public  VM with Block Device Mapping with multiple volumes  ##

resource "outscale_vm" "outscale_vm2" {
  image_id           = var.image_id
  vm_type            = var.vm_type
  keypair_name       = outscale_keypair.my_keypair.keypair_name
  security_group_ids = [outscale_security_group.sg_snap.security_group_id]
  block_device_mappings {
    device_name = "/dev/sda1" # resizing bootdisk volume
    bsu {
      volume_size           = "101"
      volume_type           = "standard"
      delete_on_vm_deletion = true
    }
  }
  block_device_mappings {
    device_name = "/dev/sdb"
    bsu {
      volume_size           = 16
      volume_type           = "io1"
      iops                  = 200
      delete_on_vm_deletion = true
    }
  }
  block_device_mappings {
    device_name = "/dev/sdc"
    bsu {
      volume_size           = 23
      volume_type           = "gp2"
      snapshot_id           = outscale_snapshot.my_snapshot.snapshot_id
      delete_on_vm_deletion = true
    }
  }
}

resource "outscale_volume" "volume01" {
  subregion_name = "${var.region}a"
  size           = 41
  volume_type    = "io1"
  iops           = 1000
}

resource "outscale_volume_link" "volume_link01" {
  device_name = "/dev/sdd"
  volume_id   = outscale_volume.volume01.id
  vm_id       = outscale_vm.outscale_vm2.id
}
