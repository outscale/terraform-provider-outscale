resource "outscale_keypair" "outscale_keypair" {
    keypair_name = "test-keypair-${random_string.suffix[0].result}"
}
