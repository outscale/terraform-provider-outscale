resource "outscale_volume" "my_data" {
  subregion_name = "${var.region}a"
  volume_type    = var.volume_type
  iops           = var.volume_iops
  size           = var.volume_size_gib
}

resource "outscale_volumes_link" "my_volume_link" {
  device_name = "/dev/xvdb"
  volume_id   = outscale_volume.my_data.volume_id
  vm_id       = outscale_vm.my_vm.vm_id
}
