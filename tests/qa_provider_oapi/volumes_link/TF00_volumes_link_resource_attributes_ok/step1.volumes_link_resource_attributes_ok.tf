# Scenario: Succesfull creation of a blank volume
# Given a configuration file declaring a volume without snapshot_id
# When running terraform apply 
# Then the volume is created. Can be seen in cockpit and attached to a vm. Seen as empty volume.

resource "outscale_volume" "outscale_volume" {
    subregion_name = format("%s%s", var.region, "a")
    size           = 40
    tags {
    key = "type"
    value = "standard"
         }
}

# the instance is created at the same time to get the attributes of both resources prior to perform the link

resource "outscale_vm" "outscale_vm" {
    image_id           = var.image_id
    vm_type            = var.vm_type
    keypair_name       = var.keypair_name
    security_group_ids = [var.security_group_id]
}

resource "outscale_volumes_link" "outscale_volumes_link" {
    #device_name = "/dev/sdc"
    device_name = "/dev/xvdc"
    volume_id   = outscale_volume.outscale_volume.id
    vm_id       = outscale_vm.outscale_vm.id
    #vm_id       = "i-ac095195"
}
