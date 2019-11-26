resource "outscale_volume" "outscale_volume_snap" {
    sub_region_name = format("%s%s", var.region, "a")
    size            = 40
    snapshot_id     = "snap-439943a0"
}

resource "outscale_snapshot" "outscale_snapshot" {
    volume_id = outscale_volume.outscale_volume_snap.volume_id
}

resource "outscale_snapshot_copy" "outscale_snapshot_copy" {
    description             = "backup snapshot"
    source_snapshot_id      = outscale_snapshot.outscale_snapshot.snapshot_id
    source_region_name      = format("%s%s", var.region, "a") 
    destination_region_name = format("%s%s", var.region, "a") 
} 
