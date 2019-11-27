resource "outscale_keypair" "outscale_keypair" {
    count = 1

    keypair_name = "keyname_test_"
}

output "keypair" {
    value =outscale_keypair.outscale_keypair.key_name
}
