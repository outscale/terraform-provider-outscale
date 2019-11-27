resource "outscale_vm" "outscale_instance" {
    image_id               = "ami-880caa66"
    #image_id              = "ami-5ad76458"
    type                   = "c4.large"
    #instance_type         = "c4.large"
    deletion_protection    = "true"
    keypair_name           = "integ_sut_keypair"
    firewall_rules_set_ids = ["sg-c73d3b6b"]
    

    provisioner "local-exec" {
        command = "date"
    }
    provisioner "local-exec" {
        command = "date; echo ${self.image_id} ${self.instance_type} ${self.id
    }
}

resource "outscale_image" "outscale_image" {
    name      = "terraform test"
    vm_id     = outscale_vm.outscale_instance.instance_id
    no_reboot = "true"
}

resource "outscale_image_launch_permission" "outscale_image_launch_permission" {
    image_id = outscale_image.outscale_image.image_id
    #permission {
    #    create {
    #        account_id = "339215505907"
    #    }
    #}
        permission_additions = [{
		account_id = "339215505907"
	}]
}
