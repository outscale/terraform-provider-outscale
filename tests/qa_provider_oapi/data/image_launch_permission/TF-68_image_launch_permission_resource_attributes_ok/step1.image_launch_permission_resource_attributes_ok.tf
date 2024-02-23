resource "outscale_security_group" "my_sgImgl" {
   description = "test sg_group"
   security_group_name = "SG-inteImgl"
}

resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF68"
}

resource "outscale_vm" "outscale_vm" {
    image_id           = var.image_id
    vm_type            = var.vm_type
    keypair_name       = outscale_keypair.my_keypair.keypair_name
    security_group_ids = [outscale_security_group.my_sgImgl.security_group_id]
}

resource "outscale_image" "outscale_image" {
    image_name = "terraform-TF68"
    vm_id      = outscale_vm.outscale_vm.vm_id
    no_reboot  = "true"
}

resource "outscale_image_launch_permission" "outscale_image_launch_permission" {
    image_id = outscale_image.outscale_image.image_id

    permission_additions {
		account_ids = ["123456789012"]
	}
}
