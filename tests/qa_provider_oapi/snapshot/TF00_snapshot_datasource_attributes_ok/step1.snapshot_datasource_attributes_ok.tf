resource "outscale_volume" "outscale_volume" {
    subregion_name = format("%s%s", var.region, "a")
    size           = 40
}

resource "outscale_snapshot" "outscale_snapshot" {
    volume_id = outscale_volume.outscale_volume.volume_id
    tags {
     key = "name"
     value = "Snapshot_datasource"
    }

}

data "outscale_snapshot" "outscale_snapshot" {
    filter {
        name   = "snapshot_ids"
        values = [outscale_snapshot.outscale_snapshot.snapshot_id]
    }
}
