resource "outscale_vm" "basic" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
}

data "outscale_vm" "basic_web" {
	filter {
    name = "instance-id"
    values = ["${outscale_vm.basic.id}"]
  }
}

output "datasource_ip" {
  value = "${data.outscale_vm.basic_web.*.instance_set.ip_address}"
}
