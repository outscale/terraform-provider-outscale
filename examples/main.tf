resource "outscale_vm" "basic" {
  count = 1
  image_id = "ami-880caa66"
  instance_type = "t2.micro"
  disable_api_termination = "true"
  key_name = "integ_sut_keypair"
  security_group = ["sg-c73d3b6b"]
}

data "outscale_vm" "basic_web" {
	filter {
    name = "instance-id"
    values = ["${outscale_vm.basic.id}"]
  }
}

# output "datasource_ip" {
#   value = "${data.outscale_vm.basic_web}"
# }
