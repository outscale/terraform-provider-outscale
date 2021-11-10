data "outscale_product_types" "product_type_2" {
   filter {
        name     = "product_type_ids"
        values   = ["0001","0002"]
    }
}
