resource "outscale_vm" "outscale_vm_centos" {
    count = 1                                             # plus testWebsite one already created

    image_id                = "ami-880caa66"
    type                    = "c4.large"
    keypair_name            = "integ_sut_keypair"
    #firewall_rules_set_ids = ["sg-c73d3b6b"]
    firewall_rules_set_id   = "sg-c73d3b6b"    # tempo test
}

data "outscale_vms" "outscale_vms" {
   filter {
      name = "image-id"
      values = ["ami-880caa66"]
   }
}
