resource "outscale_volume" "outscale_volume" {
    subregion_name = format("%s%s", var.region, "b")
    size            = 40
}
resource "outscale_volume" "outscale_volume2" {
    subregion_name = format("%s%s", var.region, "a")
    size            = 40
}
data "outscale_volumes" "outscale_volumes" {
    filter {
        name   = "volume_ids"
        values = [outscale_volume.outscale_volume.volume_id, outscale_volume.outscale_volume2.volume_id]
    }
}
