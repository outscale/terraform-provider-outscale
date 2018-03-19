resource "outscale_vm" "outscale_vm" {
  count = 1

  image_id                = "ami-880caa66"
  instance_type           = "c4.large"
  disable_api_termination = true
}

resource "outscale_vm_attributes" "outscale_vm_attributes" {
  instance_id             = "${outscale_vm.outscale_vm.0.id}"
  attribute               = "disableApiTermination"
  disable_api_termination = false
}

output "instance_id" {
  value = "${outscale_vm.outscale_vm.0.id}"
}
