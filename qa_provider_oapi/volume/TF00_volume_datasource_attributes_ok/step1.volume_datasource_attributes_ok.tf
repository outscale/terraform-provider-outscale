resource "outscale_volume" "outscale_volume" {
    subregion_name = format("%s%s", var.region, "a")
    size           = 10
    iops           = 100
    volume_type    = "io1"
    tags {
    key ="name"
    value = "test-Volume"
         }
}

data "outscale_volume" "outscale_volume" {
    filter {
        name   = "volume_ids"
        values = [outscale_volume.outscale_volume.volume_id]
    }    
}
