resource "outscale_vm" "outscale_vm" {
    #image_id               = "ami-880caa66"
    #image_id               = "ami-7f57f68f"      #dv1
    image_id               = "ami-5c450b62" #IN2
    vm_type                = "c4.large"
    #keypair_name           = "integ_sut_keypair"
    #keypair_name           = "terraform"
    keypair_name           = "testkp"
    #firewall_rules_set_ids = ["sg-c73d3b6b"]
    #firewall_rules_set_id  = "sg-c73d3b6b"    # tempo tests
    #security_group_ids     = ["sg-419f2c0c"]   # SV1
    security_group_ids     = ["sg-9752b7a6"]   # IN2

    tags = {                               # NOK should be tag not tags
        key   = "name7"
        value = "testDataSource7"          # NOK delete doesn't delete tag
        #name8 = "testDataSource8"          # tfa doesn't display correctly
    }                                      # tfs displays nothing
}

resource "outscale_vm" "outscale_vm2" {
    #image_id               = "ami-880caa66"
    #image_id               = "ami-7f57f68f"      #dv1
    image_id               = "ami-5c450b62" #IN2
    vm_type                = "c4.large"
    #keypair_name           = "integ_sut_keypair"
    #keypair_name           = "terraform"
    keypair_name           = "testkp"
    #firewall_rules_set_ids = ["sg-c73d3b6b"]
    #firewall_rules_set_id  = "sg-c73d3b6b"    # tempo tests
    #security_group_ids     = ["sg-419f2c0c"]   # SV1
    security_group_ids     = ["sg-9752b7a6"]   # IN2

    tags = {                               # NOK should be tag not tags
        key   = "name72"
        value = "testDataSource72"          # NOK delete doesn't delete tag
        #name8 = "testDataSource8"          # tfa doesn't display correctly
    }                                      # tfs displays nothing
}

data "outscale_tags" "outscale_tags" {
    filter {
        name = "resource_ids"
        values = [outscale_vm.outscale_vm.vm_id, outscale_vm.outscale_vm2.vm_id]
    }
}
