# Scenario: Succesfull creation of a blank volume
# Given a configuration file declaring a volume without snapshot_id
# When running terraform apply 
# Then the volume is created. Can be seen in cockpit and attached to a vm. Seen as empty volume.

resource "outscale_volume" "volume-1" {
    subregion_name  = "${var.region}a"
    size            = 10
    tags {
       key          = "name"
       value        = "test1"
    }
}

resource "outscale_snapshot" "snapshot-1" {
    volume_id   = outscale_volume.volume-1.volume_id
    tags {
        key     = "name"
        value   = "Snapsho_Creation_test"
    }
}


resource "outscale_volume" "volume-2" {
    subregion_name  = "${var.region}a"
    size            = 25
    volume_type     = "standard"
    snapshot_id     = outscale_snapshot.snapshot-1.snapshot_id
}
