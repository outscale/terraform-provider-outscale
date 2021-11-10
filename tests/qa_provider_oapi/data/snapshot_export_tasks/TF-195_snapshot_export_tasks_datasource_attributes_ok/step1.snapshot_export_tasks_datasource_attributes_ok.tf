resource "outscale_volume" "outscale_volume_snap" {
    subregion_name   = "${var.region}a"
    size             = 5
}
resource "outscale_snapshot" "outscale_snapshot" {
    volume_id = outscale_volume.outscale_volume_snap.volume_id
}
resource "outscale_snapshot_export_task" "outscale_snapshot_export_task" {
    snapshot_id                     = outscale_snapshot.outscale_snapshot.snapshot_id
    osu_export {
         disk_image_format = "raw"
         osu_bucket        = var.osu_bucket_name 
         osu_prefix        = "prefix-195"
         osu_api_key {
           api_key_id      = var.access_key_id
           secret_key      = var.secret_key_id
          }
     }
  tags {
      key = "test"
      value = "test"
    }
}

resource "outscale_snapshot_export_task" "outscale_snapshot_export_task-2" {
    snapshot_id                     = outscale_snapshot.outscale_snapshot.snapshot_id
    osu_export {
         disk_image_format = "raw"
          osu_bucket        = "terraform-export-snap-3"
           osu_prefix        = "new-export-4"
          }
  tags {
      key = "test-2"
      value = "test-2"
    }
tags {
      key = "test-11"
      value = "test-11"
    }
}

data "outscale_snapshot_export_tasks" "data-2"  {
   filter {
       name = "task_ids"
       values = [outscale_snapshot_export_task.outscale_snapshot_export_task.task_id, outscale_snapshot_export_task.outscale_snapshot_export_task-2.task_id]
   }
}
