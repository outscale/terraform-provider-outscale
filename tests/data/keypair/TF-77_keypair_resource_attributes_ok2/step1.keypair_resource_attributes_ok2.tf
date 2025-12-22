resource "outscale_keypair" "outscale_keypair" {
    keypair_name = "test-keypair-${random_string.suffix[0].result}"
    public_key   = file("keypair/TF-77_keypair_resource_attributes_ok2/keypair_public_test.pub")
}
