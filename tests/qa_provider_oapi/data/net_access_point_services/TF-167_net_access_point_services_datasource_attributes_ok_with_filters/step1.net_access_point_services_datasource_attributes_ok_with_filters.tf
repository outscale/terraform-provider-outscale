data "outscale_net_access_point_services" "data-1" {
   filter {
        name     = "service_names"
        values   = [var.service_name]
   }
}

data "outscale_net_access_point_services" "data-2" {
   filter {
        name     = "service_ids"
        values   = [data.outscale_net_access_point_services.data-1.services[0].service_id]
 }
}
