resource "outscale_vm" "outscale_vm" {
  count = 1

  image_id = "ami-880caa66"

  instance_type = "c4.large"

  #key_name = "integ_sut_keypair"


  #security_group = ["sg-c73d3b6b"]

  disable_api_termination = false
  ebs_optimized           = true
}

resource "outscale_vm_attributes" "outscale_vm_attributes" {
  instance_id = "${outscale_vm.outscale_vm.0.id}"

  # This will not make the change in outscale_vm tf state because 
  # disable api termination is an argument for this resource, but it does
  # the change and we can check it on cockpit

  attribute               = "disableApiTermination"
  disable_api_termination = false

  # Changes are made succesfully, but this needs to make a tf refresh to
  # see the changes in tf state


  # attribute     = "instanceType"
  # instance_type = "c4.large"


  # Changes are made succesfully, and show its as well in tf state

  attribute     = "ebsOptimized"
  ebs_optimized = true

  # This doesent appear in terraform state as outscale_vm_attributes
  # But it does the change, to see the changes we need to do an tf refresh

  # attribute = "blockDeviceMapping"
  # block_device_mapping {
  #   device_name = "/dev/sda1"

  #   ebs {
  #     delete_on_termination = false
  #   }
  # }
}
