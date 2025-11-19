resource "outscale_volume" "outscale_volume" {
    subregion_name = "${var.region}a"
    size            = 12
    volume_type     = "standard"
}
