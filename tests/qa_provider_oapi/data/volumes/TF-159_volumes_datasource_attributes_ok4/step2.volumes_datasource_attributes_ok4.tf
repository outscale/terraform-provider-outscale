resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF159"
}
resource "outscale_volume" "outscale_volume3" {
    subregion_name = "${var.region}a"
    size           = 40
    iops           = 100
    volume_type    = "io1"
    tags {
    key = "type"
    value = "io1"
         }
}

# the instance is created at the same time to get the attributes of both resources prior to perform the link

resource "outscale_vm" "outscale_vm" {
    image_id           = var.image_id
    vm_type            = var.vm_type
    keypair_name       = outscale_keypair.my_keypair.keypair_name
}

resource "outscale_volumes_link" "outscale_volumes_link" {
    device_name = "/dev/xvdc"
    volume_id   = outscale_volume.outscale_volume3.id
    vm_id       = outscale_vm.outscale_vm.id
} 



resource "outscale_volume" "outscale_volume" {
    subregion_name = "${var.region}a"
    size = 15
    volume_type    = "standard"
    tags {
        key = "Name"
        value = "volume-standard-1"
    }
}
resource "outscale_volume" "outscale_volume2" {
    subregion_name = "${var.region}a"
    size = 13
    tags {
        key = "Name"
        value = "volume-standard-2"
    }
}
data "outscale_volumes" "outscale_volumes" {
    filter {
        name = "volume_ids"
        values = [outscale_volume.outscale_volume.volume_id, outscale_volume.outscale_volume2.volume_id, outscale_volume.outscale_volume3.volume_id]
    }
}
