/*
resource "outscale_vm" "outscale_instance" {
    count = 1

    image_id                    = "ami-880caa66"
    #image_id                    = "ami-5ad76458"
    instance_type               = "c4.large"
    #instance_type               = "c4.large"
    disable_api_termination     = "true"
    key_name                    = "integ_sut_keypair"
    security_group              = ["sg-c73d3b6b"]
    

    provisioner "local-exec" {
        command = "date"
    }
    provisioner "local-exec" {
        command = "date; echo ${self.image_id} ${self.instance_type} ${self.id}"
    }
}

resource "outscale_image" "outscale_image" {
    name        = "terraform test"
    instance_id = "${outscale_vm.outscale_instance.instance_id}"
    no_reboot   = "true"
}

resource "outscale_image_launch_permission" "outscale_image_launch_permission" {
    image_id    = "${outscale_image.outscale_image.image_id}"
    permission {
        create {
            account_id = "339215505907"
        }
    }
}
*/