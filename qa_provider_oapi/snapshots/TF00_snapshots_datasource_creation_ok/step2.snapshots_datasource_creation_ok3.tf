/*
resource "outscale_volume" "outscale_volume" {
    count = 1

    sub_region_name = format("%s%s", var.region, "a")
    size = 40
}

resource "outscale_snapshot" "outscale_snapshot" {
    count = 1

    volume_id = outscale_volume.outscale_volume.volume_id
}

resource "outscale_snapshot" "outscale_snapshot2" {
    count = 1

    volume_id = outscale_volume.outscale_volume.volume_id
}
*/

data "outscale_snapshot" "outscale_snapshot" {
    filter {
        volume_id = "vol-fc4726af"
    }
}
