resource "outscale_volume" "outscale_volume" {
  subregion_name = "${var.region}a"
  size           = 10
  tags {
    key = "name"
    value = "volume-standard"
  }
}
