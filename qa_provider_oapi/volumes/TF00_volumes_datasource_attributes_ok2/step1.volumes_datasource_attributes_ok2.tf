resource "outscale_volume" "outscale_volume" {
    subregion_name = format("%s%s", var.region, "b")
    size            = 12
    volume_type    = "gp2"
    tags {
        key = "Name"
        value = "volume-gp2-1"
    }

}
resource "outscale_volume" "outscale_volume2" {
    subregion_name = format("%s%s", var.region, "a")
    size            = 13
    volume_type    = "gp2"
    tags {
        key = "Name"
        value = "volume-gp2-1"
    }

}
data "outscale_volumes" "outscale_volumes" {
    filter {
        name   = "volume_ids"
        values = [outscale_volume.outscale_volume.volume_id, outscale_volume.outscale_volume2.volume_id]
    }
}
