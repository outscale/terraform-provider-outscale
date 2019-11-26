/*
resource "outscale_volume" "outscale_volume_snap" {
    availability_zone   = format("%s%s", var.region, "a")
    size                = 40
    snapshot_id         = "snap-439943a0"
}
*/

resource "outscale_snapshot" "outscale_snapshot" {
    volume_id           = ""
}
