resource "outscale_volume" "outscale_volume" {
    subregion_name = "${var.region}a"
    size = 15
    volume_type    = "standard"
    tags {
        key = "Name"
        value = "volume-standard-1"
    }
}
resource "outscale_volume" "outscale_volume2" {
    subregion_name = "${var.region}a"
    size = 13
    tags {
        key = "Name"
        value = "volume-standard-2"
    }
}
data "outscale_volumes" "outscale_volumes" {
    filter {
        name = "volume_ids"
        values = [outscale_volume.outscale_volume.volume_id, outscale_volume.outscale_volume2.volume_id]
    }
}
