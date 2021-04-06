resource "outscale_volume" "my_data" {
  subregion_name = "${var.region}a"
  volume_type    = var.volume_type
  iops           = var.volume_iops
  size           = var.volume_size_gib
}
