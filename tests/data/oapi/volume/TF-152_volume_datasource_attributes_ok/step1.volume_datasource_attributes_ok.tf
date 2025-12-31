resource "outscale_volume" "outscale_volume" {
    subregion_name = "${var.region}a"
    size           = 14
    iops           = 140
    volume_type    = "io1"
   tags {
    key ="name"
    value = "test-Volume"
         }
   tags {
      key = "Key"
      value = "value-tags"
     }
}

data "outscale_volume" "outscale_volume" {
    filter {
        name   = "volume_ids"
        values = [outscale_volume.outscale_volume.volume_id]
    }    
}
