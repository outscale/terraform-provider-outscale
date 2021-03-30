resource "outscale_security_group" "outscale_security_group" {
    count = 1

    description         = "test group"
    security_group_name = "sg1-test-group_test-r"
    #net_id
}
