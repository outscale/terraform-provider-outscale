# Scenario: Succesfull creation of a blank volume
# Given a configuration file declaring a volume without snapshot_id
# When running terraform apply 
# Then the volume is created. Can be seen in cockpit and attached to a vm. Seen as empty volume.

resource "outscale_volume" "outscale_volume" {
    subregion_name = format("%s%s", var.region, "a")
    size            = 10
    volume_type     = "gp2"
    snapshot_id     = var.snapshot_id
    tags {
     key = "name"
     value = "test1"
    }
   
}
