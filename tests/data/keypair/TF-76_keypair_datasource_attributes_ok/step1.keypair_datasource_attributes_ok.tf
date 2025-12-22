resource "outscale_keypair" "outscale_keypair" {
    keypair_name = "keyname_test_ford"
}

data "outscale_keypair" "outscale_keypair" {
    filter {
        name   = "keypair_names"
        values = [outscale_keypair.outscale_keypair.keypair_name]
    }    
depends_on = [outscale_keypair.outscale_keypair]
}
