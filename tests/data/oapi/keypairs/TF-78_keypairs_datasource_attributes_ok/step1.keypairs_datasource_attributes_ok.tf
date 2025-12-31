resource "outscale_keypair" "outscale_keypair" {
    keypair_name = "test-keypair-${random_string.suffix[0].result}"
}

resource "outscale_keypair" "outscale_keypair2" {
    keypair_name = "test-keypair-${random_string.suffix[1].result}"
}

data "outscale_keypairs" "outscale_keypairs" {
    filter {
         name = "keypair_names"
         values = [outscale_keypair.outscale_keypair.keypair_name, outscale_keypair.outscale_keypair2.keypair_name]
    }
depends_on= [outscale_keypair.outscale_keypair,outscale_keypair.outscale_keypair2]
}
