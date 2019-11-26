resource "outscale_volume" "outscale_volume_snap" {
    subregion_name = format("%s%s", var.region, "a")
    size           = 40
}

resource "outscale_snapshot" "outscale_snapshot" {
    volume_id = outscale_volume.outscale_volume_snap.volume_id
}
