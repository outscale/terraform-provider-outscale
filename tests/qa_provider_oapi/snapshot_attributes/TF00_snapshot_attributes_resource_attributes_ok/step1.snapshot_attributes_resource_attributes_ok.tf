resource "outscale_volume" "outscale_volume" {
    subregion_name = format("%s%s", var.region, "a")
    size            = 40
    #snapshot_id     = "snap-439943a0"
}

resource "outscale_snapshot" "outscale_snapshot" {
    volume_id = outscale_volume.outscale_volume.volume_id
}

resource "outscale_snapshot_attributes" "outscale_snapshot_attributes" {
    snapshot_id = outscale_snapshot.outscale_snapshot.snapshot_id
	
    permissions_to_create_volume_additions {
			account_ids = ["339215505907"]
	}
}

data "outscale_snapshot" "outscale_snapshot" {
    depends_on  = ["outscale_snapshot_attributes.outscale_snapshot_attributes"]
    snapshot_id = outscale_snapshot.outscale_snapshot.snapshot_id
}
