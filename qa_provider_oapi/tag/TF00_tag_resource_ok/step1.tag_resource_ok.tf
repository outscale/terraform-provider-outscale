resource "outscale_tag" "outscale_tag" {
    resource_ids = ["i-fc4bd3b6"]

    tags = {                               # NOK should be tag not tags
        name7 = "testDataSource7"          # NOK delete doesn't delete tag
        #name8 = "testDataSource8"          # tfa doesn't display correctly
    }                                      # tfs displays nothing
}



/*
data "outscale_vm" "vmd" {
    instance_id = "i-fc4bd3b6"    # centos created with cockpit
    #instance_id = "i-e4626d0a"    # windows created with cockpit
}

output "datasource_arch" {
    value = "${data.outscale_vm.vmd.instances_set.0.architecture}"
}

data "outscale_tag" "tag" {
   filter {
      name = "name7"
      values = ["testDataSource7"]      # windows created with cockpit
   }
   
   filter {
      name = "instance-id"
      values = ["i-fc4bd3b6"]
   }
}

output "tag" {
    value = "${data.outscale_tag.tag.resource_type}"
}
*/
