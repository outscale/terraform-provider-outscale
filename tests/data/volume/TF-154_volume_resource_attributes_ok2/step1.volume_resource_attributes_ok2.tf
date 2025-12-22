resource "outscale_volume" "outscale_volume" {
    subregion_name = "${var.region}a"
    size            = 10
    iops            = 100
    volume_type     = "io1"
    tags {
      key = "name"
      value = "volume-io1"
    }
}
