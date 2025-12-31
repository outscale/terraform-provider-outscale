resource "outscale_volume" "outscale_volume" {
    subregion_name = "${var.region}a"
    size            = 10
}

resource "outscale_snapshot" "outscale_snapshot" {
    volume_id = outscale_volume.outscale_volume.volume_id
    tags {
     key = "name"
     value = "test snapshot 1"
    }
tags {
      key = "Key"
      value = "value-tags"
     }
}

resource "outscale_snapshot" "outscale_snapshot2" {
    volume_id = outscale_volume.outscale_volume.volume_id
    tags {
     key = "name"
     value = "test snapshot 1"
    }
}

data "outscale_snapshots" "outscale_snapshots" {
    filter {
        name = "snapshot_ids"
        values = [outscale_snapshot.outscale_snapshot.snapshot_id,outscale_snapshot.outscale_snapshot2.snapshot_id]
    }    
}
