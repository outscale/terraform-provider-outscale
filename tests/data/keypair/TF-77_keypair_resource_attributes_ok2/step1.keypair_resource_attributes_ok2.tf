resource "outscale_keypair" "outscale_keypair" {
    keypair_name = "keyname_test_import"
    public_key   = file("data/keypair/TF-77_keypair_resource_attributes_ok2/keypair_public_test.pub")
}
