# Scenario: Succesfull creation of a blank volume
# Given a configuration file declaring a volume without snapshot_id
# When running terraform apply 
# Then the volume is created. Can be seen in cockpit and attached to a vm. Seen as empty volume.

resource "outscale_volume" "outscale_volume" {
    subregion_name = "${var.region}a"
    size           = 10
    tags {
     key = "name"
     value = "volume-standard"
    }

}




