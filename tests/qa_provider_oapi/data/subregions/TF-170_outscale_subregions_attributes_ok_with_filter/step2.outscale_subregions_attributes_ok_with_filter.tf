data "outscale_subregions" "subregion-2" {
   filter {
        name     = "subregion_names"
        values   = ["${var.region}a", "${var.region}b"]
        }
}
