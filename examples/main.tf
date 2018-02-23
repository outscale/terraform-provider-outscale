resource "outscale_vm" "basic" {
	# count = 2
	image_id = "ami-8a6a0120"
	instance_type = "c4.large"
	key_name = "terraform-basic"
	subnet_id = "subnet-861fbecc"
}

# resource "outscale_public_ip" "bar" {
# 	count = 2
# }

# resource "outscale_public_ip_link" "by_allocation_id" {
# 	allocation_id = "${outscale_public_ip.bar.0.id}"
# 	instance_id = "${outscale_vm.basic.0.id}"
# 	depends_on = ["outscale_vm.basic"]
# }

# resource "outscale_public_ip_link" "by_public_ip" {
# 	public_ip = "${outscale_public_ip.bar.1.public_ip}"
# 	instance_id = "${outscale_vm.basic.1.id}"
#   depends_on = ["outscale_vm.basic"]
# }
# #resource "outscale_public_ip_link" "to_eni" {
# #	allocation_id = "${outscale_public_ip.bar.0.id}"
# #	network_interface_id = "eni-f2a898a3"
# #}