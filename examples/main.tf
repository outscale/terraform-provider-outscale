# resource "outscale_keypair" "outscale_keypair" {
#   key_name = "keyname_test_"
# }

data "outscale_keypair" "outscale_keypair" {
  key_name = "${outscale_keypair.outscale_keypair.key_name}"
}

resource "outscale_keypair" "outscale_keypair" {
  key_name = "keyname_test_"
}

resource "outscale_keypair" "outscale_keypair2" {
  key_name = "keyname_test_2"
}

data "outscale_keypairs" "outscale_keypairs" {
  key_name = [
    "${outscale_keypair.outscale_keypair.key_name}",
    "${outscale_keypair.outscale_keypair2.key_name}",
  ]
}
