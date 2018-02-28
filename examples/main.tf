resource "outscale_vm" "basic" {
  image_id = "ami-8a6a0120"
  instance_type = "t2.micro"
  disable_api_termination = "false"
  key_name = "terraform-basic"
  security_group = ["sg-6ed31f3e"]
}

resource "outscale_vm" "basic2" {
  image_id = "ami-8a6a0120"
  instance_type = "t2.micro"
  disable_api_termination = "false"
  key_name = "terraform-basic"
  security_group = ["sg-6ed31f3e"]
}

# resource "outscale_public_ip" "bar" {
# 	count = 2
# }

# resource "outscale_public_ip_link" "by_allocation_id" {
# 	allocation_id = "${outscale_public_ip.bar.0.id}"
# 	instance_id = "${outscale_vm.basic.0.id}"
# 	depends_on = ["outscale_vm.basic"]
# }

data "outscale_vms" "basic_webs" {
	filter {
    name = "image-id"
    values = ["${outscale_vm.basic.image_id}"]
  }
}

# output "datasource_ip" {
#   value = "${data.outscale_vm.basic_web}"
# }
# #resource "outscale_public_ip_link" "to_eni" {
# #	allocation_id = "${outscale_public_ip.bar.0.id}"
# #	network_interface_id = "eni-f2a898a3"
# #}