# Scenario: Succesfull creation of a blank volume
# Given a configuration file declaring a volume without snapshot_id
# When running terraform apply
# Then the volume is created. Can be seen in cockpit and attached to a vm. Seen as empty volume.

resource "outscale_keypair" "my_keypair" {
 keypair_name = "test-keypair-${random_string.suffix[0].result}"
}

resource "outscale_volume" "outscale_volume" {
    subregion_name = "${var.region}a"
    size           = 40
    tags {
    key = "type"
    value = "standard"
         }
}

# the instance is created at the same time to get the attributes of both resources prior to perform the link
resource "outscale_security_group" "security_group_TF160" {
  description         = "test-terraform-TF160"
  security_group_name = "test-sg-${random_string.suffix[0].result}"
}

resource "outscale_vm" "outscale_vm" {
    image_id           = var.image_id
    vm_type            = var.vm_type
    keypair_name       = outscale_keypair.my_keypair.keypair_name
    security_group_ids = [outscale_security_group.security_group_TF160.security_group_id]
}

resource "outscale_volume_link" "outscale_volume_link" {
    device_name = "/dev/xvdc"
    volume_id   = outscale_volume.outscale_volume.id
    vm_id       = outscale_vm.outscale_vm.id
}
