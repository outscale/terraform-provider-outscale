resource "outscale_security_group" "outscale_security_group" {
    description         = "test group"
    security_group_name = "test-sg-${random_string.suffix[0].result}"
}
