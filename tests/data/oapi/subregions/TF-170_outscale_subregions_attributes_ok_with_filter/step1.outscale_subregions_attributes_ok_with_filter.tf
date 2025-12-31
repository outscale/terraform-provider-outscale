data "outscale_subregions" "subregions-1" {
   filter {
        name     = "subregion_names"
        values   = ["${var.region}a"]
        }
}
