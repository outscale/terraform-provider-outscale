resource "outscale_vm" "outscale_vm" {
    count = 1

    image_id                = "ami-880caa66"
    vm_type                 = "c4.large"
    keypair_name            = "integ_sut_keypair"
    #firewall_rules_set_ids = ["sg-c73d3b6b"]
    firewall_rules_set_id = ["sg-c73d3b6b"]

    provisioner "local-exec" {
        command = "date; who -b"
    }
}

resource "outscale_image" "outscale_image" {
    name       = "image_${outscale_vm.outscale_vm.id}"
    vm_id      = "${outscale_vm.outscale_vm.id}"
    #no_reboot = "false"                 # default value
    no_reboot  = "true"                 # test value
}

resource "null_resource" "nullr" {
    count = 1

    depends_on = ["outscale_image.outscale_image"]

    provisioner "remote-exec" {
        inline = [
            "date; who -b",
            #""
        ]
        connection {
            type        = "ssh"
            user        = "centos"
            
            private_key = "${file("outscale_integ_sut_keypair.rsa.txt")}"
            #private_key = "${file("outscale_ct_test.rsa.txt")}"
            host        = "${outscale_vm.outscale_vm.instances_set.0.ip_address}"
        }
     }
}

