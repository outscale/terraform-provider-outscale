 	
data "outscale_vm_types" "vm_types_4" {
 filter {
        name     = "vm_type_names"
        values   = ["m3.large"]
    }
}
