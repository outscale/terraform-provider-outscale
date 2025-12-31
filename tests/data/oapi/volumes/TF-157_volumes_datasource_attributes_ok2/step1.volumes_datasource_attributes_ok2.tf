resource "outscale_volume" "outscale_volume" {
  subregion_name = "${var.region}b"
  size           = 12
  volume_type    = "standard"
  tags {
    key   = "Name"
    value = "volume-standard-1"
  }

}
resource "outscale_volume" "outscale_volume2" {
  subregion_name = "${var.region}a"
  size           = 13
  volume_type    = "io1"
  iops           = 120
  tags {
    key   = "Name"
    value = "volume-io1-1"
  }

}
data "outscale_volumes" "outscale_volumes" {
  filter {
    name   = "volume_ids"
    values = [outscale_volume.outscale_volume.volume_id, outscale_volume.outscale_volume2.volume_id]
  }
}
