resource "outscale_volume" "outscale_volume" {
    subregion_name  = "${var.region}a"
    size            = 40
}

resource "outscale_snapshot" "outscale_snapshot" {
    volume_id = outscale_volume.outscale_volume.volume_id
}

resource "outscale_snapshot_attributes" "outscale_snapshot_attributes" {
    snapshot_id = outscale_snapshot.outscale_snapshot.snapshot_id
	
    permissions_to_create_volume_additions {
			account_ids = ["123456789012"]
	}
}

data "outscale_snapshot" "outscale_snapshot" {
    depends_on  = [outscale_snapshot_attributes.outscale_snapshot_attributes]
    snapshot_id = outscale_snapshot.outscale_snapshot.snapshot_id
}
