resource "outscale_keypair" "outscale_keypair" {
    keypair_name = "keyname_test_import"
    public_key   = "${file("keypair_public_test.pub")}"
}