resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF147"
}

resource "outscale_security_group" "sg_usd" {
    description         = "test vms"
    security_group_name = "test-sgusd"
}

## Test Private VM with user_data "private_only"  ##
resource "outscale_vm" "outscale_vm" {
    image_id             = var.image_id
    vm_type              = var.vm_type
    keypair_name         = outscale_keypair.my_keypair.keypair_name
    security_group_names = [outscale_security_group.sg_usd.security_group_name]
    user_data            = "LS0tLS1CRUdJTiBPVVRTQ0FMRSBTRUNUSU9OLS0tLS0KCnByaXZhdGVfb25seT10cnVlCgotLS0tLUVORCBPVVRTQ0FMRSBTRUNUSU9OLS0tLS0=" 
    tags {
       key = "name"
       value = "test-VM-private_only"   
    } 
 }

