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

data "outscale_vm" "basic_web" {
	filter {
    name = "instance-id"
    values = ["${outscale_vm.basic.id}"]
  }
}

data "outscale_vms" "basic_webs" {
	filter {
    name = "image-id"
    values = ["${outscale_vm.basic.image_id}"]
  }
}

# output "datasource_ip" {
#   value = "${data.outscale_vm.basic_web}"
# }
