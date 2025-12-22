resource "outscale_keypair" "outscale_keypair" {
    keypair_name = "test-keypair-${random_string.suffix[0].result}"
}

data "outscale_keypair" "outscale_keypair" {
    filter {
        name   = "keypair_names"
        values = [outscale_keypair.outscale_keypair.keypair_name]
    }
depends_on = [outscale_keypair.outscale_keypair]
}
