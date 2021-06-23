resource "outscale_keypair" "outscale_keypair" {
    keypair_name = "keyname_test_"
}

resource "outscale_keypair" "outscale_keypair2" {
    keypair_name = "keyname_test_2"
}

data "outscale_keypairs" "outscale_keypairs" {
    filter {
         name = "keypair_names"
         values = [outscale_keypair.outscale_keypair.keypair_name, outscale_keypair.outscale_keypair2.keypair_name]
    }     
depends_on= [outscale_keypair.outscale_keypair,outscale_keypair.outscale_keypair2]
}
