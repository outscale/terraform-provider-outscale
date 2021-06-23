data "outscale_vm_types" "vm_types_7" {
   filter {
        name     = "volume_sizes"
        values   = [80,40]
    }
  filter {
        name     = "volume_counts"
        values   = [2]
    }
  filter {
        name     = "memory_sizes"
        values   = [15]
    }
   filter {
        name     = "bsu_optimized"
        values   = [true]
    }
   filter {
        name     = "vcore_counts"
        values   = [4]
    }
}
