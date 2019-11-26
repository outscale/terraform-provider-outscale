resource "outscale_vm" "outscale_vm" {
    image_id           = var.image_id
    vm_type            = var.vm_type
    keypair_name       = var.keypair_name
    security_group_ids = [var.security_group_id]
}

resource "outscale_image" "outscale_image" {
    image_name = "terraform test image launch permission"
    #vm_id      = outscale_vm.outscale_vm.vm_id
    vm_id      = outscale_vm.outscale_vm.id         # .id test purpose only
    no_reboot  = "true"
}

resource "outscale_image_launch_permission" "outscale_image_launch_permission" {
    image_id = outscale_image.outscale_image.image_id

    permission_additions {
		account_ids = ["313532087625"]
		#account_id = "313532087625"                      # _id test purpose only
	}
}
