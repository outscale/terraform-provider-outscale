resource "outscale_keypair" "outscale_keypair" {
    count = 1

    keypair_name = "keyname_test_"
}

resource "outscale_keypair" "outscale_keypair2" {
    count = 1

    keypair_name = "keyname_test_"
}
