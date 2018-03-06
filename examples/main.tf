# resource "outscale_image" "outscale_image" {
#     name = "terraform testtest"
#     instance_id = "i-55031358"
#     no_reboot = "true"
# }

data "outscale_vms" "basic_web" {}
