resource "outscale_net" "test" {
    ip_range  = "172.16.0.0/16"

    tags {
        key   = "Name"
        value = "terraform-testacc-data-source"
    }
}

resource "outscale_subnet" "test" {
    ip_range  = "172.16.0.0/24"
    net_id    = outscale_net.test.id
    tags {
        key   = "Name"
        value = "terraform-testacc-data-source"
    }
}

resource "outscale_route_table" "test" {
    net_id    = outscale_net.test.id
    tags {
        key   = "Name"
        value = "terraform-testacc-routetable-data-source"
    }
}
resource "outscale_route_table" "test2" {
    net_id    = outscale_net.test.id
}
resource "outscale_route_table_link" "a" {
    subnet_id = outscale_subnet.test.id
    route_table_id = outscale_route_table.test.id
}

data "outscale_route_table" "by_filter_1" {
    filter {
	name = "link_route_table_link_route_table_ids"
        values = [outscale_route_table_link.a.id]
    }
}

data "outscale_route_table" "by_filter_2" {
    filter {
        name = "link_route_table_ids"
        values = [outscale_route_table_link.a.route_table_id]
    }
depends_on = [outscale_route_table_link.a]
}  
