resource "random_string" "suffix" {
    count = 4
    length  = 8
    special = false
    upper   = false
}

resource "random_integer" "bgp_asn" {
    count = 2
    min = 1
    max = 50620
}
